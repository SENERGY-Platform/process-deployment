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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"runtime/debug"
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
		debug.PrintStack()
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
