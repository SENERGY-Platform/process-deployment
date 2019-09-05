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
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

type DeviceManagerFactory interface {
	New(ctx context.Context, config config.Config) (DeviceManager, error)
}

type DeviceManager interface {
	GetDeploymentOptions(token jwt_http_router.JwtImpersonate, descriptions []model.DeviceDescription) ([]model.DeviceOption, error)
}
