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
	"errors"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/segmentio/kafka-go"
)

func NewConsumer(ctx context.Context, config config.Config, topic string, listener func(delivery []byte) error) error {
	if config.InitTopics {
		err := InitTopic(config.KafkaUrl, topic)
		if err != nil {
			config.GetLogger().Error("unable to create topic", "topic", topic, "error", err)
			return err
		}
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		CommitInterval:         0, //synchronous commits
		Brokers:                []string{config.KafkaUrl},
		GroupID:                config.ConsumerGroup,
		Topic:                  topic,
		MaxWait:                1 * time.Second,
		Logger:                 log.New(io.Discard, "", 0),
		ErrorLogger:            log.New(os.Stdout, "[KAFKA-ERROR] ", 0),
		WatchPartitionChanges:  true,
		PartitionWatchInterval: time.Minute,
	})
	go func() {
		defer r.Close()
		defer config.GetLogger().Info("close consumer for topic", "topic", topic)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m, err := r.FetchMessage(ctx)
				if err == io.EOF || errors.Is(err, context.Canceled) {
					return
				}
				if err != nil {
					config.GetLogger().Error("FATAL: while fetching message", "topic", topic, "error", err)
					log.Fatal("ERROR: while consuming topic ", topic, err)
					return
				}

				err = retry(func() error {
					return listener(m.Value)
				}, func(n int64) time.Duration {
					return time.Duration(n) * time.Second
				}, 10*time.Minute)

				if err != nil {
					config.GetLogger().Error("FATAL: unable to handle message", "topic", topic, "error", err)
					log.Fatal("ERROR: unable to handle message (no commit)", err)
				} else {
					err = r.CommitMessages(ctx, m)
					if err != nil {
						config.GetLogger().Error("FATAL: while committing consumption", "topic", topic, "error", err)
						log.Fatal("ERROR: while committing consumption ", topic, err)
						return
					}
				}
			}
		}
	}()
	return nil
}

func retry(f func() error, waitProvider func(n int64) time.Duration, timeout time.Duration) (err error) {
	err = errors.New("initial")
	start := time.Now()
	for i := int64(1); err != nil && time.Since(start) < timeout; i++ {
		err = f()
		if err != nil {
			slog.Default().Error("kafka listener error", "error", err)
			wait := waitProvider(i)
			if time.Since(start)+wait < timeout {
				slog.Default().Error("retry after", "wait", wait, "error", err)
				time.Sleep(wait)
			} else {
				return err
			}
		}
	}
	return err
}
