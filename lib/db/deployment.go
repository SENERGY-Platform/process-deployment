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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"go.mongodb.org/mongo-driver/mongo"
	"runtime/debug"
)

const deploymentIdFiledName = "Id"
const deploymentOwnerFiledName = "Owner"

var deploymentIdKey string
var deploymentOwnerKey string

func init() {
	CreateCollections = append(CreateCollections, func(db *Mongo, config config.Config) error {
		var err error
		deploymentIdKey, err = getBsonFieldName(model.DeploymentCommand{}, deploymentIdFiledName)
		if err != nil {
			debug.PrintStack()
			return err
		}
		deploymentOwnerKey, err = getBsonFieldName(model.DeploymentCommand{}, deploymentOwnerFiledName)
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

func (this *Mongo) CheckDeploymentAccess(user string, deploymentId string) (error, int) {
	panic("implement me") //TODO
}

func (this *Mongo) GetDeployment(user string, deploymentId string) (model.Deployment, error, int) {
	panic("implement me") //TODO
}

func (this *Mongo) DeleteDeployment(id string) error {
	panic("implement me") //TODO
}

func (this *Mongo) SetDeployment(id string, owner string, deployment model.Deployment) error {
	panic("implement me") //TODO
}
