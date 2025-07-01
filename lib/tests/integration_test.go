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

package tests

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/camunda-engine-wrapper/lib/shards"
	"github.com/SENERGY-Platform/camunda-engine-wrapper/lib/shards/cache"
	"github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/process-deployment/lib"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/events/kafka"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/resources/integrationtest"
	kafkalib "github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}
	conf.Debug = true
	conf.ConnectivityTest = false
	//TODO: use as default config values
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

	_, zkIp, err := docker.Zookeeper(ctx, &wg)
	if err != nil {
		t.Error(err)
		return
	}

	zkUrl := zkIp + ":2181"

	conf.KafkaUrl, err = docker.Kafka(ctx, &wg, zkUrl)
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

	_, eventDeploymentIp, err := docker.EventDeployment(ctx, &wg, conf.DeviceRepoUrl, conf.MongoUrl)
	if err != nil {
		t.Error(err)
		return
	}
	conf.EventDeploymentUrl = "http://" + eventDeploymentIp + ":8080"

	eventTriggerUrl := conf.ProcessEngineWrapperUrl + "/v2/event-trigger"
	err = docker.EventWorker(ctx, &wg, conf.DeviceRepoUrl, eventTriggerUrl, conf.KafkaUrl, conf.MongoUrl)
	if err != nil {
		t.Error(err)
		return
	}

	_, memcachedIp, err := docker.Memcached(ctx, &wg)
	if err != nil {
		t.Error(err)
		return
	}

	err = docker.TaskWorker(ctx, &wg, conf.DeviceRepoUrl, conf.KafkaUrl, incidentApiUrl, pgConn, memcachedIp+":11211")
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

	devices := client.NewClient(conf.DeviceRepoUrl, nil)

	_, err, _ = devices.SetDeviceClass(client.InternalAdminToken, models.DeviceClass{
		Id:    "urn:infai:ses:device-class:ff64280a-58e6-4cf9-9a44-e70d3831a79d",
		Image: "",
		Name:  "dc",
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err, _ = devices.SetCharacteristic(client.InternalAdminToken, models.Characteristic{
		Id:   "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
		Name: "Degree Celsius",
		Type: models.Float,
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err, _ = devices.SetCharacteristic(client.InternalAdminToken, models.Characteristic{
		Id:   "urn:infai:ses:characteristic:a49a48fc-3a2c-4149-ac7f-1a5482d4c6e1",
		Name: "Degree Celsius (int)",
		Type: models.Integer,
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err, _ = devices.SetCharacteristic(client.InternalAdminToken, models.Characteristic{
		Id:          "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
		Name:        "Kelvin",
		DisplayUnit: "",
		Type:        models.Integer,
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err, _ = devices.SetConcept(client.InternalAdminToken, models.Concept{
		Id:   "urn:infai:ses:concept:0bc81398-3ed6-4e2b-a6c4-b754583aac37",
		Name: "Temperature",
		CharacteristicIds: []string{
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
			"urn:infai:ses:characteristic:a49a48fc-3a2c-4149-ac7f-1a5482d4c6e1",
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
		},
		BaseCharacteristicId: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
		Conversions: []models.ConverterExtension{
			{
				From:            "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
				To:              "urn:infai:ses:characteristic:a49a48fc-3a2c-4149-ac7f-1a5482d4c6e1",
				Distance:        1,
				Formula:         "atoi(ntoa(x))",
				PlaceholderName: "x",
			},
			{
				From:            "urn:infai:ses:characteristic:a49a48fc-3a2c-4149-ac7f-1a5482d4c6e1",
				To:              "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
				Distance:        1,
				Formula:         "x",
				PlaceholderName: "x",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err, _ = devices.SetAspect(client.InternalAdminToken, models.Aspect{
		Id:   "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6",
		Name: "Air",
		SubAspects: []models.Aspect{
			{
				Id:   "urn:infai:ses:aspect:f65876d6-77db-4123-b38b-f0eef5169378",
				Name: "Indoor",
			},
			{
				Id:   "urn:infai:ses:aspect:6eb3fb38-3592-4e2c-b2e2-5e303ddc55c0",
				Name: "Outdoor",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err, _ = devices.SetFunction(client.InternalAdminToken, models.Function{
		Id:          "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
		Name:        "Set Target Temperature",
		DisplayName: "Temperature",
		Description: "Temperature",
		ConceptId:   "urn:infai:ses:concept:0bc81398-3ed6-4e2b-a6c4-b754583aac37",
		RdfType:     "https://senergy.infai.org/ontology/ControllingFunction",
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err, _ = devices.SetFunction(client.InternalAdminToken, models.Function{
		Id:          "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
		Name:        "Get Temperature",
		DisplayName: "Get Temperature",
		Description: "Get Temperature",
		ConceptId:   "urn:infai:ses:concept:0bc81398-3ed6-4e2b-a6c4-b754583aac37",
		RdfType:     "https://senergy.infai.org/ontology/MeasuringFunction",
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err, _ = devices.SetProtocol(client.InternalAdminToken, models.Protocol{
		Id:      "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
		Name:    "connector",
		Handler: "http://unknown:88", //"connector", //to provoke an incident
		ProtocolSegments: []models.ProtocolSegment{
			{
				Id:   "urn:infai:ses:protocol-segment:9956d8b5-46fa-4381-a227-c1df69808997",
				Name: "metadata",
			},
			{
				Id:   "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
				Name: "data",
			},
		},
		Constraints: []string{"senergy_connector_local_id"},
	})
	if err != nil {
		t.Error(err)
		return
	}

	dt, err, _ := devices.SetDeviceType(client.InternalAdminToken, models.DeviceType{
		Name:          "canary-device-type",
		Description:   "used for canary service github.com/SENERGY-Platform/canary",
		DeviceClassId: "urn:infai:ses:device-class:ff64280a-58e6-4cf9-9a44-e70d3831a79d",
		Attributes:    []models.Attribute{},
		Services: []models.Service{
			{
				LocalId:     "cmd",
				Name:        "cmd",
				Description: "canary cmd service, needed to test online state by subscription",
				Interaction: models.REQUEST,
				ProtocolId:  "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
				Inputs: []models.Content{
					{
						ContentVariable: models.ContentVariable{
							Name:             "value",
							Type:             models.Integer,
							CharacteristicId: "urn:infai:ses:characteristic:a49a48fc-3a2c-4149-ac7f-1a5482d4c6e1",
							FunctionId:       "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
						},
						Serialization:     models.JSON,
						ProtocolSegmentId: "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
					},
				},
			},
			{
				LocalId:     "sensor",
				Name:        "sensor",
				Description: "canary sensor service, needed to test device data handling",
				Interaction: models.EVENT,
				ProtocolId:  "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
				Outputs: []models.Content{
					{
						ContentVariable: models.ContentVariable{
							Name:             "value",
							Type:             models.Integer,
							CharacteristicId: "urn:infai:ses:characteristic:a49a48fc-3a2c-4149-ac7f-1a5482d4c6e1",
							FunctionId:       "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
							AspectId:         "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6",
						},
						Serialization:     models.JSON,
						ProtocolSegmentId: "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
					},
				},
			},
		},
	}, client.DeviceTypeUpdateOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	sensorServiceId := ""
	cmdServiceId := ""
	for _, service := range dt.Services {
		if service.LocalId == "sensor" {
			sensorServiceId = service.Id
		}
		if service.LocalId == "cmd" {
			cmdServiceId = service.Id
		}
	}

	tokenObj, err := auth.CreateToken("test", "test")
	if err != nil {
		t.Error(err)
		return
	}
	token := tokenObj.Token

	device, err, _ := devices.CreateDevice(token, models.Device{
		LocalId:      "d1",
		Name:         "d1",
		DeviceTypeId: dt.Id,
	})

	time.Sleep(10 * time.Second)

	t.Run("test event deployment", func(t *testing.T) {
		var deplId string
		t.Run("deploy", func(t *testing.T) {
			deplId, err = integrationtest.DeployEventProcess(token, deploymentUrl, device.Id, sensorServiceId)
			if err != nil {
				t.Error(err)
				return
			}
		})
		time.Sleep(5 * time.Second)
		t.Run("trigger event", func(t *testing.T) {
			servicetopic := ServiceIdToTopic(sensorServiceId)
			err := kafka.InitTopic(conf.KafkaUrl, servicetopic)
			if err != nil {
				t.Error(err)
				return
			}
			writer := &kafkalib.Writer{
				Addr:        kafkalib.TCP(conf.KafkaUrl),
				Topic:       servicetopic,
				Async:       false,
				BatchSize:   1,
				Balancer:    &kafkalib.Hash{},
				ErrorLogger: log.New(os.Stderr, "KAFKA", 0),
			}
			defer writer.Close()

			type EventTestType struct {
				DeviceId  string                 `json:"device_id"`
				ServiceId string                 `json:"service_id"`
				Value     map[string]interface{} `json:"value"`
			}
			pl, err := json.Marshal(EventTestType{
				DeviceId:  device.Id,
				ServiceId: sensorServiceId,
				Value:     map[string]interface{}{"value": 10},
			})
			if err != nil {
				t.Error(err)
				return
			}

			err = writer.WriteMessages(ctx, kafkalib.Message{
				Key:   []byte(device.Id),
				Value: pl,
			})
			if err != nil {
				t.Error(err)
				return
			}
		})
		time.Sleep(5 * time.Second)
		t.Run("check process start", func(t *testing.T) {
			req, err := http.NewRequest("GET", conf.ProcessEngineWrapperUrl+"/v2/history/process-instances", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			var instances []HistoricProcessInstance
			err = json.NewDecoder(resp.Body).Decode(&instances)
			if err != nil {
				t.Error(err)
				return
			}
			index := slices.IndexFunc(instances, func(instance HistoricProcessInstance) bool {
				return instance.ProcessDefinitionKey == "canary_event_process"
			})
			if index == -1 {
				t.Error("no process instance found")
				return
			}
			if instances[index].State != "COMPLETED" {
				t.Error("unexpected instance state", instances[index].State)
				return
			}
		})

		t.Run("delete deployment", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, deploymentUrl+"/v3/deployments/"+url.PathEscape(deplId), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}
		})

		t.Run("check process instances deleted", func(t *testing.T) {
			req, err := http.NewRequest("GET", conf.ProcessEngineWrapperUrl+"/v2/history/process-instances", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			var instances []HistoricProcessInstance
			err = json.NewDecoder(resp.Body).Decode(&instances)
			if err != nil {
				t.Error(err)
				return
			}
			if len(instances) > 0 {
				t.Error("unexpected process instances", instances)
				return
			}
		})

		t.Run("check process deployments deleted", func(t *testing.T) {
			req, err := http.NewRequest("GET", conf.ProcessEngineWrapperUrl+"/v2/deployments", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			var deployments []interface{}
			err = json.NewDecoder(resp.Body).Decode(&deployments)
			if err != nil {
				t.Error(err)
				return
			}
			if len(deployments) > 0 {
				t.Error("unexpected process deployments", deployments)
				return
			}
		})
	})

	t.Run("test incident", func(t *testing.T) {
		var deplId string
		t.Run("deploy", func(t *testing.T) {
			deplId, err = integrationtest.DeployTaskProcess(token, deploymentUrl, device.Id, cmdServiceId)
			if err != nil {
				t.Error(err)
				return
			}
		})

		t.Run("start", func(t *testing.T) {
			req, err := http.NewRequest("GET", conf.ProcessEngineWrapperUrl+"/v2/deployments/"+url.PathEscape(deplId)+"/start", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}
		})

		time.Sleep(5 * time.Second)

		t.Run("check incidents", func(t *testing.T) {
			req, err := http.NewRequest("GET", incidentApiUrl+"/incidents", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}

			var incidents []interface{}
			err = json.NewDecoder(resp.Body).Decode(&incidents)
			if err != nil {
				t.Error(err)
				return
			}

			if len(incidents) != 1 {
				t.Error("unexpected incidents", incidents)
				return
			}
			t.Logf("incidents: %#v", incidents)
		})

		t.Run("delete deployment", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, deploymentUrl+"/v3/deployments/"+url.PathEscape(deplId), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}
		})

		t.Run("check deleted incidents", func(t *testing.T) {
			req, err := http.NewRequest("GET", incidentApiUrl+"/incidents", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}

			var incidents []interface{}
			err = json.NewDecoder(resp.Body).Decode(&incidents)
			if err != nil {
				t.Error(err)
				return
			}

			if len(incidents) != 0 {
				t.Error("unexpected incidents", incidents)
				return
			}
			t.Logf("incidents: %#v", incidents)
		})

	})

	t.Run("test incident handling", func(t *testing.T) {
		var deplId string
		t.Run("deploy", func(t *testing.T) {
			deplId, err = integrationtest.DeployRestartProcess(token, deploymentUrl, device.Id, cmdServiceId)
			if err != nil {
				t.Error(err)
				return
			}
		})

		t.Run("start", func(t *testing.T) {
			req, err := http.NewRequest("GET", conf.ProcessEngineWrapperUrl+"/v2/deployments/"+url.PathEscape(deplId)+"/start", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}
		})

		time.Sleep(1 * time.Second)

		var oldIncidentsLen int
		t.Run("check incidents", func(t *testing.T) {
			req, err := http.NewRequest("GET", incidentApiUrl+"/incidents", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}

			var incidents []interface{}
			err = json.NewDecoder(resp.Body).Decode(&incidents)
			if err != nil {
				t.Error(err)
				return
			}
			oldIncidentsLen = len(incidents)
			if len(incidents) < 2 {
				t.Error("unexpected incidents", len(incidents), incidents)
				return
			}
			t.Logf("incidents: %v %#v", len(incidents), incidents)
		})

		t.Run("delete deployment", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, deploymentUrl+"/v3/deployments/"+url.PathEscape(deplId), nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}
		})

		time.Sleep(2 * time.Second)

		t.Run("check deleted incidents", func(t *testing.T) {
			req, err := http.NewRequest("GET", incidentApiUrl+"/incidents", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Error("unexpected status code", resp.StatusCode)
				return
			}

			var incidents []interface{}
			err = json.NewDecoder(resp.Body).Decode(&incidents)
			if err != nil {
				t.Error(err)
				return
			}

			//because a race condition, some incidents may remain. we accept that currently because they can be deleted manually
			if !(len(incidents) < oldIncidentsLen) || len(incidents) > 2 {
				t.Error("unexpected incidents", len(incidents), incidents)
				return
			}
			t.Logf("incidents: %v %#v", len(incidents), incidents)
		})

	})
}

func ServiceIdToTopic(id string) string {
	id = strings.ReplaceAll(id, "#", "_")
	id = strings.ReplaceAll(id, ":", "_")
	return id
}

type HistoricProcessInstance struct {
	Id                       string  `json:"id"`
	SuperProcessInstanceId   string  `json:"superProcessInstanceId"`
	SuperCaseInstanceId      string  `json:"superCaseInstanceId"`
	CaseInstanceId           string  `json:"caseInstanceId"`
	ProcessDefinitionName    string  `json:"processDefinitionName"`
	ProcessDefinitionKey     string  `json:"processDefinitionKey"`
	ProcessDefinitionVersion float64 `json:"processDefinitionVersion"`
	ProcessDefinitionId      string  `json:"processDefinitionId"`
	BusinessKey              string  `json:"businessKey"`
	StartTime                string  `json:"startTime"`
	EndTime                  string  `json:"endTime"`
	DurationInMillis         float64 `json:"durationInMillis"`
	StartUserId              string  `json:"startUserId"`
	StartActivityId          string  `json:"startActivityId"`
	DeleteReason             string  `json:"deleteReason"`
	TenantId                 string  `json:"tenantId"`
	State                    string  `json:"state"`
}
