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

package stringifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/executionmodel"
	"github.com/beevik/etree"
)

func (this *Stringifier) Task(doc *etree.Document, element deploymentmodel.Element, token auth.Token) (err error) {
	task := element.Task
	if task == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil && err == nil {
			this.conf.GetLogger().Error("recovered from panic", "error", r)
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	command := executionmodel.Task{
		Retries:     task.Retries,
		Version:     deploymentmodel.CurrentVersion,
		PreferEvent: task.PreferEvent,
	}

	if task.Selection.SelectedPath != nil {
		command.ConfigurablesV2 = task.Selection.SelectedPath.Configurables
		if isControllingFunction(task.Selection.SelectedPath.FunctionId) {
			command.InputPaths = []string{task.Selection.SelectedPath.Path}
		} else {
			command.OutputPath = task.Selection.SelectedPath.Path
		}
	}

	if task.Selection.FilterCriteria.CharacteristicId != nil {
		command.CharacteristicId = *task.Selection.FilterCriteria.CharacteristicId
	}

	if task.Selection.FilterCriteria.FunctionId != nil {
		command.Function = devicemodel.Function{Id: *task.Selection.FilterCriteria.FunctionId}
	}

	if task.Selection.FilterCriteria.DeviceClassId != nil {
		command.DeviceClass = &devicemodel.DeviceClass{Id: *task.Selection.FilterCriteria.DeviceClassId}
	}

	//both aspect spellings of the criteria are passed on as they were selected: Aspect is
	//deprecated and an alias for an Aspects list with a single element, and a command that
	//carries it alone stays readable for a worker that only knows the single field. The
	//worker folds the two on read, at its camunda boundary.
	if task.Selection.FilterCriteria.AspectId != nil {
		temp, err := this.aspectNodeProvider(token, *task.Selection.FilterCriteria.AspectId)
		if err != nil {
			this.conf.GetLogger().Error("unable to load aspect node", "aspectId", *task.Selection.FilterCriteria.AspectId, "error", err)
			return err
		}
		command.Aspect = &temp
	}

	command.Aspects, err = this.aspectNodes(token, task.Selection.FilterCriteria.AspectIds)
	if err != nil {
		return err
	}

	xpath := "//bpmn:serviceTask[@id='" + element.BpmnId + "']//camunda:inputParameter[@name='" + executionmodel.CAMUNDA_VARIABLES_PAYLOAD + "']"

	cmd := executionmodel.Task{}
	cmdPayload := doc.FindElement(xpath)
	err = json.Unmarshal([]byte(cmdPayload.Text()), &cmd)
	if err != nil {
		return err
	}

	command.Input = cmd.Input
	command.Output = cmd.Output

	if task.Selection.SelectedDeviceId != nil && *task.Selection.SelectedDeviceId != "" {
		command.DeviceId = *task.Selection.SelectedDeviceId
	}
	if task.Selection.SelectedServiceId != nil && *task.Selection.SelectedServiceId != "" {
		command.ServiceId = *task.Selection.SelectedServiceId
	}
	if task.Selection.SelectedDeviceGroupId != nil && *task.Selection.SelectedDeviceGroupId != "" {
		command.DeviceGroupId = *task.Selection.SelectedDeviceGroupId
	}

	commandStr, err := json.MarshalIndent(command, "", "\t")

	if err != nil {
		return err
	}

	doc.FindElement(xpath).SetCData(string(commandStr))

	for name, value := range task.Parameter {
		xpath := "//bpmn:serviceTask[@id='" + element.BpmnId + "']//camunda:inputParameter[@name='" + name + "']"
		doc.FindElement(xpath).SetText(value)
	}
	return nil
}

// aspectNodes resolves the aspects a criteria names, in a stable order, so that the same
// criteria always produces the same command payload. Sorting by id also puts the node the
// deprecated single field would carry first, which is how the platform picks it elsewhere.
func (this *Stringifier) aspectNodes(token auth.Token, aspectIds []string) (result []devicemodel.AspectNode, err error) {
	for _, aspectId := range slices.Sorted(slices.Values(aspectIds)) {
		if aspectId == "" {
			continue
		}
		node, err := this.aspectNodeProvider(token, aspectId)
		if err != nil {
			this.conf.GetLogger().Error("unable to load aspect node", "aspectId", aspectId, "error", err)
			return nil, err
		}
		result = append(result, node)
	}
	return result, nil
}

func isControllingFunction(functionId string) bool {
	return strings.HasPrefix(functionId, devicemodel.CONTROLLING_FUNCTION_PREFIX)
}
