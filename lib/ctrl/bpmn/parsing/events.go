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

package parsing

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

func BpmnToMsgEvent(event *etree.Element) (ok bool, result deploymentmodel.MsgEvent, order int64, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if event.SelectAttrValue("messageRef", "") != "" {
		return false, result, 0, nil
	}
	if msgEvent := event.FindElement(".//bpmn:messageEventDefinition"); msgEvent != nil && msgEvent.SelectAttrValue("messageRef", "") != "" {
		return false, result, 0, nil
	}

	id := event.SelectAttr("id").Value
	label := event.SelectAttrValue("name", id)
	documentation := executionmodel.EventDocumentation{}
	documentations := event.FindElements(".//bpmn:documentation")
	if len(documentations) > 0 {
		err = json.Unmarshal([]byte(documentations[0].Text()), &documentation)
		if err != nil {
			return false, result, 0, err
		}
	}
	result = deploymentmodel.MsgEvent{
		Label:         label,
		BpmnElementId: id,
	}

	if documentation.CharacteristicId != "" {
		result.TriggerConversion = &deploymentmodel.Conversion{
			To: documentation.CharacteristicId,
		}
	}

	return true, result, documentation.Order, nil
}

func BpmnToTimeEvent(event *etree.Element, eventDefinition *etree.Element) (result deploymentmodel.TimeEvent, order int64, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	id := event.SelectAttr("id").Value
	label := event.SelectAttrValue("name", id)
	documentation := executionmodel.Documentation{}
	documentations := event.FindElements(".//bpmn:documentation")
	if len(documentations) > 0 {
		err = json.Unmarshal([]byte(documentations[0].Text()), &documentation)
		if err != nil {
			return result, 0, err
		}
	}
	result = deploymentmodel.TimeEvent{
		BpmnElementId: id,
		Label:         label,
	}
	for _, child := range eventDefinition.ChildElements() {
		result.Kind = child.Tag
		result.Time = child.Text()
	}
	if result.Kind == "" {
		if event.Tag == "startEvent" {
			result.Kind = "timeCycle"
		} else {
			result.Kind = "timeDuration"
		}
	}
	return result, documentation.Order, nil
}
