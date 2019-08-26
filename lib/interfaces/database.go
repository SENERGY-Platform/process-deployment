package interfaces

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
)

type DatabaseFactory interface {
	New(ctx context.Context, config config.Config) (Database, error)
}

type Database interface {
}
