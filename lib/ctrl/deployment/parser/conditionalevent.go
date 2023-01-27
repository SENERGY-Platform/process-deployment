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
	"github.com/beevik/etree"
	"strconv"
)

func init() {
	ElementParsers = append(ElementParsers, ElementParser{
		Is: func(this *Parser, element *etree.Element) bool {
			return this.isConditionalEvent(element)
		},
		Parse: func(this *Parser, dom *etree.Element) (element deploymentmodel.Element, err error) {
			return this.getConditionalEvent(dom)
		},
	})
}

func (this *Parser) isConditionalEvent(element *etree.Element) bool {
	msgEvent := element.FindElement(".//bpmn:messageEventDefinition")
	if msgEvent == nil {
		return false
	}
	if msgEvent.SelectAttrValue("messageRef", "") != "" {
		return false
	}
	aspect := element.SelectAttr("aspect")
	if aspect == nil || aspect.Value == "" {
		return false
	}

	function := element.SelectAttr("function")
	if function == nil || function.Value == "" {
		return false
	}

	characteristic := element.SelectAttr("characteristic")
	if characteristic == nil || characteristic.Value == "" {
		return false
	}

	script := element.SelectAttr("script")
	if script == nil || script.Value == "" {
		return false
	}

	return true
}

func (this *Parser) getConditionalEvent(element *etree.Element) (result deploymentmodel.Element, err error) {
	id := element.SelectAttr("id").Value
	label := element.SelectAttrValue("name", id)

	filterCriteria := deploymentmodel.FilterCriteria{}

	aspect := element.SelectAttr("aspect")
	if aspect != nil {
		filterCriteria.AspectId = &aspect.Value
	}

	function := element.SelectAttr("function")
	if function != nil {
		filterCriteria.FunctionId = &function.Value
	}

	characteristic := element.SelectAttr("characteristic")
	if characteristic != nil {
		filterCriteria.CharacteristicId = &characteristic.Value
	}

	scriptAttr := element.SelectAttr("script")
	var script string
	if scriptAttr != nil {
		script = scriptAttr.Value
	}

	valueVariableName := "value"
	valueVariableNameAttr := element.SelectAttr("value_variable_name")
	if valueVariableNameAttr != nil {
		valueVariableName = valueVariableNameAttr.Value
	}

	variables := map[string]string{}
	variablesAttr := element.SelectAttr("variables")
	if variablesAttr != nil && variablesAttr.Value != "" {
		err = json.Unmarshal([]byte(variablesAttr.Value), &variables)
		if err != nil {
			return result, err
		}
	}

	qos := 0
	qosAttr := element.SelectAttr("qos")
	if qosAttr != nil && qosAttr.Value != "" {
		qos, err = strconv.Atoi(qosAttr.Value)
		if err != nil {
			return result, err
		}
	}

	result = deploymentmodel.Element{
		Name:   label,
		BpmnId: id,
		Order:  this.getOrder(element),
		ConditionalEvent: &deploymentmodel.ConditionalEvent{
			Script:        script,
			ValueVariable: valueVariableName,
			Variables:     variables,
			Qos:           qos,
			EventId:       "",
			Selection: deploymentmodel.Selection{
				FilterCriteria: filterCriteria,
			},
		},
	}
	return result, nil
}
