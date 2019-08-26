package interfaces

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

type ConnectionlogFactory interface {
	New(ctx context.Context, config config.Config) (Connectionlog, error)
}

type Connectionlog interface {
	CheckDeviceStates(jwtimpersonate jwt_http_router.JwtImpersonate, ids []string) (result map[string]bool, err error)
}
