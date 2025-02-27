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

package executionmodel

import (
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"
)

type Task struct {
	Version int64 `json:"version"`

	//modeling time
	Function         devicemodel.Function `json:"function"`
	CharacteristicId string               `json:"characteristic_id"`

	//optional modeling time (used to limit/filter device and service selection in deployment)
	DeviceClass *devicemodel.DeviceClass `json:"device_class,omitempty"`
	Aspect      *devicemodel.AspectNode  `json:"aspect,omitempty"`

	//deployment time
	DeviceGroupId string `json:"device_group_id,omitempty"`
	DeviceId      string `json:"device_id,omitempty"`
	ServiceId     string `json:"service_id,omitempty"`
	ProtocolId    string `json:"protocol_id,omitempty"`

	//version >= 3
	InputPaths      []string                            `json:"input_paths,omitempty"`
	OutputPath      string                              `json:"output_path,omitempty"`
	ConfigurablesV2 []deviceselectionmodel.Configurable `json:"configurables_v2,omitempty"`

	PreferEvent bool `json:"prefer_event,omitempty"`

	//runtime
	Input  interface{} `json:"input,omitempty"`
	Output interface{} `json:"output,omitempty"`

	Retries int64 `json:"retries,omitempty"`
}

type Documentation struct {
	Order int64 `json:"order"`
}

type EventDocumentation struct {
	Order            int64  `json:"order"`
	CharacteristicId string `json:"characteristic_id,omitempty"`
}

type Overwrite struct {
	DeviceId   string `json:"device_id,omitempty"`
	ServiceId  string `json:"service_id,omitempty"`
	ProtocolId string `json:"protocol_id,omitempty"`
}

type NotificationPayload struct {
	//information direct from model
	Message string `json:"message"`
	UserId  string `json:"userId"`
	Title   string `json:"title"`
	IsRead  bool   `json:"isRead"`
	Topic   string `json:"topic"`
}
