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

func ExampleSaveDependencies() {
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

	err = db.SetDependencies(model.Dependencies{
		DeploymentId: "id1",
		Owner:        "user1",
		Devices: []model.DeviceDependency{{
			DeviceId:      "d1",
			Name:          "d1",
			BpmnResources: []model.BpmnResource{{Id: "r1"}, {Id: "r2"}},
		}},
		Events: nil,
	})
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	err = db.SetDependencies(model.Dependencies{
		DeploymentId: "id2",
		Owner:        "user1",
		Devices: []model.DeviceDependency{{
			DeviceId:      "d1",
			Name:          "d1",
			BpmnResources: []model.BpmnResource{{Id: "r1"}, {Id: "r2"}},
		}},
		Events: nil,
	})
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	err = db.SetDependencies(model.Dependencies{
		DeploymentId: "id3",
		Owner:        "user2",
		Devices: []model.DeviceDependency{{
			DeviceId:      "d1",
			Name:          "d1",
			BpmnResources: []model.BpmnResource{{Id: "r1"}, {Id: "r2"}},
		}},
		Events: nil,
	})
	if err != nil {
		debug.PrintStack()
		fmt.Println(err)
		return
	}

	fmt.Println("GET:")
	fmt.Println(db.GetDependencies("nope", "nope"))
	fmt.Println(db.GetDependencies("user1", "nope"))
	fmt.Println(db.GetDependencies("nope", "id1"))
	fmt.Println(db.GetDependencies("user1", "id1"))

	fmt.Println("LIST:")
	fmt.Println(db.GetDependenciesList("nope", 100, 0))
	fmt.Println(db.GetDependenciesList("user1", 100, 0))

	fmt.Println("SELECT:")
	fmt.Println(db.GetSelectedDependencies("nope", []string{"nope"}))
	fmt.Println(db.GetSelectedDependencies("nope", []string{"nope", "id1"}))
	fmt.Println(db.GetSelectedDependencies("nope", []string{"id2", "id1"}))
	fmt.Println(db.GetSelectedDependencies("user1", []string{"nope", "id1"}))
	fmt.Println(db.GetSelectedDependencies("user1", []string{"id2", "id1"}))

	fmt.Println("DELETE:")
	fmt.Println(db.GetDependencies("user1", "id1"))
	fmt.Println(db.DeleteDependencies("nope"))
	fmt.Println(db.GetDependencies("user1", "id1"))
	fmt.Println(db.DeleteDependencies("id1"))
	fmt.Println(db.GetDependencies("user1", "id1"))

	//output:
	//GET:
	//{  [] []} not found 404
	//{  [] []} not found 404
	//{  [] []} access denied 403
	//{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} <nil> 200
	//LIST:
	//[] <nil> 200
	//[{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} {id2 user1 [{d1 d1 [{r1 } {r2 }]}] []}] <nil> 200
	//SELECT:
	//[] <nil> 200
	//[] <nil> 200
	//[] <nil> 200
	//[{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []}] <nil> 200
	//[{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} {id2 user1 [{d1 d1 [{r1 } {r2 }]}] []}] <nil> 200
	//DELETE:
	//{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} <nil> 200
	//<nil>
	//{id1 user1 [{d1 d1 [{r1 } {r2 }]}] []} <nil> 200
	//<nil>
	//{  [] []} not found 404

}
