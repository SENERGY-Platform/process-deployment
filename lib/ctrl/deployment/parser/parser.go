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
	"errors"
	"fmt"

	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/beevik/etree"
)

type Parser struct {
	conf config.Config
}

func New(conf config.Config) *Parser {
	return &Parser{conf: conf}
}

func (this *Parser) PrepareDeployment(xml string) (result deploymentmodel.Deployment, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			this.conf.GetLogger().Error("recovered from panic", "error", r)
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	doc := etree.NewDocument()
	err = doc.ReadFromString(xml)
	if err != nil {
		return
	}
	return this.getDeployment(doc, deploymentmodel.Diagram{XmlRaw: xml})
}

func (this *Parser) EstimateStartParameter(xml string) (result []deploymentmodel.ProcessStartParameter, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			this.conf.GetLogger().Error("recovered from panic", "error", r)
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	doc := etree.NewDocument()
	err = doc.ReadFromString(xml)
	if err != nil {
		return
	}
	return this.estimateStartParameter(doc)
}

func (this *Parser) estimateStartParameter(doc *etree.Document) (result []deploymentmodel.ProcessStartParameter, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			this.conf.GetLogger().Error("recovered from panic", "error", r)
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	elements := doc.FindElements("//bpmn:startEvent/bpmn:extensionElements/camunda:formData/camunda:formField")
	for _, element := range elements {
		id := element.SelectAttrValue("id", "")
		if id != "" {
			label := element.SelectAttrValue("label", "")
			paramtype := element.SelectAttrValue("type", "string")
			defaultValue := element.SelectAttrValue("defaultValue", "")
			properties := map[string]string{}
			for _, propterty := range element.FindElements(".//camunda:property") {
				propertyName := propterty.SelectAttrValue("id", "")
				propertyValue := propterty.SelectAttrValue("value", "")
				if propertyName != "" {
					properties[propertyName] = propertyValue
				}
			}
			result = append(result, deploymentmodel.ProcessStartParameter{
				Id:         id,
				Label:      label,
				Type:       paramtype,
				Default:    defaultValue,
				Properties: properties,
			})
		}
	}
	return result, nil
}
