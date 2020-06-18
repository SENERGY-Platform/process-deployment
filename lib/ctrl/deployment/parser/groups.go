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
	"github.com/beevik/etree"
)

type Group struct {
	Group    string
	Elements []string
}

func (this *Parser) getGroups(doc *etree.Document) (groups []Group, err error) {
	colab := doc.FindElement("//bpmn:collaboration")
	if colab != nil {
		participants := colab.FindElements("./bpmn:participant")
		for _, participant := range participants {
			id := participant.SelectAttr("processRef").Value
			err = this.appendPoolToGroups(doc, id, &groups)
			if err != nil {
				return groups, err
			}
		}
	} else {
		err = this.appendNoPoolToGroups(doc, &groups)
		if err != nil {
			return groups, err
		}
	}
	return
}

func (this *Parser) appendNoPoolToGroups(doc *etree.Document, groups *[]Group) (err error) {
	process := doc.FindElement("//bpmn:process")
	ids, err := this.geProcessChildIds(process)
	if err != nil {
		return err
	}
	*groups = append(*groups, Group{
		Elements: ids,
	})
	return
}

func (this *Parser) appendPoolToGroups(doc *etree.Document, id string, groups *[]Group) (err error) {
	pool := doc.FindElement("//bpmn:process[@id='" + id + "']")
	laneSet := pool.FindElement("./laneSet")
	if laneSet != nil {
		for _, laneElement := range laneSet.FindElements("./lane") {
			laneId := laneElement.SelectAttrValue("id", "")
			ids, err := this.getLaneIds(laneElement)
			if err != nil {
				return err
			}
			*groups = append(*groups, Group{
				Group:    "lane:" + laneId,
				Elements: ids,
			})
		}
	} else {
		ids, err := this.geProcessChildIds(pool)
		if err != nil {
			return err
		}
		*groups = append(*groups, Group{
			Group:    "pool:" + id,
			Elements: ids,
		})
	}
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
