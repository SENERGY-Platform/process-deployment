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

package bpmn

import (
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
)

type LaneByOrder []deploymentmodel.LaneElement

func (a LaneByOrder) Len() int           { return len(a) }
func (a LaneByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LaneByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

type ElementByOrder []deploymentmodel.Element

func (a ElementByOrder) Len() int           { return len(a) }
func (a ElementByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ElementByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }
