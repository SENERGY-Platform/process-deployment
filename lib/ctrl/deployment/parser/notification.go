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
)

func init() {
	ElementParsers = append(ElementParsers, ElementParser{
		Is: func(this *Parser, element *etree.Element) bool {
			return this.isNotification(element)
		},
		Parse: func(this *Parser, dom *etree.Element) (element deploymentmodel.Element, err error) {
			return this.getNotification(dom)
		},
	})
}

func (this *Parser) isNotification(element *etree.Element) bool {
	if element.Tag != "serviceTask" && element.Tag != "bpmn:serviceTask" {
		return false
	}
	connector := element.FindElements("./bpmn:extensionElements/camunda:connector")
	if len(connector) != 1 {
		return false
	}

	connectorId := connector[0].FindElement("./camunda:connectorId")
	if connectorId == nil {
		return false
	}
	if connectorId.Text() != "http-connector" {
		return false
	}

	deploymentIdentifier := connector[0].FindElement("./camunda:inputOutput/camunda:inputParameter[@name='deploymentIdentifier']")
	if deploymentIdentifier == nil {
		return false
	}
	if deploymentIdentifier.Text() != "notification" {
		return false
	}
	return true
}

func (this *Parser) getNotification(element *etree.Element) (result deploymentmodel.Element, err error) {
	cmd := executionmodel.NotificationPayload{}
	cmdPayload := element.FindElement(".//camunda:inputParameter[@name='" + executionmodel.CAMUNDA_VARIABLES_PAYLOAD + "']")
	err = json.Unmarshal([]byte(cmdPayload.Text()), &cmd)
	if err != nil {
		return result, err
	}

	id := element.SelectAttr("id").Value
	label := element.SelectAttrValue("name", id)

	result = deploymentmodel.Element{
		Name:   label,
		BpmnId: id,
		Order:  this.getOrder(element),
		Notification: &deploymentmodel.Notification{
			Title:   cmd.Title,
			Message: cmd.Message,
		},
	}
	return result, nil
}
