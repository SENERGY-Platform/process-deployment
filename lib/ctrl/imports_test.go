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

	dsmodel "github.com/SENERGY-Platform/device-selection/v2/pkg/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deviceselectionmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/importmodel"
)

// The selection answer and the deployment model name the import shapes through different
// packages, and both have to resolve to the same type. A cast layer between them drops every
// field it does not name - that is how an import type's Cost and the aspect list of a content
// variable went missing before. This assertion is the guard: if either side stops aliasing
// models/go, this stops compiling instead of quietly narrowing the answer.
var _ = func(selectable dsmodel.Selectable) (*importmodel.Import, *importmodel.ImportType) {
	return selectable.Import, selectable.ImportType
}

var _ = func(option deploymentmodel.SelectionOption) (*importmodel.Import, *importmodel.ImportType) {
	return option.Import, option.ImportType
}

// TestImportSelectableKeepsEveryField pins the behaviour the assertions above only make
// possible: an import selectable reaches the deployment option whole.
func TestImportSelectableKeepsEveryField(t *testing.T) {
	restart := true
	selectedImport := &importmodel.Import{
		Id:           "urn:infai:ses:import:1",
		Name:         "import",
		ImportTypeId: "urn:infai:ses:import-type:1",
		Image:        "image",
		KafkaTopic:   "topic",
		Restart:      &restart,
		Configs:      []importmodel.ImportConfig{{Name: "config", Value: "set"}},
	}
	selectedType := &importmodel.ImportType{
		Id:             "urn:infai:ses:import-type:1",
		Name:           "type",
		Description:    "described",
		Image:          "image",
		DefaultRestart: true,
		Owner:          "user1",
		Cost:           42,
		Configs: []importmodel.ImportTypeConfig{{
			Name:         "config",
			Description:  "what it does",
			Type:         "https://schema.org/Text",
			DefaultValue: "fallback",
		}},
		Output: importmodel.ImportContentVariable{
			Name: "root",
			Type: "https://schema.org/StructuredValue",
			SubContentVariables: []importmodel.ImportContentVariable{{
				Name:             "value",
				Type:             "https://schema.org/Float",
				CharacteristicId: "urn:infai:ses:characteristic:1",
				FunctionId:       "urn:infai:ses:measuring-function:1",
				AspectId:         "urn:infai:ses:aspect:a",
				AspectIds:        []string{"urn:infai:ses:aspect:a", "urn:infai:ses:aspect:b"},
			}},
		},
	}

	options := getSelectionOptions(
		[]deviceselectionmodel.Selectable{{Import: selectedImport, ImportType: selectedType}},
		deploymentmodel.FilterCriteria{},
	)

	if len(options) != 1 {
		t.Fatal(options)
	}
	if !reflect.DeepEqual(options[0].Import, selectedImport) {
		t.Error(options[0].Import)
	}
	if !reflect.DeepEqual(options[0].ImportType, selectedType) {
		t.Error(options[0].ImportType)
	}

	//the two fields a narrowing copy lost, named so that a reintroduced one fails loudly
	if options[0].ImportType.Cost != 42 {
		t.Error("the cost of an import type is dropped,", options[0].ImportType.Cost)
	}
	nested := options[0].ImportType.Output.SubContentVariables[0]
	if !reflect.DeepEqual(nested.AspectIds, []string{"urn:infai:ses:aspect:a", "urn:infai:ses:aspect:b"}) {
		t.Error("the aspect list of a content variable is dropped,", nested.AspectIds)
	}
}
