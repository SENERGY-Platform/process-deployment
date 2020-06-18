/*
 * Copyright 2019 InfAI (CC SES)
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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"github.com/SENERGY-Platform/process-deployment/lib/tests"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ExampleKeepRetries() {
	conf, err := config.LoadConfig("../config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	port, err := tests.GetFreePort()
	if err != nil {
		fmt.Println(err)
		return
	}
	conf.ApiPort = strconv.Itoa(port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = Start(ctx, conf, mocks.Kafka, mocks.Database, mocks.Devices, mocks.ProcessModelRepo)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(100 * time.Millisecond) //wait for api startup
	prepareMockRepos()

	file, err := ioutil.ReadFile("./tests/resources/retry.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	preparereq := messages.PrepareRequest{Xml: string(file), Svg: "<svg/>"}
	temp, err := json.Marshal(preparereq)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:"+conf.ApiPort+"/prepared-deployments",
		bytes.NewBuffer(temp),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", TestToken)

	log.Println("request prepared deployment")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	deployment := deploymentmodel.Deployment{}

	err = json.NewDecoder(resp.Body).Decode(&deployment)

	if err != nil {
		fmt.Println(err)
		return
	}
	deployment.XmlRaw = ""
	deployment.Xml = ""

	msg, err := json.Marshal(deployment)

	fmt.Println(err, string(msg))

	//deployment
	tempNewId := config.NewId
	defer func() {
		config.NewId = tempNewId
	}()
	config.NewId = func() string {
		return "uuid"
	}

	createSimpleObj := deploymentmodel.Deployment{}
	err = json.Unmarshal(msg, &createSimpleObj)
	if err != nil {
		fmt.Println(err)
		return
	}
	createSimpleObj.Id = "id-example"
	createSimpleObj.Name = "simple"
	createSimpleObj.XmlRaw = string(file)
	if createSimpleObj.Elements[0].Task.Selection.Device == nil {
		createSimpleObj.Elements[0].Task.Selection.Device = &devicemodel.Device{}
	}
	if createSimpleObj.Elements[0].Task.Selection.Service == nil {
		createSimpleObj.Elements[0].Task.Selection.Service = &devicemodel.Service{}
	}
	createSimpleObj.Elements[0].Task.Selection.Device.Id = "device_id_1"
	createSimpleObj.Elements[0].Task.Selection.Service.Id = "service_id_1"
	createSimpleObj.Elements[0].Task.Selection.Service.ProtocolId = "pid"

	createSimpleObj.Elements[1].MultiTask.Selections = append(createSimpleObj.Elements[1].MultiTask.Selections,
		deploymentmodel.Selection{
			Device:  &devicemodel.Device{Id: "device_id_1"},
			Service: &devicemodel.Service{Id: "service_id_1", ProtocolId: "pid"},
		}, deploymentmodel.Selection{
			Device:  &devicemodel.Device{Id: "device_id_2"},
			Service: &devicemodel.Service{Id: "service_id_2", ProtocolId: "pid"},
		})

	if createSimpleObj.Elements[2].MsgEvent.Device == nil {
		createSimpleObj.Elements[2].MsgEvent.Device = &devicemodel.Device{}
	}
	if createSimpleObj.Elements[2].MsgEvent.Service == nil {
		createSimpleObj.Elements[2].MsgEvent.Service = &devicemodel.Service{}
	}
	createSimpleObj.Elements[2].MsgEvent.Device.Id = "device_id_1"
	createSimpleObj.Elements[2].MsgEvent.Service.Id = "service_id_2"
	createSimpleObj.Elements[2].MsgEvent.Service.ProtocolId = "pid"
	createSimpleObj.Elements[2].MsgEvent.Path = "$.value+"
	createSimpleObj.Elements[2].MsgEvent.Value = "42"
	createSimpleObj.Elements[2].MsgEvent.Operation = "operation"

	createSimple, err := json.Marshal(createSimpleObj)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("CREATE:")
	err = subExampleDeploymentCreate(conf, createSimple)
	if err != nil {
		fmt.Println(err)
		return
	}

	updateSimpleObj := deploymentmodel.Deployment{}
	err = json.Unmarshal(createSimple, &updateSimpleObj)
	if err != nil {
		fmt.Println(err)
		return
	}
	updateSimpleObj.Id = config.NewId()
	updateSimpleObj.Name = "somthing else"
	updateSimpleObj.Elements[0].Task.Selection.Device.Id = "device_id_2"

	fmt.Println("UPDATE:")
	err = subExampleDeploymentUpdate(conf, updateSimpleObj)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("LIST:")
	err = subExampleGetDependenciesList(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = subExampleGetSelectedDependencies(conf, []string{config.NewId()})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = subExampleGetSelectedDependencies(conf, []string{config.NewId(), "expect_error"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("DELETE:")
	err = subExampleDeploymentDelete(conf, updateSimpleObj.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	deploymentMsgs := mocks.Kafka.GetProduced(conf.DeploymentTopic)

	fmt.Println("DEPLOYMENTS:", len(deploymentMsgs))
	for _, msg := range deploymentMsgs {
		fmt.Println(msg)
	}

	//output:
	//<nil> {"id":"","executable":true,"xml_raw":"","xml":"","svg":"\u003csvg/\u003e","name":"simple","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","retries":3,"input":"000","selectables":[{"device":{"id":"device1"},"services":[{"id":"service1"}]}],"selection":{"device":null,"service":null},"parameter":{"inputs":"\"ff0\""},"configurables":null}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","retries":2,"input":"000","selectables":[{"device":{"id":"device1"},"services":[{"id":"service1"}]}],"selections":null,"parameter":{"inputs":"\"fff\""},"configurables":null}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":null,"service":null,"path":"","value":"","operation":"","event_id":""}}],"lanes":null}
	//CREATE:
	//<nil> 200{"deployment_id":"uuid","owner":"connectivity-test","devices":[{"device_id":"device_id_1","name":"d1","bpmn_resources":[{"id":"IntermediateThrowEvent_0905jg5"},{"id":"Task_096xjeg"},{"id":"Task_0wjr1fj"}]},{"device_id":"device_id_2","name":"d2","bpmn_resources":[{"id":"Task_096xjeg"}]}],"events":[{"event_id":"uuid","bpmn_resources":[{"id":"IntermediateThrowEvent_0905jg5"}]}]}
	//UPDATE:
	//<nil> 200{"deployment_id":"uuid","owner":"connectivity-test","devices":[{"device_id":"device_id_1","name":"d1","bpmn_resources":[{"id":"IntermediateThrowEvent_0905jg5"},{"id":"Task_096xjeg"}]},{"device_id":"device_id_2","name":"d2","bpmn_resources":[{"id":"Task_096xjeg"},{"id":"Task_0wjr1fj"}]}],"events":[{"event_id":"uuid","bpmn_resources":[{"id":"IntermediateThrowEvent_0905jg5"}]}]}
	//LIST:
	//<nil> 200[{"deployment_id":"uuid","owner":"connectivity-test","devices":[{"device_id":"device_id_1","name":"d1","bpmn_resources":[{"id":"IntermediateThrowEvent_0905jg5"},{"id":"Task_096xjeg"}]},{"device_id":"device_id_2","name":"d2","bpmn_resources":[{"id":"Task_096xjeg"},{"id":"Task_0wjr1fj"}]}],"events":[{"event_id":"uuid","bpmn_resources":[{"id":"IntermediateThrowEvent_0905jg5"}]}]}]
	//<nil> 200[{"deployment_id":"uuid","owner":"connectivity-test","devices":[{"device_id":"device_id_1","name":"d1","bpmn_resources":[{"id":"IntermediateThrowEvent_0905jg5"},{"id":"Task_096xjeg"}]},{"device_id":"device_id_2","name":"d2","bpmn_resources":[{"id":"Task_096xjeg"},{"id":"Task_0wjr1fj"}]}],"events":[{"event_id":"uuid","bpmn_resources":[{"id":"IntermediateThrowEvent_0905jg5"}]}]}]
	//<nil> 404unknown id
	//DELETE:
	//<nil> 404dependencies not found
	//DEPLOYMENTS: 3
	//{"command":"PUT","id":"uuid","owner":"connectivity-test","deployment":{"id":"uuid","executable":true,"xml_raw":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\"\n                  xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\"\n                  xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\"\n                  xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\"\n                  targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\"\n                           targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\"\n                           targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid2\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;,\n                        \"retries\": 2\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;,\n                        \"retries\": 3\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e","xml":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\" xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\" xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\" xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\" targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\" name=\"simple\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\" targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\" targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;fid2\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;concept_id\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;characteristic_id\u0026quot;:\u0026quot;example_hex\u0026quot;,\u0026quot;input\u0026quot;:\u0026quot;000\u0026quot;,\u0026quot;retries\u0026quot;:2}\u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003ccamunda:executionListener event=\"start\"\u003e\u003ccamunda:script scriptFormat=\"groovy\"\u003eexecution.setVariable(\u0026quot;collection\u0026quot;, [\u0026quot;{\\\u0026quot;device_id\\\u0026quot;:\\\u0026quot;device_id_1\\\u0026quot;,\\\u0026quot;service_id\\\u0026quot;:\\\u0026quot;service_id_1\\\u0026quot;,\\\u0026quot;protocol_id\\\u0026quot;:\\\u0026quot;pid\\\u0026quot;}\u0026quot;,\u0026quot;{\\\u0026quot;device_id\\\u0026quot;:\\\u0026quot;device_id_2\\\u0026quot;,\\\u0026quot;service_id\\\u0026quot;:\\\u0026quot;service_id_2\\\u0026quot;,\\\u0026quot;protocol_id\\\u0026quot;:\\\u0026quot;pid\\\u0026quot;}\u0026quot;])\u003c/camunda:script\u003e\u003c/camunda:executionListener\u003e\u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\" camunda:collection=\"collection\" camunda:elementVariable=\"overwrite\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;fid\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;concept_id\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;characteristic_id\u0026quot;:\u0026quot;example_hex\u0026quot;,\u0026quot;device_id\u0026quot;:\u0026quot;device_id_1\u0026quot;,\u0026quot;service_id\u0026quot;:\u0026quot;service_id_1\u0026quot;,\u0026quot;protocol_id\u0026quot;:\u0026quot;pid\u0026quot;,\u0026quot;input\u0026quot;:\u0026quot;000\u0026quot;,\u0026quot;retries\u0026quot;:3}\u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition messageRef=\"e_uuid\"/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmn:message id=\"e_uuid\" name=\"uuid\"/\u003e\u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e","svg":"\u003csvg/\u003e","name":"simple","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","retries":3,"input":"000","selectables":[{"device":{"id":"device1"},"services":[{"id":"service1"}]}],"selection":{"device":{"id":"device_id_1","local_id":"d1url","name":"d1"},"service":{"id":"service_id_1","local_id":"s1url","name":"s1","protocol_id":"pid"}},"parameter":{"inputs":"\"ff0\""},"configurables":null}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","retries":2,"input":"000","selectables":[{"device":{"id":"device1"},"services":[{"id":"service1"}]}],"selections":[{"device":{"id":"device_id_1","local_id":"d1url","name":"d1"},"service":{"id":"service_id_1","local_id":"s1url","name":"s1","protocol_id":"pid"}},{"device":{"id":"device_id_2","local_id":"d2url","name":"d2"},"service":{"id":"service_id_2","local_id":"s2url","name":"s2","protocol_id":"pid"}}],"parameter":{"inputs":"\"fff\""},"configurables":null}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":{"id":"device_id_1","local_id":"d1url","name":"d1"},"service":{"id":"service_id_2","local_id":"s2url","name":"s2","protocol_id":"pid"},"path":"$.value+","value":"42","operation":"operation","event_id":"uuid"}}],"lanes":null}}
	//{"command":"PUT","id":"uuid","owner":"connectivity-test","deployment":{"id":"uuid","executable":true,"xml_raw":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\"\n                  xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\"\n                  xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\"\n                  xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\"\n                  targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\"\n                           targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\"\n                           targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid2\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;,\n                        \"retries\": 2\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;,\n                        \"retries\": 3\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e","xml":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\" xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\" xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\" xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\" targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\" name=\"somthing else\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\" targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\" targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;fid2\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;concept_id\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;characteristic_id\u0026quot;:\u0026quot;example_hex\u0026quot;,\u0026quot;input\u0026quot;:\u0026quot;000\u0026quot;,\u0026quot;retries\u0026quot;:2}\u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003ccamunda:executionListener event=\"start\"\u003e\u003ccamunda:script scriptFormat=\"groovy\"\u003eexecution.setVariable(\u0026quot;collection\u0026quot;, [\u0026quot;{\\\u0026quot;device_id\\\u0026quot;:\\\u0026quot;device_id_1\\\u0026quot;,\\\u0026quot;service_id\\\u0026quot;:\\\u0026quot;service_id_1\\\u0026quot;,\\\u0026quot;protocol_id\\\u0026quot;:\\\u0026quot;pid\\\u0026quot;}\u0026quot;,\u0026quot;{\\\u0026quot;device_id\\\u0026quot;:\\\u0026quot;device_id_2\\\u0026quot;,\\\u0026quot;service_id\\\u0026quot;:\\\u0026quot;service_id_2\\\u0026quot;,\\\u0026quot;protocol_id\\\u0026quot;:\\\u0026quot;pid\\\u0026quot;}\u0026quot;])\u003c/camunda:script\u003e\u003c/camunda:executionListener\u003e\u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\" camunda:collection=\"collection\" camunda:elementVariable=\"overwrite\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;fid\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;concept_id\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;characteristic_id\u0026quot;:\u0026quot;example_hex\u0026quot;,\u0026quot;device_id\u0026quot;:\u0026quot;device_id_2\u0026quot;,\u0026quot;service_id\u0026quot;:\u0026quot;service_id_1\u0026quot;,\u0026quot;protocol_id\u0026quot;:\u0026quot;pid\u0026quot;,\u0026quot;input\u0026quot;:\u0026quot;000\u0026quot;,\u0026quot;retries\u0026quot;:3}\u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition messageRef=\"e_uuid\"/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmn:message id=\"e_uuid\" name=\"uuid\"/\u003e\u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e","svg":"\u003csvg/\u003e","name":"somthing else","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","retries":3,"input":"000","selectables":[{"device":{"id":"device1"},"services":[{"id":"service1"}]}],"selection":{"device":{"id":"device_id_2","local_id":"d2url","name":"d2"},"service":{"id":"service_id_1","local_id":"s1url","name":"s1","protocol_id":"pid"}},"parameter":{"inputs":"\"ff0\""},"configurables":null}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","retries":2,"input":"000","selectables":[{"device":{"id":"device1"},"services":[{"id":"service1"}]}],"selections":[{"device":{"id":"device_id_1","local_id":"d1url","name":"d1"},"service":{"id":"service_id_1","local_id":"s1url","name":"s1","protocol_id":"pid"}},{"device":{"id":"device_id_2","local_id":"d2url","name":"d2"},"service":{"id":"service_id_2","local_id":"s2url","name":"s2","protocol_id":"pid"}}],"parameter":{"inputs":"\"fff\""},"configurables":null}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":{"id":"device_id_1","local_id":"d1url","name":"d1"},"service":{"id":"service_id_2","local_id":"s2url","name":"s2","protocol_id":"pid"},"path":"$.value+","value":"42","operation":"operation","event_id":"uuid"}}],"lanes":null}}
	//{"command":"DELETE","id":"uuid","owner":"connectivity-test","deployment":null}
}
