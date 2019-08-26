package interfaces

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
)

type SourcingFactory interface {
	NewConsumer(ctx context.Context, config config.Config, topic string, listener func(delivery []byte) error) error
	NewProducer(ctx context.Context, config config.Config, topic string) (Producer, error)
}

type Producer interface {
	Produce(key string, message []byte) error
}
