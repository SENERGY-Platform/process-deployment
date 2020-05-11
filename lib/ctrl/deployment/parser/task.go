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

package parser

import (
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/executionmodel"
	"github.com/beevik/etree"
	"strings"
)

func init() {
	ElementParsers = append(ElementParsers, ElementParser{
		Is: func(this *Parser, element *etree.Element) bool {
			return this.isTask(element)
		},
		Parse: func(this *Parser, dom *etree.Element) (element deploymentmodel.Element, err error) {
			return this.getTask(dom)
		},
	})
}

func (this *Parser) isTask(element *etree.Element) bool {
	if element.Tag != "serviceTask" && element.Tag != "bpmn:serviceTask" {
		return false
	}
	topic := element.SelectAttr("camunda:topic")
	if topic == nil || topic.Value == "" {
		return false
	}
	return true
}

func (this *Parser) getTask(dom *etree.Element) (result deploymentmodel.Element, err error) {
	cmd := executionmodel.Task{}
	cmdPayload := dom.FindElement(".//camunda:inputParameter[@name='" + executionmodel.CAMUNDA_VARIABLES_PAYLOAD + "']")
	err = json.Unmarshal([]byte(cmdPayload.Text()), &cmd)
	if err != nil {
		return result, err
	}

	parameter, err := this.getParameter(dom)
	if err != nil {
		return result, err
	}

	id := dom.SelectAttr("id").Value
	label := dom.SelectAttrValue("name", id)

	filterCriteria := deploymentmodel.FilterCriteria{
		CharacteristicId: &cmd.CharacteristicId,
		FunctionId:       &cmd.Function.Id,
	}
	if cmd.Aspect != nil {
		filterCriteria.AspectId = &cmd.Aspect.Id
	}
	if cmd.DeviceClass != nil {
		filterCriteria.DeviceClassId = &cmd.DeviceClass.Id
	}

	result = deploymentmodel.Element{
		BaseInfo: deploymentmodel.BaseInfo{
			Name:   label,
			BpmnId: id,
			Order:  this.getOrder(dom),
		},
		Task: &deploymentmodel.Task{
			Retries:        cmd.Retries,
			Parameter:      parameter,
			FilterCriteria: filterCriteria,
		},
	}
	return result, nil
}

func (this *Parser) getParameter(task *etree.Element) (result map[string]string, err error) {
	result = map[string]string{}
	for _, input := range task.FindElements(".//camunda:inputParameter") {
		name := input.SelectAttr("name").Value
		value := input.Text()
		if strings.HasPrefix(name, "inputs") {
			result[name] = value
		}
	}
	return result, err
}
