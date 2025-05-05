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
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"time"
)

type DatabaseFactory interface {
	New(ctx context.Context, config config.Config) (Database, error)
}

type Database interface {
	CheckDeploymentAccess(user string, deploymentId string) (error, int)
	ListDeployments(user string, options model.DeploymentListOptions) (deployments []deploymentmodel.Deployment, err error)
	GetDeployment(user string, deploymentId string) (deployment *deploymentmodel.Deployment, err error, code int)
	GetDeploymentIds(user string) (deployments []string, err error)
	GetDependencies(user string, deploymentId string) (dependencymodel.Dependencies, error, int)
	GetDependenciesList(user string, limit int, offset int) ([]dependencymodel.Dependencies, error, int)
	GetSelectedDependencies(user string, ids []string) ([]dependencymodel.Dependencies, error, int)
	SetDependencies(dependencies dependencymodel.Dependencies) error
	DeleteDependencies(id string) error

	SetDeployment(depl messages.DeploymentCommand, syncHandler func(messages.DeploymentCommand) error) error
	DeleteDeployment(id string, syncDeleteHandler func(messages.DeploymentCommand) error) error
	RetryDeploymentSync(lockduration time.Duration, syncDeleteHandler func(messages.DeploymentCommand) error, syncHandler func(messages.DeploymentCommand) error) error
}
