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
	"github.com/SENERGY-Platform/process-deployment/lib/devices"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"
)

func TestDeviceOptions(t *testing.T) {

	mux := sync.Mutex{}
	calls := []string{}

	standardConnectorProtocolId := "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b"
	mqttProtocolId := "urn:infai:ses:protocol:c9a06d44-0cd0-465b-b0d9-560d604057a2"

	semanticmock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.Lock()
		defer mux.Unlock()
		calls = append(calls, r.URL.Path+"?"+r.URL.RawQuery)
		json.NewEncoder(w).Encode([]devicemodel.DeviceType{
			{Id: "dt1", Name: "dt1name", DeviceClass: devicemodel.DeviceClass{Id: "dc1"}, Services: []devicemodel.Service{
				testService("11", standardConnectorProtocolId, devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION),
				testService("11_b", mqttProtocolId, devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION),
				testService("12", standardConnectorProtocolId, devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
			}},
			{Id: "dt2", Name: "dt2name", DeviceClass: devicemodel.DeviceClass{Id: "dc1"}, Services: []devicemodel.Service{
				testService("21", standardConnectorProtocolId, devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
				testService("22", standardConnectorProtocolId, devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
			}},
			{Id: "dt3", Name: "dt1name", DeviceClass: devicemodel.DeviceClass{Id: "dc1"}, Services: []devicemodel.Service{
				testService("31", mqttProtocolId, devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION),
				testService("32", mqttProtocolId, devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
			}},
			{Id: "dt4", Name: "dt2name", DeviceClass: devicemodel.DeviceClass{Id: "dc1"}, Services: []devicemodel.Service{
				testService("41", mqttProtocolId, devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
				testService("42", mqttProtocolId, devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
			}},
		})
	}))

	defer semanticmock.Close()

	searchmock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwt/select/devices/device_type_id/dt1/x" {
			json.NewEncoder(w).Encode([]devices.PermSearchDevice{
				{Id: "1", Name: "1", DeviceType: "dt1"},
			})
		}
		if r.URL.Path == "/jwt/select/devices/device_type_id/dt2/x" {
			json.NewEncoder(w).Encode([]devices.PermSearchDevice{
				{Id: "2", Name: "2", DeviceType: "dt2"},
			})
		}
		if r.URL.Path == "/jwt/select/devices/device_type_id/dt3/x" {
			json.NewEncoder(w).Encode([]devices.PermSearchDevice{
				{Id: "3", Name: "3", DeviceType: "dt3"},
			})
		}
		if r.URL.Path == "/jwt/select/devices/device_type_id/dt4/x" {
			json.NewEncoder(w).Encode([]devices.PermSearchDevice{
				{Id: "4", Name: "4", DeviceType: "dt4"},
			})
		}
	}))

	defer searchmock.Close()

	protocolsEndpointMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e","name":"moses","handler":"moses","protocol_segments":[{"id":"urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a","name":"payload"}]},{"id":"urn:infai:ses:protocol:c9a06d44-0cd0-465b-b0d9-560d604057a2","name":"mqtt-connector","handler":"mqtt","interaction":"event","protocol_segments":[{"id":"urn:infai:ses:protocol-segment:ffaaf98e-7360-400c-94d4-7775683d38ca","name":"payload"},{"id":"urn:infai:ses:protocol-segment:e07a3534-d67f-4d35-bc14-352fdccf6d8d","name":"topic"}]},{"id":"urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b","name":"standard-connector","handler":"connector","protocol_segments":[{"id":"urn:infai:ses:protocol-segment:9956d8b5-46fa-4381-a227-c1df69808997","name":"metadata"},{"id":"urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65","name":"data"}]}]`))
	}))

	defer searchmock.Close()

	c := &config.ConfigStruct{
		SemanticRepoUrl: semanticmock.URL,
		PermSearchUrl:   searchmock.URL,
		DeviceRepoUrl:   protocolsEndpointMock.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	devices, err := devices.Factory.New(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}

	ctrl := &Ctrl{
		config:  c,
		devices: devices,
	}

	blockedProtocols, err := ctrl.GetBlockedProtocols()
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(blockedProtocols, []string{mqttProtocolId}) {
		t.Error(blockedProtocols)
		return
	}

	d, err := ctrl.GetOptions("token", deploymentmodel.DeviceDescriptions{{
		CharacteristicId: "chid1",
		Function:         devicemodel.Function{Id: devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION + "_1"},
		DeviceClass:      &devicemodel.DeviceClass{Id: "dc1"},
		Aspect:           &devicemodel.Aspect{Id: "a1"},
	}}.ToFilter(), blockedProtocols)

	if err != nil {
		t.Error(err)
		return
	}

	if len(d) != 1 || d[0].Device.Name != "1" || d[0].Device.Id != "1" || len(d[0].Services) != 1 || d[0].Services[0].Id != "11" {
		t.Error(d)
		return
	}

	mux.Lock()
	defer mux.Unlock()
	if !reflect.DeepEqual(calls, []string{
		"/device-types?filter=" + url.QueryEscape(`[{"function_id":"`+devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION+`_1","device_class_id":"dc1","aspect_id":"a1"}]`),
	}) {
		temp, _ := json.Marshal(calls)
		t.Error(string(temp))
	}
}

func testService(id string, protocolId string, functionType string) devicemodel.Service {
	return devicemodel.Service{
		Id:         id,
		LocalId:    id + "_l",
		Name:       id + "_name",
		Aspects:    []devicemodel.Aspect{{Id: "a1"}},
		ProtocolId: protocolId,
		Functions:  []devicemodel.Function{{Id: functionType + "_1", RdfType: functionType}},
	}
}
