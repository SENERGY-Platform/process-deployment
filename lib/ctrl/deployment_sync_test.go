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

package ctrl

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	engine "github.com/SENERGY-Platform/camunda-engine-wrapper/lib/client"
	eventdeployment "github.com/SENERGY-Platform/event-deployment/lib/client"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/db"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/mocks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

// countingEngine fails its first Deploy call and counts all calls.
type countingEngine struct {
	mux     sync.Mutex
	deploys int
	deletes int
}

func (this *countingEngine) Deploy(token string, depl engine.DeploymentMessage) (err error, code int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.deploys++
	if this.deploys == 1 {
		return errors.New("engine not reachable"), http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (this *countingEngine) DeleteDeployment(token string, userId string, deplId string) (err error, code int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.deletes++
	return nil, http.StatusOK
}

type countingEventDepl struct {
	mocks.MockEventDepl
	mux     sync.Mutex
	deploys int
	deletes int
}

func (this *countingEventDepl) Deploy(token string, depl eventdeployment.Deployment) (err error, code int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.deploys++
	return nil, http.StatusOK
}

func (this *countingEventDepl) DeleteDeployment(token string, userId string, deplId string) (err error, code int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.deletes++
	return nil, http.StatusOK
}

// A deployment whose first sync fails must be deployed again by the retry, not deleted.
func TestSyncRetryDeploysFailedDeployment(t *testing.T) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		t.Fatal(err)
	}

	mongoPort, _, err := docker.Mongo(ctx, &wg)
	if err != nil {
		t.Fatal(err)
	}
	conf.MongoUrl = "mongodb://localhost:" + mongoPort

	database, err := db.Factory.New(ctx, conf)
	if err != nil {
		t.Fatal(err)
	}

	eng := &countingEngine{}
	eventDepl := &countingEventDepl{}
	c := &Ctrl{
		config:          conf,
		db:              database,
		devices:         &mocks.DeviceRepoMock{},
		engine:          eng,
		eventdeployment: eventDepl,
	}
	c.deploymentPublisher, err = mocks.Kafka.NewDeploymentProducer(ctx, conf)
	if err != nil {
		t.Fatal(err)
	}

	deployment := deploymentmodel.Deployment{
		Id:      "deployment-id",
		Name:    "name",
		Version: deploymentmodel.CurrentVersion,
		Diagram: deploymentmodel.Diagram{Svg: "<svg/>"},
	}
	err = c.setDeployment("user1", "test", deployment)
	if err != nil {
		t.Fatal(err)
	}
	if eng.deploys != 1 {
		t.Fatalf("expected the first sync to call engine deploy once, got %v", eng.deploys)
	}

	// sync jobs become due once their timestamp is older than the lock duration (unix seconds)
	time.Sleep(2100 * time.Millisecond)
	if err := c.Sync(time.Second); err != nil {
		t.Fatal(err)
	}

	if eng.deploys != 2 {
		t.Errorf("expected the retry to deploy again, engine deploy calls: %v", eng.deploys)
	}
	if eng.deletes != 0 || eventDepl.deletes != 0 {
		t.Errorf("retry must not delete the deployment, engine deletes: %v, event deletes: %v", eng.deletes, eventDepl.deletes)
	}
	if eventDepl.deploys != 1 {
		t.Errorf("expected one event deploy call, got %v", eventDepl.deploys)
	}

	_, err, code := database.GetDeployment("user1", deployment.Id)
	if err != nil {
		t.Fatal(err, code)
	}

	client, err := mongo.Connect(ctx, mongooptions.Client().ApplyURI(conf.MongoUrl))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	stored := db.DeploymentWithSyncInfo{}
	err = client.Database(conf.MongoDatabase).Collection(conf.MongoDeploymentCollection).FindOne(ctx, bson.M{db.DeploymentBson.Id: deployment.Id}).Decode(&stored)
	if err != nil {
		t.Fatal(err)
	}
	if stored.SyncTodo || stored.SyncDelete {
		t.Errorf("expected deployment to be synced and not deleted, got %+v", stored.SyncInfo)
	}
}
