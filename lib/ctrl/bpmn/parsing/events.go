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
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
)

func BpmnToMsgEvent(event *etree.Element) (result model.MsgEvent, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	id := event.SelectAttr("id").Value
	label := event.SelectAttrValue("name", id)
	documentation := model.Documentation{}
	documentations := event.FindElements(".//bpmn:documentation")
	if len(documentations) > 0 {
		err = json.Unmarshal([]byte(documentations[0].Text()), &documentation)
		if err != nil {
			return result, err
		}
	}
	result = model.MsgEvent{
		Label:         label,
		BpmnElementId: id,
		Order:         documentation.Order,
	}

	return result, nil
}

//TODO: test
func BpmnToTimeEvent(event *etree.Element, eventDefinition *etree.Element) (result model.TimeEvent, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	id := event.SelectAttr("id").Value
	label := event.SelectAttrValue("name", id)
	documentation := model.Documentation{}
	documentations := event.FindElements(".//bpmn:documentation")
	if len(documentations) > 0 {
		err = json.Unmarshal([]byte(documentations[0].Text()), &documentation)
		if err != nil {
			return result, err
		}
	}
	result = model.TimeEvent{
		BpmnElementId: id,
		Label:         label,
		Order:         documentation.Order,
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
	return result, nil
}
