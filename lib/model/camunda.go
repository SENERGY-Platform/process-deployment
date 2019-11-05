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
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
)

type Command struct {
	//modeling time
	Function         devicemodel.Function `json:"function"`
	CharacteristicId string               `json:"characteristic_id"`

	//optional modeling time (used to limit/filter device and service selection in deployment)
	DeviceClass *devicemodel.DeviceClass `json:"device_class,omitempty"`
	Aspect      *devicemodel.Aspect      `json:"aspect,omitempty"`

	//deployment time
	DeviceId   string                `json:"device_id,omitempty"`
	Device     *devicemodel.Device   `json:"device,omitempty"`
	ServiceId  string                `json:"service_id,omitempty"`
	Service    *devicemodel.Service  `json:"service,omitempty"`
	Protocol   *devicemodel.Protocol `json:"protocol,omitempty"`
	ProtocolId string                `json:"protocol_id,omitempty"`

	//runtime
	Input  interface{} `json:"input,omitempty"`
	Output interface{} `json:"output,omitempty"`

	Retries int64 `json:"retries,omitempty"`
}

type Documentation struct {
	Order int64 `json:"order"`
}

type Overwrite struct {
	DeviceId   string                `json:"device_id,omitempty"`
	Device     *devicemodel.Device   `json:"device,omitempty"`
	ServiceId  string                `json:"service_id,omitempty"`
	Service    *devicemodel.Service  `json:"service,omitempty"`
	Protocol   *devicemodel.Protocol `json:"protocol,omitempty"`
	ProtocolId string                `json:"protocol_id,omitempty"`
}
