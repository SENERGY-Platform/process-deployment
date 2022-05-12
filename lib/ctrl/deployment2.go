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
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/importmodel"
	"net/http"
)

func (this *Ctrl) PrepareDeployment(token auth.Token, xml string, svg string, withOptions bool) (result deploymentmodel.Deployment, err error, code int) {
	result, err = this.deploymentParser.PrepareDeployment(xml)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	if withOptions {
		err = this.SetDeploymentOptions(token, &result)
		if err != nil {
			return result, err, http.StatusInternalServerError
		}
	}
	result.Diagram.Svg = svg
	this.SetExecutableFlag(&result)
	return result, nil, http.StatusOK
}

func (this *Ctrl) GetDeployment(token auth.Token, id string, withOptions bool) (result deploymentmodel.Deployment, err error, code int) {
	temp, err, code := this.db.GetDeployment(token.GetUserId(), id)
	if err != nil {
		return result, err, code
	}
	if temp == nil {
		return result, errors.New("found deployment is not of requested version"), http.StatusBadRequest
	}
	result = *temp
	if withOptions {
		err = this.SetDeploymentOptions(token, &result)
		if err != nil {
			return result, err, http.StatusInternalServerError
		}
	}
	return
}

func (this *Ctrl) CreateDeployment(token auth.Token, deployment deploymentmodel.Deployment, source string) (result deploymentmodel.Deployment, err error, code int) {
	deployment.Id = config.NewId()
	return this.setDeployment(token, deployment, source)
}

func (this *Ctrl) UpdateDeployment(token auth.Token, id string, deployment deploymentmodel.Deployment, source string) (result deploymentmodel.Deployment, err error, code int) {
	if id != deployment.Id {
		return deployment, errors.New("path id != body id"), http.StatusBadRequest
	}

	err, code = this.db.CheckDeploymentAccess(token.GetUserId(), id)
	if err != nil {
		return result, err, code
	}

	return this.setDeployment(token, deployment, source)
}

func (this *Ctrl) RemoveDeployment(token auth.Token, id string) (err error, code int) {
	err, code = this.db.CheckDeploymentAccess(token.GetUserId(), id)
	if err != nil {
		return err, code
	}

	err = this.publishDeploymentDelete(token.GetUserId(), id)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, 200
}

func (this *Ctrl) SetExecutableFlag(deployment *deploymentmodel.Deployment) {
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

func (this *Ctrl) SetDeploymentOptions(token auth.Token, deployment *deploymentmodel.Deployment) (err error) {
	bulk := this.getDeploymentBulkSelectableRequest(deployment)
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
			element.MessageEvent.Selection.SelectionOptions = removeConfigurables(getSelectionOptions(selectable, element.MessageEvent.Selection.FilterCriteria))
		}
		deployment.Elements[index] = element
	}
	return nil
}

func removeConfigurables(options []deploymentmodel.SelectionOption) (result []deploymentmodel.SelectionOption) {
	for _, option := range options {
		optionCopy := option
		newPathOptions := map[string][]deviceselectionmodel.PathOption{}
		for serviceId, pathOptionList := range optionCopy.PathOptions {
			for _, pathOption := range pathOptionList {
				pathOptionCopy := pathOption
				pathOptionCopy.Configurables = nil
				newPathOptions[serviceId] = append(newPathOptions[serviceId], pathOptionCopy)
			}
		}
		optionCopy.PathOptions = newPathOptions
		result = append(result, optionCopy)
	}
	return
}

func (this *Ctrl) getDeploymentBulkSelectableRequest(deployment *deploymentmodel.Deployment) (bulk deviceselectionmodel.BulkRequest) {
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
				Id:             element.BpmnId,
				Criteria:       deviceselectionmodel.FilterCriteriaAndSet{element.MessageEvent.Selection.FilterCriteria.ToDeviceTypeFilter()},
				IncludeGroups:  this.config.EnableDeviceGroupsForEvents,
				IncludeImports: this.config.EnableImportsForEvents,
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
		var selectableImport *importmodel.Import
		var importType *importmodel.ImportType
		if selectable.DeviceGroup != nil {
			devicegroup = &deploymentmodel.DeviceGroup{
				Id:   selectable.DeviceGroup.Id,
				Name: selectable.DeviceGroup.Name,
			}
		}
		if selectable.Device != nil {
			for _, service := range selectable.Services {
				if serviceMatchesCriteria(service, criteria, selectable.ServicePathOptions) {
					serviceDesc = append(serviceDesc, deploymentmodel.Service{
						Id:   service.Id,
						Name: service.Name,
					})
				}
			}
			devicename := selectable.Device.DisplayName
			if devicename == "" {
				devicename = selectable.Device.Name
			}
			device = &deploymentmodel.Device{
				Id:   selectable.Device.Id,
				Name: devicename,
			}
		}
		if selectable.Import != nil && selectable.ImportType != nil {
			selectableImport = selectable.Import
			importType = selectable.ImportType
		}

		result = append(result, deploymentmodel.SelectionOption{
			Device:      device,
			DeviceGroup: devicegroup,
			Services:    serviceDesc,
			Import:      selectableImport,
			ImportType:  importType,
			PathOptions: selectable.ServicePathOptions,
		})
	}
	return result
}

func serviceMatchesCriteria(service devicemodel.Service, criteria deploymentmodel.FilterCriteria, servicePathOptions map[string][]deviceselectionmodel.PathOption) bool {
	implementsFunction := false
	matchesAspect := false
	pathOptions, ok := servicePathOptions[service.Id]
	if !ok {
		return false
	}
	for _, option := range pathOptions {
		aspects := append(option.AspectNode.AncestorIds, option.AspectNode.Id)
		for _, aspect := range aspects {
			if criteria.AspectId != nil && *criteria.AspectId == aspect {
				matchesAspect = true
				break
			}
		}
		if criteria.FunctionId != nil && *criteria.FunctionId == option.FunctionId {
			implementsFunction = true
		}
		if (criteria.AspectId == nil || matchesAspect) && (criteria.FunctionId == nil || implementsFunction) {
			return true
		}
	}
	return (criteria.AspectId == nil || matchesAspect) && (criteria.FunctionId == nil || implementsFunction)
}

func (this *Ctrl) setDeployment(token auth.Token, deployment deploymentmodel.Deployment, source string) (result deploymentmodel.Deployment, err error, code int) {
	if err := deployment.Validate(deploymentmodel.ValidateRequest); err != nil {
		return deployment, err, http.StatusBadRequest
	}

	//ensure selected devices and services exist and have the given content and are executable for the requesting user (if not using id ref)
	err, code = this.EnsureDeploymentSelectionAccess(token, &deployment)
	if err != nil {
		return deployment, err, code
	}

	err = this.completeEvents(&deployment)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	userid := token.GetUserId()

	deployment.Diagram.XmlDeployed, err = this.deploymentStringifier.Deployment(deployment, userid, token)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}

	if err = this.publishDeployment(userid, deployment.Id, deployment, source); err != nil {
		return deployment, err, http.StatusInternalServerError
	}
	return deployment, nil, 200
}

//ensures selection correctness
func (this *Ctrl) EnsureDeploymentSelectionAccess(token auth.Token, deployment *deploymentmodel.Deployment) (err error, code int) {
	deviceIds := []string{}
	deviceGroupIds := []string{}
	importIds := []string{}
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
		if element.MessageEvent != nil && element.MessageEvent.Selection.SelectedImportId != nil {
			importIds = append(importIds, *element.MessageEvent.Selection.SelectedImportId)
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

	importaccess, err := this.imports.CheckAccess(token, importIds, false)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	for !importaccess {
		return errors.New("import access denied"), http.StatusForbidden
	}
	return nil, http.StatusOK
}
