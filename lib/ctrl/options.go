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
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

func (this *Ctrl) SetDeploymentOptions(token jwt_http_router.JwtImpersonate, deployment *model.Deployment) (err error) {
	for _, element := range deployment.Elements {
		if element.Task != nil {
			options, err := this.GetOptions(token, []model.DeviceDescription{element.Task.DeviceDescription})
			if err != nil {
				return err
			}
			element.Task.DeviceOptions = options
		}
		if element.MultiTask != nil {
			options, err := this.GetOptions(token, []model.DeviceDescription{element.MultiTask.DeviceDescription})
			if err != nil {
				return err
			}
			element.MultiTask.DeviceOptions = options
		}
	}
	for _, lane := range deployment.Lanes {
		if lane.Lane != nil {
			options, err := this.GetOptions(token, lane.Lane.DeviceDescriptions)
			if err != nil {
				return err
			}
			lane.Lane.DeviceOptions = options
		}
		if lane.MultiLane != nil {
			options, err := this.GetOptions(token, lane.MultiLane.DeviceDescriptions)
			if err != nil {
				return err
			}
			lane.MultiLane.DeviceOptions = options
		}
	}
	return nil
}

func (this *Ctrl) GetOptions(token jwt_http_router.JwtImpersonate, descriptions []model.DeviceDescription) ([]model.DeviceOption, error) {
	return this.semanticRepo.GetDeploymentOptions(token, descriptions)
}
