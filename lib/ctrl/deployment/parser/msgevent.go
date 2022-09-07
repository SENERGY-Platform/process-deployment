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
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/beevik/etree"
	"strconv"
)

func init() {
	ElementParsers = append(ElementParsers, ElementParser{
		Is: func(this *Parser, element *etree.Element) bool {
			return this.isMsgEvent(element)
		},
		Parse: func(this *Parser, dom *etree.Element) (element deploymentmodel.Element, err error) {
			return this.getMsgEvent(dom)
		},
	})
}

func (this *Parser) isMsgEvent(element *etree.Element) bool {
	msgEvent := element.FindElement(".//bpmn:messageEventDefinition")
	if msgEvent == nil {
		return false //is not msg event
	}
	if msgEvent.SelectAttrValue("messageRef", "") != "" {
		return false //msg event uses user defined events
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
	return true
}

func (this *Parser) getMsgEvent(element *etree.Element) (result deploymentmodel.Element, err error) {
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

	useMarshallerBool := false
	useMarshaller := element.SelectAttr("use_marshaller")
	if useMarshaller != nil {
		useMarshallerBool, err = strconv.ParseBool(useMarshaller.Value)
		if err != nil {
			return result, fmt.Errorf("expect boolean in senergy:use_marshaller: %w", err)
		}
	}

	result = deploymentmodel.Element{
		Name:   label,
		BpmnId: id,
		Order:  this.getOrder(element),
		MessageEvent: &deploymentmodel.MessageEvent{
			EventId:       "",
			UseMarshaller: useMarshallerBool,
			Selection: deploymentmodel.Selection{
				FilterCriteria: filterCriteria,
			},
		},
	}
	return result, nil
}
