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

package lib

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/api"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	"github.com/SENERGY-Platform/process-deployment/lib/db"
	"github.com/SENERGY-Platform/process-deployment/lib/devices"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/kafka"
)

func StartDefault(ctx context.Context, config config.Config) error {
	return Start(ctx, config, kafka.Factory, db.Factory, devices.Factory)
}

func Start(
	ctx context.Context,
	config config.Config,
	sourcing interfaces.SourcingFactory,
	database interfaces.DatabaseFactory,
	devices interfaces.DevicesFactory) error {

	db, err := database.New(ctx, config)
	if err != nil {
		return err
	}
	d, err := devices.New(ctx, config)
	if err != nil {
		return err
	}
	controller, err := ctrl.New(ctx, config, sourcing, db, d)
	if err != nil {
		return err
	}
	return api.Start(ctx, config, controller)
}
