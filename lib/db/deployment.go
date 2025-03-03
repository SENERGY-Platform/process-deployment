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
)

const deploymentIdFiledName = "Id"
const deploymentOwnerFiledName = "Owner"

var deploymentIdKey string
var deploymentOwnerKey string

func init() {
	CreateCollections = append(CreateCollections, func(db *Mongo, config config.Config) error {
		var err error
		deploymentIdKey, err = getBsonFieldName(messages.DeploymentCommand{}, deploymentIdFiledName)
		if err != nil {
			debug.PrintStack()
			return err
		}
		deploymentOwnerKey, err = getBsonFieldName(messages.DeploymentCommand{}, deploymentOwnerFiledName)
		if err != nil {
			debug.PrintStack()
			return err
		}
		collection := db.client.Database(db.config.MongoTable).Collection(db.config.MongoDeploymentCollection)
		err = db.ensureIndex(collection, "deploymentidindex", deploymentIdKey, true, true)
		if err != nil {
			debug.PrintStack()
			return err
		}
		err = db.ensureIndex(collection, "deploymentownerindex", deploymentOwnerKey, true, false)
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
	err = this.deploymentsCollection().FindOne(ctx, bson.M{deploymentIdKey: deploymentId}).Decode(&wrapper)
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

func (this *Mongo) DeleteDeployment(id string) error {
	ctx, _ := getTimeoutContext()
	_, err := this.deploymentsCollection().DeleteOne(ctx, bson.M{deploymentIdKey: id})
	return err
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
		listOptions.SortBy = deploymentIdFiledName + ".asc"
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

	c, err := this.deploymentsCollection().Find(ctx, bson.M{deploymentOwnerKey: user}, dbOptions)
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

func (this *Mongo) GetDeployment(user string, deploymentId string) (deployment *deploymentmodel.Deployment, err error, code int) {
	ctx, _ := getTimeoutContext()
	response := this.deploymentsCollection().FindOne(ctx, bson.M{deploymentIdKey: deploymentId})

	//first decode to check version to prevent crash in second decode

	versionWrapper := DeploymentCommandIdVersionWrapper{}
	err = response.Decode(&versionWrapper)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("not found"), http.StatusNotFound
	}
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if versionWrapper.Owner != user {
		return nil, errors.New("access denied"), http.StatusForbidden
	}

	wrapper := messages.DeploymentCommand{}
	err = response.Decode(&wrapper)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return wrapper.Deployment, nil, 200
}

func (this *Mongo) GetDeploymentIds(user string) (deployments []string, err error) {
	ctx, _ := getTimeoutContext()
	cursor, err := this.deploymentsCollection().Find(ctx, bson.M{deploymentOwnerKey: user})
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

func (this *Mongo) SetDeployment(id string, owner string, deployment *deploymentmodel.Deployment) error {
	if deployment.Version != deploymentmodel.CurrentVersion {
		return ErrorUnexpectedDeploymentVersion
	}
	ctx, _ := getTimeoutContext()
	_, err := this.deploymentsCollection().ReplaceOne(ctx, bson.M{deploymentIdKey: id}, messages.DeploymentCommand{Id: id, Owner: owner, Deployment: deployment, Version: deployment.Version}, options.Replace().SetUpsert(true))
	return err
}

type DeploymentCommandIdVersionWrapper struct {
	Id      string `json:"id"`
	Owner   string `json:"owner"`
	Version int64  `json:"version"`
}
