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
	//<nil> {"id":"","xml":"","name":"simple","elements":[{"label":"taskLabel","characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""},"order":1,"bpmn_element_id":"Task_0wjr1fj","input":"000","device_options":null,"selections":null,"Parameter":{"inputs":"\"ff0\""}},{"label":"multiTaskLabel","characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""},"order":2,"bpmn_element_id":"Task_096xjeg","input":"000","device_options":null,"selections":null,"Parameter":{"inputs":"\"fff\""}},{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_0905jg5","order":3,"device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}],"lanes":null}
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
	//
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
	//<nil> {"id":"","xml":"","name":"timeAndReceive","elements":[{"label":"Task_1uyyxb0","bpmn_element_id":"Task_1uyyxb0","order":0,"device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""},{"bpmn_element_id":"IntermediateThrowEvent_10mhx3e","kind":"timeDuration","time":"PT1M","label":"eine Minute","order":0}],"lanes":null}
}
