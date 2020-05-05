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
	"github.com/SENERGY-Platform/process-deployment/lib/model/dependencymodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"net/http"
	"runtime/debug"
)

const dependenciesIdFiledName = "DeploymentId"
const dependenciesOwnerFiledName = "Owner"

var dependenciesIdKey string
var dependenciesOwnerKey string

func init() {
	CreateCollections = append(CreateCollections, func(db *Mongo, config config.Config) error {
		var err error
		dependenciesIdKey, err = getBsonFieldName(dependencymodel.Dependencies{}, dependenciesIdFiledName)
		if err != nil {
			debug.PrintStack()
			return err
		}
		dependenciesOwnerKey, err = getBsonFieldName(dependencymodel.Dependencies{}, dependenciesOwnerFiledName)
		if err != nil {
			debug.PrintStack()
			return err
		}
		collection := db.client.Database(db.config.MongoTable).Collection(db.config.MongoDependenciesCollection)
		err = db.ensureIndex(collection, "dependenciesidindex", dependenciesIdKey, true, true)
		if err != nil {
			debug.PrintStack()
			return err
		}
		err = db.ensureIndex(collection, "dependenciesownerindex", dependenciesOwnerKey, true, false)
		if err != nil {
			debug.PrintStack()
			return err
		}
		return nil
	})
}

func (this *Mongo) dependenciesCollection() *mongo.Collection {
	return this.client.Database(this.config.MongoTable).Collection(this.config.MongoDependenciesCollection)
}

func (this *Mongo) SetDependencies(dependencies dependencymodel.Dependencies) error {
	ctx, _ := getTimeoutContext()
	_, err := this.dependenciesCollection().ReplaceOne(ctx, bson.M{dependenciesIdKey: dependencies.DeploymentId}, dependencies, options.Replace().SetUpsert(true))
	return err
}

func (this *Mongo) GetDependencies(user string, deploymentId string) (result dependencymodel.Dependencies, err error, code int) {
	ctx, _ := getTimeoutContext()
	err = this.dependenciesCollection().FindOne(ctx, bson.M{dependenciesIdKey: deploymentId}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return result, errors.New("not found"), http.StatusNotFound
	}
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	if result.Owner != user {
		return dependencymodel.Dependencies{}, errors.New("access denied"), http.StatusForbidden
	}
	return result, nil, 200
}

func (this *Mongo) GetDependenciesList(user string, limit int, offset int) (result []dependencymodel.Dependencies, err error, code int) {
	opt := options.Find()
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(offset))
	opt.SetSort(bsonx.Doc{{deploymentIdKey, bsonx.Int32(1)}})
	ctx, _ := getTimeoutContext()
	cursor, err := this.dependenciesCollection().Find(ctx, bson.M{dependenciesOwnerKey: user}, opt)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer cursor.Close(context.Background())
	for cursor.Next(ctx) {
		dependency := dependencymodel.Dependencies{}
		err = cursor.Decode(&dependency)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		result = append(result, dependency)
	}
	err = cursor.Err()
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return result, nil, 200
}

func (this *Mongo) GetSelectedDependencies(user string, ids []string) (result []dependencymodel.Dependencies, err error, code int) {
	opt := options.Find()
	opt.SetSort(bsonx.Doc{{deploymentIdKey, bsonx.Int32(1)}})
	ctx, _ := getTimeoutContext()
	cursor, err := this.dependenciesCollection().Find(ctx, bson.M{dependenciesIdKey: bson.M{"$in": ids}, dependenciesOwnerKey: user}, opt)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer cursor.Close(context.Background())
	for cursor.Next(ctx) {
		dependency := dependencymodel.Dependencies{}
		err = cursor.Decode(&dependency)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		result = append(result, dependency)
	}
	err = cursor.Err()
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return result, nil, 200
}

func (this *Mongo) DeleteDependencies(id string) error {
	ctx, _ := getTimeoutContext()
	_, err := this.dependenciesCollection().DeleteOne(ctx, bson.M{dependenciesIdKey: id})
	return err
}
