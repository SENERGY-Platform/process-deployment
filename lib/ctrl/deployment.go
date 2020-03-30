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
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/bpmn"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/bpmn/stringify"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
	"sort"
	"time"
)

func (this *Ctrl) HandleDeployment(cmd model.DeploymentCommand) error {
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

func (this *Ctrl) PrepareDeployment(token jwt_http_router.JwtImpersonate, xml string, svg string) (result model.Deployment, err error, code int) {
	startParsing := time.Now()
	result, err = bpmn.PrepareDeployment(xml)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	durParsing := time.Now().Sub(startParsing)
	log.Println("DEBUG: prepare deployment parsing time:", durParsing, durParsing.Milliseconds())
	startSelectables := time.Now()
	err = this.SetDeploymentOptions(token, &result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	durSelectables := time.Now().Sub(startSelectables)
	log.Println("DEBUG: prepare deployment selectables time:", durSelectables, durSelectables.Milliseconds())
	result.Svg = svg
	this.SetExecutableFlag(&result)
	return result, nil, http.StatusOK
}

func (this *Ctrl) GetDeployment(jwt jwt_http_router.Jwt, id string) (result model.Deployment, err error, code int) {
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

func (this *Ctrl) CreateDeployment(jwt jwt_http_router.Jwt, deployment model.Deployment) (result model.Deployment, err error, code int) {
	deployment.Id = config.NewId()
	return this.setDeployment(jwt, deployment)
}

func (this *Ctrl) UpdateDeployment(jwt jwt_http_router.Jwt, id string, deployment model.Deployment) (result model.Deployment, err error, code int) {
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

func (this *Ctrl) DeleteDeployment(command model.DeploymentCommand) error {
	return this.db.DeleteDeployment(command.Id)
}

func (this *Ctrl) SaveDeployment(command model.DeploymentCommand) error {
	return this.db.SetDeployment(command.Id, command.Owner, command.Deployment)
}

func (this *Ctrl) setDeployment(jwt jwt_http_router.Jwt, deployment model.Deployment) (result model.Deployment, err error, code int) {
	if err := deployment.Validate(false); err != nil {
		return deployment, err, http.StatusBadRequest
	}

	//ensure selected devices and services exist and have the given content and are executable for the requesting user (if not using id ref)
	err, code = this.ensureDeploymentSelectionCorrectness(jwt.Impersonate, &deployment)
	if err != nil {
		return deployment, err, code
	}

	err = this.completeEvents(&deployment)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	deployment.Xml, err = stringify.Deployment(deployment, this.config.DeploymentAsRef, this.devices, jwt.UserId, this.config.NotificationUrl)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	sort.Sort(bpmn.LaneByOrder(deployment.Lanes))
	sort.Sort(bpmn.ElementByOrder(deployment.Elements))

	if err = this.publishDeployment(jwt.UserId, deployment); err != nil {
		return deployment, err, http.StatusInternalServerError
	}
	return deployment, nil, 200
}

func (this *Ctrl) publishDeployment(owner string, deployment model.Deployment) error {
	if err := deployment.Validate(true); err != nil {
		return err
	}
	cmd := model.DeploymentCommand{
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
	cmd := model.DeploymentCommand{
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
func (this *Ctrl) ensureDeploymentSelectionCorrectness(token jwt_http_router.JwtImpersonate, deployment *model.Deployment) (err error, code int) {
	deviceCache := map[string]devicemodel.Device{}
	serviceCache := map[string]devicemodel.Service{}

	for _, lane := range deployment.Lanes {
		if lane.Lane != nil {
			if lane.Lane.Selection != nil && lane.Lane.Selection.Id != "" {
				lane.Lane.Selection, err, code = this.getCachedDevice(token, &deviceCache, lane.Lane.Selection.Id)
				if err != nil {
					return err, code
				}
			}
			for _, element := range lane.Lane.Elements {
				if element.LaneTask != nil {
					element.LaneTask.SelectedService, err, code = this.getCachedService(token, &serviceCache, element.LaneTask.SelectedService.Id)
					if err != nil {
						return err, code
					}
				}
				if element.MsgEvent != nil {
					element.MsgEvent.Device, err, code = this.getCachedDevice(token, &deviceCache, element.MsgEvent.Device.Id)
					if err != nil {
						return err, code
					}
					element.MsgEvent.Service, err, code = this.getCachedService(token, &serviceCache, element.MsgEvent.Service.Id)
					if err != nil {
						return err, code
					}
				}
				if element.ReceiveTaskEvent != nil {
					element.ReceiveTaskEvent.Device, err, code = this.getCachedDevice(token, &deviceCache, element.ReceiveTaskEvent.Device.Id)
					if err != nil {
						return err, code
					}
					element.ReceiveTaskEvent.Service, err, code = this.getCachedService(token, &serviceCache, element.ReceiveTaskEvent.Service.Id)
					if err != nil {
						return err, code
					}
				}
			}
		}
		if lane.MultiLane != nil {
			for index, selection := range lane.MultiLane.Selections {
				selection, err, code = this.getCachedDevice(token, &deviceCache, selection.Id)
				if err != nil {
					return err, code
				}
				lane.MultiLane.Selections[index] = selection
			}
			for _, element := range lane.MultiLane.Elements {
				if element.LaneTask != nil {
					element.LaneTask.SelectedService, err, code = this.getCachedService(token, &serviceCache, element.LaneTask.SelectedService.Id)
					if err != nil {
						return err, code
					}
				}
				if element.MsgEvent != nil {
					element.MsgEvent.Device, err, code = this.getCachedDevice(token, &deviceCache, element.MsgEvent.Device.Id)
					if err != nil {
						return err, code
					}
					element.MsgEvent.Service, err, code = this.getCachedService(token, &serviceCache, element.MsgEvent.Service.Id)
					if err != nil {
						return err, code
					}
				}
				if element.ReceiveTaskEvent != nil {
					element.ReceiveTaskEvent.Device, err, code = this.getCachedDevice(token, &deviceCache, element.ReceiveTaskEvent.Device.Id)
					if err != nil {
						return err, code
					}
					element.ReceiveTaskEvent.Service, err, code = this.getCachedService(token, &serviceCache, element.ReceiveTaskEvent.Service.Id)
					if err != nil {
						return err, code
					}
				}
			}
		}
	}

	for _, element := range deployment.Elements {
		if element.Task != nil {
			element.Task.Selection.Device, err, code = this.getCachedDevice(token, &deviceCache, element.Task.Selection.Device.Id)
			if err != nil {
				return err, code
			}
			element.Task.Selection.Service, err, code = this.getCachedService(token, &serviceCache, element.Task.Selection.Service.Id)
			if err != nil {
				return err, code
			}
		}
		if element.MultiTask != nil {
			for index, selection := range element.MultiTask.Selections {
				selection.Device, err, code = this.getCachedDevice(token, &deviceCache, selection.Device.Id)
				if err != nil {
					return err, code
				}
				selection.Service, err, code = this.getCachedService(token, &serviceCache, selection.Service.Id)
				if err != nil {
					return err, code
				}
				element.MultiTask.Selections[index] = selection
			}
		}
		if element.MsgEvent != nil {
			element.MsgEvent.Device, err, code = this.getCachedDevice(token, &deviceCache, element.MsgEvent.Device.Id)
			if err != nil {
				return err, code
			}
			element.MsgEvent.Service, err, code = this.getCachedService(token, &serviceCache, element.MsgEvent.Service.Id)
			if err != nil {
				return err, code
			}
		}
		if element.ReceiveTaskEvent != nil {
			element.ReceiveTaskEvent.Device, err, code = this.getCachedDevice(token, &deviceCache, element.ReceiveTaskEvent.Device.Id)
			if err != nil {
				return err, code
			}
			element.ReceiveTaskEvent.Service, err, code = this.getCachedService(token, &serviceCache, element.ReceiveTaskEvent.Service.Id)
			if err != nil {
				return err, code
			}
		}
	}
	return nil, 200
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

func (this *Ctrl) SetExecutableFlag(deployment *model.Deployment) {
	deployment.Executable = true
	for _, lane := range deployment.Lanes {
		this.setExecutableFlagByLane(deployment, lane)
	}
	for _, element := range deployment.Elements {
		this.SetExecutableFlagByElement(deployment, element)
	}
}

func (this *Ctrl) setExecutableFlagByLane(deployment *model.Deployment, lane model.LaneElement) {
	if lane.Lane != nil && len(lane.Lane.Selectables) == 0 {
		deployment.Executable = false
	}
	if lane.MultiLane != nil && len(lane.MultiLane.Selectables) == 0 {
		deployment.Executable = false
	}
}

func (this *Ctrl) SetExecutableFlagByElement(deployment *model.Deployment, element model.Element) {
	if element.Task != nil && len(element.Task.Selectables) == 0 {
		deployment.Executable = false
	}
	if element.MultiTask != nil && len(element.MultiTask.Selectables) == 0 {
		deployment.Executable = false
	}
}
