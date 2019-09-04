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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FactoryType struct{}

var Factory = FactoryType{}

type Mongo struct {
	config config.Config
	client *mongo.Client
}

var CreateCollections = []func(db *Mongo, config config.Config) error{}

func (f FactoryType) New(ctx context.Context, config config.Config) (result interfaces.Database, err error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoUrl))
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		client.Disconnect(context.Background())
	}()
	db := &Mongo{config: config, client: client}
	for _, creators := range CreateCollections {
		err = creators(db, config)
		if err != nil {
			client.Disconnect(context.Background())
			return nil, err
		}
	}
	return db, nil
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

func (this *Mongo) GetDependencies(user string, deploymentId string) (model.Dependencies, error, int) {
	panic("implement me") //TODO
}

func (this *Mongo) SetDependencies(dependencies model.Dependencies) error {
	panic("implement me") //TODO
}

func (this *Mongo) DeleteDependencies(id string) error {
	panic("implement me") //TODO
}
