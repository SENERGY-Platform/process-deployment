package lib

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
)

func New(ctx context.Context, config config.Config, sourcing interfaces.SourcingFactory, database interfaces.DatabaseFactory, connectionlog interfaces.ConnectionlogFactory, devicerepo interfaces.DeviceRepoFactory) error {

}
