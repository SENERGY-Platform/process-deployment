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
	"io/ioutil"
)

func ExampleSimpleBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/simple.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := BpmnToDeployment(string(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	result.Xml = ""
	temp, err := json.Marshal(result)
	fmt.Println(err, string(temp))

	//output:
	//<nil> {"id":"","xml":"","name":"simple","elements":[{"order":1,"task":{"label":"taskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_0wjr1fj","input":"000","device_options":null,"selection":{"selected_device":{"id":"","local_id":"","name":"","device_type_id":""},"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}},"Parameter":{"inputs":"\"ff0\""}}},{"order":2,"multi_task":{"label":"multiTaskLabel","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"bpmn_element_id":"Task_096xjeg","input":"000","device_options":null,"selections":null,"Parameter":{"inputs":"\"fff\""}}},{"order":3,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}}],"lanes":null}
}

func ExampleLaneBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/lanes.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := BpmnToDeployment(string(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	result.Xml = ""
	temp, err := json.Marshal(result)
	fmt.Println(err, string(temp))

	//output:
	//<nil> {"id":"","xml":"","name":"lanes","elements":null,"lanes":[{"order":0,"multi_lane":{"label":"multiTaskLane","bpmn_element_id":"Lane_12774cv","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"device_options":null,"selections":null,"elements":[{"order":3,"task":{"label":"multi_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_084s3g5","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"Parameter":{}}},{"order":4,"task":{"label":"multi_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_098jmqp","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"Parameter":{}}}]}},{"order":0,"lane":{"label":"MixedLane","bpmn_element_id":"Lane_0odlj5k","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"device_options":null,"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":5,"task":{"label":"mixed_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1npvonw","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"Parameter":{}}},{"order":6,"task":{"label":"mixed_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1mnjsed","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"Parameter":{}}}]}},{"order":1,"lane":{"label":"taskLane","bpmn_element_id":"Lane_0v679jg","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}}],"device_options":null,"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":1,"task":{"label":"lane_task_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_0nmb2on","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"Parameter":{}}},{"order":2,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_1tchutl","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}}]}}]}
}

func ExampleTimeAndReceiverBpmnToDeployment() {
	file, err := ioutil.ReadFile("../../tests/resources/timeAndReceive.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := BpmnToDeployment(string(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	result.Xml = ""
	temp, err := json.Marshal(result)
	fmt.Println(err, string(temp))

	//output:
	//<nil> {"id":"","xml":"","name":"timeAndReceive","elements":[{"order":0,"receive_task_event":{"label":"Task_1uyyxb0","bpmn_element_id":"Task_1uyyxb0","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}},{"order":0,"time_event":{"bpmn_element_id":"IntermediateThrowEvent_10mhx3e","kind":"timeDuration","time":"PT1M","label":"eine Minute"}}],"lanes":null}
}
