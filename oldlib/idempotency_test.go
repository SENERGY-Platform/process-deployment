/*
 * Copyright 2018 InfAI (CC SES)
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

package oldlib

import (
	"github.com/SENERGY-Platform/process-deployment/oldlib/model"
	"github.com/SENERGY-Platform/process-deployment/oldlib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/ory/dockertest"
	"gopkg.in/mgo.v2"
	"log"
	"testing"
)

const owner = "owner"
const ownerJwt = jwt_http_router.JwtImpersonate("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJvd25lciIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.M33n6BgW1v-RcR0XaI4z288FwnctuijTuaHDIKBnKpI")

func Test(t *testing.T) {
	closer, mongoport, _, err := testHelper_getMongoDependency()
	defer closer()
	if err != nil {
		t.Error(err)
		return
	}

	err = util.LoadConfig("../config.json")
	if err != nil {
		t.Error(err)
		return
	}

	util.Config.MongoUrl = "mongodb://localhost:" + mongoport

	err = testHelper_putProcess("1", "f1")
	if err != nil {
		t.Error(err)
		return
	}

	metadata, err := GetAllMetadata(owner)
	if err != nil {
		t.Error(err)
		return
	}
	if len(metadata) != 1 || metadata[0].Process != "1" || metadata[0].Abstract.MsgEvents[0].FilterId != "f1" {
		t.Error("unexpected result: ", metadata)
		return
	}

	err = testHelper_putProcess("1", "f2")
	if err != nil {
		t.Error(err)
		return
	}
	metadata, err = GetAllMetadata(owner)
	if err != nil {
		t.Error(err)
		return
	}
	if len(metadata) != 1 || metadata[0].Process != "1" || metadata[0].Abstract.MsgEvents[0].FilterId != "f2" {
		t.Error("unexpected result: ", metadata)
		return
	}

	err = testHelper_deleteProcess("1")
	if err != nil {
		t.Error(err)
		return
	}
	metadata, err = GetAllMetadata(owner)
	if err != nil {
		t.Error(err)
		return
	}
	if len(metadata) != 0 {
		t.Error("unexpected result: ", metadata)
		return
	}

	err = testHelper_deleteProcess("1")
	if err != nil {
		t.Error(err)
		return
	}
	metadata, err = GetAllMetadata(owner)
	if err != nil {
		t.Error(err)
		return
	}
	if len(metadata) != 0 {
		t.Error("unexpected result: ", metadata)
		return
	}

}

func testHelper_getMongoDependency() (closer func(), hostPort string, ipAddress string, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return func() {}, "", "", err
	}
	log.Println("start mongodb")
	mongo, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		return func() {}, "", "", err
	}
	hostPort = mongo.GetPort("27017/tcp")
	err = pool.Retry(func() error {
		log.Println("try mongodb connection...")
		sess, err := mgo.Dial("mongodb://localhost:" + hostPort)
		if err != nil {
			return err
		}
		defer sess.Close()
		return sess.Ping()
	})
	return func() { mongo.Close() }, hostPort, mongo.Container.NetworkSettings.IPAddress, err
}

func testHelper_putProcess(vid string, filterId string) error {
	return handleDeploymentMetadataUpdate(DeploymentCommand{
		Id:      vid,
		Command: "PUT",
		Owner:   owner,
		Deployment: model.DeploymentRequest{
			Process: model.AbstractProcess{
				MsgEvents: []model.MsgEvent{
					{
						FilterId: filterId,
					},
				},
			},
		},
	})
}

func testHelper_deleteProcess(vid string) error {
	return handleDeploymentMetadataDelete(DeploymentCommand{
		Id:      vid,
		Command: "DELETE",
		Owner:   owner,
	})
}
