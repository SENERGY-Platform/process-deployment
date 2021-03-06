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

package processmodel

type ProcessModel struct {
	Id          string `json:"_id" bson:"_id"`
	Date        int64  `json:"date" bson:"date"`
	Owner       string `json:"owner" bson:"owner"`
	BpmnXml     string `json:"bpmn_xml" bson:"bpmn_xml"`
	SvgXml      string `json:"svgXML" bson:"svgXML"`
	Publish     bool   `json:"publish" bson:"publish"`
	PublishDate string `json:"publish_date" bson:"publish_date"`
	Description string `json:"description" bson:"description"`
}
