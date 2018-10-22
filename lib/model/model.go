/*
 * Copyright 2018 InfAI (CC SES)
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

import "github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

type BpmnAbstractMsg struct {
	TaskId     string            `json:"-"`
	Parameter  []TaskParameter   `json:"-"`
	ReuseId    string            `json:"-"`
	DeviceType string            `json:"device_type"`
	Service    string            `json:"service"`
	Label      string            `json:"label"`
	Values     BpmnValueSkeleton `json:"values"`
}

type BpmnValueSkeleton struct {
	Inputs  map[string]interface{} `json:"inputs,omitempty"`
	Outputs map[string]interface{} `json:"outputs,omitempty"`
}

type BpmnMsg struct {
	InstanceId string                 `json:"instance_id,omitempty"`
	ServiceId  string                 `json:"service_id,omitempty"`
	Inputs     map[string]interface{} `json:"inputs,omitempty"`
	Outputs    map[string]interface{} `json:"outputs,omitempty"`
	ErrorMsg   string                 `json:"error_msg,omitempty"`
}

type InputOutput struct {
	Name    string        `json:"name,omitempty"`
	FieldId string        `json:"field_id"`
	Type    Type          `json:"type"`
	Value   string        `json:"value,omitempty"`
	Values  []InputOutput `json:"values,omitempty"`
}

type Type struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Base string `json:"base"`
}

type AbstractTask struct {
	Tasks        []Task                 `json:"tasks"`
	DeviceTypeId string                 `json:"device_type_id"`
	Options      []model.DeviceInstance `json:"options"`
	Selected     model.DeviceInstance   `json:"selected"`
	State        string                 `json:"state" bson:"-"`
}

type TaskParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Task struct {
	Id        string            `json:"id"`
	Values    BpmnValueSkeleton `json:"values"`
	Parameter []TaskParameter   `json:"parameter"`
	ServiceId string            `json:"service_id"`
	Label     string            `json:"label"`
}

type AbstractDataExportTask struct {
	Id        string            `json:"id"`
	MeasurementField			string		`json:"measurementField"`
	MeasurementIdSelected		string		`json:"measurement"`
	AnalysisAction string `json:"analysisAction"`
	Interval Interval `json:"interval"`
	UserId string `json:"userid"`
}

type Interval struct {
	Value string `json:"value"`
	TimeType string `json:"timeType"`
}

type MsgEvent struct {
	FilterId string `json:"filter_id,omitempty"`
	ShapeId  string `json:"shape_id"`
	Filter   Filter `json:"filter"`
	State    string `json:"state,omitempty" bson:"-"`
}

type TimeEvent struct {
	ShapeId string `json:"shape_id"`
	Kind    string `json:"kind"`
	Time    string `json:"time"`
	Label   string `json:"label"`
}

type AbstractProcess struct {
	Xml           string         `json:"xml"`
	Name          string         `json:"name"`
	AbstractTasks []AbstractTask `json:"abstract_tasks"`
	AbstractDataExportTasks []AbstractDataExportTask `json:"abstract_data_export_tasks"`
	ReceiveTasks  []MsgEvent     `json:"receive_tasks"`
	MsgEvents     []MsgEvent     `json:"msg_events"`
	TimeEvents    []TimeEvent    `json:"time_events"`
}

type Rule struct {
	Path     string `json:"path" bson:"path,omitempty"`         // github.com/NodePrime/jsonpath
	Scope    string `json:"scope" bson:"scope,omitempty"`       // 'any' || 'all' || 'none' || 'max <number>' || 'min <number>' || '<number>'
	Operator string `json:"operator" bson:"operator,omitempty"` // '==' || '!=' || '>=' || '<=' || '<' || '>' || 'regex'
	Value    string `json:"value" bson:"value,omitempty"`
}

type Filter struct {
	Scope     string `json:"scope" bson:"scope,omitempty"` // 'any' || 'all' || 'none' || 'max <number>' || 'min <number>' || '<number>'
	Rules     []Rule `json:"rules" bson:"rules,omitempty"`
	Topic     string `json:"topic" bson:"topic,omitempty"`
	DeviceId  string `json:"device_id,omitempty" bson:"device_id,omitempty"`
	ServiceId string `json:"service_id,omitempty" bson:"service_id,omitempty"`
}

type DeploymentRequest struct {
	Process AbstractProcess `json:"process"`
	Svg     string          `json:"svg"`
}
