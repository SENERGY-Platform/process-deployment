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
	"github.com/SENERGY-Platform/process-deployment/lib/api/util"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func init() {
	endpoints = append(endpoints, Deployments2Endpoints)
}

func Deployments2Endpoints(router *httprouter.Router, config config.Config, ctrl *ctrl.Ctrl) {
	router.POST("/v2/prepared-deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		msg := messages.PrepareRequest{}
		err := json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.PrepareDeploymentV2(util.GetAuthToken(request), msg.Xml, msg.Svg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/v2/prepared-deployments/:modelId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("modelId")
		process, err, code := ctrl.GetProcessModel(util.GetAuthToken(request), id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		start := time.Now()
		result, err, code := ctrl.PrepareDeploymentV2(util.GetAuthToken(request), process.BpmnXml, process.SvgXml)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		dur := time.Now().Sub(start)
		log.Println("DEBUG: prepare deployment complete time:", dur, dur.Milliseconds())
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.POST("/v2/deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		source := request.URL.Query().Get("source")
		deployment := deploymentmodel.Deployment{}
		err := json.NewDecoder(request.Body).Decode(&deployment)
		if err != nil {
			log.Println("ERROR: unable to parse request", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.CreateDeploymentV2(util.GetAuthToken(request), deployment, source)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.PUT("/v2/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		source := request.URL.Query().Get("source")
		id := params.ByName("id")
		deployment := deploymentmodel.Deployment{}
		err := json.NewDecoder(request.Body).Decode(&deployment)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.UpdateDeploymentV2(util.GetAuthToken(request), id, deployment, source)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/v2/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		result, err, code := ctrl.GetDeploymentV2(util.GetAuthToken(request), id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.DELETE("/v2/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		err, code := ctrl.RemoveDeploymentV2(util.GetAuthToken(request), id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(true)
	})
}
