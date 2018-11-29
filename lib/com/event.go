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
	"github.com/SmartEnergyPlatform/process-deployment/lib/util"
	"github.com/SmartEnergyPlatform/util/http/request"
	"net/url"
)


func GetEventState(filterId string) (state string, err error) {
	result := struct {
		State string `json:"state"`
	}{}
	err = request.Get(util.Config.EventManagerUrl+"/filter/"+url.QueryEscape(filterId), &result)
	state = result.State
	return
}
