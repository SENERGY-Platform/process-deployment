/*
 * Copyright 2020 InfAI (CC SES)
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
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"strings"
)

func (this *Ctrl) completeEvents(deployment *deploymentmodel.Deployment) error {
	err := this.setDeploymentEventIds(deployment)
	if err != nil {
		return err
	}
	err = this.completeDeploymentEventCasts(deployment)
	if err != nil {
		return err
	}
	return nil
}

func (this *Ctrl) setDeploymentEventIds(deployment *deploymentmodel.Deployment) error {
	panic("not implemented") //TODO
}

func (this *Ctrl) completeDeploymentEventCasts(deployment *deploymentmodel.Deployment) (err error) {
	panic("not implemented") //TODO
}

func getCharacteristicOfPathInService(service *devicemodel.Service, path string) (string, error) {
	pathSegments := strings.Split(path, ".")
	if len(pathSegments) <= 1 || pathSegments[0] != "value" {
		return "", errors.New("expect 'value' as prefix of msg_vent path")
	}
	pathSegments = pathSegments[1:]
	for _, output := range service.Outputs {
		if output.ContentVariable.Name == pathSegments[0] {
			return getCharacteristicOfPathInVariable(output.ContentVariable, pathSegments)
		}
	}
	return "", errors.New("no characteristic found for " + path + " in " + service.Id)
}

func getCharacteristicOfPathInVariable(variable devicemodel.ContentVariable, path []string) (string, error) {
	for {
		name := path[0]
		rest := path[1:]
		if variable.Name != name {
			return "", errors.New("event path not found: " + strings.Join(path, "."))
		}
		if len(rest) == 0 {
			if variable.CharacteristicId == "" {
				return "", errors.New("path does not reference characteristic" + strings.Join(path, "."))
			} else {
				return variable.CharacteristicId, nil
			}
		} else {
			path = rest
			next := devicemodel.ContentVariable{}
			for _, sub := range variable.SubContentVariables {
				if sub.Name == path[0] {
					next = sub
				}
			}
			variable = next
		}
	}
}
