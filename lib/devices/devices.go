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

package devices

import (
	"context"
	devicerepo "github.com/SENERGY-Platform/device-repository/lib/client"
	permv2 "github.com/SENERGY-Platform/permissions-v2/pkg/client"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/service-commons/pkg/cache"
	"github.com/SENERGY-Platform/service-commons/pkg/signal"
	"time"
)

var CacheExpiration = 60 * time.Second

type RepositoryFactory struct{}

func (this *RepositoryFactory) New(ctx context.Context, config config.Config) (interfaces.Devices, error) {
	c, err := cache.New(cache.Config{
		CacheInvalidationSignalHooks: map[cache.Signal]cache.ToKey{
			signal.Known.CacheInvalidationAll:        nil,
			signal.Known.AspectCacheInvalidation:     nil,
			signal.Known.DeviceTypeCacheInvalidation: nil,
			signal.Known.DeviceGroupInvalidation: func(signalValue string) (cacheKey string) {
				return "device-groups." + signalValue
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &Repository{
		config:       config,
		cache:        c,
		defaultToken: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjb25uZWN0aXZpdHktdGVzdCJ9.OnihzQ7zwSq0l1Za991SpdsxkktfrdlNl-vHHpYpXQw",
		devicerepo:   devicerepo.NewClient(config.DeviceRepoUrl, nil),
		permv2:       permv2.New(config.PermissionsV2Url),
	}, nil
}

var Factory = &RepositoryFactory{}

type Repository struct {
	config       config.Config
	cache        *cache.Cache
	defaultToken string
	devicerepo   devicerepo.Interface
	permv2       permv2.Client
}
