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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	mock "github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"github.com/beevik/etree"
	"html"
	"io/ioutil"
	"runtime/debug"
	"testing"
)

func ExampleEventCastBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/event_cast.bpmn")
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
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"","name":"simple","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","input":"000","selectables":null,"selection":{"device":null,"service":null},"parameter":{"inputs":"\"ff0\""},"configurables":null}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","input":"000","selectables":null,"selections":null,"parameter":{"inputs":"\"fff\""},"configurables":null}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":null,"service":null,"path":"","value":"","operation":"","trigger_conversion":{"from":"","to":"example_hex"},"event_id":""}}],"lanes":null}
}

func ExampleEventCastBpmnDeploymentToXmlWithRefs() {
	file, err := ioutil.ReadFile("../../tests/resources/event_cast.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "test_id"
	}

	deployment := model.Deployment{}
	err = json.Unmarshal(file, &deployment)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	mock.Devices.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, true, mock.Devices)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(html.UnescapeString(deployment.Xml))

	//output:
	// <bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:camunda="http://camunda.org/schema/1.0/bpmn" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn">
	//     <bpmn:process id="simple" isExecutable="true" name="simple">
	//         <bpmn:startEvent id="StartEvent_1">
	//             <bpmn:outgoing>SequenceFlow_0ixns30</bpmn:outgoing>
	//         </bpmn:startEvent>
	//         <bpmn:sequenceFlow id="SequenceFlow_0ixns30" sourceRef="StartEvent_1" targetRef="Task_096xjeg"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0htq2f6" sourceRef="Task_096xjeg" targetRef="IntermediateThrowEvent_0905jg5"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0npfu5a" sourceRef="IntermediateThrowEvent_0905jg5" targetRef="Task_0wjr1fj"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_1a1qwlk" sourceRef="Task_0wjr1fj" targetRef="EndEvent_0yi4y22"/>
	//         <bpmn:endEvent id="EndEvent_0yi4y22">
	//             <bpmn:incoming>SequenceFlow_1a1qwlk</bpmn:incoming>
	//         </bpmn:endEvent>
	//         <bpmn:serviceTask id="Task_096xjeg" name="multiTaskLabel" camunda:type="external" camunda:topic="task">
	//             <bpmn:documentation>{"order": 2}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid2","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="inputs">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             <camunda:executionListener event="start"><camunda:script scriptFormat="groovy">execution.setVariable("collection", ["{\"device\":{\"id\":\"device_id_1\"},\"service\":{\"id\":\"service_id_1\",\"protocol_id\":\"pid\"},\"protocol\":{\"id\":\"pid\",\"name\":\"protocol1\",\"handler\":\"p\",\"protocol_segments\":null}}","{\"device\":{\"id\":\"device_id_2\"},\"service\":{\"id\":\"service_id_2\",\"protocol_id\":\"pid\"},\"protocol\":{\"id\":\"pid\",\"name\":\"protocol1\",\"handler\":\"p\",\"protocol_segments\":null}}"])</camunda:script></camunda:executionListener></bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0ixns30</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0htq2f6</bpmn:outgoing>
	//             <bpmn:multiInstanceLoopCharacteristics isSequential="true" camunda:collection="collection" camunda:elementVariable="overwrite"/>
	//         </bpmn:serviceTask>
	//         <bpmn:serviceTask id="Task_0wjr1fj" name="taskLabel" camunda:type="external" camunda:topic="task">
	//             <bpmn:documentation>{"order": 1}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","device_id":"device_id_1","service_id":"service_id_1","protocol_id":"pid","input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="inputs">"ff0"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             </bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0npfu5a</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_1a1qwlk</bpmn:outgoing>
	//         </bpmn:serviceTask>
	//         <bpmn:intermediateCatchEvent id="IntermediateThrowEvent_0905jg5" name="eventName">
	//             <bpmn:documentation>{"order": 3}</bpmn:documentation>
	//             <bpmn:incoming>SequenceFlow_0htq2f6</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0npfu5a</bpmn:outgoing>
	//             <bpmn:messageEventDefinition messageRef="e_test_id"/>
	//         </bpmn:intermediateCatchEvent>
	//     </bpmn:process>
	//     <bpmn:message id="e_test_id" name="test_id"/><bpmndi:BPMNDiagram id="BPMNDiagram_1">
	//         <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="simple">
	//             <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
	//                 <dc:Bounds x="173" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="EndEvent_0yi4y22_di" bpmnElement="EndEvent_0yi4y22">
	//                 <dc:Bounds x="645" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_1x4556d_di" bpmnElement="Task_096xjeg">
	//                 <dc:Bounds x="259" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="IntermediateCatchEvent_072g4ud_di" bpmnElement="IntermediateThrowEvent_0905jg5">
	//                 <dc:Bounds x="409" y="102" width="36" height="36"/>
	//                 <bpmndi:BPMNLabel>
	//                     <dc:Bounds x="399" y="145" width="57" height="14"/>
	//                 </bpmndi:BPMNLabel>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_0ptq5va_di" bpmnElement="Task_0wjr1fj">
	//                 <dc:Bounds x="495" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0ixns30_di" bpmnElement="SequenceFlow_0ixns30">
	//                 <di:waypoint x="209" y="120"/>
	//                 <di:waypoint x="259" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0htq2f6_di" bpmnElement="SequenceFlow_0htq2f6">
	//                 <di:waypoint x="359" y="120"/>
	//                 <di:waypoint x="409" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0npfu5a_di" bpmnElement="SequenceFlow_0npfu5a">
	//                 <di:waypoint x="445" y="120"/>
	//                 <di:waypoint x="495" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_1a1qwlk_di" bpmnElement="SequenceFlow_1a1qwlk">
	//                 <di:waypoint x="595" y="120"/>
	//                 <di:waypoint x="645" y="120"/>
	//             </bpmndi:BPMNEdge>
	//         </bpmndi:BPMNPlane>
	//     </bpmndi:BPMNDiagram>
	// </bpmn:definitions>
}

func TestEmptyTimeEvent(t *testing.T) {
	file, err := ioutil.ReadFile("../../tests/resources/empty_time_event.bpmn")
	if err != nil {
		t.Error(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		t.Error(err)
		return
	}
	err = UseDeploymentSelections(&result, false, mock.Devices)
	if err != nil {
		t.Error(err)
		return
	}
}

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
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"","name":"simple","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","input":"000","selectables":null,"selection":{"device":null,"service":null},"parameter":{"inputs":"\"ff0\""},"configurables":null}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","input":"000","selectables":null,"selections":null,"parameter":{"inputs":"\"fff\""},"configurables":null}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":null,"service":null,"path":"","value":"","operation":"","event_id":""}}],"lanes":null}
}

func ExampleLineWithOnlyTimer() {
	file, err := ioutil.ReadFile("../../tests/resources/lane_only_timer.bpmn")
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
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"","name":"Lane_Timer_FJ","elements":null,"lanes":[{"order":0,"lane":{"label":"Lane_Timer_fj","bpmn_element_id":"Lane_Timer_fj","device_descriptions":null,"selectables":null,"selection":null,"elements":[{"order":0,"time_event":{"bpmn_element_id":"IntermediateThrowEvent_1opksgz","kind":"timeDuration","time":"","label":"IntermediateThrowEvent_1opksgz"}}]}}]}

}

func ExampleIgnoredPredefinedEventBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/ignore_msg_event.bpmn")
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
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"","name":"ignore_msg_event","elements":[{"order":0,"msg_event":{"label":"IntermediateThrowEvent_1277ux9","bpmn_element_id":"IntermediateThrowEvent_1277ux9","device":null,"service":null,"path":"","value":"","operation":"","event_id":""}}],"lanes":null}
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
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"","name":"lanes","elements":null,"lanes":[{"order":0,"multi_lane":{"label":"multiTaskLane","bpmn_element_id":"Lane_12774cv","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"selectables":null,"selections":null,"elements":[{"order":3,"task":{"label":"multi_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_084s3g5","multi_task":true,"selected_service":null,"parameter":{},"configurables":null}},{"order":4,"task":{"label":"multi_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_098jmqp","multi_task":true,"selected_service":null,"parameter":{},"configurables":null}}]}},{"order":0,"lane":{"label":"MixedLane","bpmn_element_id":"Lane_0odlj5k","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"selectables":null,"selection":null,"elements":[{"order":5,"task":{"label":"mixed_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1npvonw","multi_task":true,"selected_service":null,"parameter":{},"configurables":null}},{"order":6,"task":{"label":"mixed_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1mnjsed","multi_task":false,"selected_service":null,"parameter":{},"configurables":null}}]}},{"order":1,"lane":{"label":"taskLane","bpmn_element_id":"Lane_0v679jg","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}}],"selectables":null,"selection":null,"elements":[{"order":1,"task":{"label":"lane_task_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_0nmb2on","multi_task":false,"selected_service":null,"parameter":{},"configurables":null}},{"order":2,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_1tchutl","device":null,"service":null,"path":"","value":"","operation":"","event_id":""}}]}}]}
}

func ExampleTimeAndReceiverBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/time_and_receive.bpmn")
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
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"","name":"timeAndReceive","elements":[{"order":0,"receive_task_event":{"label":"Task_1uyyxb0","bpmn_element_id":"Task_1uyyxb0","device":null,"service":null,"path":"","value":"","operation":"","event_id":""}},{"order":0,"time_event":{"bpmn_element_id":"IntermediateThrowEvent_10mhx3e","kind":"timeDuration","time":"PT1M","label":"eine Minute"}}],"lanes":null}
}

func ExampleSimpleBpmnDeploymentToXmlWithRefs() {
	file, err := ioutil.ReadFile("../../tests/resources/simple.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "test_id"
	}

	deployment := model.Deployment{}
	err = json.Unmarshal(file, &deployment)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	mock.Devices.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, true, mock.Devices)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(html.UnescapeString(deployment.Xml))

	//output:
	//<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:camunda="http://camunda.org/schema/1.0/bpmn" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn">
	//     <bpmn:process id="simple" isExecutable="true" name="simple">
	//         <bpmn:startEvent id="StartEvent_1">
	//             <bpmn:outgoing>SequenceFlow_0ixns30</bpmn:outgoing>
	//         </bpmn:startEvent>
	//         <bpmn:sequenceFlow id="SequenceFlow_0ixns30" sourceRef="StartEvent_1" targetRef="Task_096xjeg"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0htq2f6" sourceRef="Task_096xjeg" targetRef="IntermediateThrowEvent_0905jg5"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0npfu5a" sourceRef="IntermediateThrowEvent_0905jg5" targetRef="Task_0wjr1fj"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_1a1qwlk" sourceRef="Task_0wjr1fj" targetRef="EndEvent_0yi4y22"/>
	//         <bpmn:endEvent id="EndEvent_0yi4y22">
	//             <bpmn:incoming>SequenceFlow_1a1qwlk</bpmn:incoming>
	//         </bpmn:endEvent>
	//         <bpmn:serviceTask id="Task_096xjeg" name="multiTaskLabel" camunda:type="external" camunda:topic="task">
	//             <bpmn:documentation>{"order": 2}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid2","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="inputs">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             <camunda:executionListener event="start"><camunda:script scriptFormat="groovy">execution.setVariable("collection", ["{\"device\":{\"id\":\"device_id_1\"},\"service\":{\"id\":\"service_id_1\",\"protocol_id\":\"pid\"},\"protocol\":{\"id\":\"pid\",\"name\":\"protocol1\",\"handler\":\"p\",\"protocol_segments\":null}}","{\"device\":{\"id\":\"device_id_2\"},\"service\":{\"id\":\"service_id_2\",\"protocol_id\":\"pid\"},\"protocol\":{\"id\":\"pid\",\"name\":\"protocol1\",\"handler\":\"p\",\"protocol_segments\":null}}"])</camunda:script></camunda:executionListener></bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0ixns30</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0htq2f6</bpmn:outgoing>
	//             <bpmn:multiInstanceLoopCharacteristics isSequential="true" camunda:collection="collection" camunda:elementVariable="overwrite"/>
	//         </bpmn:serviceTask>
	//         <bpmn:serviceTask id="Task_0wjr1fj" name="taskLabel" camunda:type="external" camunda:topic="task">
	//             <bpmn:documentation>{"order": 1}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","device_id":"device_id_1","service_id":"service_id_1","protocol_id":"pid","input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="inputs">"ff0"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             </bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0npfu5a</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_1a1qwlk</bpmn:outgoing>
	//         </bpmn:serviceTask>
	//         <bpmn:intermediateCatchEvent id="IntermediateThrowEvent_0905jg5" name="eventName">
	//             <bpmn:documentation>{"order": 3}</bpmn:documentation>
	//             <bpmn:incoming>SequenceFlow_0htq2f6</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0npfu5a</bpmn:outgoing>
	//             <bpmn:messageEventDefinition messageRef="e_test_id"/>
	//         </bpmn:intermediateCatchEvent>
	//     </bpmn:process>
	//     <bpmn:message id="e_test_id" name="test_id"/><bpmndi:BPMNDiagram id="BPMNDiagram_1">
	//         <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="simple">
	//             <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
	//                 <dc:Bounds x="173" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="EndEvent_0yi4y22_di" bpmnElement="EndEvent_0yi4y22">
	//                 <dc:Bounds x="645" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_1x4556d_di" bpmnElement="Task_096xjeg">
	//                 <dc:Bounds x="259" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="IntermediateCatchEvent_072g4ud_di" bpmnElement="IntermediateThrowEvent_0905jg5">
	//                 <dc:Bounds x="409" y="102" width="36" height="36"/>
	//                 <bpmndi:BPMNLabel>
	//                     <dc:Bounds x="399" y="145" width="57" height="14"/>
	//                 </bpmndi:BPMNLabel>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_0ptq5va_di" bpmnElement="Task_0wjr1fj">
	//                 <dc:Bounds x="495" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0ixns30_di" bpmnElement="SequenceFlow_0ixns30">
	//                 <di:waypoint x="209" y="120"/>
	//                 <di:waypoint x="259" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0htq2f6_di" bpmnElement="SequenceFlow_0htq2f6">
	//                 <di:waypoint x="359" y="120"/>
	//                 <di:waypoint x="409" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0npfu5a_di" bpmnElement="SequenceFlow_0npfu5a">
	//                 <di:waypoint x="445" y="120"/>
	//                 <di:waypoint x="495" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_1a1qwlk_di" bpmnElement="SequenceFlow_1a1qwlk">
	//                 <di:waypoint x="595" y="120"/>
	//                 <di:waypoint x="645" y="120"/>
	//             </bpmndi:BPMNEdge>
	//         </bpmndi:BPMNPlane>
	//     </bpmndi:BPMNDiagram>
	//</bpmn:definitions>

}

func ExampleSimpleBpmnDeploymentToXml() {
	file, err := ioutil.ReadFile("../../tests/resources/simple.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "test_id"
	}

	deployment := model.Deployment{}
	err = json.Unmarshal(file, &deployment)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	mock.Devices.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, false, mock.Devices)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(html.UnescapeString(deployment.Xml))

	//output:
	//<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:camunda="http://camunda.org/schema/1.0/bpmn" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn">
	//     <bpmn:process id="simple" isExecutable="true" name="simple">
	//         <bpmn:startEvent id="StartEvent_1">
	//             <bpmn:outgoing>SequenceFlow_0ixns30</bpmn:outgoing>
	//         </bpmn:startEvent>
	//         <bpmn:sequenceFlow id="SequenceFlow_0ixns30" sourceRef="StartEvent_1" targetRef="Task_096xjeg"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0htq2f6" sourceRef="Task_096xjeg" targetRef="IntermediateThrowEvent_0905jg5"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0npfu5a" sourceRef="IntermediateThrowEvent_0905jg5" targetRef="Task_0wjr1fj"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_1a1qwlk" sourceRef="Task_0wjr1fj" targetRef="EndEvent_0yi4y22"/>
	//         <bpmn:endEvent id="EndEvent_0yi4y22">
	//             <bpmn:incoming>SequenceFlow_1a1qwlk</bpmn:incoming>
	//         </bpmn:endEvent>
	//         <bpmn:serviceTask id="Task_096xjeg" name="multiTaskLabel" camunda:type="external" camunda:topic="task">
	//             <bpmn:documentation>{"order": 2}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid2","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="inputs">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             <camunda:executionListener event="start"><camunda:script scriptFormat="groovy">execution.setVariable("collection", ["{\"device_id\":\"device_id_1\",\"service_id\":\"service_id_1\",\"protocol_id\":\"pid\"}","{\"device_id\":\"device_id_2\",\"service_id\":\"service_id_2\",\"protocol_id\":\"pid\"}"])</camunda:script></camunda:executionListener></bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0ixns30</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0htq2f6</bpmn:outgoing>
	//             <bpmn:multiInstanceLoopCharacteristics isSequential="true" camunda:collection="collection" camunda:elementVariable="overwrite"/>
	//         </bpmn:serviceTask>
	//         <bpmn:serviceTask id="Task_0wjr1fj" name="taskLabel" camunda:type="external" camunda:topic="task">
	//             <bpmn:documentation>{"order": 1}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","device":{"id":"device_id_1"},"service":{"id":"service_id_1","protocol_id":"pid"},"protocol":{"id":"pid","name":"protocol1","handler":"p","protocol_segments":null},"input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="inputs">"ff0"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             </bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0npfu5a</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_1a1qwlk</bpmn:outgoing>
	//         </bpmn:serviceTask>
	//         <bpmn:intermediateCatchEvent id="IntermediateThrowEvent_0905jg5" name="eventName">
	//             <bpmn:documentation>{"order": 3}</bpmn:documentation>
	//             <bpmn:incoming>SequenceFlow_0htq2f6</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0npfu5a</bpmn:outgoing>
	//             <bpmn:messageEventDefinition messageRef="e_test_id"/>
	//         </bpmn:intermediateCatchEvent>
	//     </bpmn:process>
	//     <bpmn:message id="e_test_id" name="test_id"/><bpmndi:BPMNDiagram id="BPMNDiagram_1">
	//         <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="simple">
	//             <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
	//                 <dc:Bounds x="173" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="EndEvent_0yi4y22_di" bpmnElement="EndEvent_0yi4y22">
	//                 <dc:Bounds x="645" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_1x4556d_di" bpmnElement="Task_096xjeg">
	//                 <dc:Bounds x="259" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="IntermediateCatchEvent_072g4ud_di" bpmnElement="IntermediateThrowEvent_0905jg5">
	//                 <dc:Bounds x="409" y="102" width="36" height="36"/>
	//                 <bpmndi:BPMNLabel>
	//                     <dc:Bounds x="399" y="145" width="57" height="14"/>
	//                 </bpmndi:BPMNLabel>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_0ptq5va_di" bpmnElement="Task_0wjr1fj">
	//                 <dc:Bounds x="495" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0ixns30_di" bpmnElement="SequenceFlow_0ixns30">
	//                 <di:waypoint x="209" y="120"/>
	//                 <di:waypoint x="259" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0htq2f6_di" bpmnElement="SequenceFlow_0htq2f6">
	//                 <di:waypoint x="359" y="120"/>
	//                 <di:waypoint x="409" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0npfu5a_di" bpmnElement="SequenceFlow_0npfu5a">
	//                 <di:waypoint x="445" y="120"/>
	//                 <di:waypoint x="495" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_1a1qwlk_di" bpmnElement="SequenceFlow_1a1qwlk">
	//                 <di:waypoint x="595" y="120"/>
	//                 <di:waypoint x="645" y="120"/>
	//             </bpmndi:BPMNEdge>
	//         </bpmndi:BPMNPlane>
	//     </bpmndi:BPMNDiagram>
	// </bpmn:definitions>

}

func ExampleDeploymentWithConfigurablesToXml() {
	file, err := ioutil.ReadFile("../../tests/resources/simple.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "test_id"
	}

	deployment := model.Deployment{}
	err = json.Unmarshal(file, &deployment)
	if err != nil {
		fmt.Println(err)
		return
	}

	for index, _ := range deployment.Elements {
		if deployment.Elements[index].Task != nil {
			deployment.Elements[index].Task.Configurables = []model.Configurable{{
				CharacteristicId: "foo",
				Values: []model.ConfigurableCharacteristicValue{
					{
						Label:     "bar",
						Path:      "batz",
						Value:     42,
						ValueType: devicemodel.Integer,
					},
				},
			}}
		}
	}

	err = deployment.Validate(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	mock.Devices.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, false, mock.Devices)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(html.UnescapeString(deployment.Xml))

	//output:
	// <bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:camunda="http://camunda.org/schema/1.0/bpmn" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn">
	//     <bpmn:process id="simple" isExecutable="true" name="simple">
	//         <bpmn:startEvent id="StartEvent_1">
	//             <bpmn:outgoing>SequenceFlow_0ixns30</bpmn:outgoing>
	//         </bpmn:startEvent>
	//         <bpmn:sequenceFlow id="SequenceFlow_0ixns30" sourceRef="StartEvent_1" targetRef="Task_096xjeg"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0htq2f6" sourceRef="Task_096xjeg" targetRef="IntermediateThrowEvent_0905jg5"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0npfu5a" sourceRef="IntermediateThrowEvent_0905jg5" targetRef="Task_0wjr1fj"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_1a1qwlk" sourceRef="Task_0wjr1fj" targetRef="EndEvent_0yi4y22"/>
	//         <bpmn:endEvent id="EndEvent_0yi4y22">
	//             <bpmn:incoming>SequenceFlow_1a1qwlk</bpmn:incoming>
	//         </bpmn:endEvent>
	//         <bpmn:serviceTask id="Task_096xjeg" name="multiTaskLabel" camunda:type="external" camunda:topic="task">
	//             <bpmn:documentation>{"order": 2}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid2","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="inputs">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             <camunda:executionListener event="start"><camunda:script scriptFormat="groovy">execution.setVariable("collection", ["{\"device_id\":\"device_id_1\",\"service_id\":\"service_id_1\",\"protocol_id\":\"pid\"}","{\"device_id\":\"device_id_2\",\"service_id\":\"service_id_2\",\"protocol_id\":\"pid\"}"])</camunda:script></camunda:executionListener></bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0ixns30</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0htq2f6</bpmn:outgoing>
	//             <bpmn:multiInstanceLoopCharacteristics isSequential="true" camunda:collection="collection" camunda:elementVariable="overwrite"/>
	//         </bpmn:serviceTask>
	//         <bpmn:serviceTask id="Task_0wjr1fj" name="taskLabel" camunda:type="external" camunda:topic="task">
	//             <bpmn:documentation>{"order": 1}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","device":{"id":"device_id_1"},"service":{"id":"service_id_1","protocol_id":"pid"},"protocol":{"id":"pid","name":"protocol1","handler":"p","protocol_segments":null},"configurables":[{"characteristic_id":"foo","values":[{"label":"bar","path":"batz","value":42,"value_type":"https://schema.org/Integer"}]}],"input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="inputs">"ff0"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             </bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0npfu5a</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_1a1qwlk</bpmn:outgoing>
	//         </bpmn:serviceTask>
	//         <bpmn:intermediateCatchEvent id="IntermediateThrowEvent_0905jg5" name="eventName">
	//             <bpmn:documentation>{"order": 3}</bpmn:documentation>
	//             <bpmn:incoming>SequenceFlow_0htq2f6</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0npfu5a</bpmn:outgoing>
	//             <bpmn:messageEventDefinition messageRef="e_test_id"/>
	//         </bpmn:intermediateCatchEvent>
	//     </bpmn:process>
	//     <bpmn:message id="e_test_id" name="test_id"/><bpmndi:BPMNDiagram id="BPMNDiagram_1">
	//         <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="simple">
	//             <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
	//                 <dc:Bounds x="173" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="EndEvent_0yi4y22_di" bpmnElement="EndEvent_0yi4y22">
	//                 <dc:Bounds x="645" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_1x4556d_di" bpmnElement="Task_096xjeg">
	//                 <dc:Bounds x="259" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="IntermediateCatchEvent_072g4ud_di" bpmnElement="IntermediateThrowEvent_0905jg5">
	//                 <dc:Bounds x="409" y="102" width="36" height="36"/>
	//                 <bpmndi:BPMNLabel>
	//                     <dc:Bounds x="399" y="145" width="57" height="14"/>
	//                 </bpmndi:BPMNLabel>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_0ptq5va_di" bpmnElement="Task_0wjr1fj">
	//                 <dc:Bounds x="495" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0ixns30_di" bpmnElement="SequenceFlow_0ixns30">
	//                 <di:waypoint x="209" y="120"/>
	//                 <di:waypoint x="259" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0htq2f6_di" bpmnElement="SequenceFlow_0htq2f6">
	//                 <di:waypoint x="359" y="120"/>
	//                 <di:waypoint x="409" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0npfu5a_di" bpmnElement="SequenceFlow_0npfu5a">
	//                 <di:waypoint x="445" y="120"/>
	//                 <di:waypoint x="495" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_1a1qwlk_di" bpmnElement="SequenceFlow_1a1qwlk">
	//                 <di:waypoint x="595" y="120"/>
	//                 <di:waypoint x="645" y="120"/>
	//             </bpmndi:BPMNEdge>
	//         </bpmndi:BPMNPlane>
	//     </bpmndi:BPMNDiagram>
	// </bpmn:definitions>

}

func ExampleTimeAndReceiveBpmnDeploymentToXml() {
	file, err := ioutil.ReadFile("../../tests/resources/time_and_receive.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "test_id"
	}

	deployment := model.Deployment{}
	err = json.Unmarshal(file, &deployment)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	mock.Devices.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, false, mock.Devices)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(html.UnescapeString(deployment.Xml))

	//output:
	//<bpmn:definitions xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn">
	//     <bpmn:process id="timeAndReceive" isExecutable="true" name="timeAndReceive">
	//         <bpmn:startEvent id="StartEvent_1">
	//             <bpmn:outgoing>SequenceFlow_1nh3k2a</bpmn:outgoing>
	//         </bpmn:startEvent>
	//         <bpmn:sequenceFlow id="SequenceFlow_1nh3k2a" sourceRef="StartEvent_1" targetRef="IntermediateThrowEvent_10mhx3e"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0d8iosp" sourceRef="IntermediateThrowEvent_10mhx3e" targetRef="Task_1uyyxb0"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_006q38h" sourceRef="Task_1uyyxb0" targetRef="EndEvent_0oatj6u"/>
	//         <bpmn:intermediateCatchEvent id="IntermediateThrowEvent_10mhx3e" name="eine Minute">
	//             <bpmn:incoming>SequenceFlow_1nh3k2a</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0d8iosp</bpmn:outgoing>
	//             <bpmn:timerEventDefinition>
	//                 <bpmn:timeDuration xsi:type="bpmn:tFormalExpression">PT2M</bpmn:timeDuration>
	//             </bpmn:timerEventDefinition>
	//         </bpmn:intermediateCatchEvent>
	//         <bpmn:receiveTask id="Task_1uyyxb0" messageRef="e_test_id">
	//             <bpmn:incoming>SequenceFlow_0d8iosp</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_006q38h</bpmn:outgoing>
	//         </bpmn:receiveTask>
	//         <bpmn:endEvent id="EndEvent_0oatj6u">
	//             <bpmn:incoming>SequenceFlow_006q38h</bpmn:incoming>
	//         </bpmn:endEvent>
	//     </bpmn:process>
	//     <bpmn:message id="e_test_id" name="test_id"/><bpmndi:BPMNDiagram id="BPMNDiagram_1">
	//         <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="timeAndReceive">
	//             <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
	//                 <dc:Bounds x="173" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="IntermediateCatchEvent_0ztbgip_di" bpmnElement="IntermediateThrowEvent_10mhx3e">
	//                 <dc:Bounds x="259" y="102" width="36" height="36"/>
	//                 <bpmndi:BPMNLabel>
	//                     <dc:Bounds x="249" y="145" width="57" height="14"/>
	//                 </bpmndi:BPMNLabel>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ReceiveTask_08kr9nv_di" bpmnElement="Task_1uyyxb0">
	//                 <dc:Bounds x="345" y="80" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="EndEvent_0oatj6u_di" bpmnElement="EndEvent_0oatj6u">
	//                 <dc:Bounds x="495" y="102" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNEdge id="SequenceFlow_1nh3k2a_di" bpmnElement="SequenceFlow_1nh3k2a">
	//                 <di:waypoint x="209" y="120"/>
	//                 <di:waypoint x="259" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0d8iosp_di" bpmnElement="SequenceFlow_0d8iosp">
	//                 <di:waypoint x="295" y="120"/>
	//                 <di:waypoint x="345" y="120"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_006q38h_di" bpmnElement="SequenceFlow_006q38h">
	//                 <di:waypoint x="445" y="120"/>
	//                 <di:waypoint x="495" y="120"/>
	//             </bpmndi:BPMNEdge>
	//         </bpmndi:BPMNPlane>
	//     </bpmndi:BPMNDiagram>
	// </bpmn:definitions>
}

func ExampleLanesBpmnDeploymentToXml() {
	file, err := ioutil.ReadFile("../../tests/resources/lanes.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "test_id"
	}

	deployment := model.Deployment{}
	err = json.Unmarshal(file, &deployment)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	mock.Devices.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, false, mock.Devices)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deployment.Validate(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(html.UnescapeString(deployment.Xml))

	//output:
	//<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" xmlns:camunda="http://camunda.org/schema/1.0/bpmn" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn">
	//     <bpmn:collaboration id="Collaboration_1t682rc">
	//         <bpmn:participant id="Participant_1a4rg7s" processRef="lanes"/>
	//     </bpmn:collaboration>
	//     <bpmn:process id="lanes" isExecutable="true" name="lanes">
	//         <bpmn:laneSet id="LaneSet_1ekoknz">
	//             <bpmn:lane id="Lane_0v679jg" name="taskLane">
	//                 <bpmn:documentation>{"order":1}</bpmn:documentation>
	//                 <bpmn:flowNodeRef>IntermediateThrowEvent_1tchutl</bpmn:flowNodeRef>
	//                 <bpmn:flowNodeRef>StartEvent_1</bpmn:flowNodeRef>
	//                 <bpmn:flowNodeRef>Task_0nmb2on</bpmn:flowNodeRef>
	//             </bpmn:lane>
	//             <bpmn:lane id="Lane_12774cv" name="multiTaskLane">
	//                 <bpmn:flowNodeRef>Task_084s3g5</bpmn:flowNodeRef>
	//                 <bpmn:flowNodeRef>Task_098jmqp</bpmn:flowNodeRef>
	//             </bpmn:lane>
	//             <bpmn:lane id="Lane_0odlj5k" name="MixedLane">
	//                 <bpmn:flowNodeRef>EndEvent_0yfaeyo</bpmn:flowNodeRef>
	//                 <bpmn:flowNodeRef>Task_1npvonw</bpmn:flowNodeRef>
	//                 <bpmn:flowNodeRef>Task_1mnjsed</bpmn:flowNodeRef>
	//             </bpmn:lane>
	//         </bpmn:laneSet>
	//         <bpmn:intermediateCatchEvent id="IntermediateThrowEvent_1tchutl" name="eventName">
	//             <bpmn:documentation>{"order": 2}</bpmn:documentation>
	//             <bpmn:incoming>SequenceFlow_02ma6x0</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0vo4tpu</bpmn:outgoing>
	//             <bpmn:messageEventDefinition messageRef="e_test_id"/>
	//         </bpmn:intermediateCatchEvent>
	//         <bpmn:endEvent id="EndEvent_0yfaeyo">
	//             <bpmn:incoming>SequenceFlow_0v5i9ks</bpmn:incoming>
	//         </bpmn:endEvent>
	//         <bpmn:startEvent id="StartEvent_1">
	//             <bpmn:outgoing>SequenceFlow_1t4knqk</bpmn:outgoing>
	//         </bpmn:startEvent>
	//         <bpmn:sequenceFlow id="SequenceFlow_1t4knqk" sourceRef="StartEvent_1" targetRef="Task_0nmb2on"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_02ma6x0" sourceRef="Task_0nmb2on" targetRef="IntermediateThrowEvent_1tchutl"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0vo4tpu" sourceRef="IntermediateThrowEvent_1tchutl" targetRef="Task_084s3g5"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0q1hn6b" sourceRef="Task_084s3g5" targetRef="Task_098jmqp"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0g2ladw" sourceRef="Task_098jmqp" targetRef="Task_1npvonw"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_1pelemy" sourceRef="Task_1npvonw" targetRef="Task_1mnjsed"/>
	//         <bpmn:sequenceFlow id="SequenceFlow_0v5i9ks" sourceRef="Task_1mnjsed" targetRef="EndEvent_0yfaeyo"/>
	//         <bpmn:serviceTask id="Task_0nmb2on" name="lane_task_1" camunda:type="external" camunda:topic="test">
	//             <bpmn:documentation>{"order": 1}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","device":{"id":"device_1"},"service":{"id":"service_1","protocol_id":"pid"},"protocol":{"id":"pid","name":"protocol1","handler":"p","protocol_segments":null},"input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="input">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             </bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_1t4knqk</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_02ma6x0</bpmn:outgoing>
	//         </bpmn:serviceTask>
	//         <bpmn:serviceTask id="Task_084s3g5" name="multi_lane_1" camunda:type="external" camunda:topic="test">
	//             <bpmn:documentation>{"order": 3}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid1","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="input">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             <camunda:executionListener event="start"><camunda:script scriptFormat="groovy">execution.setVariable("collection", ["{\"device_id\":\"device_1\",\"service_id\":\"service_1\",\"protocol_id\":\"pid\"}","{\"device_id\":\"device_2\",\"service_id\":\"service_1\",\"protocol_id\":\"pid\"}"])</camunda:script></camunda:executionListener></bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0vo4tpu</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0q1hn6b</bpmn:outgoing>
	//             <bpmn:multiInstanceLoopCharacteristics isSequential="true" camunda:collection="collection" camunda:elementVariable="overwrite"/>
	//         </bpmn:serviceTask>
	//         <bpmn:serviceTask id="Task_098jmqp" name="multi_lane_2" camunda:type="external" camunda:topic="test">
	//             <bpmn:documentation>{"order": 4}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid2","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="input">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             <camunda:executionListener event="start"><camunda:script scriptFormat="groovy">execution.setVariable("collection", ["{\"device_id\":\"device_1\",\"service_id\":\"service_2\",\"protocol_id\":\"pid\"}","{\"device_id\":\"device_2\",\"service_id\":\"service_2\",\"protocol_id\":\"pid\"}"])</camunda:script></camunda:executionListener></bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0q1hn6b</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0g2ladw</bpmn:outgoing>
	//             <bpmn:multiInstanceLoopCharacteristics isSequential="true" camunda:collection="collection" camunda:elementVariable="overwrite"/>
	//         </bpmn:serviceTask>
	//         <bpmn:serviceTask id="Task_1npvonw" name="mixed_lane_1" camunda:type="external" camunda:topic="test">
	//             <bpmn:documentation>{"order": 5}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid1","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","device":{"id":"device_2"},"service":{"id":"service_1","protocol_id":"pid"},"protocol":{"id":"pid","name":"protocol1","handler":"p","protocol_segments":null},"input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="input">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             </bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_0g2ladw</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_1pelemy</bpmn:outgoing>
	//             <bpmn:multiInstanceLoopCharacteristics isSequential="true"/>
	//         </bpmn:serviceTask>
	//         <bpmn:serviceTask id="Task_1mnjsed" name="mixed_lane_2" camunda:type="external" camunda:topic="test">
	//             <bpmn:documentation>{"order": 6}</bpmn:documentation>
	//             <bpmn:extensionElements>
	//                 <camunda:inputOutput>
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid2","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","device":{"id":"device_2"},"service":{"id":"service_2","protocol_id":"pid"},"protocol":{"id":"pid","name":"protocol1","handler":"p","protocol_segments":null},"input":"000"}</camunda:inputParameter>
	//                     <camunda:inputParameter name="input">"fff"</camunda:inputParameter>
	//                 </camunda:inputOutput>
	//             </bpmn:extensionElements>
	//             <bpmn:incoming>SequenceFlow_1pelemy</bpmn:incoming>
	//             <bpmn:outgoing>SequenceFlow_0v5i9ks</bpmn:outgoing>
	//         </bpmn:serviceTask>
	//     </bpmn:process>
	//     <bpmn:message id="e_test_id" name="test_id"/><bpmndi:BPMNDiagram id="BPMNDiagram_1">
	//         <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="Collaboration_1t682rc">
	//             <bpmndi:BPMNShape id="Participant_1a4rg7s_di" bpmnElement="Participant_1a4rg7s" isHorizontal="true">
	//                 <dc:Bounds x="119" y="125" width="964" height="414"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
	//                 <dc:Bounds x="199" y="196" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="Lane_0v679jg_di" bpmnElement="Lane_0v679jg" isHorizontal="true">
	//                 <dc:Bounds x="149" y="125" width="934" height="169"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="Lane_12774cv_di" bpmnElement="Lane_12774cv" isHorizontal="true">
	//                 <dc:Bounds x="149" y="294" width="934" height="125"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="Lane_0odlj5k_di" bpmnElement="Lane_0odlj5k" isHorizontal="true">
	//                 <dc:Bounds x="149" y="419" width="934" height="120"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="IntermediateCatchEvent_0d8lnnn_di" bpmnElement="IntermediateThrowEvent_1tchutl">
	//                 <dc:Bounds x="473" y="196" width="36" height="36"/>
	//                 <bpmndi:BPMNLabel>
	//                     <dc:Bounds x="463" y="166" width="57" height="14"/>
	//                 </bpmndi:BPMNLabel>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="EndEvent_0yfaeyo_di" bpmnElement="EndEvent_0yfaeyo">
	//                 <dc:Bounds x="947" y="460" width="36" height="36"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_11juy3k_di" bpmnElement="Task_0nmb2on">
	//                 <dc:Bounds x="304" y="174" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_0rbj5ci_di" bpmnElement="Task_084s3g5">
	//                 <dc:Bounds x="441" y="316" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_0k0tvsq_di" bpmnElement="Task_098jmqp">
	//                 <dc:Bounds x="625" y="316" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_0rvp7od_di" bpmnElement="Task_1npvonw">
	//                 <dc:Bounds x="625" y="438" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNShape id="ServiceTask_0cubfu5_di" bpmnElement="Task_1mnjsed">
	//                 <dc:Bounds x="807" y="438" width="100" height="80"/>
	//             </bpmndi:BPMNShape>
	//             <bpmndi:BPMNEdge id="SequenceFlow_1t4knqk_di" bpmnElement="SequenceFlow_1t4knqk">
	//                 <di:waypoint x="235" y="214"/>
	//                 <di:waypoint x="304" y="214"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_02ma6x0_di" bpmnElement="SequenceFlow_02ma6x0">
	//                 <di:waypoint x="404" y="214"/>
	//                 <di:waypoint x="473" y="214"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0vo4tpu_di" bpmnElement="SequenceFlow_0vo4tpu">
	//                 <di:waypoint x="491" y="232"/>
	//                 <di:waypoint x="491" y="316"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0q1hn6b_di" bpmnElement="SequenceFlow_0q1hn6b">
	//                 <di:waypoint x="541" y="356"/>
	//                 <di:waypoint x="625" y="356"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0g2ladw_di" bpmnElement="SequenceFlow_0g2ladw">
	//                 <di:waypoint x="675" y="396"/>
	//                 <di:waypoint x="675" y="438"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_1pelemy_di" bpmnElement="SequenceFlow_1pelemy">
	//                 <di:waypoint x="725" y="478"/>
	//                 <di:waypoint x="807" y="478"/>
	//             </bpmndi:BPMNEdge>
	//             <bpmndi:BPMNEdge id="SequenceFlow_0v5i9ks_di" bpmnElement="SequenceFlow_0v5i9ks">
	//                 <di:waypoint x="907" y="478"/>
	//                 <di:waypoint x="947" y="478"/>
	//             </bpmndi:BPMNEdge>
	//         </bpmndi:BPMNPlane>
	//     </bpmndi:BPMNDiagram>
	// </bpmn:definitions>
}

func TestNotificationsBpmnDeployment(t *testing.T) {
	file, err := ioutil.ReadFile("../../tests/resources/notifications.bpmn")
	if err != nil {
		t.Error(err)
	}

	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "test_id"
	}

	deployment, err := PrepareDeployment(string(file))
	if err != nil {
		t.Error(err)
	}

	mock.Devices.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, false, mock.Devices)
	if err != nil {
		t.Error(err)
	}

	deployment.Id = "deploymentId"

	doc := etree.NewDocument()
	err = doc.ReadFromString(deployment.Xml)
	if err != nil {
		t.Error(err)
	}

	for _, input := range doc.FindElements("//camunda:inputParameter[@name='deploymentIdentifier']") {
		if input.Text() != "notification" {
			continue
		}
		parent := input.Parent()

		// Check url
		urlParameter := parent.FindElement("camunda:inputParameter[@name='url']")
		if urlParameter.Text() != "url" {
			t.Error("url not set correctly")
		}

		// Check method
		methodParameter := parent.FindElement("camunda:inputParameter[@name='method']")
		if methodParameter.Text() != "PUT" {
			t.Error("method not set correctly")
		}

		// Check payload
		payloadParameter := parent.FindElement("camunda:inputParameter[@name='payload']")
		notificationPayload := model.NotificationPayload{}
		err = json.Unmarshal([]byte(payloadParameter.Text()), &notificationPayload)
		if err != nil {
			t.Error(err)
		}
		if notificationPayload.UserId != "uid" {
			t.Error("userId not set correctly")
		}
		if notificationPayload.IsRead != false {
			t.Error("read status not set correctly")
		}

		// Check header
		keyElement := parent.FindElement("camunda:inputParameter[@name='headers']/camunda:map/camunda:entry[@key='Content-Type']")
		if keyElement.Text() != "application/json" {
			t.Error("Content-Type header not set correctly")
		}
	}
}

func ExampleEmptyLaneBpmnDeploymentToXml() {
	file, err := ioutil.ReadFile("../../tests/resources/lane_only_timer.json")
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	temp := config.NewId
	defer func() {
		config.NewId = temp
	}()
	config.NewId = func() string {
		return "test_id"
	}

	deployment := model.Deployment{}
	err = json.Unmarshal(file, &deployment)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	err = deployment.Validate(false)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	mock.Devices.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, false, mock.Devices)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	err = deployment.Validate(true)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	fmt.Println(html.UnescapeString(deployment.Xml))

	//output:
	//<?xml version="1.0" encoding="UTF-8"?>
	//<bpmn:definitions xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn"><bpmn:collaboration id="Lane_Timer_FJ"><bpmn:participant id="Participant_1nxsivp" processRef="Lane_Timer_fj"/></bpmn:collaboration><bpmn:process id="Lane_Timer_fj" isExecutable="true" name="Lane_Timer_FJ"><bpmn:startEvent id="StartEvent_1"><bpmn:outgoing>SequenceFlow_1azj5gx</bpmn:outgoing></bpmn:startEvent><bpmn:sequenceFlow id="SequenceFlow_1azj5gx" sourceRef="StartEvent_1" targetRef="IntermediateThrowEvent_1opksgz"/><bpmn:intermediateCatchEvent id="IntermediateThrowEvent_1opksgz"><bpmn:incoming>SequenceFlow_1azj5gx</bpmn:incoming><bpmn:outgoing>SequenceFlow_0fm9shw</bpmn:outgoing><bpmn:timerEventDefinition><./bpmn:timeDuration>PT1H</./bpmn:timeDuration></bpmn:timerEventDefinition></bpmn:intermediateCatchEvent><bpmn:endEvent id="EndEvent_1h8yy4g"><bpmn:incoming>SequenceFlow_0fm9shw</bpmn:incoming></bpmn:endEvent><bpmn:sequenceFlow id="SequenceFlow_0fm9shw" sourceRef="IntermediateThrowEvent_1opksgz" targetRef="EndEvent_1h8yy4g"/></bpmn:process><bpmndi:BPMNDiagram id="BPMNDiagram_1"><bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="Lane_Timer_FJ"><bpmndi:BPMNShape id="Participant_1nxsivp_di" bpmnElement="Participant_1nxsivp" isHorizontal="true"><dc:Bounds x="123" y="80" width="337" height="90"/></bpmndi:BPMNShape><bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1"><dc:Bounds x="173" y="102" width="36" height="36"/></bpmndi:BPMNShape><bpmndi:BPMNEdge id="SequenceFlow_1azj5gx_di" bpmnElement="SequenceFlow_1azj5gx"><di:waypoint x="209" y="120"/><di:waypoint x="262" y="120"/></bpmndi:BPMNEdge><bpmndi:BPMNShape id="IntermediateCatchEvent_1j97xwq_di" bpmnElement="IntermediateThrowEvent_1opksgz"><dc:Bounds x="262" y="102" width="36" height="36"/></bpmndi:BPMNShape><bpmndi:BPMNShape id="EndEvent_1h8yy4g_di" bpmnElement="EndEvent_1h8yy4g"><dc:Bounds x="352" y="102" width="36" height="36"/></bpmndi:BPMNShape><bpmndi:BPMNEdge id="SequenceFlow_0fm9shw_di" bpmnElement="SequenceFlow_0fm9shw"><di:waypoint x="298" y="120"/><di:waypoint x="352" y="120"/></bpmndi:BPMNEdge></bpmndi:BPMNPlane></bpmndi:BPMNDiagram></bpmn:definitions>

}

func ExampleSingleLaneNameBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/single_lane_name.bpmn")
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
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"","name":"Collaboration_0wfm1tv","elements":null,"lanes":[{"order":0,"lane":{"label":"TestLaneName","bpmn_element_id":"Process_1","device_descriptions":[{"characteristic_id":"urn:infai:ses:characteristic:72b624b5-6edc-4ec4-9ad9-fa00b39915c0","function":{"id":"urn:infai:ses:controlling-function:7adc7f29-5c37-4bfc-8508-6130a143ac66","name":"brightnessFunction","concept_id":"urn:infai:ses:concept:dbe4ad57-aa1d-4d24-9bee-a44a1c670d7f","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}}],"selectables":null,"selection":null,"elements":[{"order":0,"task":{"label":"Lamp brightnessFunction","device_description":{"characteristic_id":"urn:infai:ses:characteristic:72b624b5-6edc-4ec4-9ad9-fa00b39915c0","function":{"id":"urn:infai:ses:controlling-function:7adc7f29-5c37-4bfc-8508-6130a143ac66","name":"brightnessFunction","concept_id":"urn:infai:ses:concept:dbe4ad57-aa1d-4d24-9bee-a44a1c670d7f","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},"input":0,"bpmn_element_id":"Task_1d0tawd","multi_task":false,"selected_service":null,"parameter":{"inputs":"0"},"configurables":null}}]}}]}

}

func ExampleMultiplePoolsBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/RW_9.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(result.Lanes))
	result.XmlRaw = ""
	temp, err := json.Marshal(result)
	fmt.Println(err, string(temp))

	//output:
	//2
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"","name":"Collaboration_0puddqt","elements":null,"lanes":[{"order":0,"lane":{"label":"Smart-Plug-Steuerung","bpmn_element_id":"RW_9","device_descriptions":[{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:79e7914b-f303-4a7d-90af-dee70db05fd9","name":"setOnStateFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:79de1bd9-b933-412d-b98e-4cfe19aa3250","name":"SmartPlug","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}}],"selectables":null,"selection":null,"elements":[{"order":0,"task":{"label":"Beamer anschalten","device_description":{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:79e7914b-f303-4a7d-90af-dee70db05fd9","name":"setOnStateFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:79de1bd9-b933-412d-b98e-4cfe19aa3250","name":"SmartPlug","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},"input":null,"bpmn_element_id":"Task_1f9iluy","multi_task":false,"selected_service":null,"parameter":{},"configurables":null}}]}},{"order":0,"lane":{"label":"Lampensteuerung","bpmn_element_id":"Process_1qgdhm1","device_descriptions":[{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:79e7914b-f303-4a7d-90af-dee70db05fd9","name":"setOnStateFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:2f35150b-9df7-4cad-95bc-165fa00219fd","name":"setOffStateFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}}],"selectables":null,"selection":null,"elements":[{"order":0,"task":{"label":"Lampe einschalten","device_description":{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:79e7914b-f303-4a7d-90af-dee70db05fd9","name":"setOnStateFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},"input":null,"bpmn_element_id":"Task_1n8uaf6","multi_task":false,"selected_service":null,"parameter":{},"configurables":null}},{"order":0,"task":{"label":"Lampe ausschalten","device_description":{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:2f35150b-9df7-4cad-95bc-165fa00219fd","name":"setOffStateFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},"input":null,"bpmn_element_id":"Task_103ufqc","multi_task":false,"selected_service":null,"parameter":{},"configurables":null}}]}}]}

}

func TestProcessDescription(t *testing.T) {
	file, err := ioutil.ReadFile("../../tests/resources/description/process_description.bpmn")
	if err != nil {
		t.Error(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		t.Error(err)
		return
	}

	if result.Description != "desc" {
		t.Fatal("error in description")
	}

	fmt.Println(result.Description)
}

func TestCollaborationDescription(t *testing.T) {
	file, err := ioutil.ReadFile("../../tests/resources/description/collaboration_description.bpmn")
	if err != nil {
		t.Error(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		t.Error(err)
		return
	}

	if result.Description != "desc test" {
		t.Fatal("error in description")
	}

	fmt.Println(result.Description)
}

func TestProcessNoDescription(t *testing.T) {
	file, err := ioutil.ReadFile("../../tests/resources/description/process_noDescription.bpmn")
	if err != nil {
		t.Error(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		t.Error(err)
		return
	}

	if result.Description != "" {
		t.Fatal("error in description")
	}

	fmt.Println(result.Description)
}

func TestCollaborationNoDescription(t *testing.T) {
	file, err := ioutil.ReadFile("../../tests/resources/description/collaboration_noDescription.bpmn")
	if err != nil {
		t.Error(err)
		return
	}
	result, err := PrepareDeployment(string(file))
	if err != nil {
		t.Error(err)
		return
	}

	if result.Description != "" {
		t.Fatal("error in description")
	}

	fmt.Println(result.Description)
}
