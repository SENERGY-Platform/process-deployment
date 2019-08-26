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
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/api/util"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/ctrl"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

var endpoints []func(*jwt_http_router.Router, config.Config, *ctrl.Ctrl)

func Start(ctx context.Context, config config.Config, ctrl *ctrl.Ctrl) error {
	log.Println("start api")
	router := jwt_http_router.New(jwt_http_router.JwtConfig{ForceAuth: true, ForceUser: true})
	for _, e := range endpoints {
		log.Println("add endpoints: " + runtime.FuncForPC(reflect.ValueOf(e).Pointer()).Name())
		e(router, config, ctrl)
	}
	log.Println("add logging and cors")
	corsHandler := util.NewCors(router)
	logger := util.NewLogger(corsHandler, config.LogLevel)
	log.Println("listen on port", config.ApiPort)
	server := &http.Server{Addr: ":" + config.ApiPort, Handler: logger, WriteTimeout: 10 * time.Second, ReadTimeout: 2 * time.Second, ReadHeaderTimeout: 2 * time.Second}
	server.RegisterOnShutdown(func() {
		log.Println("DEBUG: api shutdown")
	})
	go func() {
		log.Println("Listening on ", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println("ERROR: api server error", err)
			log.Fatal(err)
		}
	}()
	go func() {
		<-ctx.Done()
		log.Println("DEBUG: api shutdown result: ", server.Shutdown(context.Background()))
	}()
	return nil
}
