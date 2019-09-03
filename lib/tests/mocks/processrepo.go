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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"sync"
)

type ProcessRepoMock struct {
	mux       sync.Mutex
	processes map[string]string
}

var ProcessRepository = &ProcessRepoMock{processes: map[string]string{}}

func (this *ProcessRepoMock) New(ctx context.Context, config config.Config) (interfaces.ProcessRepository, error) {
	return this, nil
}

func (this *ProcessRepoMock) GetBpmn(id string) (xml string, exists bool, err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	xml, exists = this.processes[id]
	return
}

func (this *ProcessRepoMock) SetBpmn(id string, xml string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.processes[id] = xml
	return
}
