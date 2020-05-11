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

func (this *Parser) getDeployment(doc *etree.Document, diagram deploymentmodel.Diagram) (result deploymentmodel.Deployment, err error) {
	result.Diagram = diagram
	result.Name, err = getDeploymentName(doc)
	if err != nil {
		return result, err
	}
	result.Description, err = getDeploymentDescription(doc)
	if err != nil {
		return result, err
	}
	result.Pools, err = this.getPools(doc)
	return
}

func getDeploymentName(doc *etree.Document) (string, error) {
	colab := doc.FindElement("//bpmn:collaboration")
	if colab != nil {
		return colab.SelectAttrValue("id", "process-name"), nil
	} else {
		processId := doc.FindElement("//bpmn:process").SelectAttr("id").Value
		return doc.FindElement("//bpmn:process").SelectAttrValue("name", processId), nil
	}
}

func getDeploymentDescription(doc *etree.Document) (string, error) {
	colab := doc.FindElement("//bpmn:collaboration")
	if colab != nil {
		return colab.SelectAttrValue("description", ""), nil
	} else {
		return doc.FindElement("//bpmn:process").SelectAttrValue("description", ""), nil
	}
}
