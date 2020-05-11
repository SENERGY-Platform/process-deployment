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
	"sort"
)

func (this *Parser) getPools(doc *etree.Document) (result []deploymentmodel.Pool, err error) {
	colab := doc.FindElement("//bpmn:collaboration")
	if colab != nil {
		participants := colab.FindElements("./bpmn:participant")
		for _, participant := range participants {
			id := participant.SelectAttr("processRef").Value
			name := participant.SelectAttrValue("name", "")
			pool, err := this.getPool(doc, id, name)
			if err != nil {
				return result, err
			}
			result = append(result, pool)
		}
	} else {
		pool, err := this.getNoPool(doc)
		if err != nil {
			return result, err
		}
		result = append(result, pool)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Order < result[j].Order
	})
	return
}

func (this *Parser) sort(list interface{}) func(int, int) bool {
	cast := list.([]deploymentmodel.BaseInfo)
	return func(i int, j int) bool {
		return cast[i].Order < cast[j].Order
	}
}

func (this *Parser) getPool(doc *etree.Document, id string, name string) (result deploymentmodel.Pool, err error) {
	element := doc.FindElement("//bpmn:process[@id='" + id + "']")
	result.Name = name
	result.Order = this.getOrder(element)
	result.BpmnId = id
	result.Lanes, err = this.getLanes(element)
	return
}

func (this *Parser) getNoPool(doc *etree.Document) (result deploymentmodel.Pool, err error) {
	element := doc.FindElement("//bpmn:process")
	result.BpmnId = element.SelectAttrValue("id", "")
	result.Lanes, err = this.getNoLanes(element)
	return
}
