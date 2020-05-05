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

package dependencymodel

type Dependencies struct {
	//from db
	DeploymentId string             `json:"deployment_id" bson:"deployment_id"`
	Owner        string             `json:"owner" bson:"owner"`
	Devices      []DeviceDependency `json:"devices" bson:"devices"`
	Events       []EventDependency  `json:"events" bson:"events"`
}

type DeviceDependency struct {
	DeviceId      string         `json:"device_id" bson:"device_id"`
	Name          string         `json:"name" bson:"name"`
	BpmnResources []BpmnResource `json:"bpmn_resources" bson:"bpmn_resources"`
}

type EventDependency struct {
	EventId       string         `json:"event_id" bson:"event_id"`
	BpmnResources []BpmnResource `json:"bpmn_resources" bson:"bpmn_resources"`
}

type BpmnResource struct {
	Id    string `json:"id" bson:"id"`
	label string `json:"label" bson:"label"`
}
