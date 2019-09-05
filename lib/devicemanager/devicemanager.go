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

package devicemanager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"runtime/debug"
	"time"
)

type DeviceManagerFactory struct{}

type DeviceManager struct {
	config config.Config
}

var Factory = &DeviceManagerFactory{}

func (this *DeviceManagerFactory) New(ctx context.Context, config config.Config) (interfaces.DeviceManager, error) {
	return &DeviceManager{config: config}, nil
}

func (this *DeviceManager) GetDeploymentOptions(token jwt_http_router.JwtImpersonate, descriptions []model.DeviceDescription) (result []model.DeviceOption, err error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	payload, err := json.Marshal(descriptions)

	req, err := http.NewRequest(
		"POST",
		this.config.DeviceManagerUrl+"/device-options",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		debug.PrintStack()
		return result, err
	}
	req.Header.Set("Authorization", string(token))

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		debug.PrintStack()
		return result, errors.New("unexpected statuscode")
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
