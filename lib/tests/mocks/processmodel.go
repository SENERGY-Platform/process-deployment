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
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/processmodel"

	"net/http"
	"sync"
)

type ProcessModelRepoMock struct {
	mux    sync.Mutex
	models map[string]processmodel.ProcessModel
}

var ProcessModelRepo = &ProcessModelRepoMock{models: map[string]processmodel.ProcessModel{}}

func (this *ProcessModelRepoMock) New(ctx context.Context, config config.Config) (interfaces.ProcessRepo, error) {
	return this, nil
}

func (this *ProcessModelRepoMock) GetProcessModel(impersonate auth.Token, id string) (result processmodel.ProcessModel, err error, errCode int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if result, ok := this.models[id]; ok {
		return result, nil, 200
	} else {
		return result, errors.New("process model " + id + " not found"), http.StatusNotFound
	}
}

func (this *ProcessModelRepoMock) SetProcessModel(id string, processmodel processmodel.ProcessModel) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.models[id] = processmodel
}
