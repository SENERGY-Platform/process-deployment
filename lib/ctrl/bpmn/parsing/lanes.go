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
	"sort"
)

func BpmnToLanes(doc *etree.Document) (result []deploymentmodel.LaneElement, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	colab := doc.FindElement("//bpmn:collaboration")
	participants := colab.FindElements(".//bpmn:participant")
	for _, participant := range participants {
		processId := participant.SelectAttrValue("processRef", "")
		if processId == "" {
			debug.PrintStack()
			return result, errors.New("missing participant process ref")
		}
		process := doc.FindElement("//bpmn:process[@id='" + processId + "']")

		if len(process.FindElements(".//bpmn:lane")) > 0 {
			for _, lane := range process.FindElements(".//bpmn:lane") {
				subElements, err := getLaneSubElements(lane)
				if err != nil {
					return result, err
				}
				id, label, order, err := getLaneInfo(lane)
				element, err := createLaneElement(id, label, order, subElements)
				if err != EmptyLane {
					if err != nil {
						return result, err
					}
					result = append(result, element)
				}
			}
		} else {
			subElements, err := getLoneLaneSubElements(process)
			if err != nil {
				return result, err
			}
			id, label, order, err := getLaneInfo(process)
			label = participant.SelectAttrValue("name", label)
			element, err := createLaneElement(id, label, order, subElements)
			if err != EmptyLane {
				if err != nil {
					return result, err
				}
				result = append(result, element)
			}
		}
	}
	return
}

func getLaneInfo(lane *etree.Element) (id string, label string, order int64, err error) {
	documentation := executionmodel.Documentation{}
	documentations := lane.FindElements(".//bpmn:documentation")
	if len(documentations) > 0 {
		err = json.Unmarshal([]byte(documentations[0].Text()), &documentation)
		if err != nil {
			return
		}
	}
	order = documentation.Order
	id = lane.SelectAttr("id").Value
	label = lane.SelectAttrValue("name", id)
	return
}

var EmptyLane = errors.New("empty lane")

func createLaneElement(id string, label string, order int64, subElements []deploymentmodel.LaneSubElement) (result deploymentmodel.LaneElement, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	if len(subElements) == 0 {
		return result, EmptyLane
	}
	sort.Sort(LaneElementByOrder(subElements))

	result.Order = order
	deviceDescriptions := aggregateLaneTaskInfo(subElements)

	isMulti := isMultiTaskLane(subElements)
	if isMulti {
		result.MultiLane = &deploymentmodel.MultiLane{
			Label:              label,
			BpmnElementId:      id,
			DeviceDescriptions: deviceDescriptions,
			Elements:           subElements,
		}
	} else {
		result.Lane = &deploymentmodel.Lane{
			Label:              label,
			BpmnElementId:      id,
			DeviceDescriptions: deviceDescriptions,
			Elements:           subElements,
		}
	}
	return
}

func aggregateLaneTaskInfo(elements []deploymentmodel.LaneSubElement) (result []deploymentmodel.DeviceDescription) {
	for _, element := range elements {
		if element.LaneTask != nil {
			result = append(result, element.LaneTask.DeviceDescription)
		}
	}
	return
}

func isMultiTaskLane(elements []deploymentmodel.LaneSubElement) bool {
	if !hasTasks(elements) {
		return false
	}
	for _, element := range elements {
		if element.LaneTask != nil && !element.LaneTask.MultiTask {
			return false
		}
	}
	return true
}

func hasTasks(elements []deploymentmodel.LaneSubElement) bool {
	for _, element := range elements {
		if element.LaneTask != nil {
			return true
		}
	}
	return false
}

func getLaneSubElements(lane *etree.Element) (result []deploymentmodel.LaneSubElement, err error) {
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

func idIsDescendent(element *etree.Element, id string) bool {
	return len(element.FindElements(".//*[@id='"+id+"']")) > 0
}

//use this function if a collaboration with only one lane is used and no bpmn:flowNodeRef exists
func getLoneLaneSubElements(process *etree.Element) (result []deploymentmodel.LaneSubElement, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	root := process.FindElement("/")
	distinctIds := map[string]bool{}
	for _, flowEdge := range process.FindElements(".//bpmn:sequenceFlow") {
		distinctIds[flowEdge.SelectAttrValue("sourceRef", "")] = true
		distinctIds[flowEdge.SelectAttrValue("targetRef", "")] = true
	}
	for _, flowEdge := range process.FindElements("//bpmn:messageFlow") {
		sourceId := flowEdge.SelectAttrValue("sourceRef", "")
		if idIsDescendent(process, sourceId) {
			distinctIds[sourceId] = true
		}
		targetId := flowEdge.SelectAttrValue("targetRef", "")
		if idIsDescendent(process, targetId) {
			distinctIds[targetId] = true
		}
	}
	for id, _ := range distinctIds {
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

func getLaneSubElement(doc *etree.Element, id string) (result deploymentmodel.LaneSubElement, isDeploymentElement bool, err error) {
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
			result.LaneTask = &deploymentmodel.LaneTask{
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

type LaneElementByOrder []deploymentmodel.LaneSubElement

func (a LaneElementByOrder) Len() int           { return len(a) }
func (a LaneElementByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LaneElementByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }
