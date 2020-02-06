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
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"strings"
)

func (this *Ctrl) completeEvents(deployment *model.Deployment) error {
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

func (this *Ctrl) setDeploymentEventIds(deployment *model.Deployment) error {
	for _, lane := range deployment.Lanes {
		if lane.Lane != nil {
			for _, element := range lane.Lane.Elements {
				if element.MsgEvent != nil {
					element.MsgEvent.EventId = config.NewId()
				}
				if element.ReceiveTaskEvent != nil {
					element.ReceiveTaskEvent.EventId = config.NewId()
				}
			}
		}
		if lane.MultiLane != nil {
			for _, element := range lane.MultiLane.Elements {
				if element.MsgEvent != nil {
					element.MsgEvent.EventId = config.NewId()
				}
				if element.ReceiveTaskEvent != nil {
					element.ReceiveTaskEvent.EventId = config.NewId()
				}
			}
		}
	}

	for _, element := range deployment.Elements {
		if element.MsgEvent != nil {
			element.MsgEvent.EventId = config.NewId()
		}
		if element.ReceiveTaskEvent != nil {
			element.ReceiveTaskEvent.EventId = config.NewId()
		}
	}
	return nil
}

func (this *Ctrl) completeDeploymentEventCasts(deployment *model.Deployment) (err error) {
	for _, lane := range deployment.Lanes {
		if lane.Lane != nil {
			for _, element := range lane.Lane.Elements {
				if element.MsgEvent != nil {
					err = this.completeEventCast(element.MsgEvent)
				}
				if element.ReceiveTaskEvent != nil {
					err = this.completeEventCast(element.ReceiveTaskEvent)
				}
				if err != nil {
					return err
				}
			}
		}
		if lane.MultiLane != nil {
			for _, element := range lane.MultiLane.Elements {
				if element.MsgEvent != nil {
					err = this.completeEventCast(element.MsgEvent)
				}
				if element.ReceiveTaskEvent != nil {
					err = this.completeEventCast(element.ReceiveTaskEvent)
				}
				if err != nil {
					return err
				}
			}
		}
	}
	for _, element := range deployment.Elements {
		if element.MsgEvent != nil {
			err = this.completeEventCast(element.MsgEvent)
		}
		if element.ReceiveTaskEvent != nil {
			err = this.completeEventCast(element.ReceiveTaskEvent)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Ctrl) completeEventCast(event *model.MsgEvent) (err error) {
	if event.TriggerCast != nil && event.TriggerCast.To != "" {
		event.TriggerCast.From, err = getCharacteristicOfPathInService(event.Service, event.Path)
	}
	return
}

func getCharacteristicOfPathInService(service devicemodel.Service, path string) (string, error) {
	pathSegments := strings.Split(path, ".")
	for _, output := range service.Outputs {
		if output.ContentVariable.Name == pathSegments[0] {
			return getCharacteristicOfPathInVariable(output.ContentVariable, pathSegments)
		}
	}
	return "", nil
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
