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

package kafka

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

type Producer struct {
	writer *kafka.Writer
	ctx    context.Context
}

func NewProducer(ctx context.Context, config config.Config, topic string) (interfaces.Producer, error) {
	result := &Producer{ctx: ctx}
	if config.InitTopics {
		err := InitTopic(config.KafkaUrl, topic)
		if err != nil {
			log.Println("ERROR: unable to create topic", err)
			return nil, err
		}
	}

	var logger kafka.Logger
	if config.Debug {
		logger = log.New(os.Stdout, "KAFKA", 0)
	}

	result.writer = &kafka.Writer{
		Addr:        kafka.TCP(config.KafkaUrl),
		Topic:       topic,
		Logger:      logger,
		Async:       false,
		BatchSize:   1,
		Balancer:    &kafka.Hash{},
		ErrorLogger: log.New(os.Stderr, "KAFKA", 0),
	}

	go func() {
		<-ctx.Done()
		result.writer.Close()
	}()
	return result, nil
}

func (this *Producer) Produce(key string, message []byte) error {
	return this.writer.WriteMessages(this.ctx, kafka.Message{
		Key:   []byte(key),
		Value: message,
		Time:  config.TimeNow(),
	})
}
