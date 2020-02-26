package model

import "github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"

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

type DeviceTypesFilter struct {
	FunctionIds   []string `json:"function_ids"`
	AspectIds     []string `json:"aspect_ids"`
	DeviceClassId string   `json:"device_class_id"`
}

func (this DeviceDescriptions) ToFilter() (result DeviceTypesFilter) {
	aspectIndex := map[string]bool{}
	functionIndex := map[string]bool{}
	for _, description := range this {
		if description.Aspect != nil {
			if _, notDistinct := aspectIndex[description.Aspect.Id]; !notDistinct {
				result.AspectIds = append(result.AspectIds, description.Aspect.Id)
				aspectIndex[description.Aspect.Id] = true
			}
		}
		if _, notDistinct := functionIndex[description.Function.Id]; !notDistinct {
			result.FunctionIds = append(result.FunctionIds, description.Function.Id)
			functionIndex[description.Function.Id] = true
		}
		if description.DeviceClass != nil {
			result.DeviceClassId = description.DeviceClass.Id
		}

	}
	return
}
