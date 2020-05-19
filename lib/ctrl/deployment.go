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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Ctrl) HandleDeployment(cmd messages.DeploymentCommand) error {
	switch cmd.Command {
	case "PUT":
		err := this.SaveDependencies(cmd)
		if err != nil {
			return err
		}
		err = this.SaveDeployment(cmd)
		if err != nil {
			return err
		}
		return nil
	case "DELETE":
		err := this.DeleteDependencies(cmd)
		if err != nil {
			return err
		}
		err = this.DeleteDeployment(cmd)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("unable to handle command: " + cmd.Command)
	}
}

func (this *Ctrl) PrepareDeployment(token jwt_http_router.JwtImpersonate, xml string, svg string) (result deploymentmodel.Deployment, err error, code int) {
	result, err = this.deploymentParser.PrepareDeployment(xml)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	err = this.SetDeploymentOptions(token, &result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	result.Diagram.Svg = svg
	this.SetExecutableFlag(&result)
	return result, nil, http.StatusOK
}

func (this *Ctrl) GetDeployment(jwt jwt_http_router.Jwt, id string) (result deploymentmodel.Deployment, err error, code int) {
	result, err, code = this.db.GetDeployment(jwt.UserId, id)
	if err != nil {
		return result, err, code
	}
	err = this.SetDeploymentOptions(jwt.Impersonate, &result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return
}

func (this *Ctrl) CreateDeployment(jwt jwt_http_router.Jwt, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
	deployment.Id = config.NewId()
	return this.setDeployment(jwt, deployment)
}

func (this *Ctrl) UpdateDeployment(jwt jwt_http_router.Jwt, id string, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
	if id != deployment.Id {
		return deployment, errors.New("path id != body id"), http.StatusBadRequest
	}

	err, code = this.db.CheckDeploymentAccess(jwt.UserId, id)
	if err != nil {
		return result, err, code
	}

	return this.setDeployment(jwt, deployment)
}

func (this *Ctrl) RemoveDeployment(jwt jwt_http_router.Jwt, id string) (err error, code int) {
	err, code = this.db.CheckDeploymentAccess(jwt.UserId, id)
	if err != nil {
		return err, code
	}

	err = this.publishDeploymentDelete(jwt.UserId, id)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, 200
}

func (this *Ctrl) DeleteDeployment(command messages.DeploymentCommand) error {
	return this.db.DeleteDeployment(command.Id)
}

func (this *Ctrl) SaveDeployment(command messages.DeploymentCommand) error {
	return this.db.SetDeployment(command.Id, command.Owner, command.Deployment)
}

func (this *Ctrl) setDeployment(jwt jwt_http_router.Jwt, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
	if err := deployment.Validate(deploymentmodel.ValidateRequest); err != nil {
		return deployment, err, http.StatusBadRequest
	}

	//ensure selected devices and services exist and have the given content and are executable for the requesting user (if not using id ref)
	err, code = this.ensureDeploymentSelectionAccess(jwt.Impersonate, &deployment)
	if err != nil {
		return deployment, err, code
	}

	err = this.completeEvents(&deployment)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	deployment.Diagram.XmlDeployed, err = this.deploymentStringifier.Deployment(deployment, jwt.UserId)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	if err = this.publishDeployment(jwt.UserId, deployment); err != nil {
		return deployment, err, http.StatusInternalServerError
	}
	return deployment, nil, 200
}

func (this *Ctrl) publishDeployment(owner string, deployment deploymentmodel.Deployment) error {
	if err := deployment.Validate(deploymentmodel.ValidatePublish); err != nil {
		return err
	}
	cmd := messages.DeploymentCommand{
		Command:    "PUT",
		Id:         deployment.Id,
		Owner:      owner,
		Deployment: deployment,
	}
	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.deploymentPublisher.Produce(deployment.Id, msg)
}

func (this *Ctrl) publishDeploymentDelete(user string, id string) error {
	cmd := messages.DeploymentCommand{
		Command: "DELETE",
		Id:      id,
		Owner:   user,
	}
	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.deploymentPublisher.Produce(id, msg)
}

//ensures selection correctness
func (this *Ctrl) ensureDeploymentSelectionAccess(token jwt_http_router.JwtImpersonate, deployment *deploymentmodel.Deployment) (err error, code int) {
	panic("not implemented")
}

func (this *Ctrl) getCachedDevice(token jwt_http_router.JwtImpersonate, cache *map[string]devicemodel.Device, id string) (*devicemodel.Device, error, int) {
	if result, ok := (*cache)[id]; ok {
		return &result, nil, 200
	}
	result, err, code := this.devices.GetDevice(token, id)
	if err != nil {
		return &result, err, code
	}
	(*cache)[id] = result
	return &result, nil, 200
}

func (this *Ctrl) getCachedService(token jwt_http_router.JwtImpersonate, cache *map[string]devicemodel.Service, id string) (*devicemodel.Service, error, int) {
	if result, ok := (*cache)[id]; ok {
		return &result, nil, 200
	}
	result, err, code := this.devices.GetService(token, id)
	if err != nil {
		return &result, err, code
	}
	(*cache)[id] = result
	return &result, nil, 200
}

func (this *Ctrl) SetExecutableFlag(deployment *deploymentmodel.Deployment) {
	deployment.Executable = true
	for _, element := range deployment.Elements {
		if element.Task != nil && len(element.Task.Selection.SelectionOptions) < 0 {
			deployment.Executable = false
			return
		}
		if element.MessageEvent != nil && len(element.MessageEvent.Selection.SelectionOptions) < 0 {
			deployment.Executable = false
			return
		}
	}
	return
}
