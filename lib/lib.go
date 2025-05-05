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
	engine "github.com/SENERGY-Platform/camunda-engine-wrapper/lib/client"
	eventdeployment "github.com/SENERGY-Platform/event-deployment/lib/client"
	"github.com/SENERGY-Platform/process-deployment/lib/api"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	"github.com/SENERGY-Platform/process-deployment/lib/db"
	"github.com/SENERGY-Platform/process-deployment/lib/devices"
	"github.com/SENERGY-Platform/process-deployment/lib/events"
	"github.com/SENERGY-Platform/process-deployment/lib/imports"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/processrepo"
	"github.com/SENERGY-Platform/service-commons/pkg/cache/invalidator"
	kafkalib "github.com/SENERGY-Platform/service-commons/pkg/kafka"
	"log"
	"runtime/debug"
	"time"
)

func StartDefault(ctx context.Context, config config.Config) error {
	return Start(ctx, config, events.Factory, db.Factory, devices.Factory, processrepo.Factory, imports.Factory)
}

func Start(
	ctx context.Context,
	config config.Config,
	sourcing interfaces.SourcingFactory,
	database interfaces.DatabaseFactory,
	devices interfaces.DevicesFactory,
	processrepo interfaces.ProcessRepoFactory,
	imports interfaces.ImportsFactory) error {

	db, err := database.New(ctx, config)
	if err != nil {
		return err
	}
	d, err := devices.New(ctx, config)
	if err != nil {
		return err
	}
	p, err := processrepo.New(ctx, config)
	if err != nil {
		return err
	}
	i, err := imports.New(config)
	if err != nil {
		return err
	}

	var e interfaces.Engine = VoidEngine{}
	var eventdepl interfaces.EventDeployment = VoidEventDeployment{}
	if config.ProcessEngineWrapperUrl != "" && config.ProcessEngineWrapperUrl != "-" {
		e = engine.New(config.ProcessEngineWrapperUrl)
	}
	if config.EventDeploymentUrl != "" && config.EventDeploymentUrl != "-" {
		eventdepl = eventdeployment.New(config.EventDeploymentUrl)
	}

	controller, err := ctrl.New(ctx, config, sourcing, db, d, p, i, e, eventdepl)
	if err != nil {
		return err
	}
	return api.Start(ctx, config, controller)
}

func StartCacheInvalidator(ctx context.Context, conf config.Config) error {
	if conf.KafkaUrl == "" || conf.KafkaUrl == "-" {
		return nil
	}
	return invalidator.StartCacheInvalidatorAll(ctx, kafkalib.Config{
		KafkaUrl:               conf.KafkaUrl,
		StartOffset:            kafkalib.LastOffset,
		Debug:                  conf.Debug,
		PartitionWatchInterval: time.Minute,
		OnError: func(err error) {
			log.Println("ERROR:", err)
			debug.PrintStack()
		},
	}, conf.CacheInvalidationKafkaTopics, nil)
}

type VoidEngine struct{}

func (this VoidEngine) Deploy(token string, depl engine.DeploymentMessage) (err error, code int) {
	return nil, 200
}

func (this VoidEngine) DeleteDeployment(token string, userId string, deplId string) (err error, code int) {
	return nil, 200
}

type VoidEventDeployment struct{}

func (this VoidEventDeployment) Deploy(token string, depl eventdeployment.Deployment) (err error, code int) {
	return nil, 200
}

func (this VoidEventDeployment) DeleteDeployment(token string, userId string, deplId string) (err error, code int) {
	return nil, 200
}

func (this VoidEventDeployment) UpdateDeploymentsOfDeviceGroup(token string, dg eventdeployment.DeviceGroup) (err error, code int) {
	return nil, 200
}
