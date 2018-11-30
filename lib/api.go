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

package lib

import (
	"log"
	"net/http"

	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/SmartEnergyPlatform/process-deployment/lib/model"
	"github.com/SmartEnergyPlatform/process-deployment/lib/util"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/SmartEnergyPlatform/util/http/response"
	"github.com/satori/go.uuid"

	"github.com/SmartEnergyPlatform/util/http/cors"

	"encoding/json"

	"bytes"
)

func StartRest() {
	log.Println("start server on port: ", util.Config.ServerPort)
	httpHandler := getRoutes()
	corseHandler := cors.New(httpHandler)
	logger := logger.New(corseHandler, util.Config.LogLevel)
	log.Println(http.ListenAndServe(":"+util.Config.ServerPort, logger))
}

func getRoutes() (router *jwt_http_router.Router) {
	router = jwt_http_router.New(jwt_http_router.JwtConfig{
		ForceUser: util.Config.ForceUser == "true",
		ForceAuth: util.Config.ForceAuth == "true",
		PubRsa:    util.Config.JwtPubRsa,
	})

	router.GET("/deployment", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		metadata, err := GetAllMetadata(jwt.UserId)
		if err == nil {
			response.To(res).Json(metadata)
		} else {
			log.Println("error on GetAllMetadata(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/deployment/:id/dependencies", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		metadata, err := GetMetadataWithOnlineState(id, jwt.Impersonate, jwt.UserId)
		if err == nil {
			response.To(res).Json(metadata)
		} else {
			log.Println("error on GetMetadata(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.POST("/process/prepare", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		xml := buf.String()
		resp, err := GetBpmnAbstractPrepare(xml, jwt.Impersonate)
		if err != nil {
			log.Println("DEBUG: error in GetBpmnAbstractPrepare() ", err, xml)
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(resp)
	})

	router.GET("/process/clone/:processid", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("processid")
		clone, err := CloneAbstractProcess(id, jwt.Impersonate, jwt.UserId)
		if err != nil {
			log.Println("DEBUG: error in GetBpmnAbstractPrepare() ", err)
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(clone)
	})

	router.POST("/process/deploy", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		msg := model.DeploymentRequest{}
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			log.Println("ERROR: request json decode ", err)
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		msg.Process.MsgEvents, err = setFilterIds(msg.Process.MsgEvents)
		if err != nil {
			log.Println("ERROR: setFilterIds() ", err)
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		msg.Process.ReceiveTasks, err = setFilterIds(msg.Process.ReceiveTasks)
		if err != nil {
			log.Println("ERROR: setFilterIds() ", err)
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		xmlString, err := InstantiateAbstractProcess(msg.Process, jwt.Impersonate, jwt.UserId)
		if err != nil {
			log.Println("ERROR: InstantiateAbstractProcess() ", err)
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		id, err := PublishDeployment(jwt.UserId, msg, xmlString)
		if err != nil {
			log.Println("error on DeployProcess(): ", err)
			response.To(res).DefaultError("serverside error", 500)
			return
		}
		response.To(res).Text(id)
	})

	router.DELETE("/deployment/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		err := CheckAccess(id, jwt.UserId)
		if err != nil {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		err = PublishDeploymentDelete(id)
		if err != nil {
			log.Println(err)
			response.To(res).DefaultError("serverside error", 500)
			return
		}
		response.To(res).Text("ok")
	})

	return
}

func setFilterIds(events []model.MsgEvent) (result []model.MsgEvent, err error) {
	for _, event := range events {
		id := uuid.NewV4()
		if err != nil {
			return result, err
		}
		event.FilterId = "fid_" + id.String()
		result = append(result, event)
	}
	return
}
