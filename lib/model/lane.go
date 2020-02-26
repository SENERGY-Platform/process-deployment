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
	Label         string `json:"label"`
	BpmnElementId string `json:"bpmn_element_id"`

	DeviceDescriptions []DeviceDescription `json:"device_descriptions"`

	//information prepared for the user to select device and service
	Selectables []Selectable `json:"selectables"`

	//information from user to deploy
	Selection *devicemodel.Device `json:"selection"`

	Elements []LaneSubElement `json:"elements"`
}

type MultiLane struct {
	//information direct from model
	Label         string `json:"label"`
	BpmnElementId string `json:"bpmn_element_id"`

	DeviceDescriptions []DeviceDescription `json:"device_descriptions"`

	//information prepared for the user to select device and service
	Selectables []Selectable `json:"selectables"`

	//information from user to deploy
	Selections []*devicemodel.Device `json:"selections"`

	Elements []LaneSubElement `json:"elements"`
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
	Label   string `json:"label" bson:"label"`
	Retries int64  `json:"retries,omitempty"`

	DeviceDescription DeviceDescription `json:"device_description"`
	Input             interface{}       `json:"input"`

	BpmnElementId string `json:"bpmn_element_id"`
	MultiTask     bool   `json:"multi_task"`

	SelectedService *devicemodel.Service `json:"selected_service" bson:"selected_service"`

	//information to be completed by the user
	Parameter map[string]string `json:"parameter"`

	Configurables []Configurable `json:"configurables"`
}
