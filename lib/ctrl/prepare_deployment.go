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
	"net/http"
	"slices"
	"sort"

	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/importmodel"
)

func (this *Ctrl) PrepareDeployment(token auth.Token, xml string, svg string, withOptions bool) (result deploymentmodel.Deployment, err error, code int) {
	err = deploymentmodel.DeploymentXmlValidator(xml)
	if err != nil {
		return result, err, http.StatusBadRequest
	}
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
	result.IncidentHandling = &deploymentmodel.IncidentHandling{
		Restart: false,
		Notify:  true,
	}
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

func (this *Ctrl) GetDeployments(token auth.Token, options model.DeploymentListOptions) (result []deploymentmodel.Deployment, err error, code int) {
	result, err = this.db.ListDeployments(token.GetUserId(), options)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

func (this *Ctrl) CreateDeployment(token auth.Token, deployment deploymentmodel.Deployment, source string, optionals map[string]bool) (result deploymentmodel.Deployment, err error, code int) {
	deployment.Id = config.NewId()
	return this.SetDeployment(token, deployment, source, optionals)
}

func (this *Ctrl) UpdateDeployment(token auth.Token, id string, deployment deploymentmodel.Deployment, source string, optionals map[string]bool) (result deploymentmodel.Deployment, err error, code int) {
	if id != deployment.Id {
		return deployment, errors.New("path id != body id"), http.StatusBadRequest
	}

	err, code = this.db.CheckDeploymentAccess(token.GetUserId(), id)
	if err != nil {
		return result, err, code
	}

	return this.SetDeployment(token, deployment, source, optionals)
}

func (this *Ctrl) RemoveDeployment(token auth.Token, id string) (err error, code int) {
	err, code = this.db.CheckDeploymentAccess(token.GetUserId(), id)
	if err != nil {
		return err, code
	}
	err = this.deleteDeployment(id)
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
	bulk := this.getDeploymentBulkSelectableRequestV2(deployment)
	bulkResult, err, _ := this.devices.GetBulkDeviceSelectionV2(token, bulk)
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
		if element.ConditionalEvent != nil {
			selectable := selectableIndex[element.BpmnId]
			element.ConditionalEvent.Selection.SelectionOptions = removeConfigurables(getSelectionOptions(selectable, element.ConditionalEvent.Selection.FilterCriteria))
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

func (this *Ctrl) getDeploymentBulkSelectableRequestV2(deployment *deploymentmodel.Deployment) (bulk deviceselectionmodel.BulkRequestV2) {
	taskGroups := map[string][]int{}
	for index, element := range deployment.Elements {
		if element.Task != nil {
			if element.Group == nil {
				bulk = append(bulk, deviceselectionmodel.BulkRequestElementV2{
					Id: element.BpmnId,
					Criteria: deviceselectionmodel.FilterCriteriaAndSet{
						selectionCriteria(element.Task.Selection.FilterCriteria, ""),
					},
					IncludeGroups:            this.config.EnableDeviceGroupsForTasks,
					IncludeDevices:           true,
					IncludeIdModifiedDevices: this.config.EnableModifiedDevicesForDeploymentOptions,
				})
			} else {
				taskGroups[*element.Group] = append(taskGroups[*element.Group], index)
			}
		}
		if element.MessageEvent != nil {
			bulk = append(bulk, deviceselectionmodel.BulkRequestElementV2{
				Id: element.BpmnId,
				Criteria: deviceselectionmodel.FilterCriteriaAndSet{
					selectionCriteria(element.MessageEvent.Selection.FilterCriteria, devicemodel.EVENT),
				},
				IncludeGroups:            this.config.EnableDeviceGroupsForEvents,
				IncludeImports:           this.config.EnableImportsForEvents,
				IncludeDevices:           true,
				IncludeIdModifiedDevices: this.config.EnableModifiedDevicesForDeploymentOptions,
			})
		}
		if element.ConditionalEvent != nil {
			bulk = append(bulk, deviceselectionmodel.BulkRequestElementV2{
				Id: element.BpmnId,
				Criteria: deviceselectionmodel.FilterCriteriaAndSet{
					selectionCriteria(element.ConditionalEvent.Selection.FilterCriteria, devicemodel.EVENT),
				},
				IncludeGroups:            this.config.EnableDeviceGroupsForEvents,
				IncludeImports:           this.config.EnableImportsForEvents,
				IncludeDevices:           true,
				IncludeIdModifiedDevices: this.config.EnableModifiedDevicesForDeploymentOptions,
			})
		}
	}

	for _, indexes := range taskGroups {
		filter := deviceselectionmodel.FilterCriteriaAndSet{}
		for _, index := range indexes {
			element := deployment.Elements[index]
			if element.Task != nil {
				filter = append(filter, selectionCriteria(element.Task.Selection.FilterCriteria, ""))
			}
		}
		for _, index := range indexes {
			element := deployment.Elements[index]
			if element.Task != nil {
				bulk = append(bulk, deviceselectionmodel.BulkRequestElementV2{
					Id:                       element.BpmnId,
					Criteria:                 filter,
					IncludeGroups:            this.config.EnableDeviceGroupsForTasks,
					IncludeDevices:           true,
					IncludeIdModifiedDevices: this.config.EnableModifiedDevicesForDeploymentOptions,
				})
			}
		}
	}
	return bulk
}

// selectionCriteria translates a deployment criteria into the shape the selection service is
// asked with. Both aspect spellings are passed on as they were selected: AspectId is
// deprecated and an alias for an AspectIds list with a single element, and the selection
// service folds the two at its own boundary.
func selectionCriteria(criteria deploymentmodel.FilterCriteria, interaction devicemodel.Interaction) deviceselectionmodel.FilterCriteria {
	result := deviceselectionmodel.FilterCriteria{Interaction: string(interaction)}
	if criteria.FunctionId != nil {
		result.FunctionId = *criteria.FunctionId
	}
	if criteria.DeviceClassId != nil {
		result.DeviceClassId = *criteria.DeviceClassId
	}
	if criteria.AspectId != nil {
		result.AspectId = *criteria.AspectId
	}
	result.AspectIds = criteria.AspectIds
	return result
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
			//the answer and the deployment model hold the same import types since
			//device-selection reads its import shapes from models/go, so the option takes
			//them as they arrived. A cast layer here would silently drop every field it
			//does not name - that is how Cost went missing before.
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
	sort.Slice(result, func(i, j int) bool {
		nameI := ""
		if result[i].Device != nil {
			nameI = result[i].Device.Name
		} else if result[i].DeviceGroup != nil {
			nameI = result[i].DeviceGroup.Name
		} else if result[i].Import != nil {
			nameI = result[i].Import.Name
		}
		nameJ := ""
		if result[j].Device != nil {
			nameJ = result[j].Device.Name
		} else if result[j].DeviceGroup != nil {
			nameJ = result[j].DeviceGroup.Name
		} else if result[j].Import != nil {
			nameJ = result[j].Import.Name
		}
		return nameI < nameJ
	})
	return result
}

// serviceMatchesCriteria answers whether a service of a selectable offers what the element
// asks for. Several aspects in one criteria are an AND, the way the device-repository reads
// them: one content variable, and therefore one path option, has to carry all of them.
func serviceMatchesCriteria(service devicemodel.Service, criteria deploymentmodel.FilterCriteria, servicePathOptions map[string][]deviceselectionmodel.PathOption) bool {
	aspectIds := criteriaAspectIds(criteria)
	implementsFunction := false
	matchesAspects := false
	pathOptions, ok := servicePathOptions[service.Id]
	if !ok {
		return false
	}
	for _, option := range pathOptions {
		if optionCoversAspects(option, aspectIds) {
			matchesAspects = true
		}
		if criteria.FunctionId != nil && *criteria.FunctionId == option.FunctionId {
			implementsFunction = true
		}
		if (len(aspectIds) == 0 || matchesAspects) && (criteria.FunctionId == nil || implementsFunction) {
			return true
		}
	}
	return (len(aspectIds) == 0 || matchesAspects) && (criteria.FunctionId == nil || implementsFunction)
}

// criteriaAspectIds folds the deprecated single aspect of a criteria into its aspect list.
// AspectId is an alias for a list with one element, so a criteria may carry either spelling
// and everything behind this point evaluates the list only. An empty aspect is not a filter.
func criteriaAspectIds(criteria deploymentmodel.FilterCriteria) (result []string) {
	result = slices.Clone(criteria.AspectIds)
	if criteria.AspectId != nil && *criteria.AspectId != "" && !slices.Contains(result, *criteria.AspectId) {
		result = append(result, *criteria.AspectId)
	}
	return result
}

// optionCoversAspects checks the AND: every aspect the criteria names has to be the aspect of
// one of the nodes the path option offers, or an ancestor of it, because a queried aspect
// covers its own subtree.
func optionCoversAspects(option deviceselectionmodel.PathOption, aspectIds []string) bool {
	nodes := pathOptionAspectNodes(option)
	for _, aspectId := range aspectIds {
		if !slices.ContainsFunc(nodes, func(node devicemodel.AspectNode) bool {
			return node.Id == aspectId || slices.Contains(node.AncestorIds, aspectId)
		}) {
			return false
		}
	}
	return true
}

// pathOptionAspectNodes returns the aspect nodes a path option offers. AspectNode is
// deprecated and an alias for a single element AspectNodes, so the list is preferred and the
// single node only used when it is empty - a device-selection that predates the list fills
// the single field alone, while a current one fills both.
func pathOptionAspectNodes(option deviceselectionmodel.PathOption) []devicemodel.AspectNode {
	if len(option.AspectNodes) > 0 {
		return option.AspectNodes
	}
	if option.AspectNode.Id == "" {
		return nil
	}
	return []devicemodel.AspectNode{option.AspectNode}
}

func (this *Ctrl) SetDeployment(token auth.Token, deployment deploymentmodel.Deployment, source string, optionals map[string]bool) (result deploymentmodel.Deployment, err error, code int) {
	if err := deployment.Validate(deploymentmodel.ValidateRequest, optionals, deploymentmodel.DeploymentXmlValidator); err != nil {
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

	if err := deployment.Validate(deploymentmodel.ValidatePublish, optionals, deploymentmodel.DeploymentXmlValidator); err != nil {
		return deployment, err, http.StatusBadRequest
	}

	err = this.setDeployment(userid, source, deployment)
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}
	return deployment, nil, 200
}

// ensures selection correctness
func (this *Ctrl) EnsureDeploymentSelectionAccess(token auth.Token, deployment *deploymentmodel.Deployment) (err error, code int) {
	deviceIds := []string{}
	deviceGroupIds := []string{}
	importIds := []string{}
	for _, element := range deployment.Elements {
		if element.Task != nil && element.Task.Selection.SelectedDeviceId != nil {
			deviceIds = append(deviceIds, *element.Task.Selection.SelectedDeviceId)
		}
		if element.Task != nil && element.Task.Selection.SelectedDeviceGroupId != nil {
			deviceGroupIds = append(deviceGroupIds, *element.Task.Selection.SelectedDeviceGroupId)
		}
		if element.MessageEvent != nil && element.MessageEvent.Selection.SelectedDeviceId != nil {
			deviceIds = append(deviceIds, *element.MessageEvent.Selection.SelectedDeviceId)
		}
		if element.MessageEvent != nil && element.MessageEvent.Selection.SelectedDeviceGroupId != nil {
			deviceGroupIds = append(deviceGroupIds, *element.MessageEvent.Selection.SelectedDeviceGroupId)
		}
		if element.MessageEvent != nil && element.MessageEvent.Selection.SelectedImportId != nil {
			importIds = append(importIds, *element.MessageEvent.Selection.SelectedImportId)
		}
		if element.ConditionalEvent != nil && element.ConditionalEvent.Selection.SelectedDeviceId != nil {
			deviceIds = append(deviceIds, *element.ConditionalEvent.Selection.SelectedDeviceId)
		}
		if element.ConditionalEvent != nil && element.ConditionalEvent.Selection.SelectedDeviceGroupId != nil {
			deviceGroupIds = append(deviceGroupIds, *element.ConditionalEvent.Selection.SelectedDeviceGroupId)
		}
		if element.ConditionalEvent != nil && element.ConditionalEvent.Selection.SelectedImportId != nil {
			importIds = append(importIds, *element.ConditionalEvent.Selection.SelectedImportId)
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

	if len(importIds) > 0 {
		importaccess, err := this.imports.CheckAccess(token, importIds, false)
		if err != nil {
			return err, http.StatusInternalServerError
		}
		for !importaccess {
			return errors.New("import access denied"), http.StatusForbidden
		}
	}
	return nil, http.StatusOK
}
