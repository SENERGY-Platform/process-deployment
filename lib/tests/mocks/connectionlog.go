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
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"sync"
)

type ConnectionlogMock struct {
	states map[string]bool
	mux    sync.Mutex
}

var Connectionlog = &ConnectionlogMock{states: map[string]bool{}}

func (this *ConnectionlogMock) New(ctx context.Context, config config.Config) (interfaces.Connectionlog, error) {
	return this, nil
}

func (this *ConnectionlogMock) CheckDeviceStates(jwtimpersonate jwt_http_router.JwtImpersonate, ids []string) (result map[string]bool, err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	result = map[string]bool{}
	for _, id := range ids {
		result[id] = this.states[id]
	}
	return
}

func (this *ConnectionlogMock) SetDeviceStates(states map[string]bool) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.states = states
}
