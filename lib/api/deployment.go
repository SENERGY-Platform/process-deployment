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

package api

import (
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
	"time"
)

func init() {
	endpoints = append(endpoints, DeploymentsEndpoints)
}

func DeploymentsEndpoints(router *jwt_http_router.Router, config config.Config, ctrl *ctrl.Ctrl) {
	router.POST("/prepared-deployments", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		msg := model.PrepareRequest{}
		err := json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.PrepareDeployment(jwt.Impersonate, msg.Xml, msg.Svg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/prepared-deployments/:modelId", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("modelId")
		process, err, code := ctrl.GetProcessModel(jwt.Impersonate, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		start := time.Now()
		result, err, code := ctrl.PrepareDeployment(jwt.Impersonate, process.BpmnXml, process.SvgXml)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		dur := time.Now().Sub(start)
		log.Println("DEBUG: prepare deployment complete time:", dur, dur.Milliseconds())
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.POST("/deployments", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		deployment := model.Deployment{}
		err := json.NewDecoder(request.Body).Decode(&deployment)
		if err != nil {
			log.Println("ERROR: unable to parse request", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.CreateDeployment(jwt, deployment)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.PUT("/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		deployment := model.Deployment{}
		err := json.NewDecoder(request.Body).Decode(&deployment)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.UpdateDeployment(jwt, id, deployment)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		result, err, code := ctrl.GetDeployment(jwt, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.DELETE("/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		err, code := ctrl.RemoveDeployment(jwt, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(true)
	})
}
