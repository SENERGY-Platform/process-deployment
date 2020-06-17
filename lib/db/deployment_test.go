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
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/ory/dockertest"
	"runtime/debug"
	"strings"
	"testing"
)

func TestDeployments(t *testing.T) {
	if testing.Short() {
		t.Skip("short tests only without docker")
	}
	buffer := &strings.Builder{}
	testprint := func(args ...interface{}) {
		fmt.Fprintln(buffer, args...)
	}

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
	closer, port, _, err := MongoTestServer(pool)
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

	err = db.SetDeployment("id1", "user1", deploymentmodel.Deployment{
		Id:   "id1",
		Name: "name1",
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDeployment("id2", "user1", deploymentmodel.Deployment{
		Id:   "id2",
		Name: "name2",
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDeployment("id3", "user2", deploymentmodel.Deployment{
		Id:   "id3",
		Name: "name3",
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	testprint("GET:")
	testprint(db.GetDeployment("nope", "nope"))
	testprint(db.GetDeployment("user1", "nope"))
	testprint(db.GetDeployment("nope", "id1"))
	testprint(db.GetDeployment("user1", "id1"))

	testprint("CHECK:")
	testprint(db.CheckDeploymentAccess("nope", "nope"))
	testprint(db.CheckDeploymentAccess("user1", "nope"))
	testprint(db.CheckDeploymentAccess("nope", "id1"))
	testprint(db.CheckDeploymentAccess("user1", "id1"))

	testprint("DELETE:")
	testprint(db.GetDeployment("user1", "id1"))
	testprint(db.DeleteDeployment("nope"))
	testprint(db.GetDeployment("user1", "id1"))
	testprint(db.DeleteDeployment("id1"))
	testprint(db.GetDeployment("user1", "id1"))

	expected := `
	GET:
	{  false     [] [] } not found 404 
	{  false     [] [] } not found 404 
	{  false     [] [] } access denied 403
	{id1  false    name1 [] [] } <nil> 200
	CHECK:
	not found 404
	not found 404
	access denied 403
	<nil> 200
	DELETE:
	{id1  false    name1 [] [] } <nil> 200
	<nil>
	{id1  false    name1 [] [] } <nil> 200
	<nil>
	{  false     [] [] } not found 404
	`
	compareExampleStr(t, buffer.String(), expected)
}
