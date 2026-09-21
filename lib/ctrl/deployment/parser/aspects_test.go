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

package parser

import (
	"reflect"
	"testing"

	"github.com/beevik/etree"
)

func TestSelectAspectIds(t *testing.T) {
	t.Run("reads the ids of a comma separated aspect list", func(t *testing.T) {
		actual := selectAspectIds(catchEvent(t, `senergy:aspects="urn:infai:ses:aspect:a,urn:infai:ses:aspect:b"`))
		expected := []string{"urn:infai:ses:aspect:a", "urn:infai:ses:aspect:b"}
		if !reflect.DeepEqual(actual, expected) {
			t.Error(actual)
		}
	})

	t.Run("keeps the ids in the order they were modeled", func(t *testing.T) {
		actual := selectAspectIds(catchEvent(t, `senergy:aspects="urn:infai:ses:aspect:b,urn:infai:ses:aspect:a"`))
		expected := []string{"urn:infai:ses:aspect:b", "urn:infai:ses:aspect:a"}
		if !reflect.DeepEqual(actual, expected) {
			t.Error(actual)
		}
	})

	t.Run("trims whitespace and drops empty entries", func(t *testing.T) {
		actual := selectAspectIds(catchEvent(t, `senergy:aspects=" urn:infai:ses:aspect:a , ,urn:infai:ses:aspect:b, "`))
		expected := []string{"urn:infai:ses:aspect:a", "urn:infai:ses:aspect:b"}
		if !reflect.DeepEqual(actual, expected) {
			t.Error(actual)
		}
	})

	t.Run("returns no aspects for a missing attribute", func(t *testing.T) {
		if actual := selectAspectIds(catchEvent(t, `senergy:function="urn:infai:ses:measuring-function:f"`)); actual != nil {
			t.Error(actual)
		}
	})
}

func TestHasAspect(t *testing.T) {
	t.Run("recognizes an element that names only the aspect list", func(t *testing.T) {
		if !hasAspect(catchEvent(t, `senergy:aspects="urn:infai:ses:aspect:a"`)) {
			t.Error("expected the list alone to count as an aspect")
		}
	})

	t.Run("recognizes an element that names only the deprecated single aspect", func(t *testing.T) {
		if !hasAspect(catchEvent(t, `senergy:aspect="urn:infai:ses:aspect:a"`)) {
			t.Error("expected the deprecated single aspect to keep counting")
		}
	})

	t.Run("rejects an element that names no aspect", func(t *testing.T) {
		if hasAspect(catchEvent(t, `senergy:aspect="" senergy:aspects=""`)) {
			t.Error("expected empty attributes not to count as an aspect")
		}
	})
}

func catchEvent(t *testing.T, attributes string) *etree.Element {
	t.Helper()
	doc := etree.NewDocument()
	err := doc.ReadFromString(`<bpmn:intermediateCatchEvent xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:senergy="https://senergy.infai.org" id="event_1" ` + attributes + `/>`)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Root()
}
