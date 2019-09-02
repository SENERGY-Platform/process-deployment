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
	"html"
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

	mock.DeviceRepository.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, true, mock.DeviceRepository)
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
	//     <bpmn:process id="simple" isExecutable="true">
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
	//             <camunda:executionListener event="start"><camunda:script scriptFormat="groovy">execution.setVariable("collection", ["{\"device\":{\"id\":\"device_id_1\",\"local_id\":\"\",\"name\":\"\",\"device_type_id\":\"\"},\"service\":{\"id\":\"service_id_1\",\"local_id\":\"\",\"name\":\"\",\"description\":\"\",\"aspects\":null,\"protocol_id\":\"pid\",\"inputs\":null,\"outputs\":null,\"functions\":null,\"rdf_type\":\"\"},\"protocol\":{\"id\":\"pid\",\"name\":\"protocol1\",\"handler\":\"p\",\"protocol_segments\":null}}","{\"device\":{\"id\":\"device_id_2\",\"local_id\":\"\",\"name\":\"\",\"device_type_id\":\"\"},\"service\":{\"id\":\"service_id_2\",\"local_id\":\"\",\"name\":\"\",\"description\":\"\",\"aspects\":null,\"protocol_id\":\"pid\",\"inputs\":null,\"outputs\":null,\"functions\":null,\"rdf_type\":\"\"},\"protocol\":{\"id\":\"pid\",\"name\":\"protocol1\",\"handler\":\"p\",\"protocol_segments\":null}}"])</camunda:script></camunda:executionListener></bpmn:extensionElements>
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

	mock.DeviceRepository.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, false, mock.DeviceRepository)
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
	//     <bpmn:process id="simple" isExecutable="true">
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
	//                     <camunda:inputParameter name="payload">{"function":{"id":"fid","name":"","concept_id":"","rdf_type":""},"characteristic_id":"example_hex","device":{"id":"device_id_1","local_id":"","name":"","device_type_id":""},"service":{"id":"service_id_1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"pid","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"protocol":{"id":"pid","name":"protocol1","handler":"p","protocol_segments":null},"input":"000"}</camunda:inputParameter>
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
	file, err := ioutil.ReadFile("../../tests/resources/timeAndReceive.json")
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

	mock.DeviceRepository.SetProtocol("pid", devicemodel.Protocol{Id: "pid", Handler: "p", Name: "protocol1"})

	err = UseDeploymentSelections(&deployment, false, mock.DeviceRepository)
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
	//     <bpmn:process id="timeAndReceive" isExecutable="true">
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
