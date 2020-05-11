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

package ctrl

import (
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

func (this *Ctrl) SetDeploymentOptions(token jwt_http_router.JwtImpersonate, deployment *deploymentmodel.Deployment) (err error) {
	for i, pool := range deployment.Pools {
		for j, lane := range pool.Lanes {
			deployment.Pools[i].Lanes[j].Selectables, err = this.GetOptions(token, lane.GetFilterCriteria())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Ctrl) GetOptions(token jwt_http_router.JwtImpersonate, criteria []deploymentmodel.FilterCriteria) (result []deploymentmodel.Selectable, err error) {
	if len(criteria) == 0 {
		return []deploymentmodel.Selectable{}, nil
	}
	result, err, _ = this.devices.GetFilteredDevices(token, criteria)
	return
}
