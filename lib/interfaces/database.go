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
)

type DatabaseFactory interface {
	New(ctx context.Context, config config.Config) (Database, error)
}

type Database interface {
	CheckDeploymentAccess(user string, deploymentId string) (error, int)
	GetDeployment(user string, deploymentId string) (model.Deployment, error, int)
	DeleteDeployment(id string) error
	SetDeployment(id string, owner string, deployment model.Deployment) error
	GetDependencies(user string, deploymentId string) (model.Dependencies, error, int)
	GetDependenciesList(user string, limit int, offset int) ([]model.Dependencies, error, int)
	GetSelectedDependencies(user string, ids []string) ([]model.Dependencies, error, int)
	SetDependencies(dependencies model.Dependencies) error
	DeleteDependencies(id string) error
}
