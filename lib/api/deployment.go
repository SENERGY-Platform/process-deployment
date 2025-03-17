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
	"log"
	"net/http"
	"strconv"
	"time"
)

func init() {
	endpoints = append(endpoints, &DeploymentsEndpoints{})
}

type DeploymentsEndpoints struct{}

// GetStartParameters godoc
// @Summary      get start-parameters
// @Description  get start-parameters of a process-model
// @Tags         deployment
// @Produce      json
// @Security Bearer
// @Param        modelId path string true "process-model id"
// @Success      200 {array}  deploymentmodel.ProcessStartParameter
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /v3/start-parameters/{modelId} [GET]
func (this *DeploymentsEndpoints) GetStartParameters(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("GET /v3/start-parameters/{modelId}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("modelId")
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
}

// PrepareDeployment godoc
// @Summary      prepare process deployment
// @Description  prepare process deployment
// @Tags         deployment
// @Produce      json
// @Security Bearer
// @Param        with_options query bool false "default true, omit SelectionOptions if set to false"
// @Param        message body messages.PrepareRequest true "model that should be prepared for deployment"
// @Success      200 {object} deploymentmodel.Deployment
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /v3/prepared-deployments [POST]
func (this *DeploymentsEndpoints) PrepareDeployment(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("POST /v3/prepared-deployments", func(writer http.ResponseWriter, request *http.Request) {
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
}

// PrepareDeploymentByModelId godoc
// @Summary      prepare process deployment with model-id
// @Description  prepare process deployment with model-id
// @Tags         deployment
// @Produce      json
// @Security Bearer
// @Param        modelId path string true "process-model id"
// @Param        with_options query bool false "default true, omit SelectionOptions if set to false"
// @Success      200 {object} deploymentmodel.Deployment
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /v3/prepared-deployments/{modelId} [GET]
func (this *DeploymentsEndpoints) PrepareDeploymentByModelId(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("GET /v3/prepared-deployments/{modelId}", func(writer http.ResponseWriter, request *http.Request) {
		withOptionsStr := request.URL.Query().Get("with_options")
		if withOptionsStr == "" {
			withOptionsStr = "true"
		}
		withOptions, err := strconv.ParseBool(withOptionsStr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		id := request.PathValue("modelId")
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
}

// Deploy godoc
// @Summary      deploy process
// @Description  deploy process
// @Tags         deployment
// @Produce      json
// @Security Bearer
// @Param        source query string false "source of deployment (e.g. smart-service)"
// @Param        optional_service_selection query integer false "set to true to disable validation, that a service must be selected"
// @Param        message body deploymentmodel.Deployment true "process deployment"
// @Success      200 {object} deploymentmodel.Deployment
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /v3/deployments [POST]
func (this *DeploymentsEndpoints) Deploy(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("POST /v3/deployments", func(writer http.ResponseWriter, request *http.Request) {
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
}

// UpdateDeployment godoc
// @Summary      update process deployment
// @Description  update process deployment
// @Tags         deployment
// @Produce      json
// @Security Bearer
// @Param        id path string true "deployment id"
// @Param        source query string false "source of deployment (e.g. smart-service)"
// @Param        optional_service_selection query integer false "set to true to disable validation, that a service must be selected"
// @Param        message body deploymentmodel.Deployment true "process deployment"
// @Success      200 {object} deploymentmodel.Deployment
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /v3/deployments/{id} [PUT]
func (this *DeploymentsEndpoints) UpdateDeployment(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("PUT /v3/deployments/{id}", func(writer http.ResponseWriter, request *http.Request) {
		source := request.URL.Query().Get("source")
		id := request.PathValue("id")
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
}

// ListDeployments godoc
// @Summary      list process deployments
// @Description  list process deployments
// @Tags         deployment
// @Produce      json
// @Security Bearer
// @Param        limit query integer false "default unlimited"
// @Param        offset query integer false "default 0"
// @Param        sort query string false "default name.asc"
// @Success      200 {array} deploymentmodel.Deployment
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /v3/deployments [GET]
func (this *DeploymentsEndpoints) ListDeployments(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("GET /v3/deployments", func(writer http.ResponseWriter, request *http.Request) {
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
}

// GetDeployment godoc
// @Summary      get process deployment
// @Description  get process deployment
// @Tags         deployment
// @Produce      json
// @Security Bearer
// @Param        id path string true "deployment id"
// @Param        with_options query bool false "default true, omit SelectionOptions if set to false"
// @Success      200 {object} deploymentmodel.Deployment
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /v3/deployments/{id} [GET]
func (this *DeploymentsEndpoints) GetDeployment(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("GET /v3/deployments/{id}", func(writer http.ResponseWriter, request *http.Request) {
		withOptionsStr := request.URL.Query().Get("with_options")
		if withOptionsStr == "" {
			withOptionsStr = "true"
		}
		withOptions, err := strconv.ParseBool(withOptionsStr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		id := request.PathValue("id")
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
}

// DeleteDeployment godoc
// @Summary      delete process deployment
// @Description  delete process deployment
// @Tags         deployment
// @Security Bearer
// @Param        id path string true "deployment id"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /v3/deployments/{id} [DELETE]
func (this *DeploymentsEndpoints) DeleteDeployment(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("DELETE /v3/deployments/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
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
