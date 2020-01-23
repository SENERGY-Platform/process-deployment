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
	mock "github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"time"
)

const TestToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjb25uZWN0aXZpdHktdGVzdCJ9.OnihzQ7zwSq0l1Za991SpdsxkktfrdlNl-vHHpYpXQw"

func ExampleCtrl_PrepareDeploymentById() {
	config, err := config.LoadConfig("../config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	port, err := tests.GetFreePort()
	if err != nil {
		fmt.Println(err)
		return
	}
	config.ApiPort = strconv.Itoa(port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = Start(ctx, config, mock.Kafka, mock.Database, mock.Devices, mock.ProcessModelRepo)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(100 * time.Millisecond) //wait for api startup

	file, err := ioutil.ReadFile("./tests/resources/lanes.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}

	processModelId := uuid.NewV4().String()
	mock.ProcessModelRepo.SetProcessModel(processModelId, model.ProcessModel{Id: processModelId, BpmnXml: string(file), SvgXml: "<svg/>"})

	mock.Devices.SetOptions([]model.Selectable{
		{
			Device: devicemodel.Device{
				Id: "device1",
			},
			Services: []devicemodel.Service{
				{
					Id: "service1",
				},
			},
		},
	})

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(
		"GET",
		"http://localhost:"+config.ApiPort+"/prepared-deployments/"+url.PathEscape(processModelId),
		nil,
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

	deployment := model.Deployment{}

	err = json.NewDecoder(resp.Body).Decode(&deployment)

	if err != nil {
		fmt.Println(err)
		return
	}
	deployment.XmlRaw = ""
	deployment.Xml = ""

	msg, err := json.Marshal(deployment)

	fmt.Println(err, string(msg))

	//output:
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"\u003csvg/\u003e","name":"lanes","elements":null,"lanes":[{"order":0,"multi_lane":{"label":"multiTaskLane","bpmn_element_id":"Lane_12774cv","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selections":null,"elements":[{"order":3,"task":{"label":"multi_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_084s3g5","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}},{"order":4,"task":{"label":"multi_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_098jmqp","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}}]}},{"order":0,"lane":{"label":"MixedLane","bpmn_element_id":"Lane_0odlj5k","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":5,"task":{"label":"mixed_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1npvonw","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}},{"order":6,"task":{"label":"mixed_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1mnjsed","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}}]}},{"order":1,"lane":{"label":"taskLane","bpmn_element_id":"Lane_0v679jg","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":1,"task":{"label":"lane_task_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_0nmb2on","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}},{"order":2,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_1tchutl","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}}]}}]}

}

func ExampleCtrl_PrepareDeployment() {
	config, err := config.LoadConfig("../config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	port, err := tests.GetFreePort()
	if err != nil {
		fmt.Println(err)
		return
	}
	config.ApiPort = strconv.Itoa(port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = Start(ctx, config, mock.Kafka, mock.Database, mock.Devices, mock.ProcessModelRepo)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(100 * time.Millisecond) //wait for api startup

	file, err := ioutil.ReadFile("./tests/resources/lanes.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}

	mock.Devices.SetOptions([]model.Selectable{
		{
			Device: devicemodel.Device{
				Id: "device1",
			},
			Services: []devicemodel.Service{
				{
					Id: "service1",
				},
			},
		},
	})

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	preparereq := model.PrepareRequest{Xml: string(file), Svg: "<svg/>"}
	temp, err := json.Marshal(preparereq)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:"+config.ApiPort+"/prepared-deployments",
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

	deployment := model.Deployment{}

	err = json.NewDecoder(resp.Body).Decode(&deployment)

	if err != nil {
		fmt.Println(err)
		return
	}
	deployment.XmlRaw = ""
	deployment.Xml = ""

	msg, err := json.Marshal(deployment)

	fmt.Println(err, string(msg))

	//output:
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"\u003csvg/\u003e","name":"lanes","elements":null,"lanes":[{"order":0,"multi_lane":{"label":"multiTaskLane","bpmn_element_id":"Lane_12774cv","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selections":null,"elements":[{"order":3,"task":{"label":"multi_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_084s3g5","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}},{"order":4,"task":{"label":"multi_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_098jmqp","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}}]}},{"order":0,"lane":{"label":"MixedLane","bpmn_element_id":"Lane_0odlj5k","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":5,"task":{"label":"mixed_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1npvonw","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}},{"order":6,"task":{"label":"mixed_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1mnjsed","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}}]}},{"order":1,"lane":{"label":"taskLane","bpmn_element_id":"Lane_0v679jg","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":1,"task":{"label":"lane_task_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_0nmb2on","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}},{"order":2,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_1tchutl","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}}]}}]}

}

func ExampleCtrl_PrepareDeployment2() {
	config, err := config.LoadConfig("../config.json")
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	port, err := tests.GetFreePort()
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	config.ApiPort = strconv.Itoa(port)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	defer log.Println("stop")

	//time.Sleep(1 * time.Second)

	err = Start(ctx, config, mock.Kafka, mock.Database, mock.Devices, mock.ProcessModelRepo)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	time.Sleep(100 * time.Millisecond) //wait for api startup

	file, err := ioutil.ReadFile("./tests/resources/lanes-2.bpmn")
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	mock.Devices.SetOptions([]model.Selectable{
		{
			Device: devicemodel.Device{
				Id: "device1",
			},
			Services: []devicemodel.Service{
				{
					Id: "service1",
				},
			},
		},
	})

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	preparereq := model.PrepareRequest{Xml: string(file), Svg: "<svg/>"}
	temp, err := json.Marshal(preparereq)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:"+config.ApiPort+"/prepared-deployments",
		bytes.NewBuffer(temp),
	)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", TestToken)

	log.Println("request prepared deployment")
	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	deployment := model.Deployment{}
	if resp.StatusCode != 200 {
		t, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(string(t))
	} else {
		err = json.NewDecoder(resp.Body).Decode(&deployment)
	}

	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	deployment.XmlRaw = ""
	deployment.Xml = ""

	msg, err := json.Marshal(deployment)

	fmt.Println(err, string(msg))

	//output:
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"\u003csvg/\u003e","name":"Lane_Test","elements":null,"lanes":[{"order":0,"lane":{"label":"Lane_0sswcoy","bpmn_element_id":"Lane_0sswcoy","device_descriptions":[{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:66e62e2f-39d2-4f5d-bcbd-847ac8f8e1b7","name":"newOffFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},{"characteristic_id":"urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43","function":{"id":"urn:infai:ses:controlling-function:08964911-56ff-4922-a2d2-a0d7f9f5f58d","name":"newColorFunction","concept_id":"urn:infai:ses:concept:8b1161d5-7878-4dd2-a36c-6f98f6b94bf8","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:b298a6ff-0f13-4e00-8c71-ae0c2fd8a5da","name":"newOnFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":0,"task":{"label":"Lamp newOffFunction","device_description":{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:66e62e2f-39d2-4f5d-bcbd-847ac8f8e1b7","name":"newOffFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},"input":null,"bpmn_element_id":"Task_07asof2","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}},{"order":0,"task":{"label":"Lamp newColorFunction","device_description":{"characteristic_id":"urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43","function":{"id":"urn:infai:ses:controlling-function:08964911-56ff-4922-a2d2-a0d7f9f5f58d","name":"newColorFunction","concept_id":"urn:infai:ses:concept:8b1161d5-7878-4dd2-a36c-6f98f6b94bf8","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},"input":{"b":0,"g":0,"r":0},"bpmn_element_id":"Task_04u7f2g","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{"inputs.b":"0","inputs.g":"0","inputs.r":"0"},"configurables":null}},{"order":0,"task":{"label":"Lamp newOnFunction","device_description":{"characteristic_id":"","function":{"id":"urn:infai:ses:controlling-function:b298a6ff-0f13-4e00-8c71-ae0c2fd8a5da","name":"newOnFunction","concept_id":"","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},"input":null,"bpmn_element_id":"Task_14r4c1p","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{},"configurables":null}}]}}]}
}

func ExampleCtrl_PrepareDeployment3() {
	config, err := config.LoadConfig("../config.json")
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	port, err := tests.GetFreePort()
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	config.ApiPort = strconv.Itoa(port)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	defer log.Println("stop")

	//time.Sleep(1 * time.Second)

	err = Start(ctx, config, mock.Kafka, mock.Database, mock.Devices, mock.ProcessModelRepo)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	time.Sleep(100 * time.Millisecond) //wait for api startup

	file, err := ioutil.ReadFile("./tests/resources/lanes-3.bpmn")
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	mock.Devices.SetOptions([]model.Selectable{
		{
			Device: devicemodel.Device{
				Id: "device1",
			},
			Services: []devicemodel.Service{
				{
					Id: "service1",
				},
			},
		},
	})

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	preparereq := model.PrepareRequest{Xml: string(file), Svg: "<svg/>"}
	temp, err := json.Marshal(preparereq)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:"+config.ApiPort+"/prepared-deployments",
		bytes.NewBuffer(temp),
	)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", TestToken)

	log.Println("request prepared deployment")
	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	deployment := model.Deployment{}
	if resp.StatusCode != 200 {
		t, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(string(t))
	} else {
		err = json.NewDecoder(resp.Body).Decode(&deployment)
	}

	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	deployment.XmlRaw = ""
	deployment.Xml = ""

	msg, err := json.Marshal(deployment)

	fmt.Println(err, string(msg))

	//output:
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"\u003csvg/\u003e","name":"Collaboration_0wfm1tv","elements":null,"lanes":[{"order":0,"lane":{"label":"Process_1","bpmn_element_id":"Process_1","device_descriptions":[{"characteristic_id":"urn:infai:ses:characteristic:72b624b5-6edc-4ec4-9ad9-fa00b39915c0","function":{"id":"urn:infai:ses:controlling-function:7adc7f29-5c37-4bfc-8508-6130a143ac66","name":"brightnessFunction","concept_id":"urn:infai:ses:concept:dbe4ad57-aa1d-4d24-9bee-a44a1c670d7f","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":0,"task":{"label":"Lamp brightnessFunction","device_description":{"characteristic_id":"urn:infai:ses:characteristic:72b624b5-6edc-4ec4-9ad9-fa00b39915c0","function":{"id":"urn:infai:ses:controlling-function:7adc7f29-5c37-4bfc-8508-6130a143ac66","name":"brightnessFunction","concept_id":"urn:infai:ses:concept:dbe4ad57-aa1d-4d24-9bee-a44a1c670d7f","rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"},"device_class":{"id":"urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86","name":"Lamp","rdf_type":"https://senergy.infai.org/ontology/DeviceClass"}},"input":0,"bpmn_element_id":"Task_1d0tawd","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{"inputs":"0"},"configurables":null}}]}}]}
}

func ExampleCtrl_PrepareDeploymentOfEmptyLane() {
	config, err := config.LoadConfig("../config.json")
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	port, err := tests.GetFreePort()
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}
	config.ApiPort = strconv.Itoa(port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = Start(ctx, config, mock.Kafka, mock.Database, mock.Devices, mock.ProcessModelRepo)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	time.Sleep(100 * time.Millisecond) //wait for api startup

	file, err := ioutil.ReadFile("./tests/resources/lane_only_timer.bpmn")
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	mock.Devices.SetOptions([]model.Selectable{
		{
			Device: devicemodel.Device{
				Id: "device1",
			},
			Services: []devicemodel.Service{
				{
					Id: "service1",
				},
			},
		},
	})

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	preparereq := model.PrepareRequest{Xml: string(file), Svg: "<svg/>"}
	temp, err := json.Marshal(preparereq)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:"+config.ApiPort+"/prepared-deployments",
		bytes.NewBuffer(temp),
	)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}
	req.Header.Set("Authorization", TestToken)

	log.Println("request prepared deployment")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}
	defer resp.Body.Close()

	deployment := model.Deployment{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	err = json.Unmarshal(b, &deployment)
	if err != nil {
		fmt.Println(err, string(b))
		debug.PrintStack()
		return
	}
	deployment.XmlRaw = ""
	deployment.Xml = ""

	msg, err := json.Marshal(deployment)
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	fmt.Println(string(msg))

	//output:
	//{"id":"","xml_raw":"","xml":"","svg":"\u003csvg/\u003e","name":"Lane_Timer_FJ","elements":null,"lanes":[{"order":0,"lane":{"label":"Lane_Timer_fj","bpmn_element_id":"Lane_Timer_fj","device_descriptions":null,"selectables":[],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":0,"time_event":{"bpmn_element_id":"IntermediateThrowEvent_1opksgz","kind":"timeDuration","time":"","label":"IntermediateThrowEvent_1opksgz"}}]}}]}

}
