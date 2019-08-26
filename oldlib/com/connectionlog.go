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

package com

import (
	"github.com/SENERGY-Platform/process-deployment/oldlib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"time"
)

func CheckDeviceStates(jwtimpersonate jwt_http_router.JwtImpersonate, ids []string) (result map[string]bool, err error) {
	start := time.Now()
	err = jwtimpersonate.PostJSON(util.Config.ConnectionLogUrl+"/state/device/check", ids, &result)
	if util.Config.Debug {
		log.Println("DEBUG: CheckDeviceStates()", time.Now().Sub(start))
	}
	return
}
