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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
)

func BpmnToMultitask(task *etree.Element) (result model.Element, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			err = errors.New(fmt.Sprint("Recovered Error: getAbstractTaskParameter() ", r))
		}
	}()
	cmd := model.Command{}
	cmdPayload := task.FindElement(".//camunda:inputParameter[@name='" + model.CAMUNDA_VARIABLES_PAYLOAD + "']")
	err = json.Unmarshal([]byte(cmdPayload.Text()), &cmd)
	if err != nil {
		return result, err
	}

	documentation := model.Documentation{}
	documentations := task.FindElements(".//bpmn:documentation")
	if len(documentations) > 0 {
		err = json.Unmarshal([]byte(documentations[0].Text()), &documentation)
		if err != nil {
			return result, err
		}
	}

	parameter, err := BpmnToParameter(task)

	id := task.SelectAttr("id").Value
	label := task.SelectAttrValue("name", id)

	result = model.Element{
		Order: documentation.Order,
		MultiTask: &model.MultiTask{
			Label: label,
			DeviceDescription: model.DeviceDescription{
				CharacteristicId: cmd.CharacteristicId,
				Function:         cmd.Function,
				DeviceClass:      cmd.DeviceClass,
				Aspect:           cmd.Aspect,
			},
			BpmnElementId: id,
			Input:         cmd.Input,
			Parameter:     parameter,
		},
	}
	return result, nil
}

func BpmnToTask(task *etree.Element) (result model.Element, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			err = errors.New(fmt.Sprint("Recovered Error: getAbstractTaskParameter() ", r))
		}
	}()
	cmd := model.Command{}
	cmdPayload := task.FindElement(".//camunda:inputParameter[@name='" + model.CAMUNDA_VARIABLES_PAYLOAD + "']")
	err = json.Unmarshal([]byte(cmdPayload.Text()), &cmd)
	if err != nil {
		return result, err
	}

	documentation := model.Documentation{}
	documentations := task.FindElements(".//bpmn:documentation")
	if len(documentations) > 0 {
		err = json.Unmarshal([]byte(documentations[0].Text()), &documentation)
		if err != nil {
			return result, err
		}
	}

	parameter, err := BpmnToParameter(task)

	id := task.SelectAttr("id").Value
	label := task.SelectAttrValue("name", id)

	result = model.Element{
		Order: documentation.Order,
		Task: &model.Task{
			Label: label,
			DeviceDescription: model.DeviceDescription{
				CharacteristicId: cmd.CharacteristicId,
				Function:         cmd.Function,
				DeviceClass:      cmd.DeviceClass,
				Aspect:           cmd.Aspect,
			},
			BpmnElementId: id,
			Parameter:     parameter,
			Input:         cmd.Input,
		},
	}
	return result, nil
}
