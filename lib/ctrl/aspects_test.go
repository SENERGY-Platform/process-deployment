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

package ctrl

import (
	"reflect"
	"testing"

	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"
)

const (
	parentAspect  = "urn:infai:ses:aspect:parent"
	childAspect   = "urn:infai:ses:aspect:child"
	siblingAspect = "urn:infai:ses:aspect:sibling"
	measuringFunc = "urn:infai:ses:measuring-function:f"
	serviceId     = "urn:infai:ses:service:s"
)

func TestCriteriaAspectIds(t *testing.T) {
	t.Run("folds the deprecated single aspect into the list", func(t *testing.T) {
		actual := criteriaAspectIds(criteria(aspectPtr(parentAspect), nil))
		if !reflect.DeepEqual(actual, []string{parentAspect}) {
			t.Error(actual)
		}
	})

	t.Run("does not repeat an aspect the list already contains", func(t *testing.T) {
		actual := criteriaAspectIds(criteria(aspectPtr(parentAspect), []string{parentAspect, childAspect}))
		if !reflect.DeepEqual(actual, []string{parentAspect, childAspect}) {
			t.Error(actual)
		}
	})

	t.Run("treats an empty single aspect as no aspect", func(t *testing.T) {
		if actual := criteriaAspectIds(criteria(aspectPtr(""), nil)); len(actual) != 0 {
			t.Error(actual)
		}
	})

	t.Run("leaves the aspect list of the criteria untouched", func(t *testing.T) {
		given := criteria(aspectPtr(childAspect), []string{parentAspect})
		criteriaAspectIds(given)
		if !reflect.DeepEqual(given.AspectIds, []string{parentAspect}) {
			t.Error(given.AspectIds)
		}
	})
}

func TestServiceMatchesCriteria(t *testing.T) {
	child := devicemodel.AspectNode{Id: childAspect, AncestorIds: []string{parentAspect}}
	parent := devicemodel.AspectNode{Id: parentAspect, DescendentIds: []string{childAspect}}

	t.Run("matches a path option carrying the aspect the criteria names", func(t *testing.T) {
		if !matches(criteria(nil, []string{childAspect}), child) {
			t.Error("expected the aspect of the path option to match")
		}
	})

	t.Run("matches a path option carrying a descendant of the named aspect", func(t *testing.T) {
		if !matches(criteria(nil, []string{parentAspect}), child) {
			t.Error("expected a queried aspect to cover its subtree")
		}
	})

	t.Run("rejects a path option carrying only an ancestor of the named aspect", func(t *testing.T) {
		if matches(criteria(nil, []string{childAspect}), parent) {
			t.Error("expected the subtree coverage to work downwards only")
		}
	})

	t.Run("requires every aspect of a criteria on the same path option", func(t *testing.T) {
		if !matches(criteria(nil, []string{parentAspect, childAspect}), child) {
			t.Error("expected one node to satisfy an aspect and its ancestor")
		}
		if matches(criteria(nil, []string{childAspect, siblingAspect}), child) {
			t.Error("expected a criteria naming two unrelated aspects to be an AND")
		}
	})

	t.Run("satisfies a two aspect criteria from the aspect node list of one path option", func(t *testing.T) {
		sibling := devicemodel.AspectNode{Id: siblingAspect, AncestorIds: []string{parentAspect}}
		option := deviceselectionmodel.PathOption{
			FunctionId:  measuringFunc,
			AspectNode:  child,
			AspectNodes: []devicemodel.AspectNode{child, sibling},
		}
		if !serviceMatchesCriteria(devicemodel.Service{Id: serviceId}, criteria(nil, []string{childAspect, siblingAspect}), options(option)) {
			t.Error("expected the aspect node list to satisfy both aspects")
		}
	})

	t.Run("prefers the aspect node list over the deprecated single node", func(t *testing.T) {
		option := deviceselectionmodel.PathOption{
			FunctionId:  measuringFunc,
			AspectNode:  child,
			AspectNodes: []devicemodel.AspectNode{{Id: siblingAspect}},
		}
		if matches := serviceMatchesCriteria(devicemodel.Service{Id: serviceId}, criteria(nil, []string{childAspect}), options(option)); matches {
			t.Error("expected the deprecated single node to be ignored while the list is filled")
		}
	})

	t.Run("accepts the deprecated single aspect like a one element list", func(t *testing.T) {
		if !matches(criteria(aspectPtr(parentAspect), nil), child) {
			t.Error("expected the deprecated single aspect to keep matching")
		}
		if matches(criteria(aspectPtr(siblingAspect), nil), child) {
			t.Error("expected an unrelated single aspect not to match")
		}
	})

	t.Run("treats a criteria without an aspect as unfiltered by aspect", func(t *testing.T) {
		if !matches(criteria(nil, nil), devicemodel.AspectNode{}) {
			t.Error("expected an unset aspect not to be a filter")
		}
	})

	t.Run("rejects a service the selectable offers no path option for", func(t *testing.T) {
		service := devicemodel.Service{Id: "urn:infai:ses:service:other"}
		other := deviceselectionmodel.PathOption{FunctionId: measuringFunc, AspectNode: child}
		if serviceMatchesCriteria(service, criteria(nil, []string{childAspect}), options(other)) {
			t.Error("expected a service without path options to be rejected")
		}
	})
}

func aspectPtr(aspectId string) *string {
	return &aspectId
}

func criteria(aspectId *string, aspectIds []string) deploymentmodel.FilterCriteria {
	functionId := measuringFunc
	return deploymentmodel.FilterCriteria{
		FunctionId: &functionId,
		AspectId:   aspectId,
		AspectIds:  aspectIds,
	}
}

// matches asks whether a service satisfies the criteria when its one path option names the
// given aspect the way a device-selection that predates the aspect list does: in the
// deprecated single field, with no list next to it.
func matches(c deploymentmodel.FilterCriteria, aspectNode devicemodel.AspectNode) bool {
	service := devicemodel.Service{Id: serviceId}
	return serviceMatchesCriteria(service, c, options(deviceselectionmodel.PathOption{
		FunctionId: measuringFunc,
		AspectNode: aspectNode,
	}))
}

func options(pathOptions ...deviceselectionmodel.PathOption) map[string][]deviceselectionmodel.PathOption {
	return map[string][]deviceselectionmodel.PathOption{serviceId: pathOptions}
}

// The bulk request is the one place on the aspect path where the criteria is copied field by
// field, into the shape the selection service is asked with. A field that is not copied there
// is dropped silently: the request stays valid and simply asks for less.
func TestBulkSelectableRequestCarriesTheAspectList(t *testing.T) {
	conf, err := config.LoadConfig("../../config.json")
	if err != nil {
		t.Fatal(err)
	}
	ctrl := &Ctrl{config: conf}

	t.Run("forwards the aspect list of every element kind", func(t *testing.T) {
		//the expectation is written out rather than taken from the input, or mutating the
		//input would move both sides and the assertion could not fail
		deployment := &deploymentmodel.Deployment{Elements: []deploymentmodel.Element{
			{BpmnId: "task", Task: &deploymentmodel.Task{Selection: selection(criteria(nil, []string{parentAspect, childAspect}))}},
			{BpmnId: "message", MessageEvent: &deploymentmodel.MessageEvent{Selection: selection(criteria(nil, []string{parentAspect, childAspect}))}},
			{BpmnId: "conditional", ConditionalEvent: &deploymentmodel.ConditionalEvent{Selection: selection(criteria(nil, []string{parentAspect, childAspect}))}},
		}}
		expected := []string{"urn:infai:ses:aspect:parent", "urn:infai:ses:aspect:child"}
		for id, criteria := range requestedCriteria(t, ctrl, deployment) {
			if !reflect.DeepEqual(criteria.AspectIds, expected) {
				t.Error(id, criteria.AspectIds)
			}
		}
	})

	t.Run("forwards the deprecated single aspect unfolded", func(t *testing.T) {
		deployment := &deploymentmodel.Deployment{Elements: []deploymentmodel.Element{
			{BpmnId: "task", Task: &deploymentmodel.Task{Selection: selection(criteria(aspectPtr(parentAspect), nil))}},
		}}
		asked := requestedCriteria(t, ctrl, deployment)["task"]
		if asked.AspectId != parentAspect {
			t.Error(asked.AspectId)
		}
		if len(asked.AspectIds) != 0 {
			t.Error("expected no list next to the deprecated single aspect,", asked.AspectIds)
		}
	})

	t.Run("marks an event criteria with the event interaction", func(t *testing.T) {
		deployment := &deploymentmodel.Deployment{Elements: []deploymentmodel.Element{
			{BpmnId: "task", Task: &deploymentmodel.Task{Selection: selection(criteria(nil, []string{parentAspect}))}},
			{BpmnId: "message", MessageEvent: &deploymentmodel.MessageEvent{Selection: selection(criteria(nil, []string{parentAspect}))}},
			{BpmnId: "conditional", ConditionalEvent: &deploymentmodel.ConditionalEvent{Selection: selection(criteria(nil, []string{parentAspect}))}},
		}}
		asked := requestedCriteria(t, ctrl, deployment)
		if asked["task"].Interaction != "" {
			t.Error("a task is not restricted to an interaction,", asked["task"].Interaction)
		}
		for _, id := range []string{"message", "conditional"} {
			if asked[id].Interaction != string(devicemodel.EVENT) {
				t.Error(id, asked[id].Interaction)
			}
		}
	})

	t.Run("asks a grouped task with the criteria of its whole group", func(t *testing.T) {
		group := "group1"
		deployment := &deploymentmodel.Deployment{Elements: []deploymentmodel.Element{
			{BpmnId: "task1", Group: &group, Task: &deploymentmodel.Task{Selection: selection(criteria(nil, []string{parentAspect}))}},
			{BpmnId: "task2", Group: &group, Task: &deploymentmodel.Task{Selection: selection(criteria(nil, []string{childAspect}))}},
		}}
		for _, element := range ctrl.getDeploymentBulkSelectableRequestV2(deployment) {
			if len(element.Criteria) != 2 {
				t.Fatal(element.Id, len(element.Criteria))
			}
			asked := []string{element.Criteria[0].AspectIds[0], element.Criteria[1].AspectIds[0]}
			if !reflect.DeepEqual(asked, []string{parentAspect, childAspect}) {
				t.Error(element.Id, asked)
			}
		}
	})
}

func selection(c deploymentmodel.FilterCriteria) deploymentmodel.Selection {
	return deploymentmodel.Selection{FilterCriteria: c}
}

// requestedCriteria indexes the bulk request by element, so a missing element fails the test
// rather than skipping its assertion.
func requestedCriteria(t *testing.T, ctrl *Ctrl, deployment *deploymentmodel.Deployment) map[string]deviceselectionmodel.FilterCriteria {
	t.Helper()
	result := map[string]deviceselectionmodel.FilterCriteria{}
	for _, element := range ctrl.getDeploymentBulkSelectableRequestV2(deployment) {
		if len(element.Criteria) != 1 {
			t.Fatalf("%v: expected one criteria, got %v", element.Id, len(element.Criteria))
		}
		result[element.Id] = element.Criteria[0]
	}
	if len(result) != len(deployment.Elements) {
		t.Fatalf("expected one request element per deployment element, got %v", len(result))
	}
	return result
}
