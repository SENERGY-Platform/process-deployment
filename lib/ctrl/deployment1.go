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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/bpmn"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/bpmn/stringify"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
	"sort"
	"time"
)

func (this *Ctrl) PrepareDeploymentV1(token jwt_http_router.JwtImpersonate, xml string, svg string) (result deploymentmodel.Deployment, err error, code int) {
	startParsing := time.Now()
	result, err = bpmn.PrepareDeployment(xml)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	durParsing := time.Now().Sub(startParsing)
	log.Println("DEBUG: prepare deployment parsing time:", durParsing, durParsing.Milliseconds())
	startSelectables := time.Now()
	err = this.SetDeploymentOptionsV1(token, &result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	durSelectables := time.Now().Sub(startSelectables)
	log.Println("DEBUG: prepare deployment selectables time:", durSelectables, durSelectables.Milliseconds())
	result.Svg = svg
	this.SetExecutableFlagV1(&result)
	return result, nil, http.StatusOK
}

func (this *Ctrl) GetDeploymentV1(jwt jwt_http_router.Jwt, id string) (result deploymentmodel.Deployment, err error, code int) {
	temp, _, err, code := this.db.GetDeployment(jwt.UserId, id)
	if err != nil {
		return result, err, code
	}
	if temp == nil {
		return result, errors.New("found deployment is not of requested version"), http.StatusBadRequest
	}
	result = *temp
	err = this.SetDeploymentOptionsV1(jwt.Impersonate, &result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return
}

func (this *Ctrl) CreateDeploymentV1(jwt jwt_http_router.Jwt, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
	deployment.Id = config.NewId()
	return this.setDeploymentV1(jwt, deployment)
}

func (this *Ctrl) UpdateDeploymentV1(jwt jwt_http_router.Jwt, id string, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
	if id != deployment.Id {
		return deployment, errors.New("path id != body id"), http.StatusBadRequest
	}

	err, code = this.db.CheckDeploymentAccess(jwt.UserId, id)
	if err != nil {
		return result, err, code
	}

	return this.setDeploymentV1(jwt, deployment)
}

func (this *Ctrl) RemoveDeploymentV1(jwt jwt_http_router.Jwt, id string) (err error, code int) {
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

//ensures selection correctness
func (this *Ctrl) ensureDeploymentSelectionCorrectness(token jwt_http_router.JwtImpersonate, deployment *deploymentmodel.Deployment) (err error, code int) {
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

func (this *Ctrl) SetExecutableFlagV1(deployment *deploymentmodel.Deployment) {
	deployment.Executable = true
	for _, lane := range deployment.Lanes {
		this.setExecutableFlagByLaneV1(deployment, lane)
	}
	for _, element := range deployment.Elements {
		this.SetExecutableFlagByElementV1(deployment, element)
	}
}

func (this *Ctrl) setExecutableFlagByLaneV1(deployment *deploymentmodel.Deployment, lane deploymentmodel.LaneElement) {
	if lane.Lane != nil && len(lane.Lane.Selectables) == 0 {
		deployment.Executable = false
	}
	if lane.MultiLane != nil && len(lane.MultiLane.Selectables) == 0 {
		deployment.Executable = false
	}
}

func (this *Ctrl) SetExecutableFlagByElementV1(deployment *deploymentmodel.Deployment, element deploymentmodel.Element) {
	if element.Task != nil && len(element.Task.Selectables) == 0 {
		deployment.Executable = false
	}
	if element.MultiTask != nil && len(element.MultiTask.Selectables) == 0 {
		deployment.Executable = false
	}
}

func (this *Ctrl) SetDeploymentOptionsV1(token jwt_http_router.JwtImpersonate, deployment *deploymentmodel.Deployment) (err error) {
	for index, element := range deployment.Elements {
		if element.Task != nil {
			options, err := this.GetOptions(token, deploymentmodel.DeviceDescriptions{element.Task.DeviceDescription}.ToFilter())
			if err != nil {
				return err
			}
			element.Task.Selectables = options
		}
		if element.MultiTask != nil {
			options, err := this.GetOptions(token, deploymentmodel.DeviceDescriptions{element.MultiTask.DeviceDescription}.ToFilter())
			if err != nil {
				return err
			}
			element.MultiTask.Selectables = options
		}
		deployment.Elements[index] = element
	}
	for index, lane := range deployment.Lanes {
		if lane.Lane != nil {
			options, err := this.GetOptions(token, lane.Lane.DeviceDescriptions.ToFilter())
			if err != nil {
				return err
			}
			lane.Lane.Selectables = options
		}
		if lane.MultiLane != nil {
			options, err := this.GetOptions(token, lane.MultiLane.DeviceDescriptions.ToFilter())
			if err != nil {
				return err
			}
			lane.MultiLane.Selectables = options
		}
		deployment.Lanes[index] = lane
	}
	return nil
}

func (this *Ctrl) setDeploymentV1(jwt jwt_http_router.Jwt, deployment deploymentmodel.Deployment) (result deploymentmodel.Deployment, err error, code int) {
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

	deployment.Xml, err = stringify.Deployment(deployment, this.devices, jwt.UserId, this.config.NotificationUrl)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	sort.Sort(bpmn.LaneByOrder(deployment.Lanes))
	sort.Sort(bpmn.ElementByOrder(deployment.Elements))

	if err = this.publishDeploymentV1(jwt.UserId, deployment.Id, deployment); err != nil {
		return deployment, err, http.StatusInternalServerError
	}
	return deployment, nil, 200
}
