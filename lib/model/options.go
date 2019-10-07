package model

import "github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"

type Selectable struct {
	Device   devicemodel.Device    `json:"device"`
	Services []devicemodel.Service `json:"services"`
}

type Selection struct {
	Device  devicemodel.Device  `json:"device"`
	Service devicemodel.Service `json:"service"`
}

type DeviceDescription struct {
	CharacteristicId string                   `json:"characteristic_id"`
	Function         devicemodel.Function     `json:"function"`
	DeviceClass      *devicemodel.DeviceClass `json:"device_class,omitempty"`
	Aspect           *devicemodel.Aspect      `json:"aspect,omitempty"`
}

type DeviceTypesFilter struct {
	FunctionId    string `json:"function_id"`
	DeviceClassId string `json:"device_class_id"`
	AspectId      string `json:"aspect_id"`
}

func (description *DeviceDescription) ToFilter() (filter DeviceTypesFilter) {
	filter.FunctionId = description.Function.Id
	if description.Aspect != nil {
		filter.AspectId = description.Aspect.Id
	}
	if description.DeviceClass != nil {
		filter.DeviceClassId = description.DeviceClass.Id
	}
	return
}
