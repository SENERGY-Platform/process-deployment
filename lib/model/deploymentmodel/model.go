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

type Deployment struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Diagram    Diagram `json:"diagram"`
	Pools      []Pool  `json:"pools"`
	Executable bool    `json:"executable"`
}

type Diagram struct {
	XmlRaw      string `json:"xml_raw"`
	XmlDeployed string `json:"xml_deployed"`
	Svg         string `json:"svg"`
}

type BaseInfo struct {
	Name   string `json:"name"`
	BpmnId string `json:"bpmn_id"`
	Order  int64  `json:"order"`
}

type Pool struct {
	BaseInfo
	Lanes []Lane `json:"lanes"`
}

type Lane struct {
	BaseInfo
	Elements         []LaneElement    `json:"elements"`
	FilterCriteria   []FilterCriteria `json:"filter_criteria"`
	Selectables      []Selectable     `json:"selectables"`
	SelectedDeviceId string           `json:"selected_device_id"`
}

type LaneElement struct {
	BaseInfo
	TimeEvent    *TimeEvent    `json:"time_event"`
	Notification *Notification `json:"notification"`
	MessageEvent *MessageEvent `json:"message_event"`
	Task         *Task         `json:"task"`
}

type TimeEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type MessageEvent struct {
	DeviceId       string         `json:"device_id"`
	ServiceId      string         `json:"service_id"`
	Value          string         `json:"value"`
	FlowId         string         `json:"flow_id"`
	EventId        string         `json:"event_id"`
	FilterCriteria FilterCriteria `json:"filter_criteria"`
}

type Task struct {
	Retries           int64             `json:"retries"`
	Input             interface{}       `json:"input"`
	Parameter         map[string]string `json:"parameter"`
	Configurables     []Configurable    `json:"configurables"`
	SelectedServiceId string            `json:"selected_service_id"`
}

type Configurable struct {
	CharacteristicId string              `json:"characteristic_id"`
	Values           []ConfigurableValue `json:"values"`
}

type ConfigurableValue struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Value     interface{} `json:"value"`
	ValueType string      `json:"value_type"`
}
