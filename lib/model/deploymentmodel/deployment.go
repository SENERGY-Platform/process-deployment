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
	Id          string        `json:"id"`
	Executable  bool          `json:"executable"`
	XmlRaw      string        `json:"xml_raw"`
	Xml         string        `json:"xml"`
	Svg         string        `json:"svg"`
	Name        string        `json:"name"`
	Elements    []Element     `json:"elements"`
	Lanes       []LaneElement `json:"lanes"`
	Description string        `json:"description,omitempty"`
}

type Element struct {
	Order            int64      `json:"order"`
	Task             *Task      `json:"task,omitempty"`
	MultiTask        *MultiTask `json:"multi_task,omitempty"`
	ReceiveTaskEvent *MsgEvent  `json:"receive_task_event,omitempty"`
	MsgEvent         *MsgEvent  `json:"msg_event,omitempty"`
	TimeEvent        *TimeEvent `json:"time_event,omitempty"`
}

type LaneElement struct {
	Order     int64      `json:"order"`
	MultiLane *MultiLane `json:"multi_lane,omitempty"`
	Lane      *Lane      `json:"lane,omitempty"`
}
