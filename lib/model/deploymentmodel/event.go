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

package deploymentmodel

import "github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"

type MsgEvent struct {
	//from model
	Label         string `json:"label"`
	BpmnElementId string `json:"bpmn_element_id"`

	//from user selection
	Device            *devicemodel.Device  `json:"device"`
	Service           *devicemodel.Service `json:"service"`
	Path              string               `json:"path"`
	Value             string               `json:"value"`
	Operation         string               `json:"operation"`
	TriggerConversion *Conversion          `json:"trigger_conversion,omitempty"`

	//generated
	EventId string `json:"event_id"`
}

type TimeEvent struct {
	BpmnElementId string `json:"bpmn_element_id"`
	Kind          string `json:"kind"`
	Time          string `json:"time"`
	Label         string `json:"label"`
}

type Conversion struct {
	From string `json:"from"`
	To   string `json:"to"`
}
