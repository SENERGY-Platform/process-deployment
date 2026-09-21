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

package devices

import (
	"github.com/SENERGY-Platform/device-selection/v2/pkg/client"

	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"
)

// GetBulkDeviceSelectionV2 asks the selection service which devices, groups and imports can
// serve the criteria of every element of a deployment, in one request. complete_services is
// set because the deployment options need the full import type and its path options; for
// devices the flag does nothing, its name is a legacy artefact of the endpoint.
func (this *Repository) GetBulkDeviceSelectionV2(token auth.Token, bulk deviceselectionmodel.BulkRequestV2) (result deviceselectionmodel.BulkResult, err error, code int) {
	if this.config.Debug {
		this.config.GetLogger().Debug("send GetBulkDeviceSelectionV2()", "elements", len(bulk))
	}
	result, code, err = this.deviceselection.GetBulkSelectablesV2(token.Token, bulk, &client.GetBulkSelectablesOptions{CompleteServices: true})
	return result, err, code
}
