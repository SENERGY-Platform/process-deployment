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
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
	"strings"
)

func MsgEvent(doc *etree.Document, event *model.MsgEvent) (err error) {
	if event == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	msgRef := strings.Replace("e_"+event.EventId, "-", "_", -1)
	bpmnMsg := doc.CreateElement("bpmn:message")
	doc.SelectElement("bpmn:definitions").InsertChild(doc.SelectElement("bpmn:definitions").SelectElement("bpmndi:BPMNDiagram"), bpmnMsg)
	bpmnMsg.CreateAttr("id", msgRef)
	bpmnMsg.CreateAttr("name", event.EventId)
	doc.FindElement("//*[@id='"+event.BpmnElementId+"']/bpmn:messageEventDefinition").CreateAttr("messageRef", msgRef)
	return nil
}
