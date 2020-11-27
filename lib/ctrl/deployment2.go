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
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"
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

func (this *Ctrl) CreateDeploymentV2(jwt jwt_http_router.Jwt, deployment deploymentmodel.Deployment, source string) (result deploymentmodel.Deployment, err error, code int) {
	deployment.Id = config.NewId()
	return this.setDeploymentV2(jwt, deployment, source)
}

func (this *Ctrl) UpdateDeploymentV2(jwt jwt_http_router.Jwt, id string, deployment deploymentmodel.Deployment, source string) (result deploymentmodel.Deployment, err error, code int) {
	if id != deployment.Id {
		return deployment, errors.New("path id != body id"), http.StatusBadRequest
	}

	err, code = this.db.CheckDeploymentAccess(jwt.UserId, id)
	if err != nil {
		return result, err, code
	}

	return this.setDeploymentV2(jwt, deployment, source)
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
		if element.Task != nil && len(element.Task.Selection.SelectionOptions) == 0 {
			deployment.Executable = false
			return
		}
		if element.MessageEvent != nil && len(element.MessageEvent.Selection.SelectionOptions) == 0 {
			deployment.Executable = false
			return
		}
	}
	return
}

func (this *Ctrl) SetDeploymentOptionsV2(token jwt_http_router.JwtImpersonate, deployment *deploymentmodel.Deployment) (err error) {
	bulk := this.getDeploymentV2BulkSelectableRequest(deployment)
	bulkResult, err, _ := this.devices.GetBulkDeviceSelection(token, bulk)
	if err != nil {
		return err
	}
	selectableIndex := map[string][]deviceselectionmodel.Selectable{}
	for _, element := range bulkResult {
		selectableIndex[element.Id] = element.Selectables
	}

	for index, element := range deployment.Elements {
		if element.Task != nil {
			selectable := selectableIndex[element.BpmnId]
			element.Task.Selection.SelectionOptions = getSelectionOptions(selectable, element.Task.Selection.FilterCriteria)
		}
		if element.MessageEvent != nil {
			selectable := selectableIndex[element.BpmnId]
			element.MessageEvent.Selection.SelectionOptions = getSelectionOptions(selectable, element.MessageEvent.Selection.FilterCriteria)
		}
		deployment.Elements[index] = element
	}
	return nil
}

func (this *Ctrl) getDeploymentV2BulkSelectableRequest(deployment *deploymentmodel.Deployment) (bulk deviceselectionmodel.BulkRequest) {
	useEventFilter := devicemodel.EVENT
	taskGroups := map[string][]int{}
	for index, element := range deployment.Elements {
		if element.Task != nil {
			if element.Group == nil {
				bulk = append(bulk, deviceselectionmodel.BulkRequestElement{
					Id:                element.BpmnId,
					FilterInteraction: &useEventFilter,
					Criteria:          deviceselectionmodel.FilterCriteriaAndSet{element.Task.Selection.FilterCriteria.ToDeviceTypeFilter()},
					IncludeGroups:     this.config.EnableDeviceGroupsForTasks,
				})
			} else {
				taskGroups[*element.Group] = append(taskGroups[*element.Group], index)
			}
		}
		if element.MessageEvent != nil {
			bulk = append(bulk, deviceselectionmodel.BulkRequestElement{
				Id:            element.BpmnId,
				Criteria:      deviceselectionmodel.FilterCriteriaAndSet{element.MessageEvent.Selection.FilterCriteria.ToDeviceTypeFilter()},
				IncludeGroups: this.config.EnableDeviceGroupsForEvents,
			})
		}
	}

	for _, indexes := range taskGroups {
		filter := deviceselectionmodel.FilterCriteriaAndSet{}
		for _, index := range indexes {
			element := deployment.Elements[index]
			if element.Task != nil {
				filter = append(filter, element.Task.Selection.FilterCriteria.ToDeviceTypeFilter())
			}
		}
		for _, index := range indexes {
			element := deployment.Elements[index]
			if element.Task != nil {
				bulk = append(bulk, deviceselectionmodel.BulkRequestElement{
					Id:                element.BpmnId,
					FilterInteraction: &useEventFilter,
					Criteria:          filter,
				})
			}
		}
	}
	return bulk
}

func getSelectionOptions(selectables []deviceselectionmodel.Selectable, criteria deploymentmodel.FilterCriteria) (result []deploymentmodel.SelectionOption) {
	for _, selectable := range selectables {
		serviceDesc := []deploymentmodel.Service{}
		var device *deploymentmodel.Device
		var devicegroup *deploymentmodel.DeviceGroup
		if selectable.DeviceGroup != nil {
			devicegroup = &deploymentmodel.DeviceGroup{
				Id:   selectable.DeviceGroup.Id,
				Name: selectable.DeviceGroup.Name,
			}
		}
		if selectable.Device != nil {
			for _, service := range selectable.Services {
				if serviceMatchesCriteria(service, criteria) {
					serviceDesc = append(serviceDesc, deploymentmodel.Service{
						Id:   service.Id,
						Name: service.Name,
					})
				}
			}
			device = &deploymentmodel.Device{
				Id:   selectable.Device.Id,
				Name: selectable.Device.Name,
			}
		}

		result = append(result, deploymentmodel.SelectionOption{
			Device:      device,
			DeviceGroup: devicegroup,
			Services:    serviceDesc,
		})
	}
	return result
}

func serviceMatchesCriteria(service devicemodel.Service, criteria deploymentmodel.FilterCriteria) bool {
	implementsFunction := false
	matchesAspect := false
	for _, function := range service.FunctionIds {
		if criteria.FunctionId != nil && *criteria.FunctionId == function {
			implementsFunction = true
			break
		}
	}
	for _, aspect := range service.AspectIds {
		if criteria.AspectId != nil && *criteria.AspectId == aspect {
			matchesAspect = true
			break
		}
	}
	return (criteria.AspectId == nil || matchesAspect) && (criteria.FunctionId == nil || implementsFunction)
}

func (this *Ctrl) setDeploymentV2(jwt jwt_http_router.Jwt, deployment deploymentmodel.Deployment, source string) (result deploymentmodel.Deployment, err error, code int) {
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

	if err = this.publishDeploymentV2(jwt.UserId, deployment.Id, deployment, source); err != nil {
		return deployment, err, http.StatusInternalServerError
	}
	return deployment, nil, 200
}

//ensures selection correctness
func (this *Ctrl) ensureDeploymentSelectionAccess(token jwt_http_router.JwtImpersonate, deployment *deploymentmodel.Deployment) (err error, code int) {
	deviceIds := []string{}
	deviceGroupIds := []string{}
	for _, element := range deployment.Elements {
		if element.Task != nil && element.Task.Selection.SelectedDeviceId != nil {
			deviceIds = append(deviceIds, *element.Task.Selection.SelectedDeviceId)
		}
		if element.MessageEvent != nil && element.MessageEvent.Selection.SelectedDeviceId != nil {
			deviceIds = append(deviceIds, *element.MessageEvent.Selection.SelectedDeviceId)
		}
		if element.Task != nil && element.Task.Selection.SelectedDeviceGroupId != nil {
			deviceGroupIds = append(deviceGroupIds, *element.Task.Selection.SelectedDeviceGroupId)
		}
		if element.MessageEvent != nil && element.MessageEvent.Selection.SelectedDeviceGroupId != nil {
			deviceGroupIds = append(deviceGroupIds, *element.MessageEvent.Selection.SelectedDeviceGroupId)
		}
	}

	deviceaccess, err := this.devices.CheckAccess(token, "devices", deviceIds)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	for _, access := range deviceaccess {
		if !access {
			return errors.New("device access denied"), http.StatusForbidden
		}
	}

	devicegroupaccess, err := this.devices.CheckAccess(token, "device-groups", deviceGroupIds)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	for _, access := range devicegroupaccess {
		if !access {
			return errors.New("device-groupaccess denied"), http.StatusForbidden
		}
	}
	return nil, http.StatusOK
}
