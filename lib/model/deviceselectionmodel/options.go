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

package deviceselectionmodel

import "github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"

type Selectable struct {
	Device      *devicemodel.Device      `json:"device"`
	Services    []devicemodel.Service    `json:"services"`
	DeviceGroup *devicemodel.DeviceGroup `json:"device_group,omitempty"`
}

type FilterCriteriaAndSet []FilterCriteria

type FilterCriteriaOrSet []FilterCriteria

type FilterCriteria struct {
	FunctionId    string `json:"function_id"`
	DeviceClassId string `json:"device_class_id"`
	AspectId      string `json:"aspect_id"`
}

type BulkRequestElement struct {
	Id                string                   `json:"id"`
	FilterInteraction *devicemodel.Interaction `json:"filter_interaction"`
	FilterProtocols   []string                 `json:"filter_protocols"`
	Criteria          FilterCriteriaAndSet     `json:"criteria"`
	IncludeGroups     bool                     `json:"include_groups"`
}

type BulkRequest []BulkRequestElement

type BulkResult []BulkResultElement

type BulkResultElement struct {
	Id          string       `json:"id"`
	Selectables []Selectable `json:"selectables"`
}
