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

import "github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"

type Lane struct {
	//information direct from model
	Label         string `json:"label" bson:"label"`
	BpmnElementId string `json:"bpmn_element_id" bson:"bpmn_element_id"`

	CharacteristicIds []string                  `json:"characteristic_ids" bson:"characteristic_ids"`
	Functions         []devicemodel.Function    `json:"function_id" bson:"function_id"`
	DeviceClasses     []devicemodel.DeviceClass `json:"device_class,omitempty" bson:"device_class,omitempty"`
	Aspects           []devicemodel.Aspect      `json:"aspect,omitempty" bson:"aspect,omitempty"`

	//information prepared for the user to select device and service
	DeviceOptions []devicemodel.Device `json:"device_options" bson:"device_options"`

	//information from user to deploy
	Selection devicemodel.Device `json:"selection" bson:"selection"`

	Elements []LaneSubElement `json:"elements" bson:"elements"`
}

type MultiLane struct {
	//information direct from model
	Label         string `json:"label" bson:"label"`
	BpmnElementId string `json:"bpmn_element_id" bson:"bpmn_element_id"`

	CharacteristicIds []string                  `json:"characteristic_ids" bson:"characteristic_ids"`
	Functions         []devicemodel.Function    `json:"function_id" bson:"function_id"`
	DeviceClasses     []devicemodel.DeviceClass `json:"device_class,omitempty" bson:"device_class,omitempty"`
	Aspects           []devicemodel.Aspect      `json:"aspect,omitempty" bson:"aspect,omitempty"`

	//information prepared for the user to select device and service
	DeviceOptions []devicemodel.Device `json:"device_options" bson:"device_options"`

	//information from user to deploy
	Selections []devicemodel.Device `json:"selections" bson:"selections"`

	Elements []LaneSubElement `json:"elements" bson:"elements"`
}

type LaneSubElement struct {
	Order            int64      `json:"order"`
	LaneTask         *LaneTask  `json:"task,omitempty"`
	MsgEvent         *MsgEvent  `json:"msg_event,omitempty"`
	ReceiveTaskEvent *MsgEvent  `json:"receive_task_event,omitempty"`
	TimeEvent        *TimeEvent `json:"time_event,omitempty"`
}

type LaneTask struct {
	//information direct from model
	Label string `json:"label" bson:"label"`

	CharacteristicId string                   `json:"characteristic_id" bson:"characteristic_id"`
	Function         devicemodel.Function     `json:"function" bson:"function"`
	DeviceClass      *devicemodel.DeviceClass `json:"device_class,omitempty" bson:"device_class,omitempty"`
	Aspect           *devicemodel.Aspect      `json:"aspect,omitempty" bson:"aspect,omitempty"`
	Input            interface{}              `json:"input"`

	BpmnElementId string `json:"bpmn_element_id" bson:"bpmn_element_id"`
	MultiTask     bool   `json:"multi_task" bson:"multi_task"`

	//prepared
	ServiceOptions  []devicemodel.Service `json:"service_options" bson:"-"`
	SelectedService devicemodel.Service   `json:"selected_service" bson:"selected_service"`

	//information to be completed by the user
	Parameter map[string]string
}
