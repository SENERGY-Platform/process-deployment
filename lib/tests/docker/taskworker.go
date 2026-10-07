/*
 * Copyright 2025 InfAI (CC SES)
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

package docker

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"maps"
	"sync"
)

func TaskWorker(ctx context.Context, wg *sync.WaitGroup, deviceRepoUrl string, kafkaUrl string, incidentApiUrl string, shardsDbUrl string, memcachedUrl string) (err error) {
	return TaskWorkerWithEnv(ctx, wg, deviceRepoUrl, kafkaUrl, incidentApiUrl, shardsDbUrl, memcachedUrl, nil, nil)
}

// TaskWorkerWithEnv starts the worker of TaskWorker with env overriding its environment, for
// example with a reachable AUTH_ENDPOINT and MARSHALLER_URL, so that tasks are executed instead
// of ending as incidents. hostAccessPorts are reachable from the container under
// testcontainers.HostInternal.
func TaskWorkerWithEnv(ctx context.Context, wg *sync.WaitGroup, deviceRepoUrl string, kafkaUrl string, incidentApiUrl string, shardsDbUrl string, memcachedUrl string, env map[string]string, hostAccessPorts []int) (err error) {
	log.Println("start task-worker")
	containerEnv := map[string]string{
		"DEBUG":                      "true",
		"DEVICE_REPO_URL":            deviceRepoUrl,
		"KAFKA_URL":                  kafkaUrl,
		"AUTH_ENDPOINT":              "-", //may be left empty because we want incidents
		"COMPLETION_STRATEGY":        "pessimistic",
		"CAMUNDA_TOPIC":              "pessimistic",
		"INCIDENT_API_URL":           incidentApiUrl,
		"USE_HTTP_INCIDENT_PRODUCER": "true",
		"MARSHALLER_URL":             "-", //may be left empty because we want incidents
		"SHARDS_DB":                  shardsDbUrl,
		"SUB_RESULT_DATABASE_URLS":   memcachedUrl,
		"TIMESCALE_WRAPPER_URL":      "-",
	}
	maps.Copy(containerEnv, env)
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:           "ghcr.io/senergy-platform/external-task-worker:dev",
			Env:             containerEnv,
			HostAccessPorts: hostAccessPorts,
			AlwaysPullImage: true,
		},
		Started: true,
	})
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Println("DEBUG: remove container task-worker", c.Terminate(context.Background()))
		}()
		<-ctx.Done()
		/*
			reader, err := c.Logs(context.Background())
			if err != nil {
				log.Println("ERROR: unable to get container log")
				return
			}
			buf := new(strings.Builder)
			io.Copy(buf, reader)
			fmt.Println("TASK-WORKER LOGS: ------------------------------------------")
			fmt.Println(buf.String())
			fmt.Println("\n---------------------------------------------------------------")
		*/
	}()

	return err
}
