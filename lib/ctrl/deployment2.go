/*
 * Copyright 2020 InfAI (CC SES)
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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
)

func (this *Ctrl) PrepareDeploymentV2(token jwt_http_router.JwtImpersonate, xml string, svg string) (result deploymentmodel.Deployment, err error, code int) {
	result, err = this.deploymentParser.PrepareDeployment(xml)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	err = this.SetDeploymentOptionsV2(token, &result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	result.Diagram.Svg = svg
	this.SetExecutableFlagV2(&result)
	return result, nil, http.StatusOK
}

func (this *Ctrl) GetDeploymentV2(jwt jwt_http_router.Jwt, id string) (result deploymentmodel.Deployment, err error, code int) {
	_, temp, err, code := this.db.GetDeployment(jwt.UserId, id)
	if err != nil {
		return result, err, code
	}
	if temp == nil {
		return result, errors.New("found deployment is not of requested version"), http.StatusBadRequest
	}
	result = *temp
	err = this.SetDeploymentOptionsV2(jwt.Impersonate, &result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return
}

func (this *Ctrl) CreateDeploymentV2(jwt jwt_http_router.Jwt, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
	deployment.Id = config.NewId()
	return this.setDeploymentV2(jwt, deployment)
}

func (this *Ctrl) UpdateDeploymentV2(jwt jwt_http_router.Jwt, id string, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
	if id != deployment.Id {
		return deployment, errors.New("path id != body id"), http.StatusBadRequest
	}

	err, code = this.db.CheckDeploymentAccess(jwt.UserId, id)
	if err != nil {
		return result, err, code
	}

	return this.setDeploymentV2(jwt, deployment)
}

func (this *Ctrl) RemoveDeploymentV2(jwt jwt_http_router.Jwt, id string) (err error, code int) {
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

func (this *Ctrl) SetExecutableFlagV2(deployment *deploymentmodel.Deployment) {
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

func (this *Ctrl) SetDeploymentOptionsV2(token jwt_http_router.JwtImpersonate, deployment *deploymentmodel.Deployment) (err error) {
	panic("not implemented")
	return nil
}

func (this *Ctrl) setDeploymentV2(jwt jwt_http_router.Jwt, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
	if err := deployment.Validate(deploymentmodel.ValidateRequest); err != nil {
		return deployment, err, http.StatusBadRequest
	}

	//ensure selected devices and services exist and have the given content and are executable for the requesting user (if not using id ref)
	err, code = this.ensureDeploymentSelectionAccess(jwt.Impersonate, &deployment)
	if err != nil {
		return deployment, err, code
	}

	err = this.completeEventsV2(&deployment)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	deployment.Diagram.XmlDeployed, err = this.deploymentStringifier.Deployment(deployment, jwt.UserId)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	if err = this.publishDeploymentV2(jwt.UserId, deployment.Id, deployment); err != nil {
		return deployment, err, http.StatusInternalServerError
	}
	return deployment, nil, 200
}

//ensures selection correctness
func (this *Ctrl) ensureDeploymentSelectionAccess(token jwt_http_router.JwtImpersonate, deployment *deploymentmodel.Deployment) (err error, code int) {
	panic("not implemented")
}