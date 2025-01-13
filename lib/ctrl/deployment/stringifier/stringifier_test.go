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

package stringifier

import (
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/helper"
	"io/ioutil"
	"os"
	"runtime/debug"
	"testing"
)

const RESOURCE_BASE_DIR = "../../../tests/resources/"

func TestStringifyDeployment(t *testing.T) {
	infos, err := os.ReadDir(RESOURCE_BASE_DIR)
	if err != nil {
		t.Error(err)
		return
	}
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() && isValidaForStringifierTest(RESOURCE_BASE_DIR+name) {
			t.Run(name, func(t *testing.T) {
				testDeployment(t, name)
			})
		}
	}
}

func isValidaForStringifierTest(dir string) bool {
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
	return files["selected.json"] && files["deployed.bpmn"] && files["raw.bpmn"]
}

func testDeployment(t *testing.T, exampleName string) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Error(r, string(debug.Stack()))
		}
	}()
	conf, err := config.LoadConfig("../../../../config.json")
	if err != nil {
		t.Error(err)
		return
	}
	stringifier := New(conf, func(token auth.Token, aspectNodeId string) (aspectNode devicemodel.AspectNode, err error) {
		return devicemodel.AspectNode{Id: aspectNodeId}, nil
	})
	deploymentJson, err := ioutil.ReadFile(RESOURCE_BASE_DIR + exampleName + "/selected.json")
	if err != nil {
		t.Error(err)
		return
	}
	var deployment deploymentmodel.Deployment
	err = json.Unmarshal(deploymentJson, &deployment)
	if err != nil {
		t.Error(err)
		return
	}

	raw, err := ioutil.ReadFile(RESOURCE_BASE_DIR + exampleName + "/raw.bpmn")
	if err != nil {
		t.Error(err)
		return
	}
	deployment.Diagram.XmlRaw = string(raw)

	actual, err := stringifier.Deployment(deployment, "user1", auth.Token{})
	if err != nil {
		t.Error(err)
		return
	}

	expected, err := ioutil.ReadFile(RESOURCE_BASE_DIR + exampleName + "/deployed.bpmn")
	if err != nil {
		t.Error(err)
		return
	}

	equal, err := helper.XmlIsEqual(actual, string(expected))
	if err != nil {
		t.Error(err)
		return
	}
	if !equal {
		t.Error(actual, "\n\n", string(expected))
		return
	}
}
