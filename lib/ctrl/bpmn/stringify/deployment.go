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
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
)

func Deployment(deployment model.Deployment, selectionAsRef bool, deviceRepo interfaces.Devices, userId string, notification_url string) (xml string, err error) {
	doc := etree.NewDocument()
	err = doc.ReadFromString(deployment.XmlRaw)
	if err != nil {
		return
	}

	doc.FindElement("//bpmn:process").CreateAttr("name", deployment.Name)
	doc.FindElement("//bpmn:process").CreateAttr("isExecutable", "true")

	for _, element := range deployment.Elements {
		err = Element(doc, element, selectionAsRef, deviceRepo)
		if err != nil {
			return "", err
		}
	}

	for _, lane := range deployment.Lanes {
		err = LaneElement(doc, lane, selectionAsRef, deviceRepo)
		if err != nil {
			return "", err
		}
	}

	for _, input := range doc.FindElements("//camunda:inputParameter[@name='deploymentIdentifier']") {
		if input.Text() != "notification" {
			continue
		}
		parent := input.Parent()

		// Set url input
		urlParameter := parent.CreateElement("camunda:inputParameter")
		urlParameter.CreateAttr("name", "url")
		urlParameter.SetText(notification_url)

		// Set method input
		methodParameter := parent.CreateElement("camunda:inputParameter")
		methodParameter.CreateAttr("name", "method")
		methodParameter.SetText("PUT")

		// Set header input
		headersParameter := parent.CreateElement("camunda:inputParameter")
		headersParameter.CreateAttr("name", "headers")
		mapElement := headersParameter.CreateElement("camunda:map")
		keyElement := mapElement.CreateElement("camunda:entry")
		keyElement.CreateAttr("key", "Content-Type")
		keyElement.SetText("application/json")

		// Set payload input
		payloadParameter := parent.FindElement("camunda:inputParameter[@name='payload']")
		notificationPayload := model.NotificationPayload{}
		err = json.Unmarshal([]byte(payloadParameter.Text()), &notificationPayload)
		if err != nil {
			return "", err
		}
		notificationPayload.UserId = userId
		notificationPayload.IsRead = false
		txt, err := json.Marshal(notificationPayload)
		if err != nil {
			return "", err
		}
		payloadParameter.SetText(string(txt))
	}

	xml, err = doc.WriteToString()
	return
}
