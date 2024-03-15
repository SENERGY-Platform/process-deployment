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
	"github.com/SENERGY-Platform/service-commons/pkg/accesslog"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

var endpoints []func(*httprouter.Router, config.Config, *ctrl.Ctrl)

func Start(ctx context.Context, config config.Config, ctrl *ctrl.Ctrl) error {
	log.Println("start api")

	timeout, err := time.ParseDuration(config.HttpServerTimeout)
	if err != nil {
		log.Println("WARNING: invalid http server timeout --> no timeouts\n", err)
		err = nil
	}

	readtimeout, err := time.ParseDuration(config.HttpServerReadTimeout)
	if err != nil {
		log.Println("WARNING: invalid http server read timeout --> no timeouts\n", err)
		err = nil
	}

	router := httprouter.New()
	for _, e := range endpoints {
		log.Println("add endpoints: " + runtime.FuncForPC(reflect.ValueOf(e).Pointer()).Name())
		e(router, config, ctrl)
	}
	log.Println("add logging and cors")
	corsHandler := util.NewCors(router)
	logger := accesslog.New(corsHandler)
	server := &http.Server{Addr: ":" + config.ApiPort, Handler: logger, WriteTimeout: timeout, ReadTimeout: readtimeout}
	go func() {
		log.Println("Listening on ", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Println("ERROR: api server error", err)
				log.Fatal(err)
			} else {
				log.Println("closing api server")
			}
		}
	}()
	go func() {
		<-ctx.Done()
		log.Println("DEBUG: api shutdown", server.Shutdown(context.Background()))
	}()
	return nil
}
