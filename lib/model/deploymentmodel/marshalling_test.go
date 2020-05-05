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
	deployment := Deployment{
		Name:       "example deployment",
		Diagram:    Diagram{},
		Executable: false,
		Pools: []Pool{
			{
				BaseInfo: BaseInfo{
					Name:   "pool name",
					BpmnId: "pool_bpmn_id",
					Order:  42,
				},
				Lanes: []Lane{
					{
						BaseInfo: BaseInfo{
							Name:   "lane name",
							BpmnId: "lane_bpmn_id",
							Order:  13,
						},
						Elements: []Element{
							{
								BaseInfo: BaseInfo{
									Name:   "task",
									BpmnId: "task_bpmn_id",
									Order:  0,
								},
								Task: &Task{
									Retries:   3,
									Input:     "000",
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
									SelectedServiceId: "",
								},
							},
						},
						FilterCriteria: []FilterCriteria{
							{
								CharacteristicId: strptr("example_color_hex"),
								FunctionId:       strptr("function_id"),
								DeviceClassId:    nil,
								AspectId:         nil,
							},
						},
						Selectables:      nil,
						SelectedDeviceId: "",
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
	//     "diagram": {
	//         "xml_raw": "",
	//         "xml_deployed": "",
	//         "svg": ""
	//     },
	//     "pools": [
	//         {
	//             "name": "pool name",
	//             "bpmn_id": "pool_bpmn_id",
	//             "order": 42,
	//             "lanes": [
	//                 {
	//                     "name": "lane name",
	//                     "bpmn_id": "lane_bpmn_id",
	//                     "order": 13,
	//                     "elements": [
	//                         {
	//                             "name": "task",
	//                             "bpmn_id": "task_bpmn_id",
	//                             "order": 0,
	//                             "time_event": null,
	//                             "notification": null,
	//                             "message_event": null,
	//                             "task": {
	//                                 "retries": 3,
	//                                 "input": "000",
	//                                 "parameter": {
	//                                     "inputs": "\"ff0\""
	//                                 },
	//                                 "configurables": [
	//                                     {
	//                                         "characteristic_id": "example_duration",
	//                                         "values": [
	//                                             {
	//                                                 "name": "duration",
	//                                                 "path": "path.to.duration",
	//                                                 "value": 42,
	//                                                 "value_type": "https://schema.org/Integer"
	//                                             }
	//                                         ]
	//                                     }
	//                                 ],
	//                                 "selected_service_id": ""
	//                             }
	//                         }
	//                     ],
	//                     "filter_criteria": [
	//                         {
	//                             "characteristic_id": "example_color_hex",
	//                             "function_id": "function_id",
	//                             "device_class_id": null,
	//                             "aspect_id": null
	//                         }
	//                     ],
	//                     "selectables": null,
	//                     "selected_device_id": ""
	//                 }
	//             ]
	//         }
	//     ],
	//     "executable": false
	//}
}

func strptr(in string) (out *string) {
	return &in
}
