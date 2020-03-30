package model

import (
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"reflect"
)

type Selectable struct {
	Device   devicemodel.Device    `json:"device"`
	Services []devicemodel.Service `json:"services"`
}

type Selection struct {
	Device  *devicemodel.Device  `json:"device"`
	Service *devicemodel.Service `json:"service"`
}

type DeviceDescriptions []DeviceDescription
type DeviceDescription struct {
	CharacteristicId string                   `json:"characteristic_id"`
	Function         devicemodel.Function     `json:"function"`
	DeviceClass      *devicemodel.DeviceClass `json:"device_class,omitempty"`
	Aspect           *devicemodel.Aspect      `json:"aspect,omitempty"`
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
