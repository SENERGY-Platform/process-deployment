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
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	deploymentmodel2 "github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel/v2"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"runtime/debug"
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
	if err == mongo.ErrNoDocuments {
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

func (this *Mongo) GetDeployment(user string, deploymentId string) (deploymentV1 *deploymentmodel.Deployment, deploymentV2 *deploymentmodel2.Deployment, err error, code int) {
	ctx, _ := getTimeoutContext()
	wrapper := messages.DeploymentCommand{}
	err = this.deploymentsCollection().FindOne(ctx, bson.M{deploymentIdKey: deploymentId}).Decode(&wrapper)
	if err == mongo.ErrNoDocuments {
		return nil, nil, errors.New("not found"), http.StatusNotFound
	}
	if err != nil {
		return nil, nil, err, http.StatusInternalServerError
	}
	if wrapper.Owner != user {
		return nil, nil, errors.New("access denied"), http.StatusForbidden
	}

	return wrapper.Deployment, wrapper.DeploymentV2, nil, 200
}

func (this *Mongo) SetDeployment(id string, owner string, deploymentV1 *deploymentmodel.Deployment, deploymentV2 *deploymentmodel2.Deployment) error {
	ctx, _ := getTimeoutContext()
	_, err := this.deploymentsCollection().ReplaceOne(ctx, bson.M{deploymentIdKey: id}, messages.DeploymentCommand{Id: id, Owner: owner, Deployment: deploymentV1, DeploymentV2: deploymentV2}, options.Replace().SetUpsert(true))
	return err
}
