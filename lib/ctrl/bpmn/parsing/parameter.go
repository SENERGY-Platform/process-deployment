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

package parsing

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"strings"
)

func BpmnToParameter(task *etree.Element) (result map[string]string, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			err = errors.New(fmt.Sprint("Recovered Error: getAbstractTaskParameter() ", r))
		}
	}()
	result = map[string]string{}
	for _, input := range task.FindElements(".//camunda:inputParameter") {
		name := input.SelectAttr("name").Value
		value := input.Text()
		if strings.HasPrefix(name, "inputs") {
			result[name] = value
		}
	}
	return result, err
}
