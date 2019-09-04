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
	"net/http"
	"sync"
)

type DatabaseMock struct {
	deployments  map[string]model.DeploymentCommand
	dependencies map[string]model.Dependencies
	mux          sync.Mutex
}

var Database = &DatabaseMock{
	deployments:  map[string]model.DeploymentCommand{},
	dependencies: map[string]model.Dependencies{},
}

func (this *DatabaseMock) New(ctx context.Context, config config.Config) (interfaces.Database, error) {
	return this, nil
}

func (this *DatabaseMock) CheckDeploymentAccess(user string, deploymentId string) (error, int) {
	return nil, 200
}

func (this *DatabaseMock) GetDeployment(user string, deploymentId string) (model.Deployment, error, int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	deployment, ok := this.deployments[deploymentId]
	if !ok {
		return model.Deployment{}, errors.New("deployment not found"), http.StatusNotFound
	}
	if deployment.Owner != user {
		return model.Deployment{}, errors.New("access denied"), http.StatusForbidden
	}
	return deployment.Deployment, nil, 200
}

func (this *DatabaseMock) DeleteDeployment(id string) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	delete(this.deployments, id)
	return nil
}

func (this *DatabaseMock) SetDeployment(id string, owner string, deployment model.Deployment) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.deployments[id] = model.DeploymentCommand{Id: id, Owner: owner, Deployment: deployment}
	return nil
}

func (this *DatabaseMock) GetDependencies(user string, deploymentId string) (model.Dependencies, error, int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	dependencies, ok := this.dependencies[deploymentId]
	if !ok {
		return model.Dependencies{}, errors.New("dependencies not found"), http.StatusNotFound
	}
	if dependencies.Owner != user {
		return model.Dependencies{}, errors.New("access denied"), http.StatusForbidden
	}
	return dependencies, nil, 200
}

func (this *DatabaseMock) GetDependenciesList(user string, limit int, offset int) (result []model.Dependencies, err error, code int) {
	count := 0
	for _, dependencie := range this.dependencies {
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

func (this *DatabaseMock) GetSelectedDependencies(user string, ids []string) (result []model.Dependencies, err error, code int) {
	for _, id := range ids {
		dependency, ok := this.dependencies[id]
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

func (this *DatabaseMock) SetDependencies(dependencies model.Dependencies) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.dependencies[dependencies.DeploymentId] = dependencies
	return nil
}

func (this *DatabaseMock) DeleteDependencies(id string) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	delete(this.dependencies, id)
	return nil
}
