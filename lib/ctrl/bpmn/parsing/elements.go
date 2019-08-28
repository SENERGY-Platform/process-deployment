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
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
)

func BpmnToElements(doc *etree.Document) (result []model.Element, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if len(doc.FindElements("//bpmn:lane")) > 0 {
		//process uses lanes
		return result, nil
	}
	for _, task := range doc.FindElements("//bpmn:serviceTask") {
		topic := task.SelectAttr("camunda:topic")
		if topic != nil && topic.Value != "" {
			if len(task.FindElements("//bpmn:multiInstanceLoopCharacteristics")) > 0 {
				multitask, err := BpmnToMultitask(task)
				if err != nil {
					return result, err
				}
				result = append(result, multitask)
			} else {
				simpletask, err := BpmnToTask(task)
				if err != nil {
					return result, err
				}
				result = append(result, simpletask)
			}
		}
	}

	for _, event := range doc.FindElements("//bpmn:receiveTask") {
		msgEvent, err := BpmnToMsgEvent(event.Parent(), event)
		if err != nil {
			return result, err
		}
		result = append(result, msgEvent)
	}

	for _, event := range doc.FindElements("//bpmn:messageEventDefinition") {
		msgEvent, err := BpmnToMsgEvent(event.Parent(), event)
		if err != nil {
			return result, err
		}
		result = append(result, msgEvent)
	}

	for _, event := range doc.FindElements("//bpmn:timerEventDefinition") {
		timeEvent, err := BpmnToTimeEvent(event.Parent(), event)
		if err != nil {
			return result, err
		}
		result = append(result, timeEvent)
	}

	return result, nil
}
