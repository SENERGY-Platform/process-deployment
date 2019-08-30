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

package ctrl

import (
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/bpmn"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (this *Ctrl) PrepareDeployment(id string) (result model.Deployment, err error, code int) {
	xml, exists, err := this.GetBpmn(id)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	if !exists {
		return result, errors.New("process modell not found"), http.StatusNotFound
	}
	result, err = bpmn.PrepareDeployment(xml)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	err = this.SetDeploymentOptions(&result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

func (this *Ctrl) Deploy(deployment model.Deployment) (result model.Deployment, err error, code int) {
	//TODO
	deployment.Id = uuid.NewV4().String()
	if err := deployment.Validate(false); err != nil {
		return deployment, err, http.StatusBadRequest
	}

	//TODO: ensure selected devices and services exist and have the given content and are executable for the requesting user (if not using id ref)

	//TODO: use bpmn sub lib

	if err := deployment.Validate(true); err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	//TODO: use cqrs to deploy (maybe move validation to producing cqrs function)

	panic("not implemented")
}
