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
	"github.com/SENERGY-Platform/process-deployment/lib/model"
)

func (this *Ctrl) SaveDependencies(command model.DeploymentCommand) error {
	dependencies, err := this.deploymentToDependencies(command.Deployment)
	if err != nil {
		return err
	}
	dependencies.Owner = command.Owner
	return this.db.SetDependencies(dependencies)
}

func (this *Ctrl) DeleteDependencies(command model.DeploymentCommand) error {
	return this.db.DeleteDependencies(command.Id)
}

func (this *Ctrl) deploymentToDependencies(deployment model.Deployment) (result model.Dependencies, err error) {
	result.DeploymentId = deployment.Id
	for _, lane := range deployment.Lanes {
		if lane.Lane != nil {
			dependencie := model.DeviceDependency{
				DeviceId: lane.Lane.Selection.Id,
				Name:     lane.Lane.Selection.Name,
				BpmnResources: []model.BpmnResource{{
					Id: lane.Lane.BpmnElementId,
				}},
			}
			for _, element := range lane.Lane.Elements {
				if element.LaneTask != nil {
					dependencie.BpmnResources = append(dependencie.BpmnResources, model.BpmnResource{Id: element.LaneTask.BpmnElementId})
				}
				if element.MsgEvent != nil {
					result.Devices = append(result.Devices, model.DeviceDependency{
						DeviceId: element.MsgEvent.Device.Id,
						Name:     element.MsgEvent.Device.Name,
						BpmnResources: []model.BpmnResource{{
							Id: element.MsgEvent.BpmnElementId,
						}},
					})
					result.Events = append(result.Events, model.EventDependency{
						EventId: element.MsgEvent.EventId,
						BpmnResources: []model.BpmnResource{{
							Id: element.MsgEvent.BpmnElementId,
						}},
					})
				}
				if element.ReceiveTaskEvent != nil {
					result.Devices = append(result.Devices, model.DeviceDependency{
						DeviceId: element.ReceiveTaskEvent.Device.Id,
						Name:     element.ReceiveTaskEvent.Device.Name,
						BpmnResources: []model.BpmnResource{{
							Id: element.ReceiveTaskEvent.BpmnElementId,
						}},
					})
					result.Events = append(result.Events, model.EventDependency{
						EventId: element.ReceiveTaskEvent.EventId,
						BpmnResources: []model.BpmnResource{{
							Id: element.ReceiveTaskEvent.BpmnElementId,
						}},
					})
				}
			}
			result.Devices = append(result.Devices, dependencie)
		}
		if lane.MultiLane != nil {
			for _, selection := range lane.MultiLane.Selections {
				dependencie := model.DeviceDependency{
					DeviceId: selection.Id,
					Name:     selection.Name,
					BpmnResources: []model.BpmnResource{{
						Id: lane.MultiLane.BpmnElementId,
					}},
				}
				for _, element := range lane.MultiLane.Elements {
					if element.LaneTask != nil {
						dependencie.BpmnResources = append(dependencie.BpmnResources, model.BpmnResource{Id: element.LaneTask.BpmnElementId})
					}
				}
				result.Devices = append(result.Devices, dependencie)
			}
			for _, element := range lane.MultiLane.Elements {
				if element.MsgEvent != nil {
					result.Devices = append(result.Devices, model.DeviceDependency{
						DeviceId: element.MsgEvent.Device.Id,
						Name:     element.MsgEvent.Device.Name,
						BpmnResources: []model.BpmnResource{{
							Id: element.MsgEvent.BpmnElementId,
						}},
					})
					result.Events = append(result.Events, model.EventDependency{
						EventId: element.MsgEvent.EventId,
						BpmnResources: []model.BpmnResource{{
							Id: element.MsgEvent.BpmnElementId,
						}},
					})
				}
				if element.ReceiveTaskEvent != nil {
					result.Devices = append(result.Devices, model.DeviceDependency{
						DeviceId: element.ReceiveTaskEvent.Device.Id,
						Name:     element.ReceiveTaskEvent.Device.Name,
						BpmnResources: []model.BpmnResource{{
							Id: element.ReceiveTaskEvent.BpmnElementId,
						}},
					})
					result.Events = append(result.Events, model.EventDependency{
						EventId: element.ReceiveTaskEvent.EventId,
						BpmnResources: []model.BpmnResource{{
							Id: element.ReceiveTaskEvent.BpmnElementId,
						}},
					})
				}
			}
		}
	}

	for _, element := range deployment.Elements {
		if element.Task != nil {
			result.Devices = append(result.Devices, model.DeviceDependency{
				DeviceId: element.Task.Selection.SelectedDevice.Id,
				Name:     element.Task.Selection.SelectedDevice.Name,
				BpmnResources: []model.BpmnResource{{
					Id: element.Task.BpmnElementId,
				}},
			})
		}
		if element.MultiTask != nil {
			for _, selection := range element.MultiTask.Selections {
				result.Devices = append(result.Devices, model.DeviceDependency{
					DeviceId: selection.SelectedDevice.Id,
					Name:     selection.SelectedDevice.Name,
					BpmnResources: []model.BpmnResource{{
						Id: element.MultiTask.BpmnElementId,
					}},
				})
			}
		}
		if element.MsgEvent != nil {
			result.Devices = append(result.Devices, model.DeviceDependency{
				DeviceId: element.MsgEvent.Device.Id,
				Name:     element.MsgEvent.Device.Name,
				BpmnResources: []model.BpmnResource{{
					Id: element.MsgEvent.BpmnElementId,
				}},
			})
			result.Events = append(result.Events, model.EventDependency{
				EventId: element.MsgEvent.EventId,
				BpmnResources: []model.BpmnResource{{
					Id: element.MsgEvent.BpmnElementId,
				}},
			})
		}
		if element.ReceiveTaskEvent != nil {
			result.Devices = append(result.Devices, model.DeviceDependency{
				DeviceId: element.ReceiveTaskEvent.Device.Id,
				Name:     element.ReceiveTaskEvent.Device.Name,
				BpmnResources: []model.BpmnResource{{
					Id: element.ReceiveTaskEvent.BpmnElementId,
				}},
			})
			result.Events = append(result.Events, model.EventDependency{
				EventId: element.ReceiveTaskEvent.EventId,
				BpmnResources: []model.BpmnResource{{
					Id: element.ReceiveTaskEvent.BpmnElementId,
				}},
			})
		}
	}
	result.Devices = reduceDeviceDependencies(result.Devices)
	return result, nil
}

func reduceDeviceDependencies(dependencies []model.DeviceDependency) (result []model.DeviceDependency) {
	name := map[string]string{}
	resources := map[string][]model.BpmnResource{}
	for _, device := range dependencies {
		name[device.DeviceId] = device.Name
		resources[device.DeviceId] = append(resources[device.DeviceId], device.BpmnResources...)
	}
	for id, name := range name {
		result = append(result, model.DeviceDependency{
			DeviceId:      id,
			Name:          name,
			BpmnResources: resources[id],
		})
	}
	return
}
