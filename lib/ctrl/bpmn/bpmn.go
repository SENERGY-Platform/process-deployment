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
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/bpmn/parsing"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
	"sort"
)

func BpmnToDeployment(xml string) (result model.Deployment, err error) {
	result.Xml = xml
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
	result.Name = doc.FindElement("//bpmn:process").SelectAttr("id").Value
	result.Lanes, err = parsing.BpmnToLanes(doc)
	if err != nil {
		return
	}
	result.Elements, err = parsing.BpmnToElements(doc)
	if err != nil {
		return
	}

	sort.Sort(LaneByOrder(result.Lanes))
	sort.Sort(ElementByOrder(result.Elements))

	return
}

func SetDeploymentXml(deployment *model.Deployment) (err error) {
	panic("not implemented") //TODO
}
