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

package bpmn

import (
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/bpmn/parsing"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/bpmn/stringify"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
	"sort"
)

func PrepareDeployment(xml string) (result model.Deployment, err error) {
	result.XmlRaw = xml
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()

	doc := etree.NewDocument()
	err = doc.ReadFromString(xml)
	if err != nil {
		return
	}

	if len(doc.FindElements("//bpmn:collaboration")) > 0 {
		result.Name = doc.FindElement("//bpmn:collaboration").SelectAttrValue("id", "process-name")
		result.Lanes, err = parsing.BpmnToLanes(doc)
		if err != nil {
			return
		}
		sort.Sort(LaneByOrder(result.Lanes))
	} else {
		result.Name = doc.FindElement("//bpmn:process").SelectAttr("id").Value
		result.Elements, err = parsing.BpmnToElements(doc)
		if err != nil {
			return
		}
		sort.Sort(ElementByOrder(result.Elements))
	}
	return
}

func UseDeploymentSelections(deployment *model.Deployment, selectionAsRef bool, deviceRepo interfaces.Devices) (err error) {
	setMsgEventIds(deployment)

	// This function only gets called from tests. That's why it's okay to supply test data here
	deployment.Xml, err = stringify.Deployment(*deployment, selectionAsRef, deviceRepo, "uid", "url")
	return
}

func setMsgEventIds(deployment *model.Deployment) {
	for _, element := range deployment.Elements {
		if element.MsgEvent != nil {
			element.MsgEvent.EventId = config.NewId()
		}
		if element.ReceiveTaskEvent != nil {
			element.ReceiveTaskEvent.EventId = config.NewId()
		}
	}
	for _, lane := range deployment.Lanes {
		if lane.MultiLane != nil {
			for _, element := range lane.MultiLane.Elements {
				if element.MsgEvent != nil {
					element.MsgEvent.EventId = config.NewId()
				}
				if element.ReceiveTaskEvent != nil {
					element.ReceiveTaskEvent.EventId = config.NewId()
				}
			}
		}
		if lane.Lane != nil {
			for _, element := range lane.Lane.Elements {
				if element.MsgEvent != nil {
					element.MsgEvent.EventId = config.NewId()
				}
				if element.ReceiveTaskEvent != nil {
					element.ReceiveTaskEvent.EventId = config.NewId()
				}
			}
		}
	}
	return
}
