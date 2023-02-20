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
	"log"
	"runtime/debug"
)

func (this *Stringifier) updateStartParameter(doc *etree.Document, parameter []deploymentmodel.ProcessStartParameter, setIgnoreOnStartProperty bool) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	index := map[string]deploymentmodel.ProcessStartParameter{}
	for _, param := range parameter {
		index[param.Id] = param
	}
	elements := doc.FindElements("//bpmn:startEvent/bpmn:extensionElements/camunda:formData/camunda:formField")
	for _, element := range elements {
		id := element.SelectAttrValue("id", "")
		if id != "" {
			param, ok := index[id]
			if ok {
				element.CreateAttr("defaultValue", param.Default)
			}
			if setIgnoreOnStartProperty {
				properties := map[string]string{}
				for _, property := range element.FindElements(".//camunda:property") {
					propertyName := property.SelectAttrValue("id", "")
					propertyValue := property.SelectAttrValue("value", "")
					if propertyName == "ignore_on_start" && propertyValue != "true" {
						propertyValue = "true"
						property.CreateAttr("value", "true")
					}
					if propertyName != "" {
						properties[propertyName] = propertyValue
					}
				}
				propertiesListElement := element.FindElement("./camunda:properties")
				if propertiesListElement == nil {
					propertiesListElement = element.CreateElement("camunda:properties")
				}
				if _, ok := properties["ignore_on_start"]; !ok {
					propertyElement := propertiesListElement.CreateElement("camunda:property")
					propertyElement.CreateAttr("id", "ignore_on_start")
					propertyElement.CreateAttr("value", "true")
				}
			}
		}
	}
	return nil
}
