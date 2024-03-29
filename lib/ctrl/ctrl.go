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
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/deployment/parser"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl/deployment/stringifier"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"log"
	"runtime/debug"
)

type Ctrl struct {
	config                config.Config
	db                    interfaces.Database
	devices               interfaces.Devices
	deploymentPublisher   interfaces.Producer
	processrepo           interfaces.ProcessRepo
	deploymentParser      interfaces.DeploymentParser
	deploymentStringifier interfaces.DeploymentStringifier
	imports               interfaces.Imports
}

func New(ctx context.Context, config config.Config, sourcing interfaces.SourcingFactory, db interfaces.Database, devices interfaces.Devices, processrepo interfaces.ProcessRepo, imports interfaces.Imports) (result *Ctrl, err error) {
	result = &Ctrl{
		config:                config,
		db:                    db,
		devices:               devices,
		processrepo:           processrepo,
		imports:               imports,
		deploymentParser:      parser.New(config),
		deploymentStringifier: stringifier.New(config, devices.GetAspectNode),
	}
	result.deploymentPublisher, err = sourcing.NewProducer(ctx, config, config.DeploymentTopic)
	if err != nil {
		return result, err
	}
	err = sourcing.NewConsumer(ctx, config, config.DeploymentTopic, func(delivery []byte) error {
		version := VersionWrapper{}
		err := json.Unmarshal(delivery, &version)
		if err != nil {
			log.Println("ERROR: consumed invalid message --> ignore", err)
			debug.PrintStack()
			return nil
		}
		if version.Version != deploymentmodel.CurrentVersion {
			log.Println("ERROR: consumed unexpected deployment version", version.Version)
			if version.Command == "DELETE" {
				log.Println("handle legacy delete")
				return result.HandleDeployment(messages.DeploymentCommand{
					Command: version.Command,
					Id:      version.Id,
					Version: version.Version,
				})
			}
			return nil
		}
		deployment := messages.DeploymentCommand{}
		err = json.Unmarshal(delivery, &deployment)
		if err != nil {
			log.Println("ERROR: consumed invalid message --> ignore", err)
			debug.PrintStack()
			return err
		}
		return result.HandleDeployment(deployment)
	})
	if err != nil {
		return result, err
	}
	err = sourcing.NewConsumer(ctx, config, config.UsersTopic, func(delivery []byte) error {
		msg := messages.UserCommandMsg{}
		err := json.Unmarshal(delivery, &msg)
		if err != nil {
			debug.PrintStack()
			return err
		}
		return result.HandleUsersCommand(msg)
	})
	if err != nil {
		return result, err
	}
	return result, err
}

type VersionWrapper struct {
	Command string `json:"command"`
	Id      string `json:"id"`
	Version int64  `json:"version"`
}
