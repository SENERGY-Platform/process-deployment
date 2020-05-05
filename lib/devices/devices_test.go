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

	_, err, _ = repo.GetFilteredDeviceTypes("token", []deploymentmodel.FilterCriteria{{
		CharacteristicId: strptr("chid1"),
		FunctionId:       strptr("fid"),
	}})

	if err != nil {
		t.Error(err)
		return
	}

	dt, err, _ := repo.GetFilteredDeviceTypes("token", []deploymentmodel.FilterCriteria{{
		CharacteristicId: strptr("chid1"),
		FunctionId:       strptr("fid"),
		DeviceClassId:    strptr("dc1"),
		AspectId:         strptr("a1"),
	}})

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

	_, err, _ = repo.GetFilteredDeviceTypes("token", []deploymentmodel.FilterCriteria{{
		CharacteristicId: strptr("chid1"),
		FunctionId:       strptr("fid"),
	}})

	if err != nil {
		t.Error(err)
		return
	}

	dt, err, _ := repo.GetFilteredDeviceTypes("token", []deploymentmodel.FilterCriteria{{
		CharacteristicId: strptr("chid1"),
		FunctionId:       strptr("fid"),
		DeviceClassId:    strptr("dc1"),
		AspectId:         strptr("a1"),
	}})

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

func strptr(in string) *string {
	return &in
}
