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

package parser

import (
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/beevik/etree"
)

func (this *Parser) getLanes(process *etree.Element) (result []deploymentmodel.Lane, err error) {
	laneSet := process.FindElement("./laneSet")
	if laneSet != nil {
		for i, laneElement := range laneSet.FindElements("./lane") {
			laneId := laneElement.SelectAttrValue("id", "")
			name := laneElement.SelectAttrValue("name", "")
			ids, err := this.getLaneIds(laneElement)
			if err != nil {
				return result, err
			}
			lane, err := this.getLane(process, laneId, name, int64(i), ids)
			if err != nil {
				return result, err
			}
			result = append(result, lane)
		}
	} else {
		ids, err := this.geProcessChildIds(process)
		if err != nil {
			return result, err
		}
		lane, err := this.getLane(process, "", "", 0, ids)
		if err != nil {
			return result, err
		}
		result = append(result, lane)
	}
	return
}

func (this *Parser) getNoLanes(process *etree.Element) (result []deploymentmodel.Lane, err error) {
	ids, err := this.geProcessChildIds(process)
	if err != nil {
		return result, err
	}
	for _, id := range ids {
		lane, err := this.getLane(process, "", "", 0, []string{id})
		if err != nil {
			return result, err
		}
		if len(lane.Elements) > 0 {
			lane.Order = lane.Elements[0].Order
			result = append(result, lane)
		}
	}
	return
}

func (this *Parser) getLane(process *etree.Element, id string, name string, order int64, elementIds []string) (result deploymentmodel.Lane, err error) {
	result.BpmnId = id
	result.Name = name
	result.Order = order
	result.Elements, err = this.getElements(process, elementIds)
	return
}

func (this *Parser) getLaneIds(lane *etree.Element) (result []string, err error) {
	for _, ref := range lane.FindElements("./bpmn:flowNodeRef") {
		id := ref.Text()
		result = append(result, id)
	}
	return result, nil
}

func (this *Parser) geProcessChildIds(process *etree.Element) (result []string, err error) {
	for _, element := range process.FindElements(".//*[@id]") {
		result = append(result, element.SelectAttr("id").Value)
	}
	return
}
