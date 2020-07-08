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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/coocood/freecache"
	"log"
)

var L1Expiration = 60         // 60sec
var L1Size = 20 * 1024 * 1024 //20MB

type RepositoryFactory struct{}

type Repository struct {
	config       config.Config
	l1           *freecache.Cache
	defaultToken string
}

func (this *RepositoryFactory) New(ctx context.Context, config config.Config) (interfaces.Devices, error) {
	return &Repository{
		config:       config,
		l1:           freecache.NewCache(L1Size),
		defaultToken: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjb25uZWN0aXZpdHktdGVzdCJ9.OnihzQ7zwSq0l1Za991SpdsxkktfrdlNl-vHHpYpXQw",
	}, nil
}

var Factory = &RepositoryFactory{}

func (this *Repository) GetFilteredDevices(token jwt_http_router.JwtImpersonate, descriptions devicemodel.DeviceTypesFilter, protocolBlockList []string) (result []devicemodel.Selectable, err error, code int) {
	filteredProtocols := map[string]bool{}
	for _, protocolId := range protocolBlockList {
		filteredProtocols[protocolId] = true
	}
	deviceTypes, err, code := this.GetFilteredDeviceTypes(token, descriptions)
	if err != nil {
		return result, err, code
	}
	if this.config.Debug {
		log.Println("DEBUG: GetFilteredDevices()::GetFilteredDeviceTypes()", deviceTypes)
	}
	for _, dt := range deviceTypes {
		services := []devicemodel.Service{}
		serviceIndex := map[string]devicemodel.Service{}
		for _, service := range dt.Services {
			for _, desc := range descriptions {
				for _, function := range service.Functions {
					if !(function.RdfType == devicemodel.SES_ONTOLOGY_MEASURING_FUNCTION && filteredProtocols[service.ProtocolId]) {
						if function.Id == desc.FunctionId {
							if desc.AspectId == "" {
								serviceIndex[service.Id] = service
							} else {
								for _, aspect := range service.Aspects {
									if aspect.Id == desc.AspectId {
										serviceIndex[service.Id] = service
									}
								}
							}
						}
					}
				}
			}
		}
		for _, service := range serviceIndex {
			services = append(services, service)
		}
		if len(services) > 0 {
			devices, err, code := this.GetDevicesOfType(token, dt.Id)
			if err != nil {
				return result, err, code
			}
			if this.config.Debug {
				log.Println("DEBUG: GetFilteredDevices()::GetDevicesOfType()", dt.Id, devices)
			}
			for _, device := range devices {
				result = append(result, devicemodel.Selectable{
					Device:   device,
					Services: services,
				})
			}
		}
	}
	if this.config.Debug {
		log.Println("DEBUG: GetFilteredDevices()", result)
	}
	return result, nil, 200
}
