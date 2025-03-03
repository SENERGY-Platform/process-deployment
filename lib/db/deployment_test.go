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
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
)

func TestDeploymentsV1(t *testing.T) {
	if testing.Short() {
		t.Skip("short tests only without docker")
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

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

	port, _, err := docker.Mongo(ctx, &wg)
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	config.MongoUrl = "mongodb://localhost:" + port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := Factory.New(ctx, config)
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDeployment("id1", "user1", &deploymentmodel.Deployment{
		Id:   "id1",
		Name: "name1",
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDeployment("id2", "user1", &deploymentmodel.Deployment{
		Id:   "id2",
		Name: "name2",
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDeployment("id3", "user2", &deploymentmodel.Deployment{
		Id:   "id3",
		Name: "name3",
	})
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	testprint("GET V1:")
	testprint(testGetDeploymentV1(db, "nope", "nope"))
	testprint(testGetDeploymentV1(db, "user1", "nope"))
	testprint(testGetDeploymentV1(db, "nope", "id1"))
	testprint(testGetDeploymentV1(db, "user1", "id1"))

	testprint("GET V2:")
	testprint(testGetDeploymentV2(db, "nope", "nope"))
	testprint(testGetDeploymentV2(db, "user1", "nope"))
	testprint(testGetDeploymentV2(db, "nope", "id1"))
	testprint(testGetDeploymentV2(db, "user1", "id1"))

	testprint("CHECK:")
	testprint(db.CheckDeploymentAccess("nope", "nope"))
	testprint(db.CheckDeploymentAccess("user1", "nope"))
	testprint(db.CheckDeploymentAccess("nope", "id1"))
	testprint(db.CheckDeploymentAccess("user1", "id1"))

	testprint("DELETE:")
	testprint(testGetDeploymentV1(db, "user1", "id1"))
	testprint(db.DeleteDeployment("nope"))
	testprint(testGetDeploymentV1(db, "user1", "id1"))
	testprint(db.DeleteDeployment("id1"))
	testprint(testGetDeploymentV1(db, "user1", "id1"))

	expected := `
	GET V1:
	{ false     [] [] } not found 404 
	{ false     [] [] } not found 404 
	{ false     [] [] } access denied 403
	{id1 false    name1 [] [] } <nil> 200
	GET V2:
	{   {  } [] false} not found 404 
	{   {  } [] false} not found 404 
	{   {  } [] false} access denied 403
	{   {  } [] false} <nil> 200
	CHECK:
	not found 404
	not found 404
	access denied 403
	<nil> 200
	DELETE:
	{id1 false    name1 [] [] } <nil> 200
	<nil>
	{id1 false    name1 [] [] } <nil> 200
	<nil>
	{ false     [] [] } not found 404  
	`
	compareExampleStr(t, buffer.String(), expected)
}

func TestDeploymentsV2(t *testing.T) {
	if testing.Short() {
		t.Skip("short tests only without docker")
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

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

	port, _, err := docker.Mongo(ctx, &wg)
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	config.MongoUrl = "mongodb://localhost:" + port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := Factory.New(ctx, config)
	if err != nil {
		debug.PrintStack()
		testprint(err)
		return
	}

	err = db.SetDeployment("id1", "user1", &deploymentmodel.Deployment{
		Id:   "id1",
		Name: "name1",
	})
	if err != nil {
		testprint(err)
		return
	}

	err = db.SetDeployment("id2", "user1", &deploymentmodel.Deployment{
		Id:   "id2",
		Name: "name2",
	})
	if err != nil {
		testprint(err)
		return
	}

	err = db.SetDeployment("id3", "user2", &deploymentmodel.Deployment{
		Id:   "id3",
		Name: "name3",
	})
	if err != nil {
		testprint(err)
		return
	}

	list, err := db.ListDeployments("user1", model.DeploymentListOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(list, []deploymentmodel.Deployment{
		{
			Id:   "id1",
			Name: "name1",
		},
		{
			Id:   "id2",
			Name: "name2",
		},
	}) {
		t.Errorf("%#v", list)
		return
	}

	list, err = db.ListDeployments("user1", model.DeploymentListOptions{SortBy: "name.desc"})
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(list, []deploymentmodel.Deployment{
		{
			Id:   "id2",
			Name: "name2",
		},
		{
			Id:   "id1",
			Name: "name1",
		},
	}) {
		t.Errorf("%#v", list)
		return
	}

	list, err = db.ListDeployments("user1", model.DeploymentListOptions{Limit: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(list, []deploymentmodel.Deployment{
		{
			Id:   "id1",
			Name: "name1",
		},
	}) {
		t.Errorf("%#v", list)
		return
	}

	list, err = db.ListDeployments("user1", model.DeploymentListOptions{Offset: 1})
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(list, []deploymentmodel.Deployment{
		{
			Id:   "id2",
			Name: "name2",
		},
	}) {
		t.Errorf("%#v", list)
		return
	}

	testprint("GET V1:")
	testprint(testGetDeploymentV1(db, "nope", "nope"))
	testprint(testGetDeploymentV1(db, "user1", "nope"))
	testprint(testGetDeploymentV1(db, "nope", "id1"))
	testprint(testGetDeploymentV1(db, "user1", "id1"))

	testprint("GET V2:")
	testprint(testGetDeploymentV2(db, "nope", "nope"))
	testprint(testGetDeploymentV2(db, "user1", "nope"))
	testprint(testGetDeploymentV2(db, "nope", "id1"))
	testprint(testGetDeploymentV2(db, "user1", "id1"))

	testprint("CHECK:")
	testprint(db.CheckDeploymentAccess("nope", "nope"))
	testprint(db.CheckDeploymentAccess("user1", "nope"))
	testprint(db.CheckDeploymentAccess("nope", "id1"))
	testprint(db.CheckDeploymentAccess("user1", "id1"))

	testprint("DELETE:")
	testprint(testGetDeploymentV2(db, "user1", "id1"))
	testprint(db.DeleteDeployment("nope"))
	testprint(testGetDeploymentV2(db, "user1", "id1"))
	testprint(db.DeleteDeployment("id1"))
	testprint(testGetDeploymentV2(db, "user1", "id1"))

	expected := `
	GET V1:
	{ false     [] [] } not found 404 
	{ false     [] [] } not found 404 
	{ false     [] [] } access denied 403
	{ false     [] [] } <nil> 200 
	GET V2:
	{   {  } [] false} not found 404 
	{   {  } [] false} not found 404 
	{   {  } [] false} access denied 403
	{id1 name1  {  } [] false} <nil> 200
	CHECK:
	not found 404
	not found 404
	access denied 403
	<nil> 200
	DELETE:
	{id1 name1  {  } [] false} <nil> 200
	<nil>
	{id1 name1  {  } [] false} <nil> 200
	<nil>
	{   {  } [] false} not found 404
	`
	compareExampleStr(t, buffer.String(), expected)
}

func testGetDeploymentV1(db interfaces.Database, userId string, deploymentId string) (deployment deploymentmodel.Deployment, err error, code int) {
	temp, err, code := db.GetDeployment(userId, deploymentId)
	if temp != nil {
		deployment = *temp
	}
	return deployment, err, code
}

func testGetDeploymentV2(db interfaces.Database, userId string, deploymentId string) (deployment deploymentmodel.Deployment, err error, code int) {
	temp, err, code := db.GetDeployment(userId, deploymentId)
	if temp != nil {
		deployment = *temp
	}
	return deployment, err, code
}

func compareExampleStr(t *testing.T, actual string, expected string) {
	actual = strings.TrimSpace(actual)
	expected = strings.TrimSpace(expected)
	actualLines := strings.Split(actual, "\n")
	expectedLines := strings.Split(expected, "\n")
	if len(actualLines) != len(expectedLines) {
		t.Fatal("GOT:\n", actual, "\nWANT:\n", expected)
	}
	for index, actualLine := range actualLines {
		if strings.TrimSpace(actualLine) != strings.TrimSpace(expectedLines[index]) {
			t.Fatal("LINE: ", index+1, "\nGOT:\n", strings.TrimSpace(actualLine), "\nWANT:\n", strings.TrimSpace(expectedLines[index]), "\n\n", actual)
		}
	}
}
