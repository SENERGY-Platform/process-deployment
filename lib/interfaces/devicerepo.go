package interfaces

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
)

type DeviceRepoFactory interface {
	New(ctx context.Context, config config.Config) (DeviceRepo, error)
}

type DeviceRepo interface {
}
