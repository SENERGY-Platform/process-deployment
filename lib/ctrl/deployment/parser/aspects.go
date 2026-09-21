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
	"strings"

	"github.com/beevik/etree"
)

// senergy:aspect names one aspect and senergy:aspects a list of them. The single attribute is
// deprecated and an alias for a list with one element; both are read as they were written and
// folded by whoever evaluates the criteria, so an event modeled before the lists keeps its
// results.
const (
	aspectAttr     = "aspect"
	aspectListAttr = "aspects"
)

// hasAspect reports whether an element names an aspect in either spelling. An event without
// one is not recognized as an event at all, so this gate has to accept the list on its own.
func hasAspect(element *etree.Element) bool {
	return selectAspectId(element) != "" || len(selectAspectIds(element)) > 0
}

// selectAspectId returns the deprecated single aspect of an element, empty if it names none.
func selectAspectId(element *etree.Element) string {
	return strings.TrimSpace(element.SelectAttrValue(aspectAttr, ""))
}

// selectAspectIds returns the aspect list of an element. The attribute holds the ids separated
// by commas, the way ids are passed to the device-repository; an aspect urn contains none.
func selectAspectIds(element *etree.Element) (result []string) {
	for _, aspectId := range strings.Split(element.SelectAttrValue(aspectListAttr, ""), ",") {
		aspectId = strings.TrimSpace(aspectId)
		if aspectId != "" {
			result = append(result, aspectId)
		}
	}
	return result
}
