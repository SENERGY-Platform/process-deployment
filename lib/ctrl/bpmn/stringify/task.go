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

package stringify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
	"text/template"
)

func Task(doc *etree.Document, task *model.Task, selectAsRef bool, deviceRepo interfaces.Devices) (err error) {
	if task == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	command := model.Command{
		Retries:          task.Retries,
		Function:         task.DeviceDescription.Function,
		CharacteristicId: task.DeviceDescription.CharacteristicId,
		DeviceClass:      task.DeviceDescription.DeviceClass,
		Aspect:           task.DeviceDescription.Aspect,
		Input:            task.Input,
		Configurables:    task.Configurables,
	}

	if selectAsRef {
		command.DeviceId = task.Selection.Device.Id
		command.ServiceId = task.Selection.Service.Id
		command.ProtocolId = task.Selection.Service.ProtocolId
	} else {
		command.Device = task.Selection.Device
		command.Service = task.Selection.Service
		protocol, err, _ := deviceRepo.GetProtocol(task.Selection.Service.ProtocolId)
		if err != nil {
			return err
		}
		command.Protocol = &protocol
	}

	commandStr, err := json.Marshal(command)

	if err != nil {
		return err
	}

	xpath := "//bpmn:serviceTask[@id='" + task.BpmnElementId + "']//camunda:inputParameter[@name='" + model.CAMUNDA_VARIABLES_PAYLOAD + "']"
	doc.FindElement(xpath).SetText(string(commandStr))

	for name, value := range task.Parameter {
		xpath := "//bpmn:serviceTask[@id='" + task.BpmnElementId + "']//camunda:inputParameter[@name='" + name + "']"
		doc.FindElement(xpath).SetText(value)
	}
	return nil
}

func LaneTask(doc *etree.Document, task *model.LaneTask, device *devicemodel.Device, selectAsRef bool, deviceRepo interfaces.Devices) (err error) {
	if task == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	command := model.Command{
		Retries:          task.Retries,
		Function:         task.DeviceDescription.Function,
		CharacteristicId: task.DeviceDescription.CharacteristicId,
		DeviceClass:      task.DeviceDescription.DeviceClass,
		Aspect:           task.DeviceDescription.Aspect,
		Input:            task.Input,
		Configurables:    task.Configurables,
	}

	if selectAsRef {
		command.DeviceId = device.Id
		command.ServiceId = task.SelectedService.Id
		command.ProtocolId = task.SelectedService.ProtocolId
	} else {
		command.Device = device
		command.Service = task.SelectedService
		protocol, err, _ := deviceRepo.GetProtocol(task.SelectedService.ProtocolId)
		if err != nil {
			return err
		}
		command.Protocol = &protocol
	}

	commandStr, err := json.Marshal(command)

	if err != nil {
		return err
	}

	xpath := "//bpmn:serviceTask[@id='" + task.BpmnElementId + "']//camunda:inputParameter[@name='" + model.CAMUNDA_VARIABLES_PAYLOAD + "']"
	doc.FindElement(xpath).SetText(string(commandStr))

	for name, value := range task.Parameter {
		xpath := "//bpmn:serviceTask[@id='" + task.BpmnElementId + "']//camunda:inputParameter[@name='" + name + "']"
		doc.FindElement(xpath).SetText(value)
	}
	return nil
}

func MultiTask(doc *etree.Document, task *model.MultiTask, selectAsRef bool, deviceRepo interfaces.Devices) (err error) {
	if task == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	command := model.Command{
		Retries:          task.Retries,
		Function:         task.DeviceDescription.Function,
		CharacteristicId: task.DeviceDescription.CharacteristicId,
		DeviceClass:      task.DeviceDescription.DeviceClass,
		Aspect:           task.DeviceDescription.Aspect,
		Input:            task.Input,
		Configurables:    task.Configurables,
	}

	commandStr, err := json.Marshal(command)
	if err != nil {
		return err
	}

	xpath := "//bpmn:serviceTask[@id='" + task.BpmnElementId + "']//camunda:inputParameter[@name='" + model.CAMUNDA_VARIABLES_PAYLOAD + "']"
	doc.FindElement(xpath).SetText(string(commandStr))

	for name, value := range task.Parameter {
		xpath := "//bpmn:serviceTask[@id='" + task.BpmnElementId + "']//camunda:inputParameter[@name='" + name + "']"
		doc.FindElement(xpath).SetText(value)
	}

	script, err := createOverwriteVariableScript(task.Selections, selectAsRef, deviceRepo)
	if err != nil {
		return err
	}

	loopElement := doc.FindElement("//bpmn:serviceTask[@id='" + task.BpmnElementId + "']/bpmn:multiInstanceLoopCharacteristics")
	loopElement.CreateAttr("camunda:collection", model.CAMUNDE_VARIABLES_OVERWRITE_COLLECTION)
	loopElement.CreateAttr("camunda:elementVariable", model.CAMUNDA_VARIABLES_OVERWRITE)

	scriptElement := doc.CreateElement("camunda:script")
	scriptElement.CreateAttr("scriptFormat", "groovy")
	scriptElement.SetText(script)

	executionListener := doc.CreateElement("camunda:executionListener")
	executionListener.CreateAttr("event", "start")
	executionListener.InsertChild(nil, scriptElement)
	doc.FindElement("//bpmn:serviceTask[@id='"+task.BpmnElementId+"']/bpmn:extensionElements").InsertChild(nil, executionListener)
	return nil
}

func LaneMultiTask(doc *etree.Document, task *model.LaneTask, devices []*devicemodel.Device, selectAsRef bool, deviceRepo interfaces.Devices) (err error) {
	if task == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	command := model.Command{
		Retries:          task.Retries,
		Function:         task.DeviceDescription.Function,
		CharacteristicId: task.DeviceDescription.CharacteristicId,
		DeviceClass:      task.DeviceDescription.DeviceClass,
		Aspect:           task.DeviceDescription.Aspect,
		Input:            task.Input,
		Configurables:    task.Configurables,
	}

	commandStr, err := json.Marshal(command)
	if err != nil {
		return err
	}

	xpath := "//bpmn:serviceTask[@id='" + task.BpmnElementId + "']//camunda:inputParameter[@name='" + model.CAMUNDA_VARIABLES_PAYLOAD + "']"
	doc.FindElement(xpath).SetText(string(commandStr))

	for name, value := range task.Parameter {
		xpath := "//bpmn:serviceTask[@id='" + task.BpmnElementId + "']//camunda:inputParameter[@name='" + name + "']"
		doc.FindElement(xpath).SetText(value)
	}

	overwrites := []model.Overwrite{}
	for _, device := range devices {
		overwrite := model.Overwrite{}
		if selectAsRef {
			protocol, err, _ := deviceRepo.GetProtocol(task.SelectedService.ProtocolId)
			if err != nil {
				return err
			}
			overwrite.Device = device
			overwrite.Service = task.SelectedService
			overwrite.Protocol = &protocol
		} else {
			overwrite.DeviceId = device.Id
			overwrite.ServiceId = task.SelectedService.Id
			overwrite.ProtocolId = task.SelectedService.ProtocolId
		}
		overwrites = append(overwrites, overwrite)
	}

	script, err := overwritesToScript(overwrites)
	if err != nil {
		return err
	}

	loopElement := doc.FindElement("//bpmn:serviceTask[@id='" + task.BpmnElementId + "']/bpmn:multiInstanceLoopCharacteristics")
	loopElement.CreateAttr("camunda:collection", model.CAMUNDE_VARIABLES_OVERWRITE_COLLECTION)
	loopElement.CreateAttr("camunda:elementVariable", model.CAMUNDA_VARIABLES_OVERWRITE)

	scriptElement := doc.CreateElement("camunda:script")
	scriptElement.CreateAttr("scriptFormat", "groovy")
	scriptElement.SetText(script)

	executionListener := doc.CreateElement("camunda:executionListener")
	executionListener.CreateAttr("event", "start")
	executionListener.InsertChild(nil, scriptElement)
	doc.FindElement("//bpmn:serviceTask[@id='"+task.BpmnElementId+"']/bpmn:extensionElements").InsertChild(nil, executionListener)
	return nil
}

func createOverwriteVariableScript(selections []model.Selection, selectAsRef bool, deviceRepo interfaces.Devices) (script string, err error) {
	overwrites := []model.Overwrite{}
	for _, selection := range selections {
		overwrite := model.Overwrite{}
		if selectAsRef {
			protocol, err, _ := deviceRepo.GetProtocol(selection.Service.ProtocolId)
			if err != nil {
				return "", err
			}
			overwrite.Device = selection.Device
			overwrite.Service = selection.Service
			overwrite.Protocol = &protocol
		} else {
			overwrite.DeviceId = selection.Device.Id
			overwrite.ServiceId = selection.Service.Id
			overwrite.ProtocolId = selection.Service.ProtocolId
		}
		overwrites = append(overwrites, overwrite)
	}
	return overwritesToScript(overwrites)
}

func overwritesToScript(overwrites []model.Overwrite) (script string, err error) {
	overwriteSelections := []string{}
	for _, overwrite := range overwrites {
		overwriteMsg, err := json.Marshal(overwrite)
		if err != nil {
			return "", err
		}
		overwriteSelections = append(overwriteSelections, string(overwriteMsg))
	}
	collection, err := json.Marshal(overwriteSelections)
	if err != nil {
		return "", err
	}
	templ, err := template.New("script").Parse(SCRIPT_TEMPLATE)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = templ.Execute(&buffer, map[string]string{
		"CollectionName": model.CAMUNDE_VARIABLES_OVERWRITE_COLLECTION,
		"Collection":     string(collection),
	})
	return buffer.String(), err
}

const SCRIPT_TEMPLATE = `execution.setVariable("{{.CollectionName}}", {{.Collection}})`
