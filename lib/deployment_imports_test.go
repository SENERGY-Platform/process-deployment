/*
 * Copyright 2021 InfAI (CC SES)
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

package lib

import (
	"context"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/importmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"os"
	"testing"
)

func TestImportDeployments(t *testing.T) {
	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "uuid"
	}

	conf, err := config.LoadConfig("../config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	conf.Debug = true
	conf.InitTopics = true
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	control, err := ctrl.New(ctx, conf, mocks.Kafka, mocks.Database, mocks.Devices, mocks.ProcessModelRepo, mocks.Imports, mocks.Engine, mocks.EventDepl)
	if err != nil {
		fmt.Println(err)
		return
	}

	t.Run("Test valid deployment", func(t *testing.T) {
		deployment, err := getValidImportDeployment()
		if err != nil {
			t.Error(err.Error())
		}
		err = deployment.Validate(deploymentmodel.ValidateRequest, map[string]bool{}, deploymentmodel.DeploymentXmlValidator)
		if err != nil {
			t.Error(err.Error())
		}
	})

	t.Run("Test missing path", func(t *testing.T) {
		deployment, err := getValidImportDeployment()
		if err != nil {
			t.Error(err.Error())
		}
		emptyString := ""
		if deployment.Elements[0].MessageEvent.Selection.SelectedPath == nil {
			deployment.Elements[0].MessageEvent.Selection.SelectedPath = &deviceselectionmodel.PathOption{}
		}
		deployment.Elements[0].MessageEvent.Selection.SelectedPath.CharacteristicId = emptyString
		err = deployment.Validate(deploymentmodel.ValidateRequest, map[string]bool{}, deploymentmodel.DeploymentXmlValidator)
		if err == nil {
			t.Error("Did not detect missing path")
		}
	})

	t.Run("Test missing characteristic", func(t *testing.T) {
		deployment, err := getValidImportDeployment()
		if err != nil {
			t.Error(err.Error())
		}
		emptyString := ""
		if deployment.Elements[0].MessageEvent.Selection.SelectedPath == nil {
			deployment.Elements[0].MessageEvent.Selection.SelectedPath = &deviceselectionmodel.PathOption{}
		}
		deployment.Elements[0].MessageEvent.Selection.SelectedPath.CharacteristicId = emptyString
		err = deployment.Validate(deploymentmodel.ValidateRequest, map[string]bool{}, deploymentmodel.DeploymentXmlValidator)
		if err == nil {
			t.Error("Did not detect missing characteristic")
		}
	})

	t.Run("Test missing access", func(t *testing.T) {
		deployment, err := getValidImportDeployment()
		if err != nil {
			t.Error(err.Error())
		}
		err, _ = control.EnsureDeploymentSelectionAccess(auth.Token{}, &deployment)
		if err == nil {
			t.Error("Allowed deploying import with no access to")
		}
	})

	t.Run("Test access ok", func(t *testing.T) {
		deployment, err := getValidImportDeployment()
		if err != nil {
			t.Error(err.Error())
		}
		mocks.Imports.SetImports([]importmodel.Import{{
			Id:           "urn:infai:ses:import:123",
			ImportTypeId: "urn:infai:ses:import-type:321",
		}})
		err, _ = control.EnsureDeploymentSelectionAccess(auth.Token{}, &deployment)
		if err != nil {
			t.Error(err.Error())
		}
	})

}

func getValidImportDeployment() (deploymentmodel.Deployment, error) {
	characteristicId := "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	aspectId := "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6"
	selectedImportId := "urn:infai:ses:import:123"
	selectedCharacteristicId := "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"
	selectedPath := "value.value"
	xmlRaw, err := os.ReadFile("./tests/resources/import.xml")
	if err != nil {
		return deploymentmodel.Deployment{}, err
	}
	xmlRawString := string(xmlRaw)
	return deploymentmodel.Deployment{
		Version: deploymentmodel.CurrentVersion,
		Id:      "deployment-import-id",
		Diagram: deploymentmodel.Diagram{XmlRaw: xmlRawString},
		Name:    "testEvent",
		Elements: []deploymentmodel.Element{{
			BpmnId: "bpmid",
			MessageEvent: &deploymentmodel.MessageEvent{
				Value:  "1",
				FlowId: "flow123",
				Selection: deploymentmodel.Selection{
					FilterCriteria: deploymentmodel.FilterCriteria{
						CharacteristicId: &characteristicId,
						FunctionId:       &functionId,
						AspectId:         &aspectId,
					},
					SelectedImportId: &selectedImportId,
					SelectedPath: &deviceselectionmodel.PathOption{
						Path:             selectedPath,
						CharacteristicId: selectedCharacteristicId,
					},
				},
			},
		},
		}}, nil
}
