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
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"sync"
)

func EventDeployment(ctx context.Context, wg *sync.WaitGroup, deviceRepoUrl string, mongoUrl string) (hostPort string, ipAddress string, err error) {
	log.Println("start event-deployment")
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "ghcr.io/senergy-platform/event-deployment:dev",
			Env: map[string]string{
				"DEBUG":                            "true",
				"DISABLE_KAFKA":                    "true",
				"DEVICE_REPOSITORY_URL":            deviceRepoUrl,
				"AUTH_ENDPOINT":                    "-",
				"CONDITIONAL_EVENT_REPO_MONGO_URL": mongoUrl,
				"ENABLE_ANALYTICS_EVENTS":          "false",
			},
			ExposedPorts:    []string{"8080/tcp"},
			WaitingFor:      wait.ForListeningPort("8080/tcp"),
			AlwaysPullImage: true,
		},
		Started: true,
	})
	if err != nil {
		return "", "", err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Println("DEBUG: remove container event-deployment", c.Terminate(context.Background()))
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
			fmt.Println("EVENT-DEPLOYMENT LOGS: ------------------------------------------")
			fmt.Println(buf.String())
			fmt.Println("\n---------------------------------------------------------------")
		*/
	}()

	ipAddress, err = c.ContainerIP(ctx)
	if err != nil {
		return "", "", err
	}
	temp, err := c.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return "", "", err
	}
	hostPort = temp.Port()

	return hostPort, ipAddress, err
}

func EventWorker(ctx context.Context, wg *sync.WaitGroup, deviceRepoUrl string, eventTriggerUrl string, kafkaUrl string, mongoUrl string) (err error) {
	log.Println("start event-worker")
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "ghcr.io/senergy-platform/event-worker:dev",
			Env: map[string]string{
				"DEBUG":             "true",
				"AUTH_ENDPOINT":     "-",
				"DEVICE_REPO_URL":   deviceRepoUrl,
				"NOTIFICATION_URL":  "-",
				"EVENT_TRIGGER_URL": eventTriggerUrl,
				"KAFKA_URL":         kafkaUrl,
				"WATCHED_PROCESS_DEPLOYMENT_DONE_HANDLER": "github.com/SENERGY-Platform/process-deployment",
				"CLOUD_EVENT_REPO_MONGO_URL":              mongoUrl,
			},
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
			log.Println("DEBUG: remove container event-worker", c.Terminate(context.Background()))
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
			fmt.Println("EVENT-WORKER LOGS: ------------------------------------------")
			fmt.Println(buf.String())
			fmt.Println("\n---------------------------------------------------------------")
		*/
	}()

	return err
}
