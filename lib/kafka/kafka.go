package kafka

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
)

type FactoryType struct{}

var Factory = FactoryType{}

func (FactoryType) NewConsumer(ctx context.Context, config config.Config, topic string, listener func(delivery []byte) error) error {
	return NewConsumer(ctx, config, topic, listener)
}

func (FactoryType) NewProducer(ctx context.Context, config config.Config, topic string) (interfaces.Producer, error) {
	return NewProducer(ctx, config, topic)
}
