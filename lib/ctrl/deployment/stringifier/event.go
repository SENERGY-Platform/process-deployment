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
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/beevik/etree"
	uuid "github.com/satori/go.uuid"
	"log"
	"runtime/debug"
	"strings"
)

func (this *Stringifier) MessageEvent(doc *etree.Document, element deploymentmodel.Element) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if element.MessageEvent.EventId == "" {
		element.MessageEvent.EventId = "generated_" + uuid.NewV4().String() //element.MessageEvent is pointer so edit here edits value also for caller
	}
	msgRef := strings.Replace("generated_ref_"+element.MessageEvent.EventId, "-", "_", -1)
	bpmnMsg := doc.CreateElement("bpmn:message")
	doc.SelectElement("bpmn:definitions").InsertChildAt(doc.SelectElement("bpmn:definitions").SelectElement("bpmndi:BPMNDiagram").Index(), bpmnMsg)
	bpmnMsg.CreateAttr("id", msgRef)
	bpmnMsg.CreateAttr("name", element.MessageEvent.EventId)
	doc.FindElement("//*[@id='"+element.BpmnId+"']/bpmn:messageEventDefinition").CreateAttr("messageRef", msgRef)
	return nil
}

func (this *Stringifier) ConditionalEvent(doc *etree.Document, element deploymentmodel.Element) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if element.ConditionalEvent.EventId == "" {
		element.ConditionalEvent.EventId = "generated_" + uuid.NewV4().String() //element.MessageEvent is pointer so edit here edits value also for caller
	}
	msgRef := strings.Replace("generated_ref_"+element.ConditionalEvent.EventId, "-", "_", -1)
	bpmnMsg := doc.CreateElement("bpmn:message")
	doc.SelectElement("bpmn:definitions").InsertChildAt(doc.SelectElement("bpmn:definitions").SelectElement("bpmndi:BPMNDiagram").Index(), bpmnMsg)
	bpmnMsg.CreateAttr("id", msgRef)
	bpmnMsg.CreateAttr("name", element.ConditionalEvent.EventId)
	doc.FindElement("//*[@id='"+element.BpmnId+"']/bpmn:messageEventDefinition").CreateAttr("messageRef", msgRef)
	return nil
}

func (this *Stringifier) TimeEvent(doc *etree.Document, element deploymentmodel.Element) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	timeDefinition := doc.FindElement("//*[@id='" + element.BpmnId + "']/bpmn:timerEventDefinition")
	if definitionElement := timeDefinition.FindElement("./bpmn:" + element.TimeEvent.Type); definitionElement != nil {
		definitionElement.SetText(element.TimeEvent.Time)
	} else {
		definitionElement = doc.CreateElement("bpmn:" + element.TimeEvent.Type)
		definitionElement.CreateAttr("xsi:type", "bpmn:tFormalExpression")
		definitionElement.SetText(element.TimeEvent.Time)
		timeDefinition.InsertChildAt(0, definitionElement)
	}
	return nil
}
