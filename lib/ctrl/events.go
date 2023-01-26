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

package ctrl

import (
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	uuid "github.com/satori/go.uuid"
)

func (this *Ctrl) completeEvents(deployment *deploymentmodel.Deployment) error {
	for index, element := range deployment.Elements {
		if element.MessageEvent != nil {
			element.MessageEvent.EventId = uuid.NewV4().String()
			deployment.Elements[index] = element
		}
		if element.ConditionalEvent != nil {
			element.ConditionalEvent.EventId = uuid.NewV4().String()
			deployment.Elements[index] = element
		}
	}
	return nil
}
