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

package bpmn

import (
	"encoding/json"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"io/ioutil"
)

func ExampleSimpleBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/simple.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	result.XmlRaw = ""
	temp, err := json.Marshal(result)
	fmt.Println(err, string(temp))

	//output:
	//<nil> {"id":"","xml_raw":"","xml":"","name":"simple","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","input":"000","device_options":null,"selection":{"selected_device":{"id":"","local_id":"","name":"","device_type_id":""},"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}},"parameter":{"inputs":"\"ff0\""}}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","input":"000","device_options":null,"selections":null,"parameter":{"inputs":"\"fff\""}}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}}],"lanes":null}
}

func ExampleLaneBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/lanes.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	result.XmlRaw = ""
	temp, err := json.Marshal(result)
	fmt.Println(err, string(temp))

	//output:
	//<nil> {"id":"","xml_raw":"","xml":"","name":"lanes","elements":null,"lanes":[{"order":0,"multi_lane":{"label":"multiTaskLane","bpmn_element_id":"Lane_12774cv","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"device_options":null,"selections":null,"elements":[{"order":3,"task":{"label":"multi_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_084s3g5","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}},{"order":4,"task":{"label":"multi_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_098jmqp","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}}]}},{"order":0,"lane":{"label":"MixedLane","bpmn_element_id":"Lane_0odlj5k","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"device_options":null,"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":5,"task":{"label":"mixed_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1npvonw","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}},{"order":6,"task":{"label":"mixed_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1mnjsed","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}}]}},{"order":1,"lane":{"label":"taskLane","bpmn_element_id":"Lane_0v679jg","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}}],"device_options":null,"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":1,"task":{"label":"lane_task_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_0nmb2on","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}},{"order":2,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_1tchutl","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}}]}}]}
}

func ExampleTimeAndReceiverBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/timeAndReceive.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	result.XmlRaw = ""
	temp, err := json.Marshal(result)
	fmt.Println(err, string(temp))

	//output:
	//<nil> {"id":"","xml_raw":"","xml":"","name":"timeAndReceive","elements":[{"order":0,"receive_task_event":{"label":"Task_1uyyxb0","bpmn_element_id":"Task_1uyyxb0","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}},{"order":0,"time_event":{"bpmn_element_id":"IntermediateThrowEvent_10mhx3e","kind":"timeDuration","time":"PT1M","label":"eine Minute"}}],"lanes":null}
}

func ExampleSimpleBpmnDeploymentToXml() {
	deploymentStr := `{
   "id":"uuid",
   "xml_raw":"\u003cbpmn:definitions xmlns:bpmn=\"http://www.omg.org/spec/BPMN/20100524/MODEL\"\n                  xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\"\n                  xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:camunda=\"http://camunda.org/schema/1.0/bpmn\"\n                  xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" id=\"Definitions_1\"\n                  targetNamespace=\"http://bpmn.io/schema/bpmn\"\u003e\n    \u003cbpmn:process id=\"simple\" isExecutable=\"true\"\u003e\n        \u003cbpmn:startEvent id=\"StartEvent_1\"\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0ixns30\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:startEvent\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0ixns30\" sourceRef=\"StartEvent_1\" targetRef=\"Task_096xjeg\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0htq2f6\" sourceRef=\"Task_096xjeg\"\n                           targetRef=\"IntermediateThrowEvent_0905jg5\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_0npfu5a\" sourceRef=\"IntermediateThrowEvent_0905jg5\"\n                           targetRef=\"Task_0wjr1fj\"/\u003e\n        \u003cbpmn:sequenceFlow id=\"SequenceFlow_1a1qwlk\" sourceRef=\"Task_0wjr1fj\" targetRef=\"EndEvent_0yi4y22\"/\u003e\n        \u003cbpmn:endEvent id=\"EndEvent_0yi4y22\"\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_1a1qwlk\u003c/bpmn:incoming\u003e\n        \u003c/bpmn:endEvent\u003e\n        \u003cbpmn:serviceTask id=\"Task_096xjeg\" name=\"multiTaskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 2}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid2\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;fff\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0ixns30\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0htq2f6\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:multiInstanceLoopCharacteristics isSequential=\"true\"/\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:serviceTask id=\"Task_0wjr1fj\" name=\"taskLabel\" camunda:type=\"external\" camunda:topic=\"task\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 1}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:extensionElements\u003e\n                \u003ccamunda:inputOutput\u003e\n                    \u003ccamunda:inputParameter name=\"payload\"\u003e{\n                        \u0026quot;function\u0026quot;:{\u0026quot;id\u0026quot;: \u0026quot;fid\u0026quot;},\n                        \u0026quot;characteristic_id\u0026quot;: \u0026quot;example_hex\u0026quot;,\n                        \u0026quot;input\u0026quot;: \u0026quot;000\u0026quot;\n                        }\n                    \u003c/camunda:inputParameter\u003e\n                    \u003ccamunda:inputParameter name=\"inputs\"\u003e\u0026quot;ff0\u0026quot;\u003c/camunda:inputParameter\u003e\n                \u003c/camunda:inputOutput\u003e\n            \u003c/bpmn:extensionElements\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0npfu5a\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_1a1qwlk\u003c/bpmn:outgoing\u003e\n        \u003c/bpmn:serviceTask\u003e\n        \u003cbpmn:intermediateCatchEvent id=\"IntermediateThrowEvent_0905jg5\" name=\"eventName\"\u003e\n            \u003cbpmn:documentation\u003e{\u0026quot;order\u0026quot;: 3}\u003c/bpmn:documentation\u003e\n            \u003cbpmn:incoming\u003eSequenceFlow_0htq2f6\u003c/bpmn:incoming\u003e\n            \u003cbpmn:outgoing\u003eSequenceFlow_0npfu5a\u003c/bpmn:outgoing\u003e\n            \u003cbpmn:messageEventDefinition/\u003e\n        \u003c/bpmn:intermediateCatchEvent\u003e\n    \u003c/bpmn:process\u003e\n    \u003cbpmndi:BPMNDiagram id=\"BPMNDiagram_1\"\u003e\n        \u003cbpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"simple\"\u003e\n            \u003cbpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\"\u003e\n                \u003cdc:Bounds x=\"173\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"EndEvent_0yi4y22_di\" bpmnElement=\"EndEvent_0yi4y22\"\u003e\n                \u003cdc:Bounds x=\"645\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_1x4556d_di\" bpmnElement=\"Task_096xjeg\"\u003e\n                \u003cdc:Bounds x=\"259\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"IntermediateCatchEvent_072g4ud_di\" bpmnElement=\"IntermediateThrowEvent_0905jg5\"\u003e\n                \u003cdc:Bounds x=\"409\" y=\"102\" width=\"36\" height=\"36\"/\u003e\n                \u003cbpmndi:BPMNLabel\u003e\n                    \u003cdc:Bounds x=\"399\" y=\"145\" width=\"57\" height=\"14\"/\u003e\n                \u003c/bpmndi:BPMNLabel\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNShape id=\"ServiceTask_0ptq5va_di\" bpmnElement=\"Task_0wjr1fj\"\u003e\n                \u003cdc:Bounds x=\"495\" y=\"80\" width=\"100\" height=\"80\"/\u003e\n            \u003c/bpmndi:BPMNShape\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0ixns30_di\" bpmnElement=\"SequenceFlow_0ixns30\"\u003e\n                \u003cdi:waypoint x=\"209\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"259\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0htq2f6_di\" bpmnElement=\"SequenceFlow_0htq2f6\"\u003e\n                \u003cdi:waypoint x=\"359\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"409\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_0npfu5a_di\" bpmnElement=\"SequenceFlow_0npfu5a\"\u003e\n                \u003cdi:waypoint x=\"445\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"495\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n            \u003cbpmndi:BPMNEdge id=\"SequenceFlow_1a1qwlk_di\" bpmnElement=\"SequenceFlow_1a1qwlk\"\u003e\n                \u003cdi:waypoint x=\"595\" y=\"120\"/\u003e\n                \u003cdi:waypoint x=\"645\" y=\"120\"/\u003e\n            \u003c/bpmndi:BPMNEdge\u003e\n        \u003c/bpmndi:BPMNPlane\u003e\n    \u003c/bpmndi:BPMNDiagram\u003e\n\u003c/bpmn:definitions\u003e",
   "name":"simple",
   "elements":[
      {
         "order":1,
         "task":{
            "label":"taskLabel",
            "device_description":{
               "characteristic_id":"example_hex",
               "function":{
                  "id":"fid"
               }
            },
            "bpmn_element_id":"Task_0wjr1fj",
            "input":"000",
            "device_options":null,
            "selection":{
               "selected_device":{
                  "id":"device_id_1"
               },
               "selected_service":{
                  "id":"service_id_1"
               }
            },
            "parameter":{
               "inputs":"\"ff0\""
            }
         }
      },
      {
         "order":2,
         "multi_task":{
            "label":"multiTaskLabel",
            "device_description":{
               "characteristic_id":"example_hex",
               "function":{
                  "id":"fid2"
               }
            },
            "bpmn_element_id":"Task_096xjeg",
            "input":"000",
            "device_options":null,
            "selections":[
				{
				   "selected_device":{
					  "id":"device_id_1"
				   },
				   "selected_service":{
					  "id":"service_id_1"
				   }
				},
				{
				   "selected_device":{
					  "id":"device_id_2"
				   },
				   "selected_service":{
					  "id":"service_id_2"
				   }
				}
			],
            "parameter":{
               "inputs":"\"fff\""
            }
         }
      },
      {
         "order":3,
         "msg_event":{
            "label":"eventName",
            "bpmn_element_id":"IntermediateThrowEvent_0905jg5",
            "device":{
               "id":"device_id_1"
            },
            "service":{
               "id":"service_id_2"
            },
            "path":"$.value+",
            "value":"42",
            "operation":"operation",
            "event_id":""
         }
      }
   ],
   "lanes":null
}`

	deployment := model.Deployment{}
	err := json.Unmarshal([]byte(deploymentStr), &deployment)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = UseDeploymentSelections(&deployment, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(deployment.Xml)

	//output:
	//

}
