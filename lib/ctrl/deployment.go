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
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	deploymentmodel2 "github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
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

func (this *Ctrl) DeleteDeployment(command messages.DeploymentCommand) error {
	return this.db.DeleteDeployment(command.Id)
}

func (this *Ctrl) SaveDeployment(command messages.DeploymentCommand) error {
	return this.db.SetDeployment(command.Id, command.Owner, command.Deployment, command.DeploymentV2)
}

func (this *Ctrl) publishDeploymentV1(owner string, id string, deployment deploymentmodel.Deployment, source string) error {
	if err := deployment.Validate(true); err != nil {
		return err
	}
	cmd := messages.DeploymentCommand{
		Command:    "PUT",
		Id:         id,
		Owner:      owner,
		Deployment: &deployment,
		Source:     source,
	}
	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.deploymentPublisher.Produce(id, msg)
}

func (this *Ctrl) publishDeploymentV2(owner string, id string, deployment deploymentmodel2.Deployment, source string) error {
	if err := deployment.Validate(true); err != nil {
		return err
	}
	cmd := messages.DeploymentCommand{
		Command:      "PUT",
		Id:           id,
		Owner:        owner,
		DeploymentV2: &deployment,
		Source:       source,
	}
	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return this.deploymentPublisher.Produce(id, msg)
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

func (this *Ctrl) GetOptions(token jwt_http_router.JwtImpersonate, descriptions devicemodel.DeviceTypesFilter, protocolBlockList []string) (result []devicemodel.Selectable, err error) {
	if len(descriptions) == 0 {
		return []devicemodel.Selectable{}, nil
	}
	result, err, _ = this.devices.GetFilteredDevices(token, descriptions, protocolBlockList)
	return
}

func (this *Ctrl) GetBlockedProtocols() (result []string, err error) {
	protocols, err, _ := this.devices.GetProtocols()
	if err != nil {
		return result, err
	}
	return this.FilterProtocols(protocols, devicemodel.EVENT)
}

func (this *Ctrl) FilterProtocols(protocols []devicemodel.Protocol, filterBy devicemodel.Interaction) (result []string, err error) {
	for _, protocol := range protocols {
		if protocol.Interaction == filterBy {
			result = append(result, protocol.Id)
		}
	}
	return result, nil
}
