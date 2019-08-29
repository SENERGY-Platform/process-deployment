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

type Task struct {
	//information direct from model
	Label             string            `json:"label" bson:"label"`
	DeviceDescription DeviceDescription `json:"device_description"`
	BpmnElementId     string            `json:"bpmn_element_id" bson:"bpmn_element_id"`
	Input             interface{}       `json:"input"`

	//information prepared for the user to select device and service
	DeviceOptions []DeviceOption `json:"device_options" bson:"device_options"`

	//information from user to deploy
	Selection Selection `json:"selection" bson:"selection"`

	//information to be completed by the user
	Parameter map[string]string
}

type MultiTask struct {
	//information direct from model
	Label             string            `json:"label" bson:"label"`
	DeviceDescription DeviceDescription `json:"device_description"`
	BpmnElementId     string            `json:"bpmn_element_id" bson:"bpmn_element_id"`
	Input             interface{}       `json:"input"`

	//information prepared for the user to select device and service
	DeviceOptions []DeviceOption `json:"device_options" bson:"device_options"`

	//information from user to deploy
	Selections []Selection `json:"selections" bson:"selections"`

	//information to be completed by the user
	Parameter map[string]string
}
