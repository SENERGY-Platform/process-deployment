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

import (
	dsmodel "github.com/SENERGY-Platform/device-selection/v2/pkg/model"
	dsdevicemodel "github.com/SENERGY-Platform/device-selection/v2/pkg/model/devicemodel"
	"github.com/SENERGY-Platform/models/go/models"
)

//the selectables answer and the request that asks for it are shaped by the service that
//answers, so they are taken from its client package rather than from the shared model. The
//two differ where the shared model has not caught up: a selectable device carries its
//permissions, and an import type carries the declaration of its configs instead of their
//values.

type Selectable = dsmodel.Selectable

type FilterCriteria = dsdevicemodel.FilterCriteria

type FilterCriteriaAndSet = dsmodel.FilterCriteriaAndSet

type BulkRequestElementV2 = dsmodel.BulkRequestElementV2

type BulkRequestV2 = dsmodel.BulkRequestV2

type BulkResult = dsmodel.BulkResult

type BulkResultElement = dsmodel.BulkResultElement

//PathOption and Configurable live in the shared model, which is also where the deployment
//keeps the path the user selected, so these stay aliases of it. dsmodel.PathOption is an
//alias of the same type; naming the shared model here says which side owns the shape.

type PathOption = models.PathOption

type Configurable = models.Configurable
