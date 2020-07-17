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
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/executionmodel"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
)

func (this *Stringifier) Task(doc *etree.Document, element deploymentmodel.Element) (err error) {
	task := element.Task
	if task == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	command := executionmodel.Task{
		Retries: task.Retries,
	}

	if task.Selection.FilterCriteria.FunctionId != nil {
		command.Function = devicemodel.Function{Id: *task.Selection.FilterCriteria.FunctionId}
	}

	if task.Selection.FilterCriteria.CharacteristicId != nil {
		command.CharacteristicId = *task.Selection.FilterCriteria.CharacteristicId
	}

	if task.Selection.FilterCriteria.DeviceClassId != nil {
		command.DeviceClass = &devicemodel.DeviceClass{Id: *task.Selection.FilterCriteria.DeviceClassId}
	}

	if task.Selection.FilterCriteria.AspectId != nil {
		command.Aspect = &devicemodel.Aspect{Id: *task.Selection.FilterCriteria.AspectId}
	}

	command.Configurables = task.Configurables

	xpath := "//bpmn:serviceTask[@id='" + element.BpmnId + "']//camunda:inputParameter[@name='" + executionmodel.CAMUNDA_VARIABLES_PAYLOAD + "']"

	cmd := executionmodel.Task{}
	cmdPayload := doc.FindElement(xpath)
	err = json.Unmarshal([]byte(cmdPayload.Text()), &cmd)
	if err != nil {
		return err
	}

	command.Input = cmd.Input
	command.Output = cmd.Output

	command.DeviceId = task.Selection.SelectedDeviceId
	command.ServiceId = task.Selection.SelectedServiceId

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
