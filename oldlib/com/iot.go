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
	"log"
	"net/url"

	"github.com/SENERGY-Platform/process-deployment/oldlib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func GetAllowedValues(jwthttp jwt_http_router.JwtImpersonate) (result model.AllowedValues, err error) {
	err = jwthttp.GetJSON(util.Config.IotRepoUrl+"/ui/deviceType/allowedvalues", &result)
	return
}

func GetDeviceInstance(deviceInstanceId string, jwthttp jwt_http_router.JwtImpersonate) (result model.DeviceInstance, err error) {
	err = jwthttp.GetJSON(util.Config.IotRepoUrl+"/deviceInstance/"+url.QueryEscape(deviceInstanceId), &result)
	return
}

func GetDeviceType(deviceTypeId string, jwthttp jwt_http_router.JwtImpersonate) (result model.DeviceType, err error) {
	err = jwthttp.GetJSON(util.Config.IotRepoUrl+"/deviceType/"+url.QueryEscape(deviceTypeId), &result)
	return
}

func GetDeviceService(serviceId string, jwthttp jwt_http_router.JwtImpersonate) (result model.Service, err error) {
	err = jwthttp.GetJSON(util.Config.IotRepoUrl+"/service/"+url.QueryEscape(serviceId), &result)
	return
}

func GetValueType(valueTypeId string, jwthttp jwt_http_router.JwtImpersonate) (result model.ValueType, err error) {
	err = jwthttp.GetJSON(util.Config.IotRepoUrl+"/valueType/"+url.QueryEscape(valueTypeId), &result)
	return
}

func GetDeviceInstancesFromType(deviceTypeId string, jwthttp jwt_http_router.JwtImpersonate) (result []model.DeviceInstance, err error) {
	err = jwthttp.GetJSON(util.Config.IotRepoUrl+"/bydevicetype/deviceInstances/"+url.QueryEscape(deviceTypeId)+"/execute", &result)
	return
}

func GetDeviceInstancesFromService(serviceid string, jwthttp jwt_http_router.JwtImpersonate) (result []model.DeviceInstance, err error) {
	if serviceid == "" {
		return []model.DeviceInstance{}, nil
	}
	err = jwthttp.GetJSON(util.Config.IotRepoUrl+"/byservice/deviceInstances/"+url.QueryEscape(serviceid)+"/execute", &result)
	return
}

func CheckExecuteRight(id string, jwthttp jwt_http_router.JwtImpersonate) (err error) {
	return CheckRight(id, jwthttp, "x")
}

func CheckRight(id string, jwthttp jwt_http_router.JwtImpersonate, right string) (err error) {
	kind := "deviceinstance"
	result := false
	err = jwthttp.GetJSON(util.Config.PermissionsUrl+"/jwt/check/"+kind+"/"+url.QueryEscape(id)+"/"+right+"/bool", &result)
	if err != nil {
		log.Println("DEBUG: permissions.Check:", err)
		return err
	}
	return
}
