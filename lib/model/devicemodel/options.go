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

package devicemodel

import (
	"errors"
	"reflect"
)

type Selectable struct {
	Device   Device    `json:"device"`
	Services []Service `json:"services"`
}

type Selection struct {
	Device  *Device  `json:"device"`
	Service *Service `json:"service"`
}

type DeviceDescriptions []DeviceDescription
type DeviceDescription struct {
	CharacteristicId string       `json:"characteristic_id"`
	Function         Function     `json:"function"`
	DeviceClass      *DeviceClass `json:"device_class,omitempty"`
	Aspect           *Aspect      `json:"aspect,omitempty"`
}

type DeviceTypesFilter []DeviceTypeFilterElement

type DeviceTypeFilterElement struct {
	FunctionId    string `json:"function_id"`
	DeviceClassId string `json:"device_class_id"`
	AspectId      string `json:"aspect_id"`
}

func (this DeviceDescriptions) ToFilter() (result DeviceTypesFilter) {
	for _, element := range this {
		newElement := DeviceTypeFilterElement{
			FunctionId: element.Function.Id,
		}
		if element.DeviceClass != nil {
			newElement.DeviceClassId = element.DeviceClass.Id
		}
		if element.Aspect != nil {
			newElement.AspectId = element.Aspect.Id
		}
		if !IsZero(element) {
			result = append(result, newElement)
		}
	}
	return
}

func IsZero(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

func (this Selection) Validate(strict bool) error {
	if this.Device == nil {
		return errors.New("missing device selection")
	}
	if this.Service == nil {
		return errors.New("missing service selection")
	}
	if this.Device.Id == "" {
		return errors.New("missing device id in selection")
	}
	if this.Service.Id == "" {
		return errors.New("missing service id in selection")
	}
	return nil
}
