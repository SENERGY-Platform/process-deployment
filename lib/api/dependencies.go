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
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	endpoints = append(endpoints, &DependenciesEndpoints{})
}

type DependenciesEndpoints struct{}

// ListDependencies godoc
// @Summary      list dependencies
// @Description  list dependencies
// @Tags         dependencies
// @Produce      json
// @Security Bearer
// @Param        limit query integer false "default 100, will be ignored if 'ids' is set"
// @Param        offset query integer false "default 0, will be ignored if 'ids' is set"
// @Param        ids query string false "filter; ignores limit/offset; comma-seperated list"
// @Success      200 {array} dependencymodel.Dependencies
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /dependencies [GET]
func (this *DependenciesEndpoints) ListDependencies(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("GET /dependencies", func(writer http.ResponseWriter, request *http.Request) {
		result := []dependencymodel.Dependencies{}
		var err error
		var code int
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		idstring := request.URL.Query().Get("ids")
		if idstring != "" {
			ids := strings.Split(strings.Replace(idstring, " ", "", -1), ",")
			result, err, code = ctrl.GetSelectedDependencies(token, ids)
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
			result, err, code = ctrl.GetDependenciesList(token, limit, offset)
		}
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})
}

// GetDependency godoc
// @Summary      get dependencies
// @Description  get dependencies of deployment
// @Tags         dependencies
// @Produce      json
// @Security Bearer
// @Param        id path string true "deployment id"
// @Success      200 {object}  dependencymodel.Dependencies
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /dependencies/{id} [GET]
func (this *DependenciesEndpoints) GetDependency(config config.Config, router *http.ServeMux, ctrl *ctrl.Ctrl) {
	router.HandleFunc("GET /dependencies/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		token, err := auth.GetParsedToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.GetDependencies(token, id)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(result)
	})
}
