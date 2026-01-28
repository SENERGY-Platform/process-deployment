/*
 * Copyright 2025 InfAI (CC SES)
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

package events

import (
	"context"
	"encoding/json"
	"runtime/debug"

	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/events/kafka"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
)

type Event struct{}

var Factory = &Event{}

func (this *Event) NewUserCommandConsumer(ctx context.Context, config config.Config, listener func(delivery messages.UserCommandMsg) error) error {
	if config.UsersTopic == "" || config.UsersTopic == "-" {
		return nil
	}
	return kafka.NewConsumer(ctx, config, config.UsersTopic, func(delivery []byte) error {
		msg := messages.UserCommandMsg{}
		err := json.Unmarshal(delivery, &msg)
		if err != nil {
			debug.PrintStack()
			return err
		}
		return listener(msg)
	})
}

func (this *Event) NewDeviceGroupConsumer(ctx context.Context, config config.Config, listener func(groupId string) error) error {
	if config.DeviceGroupTopic == "" || config.DeviceGroupTopic == "-" {
		return nil
	}
	return kafka.NewConsumer(ctx, config, config.DeviceGroupTopic, func(delivery []byte) error {
		msg := messages.DeviceGroupCommand{}
		err := json.Unmarshal(delivery, &msg)
		if err != nil {
			debug.PrintStack()
			return err
		}
		return listener(msg.Id)
	})
}

func (this *Event) NewDeploymentProducer(ctx context.Context, config config.Config) (interfaces.DeploymentProducer, error) {
	result := &DeploymentProducer{
		Config: config,
	}
	var err error
	if config.DeploymentTopic != "" && config.DeploymentTopic != "-" {
		result.DeploymentProducer, err = kafka.NewProducer(ctx, config, config.DeploymentTopic)
		if err != nil {
			return nil, err
		}
	}

	if config.DoneTopic != "" && config.DoneTopic != "-" {
		result.DoneProducer, err = kafka.NewProducer(ctx, config, config.DoneTopic)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type DeploymentProducer struct {
	Config             config.Config
	DeploymentProducer interfaces.Producer
	DoneProducer       interfaces.Producer
}

type DoneNotification struct {
	Command string `json:"command"`
	Id      string `json:"id"`
	Handler string `json:"handler"`
}

func (this *DeploymentProducer) Produce(command messages.DeploymentCommand) error {
	if this.DeploymentProducer != nil {
		msg, err := json.Marshal(command)
		if err != nil {
			return err
		}
		err = this.DeploymentProducer.Produce(command.Id, msg)
		if err != nil {
			return err
		}
	}

	if this.DoneProducer != nil {
		message := DoneNotification{
			Command: command.Command,
			Id:      command.Id,
			Handler: "github.com/SENERGY-Platform/process-deployment",
		}
		msg, err := json.Marshal(message)
		if err != nil {
			return err
		}
		this.Config.GetLogger().Debug("send done notification", "message", string(msg))
		err = this.DoneProducer.Produce(command.Id, msg)
		if err != nil {
			return err
		}
	}
	return nil
}
