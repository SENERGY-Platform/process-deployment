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

	if task.Selection.FilterCriteria.AspectId != nil {
		temp, err := this.aspectNodeProvider(token, *task.Selection.FilterCriteria.AspectId)
		if err != nil {
			this.conf.GetLogger().Error("unable to load aspect node", "aspectId", *task.Selection.FilterCriteria.AspectId, "error", err)
			return err
		}
		command.Aspect = &temp
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

func isControllingFunction(functionId string) bool {
	return strings.HasPrefix(functionId, devicemodel.CONTROLLING_FUNCTION_PREFIX)
}
