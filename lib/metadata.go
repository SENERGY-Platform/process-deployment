/*
 * Copyright 2018 InfAI (CC SES)
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

package lib

import (
	"log"

	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/SmartEnergyPlatform/process-deployment/lib/com"
)

func GetMetadataWithOnlineState(id string, jwtimpersonate jwt_http_router.JwtImpersonate, owner string) (metadata Metadata, err error) {
	metadata, err = GetMetadata(id, owner)
	if err != nil {
		log.Println("ERROR in GetMetadataWithOnlineState::GetMetadata()", err, id, owner)
		return
	}
	metadata.Online = true

	for index, param := range metadata.Abstract.AbstractTasks {
		state, err := com.GetDeviceState(param.Selected.Id, jwtimpersonate)
		if err != nil {
			log.Println("ERROR in GetDeviceState()", err)
			return metadata, err
		}
		instance, err := com.GetDeviceInstance(param.Selected.Id, jwtimpersonate)
		if err != nil {
			log.Println("ERROR in GetDeviceInstance(): ", err)
			state = "iot_repo_error"
		}
		if instance.Id != param.Selected.Id || instance.DeviceType == "" {
			log.Println("OFFLINE: inconsistend repo: ", instance.Id, param.Selected.Id, instance.DeviceType)
			state = "inconsistent_iot_repo"
		}
		if state != "connected" && state != "" && state != "unknown" {
			log.Println("OFFLING: device state = ", state)
			metadata.Online = false
		}
		param.State = state
		metadata.Abstract.AbstractTasks[index] = param
	}

	for index, event := range metadata.Abstract.MsgEvents {
		state, err := com.GetEventState(event.FilterId)
		if err != nil {
			log.Println("ERROR in GetEventState()", err)
			return metadata, err
		}
		if state != "running" {
			log.Println("OFFLINE: event state = ", state)
			metadata.Online = false
		}
		event.State = state
		metadata.Abstract.MsgEvents[index] = event
	}

	for index, event := range metadata.Abstract.ReceiveTasks {
		state, err := com.GetEventState(event.FilterId)
		if err != nil {
			log.Println("ERROR in GetEventState()[1]", err)
			return metadata, err
		}
		if state != "running" {
			log.Println("OFFLINE: receive task state = ", state)
			metadata.Online = false
		}
		event.State = state
		metadata.Abstract.ReceiveTasks[index] = event
	}

	return
}
