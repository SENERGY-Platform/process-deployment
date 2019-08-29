package model

import "github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"

type DeviceOption struct {
	Device         devicemodel.Device    `json:"device" bson:"device"`
	ServiceOptions []devicemodel.Service `json:"service_options" bson:"-"`
}

type Selection struct {
	SelectedDevice  devicemodel.Device  `json:"selected_device" bson:"selected_device"`
	SelectedService devicemodel.Service `json:"selected_service" bson:"selected_service"`
}

type DeviceDescription struct {
	CharacteristicId string                   `json:"characteristic_id" bson:"characteristic_id"`
	Function         devicemodel.Function     `json:"function" bson:"function"`
	DeviceClass      *devicemodel.DeviceClass `json:"device_class,omitempty" bson:"device_class,omitempty"`
	Aspect           *devicemodel.Aspect      `json:"aspect,omitempty" bson:"aspect,omitempty"`
}
