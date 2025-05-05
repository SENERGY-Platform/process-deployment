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
	"encoding/json"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/events"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"log"
	"runtime/debug"
	"sync"
)

var Kafka = &KafkaMock{Produced: map[string][]string{}, Listeners: map[string][]func(msg []byte) error{}}

type KafkaMock struct {
	mux       sync.Mutex
	Produced  map[string][]string
	Listeners map[string][]func(msg []byte) error
}

type Producer struct {
	Topic string
	Kafka *KafkaMock
}

func (this *KafkaMock) NewUserCommandConsumer(ctx context.Context, config config.Config, listener func(delivery messages.UserCommandMsg) error) error {
	return this.NewConsumer(ctx, config, config.UsersTopic, func(delivery []byte) error {
		msg := messages.UserCommandMsg{}
		err := json.Unmarshal(delivery, &msg)
		if err != nil {
			debug.PrintStack()
			return err
		}
		return listener(msg)
	})
}

func (this *KafkaMock) NewDeviceGroupConsumer(ctx context.Context, config config.Config, listener func(group models.DeviceGroup) error) error {
	return this.NewConsumer(ctx, config, config.DeviceGroupTopic, func(delivery []byte) error {
		msg := messages.DeviceGroupCommand{}
		err := json.Unmarshal(delivery, &msg)
		if err != nil {
			debug.PrintStack()
			return err
		}
		return listener(msg.DeviceGroup)
	})
}

func (this *KafkaMock) NewDeploymentProducer(ctx context.Context, config config.Config) (interfaces.DeploymentProducer, error) {
	deploymentProducer, err := this.NewProducer(ctx, config, config.DeploymentTopic)
	if err != nil {
		return nil, err
	}
	doneProducer, err := this.NewProducer(ctx, config, config.DoneTopic)
	if err != nil {
		return nil, err
	}
	return &events.DeploymentProducer{
		Config:             config,
		DoneProducer:       doneProducer,
		DeploymentProducer: deploymentProducer,
	}, nil
}

func (this *KafkaMock) NewConsumer(ctx context.Context, config config.Config, topic string, listener func(delivery []byte) error) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.Listeners[topic] = append(this.Listeners[topic], listener)
	go func() {
		<-ctx.Done()
		this.mux.Lock()
		defer this.mux.Unlock()
		this.Listeners[topic] = []func(msg []byte) error{}
	}()
	return nil
}

func (this *KafkaMock) NewProducer(ctx context.Context, config config.Config, topic string) (interfaces.Producer, error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	go func() {
		<-ctx.Done()
		this.mux.Lock()
		defer this.mux.Unlock()
		this.Produced[topic] = []string{}
	}()
	return &Producer{Kafka: this, Topic: topic}, nil
}

func (this *Producer) Produce(key string, message []byte) error {
	this.Kafka.mux.Lock()
	defer this.Kafka.mux.Unlock()
	this.Kafka.Produced[this.Topic] = append(this.Kafka.Produced[this.Topic], string(message))
	for _, l := range this.Kafka.Listeners[this.Topic] {
		go func(listener func(msg []byte) error) {
			err := listener(message)
			if err != nil {
				log.Println(err)
				debug.PrintStack()
			}
		}(l)
	}
	return nil
}

func (this *KafkaMock) GetProduced(topic string) []string {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.Produced[topic]
}
