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

package mocks

import (
	"context"
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"sync"
)

type DeviceRepoMock struct {
	mux       sync.Mutex
	protocols map[string]devicemodel.Protocol
	devices   map[string]devicemodel.Device
	services  map[string]devicemodel.Service
	options   []devicemodel.Selectable
}

var Devices = &DeviceRepoMock{protocols: map[string]devicemodel.Protocol{}, devices: map[string]devicemodel.Device{}, services: map[string]devicemodel.Service{}}

func (this *DeviceRepoMock) New(ctx context.Context, config config.Config) (interfaces.Devices, error) {
	return this, nil
}

func (this *DeviceRepoMock) GetProtocol(id string) (devicemodel.Protocol, error, int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if result, ok := this.protocols[id]; ok {
		return result, nil, 200
	} else {
		return result, errors.New("protocol " + id + " not found"), http.StatusNotFound
	}
}

func (this *DeviceRepoMock) SetProtocol(id string, protocol devicemodel.Protocol) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.protocols[id] = protocol
}

func (this *DeviceRepoMock) GetDevice(token jwt_http_router.JwtImpersonate, id string) (devicemodel.Device, error, int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if result, ok := this.devices[id]; ok {
		return result, nil, 200
	} else {
		return result, errors.New("device " + id + " not found"), http.StatusNotFound
	}
}

func (this *DeviceRepoMock) SetDevice(id string, device devicemodel.Device) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.devices[id] = device
}

func (this *DeviceRepoMock) GetService(token jwt_http_router.JwtImpersonate, id string) (devicemodel.Service, error, int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if result, ok := this.services[id]; ok {
		return result, nil, 200
	} else {
		return result, errors.New("service " + id + " not found"), http.StatusNotFound
	}
}

func (this *DeviceRepoMock) SetService(id string, service devicemodel.Service) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.services[id] = service
}

func (this *DeviceRepoMock) GetFilteredDevices(token jwt_http_router.JwtImpersonate, descriptions []devicemodel.DeviceDescription) ([]devicemodel.Selectable, error, int) {
	if len(descriptions) == 0 {
		return nil, errors.New("missing descriptions"), 500
	}
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.options, nil, 200
}

func (this *DeviceRepoMock) SetOptions(options []devicemodel.Selectable) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.options = options
}
