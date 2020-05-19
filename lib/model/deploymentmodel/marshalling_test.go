/*
 * Copyright 2020 InfAI (CC SES)
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

package deploymentmodel

import (
	"encoding/json"
	"fmt"
)

func ExampleMarshalling() {
	laneId := "lane-id"
	deployment := Deployment{
		Name:       "example deployment",
		Diagram:    Diagram{},
		Executable: false,
		Elements: []Element{
			{
				BpmnId:     "task_bpmn_id",
				LaneBpmnId: &laneId,
				Name:       "task",
				Order:      0,
				Task: &Task{
					Retries:   3,
					Parameter: map[string]string{"inputs": `"ff0"`},
					Configurables: []Configurable{
						{
							CharacteristicId: "example_duration",
							Values: []ConfigurableValue{
								{
									Name:      "duration",
									Path:      "path.to.duration",
									Value:     42,
									ValueType: "https://schema.org/Integer",
								},
							},
						},
					},
					Selection: Selection{
						FilterCriteria: FilterCriteria{},
						SelectionOptions: []SelectionOption{{
							Device: Device{
								Id:   "did",
								Name: "device",
							},
							Services: []Service{{
								Id:   "sid",
								Name: "service",
							}},
						}},
						SelectedDeviceId:  "sdid",
						SelectedServiceId: "ssid",
					},
				},
			},
		},
	}

	out, err := json.MarshalIndent(deployment, "", "    ")
	fmt.Print(err, string(out))

	//output:
	//<nil>{
	//     "id": "",
	//     "name": "example deployment",
	//     "description": "",
	//     "diagram": {
	//         "xml_raw": "",
	//         "xml_deployed": "",
	//         "svg": ""
	//     },
	//     "elements": [
	//         {
	//             "bpmn_id": "task_bpmn_id",
	//             "lane_bpmn_id": "lane-id",
	//             "name": "task",
	//             "order": 0,
	//             "time_event": null,
	//             "notification": null,
	//             "message_event": null,
	//             "task": {
	//                 "retries": 3,
	//                 "parameter": {
	//                     "inputs": "\"ff0\""
	//                 },
	//                 "configurables": [
	//                     {
	//                         "characteristic_id": "example_duration",
	//                         "values": [
	//                             {
	//                                 "name": "duration",
	//                                 "path": "path.to.duration",
	//                                 "value": 42,
	//                                 "value_type": "https://schema.org/Integer"
	//                             }
	//                         ]
	//                     }
	//                 ],
	//                 "selection": {
	//                     "filter_criteria": {
	//                         "characteristic_id": null,
	//                         "function_id": null,
	//                         "device_class_id": null,
	//                         "aspect_id": null
	//                     },
	//                     "selection_options": [
	//                         {
	//                             "device": {
	//                                 "id": "did",
	//                                 "name": "device"
	//                             },
	//                             "services": [
	//                                 {
	//                                     "id": "sid",
	//                                     "name": "service"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "selected_device_id": "sdid",
	//                     "selected_service_id": "ssid"
	//                 }
	//             }
	//         }
	//     ],
	//     "executable": false
	// }
}

func strptr(in string) (out *string) {
	return &in
}
