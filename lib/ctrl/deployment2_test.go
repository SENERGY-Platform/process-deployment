/*
 * Copyright 2020 InfAI (CC SES)
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

package ctrl

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/db"
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"strings"
	"time"

	"reflect"
	"runtime/debug"
	"sync"
	"testing"
)

const RESOURCE_BASE_DIR = "../tests/resources/"

func TestDeploymentHandler(t *testing.T) {
	infos, err := os.ReadDir(RESOURCE_BASE_DIR)
	if err != nil {
		t.Error(err)
		return
	}
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() && isValidaForDeploymentHandlerTest(RESOURCE_BASE_DIR+name) {
			t.Run(name, func(t *testing.T) {
				testDeploymentHandler(t, name)
			})
		}
	}
}

func TestSetExecutableFlag(t *testing.T) {
	t.Run("true", testSetExecutableFlag(deploymentmodel.Deployment{}, true))
	t.Run("task true", testSetExecutableFlag(deploymentmodel.Deployment{
		Elements: []deploymentmodel.Element{
			{
				Task: &deploymentmodel.Task{
					Selection: deploymentmodel.Selection{SelectionOptions: []deploymentmodel.SelectionOption{deploymentmodel.SelectionOption{}}},
				},
			},
		},
	}, true))
	t.Run("task false", testSetExecutableFlag(deploymentmodel.Deployment{
		Elements: []deploymentmodel.Element{
			{
				Task: &deploymentmodel.Task{
					Selection: deploymentmodel.Selection{},
				},
			},
		},
	}, false))

	t.Run("event false", testSetExecutableFlag(deploymentmodel.Deployment{
		Elements: []deploymentmodel.Element{
			{
				MessageEvent: &deploymentmodel.MessageEvent{
					Selection: deploymentmodel.Selection{},
				},
			},
		},
	}, false))

	t.Run("event true", testSetExecutableFlag(deploymentmodel.Deployment{
		Elements: []deploymentmodel.Element{
			{
				MessageEvent: &deploymentmodel.MessageEvent{
					Selection: deploymentmodel.Selection{SelectionOptions: []deploymentmodel.SelectionOption{deploymentmodel.SelectionOption{}}},
				},
			},
		},
	}, true))

	t.Run("event task true", testSetExecutableFlag(deploymentmodel.Deployment{
		Elements: []deploymentmodel.Element{
			{
				MessageEvent: &deploymentmodel.MessageEvent{
					Selection: deploymentmodel.Selection{SelectionOptions: []deploymentmodel.SelectionOption{deploymentmodel.SelectionOption{}}},
				},
			},
			{
				Task: &deploymentmodel.Task{
					Selection: deploymentmodel.Selection{SelectionOptions: []deploymentmodel.SelectionOption{deploymentmodel.SelectionOption{}}},
				},
			},
		},
	}, true))

	t.Run("event task false 1", testSetExecutableFlag(deploymentmodel.Deployment{
		Elements: []deploymentmodel.Element{
			{
				MessageEvent: &deploymentmodel.MessageEvent{
					Selection: deploymentmodel.Selection{},
				},
			},
			{
				Task: &deploymentmodel.Task{
					Selection: deploymentmodel.Selection{SelectionOptions: []deploymentmodel.SelectionOption{deploymentmodel.SelectionOption{}}},
				},
			},
		},
	}, false))

	t.Run("event task false 2", testSetExecutableFlag(deploymentmodel.Deployment{
		Elements: []deploymentmodel.Element{
			{
				MessageEvent: &deploymentmodel.MessageEvent{
					Selection: deploymentmodel.Selection{SelectionOptions: []deploymentmodel.SelectionOption{deploymentmodel.SelectionOption{}}},
				},
			},
			{
				Task: &deploymentmodel.Task{
					Selection: deploymentmodel.Selection{},
				},
			},
		},
	}, false))
}

func testSetExecutableFlag(deployment deploymentmodel.Deployment, expected bool) func(t *testing.T) {
	return func(t *testing.T) {
		(&Ctrl{}).SetExecutableFlag(&deployment)
		if deployment.Executable != expected {
			t.Error(deployment.Executable, expected)
		}
	}
}

func isValidaForDeploymentHandlerTest(dir string) bool {
	infos, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	files := map[string]bool{}
	for _, info := range infos {
		if !info.IsDir() {
			files[info.Name()] = true
		}
	}
	return files["selected.json"] && files["dependencies.json"]
}

func testDeploymentHandler(t *testing.T, exampleName string) {
	deploymentId := "deployment-id"
	userId := "user1"
	defer func() {
		if r := recover(); r != nil {
			t.Error(r, string(debug.Stack()))
		}
	}()

	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}

	expectedDependenciesJson, err := os.ReadFile(RESOURCE_BASE_DIR + exampleName + "/dependencies.json")
	if err != nil {
		t.Error(err)
		return
	}
	expectedDependencies := dependencymodel.Dependencies{}
	err = json.Unmarshal(expectedDependenciesJson, &expectedDependencies)
	if err != nil {
		t.Error(err)
		return
	}
	expectedDependencies.DeploymentId = deploymentId
	expectedDependencies.Owner = userId

	deploymentJson, err := os.ReadFile(RESOURCE_BASE_DIR + exampleName + "/selected.json")
	if err != nil {
		t.Error(err)
		return
	}

	deployment := deploymentmodel.Deployment{}
	err = json.Unmarshal(deploymentJson, &deployment)
	if err != nil {
		t.Error(err)
		return
	}
	deployment.Id = deploymentId

	mongoPort, _, err := docker.MongoContainer(ctx, &wg)
	if err != nil {
		t.Error(err)
		return
	}
	conf.MongoUrl = "mongodb://localhost:" + mongoPort

	devicesMock := &mocks.DeviceRepoMock{}

	devicesJson, err := os.ReadFile(RESOURCE_BASE_DIR + exampleName + "/devices.json")
	if err == nil {
		devicesList := []devicemodel.Device{}
		json.Unmarshal(devicesJson, &devicesList)
		for _, d := range devicesList {
			devicesMock.SetDevice(d.Id, d)
		}
	}

	db, err := db.Factory.New(ctx, conf)
	if err != nil {
		t.Error(err)
		return
	}

	ctrl := &Ctrl{
		config:  conf,
		db:      db,
		devices: devicesMock,
	}

	err = ctrl.HandleDeployment(messages.DeploymentCommand{
		Command:    "PUT",
		Id:         deploymentId,
		Owner:      userId,
		Deployment: &deployment,
		Source:     "test",
	})
	if err != nil {
		t.Error(err)
		return
	}

	usertoken, err := ForgeUserToken(userId)
	if err != nil {
		t.Error(err)
		return
	}

	dependencies, err, _ := ctrl.GetDependencies(usertoken, deployment.Id)
	if err != nil {
		t.Error(err)
		return
	}

	checkDependencies(t, dependencies, expectedDependencies)
}

func checkDependencies(t *testing.T, actual dependencymodel.Dependencies, expected dependencymodel.Dependencies) {
	if !reflect.DeepEqual(actual, expected) {
		t.Error("\n", actual, "\n\n", expected)
	}
}

func ForgeUserToken(user string) (token auth.Token, err error) {
	// Create the Claims
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		Issuer:    "test",
		Subject:   user,
	}

	jwtoken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	unsignedTokenString, err := jwtoken.SigningString()
	if err != nil {
		log.Println("ERROR: GetUserToken::SigningString()", err)
		return token, err
	}
	tokenString := strings.Join([]string{unsignedTokenString, ""}, ".")
	token.Token = "Bearer " + tokenString
	token.Sub = user
	return token, err
}
