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
	"sort"
)

func BpmnToLanes(doc *etree.Document) (result []model.LaneElement, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if len(doc.FindElements("//bpmn:lane")) == 0 {
		//process uses no lanes
		return result, nil
	}
	for _, lane := range doc.FindElements("//bpmn:lane") {
		element, err := bpmnToLane(lane)
		if err == EmptyLane {
			continue
		}
		if err != nil {
			return result, err
		}
		result = append(result, element)
	}
	return
}

var EmptyLane = errors.New("empty lane")

func bpmnToLane(lane *etree.Element) (result model.LaneElement, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	documentation := model.Documentation{}
	documentations := lane.FindElements(".//bpmn:documentation")
	if len(documentations) > 0 {
		err = json.Unmarshal([]byte(documentations[0].Text()), &documentation)
		if err != nil {
			return result, err
		}
	}
	result.Order = documentation.Order
	subElements, err := getLaneSubElements(lane)
	if err != nil {
		return result, err
	}
	if len(subElements) == 0 {
		return result, EmptyLane
	}
	sort.Sort(LaneElementByOrder(subElements))

	id := lane.SelectAttr("id").Value
	label := lane.SelectAttrValue("name", id)
	deviceDescriptions := aggregateLaneTaskInfo(subElements)

	isMulti := isMultiTaskLane(subElements)
	if isMulti {
		result.MultiLane = &model.MultiLane{
			Label:              label,
			BpmnElementId:      id,
			DeviceDescriptions: deviceDescriptions,
			Elements:           subElements,
		}
	} else {
		result.Lane = &model.Lane{
			Label:              label,
			BpmnElementId:      id,
			DeviceDescriptions: deviceDescriptions,
			Elements:           subElements,
		}
	}
	return
}

func aggregateLaneTaskInfo(elements []model.LaneSubElement) (result []model.DeviceDescription) {
	for _, element := range elements {
		if element.LaneTask != nil {
			result = append(result, element.LaneTask.DeviceDescription)
		}
	}
	return
}

func isMultiTaskLane(elements []model.LaneSubElement) bool {
	for _, element := range elements {
		if element.LaneTask != nil && !element.LaneTask.MultiTask {
			return false
		}
	}
	return true
}

func getLaneSubElements(lane *etree.Element) (result []model.LaneSubElement, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	root := lane.FindElement("/")
	for _, ref := range lane.FindElements(".//bpmn:flowNodeRef") {
		id := ref.Text()
		subelement, ok, err := getLaneSubElement(root, id)
		if err != nil {
			return result, err
		}
		if ok {
			result = append(result, subelement)
		}
	}
	return
}

func getLaneSubElement(doc *etree.Element, id string) (result model.LaneSubElement, isDeploymentElement bool, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	if task := doc.FindElement("//bpmn:serviceTask[@id='" + id + "']"); task != nil {
		topic := task.SelectAttr("camunda:topic")
		if topic != nil && topic.Value != "" {
			simpletask, err := BpmnToTask(task)
			if err != nil {
				return result, false, err
			}
			result.Order = simpletask.Order
			result.LaneTask = &model.LaneTask{
				Label:             simpletask.Task.Label,
				Retries:           simpletask.Task.Retries,
				DeviceDescription: simpletask.Task.DeviceDescription,
				Input:             simpletask.Task.Input,
				BpmnElementId:     simpletask.Task.BpmnElementId,
				MultiTask:         len(task.FindElements(".//bpmn:multiInstanceLoopCharacteristics")) > 0,
				Parameter:         simpletask.Task.Parameter,
			}
			return result, true, nil
		}
	}

	if event := doc.FindElement("//bpmn:receiveTask[@id='" + id + "']"); event != nil {
		ok, msgEvent, order, err := BpmnToMsgEvent(event)
		if err != nil {
			return result, false, err
		}
		result.Order = order
		result.ReceiveTaskEvent = &msgEvent
		return result, ok, nil
	}

	if event := doc.FindElement("//*[@id='" + id + "']//bpmn:messageEventDefinition"); event != nil {
		ok, msgEvent, order, err := BpmnToMsgEvent(event.Parent())
		if err != nil {
			return result, false, err
		}
		result.Order = order
		result.MsgEvent = &msgEvent
		return result, ok, nil
	}

	if event := doc.FindElement("//*[@id='" + id + "']//bpmn:timerEventDefinition"); event != nil {
		timeEvent, order, err := BpmnToTimeEvent(event.Parent(), event)
		if err != nil {
			return result, false, err
		}
		result.Order = order
		result.TimeEvent = &timeEvent
		return result, true, nil
	}

	return result, false, nil
}

type LaneElementByOrder []model.LaneSubElement

func (a LaneElementByOrder) Len() int           { return len(a) }
func (a LaneElementByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LaneElementByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }
