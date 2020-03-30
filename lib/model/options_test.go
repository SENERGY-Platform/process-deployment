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

package model

import (
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
)

func ExampleToFilter() {
	fmt.Println(DeviceDescriptions{
		{
			CharacteristicId: "foobar",
			Function:         devicemodel.Function{Id: "f1"},
			DeviceClass:      &devicemodel.DeviceClass{Id: "dc1"},
			Aspect:           &devicemodel.Aspect{Id: "a1"},
		},
		{
			CharacteristicId: "foobar",
			Function:         devicemodel.Function{Id: "f1"},
			DeviceClass:      &devicemodel.DeviceClass{Id: "dc2"},
			Aspect:           &devicemodel.Aspect{Id: "a2"},
		},
		{
			CharacteristicId: "foobar",
			Function:         devicemodel.Function{Id: "f2"},
			DeviceClass:      &devicemodel.DeviceClass{Id: "dc3"},
			Aspect:           &devicemodel.Aspect{Id: "a2"},
		},
	}.ToFilter())

	//output:
	//[{f1 dc1 a1} {f1 dc2 a2} {f2 dc3 a2}]
}

func ExampleNilToFilter() {
	fmt.Println(DeviceDescriptions{
		{
			CharacteristicId: "foobar",
			Function:         devicemodel.Function{Id: "f1"},
			DeviceClass:      &devicemodel.DeviceClass{Id: "dc1"},
			Aspect:           &devicemodel.Aspect{Id: "a1"},
		},
		{},
	}.ToFilter())

	//output:
	//[{f1 dc1 a1}]
}
