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

package devicemodel

import "github.com/SENERGY-Platform/models/go/models"

type Device = models.Device

type DeviceType = models.DeviceType

type ServiceGroup = models.ServiceGroup

type Attribute = models.Attribute

type Service = models.Service

type Interaction = models.Interaction

const (
	EVENT             = models.EVENT
	REQUEST           = models.REQUEST
	EVENT_AND_REQUEST = models.EVENT_AND_REQUEST
)

type Protocol = models.Protocol

type ProtocolSegment = models.ProtocolSegment

type DeviceGroup = models.DeviceGroup

type FilterCriteria struct {
	FunctionId    string `json:"function_id"`
	AspectId      string `json:"aspect_id"`
	DeviceClassId string `json:"device_class_id"`
}

func (this FilterCriteria) Short() string {
	return this.FunctionId + "_" + this.AspectId + "_" + this.DeviceClassId
}

type DeviceGroupFilterCriteria = models.DeviceGroupFilterCriteria
