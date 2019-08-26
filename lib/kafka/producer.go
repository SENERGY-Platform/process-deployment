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
	broker, err := GetBroker(config.ZookeeperUrl)
	if err != nil {
		log.Println("ERROR: unable to get broker list", err)
		return nil, err
	}
	err = InitTopicWithConfig(config.ZookeeperUrl, 1, 1, topic)
	if err != nil {
		log.Println("ERROR: unable to create topic", err)
		return nil, err
	}
	var logger *log.Logger = nil

	if config.Debug {
		logger = log.New(os.Stdout, "KAFKA", 0)
	}

	result.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:     broker,
		Topic:       topic,
		Async:       false,
		Logger:      logger,
		ErrorLogger: log.New(os.Stderr, "KAFKA", 0),
	})
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
