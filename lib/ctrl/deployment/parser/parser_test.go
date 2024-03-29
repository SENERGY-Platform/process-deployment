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

package parser

import (
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"testing"
)

const RESOURCE_BASE_DIR = "../../../tests/resources/"

func TestParseDeployment(t *testing.T) {
	infos, err := os.ReadDir(RESOURCE_BASE_DIR)
	if err != nil {
		t.Error(err)
		return
	}
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() && isValidaForParserTest(RESOURCE_BASE_DIR+name) {
			t.Run(name, func(t *testing.T) {
				testPrepareDeployment(t, name)
			})
		}
	}
}

func TestStartParameter(t *testing.T) {
	infos, err := os.ReadDir(RESOURCE_BASE_DIR)
	if err != nil {
		t.Error(err)
		return
	}
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() && isValidaForStartParameterTest(RESOURCE_BASE_DIR+name) {
			t.Run(name, func(t *testing.T) {
				testStartParameter(t, name)
			})
		}
	}
}

func isValidaForStartParameterTest(dir string) bool {
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
	return files["raw.bpmn"] && files["expected_start_parameters.json"]
}

func isValidaForParserTest(dir string) bool {
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
	return files["raw.bpmn"] && files["parsed.json"]
}

func testStartParameter(t *testing.T, exampleName string) {
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
	parser := New(conf)
	rawBpmnFile, err := os.ReadFile(RESOURCE_BASE_DIR + exampleName + "/raw.bpmn")
	if err != nil {
		t.Error(err)
		return
	}
	result, err := parser.EstimateStartParameter(string(rawBpmnFile))
	if err != nil {
		t.Error(err)
		return
	}

	expectedStr, err := os.ReadFile(RESOURCE_BASE_DIR + exampleName + "/expected_start_parameters.json")
	if err != nil {
		t.Error(err)
		return
	}
	var expected []deploymentmodel.ProcessStartParameter
	err = json.Unmarshal(expectedStr, &expected)
	if err != nil {
		t.Error(err)
		return
	}

	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Id < expected[j].Id
	})

	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})

	if !reflect.DeepEqual(result, expected) {
		resultJson, _ := json.Marshal(result)
		expectedJson, _ := json.Marshal(expected)
		t.Error(string(resultJson), "\n", string(expectedJson))
		return
	}
}

func testPrepareDeployment(t *testing.T, exampleName string) {
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
	parser := New(conf)
	rawBpmnFile, err := os.ReadFile(RESOURCE_BASE_DIR + exampleName + "/raw.bpmn")
	if err != nil {
		t.Error(err)
		return
	}
	deployment, err := parser.PrepareDeployment(string(rawBpmnFile))
	if err != nil {
		t.Error(err)
		return
	}
	if deployment.Diagram.XmlRaw != string(rawBpmnFile) {
		t.Error(deployment.Diagram.XmlRaw)
		return
	}
	deployment.Diagram = deploymentmodel.Diagram{}
	expectedDeployment, err := os.ReadFile(RESOURCE_BASE_DIR + exampleName + "/parsed.json")
	if err != nil {
		t.Error(err)
		return
	}
	var expected deploymentmodel.Deployment
	err = json.Unmarshal(expectedDeployment, &expected)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(deployment, expected) {
		deploymentJson, _ := json.Marshal(deployment)
		expectedJson, _ := json.Marshal(expected)
		t.Error(string(deploymentJson), "\n", string(expectedJson))
		return
	}
}
