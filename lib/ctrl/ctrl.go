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
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/deployment/parser"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/deployment/stringifier"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"time"
)

type Ctrl struct {
	config                config.Config
	db                    interfaces.Database
	devices               interfaces.Devices
	deploymentPublisher   interfaces.DeploymentProducer
	processrepo           interfaces.ProcessRepo
	deploymentParser      interfaces.DeploymentParser
	deploymentStringifier interfaces.DeploymentStringifier
	imports               interfaces.Imports

	engine          interfaces.Engine
	eventdeployment interfaces.EventDeployment
}

func New(ctx context.Context, config config.Config, sourcing interfaces.SourcingFactory, db interfaces.Database, devices interfaces.Devices, processrepo interfaces.ProcessRepo, imports interfaces.Imports, engine interfaces.Engine, eventdepl interfaces.EventDeployment) (result *Ctrl, err error) {
	result = &Ctrl{
		config:                config,
		db:                    db,
		devices:               devices,
		processrepo:           processrepo,
		imports:               imports,
		deploymentParser:      parser.New(config),
		deploymentStringifier: stringifier.New(config, devices.GetAspectNode),
		engine:                engine,
		eventdeployment:       eventdepl,
	}
	result.deploymentPublisher, err = sourcing.NewDeploymentProducer(ctx, config)
	if err != nil {
		return result, err
	}

	err = sourcing.NewDeviceGroupConsumer(ctx, config, result.HandleDeviceGroupCommand)
	if err != nil {
		return result, err
	}

	err = sourcing.NewUserCommandConsumer(ctx, config, result.HandleUsersCommand)
	if err != nil {
		return result, err
	}

	syncInterval := 10 * time.Minute
	if config.SyncInterval != "" && config.SyncInterval != "-" {
		syncInterval, err = time.ParseDuration(config.SyncInterval)
	}
	syncLockDuration := time.Minute
	if config.SyncLockDuration != "" && config.SyncLockDuration != "-" {
		syncLockDuration, err = time.ParseDuration(config.SyncLockDuration)
	}

	result.StartSyncLoop(ctx, syncInterval, syncLockDuration)

	return result, err
}

type VersionWrapper struct {
	Command string `json:"command"`
	Id      string `json:"id"`
	Version int64  `json:"version"`
}
