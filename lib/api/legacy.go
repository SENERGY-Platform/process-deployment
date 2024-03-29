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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	endpoints = append(endpoints, LegacyEndpoints)
}

func LegacyEndpoints(router *httprouter.Router, config config.Config, ctrl *ctrl.Ctrl) {
	router.POST("/prepared-deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.GET("/prepared-deployments/:modelId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.POST("/deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.PUT("/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.GET("/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.DELETE("/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.POST("/v2/prepared-deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.GET("/v2/prepared-deployments/:modelId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.POST("/v2/deployments", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.PUT("/v2/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.GET("/v2/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})

	router.DELETE("/v2/deployments/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		http.Error(writer, "use /v3 endpoints", http.StatusGone)
	})
}
