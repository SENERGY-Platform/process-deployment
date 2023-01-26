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
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"sort"
)

func (this *Ctrl) GetDependencies(token auth.Token, id string) (dependencymodel.Dependencies, error, int) {
	return this.db.GetDependencies(token.GetUserId(), id)
}

func (this *Ctrl) GetDependenciesList(token auth.Token, limit int, offset int) ([]dependencymodel.Dependencies, error, int) {
	return this.db.GetDependenciesList(token.GetUserId(), limit, offset)
}

func (this *Ctrl) GetSelectedDependencies(token auth.Token, ids []string) ([]dependencymodel.Dependencies, error, int) {
	return this.db.GetSelectedDependencies(token.GetUserId(), ids)
}

func (this *Ctrl) SaveDependencies(command messages.DeploymentCommand) error {
	var dependencies dependencymodel.Dependencies
	var err error
	dependencies, err = this.deploymentToDependencies(command.Owner, *command.Deployment)
	if err != nil {
		return err
	}
	dependencies.Owner = command.Owner
	return this.db.SetDependencies(dependencies)
}

func (this *Ctrl) DeleteDependencies(command messages.DeploymentCommand) error {
	return this.db.DeleteDependencies(command.Id)
}

func (this *Ctrl) deploymentToDependencies(user string, deployment deploymentmodel.Deployment) (result dependencymodel.Dependencies, err error) {
	result.DeploymentId = deployment.Id
	result.Events = []dependencymodel.EventDependency{}
	result.Devices = []dependencymodel.DeviceDependency{}
	for _, element := range deployment.Elements {
		if element.Task != nil && element.Task.Selection.SelectedDeviceId != nil {
			dependency := this.getDeviceDependencyFromSelection(user, element.Task.Selection)
			dependency.BpmnResources = []dependencymodel.BpmnResource{{
				Id: element.BpmnId,
				//Label: element.Name,
			}}
			result.Devices = append(result.Devices, dependency)
		}
		if element.MessageEvent != nil && element.MessageEvent.Selection.SelectedDeviceId != nil {
			dependency := this.getDeviceDependencyFromSelection(user, element.MessageEvent.Selection)
			dependency.BpmnResources = []dependencymodel.BpmnResource{{
				Id: element.BpmnId,
				//Label: element.Name,
			}}
			result.Devices = append(result.Devices, dependency)
			result.Events = append(result.Events, dependencymodel.EventDependency{
				EventId: element.MessageEvent.EventId,
				BpmnResources: []dependencymodel.BpmnResource{{
					Id: element.BpmnId,
					//Label: element.Name,
				}},
			})
		}
		if element.ConditionalEvent != nil && element.ConditionalEvent.Selection.SelectedDeviceId != nil {
			dependency := this.getDeviceDependencyFromSelection(user, element.ConditionalEvent.Selection)
			dependency.BpmnResources = []dependencymodel.BpmnResource{{
				Id: element.BpmnId,
				//Label: element.Name,
			}}
			result.Devices = append(result.Devices, dependency)
			result.Events = append(result.Events, dependencymodel.EventDependency{
				EventId: element.ConditionalEvent.EventId,
				BpmnResources: []dependencymodel.BpmnResource{{
					Id: element.BpmnId,
					//Label: element.Name,
				}},
			})
		}
	}
	return
}

func (this *Ctrl) getDeviceDependencyFromSelection(user string, selection deploymentmodel.Selection) (result dependencymodel.DeviceDependency) {
	if selection.SelectedDeviceId != nil {
		result.DeviceId = *selection.SelectedDeviceId
	}
	for _, option := range selection.SelectionOptions {
		if option.Device != nil && option.Device.Id == result.DeviceId {
			result.Name = option.Device.Name
			return
		}
	}
	if result.Name == "" && result.DeviceId != "" {
		token, err := auth.CreateToken("process-deployment", user)
		if err != nil {
			return result
		}
		device, err, _ := this.devices.GetDevice(token, result.DeviceId)
		if err != nil {
			return result
		}
		result.Name = device.Name
	}
	return
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
