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

package mocks

import (
	"context"
	devicerepo "github.com/SENERGY-Platform/device-repository/lib/client"
	"github.com/SENERGY-Platform/device-repository/lib/database"
	"github.com/SENERGY-Platform/device-repository/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"

	"sync"
)

type DeviceRepoMock struct {
	mux     sync.Mutex
	options []deviceselectionmodel.Selectable
	repodb  database.Database
	repo    devicerepo.Interface
}

func (this *DeviceRepoMock) GetAspectNode(token auth.Token, id string) (aspectNode devicemodel.AspectNode, err error) {
	//TODO implement me
	panic("implement me")
}

var Devices = &DeviceRepoMock{}

func (this *DeviceRepoMock) New(ctx context.Context, config config.Config) (interfaces.Devices, error) {
	var err error
	this.repo, this.repodb, err = devicerepo.NewTestClient()
	return this, err
}

func (this *DeviceRepoMock) GetDevice(token auth.Token, id string) (devicemodel.Device, error, int) {
	return this.repo.ReadDevice(id, token.Jwt(), devicerepo.READ)
}

func (this *DeviceRepoMock) SetDevice(id string, device devicemodel.Device, userId string) error {
	device.Id = id
	err := this.repodb.SetRights("devices", device.Id, devicemodel.ResourceRights{
		UserRights: map[string]model.Right{userId: {
			Read:         true,
			Write:        true,
			Execute:      true,
			Administrate: true,
		}},
		GroupRights:          map[string]model.Right{},
		KeycloakGroupsRights: map[string]model.Right{},
	})
	if err != nil {
		return err
	}
	return this.repodb.SetDevice(context.Background(), devicerepo.DeviceWithConnectionState{Device: device})
}

func (this *DeviceRepoMock) GetService(token auth.Token, id string) (devicemodel.Service, error, int) {
	return this.repo.GetService(id)
}

func (this *DeviceRepoMock) SetService(id string, service devicemodel.Service) error {
	return this.repodb.SetDeviceType(context.Background(), models.DeviceType{
		Id:       "ref-service:" + service.Id,
		Name:     "ref-service:" + service.Name,
		Services: []models.Service{service},
	})
}

func (this *DeviceRepoMock) GetDeviceGroup(token auth.Token, id string) (result devicemodel.DeviceGroup, err error, code int) {
	panic("not implemented")
}

func (this *DeviceRepoMock) CheckAccess(token auth.Token, kind string, ids []string) (map[string]bool, error) {
	return map[string]bool{}, nil
}

func (this *DeviceRepoMock) SetOptions(options []deviceselectionmodel.Selectable) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.options = options
}

func (this *DeviceRepoMock) GetDeviceSelection(token auth.Token, descriptions deviceselectionmodel.FilterCriteriaAndSet, filterByInteraction devicemodel.Interaction) (result []deviceselectionmodel.Selectable, err error, code int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.options, nil, 200
}

func (this *DeviceRepoMock) GetBulkDeviceSelection(token auth.Token, bulk deviceselectionmodel.BulkRequest) (result deviceselectionmodel.BulkResult, err error, code int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	for _, element := range bulk {
		result = append(result, deviceselectionmodel.BulkResultElement{
			Id:          element.Id,
			Selectables: this.options,
		})
	}
	return result, nil, 200
}

func (this *DeviceRepoMock) GetBulkDeviceSelectionV2(token auth.Token, bulk deviceselectionmodel.BulkRequestV2) (result deviceselectionmodel.BulkResult, err error, code int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	for _, element := range bulk {
		result = append(result, deviceselectionmodel.BulkResultElement{
			Id:          element.Id,
			Selectables: this.options,
		})
	}
	return result, nil, 200
}
