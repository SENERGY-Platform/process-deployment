/*
 * Copyright 2018 InfAI (CC SES)
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

package com

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SmartEnergyPlatform/process-deployment/lib/util"
)

func buildPayLoad(name string, xml string, svg string, boundary string, owner string) string {
	segments := []string{}
	deploymentSource := "sepl"

	segments = append(segments, "Content-Disposition: form-data; name=\"data\"; "+"filename=\""+name+".bpmn\"\r\nContent-Type: text/xml\r\n\r\n"+xml+"\r\n")
	segments = append(segments, "Content-Disposition: form-data; name=\"diagram\"; "+"filename=\""+name+".svg\"\r\nContent-Type: image/svg+xml\r\n\r\n"+svg+"\r\n")
	segments = append(segments, "Content-Disposition: form-data; name=\"deployment-name\"\r\n\r\n"+name+" \r\n")
	segments = append(segments, "Content-Disposition: form-data; name=\"deployment-source\"\r\n\r\n"+deploymentSource+"\r\n")
	segments = append(segments, "Content-Disposition: form-data; name=\"tenant-id\"\r\n\r\n"+owner+"\r\n")

	return "--" + boundary + "\r\n" + strings.Join(segments, "--"+boundary+"\r\n") + "--" + boundary + "--\r\n"
}

func DeployProcess(name string, xml string, svg string, owner string) (processId string, err error) {
	boundary := "---------------------------" + time.Now().String()
	b := strings.NewReader(buildPayLoad(name, xml, svg, boundary, owner))
	resp, err := http.Post(util.Config.ProcessEngineUrl+"/engine-rest/deployment/create", "multipart/form-data; boundary="+boundary, b)
	if err != nil {
		log.Println("ERROR: request to processengine ", err)
		return processId, err
	}
	responseWrapper := struct {
		Id string `json:"id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&responseWrapper)
	processId = responseWrapper.Id
	if err == nil && processId == "" {
		err = errors.New("process-engine didnt deploy process: " + xml)
	}
	log.Println("DEBUG: DeployProcess() = ", responseWrapper)
	return
}

func RemoveProcess(id string) (err error) {
	client := &http.Client{}
	url := util.Config.ProcessEngineUrl + "/engine-rest/deployment/" + id + "?cascade=true"
	request, err := http.NewRequest("DELETE", url, nil)
	_, err = client.Do(request)
	return
}
