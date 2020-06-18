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
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	deploymentmodel2 "github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"runtime/debug"
	"sort"
)

func (this *Ctrl) GetDependencies(jwt jwt_http_router.Jwt, id string) (dependencymodel.Dependencies, error, int) {
	return this.db.GetDependencies(jwt.UserId, id)
}

func (this *Ctrl) GetDependenciesList(jwt jwt_http_router.Jwt, limit int, offset int) ([]dependencymodel.Dependencies, error, int) {
	return this.db.GetDependenciesList(jwt.UserId, limit, offset)
}

func (this *Ctrl) GetSelectedDependencies(jwt jwt_http_router.Jwt, ids []string) ([]dependencymodel.Dependencies, error, int) {
	return this.db.GetSelectedDependencies(jwt.UserId, ids)
}

func (this *Ctrl) SaveDependencies(command messages.DeploymentCommand) error {
	v1, v2, err := CastDeploymentVersion(command.Deployment)
	if err != nil {
		//allow unknown versions; and save empty dependencies
		log.Println("WARNING: ", err)
		debug.PrintStack()
		return this.db.SetDependencies(dependencymodel.Dependencies{
			DeploymentId: command.Id,
			Owner:        command.Owner,
			Devices:      []dependencymodel.DeviceDependency{},
			Events:       []dependencymodel.EventDependency{},
		})
	}
	var dependencies dependencymodel.Dependencies
	if v1 != nil {
		dependencies, err = this.deploymentToDependenciesV1(*v1)
	} else if v2 != nil {
		dependencies, err = this.deploymentToDependenciesV2(*v2)
	}
	if err != nil {
		return err
	}
	dependencies.Owner = command.Owner
	return this.db.SetDependencies(dependencies)
}

func (this *Ctrl) DeleteDependencies(command messages.DeploymentCommand) error {
	return this.db.DeleteDependencies(command.Id)
}

func (this *Ctrl) deploymentToDependenciesV1(deployment deploymentmodel.Deployment) (result dependencymodel.Dependencies, err error) {
	result.DeploymentId = deployment.Id
	for _, lane := range deployment.Lanes {
		if lane.Lane != nil && lane.Lane.Selection != nil && lane.Lane.Selection.Id != "" {
			dependencie := dependencymodel.DeviceDependency{
				DeviceId: lane.Lane.Selection.Id,
				Name:     lane.Lane.Selection.Name,
				BpmnResources: []dependencymodel.BpmnResource{{
					Id: lane.Lane.BpmnElementId,
				}},
			}
			for _, element := range lane.Lane.Elements {
				if element.LaneTask != nil {
					dependencie.BpmnResources = append(dependencie.BpmnResources, dependencymodel.BpmnResource{Id: element.LaneTask.BpmnElementId})
				}
				if element.MsgEvent != nil {
					result.Devices = append(result.Devices, dependencymodel.DeviceDependency{
						DeviceId: element.MsgEvent.Device.Id,
						Name:     element.MsgEvent.Device.Name,
						BpmnResources: []dependencymodel.BpmnResource{{
							Id: element.MsgEvent.BpmnElementId,
						}},
					})
					result.Events = append(result.Events, dependencymodel.EventDependency{
						EventId: element.MsgEvent.EventId,
						BpmnResources: []dependencymodel.BpmnResource{{
							Id: element.MsgEvent.BpmnElementId,
						}},
					})
				}
				if element.ReceiveTaskEvent != nil {
					result.Devices = append(result.Devices, dependencymodel.DeviceDependency{
						DeviceId: element.ReceiveTaskEvent.Device.Id,
						Name:     element.ReceiveTaskEvent.Device.Name,
						BpmnResources: []dependencymodel.BpmnResource{{
							Id: element.ReceiveTaskEvent.BpmnElementId,
						}},
					})
					result.Events = append(result.Events, dependencymodel.EventDependency{
						EventId: element.ReceiveTaskEvent.EventId,
						BpmnResources: []dependencymodel.BpmnResource{{
							Id: element.ReceiveTaskEvent.BpmnElementId,
						}},
					})
				}
			}
			result.Devices = append(result.Devices, dependencie)
		}
		if lane.MultiLane != nil {
			for _, selection := range lane.MultiLane.Selections {
				dependencie := dependencymodel.DeviceDependency{
					DeviceId: selection.Id,
					Name:     selection.Name,
					BpmnResources: []dependencymodel.BpmnResource{{
						Id: lane.MultiLane.BpmnElementId,
					}},
				}
				for _, element := range lane.MultiLane.Elements {
					if element.LaneTask != nil {
						dependencie.BpmnResources = append(dependencie.BpmnResources, dependencymodel.BpmnResource{Id: element.LaneTask.BpmnElementId})
					}
				}
				result.Devices = append(result.Devices, dependencie)
			}
			for _, element := range lane.MultiLane.Elements {
				if element.MsgEvent != nil {
					result.Devices = append(result.Devices, dependencymodel.DeviceDependency{
						DeviceId: element.MsgEvent.Device.Id,
						Name:     element.MsgEvent.Device.Name,
						BpmnResources: []dependencymodel.BpmnResource{{
							Id: element.MsgEvent.BpmnElementId,
						}},
					})
					result.Events = append(result.Events, dependencymodel.EventDependency{
						EventId: element.MsgEvent.EventId,
						BpmnResources: []dependencymodel.BpmnResource{{
							Id: element.MsgEvent.BpmnElementId,
						}},
					})
				}
				if element.ReceiveTaskEvent != nil {
					result.Devices = append(result.Devices, dependencymodel.DeviceDependency{
						DeviceId: element.ReceiveTaskEvent.Device.Id,
						Name:     element.ReceiveTaskEvent.Device.Name,
						BpmnResources: []dependencymodel.BpmnResource{{
							Id: element.ReceiveTaskEvent.BpmnElementId,
						}},
					})
					result.Events = append(result.Events, dependencymodel.EventDependency{
						EventId: element.ReceiveTaskEvent.EventId,
						BpmnResources: []dependencymodel.BpmnResource{{
							Id: element.ReceiveTaskEvent.BpmnElementId,
						}},
					})
				}
			}
		}
	}

	for _, element := range deployment.Elements {
		if element.Task != nil {
			result.Devices = append(result.Devices, dependencymodel.DeviceDependency{
				DeviceId: element.Task.Selection.Device.Id,
				Name:     element.Task.Selection.Device.Name,
				BpmnResources: []dependencymodel.BpmnResource{{
					Id: element.Task.BpmnElementId,
				}},
			})
		}
		if element.MultiTask != nil {
			for _, selection := range element.MultiTask.Selections {
				result.Devices = append(result.Devices, dependencymodel.DeviceDependency{
					DeviceId: selection.Device.Id,
					Name:     selection.Device.Name,
					BpmnResources: []dependencymodel.BpmnResource{{
						Id: element.MultiTask.BpmnElementId,
					}},
				})
			}
		}
		if element.MsgEvent != nil {
			result.Devices = append(result.Devices, dependencymodel.DeviceDependency{
				DeviceId: element.MsgEvent.Device.Id,
				Name:     element.MsgEvent.Device.Name,
				BpmnResources: []dependencymodel.BpmnResource{{
					Id: element.MsgEvent.BpmnElementId,
				}},
			})
			result.Events = append(result.Events, dependencymodel.EventDependency{
				EventId: element.MsgEvent.EventId,
				BpmnResources: []dependencymodel.BpmnResource{{
					Id: element.MsgEvent.BpmnElementId,
				}},
			})
		}
		if element.ReceiveTaskEvent != nil {
			result.Devices = append(result.Devices, dependencymodel.DeviceDependency{
				DeviceId: element.ReceiveTaskEvent.Device.Id,
				Name:     element.ReceiveTaskEvent.Device.Name,
				BpmnResources: []dependencymodel.BpmnResource{{
					Id: element.ReceiveTaskEvent.BpmnElementId,
				}},
			})
			result.Events = append(result.Events, dependencymodel.EventDependency{
				EventId: element.ReceiveTaskEvent.EventId,
				BpmnResources: []dependencymodel.BpmnResource{{
					Id: element.ReceiveTaskEvent.BpmnElementId,
				}},
			})
		}
	}
	result.Devices = reduceDeviceDependencies(result.Devices)
	sort.Sort(DeviceDependenciesByDeviceId(result.Devices))
	sort.Sort(EventDependenciesByEventId(result.Events))
	return result, nil
}

func (this *Ctrl) deploymentToDependenciesV2(deployment deploymentmodel2.Deployment) (result dependencymodel.Dependencies, err error) {
	panic("not implemented") //TODO
}

func reduceDeviceDependencies(dependencies []dependencymodel.DeviceDependency) (result []dependencymodel.DeviceDependency) {
	name := map[string]string{}
	resources := map[string][]dependencymodel.BpmnResource{}
	for _, device := range dependencies {
		name[device.DeviceId] = device.Name
		resources[device.DeviceId] = append(resources[device.DeviceId], device.BpmnResources...)
	}
	for id, name := range name {
		resourceList := resources[id]
		sort.Sort(BpmnResourcesById(resourceList))
		result = append(result, dependencymodel.DeviceDependency{
			DeviceId:      id,
			Name:          name,
			BpmnResources: resourceList,
		})
	}
	return
}

type DeviceDependenciesByDeviceId []dependencymodel.DeviceDependency

func (a DeviceDependenciesByDeviceId) Len() int           { return len(a) }
func (a DeviceDependenciesByDeviceId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DeviceDependenciesByDeviceId) Less(i, j int) bool { return a[i].DeviceId < a[j].DeviceId }

type EventDependenciesByEventId []dependencymodel.EventDependency

func (a EventDependenciesByEventId) Len() int           { return len(a) }
func (a EventDependenciesByEventId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EventDependenciesByEventId) Less(i, j int) bool { return a[i].EventId < a[j].EventId }

type BpmnResourcesById []dependencymodel.BpmnResource

func (a BpmnResourcesById) Len() int           { return len(a) }
func (a BpmnResourcesById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BpmnResourcesById) Less(i, j int) bool { return a[i].Id < a[j].Id }
