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
	"reflect"
	"sync"
	"testing"

	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"github.com/SENERGY-Platform/process-deployment/lib/tests/docker"
)

// The stored deployment is what the event side is served from, and
// models.ProcessFilterCriteria carries json tags only — mongo therefore names its fields by
// the driver's own rule rather than by the wire spelling. That is symmetric and survives, but
// it is not visible from the struct, so it is asserted rather than reasoned about.
func TestStoredDeploymentKeepsTheAspectList(t *testing.T) {
	if testing.Short() {
		t.Skip("short tests only without docker")
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		t.Fatal(err)
	}
	port, _, err := docker.Mongo(ctx, &wg)
	if err != nil {
		t.Fatal(err)
	}
	conf.MongoUrl = "mongodb://localhost:" + port

	db, err := Factory.New(ctx, conf)
	if err != nil {
		t.Fatal(err)
	}

	aspectId := "urn:infai:ses:aspect:deprecated"
	stored := &deploymentmodel.Deployment{
		Id:      "aspect-list",
		Name:    "aspect-list",
		Version: models.CurrentDeploymentModelVersion,
		Elements: []deploymentmodel.Element{
			{
				BpmnId: "conditional",
				ConditionalEvent: &deploymentmodel.ConditionalEvent{
					Selection: deploymentmodel.Selection{
						FilterCriteria: deploymentmodel.FilterCriteria{
							AspectIds: []string{"urn:infai:ses:aspect:a", "urn:infai:ses:aspect:b"},
						},
					},
				},
			},
			{
				BpmnId: "task",
				Task: &deploymentmodel.Task{
					Selection: deploymentmodel.Selection{
						FilterCriteria: deploymentmodel.FilterCriteria{AspectId: &aspectId},
					},
				},
			},
		},
	}

	err = db.SetDeployment(messages.DeploymentCommand{
		Command:    "PUT",
		Id:         stored.Id,
		Owner:      "user1",
		Deployment: stored,
		Source:     "test",
		Version:    models.CurrentDeploymentModelVersion,
	}, func(command messages.DeploymentCommand) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	read, err, _ := db.GetDeployment("user1", stored.Id)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("reads back the aspect list of a conditional event", func(t *testing.T) {
		actual := read.Elements[0].ConditionalEvent.Selection.FilterCriteria.AspectIds
		expected := []string{"urn:infai:ses:aspect:a", "urn:infai:ses:aspect:b"}
		if !reflect.DeepEqual(actual, expected) {
			t.Error(actual)
		}
	})

	t.Run("reads back the deprecated single aspect of a task without inventing a list", func(t *testing.T) {
		criteria := read.Elements[1].Task.Selection.FilterCriteria
		if criteria.AspectId == nil || *criteria.AspectId != "urn:infai:ses:aspect:deprecated" {
			t.Error(criteria.AspectId)
		}
		if len(criteria.AspectIds) != 0 {
			t.Error("expected storage not to fold the deprecated aspect into a list,", criteria.AspectIds)
		}
	})
}
