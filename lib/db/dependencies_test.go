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

package db

import (
	"context"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
	"github.com/ory/dockertest/v3"
	"runtime/debug"
	"strings"
	"testing"
)

func TestDependencies(t *testing.T) {
	if testing.Short() {
		t.Skip("short tests only without docker")
	}
	buffer := &strings.Builder{}
	testprint := func(args ...interface{}) {
		fmt.Fprintln(buffer, args...)
	}

	testprint("")
	config, err := config.LoadConfig("../../config.json")
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}
	config.Debug = true

	pool, err := dockertest.NewPool("")
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}
	closer, port, _, err := docker.MongoTestServer(pool)
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}
	defer closer()

	config.MongoUrl = "mongodb://localhost:" + port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := Factory.New(ctx, config)
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDependencies(dependencymodel.Dependencies{
		DeploymentId: "id1",
		Owner:        "user1",
		Devices: []dependencymodel.DeviceDependency{{
			DeviceId:      "d1",
			Name:          "d1",
			BpmnResources: []dependencymodel.BpmnResource{{Id: "r1"}, {Id: "r2"}},
		}},
		Events: nil,
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDependencies(dependencymodel.Dependencies{
		DeploymentId: "id2",
		Owner:        "user1",
		Devices: []dependencymodel.DeviceDependency{{
			DeviceId:      "d1",
			Name:          "d1",
			BpmnResources: []dependencymodel.BpmnResource{{Id: "r1"}, {Id: "r2"}},
		}},
		Events: nil,
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDependencies(dependencymodel.Dependencies{
		DeploymentId: "id3",
		Owner:        "user2",
		Devices: []dependencymodel.DeviceDependency{{
			DeviceId:      "d1",
			Name:          "d1",
			BpmnResources: []dependencymodel.BpmnResource{{Id: "r1"}, {Id: "r2"}},
		}},
		Events: nil,
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	testprint("GET:")
	testprint(db.GetDependencies("nope", "nope"))
	testprint(db.GetDependencies("user1", "nope"))
	testprint(db.GetDependencies("nope", "id1"))
	testprint(db.GetDependencies("user1", "id1"))

	testprint("LIST:")
	testprint(db.GetDependenciesList("nope", 100, 0))
	testprint(db.GetDependenciesList("user1", 100, 0))

	testprint("SELECT:")
	testprint(db.GetSelectedDependencies("nope", []string{"nope"}))
	testprint(db.GetSelectedDependencies("nope", []string{"nope", "id1"}))
	testprint(db.GetSelectedDependencies("nope", []string{"id2", "id1"}))
	testprint(db.GetSelectedDependencies("user1", []string{"nope", "id1"}))
	testprint(db.GetSelectedDependencies("user1", []string{"id2", "id1"}))

	testprint("DELETE:")
	testprint(db.GetDependencies("user1", "id1"))
	testprint(db.DeleteDependencies("nope"))
	testprint(db.GetDependencies("user1", "id1"))
	testprint(db.DeleteDependencies("id1"))
	testprint(db.GetDependencies("user1", "id1"))

	expected := `
	GET:
	{  [] []} not found 404
	{  [] []} not found 404
	{  [] []} access denied 403
	{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} <nil> 200
	LIST:
	[] <nil> 200
	[{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} {id2 user1 [{d1 d1 [{r1 } {r2 }]}] []}] <nil> 200
	SELECT:
	[] <nil> 200
	[] <nil> 200
	[] <nil> 200
	[{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []}] <nil> 200
	[{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} {id2 user1 [{d1 d1 [{r1 } {r2 }]}] []}] <nil> 200
	DELETE:
	{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} <nil> 200
	<nil>
	{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} <nil> 200
	<nil>
	{  [] []} not found 404
	`

	compareExampleStr(t, buffer.String(), expected)
}
