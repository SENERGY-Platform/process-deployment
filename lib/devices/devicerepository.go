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
	"errors"
	devicerepo "github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/service-commons/pkg/cache"
	"net/http"
)

func (this *Repository) GetAspectNode(token auth.Token, id string) (devicemodel.AspectNode, error) {
	resource := "aspect-nodes"
	return cache.Use(this.cache, resource+"."+id, func() (result devicemodel.AspectNode, err error) {
		result, err, _ = this.devicerepo.GetAspectNode(id)
		return result, err
	}, func(node devicemodel.AspectNode) error {
		if node.Id == "" {
			return errors.New("invalid aspect-node returned from cache")
		}
		return nil
	}, CacheExpiration)
}

func (this *Repository) GetDevice(token auth.Token, id string) (result devicemodel.Device, err error, code int) {
	resource := "devices"
	code = http.StatusOK
	result, err = cache.Use(this.cache, resource+"."+id, func() (temp devicemodel.Device, err error) {
		temp, err, code = this.devicerepo.ReadDevice(id, token.Jwt(), devicerepo.READ)
		return temp, err
	}, func(device devicemodel.Device) error {
		if device.Id == "" {
			return errors.New("invalid device returned from cache")
		}
		return nil
	}, CacheExpiration)
	return result, err, code
}

func (this *Repository) GetService(token auth.Token, id string) (result devicemodel.Service, err error, code int) {
	resource := "services"
	code = http.StatusOK
	result, err = cache.Use(this.cache, resource+"."+id, func() (temp devicemodel.Service, err error) {
		temp, err, code = this.devicerepo.GetService(id)
		return temp, err
	}, func(service devicemodel.Service) error {
		if service.Id == "" {
			return errors.New("invalid service returned from cache")
		}
		return nil
	}, CacheExpiration)
	return result, err, code
}

func (this *Repository) GetDeviceGroup(token auth.Token, id string) (result devicemodel.DeviceGroup, err error, code int) {
	resource := "device-groups"
	code = http.StatusOK
	result, err = cache.Use(this.cache, resource+"."+id, func() (temp devicemodel.DeviceGroup, err error) {
		temp, err, code = this.devicerepo.ReadDeviceGroup(id, token.Jwt(), false)
		return temp, err
	}, func(group devicemodel.DeviceGroup) error {
		if group.Id == "" {
			return errors.New("invalid group returned from cache")
		}
		return nil
	}, CacheExpiration)
	return result, err, code
}
