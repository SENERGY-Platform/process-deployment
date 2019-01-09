/*
 * Copyright 2018 InfAI (CC SES)
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

package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/SmartEnergyPlatform/process-deployment/lib/com"
	"github.com/SmartEnergyPlatform/process-deployment/lib/etree"
	"github.com/SmartEnergyPlatform/process-deployment/lib/model"
	"github.com/SmartEnergyPlatform/process-deployment/lib/util"

	"github.com/satori/go.uuid"
)

const CAMUNDA_VARIABLES_PAYLOAD = "payload"
const CAMUNDA_VARIABLES_DATA_EXPORT_CONFIG = "config"

func CloneAbstractProcess(id string, jwtimpersonate jwt_http_router.JwtImpersonate, owner string) (result model.AbstractProcess, err error) {
	metadata, err := GetMetadata(id, owner)
	if err != nil {
		return result, err
	}
	result = metadata.Abstract
	for index, param := range result.AbstractTasks {
		param.Options, err = com.GetDeviceInstancesFromType(param.DeviceTypeId, jwtimpersonate)
		if err != nil {
			return result, err
		}
		result.AbstractTasks[index] = param
	}
	return
}

func GetBpmnAbstractPrepare(xmlValue string, jwtimpersonate jwt_http_router.JwtImpersonate) (resp model.AbstractProcess, err error) {
	resp.Xml = xmlValue
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	doc := etree.NewDocument()
	err = doc.ReadFromString(xmlValue)
	if err != nil {
		return
	}
	resp.Name = doc.FindElement("//bpmn:process").SelectAttr("id").Value
	reuseIds := map[string]string{} //map[taskId]reuseId
	for _, lane := range doc.FindElements("//bpmn:lane") {
		laneId := lane.SelectAttr("id").Value
		for _, taskRef := range lane.FindElements("bpmn:flowNodeRef") {
			reuseIds[taskRef.Text()] = laneId
		}
	}
	abstractBpmnMsgs := []model.BpmnAbstractMsg{}
	abstractDataExportTasks := []model.AbstractDataExportTask{}
	for _, task := range doc.FindElements("//bpmn:serviceTask") {
		topic := task.SelectAttr("camunda:topic")
		if topic != nil {
			if topic.Value == util.Config.IotTaskTopic {
				log.Println("DEBUG: service task is execute_in_dose")
				abstract, err := GetAbstractTask(task)
				if err != nil {
					return resp, err
				}
				abstract.ReuseId = reuseIds[abstract.TaskId]
				abstractBpmnMsgs = append(abstractBpmnMsgs, abstract)
			} else {
				log.Println("DEBUG: service task is export")
				abstract, err := getAbstractDataExportTask(task, jwtimpersonate)
				if err != nil {
					return resp, err
				}
				abstractDataExportTasks = append(abstractDataExportTasks, abstract)
			}
		} else {
			placeholderTask := model.PlaceholderTask{Id: task.SelectAttr("id").Value}
			for _, inputParameter := range task.FindElements("//camunda:inputParameter") {
				placeholder, err := getPlaceholder(inputParameter)
				if err != nil {
					log.Println("ERROR: unable to find placeholder", err)
					return resp, err
				}
				placeholderTask.Parameter = append(placeholderTask.Parameter, placeholder...)
			}
			if len(placeholderTask.Parameter) > 0 {
				resp.PlaceholderTasks = append(resp.PlaceholderTasks, placeholderTask)
			}
		}
	}
	for _, eventDefiniton := range doc.FindElements("//bpmn:messageEventDefinition") {
		shapeId := eventDefiniton.Parent().SelectAttr("id").Value
		resp.MsgEvents = append(resp.MsgEvents, model.MsgEvent{Filter: createDefaultFilter(), ShapeId: shapeId})
	}
	for _, receiveTask := range doc.FindElements("//bpmn:receiveTask") {
		shapeId := receiveTask.SelectAttr("id").Value
		resp.ReceiveTasks = append(resp.ReceiveTasks, model.MsgEvent{Filter: createDefaultFilter(), ShapeId: shapeId})
	}
	for _, timeEventDefinition := range doc.FindElements("//bpmn:timerEventDefinition") {
		shapeId := timeEventDefinition.Parent().SelectAttr("id").Value
		timeEvent := model.TimeEvent{ShapeId: shapeId}
		for _, child := range timeEventDefinition.ChildElements() {
			timeEvent.Kind = child.Tag
			timeEvent.Time = child.Text()
		}
		if timeEvent.Kind == "" {
			if timeEventDefinition.Parent().Tag == "startEvent" {
				timeEvent.Kind = "timeCycle"
			} else {
				timeEvent.Kind = "timeDuration"
			}
		}
		resp.TimeEvents = append(resp.TimeEvents, timeEvent)
	}
	resp.AbstractTasks, err = getAbstractParameters(abstractBpmnMsgs, jwtimpersonate)
	resp.AbstractDataExportTasks = abstractDataExportTasks
	return
}

func createDefaultFilter() model.Filter {
	return model.Filter{Scope: "none"}
}

func GetAbstractTask(task *etree.Element) (result model.BpmnAbstractMsg, err error) {
	script := task.FindElement("//camunda:inputParameter[@name='" + CAMUNDA_VARIABLES_PAYLOAD + "']")

	err = json.Unmarshal([]byte(script.Text()), &result)
	if err != nil {
		return result, err
	}

	result.TaskId = task.SelectAttr("id").Value
	result.Parameter, err = getAbstractTaskParameter(task)
	return
}

func getAbstractTaskParameter(task *etree.Element) (result []model.TaskParameter, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			err = errors.New(fmt.Sprint("Recovered Error: getAbstractTaskParameter() ", r))
		}
	}()
	for _, input := range task.FindElements("//camunda:inputParameter") {
		name := input.SelectAttr("name").Value
		value := input.Text()
		path := strings.SplitN(name, ".", 2)
		if len(path) == 2 && path[0] == "inputs" && path[1] != "" && value == "" && len(input.ChildElements()) == 0 {
			result = append(result, model.TaskParameter{Name: name, Value: value})
		}
	}
	return result, err
}

func getAbstractParameters(abstractBpmnMsgs []model.BpmnAbstractMsg, jwtimpersonate jwt_http_router.JwtImpersonate) (result []model.AbstractTask, err error) {
	grouping := map[string][]model.BpmnAbstractMsg{}
	for _, abstract := range abstractBpmnMsgs {
		id := ""
		if abstract.ReuseId != "" {
			id = abstract.ReuseId + "." + abstract.DeviceType
		} else {
			id = abstract.TaskId + "." + abstract.DeviceType
		}
		grouping[id] = append(grouping[id], abstract)
	}

	cache := &map[string]model.AbstractTask{}
	for _, group := range grouping {
		param, err := toAbstractParameter(group[0], cache, jwtimpersonate)
		if err != nil {
			return result, err
		}
		for _, task := range group {
			param.Tasks = append(param.Tasks, model.Task{
				Id:        task.TaskId,
				Label:     task.Label,
				ServiceId: task.Service,
				Values:    task.Values,
				Parameter: task.Parameter,
			})
		}
		result = append(result, param)
	}
	return
}

func toAbstractParameter(abstract model.BpmnAbstractMsg, cache *map[string]model.AbstractTask, jwtimpersonate jwt_http_router.JwtImpersonate) (result model.AbstractTask, err error) {
	if cache == nil {
		cache = &map[string]model.AbstractTask{}
	}
	ok := false
	if result, ok = (*cache)[abstract.DeviceType]; !ok {
		result.DeviceTypeId = abstract.DeviceType
		result.Options, err = com.GetDeviceInstancesFromType(abstract.DeviceType, jwtimpersonate)
		(*cache)[abstract.DeviceType] = result
	}
	return
}

func InstantiateAbstractProcess(msg model.AbstractProcess, impersonate jwt_http_router.JwtImpersonate, userId string) (xmlString string, err error) {
	//catch etree exceptions/panics
	defer func() {
		if r := recover(); r != nil && err == nil {
			inp, _ := json.Marshal(msg)
			err = errors.New(fmt.Sprint("Recovered Error: ", r, string(inp)))
		}
	}()

	//use etree to set new process name in xml
	doc := etree.NewDocument()
	err = doc.ReadFromString(msg.Xml)
	if err != nil {
		return
	}
	processId := doc.FindElement("//bpmn:process").SelectAttr("id")
	for _, ref := range doc.FindElements("//*[@bpmnElement='" + processId.Value + "']") {
		ref.SelectAttr("bpmnElement").Value = msg.Name
	}
	processId.Value = msg.Name

	for _, param := range msg.AbstractTasks {
		taskInstance := model.BpmnMsg{
			InstanceId: param.Selected.Id,
		}

		if err = com.CheckExecuteRight(param.Selected.Id, impersonate); err != nil {
			log.Println("ERROR: com.CheckExecuteRight()", err)
			return xmlString, err
		}

		for _, abstractTask := range param.Tasks {
			taskId := abstractTask.Id
			taskInstance.ServiceId = abstractTask.ServiceId
			taskInstance.Inputs = abstractTask.Values.Inputs
			taskInstance.Outputs = abstractTask.Values.Outputs
			taskInstanceText, err := json.Marshal(taskInstance)
			if err != nil {
				return xmlString, err
			}
			xpath := "//bpmn:serviceTask[@id='" + taskId + "']//camunda:inputParameter[@name='" + CAMUNDA_VARIABLES_PAYLOAD + "']"
			doc.FindElement(xpath).WriteCData(string(taskInstanceText))

			//set new task name
			taskName := param.Selected.Name + ": " + abstractTask.Label
			doc.FindElement("//bpmn:serviceTask[@id='" + taskId + "']").SelectAttr("name").Value = taskName

			for _, taskparam := range abstractTask.Parameter {
				paramName := taskparam.Name
				paramValue := taskparam.Value
				xpath := "//bpmn:serviceTask[@id='" + taskId + "']//camunda:inputParameter[@name='" + paramName + "']"
				doc.FindElement(xpath).WriteCData(paramValue)
			}
		}
	}

	for _, abstractDataExportTask := range msg.AbstractDataExportTasks {
		// TODO check permission in serving service with user id here or in service task
		abstractDataExportTask.UserId = userId
		taskInstanceText, err := json.Marshal(abstractDataExportTask)
		if err != nil {
			return xmlString, err
		}
		xpath := "//bpmn:serviceTask[@id='" + abstractDataExportTask.Id + "']//camunda:inputParameter[@name='" + CAMUNDA_VARIABLES_DATA_EXPORT_CONFIG + "']"
		doc.FindElement(xpath).WriteCData(string(taskInstanceText))
	}

	for _, event := range msg.MsgEvents {
		if err = com.CheckExecuteRight(event.Filter.DeviceId, impersonate); err != nil {
			log.Println("ERROR: msgevent com.CheckExecuteRight()", err)
			return xmlString, err
		}
		id := uuid.NewV4()
		if err != nil {
			log.Println("ERROR in uuid:", err)
			return xmlString, err
		}
		msgRef := strings.Replace("e_"+id.String(), "-", "_", -1)
		bpmnMsg := doc.CreateElement("bpmn:message")
		doc.SelectElement("bpmn:definitions").InsertChild(doc.SelectElement("bpmn:definitions").SelectElement("bpmndi:BPMNDiagram"), bpmnMsg)
		bpmnMsg.CreateAttr("id", msgRef)
		bpmnMsg.CreateAttr("name", event.FilterId)
		doc.FindElement("//*[@id='"+event.ShapeId+"']/bpmn:messageEventDefinition").CreateAttr("messageRef", msgRef)
	}

	for _, event := range msg.ReceiveTasks {
		if err = com.CheckExecuteRight(event.Filter.DeviceId, impersonate); err != nil {
			log.Println("ERROR: receiveevent com.CheckExecuteRight()", err)
			return xmlString, err
		}
		id := uuid.NewV4()
		if err != nil {
			log.Println("ERROR in uuid:", err)
			return xmlString, err
		}
		msgRef := strings.Replace("e_"+id.String(), "-", "_", -1)
		bpmnMsg := doc.CreateElement("bpmn:message")
		doc.SelectElement("bpmn:definitions").InsertChild(doc.SelectElement("bpmn:definitions").SelectElement("bpmndi:BPMNDiagram"), bpmnMsg)
		bpmnMsg.CreateAttr("id", msgRef)
		bpmnMsg.CreateAttr("name", event.FilterId)
		doc.FindElement("//*[@id='"+event.ShapeId+"']").CreateAttr("messageRef", msgRef)
	}

	for _, event := range msg.TimeEvents {
		eventDef := doc.FindElement("//*[@id='" + event.ShapeId + "']/bpmn:timerEventDefinition")
		eventDef.Parent().CreateAttr("name", event.Label)

		//remove all children
		children := eventDef.ChildElements()
		for _, child := range children {
			eventDef.RemoveChild(child)
		}

		timeElement := doc.CreateElement("bpmn:" + event.Kind)
		timeElement.CreateAttr("xsi:type", "bpmn:tFormalExpression")
		timeElement.SetText(event.Time)
		eventDef.AddChild(timeElement)
	}

	for _, placeholderTask := range msg.PlaceholderTasks {
		xpath := "//bpmn:serviceTask[@id='" + placeholderTask.Id + "']//camunda:inputParameter"
		inputs := doc.FindElements(xpath)
		for _, input := range inputs {
			placeholderText := input.Text()
			text, err := renderPlaceholder(placeholderText, placeholderTask.Parameter)
			if err != nil {
				return xmlString, err
			}
			if placeholderText != text {
				input.SetText(text)
			}
		}
	}
	xmlString, err = doc.WriteToString()
	return
}

func getAbstractDataExportTask(task *etree.Element, jwtimpersonate jwt_http_router.JwtImpersonate) (result model.AbstractDataExportTask, err error) {
	// creates the abstract task for data exports
	// for data exports you have to select a measurement id (= data export for a device instance service) and the name of the measurement field that should be used for aggregation
	script := task.FindElement("//camunda:inputParameter[@name='" + CAMUNDA_VARIABLES_DATA_EXPORT_CONFIG + "']")
	err = json.Unmarshal([]byte(script.Text()), &result)
	if err != nil {
		return result, err
	}
	result.Id = task.SelectAttr("id").Value
	return
}
