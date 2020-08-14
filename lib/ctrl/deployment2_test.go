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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/db"
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"io/ioutil"
	"reflect"
	"runtime/debug"
	"sync"
	"testing"
)

const RESOURCE_BASE_DIR = "../tests/resources/"

func TestDeploymentHandlerV2(t *testing.T) {
	infos, err := ioutil.ReadDir(RESOURCE_BASE_DIR)
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

func TestSetExecutableFlagV2(t *testing.T) {
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
		(&Ctrl{}).SetExecutableFlagV2(&deployment)
		if deployment.Executable != expected {
			t.Error(deployment.Executable, expected)
		}
	}
}

func isValidaForDeploymentHandlerTest(dir string) bool {
	infos, err := ioutil.ReadDir(dir)
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

	expectedDependenciesJson, err := ioutil.ReadFile(RESOURCE_BASE_DIR + exampleName + "/dependencies.json")
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

	deploymentJson, err := ioutil.ReadFile(RESOURCE_BASE_DIR + exampleName + "/selected.json")
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

	db, err := db.Factory.New(ctx, conf)
	if err != nil {
		t.Error(err)
		return
	}

	ctrl := &Ctrl{
		config: conf,
		db:     db,
	}

	err = ctrl.HandleDeployment(messages.DeploymentCommand{
		Command:      "PUT",
		Id:           deploymentId,
		Owner:        userId,
		Deployment:   nil,
		DeploymentV2: &deployment,
		Source:       "test",
	})
	if err != nil {
		t.Error(err)
		return
	}

	dependencies, err, _ := ctrl.GetDependencies(jwt_http_router.Jwt{UserId: userId}, deployment.Id)
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
