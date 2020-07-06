/*
 * Copyright 2019 InfAI (CC SES)
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

package devices

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"
)

func TestCaching(t *testing.T) {
	mux := sync.Mutex{}
	calls := []string{}

	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.Lock()
		defer mux.Unlock()
		calls = append(calls, r.URL.Path)
		json.NewEncoder(w).Encode(devicemodel.Service{Id: "s1", Name: "s1name"})
	}))

	defer mock.Close()

	c := &config.ConfigStruct{
		DeviceRepoUrl: mock.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repo, err := Factory.New(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}

	service, err, _ := repo.GetService("foobar", "s1")
	if err != nil {
		t.Error(err)
		return
	}

	if service.Name != "s1name" || service.Id != "s1" {
		t.Error(service)
		return
	}

	service, err, _ = repo.GetService("foobar", "s1")
	if err != nil {
		t.Error(err)
		return
	}

	if service.Name != "s1name" || service.Id != "s1" {
		t.Error(service)
		return
	}

	service, err, _ = repo.GetService("foobar", "s1")
	if err != nil {
		t.Error(err)
		return
	}

	if service.Name != "s1name" || service.Id != "s1" {
		t.Error(service)
		return
	}

	mux.Lock()
	defer mux.Unlock()
	if !reflect.DeepEqual(calls, []string{"/services/s1"}) {
		temp, _ := json.Marshal(calls)
		t.Error(string(temp))
	}

}

func TestGetFilteredDeviceTypes(t *testing.T) {

	mux := sync.Mutex{}
	calls := []string{}

	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.Lock()
		defer mux.Unlock()
		calls = append(calls, r.URL.Path+"?"+r.URL.RawQuery)
		json.NewEncoder(w).Encode([]devicemodel.DeviceType{{Id: "dt1", Name: "dt1name"}})
	}))

	defer mock.Close()

	c := &config.ConfigStruct{
		SemanticRepoUrl: mock.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	temp, err := Factory.New(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}

	repo := temp.(*Repository)

	_, err, _ = repo.GetFilteredDeviceTypes("token", deploymentmodel.DeviceDescriptions{{
		CharacteristicId: "chid1",
		Function:         devicemodel.Function{Id: "fid"},
		DeviceClass:      nil,
		Aspect:           nil,
	}}.ToFilter())

	if err != nil {
		t.Error(err)
		return
	}

	dt, err, _ := repo.GetFilteredDeviceTypes("token", deploymentmodel.DeviceDescriptions{{
		CharacteristicId: "chid1",
		Function:         devicemodel.Function{Id: "fid"},
		DeviceClass:      &devicemodel.DeviceClass{Id: "dc1"},
		Aspect:           &devicemodel.Aspect{Id: "a1"},
	}}.ToFilter())

	if err != nil {
		t.Error(err)
		return
	}

	if len(dt) != 1 || dt[0].Name != "dt1name" || dt[0].Id != "dt1" {
		t.Error(dt)
		return
	}

	mux.Lock()
	defer mux.Unlock()
	if !reflect.DeepEqual(calls, []string{
		"/device-types?filter=" + url.QueryEscape(`[{"function_id":"fid","device_class_id":"","aspect_id":""}]`),
		"/device-types?filter=" + url.QueryEscape(`[{"function_id":"fid","device_class_id":"dc1","aspect_id":"a1"}]`),
	}) {
		temp, _ := json.Marshal(calls)
		t.Error(string(temp))
	}
}

func TestGetFilteredDevices(t *testing.T) {

	mux := sync.Mutex{}
	calls := []string{}

	semanticmock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.Lock()
		defer mux.Unlock()
		calls = append(calls, r.URL.Path+"?"+r.URL.RawQuery)
		json.NewEncoder(w).Encode([]devicemodel.DeviceType{
			{Id: "dt1", Name: "dt1name", DeviceClass: devicemodel.DeviceClass{Id: "dc1"}, Services: []devicemodel.Service{
				testService("11", "pid", devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION),
				testService("12", "pid", devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
			}},
			{Id: "dt2", Name: "dt2name", DeviceClass: devicemodel.DeviceClass{Id: "dc1"}, Services: []devicemodel.Service{
				testService("21", "pid", devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
				testService("22", "pid", devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
			}},
			{Id: "dt3", Name: "dt1name", DeviceClass: devicemodel.DeviceClass{Id: "dc1"}, Services: []devicemodel.Service{
				testService("31", "mqtt", devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION),
				testService("32", "mqtt", devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
			}},
			{Id: "dt4", Name: "dt2name", DeviceClass: devicemodel.DeviceClass{Id: "dc1"}, Services: []devicemodel.Service{
				testService("41", "mqtt", devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
				testService("42", "mqtt", devicemodel.SES_ONTOLOGY_CONTROLLING_FUNCTION),
			}},
		})
	}))

	defer semanticmock.Close()

	searchmock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwt/select/devices/device_type_id/dt1/x" {
			json.NewEncoder(w).Encode([]PermSearchDevice{
				{Id: "1", Name: "1", DeviceType: "dt1"},
			})
		}
		if r.URL.Path == "/jwt/select/devices/device_type_id/dt2/x" {
			json.NewEncoder(w).Encode([]PermSearchDevice{
				{Id: "2", Name: "2", DeviceType: "dt2"},
			})
		}
		if r.URL.Path == "/jwt/select/devices/device_type_id/dt3/x" {
			json.NewEncoder(w).Encode([]PermSearchDevice{
				{Id: "3", Name: "3", DeviceType: "dt3"},
			})
		}
		if r.URL.Path == "/jwt/select/devices/device_type_id/dt4/x" {
			json.NewEncoder(w).Encode([]PermSearchDevice{
				{Id: "4", Name: "4", DeviceType: "dt4"},
			})
		}
	}))

	defer searchmock.Close()

	c := &config.ConfigStruct{
		SemanticRepoUrl:     semanticmock.URL,
		PermSearchUrl:       searchmock.URL,
		EventBasedProtocols: []string{"mqtt"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	temp, err := Factory.New(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}

	repo := temp.(*Repository)

	d, err, _ := repo.GetFilteredDevices("token", deploymentmodel.DeviceDescriptions{{
		CharacteristicId: "chid1",
		Function:         devicemodel.Function{Id: devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION + "_1"},
		DeviceClass:      &devicemodel.DeviceClass{Id: "dc1"},
		Aspect:           &devicemodel.Aspect{Id: "a1"},
	}}.ToFilter())

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
