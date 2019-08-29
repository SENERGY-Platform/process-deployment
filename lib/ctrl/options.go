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

import "github.com/SENERGY-Platform/process-deployment/lib/model"

func (this *Ctrl) SetDeploymentOptions(deployment *model.Deployment) (err error) {
	for _, element := range deployment.Elements {
		if element.Task != nil {
			options, err := this.GetOptions([]model.DeviceDescription{element.Task.DeviceDescription})
			if err != nil {
				return err
			}
			element.Task.DeviceOptions = options
		}
		if element.MultiTask != nil {
			options, err := this.GetOptions([]model.DeviceDescription{element.MultiTask.DeviceDescription})
			if err != nil {
				return err
			}
			element.MultiTask.DeviceOptions = options
		}
	}
	for _, lane := range deployment.Lanes {
		if lane.Lane != nil {
			options, err := this.GetOptions(lane.Lane.DeviceDescriptions)
			if err != nil {
				return err
			}
			lane.Lane.DeviceOptions = options
		}
		if lane.MultiLane != nil {
			options, err := this.GetOptions(lane.MultiLane.DeviceDescriptions)
			if err != nil {
				return err
			}
			lane.MultiLane.DeviceOptions = options
		}
	}
	return nil
}

func (this *Ctrl) GetOptions(descriptions []model.DeviceDescription) ([]model.DeviceOption, error) {
	return this.semanticRepo.GetDeploymentOptions(descriptions)
}
