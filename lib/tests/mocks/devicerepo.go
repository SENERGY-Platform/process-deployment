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

package mock

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"sync"
)

type DeviceRepoMock struct {
	mux       sync.Mutex
	protocols map[string]devicemodel.Protocol
}

var DeviceRepository = &DeviceRepoMock{protocols: map[string]devicemodel.Protocol{}}

func (this *DeviceRepoMock) New(ctx context.Context, config config.Config) (interfaces.DeviceRepository, error) {
	return this, nil
}

func (this *DeviceRepoMock) GetProtocol(id string) (devicemodel.Protocol, error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.protocols[id], nil
}

func (this *DeviceRepoMock) SetProtocol(id string, protocol devicemodel.Protocol) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.protocols[id] = protocol
}
