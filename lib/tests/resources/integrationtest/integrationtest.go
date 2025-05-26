/*
 * Copyright 2025 InfAI (CC SES)
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

package integrationtest

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"text/template"
)

//go:embed task.json
var TaskDeploymentTemplate string

func getTaskDeploymentMessage(deviceId string, serviceId string) (buff *bytes.Buffer, err error) {
	templ, err := template.New("deployment").Parse(TaskDeploymentTemplate)
	if err != nil {
		return buff, err
	}
	buff = &bytes.Buffer{}
	err = templ.Execute(buff, map[string]string{"DeviceId": deviceId, "ServiceId": serviceId})
	return buff, err
}

//go:embed restart.json
var RestartDeploymentTemplate string

func getRestartDeploymentMessage(deviceId string, serviceId string) (buff *bytes.Buffer, err error) {
	templ, err := template.New("deployment").Parse(RestartDeploymentTemplate)
	if err != nil {
		return buff, err
	}
	buff = &bytes.Buffer{}
	err = templ.Execute(buff, map[string]string{"DeviceId": deviceId, "ServiceId": serviceId})
	return buff, err
}

//go:embed event.json
var EventDeploymentTemplate string

func getEventDeploymentMessage(deviceId string, serviceId string) (buff *bytes.Buffer, err error) {
	templ, err := template.New("deployment").Parse(EventDeploymentTemplate)
	if err != nil {
		return buff, err
	}
	buff = &bytes.Buffer{}
	err = templ.Execute(buff, map[string]string{"DeviceId": deviceId, "ServiceId": serviceId})
	return buff, err
}

type Wrapper struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func DeployEventProcess(token string, deploymentUrl string, deviceId string, serviceId string) (deploymentId string, err error) {
	buff, err := getEventDeploymentMessage(deviceId, serviceId)
	if err != nil {
		return "", err
	}
	return DeployProcess(token, deploymentUrl, buff)
}

func DeployTaskProcess(token string, deploymentUrl string, deviceId string, serviceId string) (deploymentId string, err error) {
	buff, err := getTaskDeploymentMessage(deviceId, serviceId)
	if err != nil {
		return "", err
	}
	return DeployProcess(token, deploymentUrl, buff)
}

func DeployRestartProcess(token string, deploymentUrl string, deviceId string, serviceId string) (deploymentId string, err error) {
	buff, err := getRestartDeploymentMessage(deviceId, serviceId)
	if err != nil {
		return "", err
	}
	return DeployProcess(token, deploymentUrl, buff)
}

func DeployProcess(token string, deploymentUrl string, payload *bytes.Buffer) (deploymentId string, err error) {
	endpoint := deploymentUrl + "/v3/deployments?source=sepl"
	method := "POST"
	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		temp, _ := io.ReadAll(resp.Body) //read error response end ensure that resp.Body is read to EOF
		return "", errors.New("unable to deploy process: " + string(temp))
	}
	wrapper := Wrapper{}
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		_, _ = io.ReadAll(resp.Body) //ensure resp.Body is read to EOF
		return "", err
	}
	return wrapper.Id, nil
}
