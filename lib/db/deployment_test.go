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
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/ory/dockertest"
	"log"
	"runtime/debug"
)

func ExampleDeployments() {
	config, err := config.LoadConfig("../../config.json")
	if err != nil {
		debug.PrintStack()
		log.Println(err)
		return
	}
	config.Debug = true

	pool, err := dockertest.NewPool("")
	if err != nil {
		debug.PrintStack()
		log.Println(err)
		return
	}
	closer, port, _, err := MongoTestServer(pool)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}
	defer closer()

	config.MongoUrl = "mongodb://localhost:" + port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := Factory.New(ctx, config)
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	err = db.SetDeployment("id1", "user1", model.Deployment{
		Id:   "id1",
		Name: "name1",
	})
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	err = db.SetDeployment("id2", "user1", model.Deployment{
		Id:   "id2",
		Name: "name2",
	})
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	err = db.SetDeployment("id3", "user2", model.Deployment{
		Id:   "id3",
		Name: "name3",
	})
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	fmt.Println("GET:")
	fmt.Println(db.GetDeployment("nope", "nope"))
	fmt.Println(db.GetDeployment("user1", "nope"))
	fmt.Println(db.GetDeployment("nope", "id1"))
	fmt.Println(db.GetDeployment("user1", "id1"))

	fmt.Println("CHECK:")
	fmt.Println(db.CheckDeploymentAccess("nope", "nope"))
	fmt.Println(db.CheckDeploymentAccess("user1", "nope"))
	fmt.Println(db.CheckDeploymentAccess("nope", "id1"))
	fmt.Println(db.CheckDeploymentAccess("user1", "id1"))

	fmt.Println("DELETE:")
	fmt.Println(db.GetDeployment("user1", "id1"))
	fmt.Println(db.DeleteDeployment("nope"))
	fmt.Println(db.GetDeployment("user1", "id1"))
	fmt.Println(db.DeleteDeployment("id1"))
	fmt.Println(db.GetDeployment("user1", "id1"))

	//output:
	//GET:
	//{     [] []} not found 404
	//{     [] []} not found 404
	//{     [] []} access denied 403
	//{id1    name1 [] []} <nil> 200
	//CHECK:
	//not found 404
	//not found 404
	//access denied 403
	//<nil> 200
	//DELETE:
	//{id1    name1 [] []} <nil> 200
	//<nil>
	//{id1    name1 [] []} <nil> 200
	//<nil>
	//{     [] []} not found 404

}
