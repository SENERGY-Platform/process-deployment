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

package mocks

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
)

type DatabaseMock struct {
}

var Database = &DatabaseMock{}

func (this *DatabaseMock) New(ctx context.Context, config config.Config) (interfaces.Database, error) {
	return this, nil
}

func (this *DatabaseMock) CheckDeploymentAccess(user string, deploymentId string) (error, int) {
	return nil, 200
}

func (this *DatabaseMock) GetDeployment(user string, deploymentId string) (model.Deployment, error, int) {
	panic("implement me")
}
