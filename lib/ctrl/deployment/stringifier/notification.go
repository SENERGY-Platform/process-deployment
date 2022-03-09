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
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/executionmodel"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
)

func (this *Stringifier) Notification(doc *etree.Document, element deploymentmodel.Element, userId string) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	input := doc.FindElement("//*[@id='" + element.BpmnId + "']//camunda:inputParameter[@name='deploymentIdentifier']")
	if input.Text() != "notification" {
		return errors.New("unexpected content in notification input parameter")
	}
	parent := input.Parent()

	// Set url input
	urlParameter := parent.CreateElement("camunda:inputParameter")
	urlParameter.CreateAttr("name", "url")
	urlParameter.SetText(this.conf.NotificationUrl)

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
	notificationPayload := executionmodel.NotificationPayload{
		Message: element.Notification.Message,
		UserId:  userId,
		Title:   element.Notification.Title,
		IsRead:  false,
	}
	txt, err := json.Marshal(notificationPayload)
	if err != nil {
		return err
	}
	payloadParameter.SetCData(string(txt))

	return nil
}
