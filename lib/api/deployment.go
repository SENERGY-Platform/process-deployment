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
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"time"
)

func init() {
	endpoints = append(endpoints, DeploymentsEndpoints)
}

func DeploymentsEndpoints(router *httprouter.Router, config config.Config, ctrl *ctrl.Ctrl) {
	router.GET("/v3/start-parameters/:modelId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("modelId")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		process, err, code := ctrl.GetProcessModel(token, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		result, err := ctrl.GetProcessStartParameters(process.BpmnXml)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.POST("/v3/prepared-deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		withOptionsStr := request.URL.Query().Get("with_options")
		if withOptionsStr == "" {
			withOptionsStr = "true"
		}
		withOptions, err := strconv.ParseBool(withOptionsStr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		msg := messages.PrepareRequest{}
		err = json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.PrepareDeployment(token, msg.Xml, msg.Svg, withOptions)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/v3/prepared-deployments/:modelId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		withOptionsStr := request.URL.Query().Get("with_options")
		if withOptionsStr == "" {
			withOptionsStr = "true"
		}
		withOptions, err := strconv.ParseBool(withOptionsStr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		id := params.ByName("modelId")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		process, err, code := ctrl.GetProcessModel(token, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		start := time.Now()
		result, err, code := ctrl.PrepareDeployment(token, process.BpmnXml, process.SvgXml, withOptions)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		dur := time.Now().Sub(start)
		log.Println("DEBUG: prepare deployment complete time:", dur, dur.Milliseconds())
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.POST("/v3/deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		source := request.URL.Query().Get("source")
		deployment := deploymentmodel.Deployment{}
		err := json.NewDecoder(request.Body).Decode(&deployment)
		if err != nil {
			log.Println("ERROR: unable to parse request", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		optionals := map[string]bool{}
		optionalServiceStr := request.URL.Query().Get("optional_service_selection")
		if optionalServiceStr != "" {
			optionals["service"], err = strconv.ParseBool(optionalServiceStr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		result, err, code := ctrl.CreateDeployment(token, deployment, source, optionals)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.PUT("/v3/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		source := request.URL.Query().Get("source")
		id := params.ByName("id")
		deployment := deploymentmodel.Deployment{}
		err := json.NewDecoder(request.Body).Decode(&deployment)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		optionals := map[string]bool{}
		optionalServiceStr := request.URL.Query().Get("optional_service_selection")
		if optionalServiceStr != "" {
			optionals["service"], err = strconv.ParseBool(optionalServiceStr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		result, err, code := ctrl.UpdateDeployment(token, id, deployment, source, optionals)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/v3/deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		listOptions := model.DeploymentListOptions{
			Limit:  0,
			Offset: 0,
		}
		var err error
		limitParam := request.URL.Query().Get("limit")
		if limitParam != "" {
			listOptions.Limit, err = strconv.ParseInt(limitParam, 10, 64)
		}
		if err != nil {
			http.Error(writer, "unable to parse limit:"+err.Error(), http.StatusBadRequest)
			return
		}

		offsetParam := request.URL.Query().Get("offset")
		if offsetParam != "" {
			listOptions.Offset, err = strconv.ParseInt(offsetParam, 10, 64)
		}
		if err != nil {
			http.Error(writer, "unable to parse offset:"+err.Error(), http.StatusBadRequest)
			return
		}

		listOptions.SortBy = request.URL.Query().Get("sort")
		if listOptions.SortBy == "" {
			listOptions.SortBy = "name.asc"
		}
		
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.GetDeployments(token, listOptions)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/v3/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		withOptionsStr := request.URL.Query().Get("with_options")
		if withOptionsStr == "" {
			withOptionsStr = "true"
		}
		withOptions, err := strconv.ParseBool(withOptionsStr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		id := params.ByName("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.GetDeployment(token, id, withOptions)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.DELETE("/v3/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		err, code := ctrl.RemoveDeployment(token, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(true)
	})
}
