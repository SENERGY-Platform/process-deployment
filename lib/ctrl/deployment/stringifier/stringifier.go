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

	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/beevik/etree"
)

type Stringifier struct {
	conf               config.Config
	aspectNodeProvider AspectNodeProvider
}

type AspectNodeProvider func(token auth.Token, aspectNodeId string) (aspectNode devicemodel.AspectNode, err error)

func New(conf config.Config, aspectNodeProvider AspectNodeProvider) *Stringifier {
	return &Stringifier{conf: conf, aspectNodeProvider: aspectNodeProvider}
}

func (this *Stringifier) Deployment(deployment deploymentmodel.Deployment, userId string, token auth.Token) (xml string, err error) {
	defer func() {
		if r := recover(); r != nil {
			this.conf.GetLogger().Error("recovered from panic", "error", r)
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	doc := etree.NewDocument()
	err = doc.ReadFromString(deployment.Diagram.XmlRaw)
	if err != nil {
		return
	}

	doc.FindElement("//bpmn:process").CreateAttr("name", deployment.Name)
	doc.FindElement("//bpmn:process").CreateAttr("isExecutable", "true")

	err = this.updateStartParameter(doc, deployment.StartParameter, deployment.IncidentHandling != nil && deployment.IncidentHandling.Restart)
	if err != nil {
		return "", err
	}

	for _, element := range deployment.Elements {
		err = this.Element(doc, element, userId, token)
		if err != nil {
			return "", err
		}
	}

	xml, err = doc.WriteToString()
	return
}

func (this *Stringifier) Element(doc *etree.Document, element deploymentmodel.Element, userId string, token auth.Token) error {
	if element.Task != nil {
		err := this.Task(doc, element, token)
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
	if element.ConditionalEvent != nil {
		err := this.ConditionalEvent(doc, element)
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
		err := this.Notification(doc, element, userId)
		if err != nil {
			return err
		}
	}
	return nil
}
