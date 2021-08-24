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

package interfaces

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	deploymentmodel2 "github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
)

type DatabaseFactory interface {
	New(ctx context.Context, config config.Config) (Database, error)
}

type Database interface {
	CheckDeploymentAccess(user string, deploymentId string) (error, int)
	DeleteDeployment(id string) error
	GetDeployment(user string, deploymentId string) (deploymentV1 *deploymentmodel.Deployment, deploymentV2 *deploymentmodel2.Deployment, err error, code int)
	SetDeployment(id string, owner string, deploymentV1 *deploymentmodel.Deployment, deploymentV2 *deploymentmodel2.Deployment) error
	GetDeploymentIds(user string) (deployments []string, err error)
	GetDependencies(user string, deploymentId string) (dependencymodel.Dependencies, error, int)
	GetDependenciesList(user string, limit int, offset int) ([]dependencymodel.Dependencies, error, int)
	GetSelectedDependencies(user string, ids []string) ([]dependencymodel.Dependencies, error, int)
	SetDependencies(dependencies dependencymodel.Dependencies) error
	DeleteDependencies(id string) error
}
