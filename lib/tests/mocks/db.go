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
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"net/http"
	"sync"
	"time"
)

type DatabaseMock struct {
	Deployments  map[string]messages.DeploymentCommand
	Dependencies map[string]dependencymodel.Dependencies
	mux          sync.Mutex
}

var Database = &DatabaseMock{
	Deployments:  map[string]messages.DeploymentCommand{},
	Dependencies: map[string]dependencymodel.Dependencies{},
}

func (this *DatabaseMock) New(ctx context.Context, config config.Config) (interfaces.Database, error) {
	return this, nil
}

func (this *DatabaseMock) CheckDeploymentAccess(user string, deploymentId string) (error, int) {
	return nil, 200
}

func (this *DatabaseMock) GetDeploymentIds(user string) (deployments []string, err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	for _, depl := range this.Deployments {
		if depl.Owner == user {
			deployments = append(deployments, depl.Id)
		}
	}
	return
}

func (this *DatabaseMock) GetDeployment(user string, deploymentId string) (deployment *deploymentmodel.Deployment, err error, code int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	temp, ok := this.Deployments[deploymentId]
	if !ok {
		return nil, errors.New("deployment not found"), http.StatusNotFound
	}
	if temp.Owner != user {
		return nil, errors.New("access denied"), http.StatusForbidden
	}
	return temp.Deployment, nil, 200
}

func (this *DatabaseMock) ListDeployments(user string, options model.DeploymentListOptions) (deployments []deploymentmodel.Deployment, err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	for _, depl := range this.Deployments {
		if depl.Owner == user && depl.Deployment != nil {
			deployments = append(deployments, *depl.Deployment)
		}
	}
	return deployments, nil
}

func (this *DatabaseMock) GetDependencies(user string, deploymentId string) (dependencymodel.Dependencies, error, int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	dependencies, ok := this.Dependencies[deploymentId]
	if !ok {
		return dependencymodel.Dependencies{}, errors.New("dependencies not found"), http.StatusNotFound
	}
	if dependencies.Owner != user {
		return dependencymodel.Dependencies{}, errors.New("access denied"), http.StatusForbidden
	}
	return dependencies, nil, 200
}

func (this *DatabaseMock) SetDeployment(depl messages.DeploymentCommand, syncHandler func(messages.DeploymentCommand) error) error {
	this.mux.Lock()
	this.Deployments[depl.Id] = depl
	this.mux.Unlock()
	return syncHandler(depl)
}

func (this *DatabaseMock) DeleteDeployment(id string, syncDeleteHandler func(messages.DeploymentCommand) error) error {
	this.mux.Lock()
	element, ok := this.Deployments[id]
	if !ok {
		this.mux.Unlock()
		return nil
	}
	delete(this.Deployments, id)
	this.mux.Unlock()
	return syncDeleteHandler(element)
}

func (this *DatabaseMock) RetryDeploymentSync(lockduration time.Duration, syncDeleteHandler func(messages.DeploymentCommand) error, syncHandler func(messages.DeploymentCommand) error) error {
	return nil
}

func (this *DatabaseMock) GetDependenciesList(user string, limit int, offset int) (result []dependencymodel.Dependencies, err error, code int) {
	count := 0
	for _, dependencie := range this.Dependencies {
		if dependencie.Owner == user {
			if count >= (limit + offset) {
				return result, nil, 200
			}
			if count >= offset {
				result = append(result, dependencie)
			}
			count = count + 1
		}
	}
	return result, nil, 200
}

func (this *DatabaseMock) GetSelectedDependencies(user string, ids []string) (result []dependencymodel.Dependencies, err error, code int) {
	for _, id := range ids {
		dependency, ok := this.Dependencies[id]
		if !ok {
			return result, errors.New("unknown id"), http.StatusNotFound
		}
		result = append(result, dependency)
	}
	for _, dependency := range result {
		if dependency.Owner != user {
			return result, errors.New("user dosnt have access to given id"), http.StatusForbidden
		}
	}
	return result, nil, 200
}

func (this *DatabaseMock) SetDependencies(dependencies dependencymodel.Dependencies) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.Dependencies[dependencies.DeploymentId] = dependencies
	return nil
}

func (this *DatabaseMock) DeleteDependencies(id string) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	delete(this.Dependencies, id)
	return nil
}
