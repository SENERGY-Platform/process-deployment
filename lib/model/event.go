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

type MsgEvent struct {
	//from model
	Label         string `json:"label" bson:"label"`
	BpmnElementId string `json:"bpmn_element_id" bson:"bpmn_element_id"`
	Order         int    `json:"order" bson:"order"`

	//from user selection
	Device    devicemodel.Device  `json:"device" bson:"device"`
	Service   devicemodel.Service `json:"service" bson:"service"`
	Path      string              `json:"path" bson:"path"`
	Value     string              `json:"value" bson:"value"`
	Operation string              `json:"operation"`

	//generated
	EventId string `json:"event_id" bson:"event_id"`
}

func (this MsgEvent) GetOrder() int {
	return this.Order
}

type TimeEvent struct {
	BpmnElementId string `json:"bpmn_element_id"`
	Kind          string `json:"kind"`
	Time          string `json:"time"`
	Label         string `json:"label"`
	Order         int    `json:"order"`
}

func (this TimeEvent) GetOrder() int {
	return this.Order
}
