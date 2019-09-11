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
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/tests"
	mock "github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const TestToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjb25uZWN0aXZpdHktdGVzdCJ9.OnihzQ7zwSq0l1Za991SpdsxkktfrdlNl-vHHpYpXQw"

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
	defer log.Println("stop")

	//time.Sleep(1 * time.Second)

	err = Start(ctx, config, mock.Kafka, mock.Database, mock.Devices)
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
	//<nil> {"id":"","xml_raw":"","xml":"","svg":"\u003csvg/\u003e","name":"lanes","elements":null,"lanes":[{"order":0,"multi_lane":{"label":"multiTaskLane","bpmn_element_id":"Lane_12774cv","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selections":null,"elements":[{"order":3,"task":{"label":"multi_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_084s3g5","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}},{"order":4,"task":{"label":"multi_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_098jmqp","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}}]}},{"order":0,"lane":{"label":"MixedLane","bpmn_element_id":"Lane_0odlj5k","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":5,"task":{"label":"mixed_lane_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid1","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1npvonw","multi_task":true,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}},{"order":6,"task":{"label":"mixed_lane_2","device_description":{"characteristic_id":"example_hex","function":{"id":"fid2","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_1mnjsed","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}}]}},{"order":1,"lane":{"label":"taskLane","bpmn_element_id":"Lane_0v679jg","device_descriptions":[{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}}],"selectables":[{"device":{"id":"device1","local_id":"","name":"","device_type_id":""},"services":[{"id":"service1","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""}]}],"selection":{"id":"","local_id":"","name":"","device_type_id":""},"elements":[{"order":1,"task":{"label":"lane_task_1","device_description":{"characteristic_id":"example_hex","function":{"id":"fid","name":"","concept_id":"","rdf_type":""}},"input":"000","bpmn_element_id":"Task_0nmb2on","multi_task":false,"selected_service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"parameter":{}}},{"order":2,"msg_event":{"label":"eventName","bpmn_element_id":"IntermediateThrowEvent_1tchutl","device":{"id":"","local_id":"","name":"","device_type_id":""},"service":{"id":"","local_id":"","name":"","description":"","aspects":null,"protocol_id":"","inputs":null,"outputs":null,"functions":null,"rdf_type":""},"path":"","value":"","operation":"","event_id":""}}]}}]}

}
