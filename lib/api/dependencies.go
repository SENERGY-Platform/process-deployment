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
	"net/http"
	"strconv"
	"strings"
)

func init() {
	endpoints = append(endpoints, DependenciesEndpoints)
}

func DependenciesEndpoints(router *jwt_http_router.Router, config config.Config, ctrl *ctrl.Ctrl) {

	/*
		query-parameter:
			ids:
				comma separated list of deployment ids
				filters dependencies by given deployments
			limit
			offset
	*/
	router.GET("/dependencies", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		result := []model.Dependencies{}
		var err error
		var code int

		idstring := request.URL.Query().Get("ids")
		if idstring != "" {
			ids := strings.Split(strings.Replace(idstring, " ", "", -1), ",")
			result, err, code = ctrl.GetSelectedDependencies(jwt, ids)
		} else {
			var limit int = 100
			var offset int = 0
			limitstr := request.URL.Query().Get("limit")
			if limitstr != "" {
				limit, err = strconv.Atoi(limitstr)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
					return
				}
			}
			offsetstr := request.URL.Query().Get("offset")
			if offsetstr != "" {
				offset, err = strconv.Atoi(offsetstr)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
					return
				}
			}
			result, err, code = ctrl.GetDependenciesList(jwt, limit, offset)
		}
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})

	router.GET("/dependencies/:id", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := params.ByName("id")
		result, err, code := ctrl.GetDependencies(jwt, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})
}
