/*
 * Copyright 2025 InfAI (CC SES)
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

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/db"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDeploymentListApi(t *testing.T) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}
	conf.Debug = true
	conf.InitTopics = true

	port, _, err := docker.Mongo(ctx, &wg)
	if err != nil {
		t.Error(err)
		return
	}
	conf.MongoUrl = "mongodb://localhost:" + port

	freePort, err := GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}
	conf.ApiPort = strconv.Itoa(freePort)
	conf.Debug = true

	err = lib.Start(ctx, conf, mocks.Kafka, db.Factory, mocks.Devices, mocks.ProcessModelRepo, mocks.Imports)
	if err != nil {
		t.Error(err)
		return
	}

	ids := []string{}
	t.Run("create deployments", func(t *testing.T) {
		deploymentsByUser := map[string][]deploymentmodel.Deployment{
			"user1": {
				{
					Name:    "name1",
					Version: deploymentmodel.CurrentVersion,
					Diagram: deploymentmodel.Diagram{
						XmlRaw: CreateBlankProcess("1"),
						Svg:    CreateBlankSvg(),
					},
				},
				{
					Name:    "name2",
					Version: deploymentmodel.CurrentVersion,
					Diagram: deploymentmodel.Diagram{
						XmlRaw: CreateBlankProcess("2"),
						Svg:    CreateBlankSvg(),
					},
				},
			},
			"user2": {
				{
					Name:    "name3",
					Version: deploymentmodel.CurrentVersion,
					Diagram: deploymentmodel.Diagram{
						XmlRaw: CreateBlankProcess("3"),
						Svg:    CreateBlankSvg(),
					},
				},
			},
		}
		for user, deployments := range deploymentsByUser {
			token, err := auth.CreateToken("test", user)
			if err != nil {
				t.Error(err)
				return
			}
			for _, deployment := range deployments {
				buff := bytes.NewBuffer([]byte{})
				err = json.NewEncoder(buff).Encode(deployment)
				if err != nil {
					t.Error(err)
					return
				}
				req, err := http.NewRequest(http.MethodPost, "http://localhost:"+conf.ApiPort+"/v3/deployments", buff)
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token.Jwt())
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}
				if resp.StatusCode != http.StatusOK {
					temp, _ := io.ReadAll(resp.Body)
					t.Error(resp.StatusCode, string(temp))
					return
				}
				result := deploymentmodel.Deployment{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				if err != nil {
					t.Error(err)
					return
				}
				t.Logf("create %#v\n", result)
				ids = append(ids, result.Id)
			}
		}
	})

	time.Sleep(1 * time.Second)

	t.Run("list deployments default", func(t *testing.T) {
		token, err := auth.CreateToken("test", "user1")
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:"+conf.ApiPort+"/v3/deployments", nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token.Jwt())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			temp, _ := io.ReadAll(resp.Body)
			t.Error(resp.StatusCode, string(temp))
			return
		}
		var deployments []deploymentmodel.Deployment
		err = json.NewDecoder(resp.Body).Decode(&deployments)
		if err != nil {
			t.Error(err)
			return
		}
		if len(deployments) != 2 {
			t.Error(len(deployments))
			return
		}
	})

	t.Run("list deployments sort name desc", func(t *testing.T) {
		token, err := auth.CreateToken("test", "user1")
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:"+conf.ApiPort+"/v3/deployments?sort=name.desc", nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token.Jwt())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			temp, _ := io.ReadAll(resp.Body)
			t.Error(resp.StatusCode, string(temp))
			return
		}
		var deployments []deploymentmodel.Deployment
		err = json.NewDecoder(resp.Body).Decode(&deployments)
		if err != nil {
			t.Error(err)
			return
		}
		expected := []deploymentmodel.Deployment{
			{
				Id:      ids[1],
				Name:    "name2",
				Version: deploymentmodel.CurrentVersion,
				Diagram: deploymentmodel.Diagram{
					XmlRaw:      CreateBlankProcess("2"),
					XmlDeployed: CreateBlankProcessWithName("2", "name2"),
					Svg:         CreateBlankSvg(),
				},
			},
			{
				Id:      ids[0],
				Name:    "name1",
				Version: deploymentmodel.CurrentVersion,
				Diagram: deploymentmodel.Diagram{
					XmlRaw:      CreateBlankProcess("1"),
					XmlDeployed: CreateBlankProcessWithName("1", "name1"),
					Svg:         CreateBlankSvg(),
				},
			},
		}
		if !reflect.DeepEqual(deployments, expected) {
			t.Errorf("\na:%#v\ne:%#v\n", deployments, expected)
			return
		}
	})

	t.Run("list deployments sort name asc", func(t *testing.T) {
		token, err := auth.CreateToken("test", "user1")
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:"+conf.ApiPort+"/v3/deployments?sort=name", nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token.Jwt())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			temp, _ := io.ReadAll(resp.Body)
			t.Error(resp.StatusCode, string(temp))
			return
		}
		var deployments []deploymentmodel.Deployment
		err = json.NewDecoder(resp.Body).Decode(&deployments)
		if err != nil {
			t.Error(err)
			return
		}
		expected := []deploymentmodel.Deployment{
			{
				Id:      ids[0],
				Name:    "name1",
				Version: deploymentmodel.CurrentVersion,
				Diagram: deploymentmodel.Diagram{
					XmlRaw:      CreateBlankProcess("1"),
					XmlDeployed: CreateBlankProcessWithName("1", "name1"),
					Svg:         CreateBlankSvg(),
				},
			},
			{
				Id:      ids[1],
				Name:    "name2",
				Version: deploymentmodel.CurrentVersion,
				Diagram: deploymentmodel.Diagram{
					XmlRaw:      CreateBlankProcess("2"),
					XmlDeployed: CreateBlankProcessWithName("2", "name2"),
					Svg:         CreateBlankSvg(),
				},
			},
		}
		if !reflect.DeepEqual(deployments, expected) {
			t.Error(deployments)
			return
		}
	})

	t.Run("list deployments limit", func(t *testing.T) {
		token, err := auth.CreateToken("test", "user1")
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:"+conf.ApiPort+"/v3/deployments?sort=name&limit=1", nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token.Jwt())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			temp, _ := io.ReadAll(resp.Body)
			t.Error(resp.StatusCode, string(temp))
			return
		}
		var deployments []deploymentmodel.Deployment
		err = json.NewDecoder(resp.Body).Decode(&deployments)
		if err != nil {
			t.Error(err)
			return
		}
		expected := []deploymentmodel.Deployment{
			{
				Id:      ids[0],
				Name:    "name1",
				Version: deploymentmodel.CurrentVersion,
				Diagram: deploymentmodel.Diagram{
					XmlRaw:      CreateBlankProcess("1"),
					XmlDeployed: CreateBlankProcessWithName("1", "name1"),
					Svg:         CreateBlankSvg(),
				},
			},
		}
		if !reflect.DeepEqual(deployments, expected) {
			t.Error(deployments)
			return
		}
	})

	t.Run("list deployments offset", func(t *testing.T) {
		token, err := auth.CreateToken("test", "user1")
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:"+conf.ApiPort+"/v3/deployments?sort=name&offset=1", nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", token.Jwt())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			temp, _ := io.ReadAll(resp.Body)
			t.Error(resp.StatusCode, string(temp))
			return
		}
		var deployments []deploymentmodel.Deployment
		err = json.NewDecoder(resp.Body).Decode(&deployments)
		if err != nil {
			t.Error(err)
			return
		}
		expected := []deploymentmodel.Deployment{
			{
				Id:      ids[1],
				Name:    "name2",
				Version: deploymentmodel.CurrentVersion,
				Diagram: deploymentmodel.Diagram{
					XmlRaw:      CreateBlankProcess("2"),
					XmlDeployed: CreateBlankProcessWithName("2", "name2"),
					Svg:         CreateBlankSvg(),
				},
			},
		}
		if !reflect.DeepEqual(deployments, expected) {
			t.Errorf("\na:%#v\ne:%#v\n", deployments, expected)
			return
		}
	})
}

func CreateBlankSvg() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.2" id="Layer_1" x="0px" y="0px" viewBox="0 0 20 16" xml:space="preserve">
<path fill="#D61F33" d="M10,0L0,16h20L10,0z M11,13.908H9v-2h2V13.908z M9,10.908v-6h2v6H9z"/>
</svg>`
}

func CreateBlankProcess(id string) string {
	templ := `<bpmn:definitions xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn"><bpmn:process id="PROCESSID" isExecutable="true"><bpmn:startEvent id="StartEvent_1"/></bpmn:process><bpmndi:BPMNDiagram id="BPMNDiagram_1"><bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="PROCESSID"><bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1"><dc:Bounds x="173" y="102" width="36" height="36"/></bpmndi:BPMNShape></bpmndi:BPMNPlane></bpmndi:BPMNDiagram></bpmn:definitions>`
	return strings.Replace(templ, "PROCESSID", "id_"+id, 1)
}

func CreateBlankProcessWithName(id string, name string) string {
	templ := `<bpmn:definitions xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn"><bpmn:process id="PROCESSID" isExecutable="true" name="PROCESSNAME"><bpmn:startEvent id="StartEvent_1"/></bpmn:process><bpmndi:BPMNDiagram id="BPMNDiagram_1"><bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="PROCESSID"><bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1"><dc:Bounds x="173" y="102" width="36" height="36"/></bpmndi:BPMNShape></bpmndi:BPMNPlane></bpmndi:BPMNDiagram></bpmn:definitions>`
	result := strings.Replace(templ, "PROCESSID", "id_"+id, 1)
	result = strings.Replace(result, "PROCESSNAME", name, 1)
	return result
}
