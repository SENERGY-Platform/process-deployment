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
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/beevik/etree"
)

func init() {
	ElementParsers = append(ElementParsers, ElementParser{
		Is: func(this *Parser, element *etree.Element) bool {
			return this.isTimeEvent(element)
		},
		Parse: func(this *Parser, dom *etree.Element) (element deploymentmodel.Element, err error) {
			return this.getTimeEvent(dom)
		},
	})
}

func (this *Parser) isTimeEvent(element *etree.Element) bool {
	return element.FindElement("./bpmn:timerEventDefinition") != nil
}

func (this *Parser) getTimeEvent(element *etree.Element) (result deploymentmodel.Element, err error) {
	id := element.SelectAttr("id").Value
	label := element.SelectAttrValue("name", id)

	timeEvent := deploymentmodel.TimeEvent{}

	eventDefinitionWrapper := element.FindElement("./bpmn:timerEventDefinition")
	for _, child := range eventDefinitionWrapper.ChildElements() {
		timeEvent.Type = child.Tag
		timeEvent.Time = child.Text()
	}
	if timeEvent.Type == "" {
		timeEvent.Time = ""
		if element.Tag == "startEvent" {
			timeEvent.Type = "timeCycle"
		} else {
			timeEvent.Type = "timeDuration"
		}
	}

	result = deploymentmodel.Element{
		Name:      label,
		BpmnId:    id,
		Order:     this.getOrder(element),
		TimeEvent: &timeEvent,
	}
	return result, nil
}
