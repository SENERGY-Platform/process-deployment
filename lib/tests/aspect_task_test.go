/*
 * Copyright 2026 InfAI (CC SES)
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

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/SENERGY-Platform/camunda-engine-wrapper/lib/shards"
	"github.com/SENERGY-Platform/camunda-engine-wrapper/lib/shards/cache"
	"github.com/SENERGY-Platform/device-repository/v2/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/process-deployment/lib"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/resources/integrationtest"
	"github.com/testcontainers/testcontainers-go"
)

const thermostatDeviceTypeId = "urn:infai:ses:device-type:d447354d-2773-4a9c-b393-f07c745b04d4"

// TestAspectTaskCommand deploys and starts a process through process-deployment, the
// engine-wrapper and camunda, lets the external-task-worker execute its task with the marshaller
// and the converter, and checks the command that arrives at the protocol handler (SNRGY-4845).
//
// The function of the task, urn:infai:ses:controlling-function:df08a869-f6b7-4c8a-ada7-2147255b6d70,
// is found in two variables of the "set" service of the device type: external_temperature_input
// carries the aspect "Reading" of the task, occupied_heating_setpoint does not. Both are
// omit_empty, so the variable without the aspect must not be part of the command at all.
//
// The protocol handler is an http url instead of a kafka topic, the worker posts the command to
// it. The test serves that url itself, next to a stand-in for keycloak that the worker and the
// marshaller fetch their tokens from.
func TestAspectTaskCommand(t *testing.T) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	commands := make(chan []byte, 10)
	server := httptest.NewServer(commandAndTokenHandler(t, commands))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port
	serverUrlInContainers := fmt.Sprintf("http://%s:%d", testcontainers.HostInternal, serverPort)

	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}
	conf.Debug = true
	conf.ConnectivityTest = false
	conf.DeploymentTopic = "-"
	conf.DoneTopic = "process-deployment-done"
	conf.DeviceGroupTopic = "device-groups"
	conf.InitTopics = true

	_, mongoIp, err := docker.Mongo(ctx, &wg)
	if err != nil {
		t.Error(err)
		return
	}
	conf.MongoUrl = "mongodb://" + mongoIp + ":27017"

	_, camundaPgIp, _, err := docker.PostgresWithNetwork(ctx, &wg, "camunda")
	if err != nil {
		t.Error(err)
		return
	}

	pgConn, err := docker.Postgres(ctx, &wg, "shards")
	if err != nil {
		t.Error(err)
		return
	}

	camundaUrl, err := docker.Camunda(ctx, &wg, camundaPgIp, "5432")
	if err != nil {
		t.Error(err)
		return
	}

	s, err := shards.New(pgConn, cache.None)
	if err != nil {
		t.Error(err)
		return
	}
	err = s.EnsureShard(camundaUrl)
	if err != nil {
		t.Error(err)
		return
	}

	_, incidentsApiIp, err := docker.IncidentsApi(ctx, &wg, pgConn, conf.MongoUrl, "")
	if err != nil {
		t.Error(err)
		return
	}
	incidentApiUrl := "http://" + incidentsApiIp + ":8080"

	_, engineWrapperIp, err := docker.EngineWrapper(ctx, &wg, incidentApiUrl, pgConn)
	if err != nil {
		t.Error(err)
		return
	}
	conf.ProcessEngineWrapperUrl = "http://" + engineWrapperIp + ":8080"

	conf.KafkaUrl, err = docker.Kafka(ctx, &wg)
	if err != nil {
		t.Error(err)
		return
	}

	_, permV2Ip, err := docker.PermissionsV2(ctx, &wg, conf.MongoUrl, conf.KafkaUrl)
	if err != nil {
		t.Error(err)
		return
	}
	conf.PermissionsV2Url = "http://" + permV2Ip + ":8080"

	_, repoIp, err := docker.DeviceRepo(ctx, &wg, conf.KafkaUrl, conf.MongoUrl, conf.PermissionsV2Url)
	if err != nil {
		t.Error(err)
		return
	}
	conf.DeviceRepoUrl = "http://" + repoIp + ":8080"

	devices := client.NewClient(conf.DeviceRepoUrl, nil)

	//the marshaller loads concepts and characteristics on start, so they are stored before it starts
	deviceType, err := storeThermostatMetadata(devices, serverUrlInContainers+"/commands")
	if err != nil {
		t.Error(err)
		return
	}
	setServiceId := ""
	for _, service := range deviceType.Services {
		if service.LocalId == "set" {
			setServiceId = service.Id
		}
	}
	if setServiceId == "" {
		t.Error("device type has no service with local id set")
		return
	}

	_, converterIp, err := docker.Converter(ctx, &wg)
	if err != nil {
		t.Error(err)
		return
	}

	_, marshallerIp, err := docker.Marshaller(ctx, &wg, conf.DeviceRepoUrl, "http://"+converterIp+":8080", serverUrlInContainers, []int{serverPort})
	if err != nil {
		t.Error(err)
		return
	}

	_, memcachedIp, err := docker.Memcached(ctx, &wg)
	if err != nil {
		t.Error(err)
		return
	}

	err = docker.TaskWorkerWithEnv(ctx, &wg, conf.DeviceRepoUrl, conf.KafkaUrl, incidentApiUrl, pgConn, memcachedIp+":11211", map[string]string{
		"AUTH_ENDPOINT":  serverUrlInContainers,
		"MARSHALLER_URL": "http://" + marshallerIp + ":8080",
	}, []int{serverPort})
	if err != nil {
		t.Error(err)
		return
	}

	freePort, err := GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}
	conf.ApiPort = strconv.Itoa(freePort)

	err = lib.StartDefault(ctx, conf)
	if err != nil {
		t.Error(err)
		return
	}

	deploymentUrl := "http://localhost:" + conf.ApiPort

	tokenObj, err := auth.CreateToken("test", "test")
	if err != nil {
		t.Error(err)
		return
	}
	token := tokenObj.Token

	device, err, _ := devices.CreateDevice(token, models.Device{
		LocalId:      "thermostat",
		Name:         "thermostat",
		DeviceTypeId: deviceType.Id,
	})
	if err != nil {
		t.Error(err)
		return
	}

	//permissions-v2 learns about the device asynchronously, the deployment checks it
	err = docker.Retry(time.Minute, func() error {
		_, err, _ := devices.ReadDevice(device.Id, token, models.Execute)
		return err
	})
	if err != nil {
		t.Error(err)
		return
	}

	//294.15 K are 21 °C, which external_temperature_input expects
	deplId, err := integrationtest.DeployAspectTaskProcess(token, deploymentUrl, device.Id, setServiceId, "294.15")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err := deleteDeployment(token, deploymentUrl, deplId)
		if err != nil {
			t.Error(err)
		}
	}()

	//the engine-wrapper deploys asynchronously
	err = docker.Retry(time.Minute, func() error {
		return startDeployment(token, conf.ProcessEngineWrapperUrl, deplId)
	})
	if err != nil {
		t.Error(err)
		return
	}

	var command []byte
	select {
	case command = <-commands:
	case <-time.After(2 * time.Minute):
		t.Errorf("no command received; incidents: %v", listIncidents(token, incidentApiUrl))
		return
	}

	var msg ProtocolMsg
	err = json.Unmarshal(command, &msg)
	if err != nil {
		t.Errorf("command is no protocol message: %v\n%s", err, command)
		return
	}
	if msg.Metadata.Service.Id != setServiceId {
		t.Errorf("command is for service %q, want %q", msg.Metadata.Service.Id, setServiceId)
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(msg.Request.Input["data"]), &data)
	if err != nil {
		t.Errorf("data segment of the command is no json object: %v\n%#v", err, msg.Request.Input)
		return
	}

	t.Run("sets external_temperature_input, which carries the aspect", func(t *testing.T) {
		value, ok := data["external_temperature_input"].(float64)
		if !ok {
			t.Errorf("external_temperature_input is missing or no number: %#v", data)
			return
		}
		if math.Abs(value-21) > 1e-9 {
			t.Errorf("external_temperature_input = %v, want 21", value)
		}
	})

	t.Run("omits occupied_heating_setpoint, which lacks the aspect", func(t *testing.T) {
		if value, ok := data["occupied_heating_setpoint"]; ok {
			t.Errorf("occupied_heating_setpoint = %#v, want it omitted: %#v", value, data)
		}
	})

	t.Run("omits every other variable of the service", func(t *testing.T) {
		if len(data) != 1 {
			t.Errorf("want external_temperature_input only, got %#v", data)
		}
	})
}

// ProtocolMsg is the part of the message of the external-task-worker to a protocol handler that
// the test looks at.
type ProtocolMsg struct {
	Request struct {
		Input map[string]string `json:"input"`
	} `json:"request"`
	Metadata struct {
		Service models.Service `json:"service"`
	} `json:"metadata"`
}

// commandAndTokenHandler receives the commands of the external-task-worker under /commands and
// stands in for keycloak: a token exchange gets a user token for the requested subject, every
// other grant, as the client credentials of the marshaller, an admin token.
func commandAndTokenHandler(t *testing.T, commands chan<- []byte) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /commands", func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		select {
		case commands <- body:
		default:
			t.Log("drop command, nobody waits for it:", string(body))
		}
		writer.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /auth/realms/master/protocol/openid-connect/token", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token := client.InternalAdminToken
		if subject := request.PostForm.Get("requested_subject"); subject != "" {
			userToken, err := auth.CreateTokenWithRoles("test", subject, []string{"user"})
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			token = userToken.Token
		}
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(map[string]interface{}{
			"access_token":       strings.TrimPrefix(token, "Bearer "),
			"token_type":         "bearer",
			"expires_in":         600,
			"refresh_expires_in": 600,
		})
	})
	return mux
}

// storeThermostatMetadata stores the device type of the test with everything it references and
// points its protocol to protocolHandler.
func storeThermostatMetadata(devices client.Interface, protocolHandler string) (deviceType models.DeviceType, err error) {
	metadata, err := integrationtest.LoadThermostatMetadata()
	if err != nil {
		return deviceType, err
	}
	for _, aspectClass := range metadata.AspectClasses {
		_, err, _ = devices.SetAspectClass(client.InternalAdminToken, aspectClass)
		if err != nil {
			return deviceType, fmt.Errorf("aspect class %v: %w", aspectClass.Id, err)
		}
	}
	for _, aspect := range metadata.Aspects {
		_, err, _ = devices.SetAspect(client.InternalAdminToken, aspect)
		if err != nil {
			return deviceType, fmt.Errorf("aspect %v: %w", aspect.Id, err)
		}
	}
	for _, characteristic := range metadata.Characteristics {
		_, err, _ = devices.SetCharacteristic(client.InternalAdminToken, characteristic)
		if err != nil {
			return deviceType, fmt.Errorf("characteristic %v: %w", characteristic.Id, err)
		}
	}
	for _, concept := range metadata.Concepts {
		_, err, _ = devices.SetConcept(client.InternalAdminToken, concept)
		if err != nil {
			return deviceType, fmt.Errorf("concept %v: %w", concept.Id, err)
		}
	}
	for _, function := range metadata.Functions {
		_, err, _ = devices.SetFunction(client.InternalAdminToken, function)
		if err != nil {
			return deviceType, fmt.Errorf("function %v: %w", function.Id, err)
		}
	}
	for _, deviceClass := range metadata.DeviceClasses {
		_, err, _ = devices.SetDeviceClass(client.InternalAdminToken, deviceClass)
		if err != nil {
			return deviceType, fmt.Errorf("device class %v: %w", deviceClass.Id, err)
		}
	}
	for _, protocol := range metadata.Protocols {
		protocol.Handler = protocolHandler
		_, err, _ = devices.SetProtocol(client.InternalAdminToken, protocol)
		if err != nil {
			return deviceType, fmt.Errorf("protocol %v: %w", protocol.Id, err)
		}
	}
	for _, dt := range metadata.DeviceTypes {
		if dt.Id != thermostatDeviceTypeId {
			continue
		}
		deviceType, err, _ = devices.SetDeviceType(client.InternalAdminToken, dt, client.DeviceTypeUpdateOptions{})
		if err != nil {
			return deviceType, fmt.Errorf("device type %v: %w", dt.Id, err)
		}
		return deviceType, nil
	}
	return deviceType, fmt.Errorf("metadata contains no device type %v", thermostatDeviceTypeId)
}

func startDeployment(token string, engineWrapperUrl string, deploymentId string) error {
	req, err := http.NewRequest(http.MethodGet, engineWrapperUrl+"/v2/deployments/"+url.PathEscape(deploymentId)+"/start", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("start deployment: %v %s", resp.StatusCode, body)
	}
	return nil
}

func deleteDeployment(token string, deploymentUrl string, deploymentId string) error {
	req, err := http.NewRequest(http.MethodDelete, deploymentUrl+"/v3/deployments/"+url.PathEscape(deploymentId), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete deployment: %v %s", resp.StatusCode, body)
	}
	return nil
}

// listIncidents explains a missing command: a task that fails in the worker or the marshaller
// ends as an incident.
func listIncidents(token string, incidentApiUrl string) string {
	req, err := http.NewRequest(http.MethodGet, incidentApiUrl+"/incidents", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
