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

package interfaces

import (
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
)

type DeploymentParser interface {
	PrepareDeployment(xml string) (deploymentmodel.Deployment, error)
	EstimateStartParameter(xml string) ([]deploymentmodel.ProcessStartParameter, error)
}

type DeploymentStringifier interface {
	Deployment(deployment deploymentmodel.Deployment, userId string, token auth.Token) (string, error)
}
