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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	"github.com/beevik/etree"
)

type Stringifier struct {
	conf config.Config
}

func New(conf config.Config) *Stringifier {
	return &Stringifier{conf: conf}
}

func (this *Stringifier) Deployment(deployment deploymentmodel.Deployment, userId string) (xml string, err error) {
	doc := etree.NewDocument()
	err = doc.ReadFromString(deployment.Diagram.XmlRaw)
	if err != nil {
		return
	}

	doc.FindElement("//bpmn:process").CreateAttr("name", deployment.Name)
	doc.FindElement("//bpmn:process").CreateAttr("isExecutable", "true")

	for _, element := range deployment.Elements {
		err = this.Element(doc, element)
		if err != nil {
			return "", err
		}
	}

	xml, err = doc.WriteToString()
	return
}

func (this *Stringifier) Element(doc *etree.Document, element deploymentmodel.Element) error {
	if element.Task != nil {
		err := this.Task(doc, element)
		if err != nil {
			return err
		}
	}
	if element.MessageEvent != nil {
		err := this.MessageEvent(doc, element)
		if err != nil {
			return err
		}
	}
	if element.TimeEvent != nil {
		err := this.TimeEvent(doc, element)
		if err != nil {
			return err
		}
	}
	if element.Notification != nil {
		err := this.Notification(doc, element)
		if err != nil {
			return err
		}
	}
	return nil
}
