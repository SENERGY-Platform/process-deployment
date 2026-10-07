/*
 * Copyright 2026 InfAI (CC SES)
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
	"fmt"
	"log"
	"sync"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Marshaller loads concepts and characteristics from the device-repository on start, so they have
// to be stored there before it is started. hostAccessPorts are reachable from the container under
// testcontainers.HostInternal, which allows an authEndpoint served by the test itself.
func Marshaller(ctx context.Context, wg *sync.WaitGroup, deviceRepoUrl string, converterUrl string, authEndpoint string, hostAccessPorts []int) (hostPort string, ipAddress string, err error) {
	log.Println("start marshaller")
	//the marshaller fetches a token on start and exits if it fails, but testcontainers forwards
	//the host ports only after the container is ready; so the start of the marshaller waits for
	//the forwarding, and readiness is checked after testcontainers is done, not by WaitingFor
	waitForHost := ""
	for _, port := range hostAccessPorts {
		waitForHost += fmt.Sprintf("until nc -z %s %d; do sleep 0.1; done; ", testcontainers.HostInternal, port)
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:      "ghcr.io/senergy-platform/marshaller:dev",
			Entrypoint: []string{"sh", "-c", waitForHost + "exec ./app"},
			Env: map[string]string{
				"DEBUG":                 "true",
				"DEVICE_REPOSITORY_URL": deviceRepoUrl,
				"CONVERTER_URL":         converterUrl,
				"AUTH_ENDPOINT":         authEndpoint,
				"AUTH_CLIENT_ID":        "marshaller",
			},
			ExposedPorts:    []string{"8080/tcp"},
			HostAccessPorts: hostAccessPorts,
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
		<-ctx.Done()
		log.Println("DEBUG: remove container marshaller", c.Terminate(context.Background()))
	}()

	err = wait.ForListeningPort("8080/tcp").WaitUntilReady(ctx, c)
	if err != nil {
		return "", "", err
	}

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
