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

package ctrl

import (
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
)

func Example_getCharacteristicOfPathInService() {
	service := devicemodel.Service{
		Id:         "service_id_3",
		LocalId:    "s3url",
		Name:       "s3",
		ProtocolId: "pid",
		Outputs: []devicemodel.Content{
			{
				ContentVariable: devicemodel.ContentVariable{
					Name: "payload",
					Type: devicemodel.Structure,
					SubContentVariables: []devicemodel.ContentVariable{
						{
							Name:             "kelvin",
							Type:             devicemodel.String,
							CharacteristicId: "example_hex",
						},
					},
				},
			},
		},
	}
	path := "value.payload.kelvin"

	result, err := getCharacteristicOfPathInService(&service, path)

	fmt.Println(err, result)

	//output:
	//<nil> example_hex
}

func Example_getCharacteristicOfPathInService2() {
	service := devicemodel.Service{
		Id:         "service_id_3",
		LocalId:    "s3url",
		Name:       "s3",
		ProtocolId: "pid",
		Outputs: []devicemodel.Content{
			{
				ContentVariable: devicemodel.ContentVariable{
					Name:             "kelvin",
					Type:             devicemodel.String,
					CharacteristicId: "example_hex",
				},
			},
		},
	}
	path := "value.kelvin"

	result, err := getCharacteristicOfPathInService(&service, path)

	fmt.Println(err, result)

	//output:
	//<nil> example_hex
}
