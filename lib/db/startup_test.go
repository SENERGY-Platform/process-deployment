/*
 * Copyright 2026 InfAI (CC SES)
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
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const replicaSetURL = "mongodb://mongo-0.mongo:27017,mongo-1.mongo:27017/?replicaSet=rs0&readPreference=primary"

func TestClientOptions_AuthWhenUserGiven(t *testing.T) {
	opts := clientOptions(&config.ConfigStruct{
		MongoUrl:        replicaSetURL,
		MongoUser:       "process-deployment",
		MongoPassword:   "s3cr3t",
		MongoAuthSource: "admin",
		MongoDatabase:   "process_deployment",
	})
	if err := opts.Validate(); err != nil {
		t.Fatal(err)
	}
	want := &options.Credential{Username: "process-deployment", Password: "s3cr3t", AuthSource: "admin"}
	if !reflect.DeepEqual(opts.Auth, want) {
		t.Errorf("auth = %+v, want %+v", opts.Auth, want)
	}
}

func TestClientOptions_NoAuthWhenUserEmpty(t *testing.T) {
	// A password without a user must not switch auth on.
	opts := clientOptions(&config.ConfigStruct{
		MongoUrl:        "mongodb://localhost:27017",
		MongoPassword:   "s3cr3t",
		MongoAuthSource: "admin",
		MongoDatabase:   "process_deployment",
	})
	if err := opts.Validate(); err != nil {
		t.Fatal(err)
	}
	if opts.Auth != nil {
		t.Errorf("auth = %+v, want nil", opts.Auth)
	}
}

func TestClientOptions_ConfiguredCredentialsReplaceURICredentials(t *testing.T) {
	opts := clientOptions(&config.ConfigStruct{
		MongoUrl:        "mongodb://old:oldpw@localhost:27017/?authSource=other&authMechanism=SCRAM-SHA-1",
		MongoUser:       "process-deployment",
		MongoPassword:   "newpw",
		MongoAuthSource: "admin",
		MongoDatabase:   "process_deployment",
	})
	if err := opts.Validate(); err != nil {
		t.Fatal(err)
	}
	want := &options.Credential{Username: "process-deployment", Password: "newpw", AuthSource: "admin"}
	if !reflect.DeepEqual(opts.Auth, want) {
		t.Errorf("auth = %+v, want %+v", opts.Auth, want)
	}
}

func TestClientOptions_URIPassedUnchanged(t *testing.T) {
	opts := clientOptions(&config.ConfigStruct{MongoUrl: replicaSetURL, MongoDatabase: "process_deployment"})
	if err := opts.Validate(); err != nil {
		t.Fatal(err)
	}
	if got := opts.GetURI(); got != replicaSetURL {
		t.Errorf("uri = %q, want %q", got, replicaSetURL)
	}
	if want := []string{"mongo-0.mongo:27017", "mongo-1.mongo:27017"}; !reflect.DeepEqual(opts.Hosts, want) {
		t.Errorf("hosts = %v, want %v", opts.Hosts, want)
	}
	if opts.ReplicaSet == nil || *opts.ReplicaSet != "rs0" {
		t.Errorf("replica set = %v, want rs0", opts.ReplicaSet)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.ConfigStruct
		wantErr error
	}{
		{"no auth", &config.ConfigStruct{MongoDatabase: "process_deployment"}, nil},
		{"user and password", &config.ConfigStruct{MongoDatabase: "process_deployment", MongoUser: "u", MongoPassword: "p"}, nil},
		{"password without user", &config.ConfigStruct{MongoDatabase: "process_deployment", MongoPassword: "p"}, nil},
		{"user without password", &config.ConfigStruct{MongoDatabase: "process_deployment", MongoUser: "u"}, errMissingPassword},
		{"empty database", &config.ConfigStruct{}, errEmptyDatabase},
		{"empty database with credentials", &config.ConfigStruct{MongoUser: "u", MongoPassword: "p"}, errEmptyDatabase},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateConfig(tt.cfg); !errors.Is(err, tt.wantErr) {
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// unreachableURL points at a port that was just free, so only the startup check can fail.
func unreachableURL(t *testing.T) string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	if err = l.Close(); err != nil {
		t.Fatal(err)
	}
	return "mongodb://" + addr + "/?directConnection=true"
}

// The startup check would fail as well, so these check for the specific validation error.
func TestNew_RejectsBeforeConnecting(t *testing.T) {
	cases := []struct {
		name string
		cfg  *config.ConfigStruct
		want error
	}{
		{"empty database", &config.ConfigStruct{MongoUser: "process-deployment", MongoPassword: "s3cr3t"}, errEmptyDatabase},
		{"user without password", &config.ConfigStruct{MongoUser: "process-deployment", MongoDatabase: "process_deployment"}, errMissingPassword},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.cfg.MongoUrl = unreachableURL(t)
			begin := time.Now()
			result, err := Factory.New(context.Background(), c.cfg)
			if !errors.Is(err, c.want) {
				t.Fatalf("err = %v, want %v", err, c.want)
			}
			if result != nil {
				t.Error("expected no db on failure")
			}
			if strings.Contains(err.Error(), "s3cr3t") {
				t.Errorf("error leaks the password: %v", err)
			}
			if elapsed := time.Since(begin); elapsed > time.Second {
				t.Errorf("validation took %v, it should not contact the server", elapsed)
			}
		})
	}
}

// TestStart_StartupCheckFailsWithoutServer runs in -short mode too: it needs no docker, only an
// unreachable port, so it still catches a startup check that was removed or stopped disconnecting.
func TestStart_StartupCheckFailsWithoutServer(t *testing.T) {
	const password = "pw-must-not-appear-7f3a"
	conf := &config.ConfigStruct{
		MongoUrl:                    unreachableURL(t),
		MongoUser:                   "process-deployment",
		MongoPassword:               password,
		MongoAuthSource:             "admin",
		MongoDatabase:               "process_deployment",
		MongoDeploymentCollection:   "deployments",
		MongoDependenciesCollection: "dependencies",
	}
	pools := &poolCounter{}
	begin := time.Now()
	mongoDb, err := start(context.Background(), conf, clientOptions(conf).SetPoolMonitor(pools.monitor()), 500*time.Millisecond)
	if err == nil {
		t.Fatal("expected an error when the server is unreachable")
	}
	if mongoDb != nil {
		t.Error("expected no db on failure")
	}
	if !strings.HasPrefix(err.Error(), "mongo startup check failed: ") {
		t.Errorf("unexpected error: %v", err)
	}
	if strings.Contains(err.Error(), password) {
		t.Error("error text contains the password")
	}
	if elapsed := time.Since(begin); elapsed > 5*time.Second {
		t.Errorf("start took %v, the timeout was not applied", elapsed)
	}
	pools.assertAllClosed(t)
}
