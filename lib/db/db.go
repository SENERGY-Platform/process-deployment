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
	"fmt"
	"runtime/debug"
	"time"

	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"go.mongodb.org/mongo-driver/bson"
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

const startupCheckTimeout = 10 * time.Second

var (
	errEmptyDatabase   = errors.New("mongo database name must not be empty")
	errMissingPassword = errors.New("mongo password must not be empty when a mongo user is set")
)

func (f FactoryType) New(ctx context.Context, config config.Config) (result interfaces.Database, err error) {
	if err = validateConfig(config); err != nil {
		return nil, err
	}
	return start(ctx, config, clientOptions(config), startupCheckTimeout)
}

// start disconnects the client on every failure path, so a failed startup leaves nothing connected.
func start(ctx context.Context, config config.Config, opts *options.ClientOptions, timeout time.Duration) (*Mongo, error) {
	client, err := connect(ctx, opts, config.MongoDatabase, timeout)
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

// connect runs listCollections on the service database because Connect is lazy and ping needs no
// authentication; unreachable servers and wrong or missing credentials then fail at startup.
func connect(ctx context.Context, opts *options.ClientOptions, database string, timeout time.Duration) (*mongo.Client, error) {
	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	client, err := mongo.Connect(checkCtx, opts)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	listOpts := options.ListCollections().SetNameOnly(true).SetAuthorizedCollections(true)
	if _, err = client.Database(database).ListCollectionNames(checkCtx, bson.D{}, listOpts); err != nil {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), timeout)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
		return nil, fmt.Errorf("mongo startup check failed: %w", err)
	}
	return client, nil
}

func validateConfig(conf config.Config) error {
	if conf.MongoDatabase == "" {
		return errEmptyDatabase
	}
	if conf.MongoUser != "" && conf.MongoPassword == "" {
		return errMissingPassword
	}
	return nil
}

// clientOptions applies the credentials after the URI so they replace any given in MONGO_URL.
func clientOptions(conf config.Config) *options.ClientOptions {
	opts := options.Client().ApplyURI(conf.MongoUrl)
	if conf.MongoUser != "" {
		opts.SetAuth(options.Credential{
			Username:   conf.MongoUser,
			Password:   conf.MongoPassword,
			AuthSource: conf.MongoAuthSource,
		})
	}
	return opts
}
