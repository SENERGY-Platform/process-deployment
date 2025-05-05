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
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

var DeploymentBson = getBsonFieldObject[messages.DeploymentCommand]()

func init() {
	CreateCollections = append(CreateCollections, func(db *Mongo, config config.Config) error {
		var err error
		collection := db.client.Database(db.config.MongoTable).Collection(db.config.MongoDeploymentCollection)
		err = db.ensureIndex(collection, "deploymentidindex", DeploymentBson.Id, true, true)
		if err != nil {
			debug.PrintStack()
			return err
		}
		err = db.ensureIndex(collection, "deploymentownerindex", DeploymentBson.Owner, true, false)
		if err != nil {
			debug.PrintStack()
			return err
		}
		return nil
	})
}

func (this *Mongo) deploymentsCollection() *mongo.Collection {
	return this.client.Database(this.config.MongoTable).Collection(this.config.MongoDeploymentCollection)
}

func (this *Mongo) CheckDeploymentAccess(user string, deploymentId string) (err error, code int) {
	ctx, _ := getTimeoutContext()
	wrapper := messages.DeploymentCommand{}
	err = this.deploymentsCollection().FindOne(ctx, bson.M{DeploymentBson.Id: deploymentId, NotDeletedFilterKey: NotDeletedFilterValue}).Decode(&wrapper)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("not found"), http.StatusNotFound
	}
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if wrapper.Owner != user {
		return errors.New("access denied"), http.StatusForbidden
	}
	return nil, 200
}

func (this *Mongo) ListDeployments(user string, listOptions model.DeploymentListOptions) (deployments []deploymentmodel.Deployment, err error) {
	ctx, _ := getTimeoutContext()
	dbOptions := options.Find()
	if listOptions.Limit > 0 {
		dbOptions.SetLimit(listOptions.Limit)
	}
	if listOptions.Offset > 0 {
		dbOptions.SetSkip(listOptions.Offset)
	}

	if listOptions.SortBy == "" {
		listOptions.SortBy = DeploymentBson.Id + ".asc"
	}
	sortby := listOptions.SortBy
	sortby = strings.TrimSuffix(sortby, ".asc")
	sortby = strings.TrimSuffix(sortby, ".desc")
	switch sortby {
	case "name":
		sortby = "deployment.name"
	}
	log.Println("sortby", sortby)
	direction := int32(1)
	if strings.HasSuffix(listOptions.SortBy, ".desc") {
		direction = int32(-1)
	}
	dbOptions.SetSort(bson.D{{sortby, direction}})

	c, err := this.deploymentsCollection().Find(ctx, bson.M{DeploymentBson.Owner: user, NotDeletedFilterKey: NotDeletedFilterValue}, dbOptions)
	if err != nil {
		return nil, err
	}
	list := []messages.DeploymentCommand{}
	err = c.All(ctx, &list)
	if err != nil {
		return nil, err
	}
	for _, element := range list {
		if element.Deployment != nil {
			deployments = append(deployments, *element.Deployment)
		}
	}
	return deployments, nil
}

func (this *Mongo) getDeployment(deploymentId string) (deployment messages.DeploymentCommand, err error, code int) {
	ctx, _ := getTimeoutContext()
	err = this.deploymentsCollection().FindOne(ctx, bson.M{DeploymentBson.Id: deploymentId, NotDeletedFilterKey: NotDeletedFilterValue}).Decode(&deployment)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return deployment, errors.New("not found"), http.StatusNotFound
	}
	if err != nil {
		return deployment, err, http.StatusInternalServerError
	}
	return deployment, nil, 200
}

func (this *Mongo) GetDeployment(user string, deploymentId string) (deployment *deploymentmodel.Deployment, err error, code int) {
	depl, err, code := this.getDeployment(deploymentId)
	if err != nil {
		return nil, err, code
	}
	if depl.Owner != user {
		return nil, errors.New("access denied"), http.StatusForbidden
	}
	return depl.Deployment, nil, 200
}

func (this *Mongo) GetDeploymentIds(user string) (deployments []string, err error) {
	ctx, _ := getTimeoutContext()
	cursor, err := this.deploymentsCollection().Find(ctx, bson.M{DeploymentBson.Owner: user, NotDeletedFilterKey: NotDeletedFilterValue})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		deployment := DeploymentCommandIdVersionWrapper{}
		err = cursor.Decode(&deployment)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, deployment.Id)
	}
	err = cursor.Err()
	return
}

var ErrorUnexpectedDeploymentVersion = errors.New("unexpected deployment version")

type DeploymentCommandIdVersionWrapper struct {
	Id      string `json:"id"`
	Owner   string `json:"owner"`
	Version int64  `json:"version"`
}

type DeploymentWithSyncInfo struct {
	messages.DeploymentCommand `bson:",inline"`
	SyncInfo                   `bson:",inline"`
}

func (this *Mongo) SetDeployment(depl messages.DeploymentCommand, syncHandler func(messages.DeploymentCommand) error) error {
	if depl.Deployment.Version != deploymentmodel.CurrentVersion {
		return ErrorUnexpectedDeploymentVersion
	}
	if depl.Id != depl.Deployment.Id {
		return errors.New("deployment id mismatch")
	}
	ctx, _ := getTimeoutContext()
	timestamp := time.Now().Unix()
	collection := this.deploymentsCollection()
	_, err := collection.ReplaceOne(ctx, bson.M{DeploymentBson.Id: depl.Id}, DeploymentWithSyncInfo{
		DeploymentCommand: depl,
		SyncInfo: SyncInfo{
			SyncTodo:          true,
			SyncDelete:        false,
			SyncUnixTimestamp: timestamp,
		},
	}, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	err = syncHandler(depl)
	if err != nil {
		log.Printf("WARNING: error in SetConcept::syncHandler %v, will be retried later\n", err)
		return nil
	}
	err = this.setSynced(ctx, collection, DeploymentBson.Id, depl.Id, timestamp)
	if err != nil {
		log.Printf("WARNING: error in SetConcept::setSynced %v, will be retried later\n", err)
		return nil
	}
	return nil
}

func (this *Mongo) DeleteDeployment(id string, syncDeleteHandler func(messages.DeploymentCommand) error) error {
	ctx, _ := getTimeoutContext()
	old, err, code := this.getDeployment(id)
	if err != nil {
		if code == http.StatusNotFound {
			return nil
		}
		return err
	}
	collection := this.deploymentsCollection()
	err = this.setDeleted(ctx, collection, DeploymentBson.Id, id)
	if err != nil {
		return err
	}
	err = syncDeleteHandler(old)
	if err != nil {
		log.Printf("WARNING: error in RemoveConcept::syncDeleteHandler %v, will be retried later\n", err)
		return nil
	}
	_, err = collection.DeleteOne(ctx, bson.M{DeploymentBson.Id: id})
	if err != nil {
		log.Printf("WARNING: error in RemoveConcept::DeleteOne %v, will be retried later\n", err)
		return nil
	}
	return nil
}

func (this *Mongo) RetryDeploymentSync(lockduration time.Duration, syncDeleteHandler func(messages.DeploymentCommand) error, syncHandler func(messages.DeploymentCommand) error) error {
	collection := this.deploymentsCollection()
	jobs, err := FetchSyncJobs[DeploymentWithSyncInfo](collection, lockduration, FetchSyncJobsDefaultBatchSize)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if job.SyncDelete {
			err = syncDeleteHandler(job.DeploymentCommand)
			if err != nil {
				log.Printf("WARNING: error in RetryConceptSync::syncDeleteHandler %v, will be retried later\n", err)
				continue
			}
			ctx, _ := getTimeoutContext()
			_, err = collection.DeleteOne(ctx, bson.M{DeploymentBson.Id: job.Id})
			if err != nil {
				log.Printf("WARNING: error in RetryConceptSync::DeleteOne %v, will be retried later\n", err)
				continue
			}
		} else if job.SyncTodo {
			err = syncHandler(job.DeploymentCommand)
			if err != nil {
				log.Printf("WARNING: error in RetryConceptSync::syncHandler %v, will be retried later\n", err)
				continue
			}
			ctx, _ := getTimeoutContext()
			err = this.setSynced(ctx, collection, DeploymentBson.Id, job.Id, job.SyncUnixTimestamp)
			if err != nil {
				log.Printf("WARNING: error in RetryConceptSync::setSynced %v, will be retried later\n", err)
				continue
			}
		}
	}
	return nil
}
