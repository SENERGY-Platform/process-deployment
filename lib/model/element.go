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
	Label             string            `json:"label"`
	DeviceDescription DeviceDescription `json:"device_description"`
	BpmnElementId     string            `json:"bpmn_element_id"`
	Retries           int64             `json:"retries,omitempty"`
	Input             interface{}       `json:"input"`

	//information prepared for the user to select device and service
	Selectables []Selectable `json:"selectables"`

	//information from user to deploy
	Selection Selection `json:"selection"`

	//information to be completed by the user
	Parameter map[string]string `json:"parameter"`
}

type MultiTask struct {
	//information direct from model
	Label             string            `json:"label"`
	DeviceDescription DeviceDescription `json:"device_description"`
	BpmnElementId     string            `json:"bpmn_element_id"`
	Retries           int64             `json:"retries,omitempty"`
	Input             interface{}       `json:"input"`

	//information prepared for the user to select device and service
	Selectables []Selectable `json:"selectables"`

	//information from user to deploy
	Selections []Selection `json:"selections"`

	//information to be completed by the user
	Parameter map[string]string `json:"parameter"`
}
