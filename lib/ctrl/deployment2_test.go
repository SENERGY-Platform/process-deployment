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
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/deployment/parser"
	"github.com/SENERGY-Platform/process-deployment/lib/devices"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime/debug"
	"testing"
)

const RESOURCE_BASE_DIR = "../tests/resources/"

func TestPrepareDeploymentV2(t *testing.T) {
	infos, err := ioutil.ReadDir(RESOURCE_BASE_DIR)
	if err != nil {
		t.Error(err)
		return
	}
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() && isValidaForPrepareTest(RESOURCE_BASE_DIR+name) {
			t.Run(name, func(t *testing.T) {
				testPrepareDeployment(t, name)
			})
		}
	}
}

func isValidaForPrepareTest(dir string) bool {
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
	return files["raw.bpmn"] && files["prepared.json"] && files["devices.json"]
}

func testPrepareDeployment(t *testing.T, exampleName string) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Error(r, string(debug.Stack()))
		}
	}()
	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}

	devicesDesc, err := ioutil.ReadFile(RESOURCE_BASE_DIR + exampleName + "/devices.json")
	if err != nil {
		t.Error(err)
		return
	}

	deviceTypesDesc, err := ioutil.ReadFile(RESOURCE_BASE_DIR + "iot/devicetypes.json")
	if err != nil {
		t.Error(err)
		return
	}

	protocolsDesc, err := ioutil.ReadFile(RESOURCE_BASE_DIR + "iot/protocols.json")
	if err != nil {
		t.Error(err)
		return
	}

	rawBpmnFile, err := ioutil.ReadFile(RESOURCE_BASE_DIR + exampleName + "/raw.bpmn")
	if err != nil {
		t.Error(err)
		return
	}

	closeTestApis := func() {}
	conf.DeviceRepoUrl, conf.SemanticRepoUrl, conf.PermSearchUrl, closeTestApis, err = createTestIotApis(protocolsDesc, deviceTypesDesc, devicesDesc)
	if err != nil {
		t.Error(err)
		return
	}
	defer closeTestApis()

	devices, err := devices.Factory.New(context.Background(), conf)
	if err != nil {
		t.Error(err)
		return
	}

	ctrl := &Ctrl{
		config:           conf,
		devices:          devices,
		deploymentParser: parser.New(conf),
	}

	deployment, err, _ := ctrl.PrepareDeploymentV2("token", string(rawBpmnFile), "")
	if err != nil {
		t.Error(err)
		return
	}
	if deployment.Diagram.XmlRaw != string(rawBpmnFile) {
		t.Error(deployment.Diagram.XmlRaw)
		return
	}
	deployment.Diagram = deploymentmodel.Diagram{}
	expectedDeployment, err := ioutil.ReadFile(RESOURCE_BASE_DIR + exampleName + "/prepared.json")
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

func createTestIotApis(protocolsDesc []byte, deviceTypesDesc []byte, devicesDesc []byte) (deviceRepoUrl string, semanticRepoUrl string, searchUrl string, close func(), err error) {
	protocols := []devicemodel.Protocol{}
	err = json.Unmarshal(protocolsDesc, &protocols)
	if err != nil {
		return
	}

	deviceTypes := []devicemodel.DeviceType{}
	err = json.Unmarshal(deviceTypesDesc, &deviceTypes)
	if err != nil {
		return
	}

	testDevices := []devices.PermSearchDevice{}
	err = json.Unmarshal(devicesDesc, &testDevices)
	if err != nil {
		return
	}

	//search fro device-type by function in semantic-repo
	semanticmock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filter := devicemodel.DeviceTypesFilter{}
		err := json.Unmarshal([]byte(r.URL.Query().Get("filter")), &filter)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		result := []devicemodel.DeviceType{}
		for _, deviceType := range deviceTypes {
			if testDeviceTypeMatchesFilterList(deviceType, filter) {
				result = append(result, deviceType)
			}
		}
		json.NewEncoder(w).Encode(result)
	}))

	semanticRepoUrl = semanticmock.URL

	//search for device by device type from perm-search
	searchmock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := []devices.PermSearchDevice{}
		for _, device := range testDevices {
			if r.URL.Path == "/jwt/select/devices/device_type_id/"+device.DeviceType+"/x" {
				result = append(result, device)
			}
		}
		json.NewEncoder(w).Encode(result)
	}))

	searchUrl = searchmock.URL

	//list protocols in device-repo
	protocolsEndpointMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(protocols)
	}))

	deviceRepoUrl = protocolsEndpointMock.URL
	close = func() {
		semanticmock.Close()
		searchmock.Close()
		protocolsEndpointMock.Close()
	}
	return
}

func testDeviceTypeMatchesFilterList(deviceType devicemodel.DeviceType, filters devicemodel.DeviceTypesFilter) bool {
	for _, filter := range filters {
		if !testDeviceTypeMatchesFilter(deviceType, filter) {
			return false
		}
	}
	return true
}

func testDeviceTypeMatchesFilter(deviceType devicemodel.DeviceType, filter devicemodel.DeviceTypeFilterElement) bool {
	if filter.DeviceClassId != "" && filter.DeviceClassId != deviceType.DeviceClass.Id {
		return false
	}
	for _, service := range deviceType.Services {
		if testServiceMatchesFilter(service, filter) {
			return true
		}
	}
	return false
}

func testServiceMatchesFilter(service devicemodel.Service, filter devicemodel.DeviceTypeFilterElement) bool {
	implementsFunction := false
	matchesAspect := false
	for _, function := range service.Functions {
		if function.Id == filter.FunctionId {
			implementsFunction = true
			break
		}
	}
	for _, aspect := range service.Aspects {
		if aspect.Id == filter.AspectId {
			matchesAspect = true
			break
		}
	}
	return (filter.AspectId == "" || matchesAspect) && (filter.FunctionId == "" || implementsFunction)
}
