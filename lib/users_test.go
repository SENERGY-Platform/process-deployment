/*
 * Copyright 2021 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestUserDelete(t *testing.T) {
	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "uuid"
	}

	conf, err := config.LoadConfig("../config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	conf.Debug = true
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := &mocks.DatabaseMock{
		Deployments:  map[string]messages.DeploymentCommand{},
		Dependencies: map[string]dependencymodel.Dependencies{},
	}

	k := &mocks.KafkaMock{Produced: map[string][]string{}, Listeners: map[string][]func(msg []byte) error{}}

	control, err := ctrl.New(ctx, conf, k, db, mocks.Devices, mocks.ProcessModelRepo, mocks.Imports)
	if err != nil {
		fmt.Println(err)
		return
	}

	user1, err := auth.CreateToken("test", "user1")
	if err != nil {
		t.Error(err)
		return
	}
	user2, err := auth.CreateToken("test", "user2")
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("create device 1", testDeploy(control, user1, "1"))
	t.Run("create device 2", testDeploy(control, user1, "2"))
	t.Run("create device 3", testDeploy(control, user2, "3"))
	t.Run("create device 4", testDeploy(control, user2, "4"))

	time.Sleep(1 * time.Second)

	t.Run("check user1 before user1 delete", testCheckDeployments(db, user1, []string{"1", "2"}))
	t.Run("check user2 before user1 delete", testCheckDeployments(db, user2, []string{"3", "4"}))

	t.Run("delete user1", func(t *testing.T) {
		p, err := k.NewProducer(ctx, conf, conf.UsersTopic)
		if err != nil {
			t.Error(err)
			return
		}
		msg, err := json.Marshal(messages.UserCommandMsg{
			Command: "DELETE",
			Id:      "user1",
		})
		if err != nil {
			t.Error(err)
			return
		}
		err = p.Produce("user1", msg)
		if err != nil {
			t.Error(err)
			return
		}
	})

	time.Sleep(1 * time.Second)

	t.Run("check user1 after user1 delete", testCheckDeployments(db, user1, []string{}))
	t.Run("check user2 after user1 delete", testCheckDeployments(db, user2, []string{"3", "4"}))
}

func testCheckDeployments(db interfaces.Database, user auth.Token, expectedIds []string) func(t *testing.T) {
	return func(t *testing.T) {
		actualIds, err := db.GetDeploymentIds(user.GetUserId())
		if err != nil {
			t.Error(err)
			return
		}
		if expectedIds == nil {
			expectedIds = []string{}
		}
		if actualIds == nil {
			actualIds = []string{}
		}
		sort.Strings(expectedIds)
		sort.Strings(actualIds)
		if !reflect.DeepEqual(expectedIds, actualIds) {
			t.Error(actualIds, expectedIds)
			return
		}
	}
}

func testDeploy(control *ctrl.Ctrl, token auth.Token, id string) func(t *testing.T) {
	return func(t *testing.T) {
		old := config.NewId
		config.NewId = func() string {
			return id
		}
		defer func() {
			config.NewId = old
		}()
		process := CreateBlankProcess()
		_, err, _ := control.CreateDeployment(token, deploymentmodel.Deployment{
			Version: deploymentmodel.CurrentVersion,
			Id:      id,
			Name:    id + "_n",
			Diagram: deploymentmodel.Diagram{
				XmlRaw:      process,
				XmlDeployed: process,
				Svg:         "<svg/>",
			},
			Elements:   nil,
			Executable: true,
		}, "test", map[string]bool{})
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func CreateBlankProcess() string {
	templ := `<bpmn:definitions xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xmlns:bpmn='http://www.omg.org/spec/BPMN/20100524/MODEL' xmlns:bpmndi='http://www.omg.org/spec/BPMN/20100524/DI' xmlns:dc='http://www.omg.org/spec/DD/20100524/DC' id='Definitions_1' targetNamespace='http://bpmn.io/schema/bpmn'><bpmn:process id='PROCESSID' isExecutable='true'><bpmn:startEvent id='StartEvent_1'/></bpmn:process><bpmndi:BPMNDiagram id='BPMNDiagram_1'><bpmndi:BPMNPlane id='BPMNPlane_1' bpmnElement='PROCESSID'><bpmndi:BPMNShape id='_BPMNShape_StartEvent_2' bpmnElement='StartEvent_1'><dc:Bounds x='173' y='102' width='36' height='36'/></bpmndi:BPMNShape></bpmndi:BPMNPlane></bpmndi:BPMNDiagram></bpmn:definitions>`
	return strings.Replace(templ, "PROCESSID", "id_"+strconv.FormatInt(time.Now().Unix(), 10), 1)
}
