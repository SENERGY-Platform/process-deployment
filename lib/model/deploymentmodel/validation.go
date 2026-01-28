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

package deploymentmodel

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/SENERGY-Platform/models/go/models"
	"github.com/beevik/etree"
)

type ValidationKind = models.ValidationKind

const (
	ValidatePublish ValidationKind = models.ValidatePublish
	ValidateRequest ValidationKind = models.ValidateRequest
)

func DeploymentXmlValidator(xml string) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			slog.Default().Error("recovered from panic in DeploymentXmlValidator", "error", r)
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	doc := etree.NewDocument()
	err = doc.ReadFromString(xml)
	if err != nil {
		return err
	}

	//may not contain scripts that access 'execution'
	scripts := []string{}
	for _, script := range doc.FindElements("//camunda:script") {
		scripts = append(scripts, script.Text())
	}
	for _, script := range doc.FindElements("//bpmn:script") {
		scripts = append(scripts, script.Text())
	}
	for _, script := range scripts {
		if strings.Contains(script, "execution.") {
			return errors.New("bpmn script contains engine access")
		}
	}

	for _, connectors := range doc.FindElements("//camunda:connectorId") {
		if connectors.Text() == "mail-send" {
			return errors.New("bpmn contains mail-send connector, which is not allowed")
		}
	}

	return nil
}
