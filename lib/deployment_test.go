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
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/tests"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

func ExampleCtrl_Deployment() {

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

	port, err := tests.GetFreePort()
	if err != nil {
		fmt.Println(err)
		return
	}
	conf.ApiPort = strconv.Itoa(port)

	ctx, _ := context.WithCancel(context.Background())

	//defer cancel()
	//defer log.Println("stop")

	err = Start(ctx, conf, mocks.Kafka, mocks.Database, mocks.Connectionlog, mocks.SemanticRepository, mocks.DeviceRepository, mocks.ProcessRepository)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(100 * time.Millisecond) //wait for api startup

	prepareMockRepos()

	createSimple, err := ioutil.ReadFile("./tests/resources/simple.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = subExampleDeploymentCreate(conf, createSimple)
	if err != nil {
		fmt.Println(err)
		return
	}

	updateSimpleObj := model.Deployment{}
	err = json.Unmarshal(createSimple, &updateSimpleObj)
	if err != nil {
		fmt.Println(err)
		return
	}
	updateSimpleObj.Id = config.NewId()
	updateSimpleObj.Name = "somthing else"
	updateSimpleObj.Elements[0].Task.Selection.SelectedDevice.Id = "device_id_2"

	err = subExampleDeploymentUpdate(conf, updateSimpleObj)
	if err != nil {
		fmt.Println(err)
		return
	}

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
	//{uuid connectivity-test [{device_id_1 d1 [{Task_0wjr1fj } {Task_096xjeg } {IntermediateThrowEvent_0905jg5 }]} {device_id_2 d2 [{Task_096xjeg }]}] [{uuid [{IntermediateThrowEvent_0905jg5 }]}] false} <nil> 200
	//{uuid connectivity-test [{device_id_2 d2 [{Task_0wjr1fj } {Task_096xjeg }]} {device_id_1 d1 [{Task_096xjeg } {IntermediateThrowEvent_0905jg5 }]}] [{uuid [{IntermediateThrowEvent_0905jg5 }]}] false} <nil> 200
	//{  [] [] false} dependencies not found 404
	//DEPLOYMENTS: 3
	//{"command":"PUT","id":"uuid","owner":"connectivity-test","deployment":{"id":"uuid","xml_raw":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\"\n                  xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\"\n                  xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\"\n                  xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\"\n                  targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\"\n                           targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\"\n                           targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid2\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e","xml":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\" xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\" xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\" xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\" targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\" targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\" targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;fid2\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;concept_id\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;characteristic_id\u0026quot;:\u0026quot;example_hex\u0026quot;,\u0026quot;input\u0026quot;:\u0026quot;000\u0026quot;}\u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003ccamunda:executionListener event=\"start\"\u003e\u003ccamunda:script scriptFormat=\"groovy\"\u003eexecution.setVariable(\u0026quot;collection\u0026quot;, [\u0026quot;{\\\u0026quot;device_id\\\u0026quot;:\\\u0026quot;device_id_1\\\u0026quot;,\\\u0026quot;service_id\\\u0026quot;:\\\u0026quot;service_id_1\\\u0026quot;,\\\u0026quot;protocol_id\\\u0026quot;:\\\u0026quot;pid\\\u0026quot;}\u0026quot;,\u0026quot;{\\\u0026quot;device_id\\\u0026quot;:\\\u0026quot;device_id_2\\\u0026quot;,\\\u0026quot;service_id\\\u0026quot;:\\\u0026quot;service_id_2\\\u0026quot;,\\\u0026quot;protocol_id\\\u0026quot;:\\\u0026quot;pid\\\u0026quot;}\u0026quot;])\u003c/camunda:script\u003e\u003c/camunda:executionListener\u003e\u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\" camunda:collection=\"collection\" camunda:elementVariable=\"overwrite\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;fid\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;concept_id\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;characteristic_id\u0026quot;:\u0026quot;example_hex\u0026quot;,\u0026quot;device\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;device_id_1\u0026quot;,\u0026quot;local_id\u0026quot;:\u0026quot;d1url\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;d1\u0026quot;,\u0026quot;device_type_id\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;service\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;service_id_1\u0026quot;,\u0026quot;local_id\u0026quot;:\u0026quot;s1url\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;s1\u0026quot;,\u0026quot;description\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;aspects\u0026quot;:null,\u0026quot;protocol_id\u0026quot;:\u0026quot;pid\u0026quot;,\u0026quot;inputs\u0026quot;:null,\u0026quot;outputs\u0026quot;:null,\u0026quot;functions\u0026quot;:null,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;protocol\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;pid\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;protocol1\u0026quot;,\u0026quot;handler\u0026quot;:\u0026quot;p\u0026quot;,\u0026quot;protocol_segments\u0026quot;:null},\u0026quot;input\u0026quot;:\u0026quot;000\u0026quot;}\u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition messageRef=\"e_uuid\"/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmn:message id=\"e_uuid\" name=\"uuid\"/\u003e\u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e","name":"simple","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","input":"000","device_options":null,"selection":{"selected_device":{"id":"device_id_1","local_id":"d1url","name":"d1","device_type_id":""},"selected_service":{"id":"service_id_1","local_id":"s1url","name":"s1","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""}},"parameter":{"inputs":"\"ff0\""}}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","input":"000","device_options":null,"selections":[{"selected_device":{"id":"device_id_1","local_id":"d1url","name":"d1","device_type_id":""},"selected_service":{"id":"service_id_1","local_id":"s1url","name":"s1","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""}},{"selected_device":{"id":"device_id_2","local_id":"d2url","name":"d2","device_type_id":""},"selected_service":{"id":"service_id_2","local_id":"s2url","name":"s2","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""}}],"parameter":{"inputs":"\"fff\""}}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":{"id":"device_id_1","local_id":"d1url","name":"d1","device_type_id":""},"service":{"id":"service_id_2","local_id":"s2url","name":"s2","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"$.value+","value":"42","operation":"operation","event_id":"uuid"}}],"lanes":null}}
	//{"command":"PUT","id":"uuid","owner":"connectivity-test","deployment":{"id":"uuid","xml_raw":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\"\n                  xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\"\n                  xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\"\n                  xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\"\n                  targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\"\n                           targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\"\n                           targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid2\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e","xml":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\" xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\" xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\" xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\" targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\" targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\" targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;fid2\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;concept_id\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;characteristic_id\u0026quot;:\u0026quot;example_hex\u0026quot;,\u0026quot;input\u0026quot;:\u0026quot;000\u0026quot;}\u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003ccamunda:executionListener event=\"start\"\u003e\u003ccamunda:script scriptFormat=\"groovy\"\u003eexecution.setVariable(\u0026quot;collection\u0026quot;, [\u0026quot;{\\\u0026quot;device_id\\\u0026quot;:\\\u0026quot;device_id_1\\\u0026quot;,\\\u0026quot;service_id\\\u0026quot;:\\\u0026quot;service_id_1\\\u0026quot;,\\\u0026quot;protocol_id\\\u0026quot;:\\\u0026quot;pid\\\u0026quot;}\u0026quot;,\u0026quot;{\\\u0026quot;device_id\\\u0026quot;:\\\u0026quot;device_id_2\\\u0026quot;,\\\u0026quot;service_id\\\u0026quot;:\\\u0026quot;service_id_2\\\u0026quot;,\\\u0026quot;protocol_id\\\u0026quot;:\\\u0026quot;pid\\\u0026quot;}\u0026quot;])\u003c/camunda:script\u003e\u003c/camunda:executionListener\u003e\u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\" camunda:collection=\"collection\" camunda:elementVariable=\"overwrite\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;fid\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;concept_id\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;characteristic_id\u0026quot;:\u0026quot;example_hex\u0026quot;,\u0026quot;device\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;device_id_2\u0026quot;,\u0026quot;local_id\u0026quot;:\u0026quot;d2url\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;d2\u0026quot;,\u0026quot;device_type_id\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;service\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;service_id_1\u0026quot;,\u0026quot;local_id\u0026quot;:\u0026quot;s1url\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;s1\u0026quot;,\u0026quot;description\u0026quot;:\u0026quot;\u0026quot;,\u0026quot;aspects\u0026quot;:null,\u0026quot;protocol_id\u0026quot;:\u0026quot;pid\u0026quot;,\u0026quot;inputs\u0026quot;:null,\u0026quot;outputs\u0026quot;:null,\u0026quot;functions\u0026quot;:null,\u0026quot;rdf_type\u0026quot;:\u0026quot;\u0026quot;},\u0026quot;protocol\u0026quot;:{\u0026quot;id\u0026quot;:\u0026quot;pid\u0026quot;,\u0026quot;name\u0026quot;:\u0026quot;protocol1\u0026quot;,\u0026quot;handler\u0026quot;:\u0026quot;p\u0026quot;,\u0026quot;protocol_segments\u0026quot;:null},\u0026quot;input\u0026quot;:\u0026quot;000\u0026quot;}\u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition messageRef=\"e_uuid\"/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmn:message id=\"e_uuid\" name=\"uuid\"/\u003e\u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e","name":"somthing else","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","input":"000","device_options":null,"selection":{"selected_device":{"id":"device_id_2","local_id":"d2url","name":"d2","device_type_id":""},"selected_service":{"id":"service_id_1","local_id":"s1url","name":"s1","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""}},"parameter":{"inputs":"\"ff0\""}}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","input":"000","device_options":null,"selections":[{"selected_device":{"id":"device_id_1","local_id":"d1url","name":"d1","device_type_id":""},"selected_service":{"id":"service_id_1","local_id":"s1url","name":"s1","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""}},{"selected_device":{"id":"device_id_2","local_id":"d2url","name":"d2","device_type_id":""},"selected_service":{"id":"service_id_2","local_id":"s2url","name":"s2","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""}}],"parameter":{"inputs":"\"fff\""}}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":{"id":"device_id_1","local_id":"d1url","name":"d1","device_type_id":""},"service":{"id":"service_id_2","local_id":"s2url","name":"s2","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"$.value+","value":"42","operation":"operation","event_id":"uuid"}}],"lanes":null}}
	//{"command":"DELETE","id":"uuid","owner":"","deployment":{"id":"","xml_raw":"","xml":"","name":"","elements":null,"lanes":null}}
}

func prepareMockRepos() {
	mocks.DeviceRepository.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})
	mocks.SemanticRepository.SetOptions([]model.DeviceOption{
		{
			Device: devicemodel.Device{
				Id: "device1",
			},
			ServiceOptions: []devicemodel.Service{
				{
					Id: "service1",
				},
			},
		},
	})
	mocks.DeviceRepository.SetDevice("device_id_1", devicemodel.Device{
		Id:      "device_id_1",
		LocalId: "d1url",
		Name:    "d1",
	})
	mocks.DeviceRepository.SetDevice("device_id_2", devicemodel.Device{
		Id:      "device_id_2",
		LocalId: "d2url",
		Name:    "d2",
	})
	mocks.DeviceRepository.SetService("service_id_1", devicemodel.Service{
		Id:         "service_id_1",
		LocalId:    "s1url",
		Name:       "s1",
		ProtocolId: "pid",
	})
	mocks.DeviceRepository.SetService("service_id_2", devicemodel.Service{
		Id:         "service_id_2",
		LocalId:    "s2url",
		Name:       "s2",
		ProtocolId: "pid",
	})
}

func subExampleDeploymentCreate(conf config.Config, deploymentJson []byte) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:"+conf.ApiPort+"/deployments",
		bytes.NewBuffer(deploymentJson),
	)
	if err != nil {
		debug.PrintStack()
		return err
	}
	req.Header.Set("Authorization", TestToken)

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		debug.PrintStack()
		return errors.New("unexpected statuscode")
	}

	fmt.Println(mocks.Database.GetDependencies("connectivity-test", config.NewId()))
	return nil
}

func subExampleDeploymentUpdate(conf config.Config, deployment model.Deployment) error {
	deploymentJson, err := json.Marshal(deployment)
	if err != nil {
		debug.PrintStack()
		return err
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(
		"PUT",
		"http://localhost:"+conf.ApiPort+"/deployments/"+deployment.Id,
		bytes.NewBuffer(deploymentJson),
	)
	if err != nil {
		debug.PrintStack()
		return err
	}
	req.Header.Set("Authorization", TestToken)

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		debug.PrintStack()
		return errors.New("unexpected statuscode")
	}

	fmt.Println(mocks.Database.GetDependencies("connectivity-test", deployment.Id))
	return nil
}

func subExampleDeploymentDelete(conf config.Config, id string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(
		"DELETE",
		"http://localhost:"+conf.ApiPort+"/deployments/"+id,
		nil,
	)
	if err != nil {
		debug.PrintStack()
		return err
	}
	req.Header.Set("Authorization", TestToken)

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		debug.PrintStack()
		return errors.New("unexpected statuscode")
	}

	fmt.Println(mocks.Database.GetDependencies("connectivity-test", id))
	return nil
}
