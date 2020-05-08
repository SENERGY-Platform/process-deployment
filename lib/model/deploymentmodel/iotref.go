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

type Selectable struct {
	Device   Device    `json:"device"`
	Services []Service `json:"services"`
}

type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	FunctionId string `json:"function_id"`
}

type FilterCriteria struct {
	CharacteristicId *string `json:"characteristic_id"`
	FunctionId       *string `json:"function_id"`
	DeviceClassId    *string `json:"device_class_id"`
	AspectId         *string `json:"aspect_id"`
}
