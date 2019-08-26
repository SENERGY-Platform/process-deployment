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

package oldlib

import (
	"log"

	"github.com/SENERGY-Platform/process-deployment/oldlib/com"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

func GetMetadataWithOnlineState(id string, jwtimpersonate jwt_http_router.JwtImpersonate, owner string) (metadata Metadata, err error) {
	metadata, err = GetMetadata(id, owner)
	if err != nil {
		log.Println("ERROR in GetMetadataWithOnlineState::GetMetadata()", err, id, owner)
		return
	}
	metadata.Online = true

	ids := []string{}
	for _, param := range metadata.Abstract.AbstractTasks {
		ids = append(ids, param.Selected.Id)
	}
	deviceStates := map[string]bool{}
	if len(ids) > 0 {
		deviceStates, err = com.CheckDeviceStates(jwtimpersonate, ids)
		if err != nil {
			log.Println("WARNING: error in CheckDeviceStates()", err)
		}
		if deviceStates == nil {
			deviceStates = map[string]bool{}
		}
	}
	for index, param := range metadata.Abstract.AbstractTasks {
		state, ok := deviceStates[param.Selected.Id]
		if ok && !state {
			log.Println("OFFLING: device state = ", state)
			metadata.Online = false
		}
		if !ok {
			param.State = "unknown"
		} else if state {
			param.State = "connected"
		} else {
			param.State = "disconnected"
		}
		metadata.Abstract.AbstractTasks[index] = param
	}

	return
}

func GetMetadataListWithOnlineState(ids []string, jwtimpersonate jwt_http_router.JwtImpersonate, owner string) (result []Metadata, err error) {
	metadataList, err := GetMetadataList(ids, owner)
	if err != nil {
		log.Println("ERROR in GetMetadataWithOnlineState::GetMetadata()", err, ids, owner)
		return
	}
	deviceIds := []string{}
	deviceSet := map[string]bool{}
	for _, metadata := range metadataList {
		for _, param := range metadata.Abstract.AbstractTasks {
			if !deviceSet[param.Selected.Id] {
				deviceSet[param.Selected.Id] = true
				deviceIds = append(deviceIds, param.Selected.Id)
			}
		}
	}

	deviceStates := map[string]bool{}
	if len(deviceIds) > 0 {
		deviceStates, err = com.CheckDeviceStates(jwtimpersonate, deviceIds)
		if err != nil {
			log.Println("WARNING: error in CheckDeviceStates()", err)
		}
		if deviceStates == nil {
			deviceStates = map[string]bool{}
		}
	}

	for _, metadata := range metadataList {
		metadata.Online = true
		for index, param := range metadata.Abstract.AbstractTasks {
			state, ok := deviceStates[param.Selected.Id]
			if ok && !state {
				//log.Println("OFFLINE: device state = ", state)
				metadata.Online = false
			}
			if !ok {
				param.State = "unknown"
			} else if state {
				param.State = "connected"
			} else {
				param.State = "disconnected"
			}
			metadata.Abstract.AbstractTasks[index] = param
		}
		result = append(result, metadata)
	}
	return
}
