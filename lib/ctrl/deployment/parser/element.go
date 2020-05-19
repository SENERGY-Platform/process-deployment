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

var ElementParsers = []ElementParser{}

type ElementParser struct {
	Is    func(this *Parser, element *etree.Element) bool
	Parse func(this *Parser, element *etree.Element) (deploymentmodel.Element, error)
}

func (this *Parser) getElements(doc *etree.Document) (result []deploymentmodel.Element, err error) {
	groups, err := this.getGroups(doc)
	if err != nil {
		return result, err
	}
	for _, group := range groups {
		for _, id := range group.Elements {
			element, isModelElement, err := this.getElement(doc, id)
			if err != nil {
				return result, err
			}
			if isModelElement {
				if group.Group != "" {
					group := group.Group
					element.Group = &group
				}
				result = append(result, element)
			}
		}
	}
	return result, err
}

func (this *Parser) getElement(doc *etree.Document, id string) (result deploymentmodel.Element, isElement bool, err error) {
	dom := doc.FindElement("//*[@id='" + id + "']")
	for _, parser := range ElementParsers {
		if parser.Is(this, dom) {
			result, err = parser.Parse(this, dom)
			if err != nil {
				return result, false, err
			}
			return result, true, nil
		}
	}
	return result, false, nil
}
