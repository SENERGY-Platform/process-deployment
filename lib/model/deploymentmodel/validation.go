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

package deploymentmodel

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
	"strings"
)

//strict for cqrs; else for user
func (this Deployment) Validate(strict bool) (err error) {
	if this.Id == "" {
		return errors.New("missing deployment id")
	}
	if this.Name == "" {
		return errors.New("missing deployment name")
	}
	if this.XmlRaw == "" {
		return errors.New("missing deployment xml_raw")
	}
	engineAccess, err := xmlContainsEngineAccess(this.XmlRaw)
	if err != nil {
		return err
	}
	if engineAccess {
		return errors.New("process tries to access execution engine")
	}
	if strict && this.Xml == "" {
		return errors.New("missing deployment xml")
	}
	for _, element := range this.Elements {
		err = element.Validate(strict)
		if err != nil {
			return err
		}
	}
	for _, lane := range this.Lanes {
		err = lane.Validate(strict)
		if err != nil {
			return err
		}
	}
	return nil
}

func xmlContainsEngineAccess(xml string) (triesAccess bool, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	doc := etree.NewDocument()
	err = doc.ReadFromString(xml)
	if err != nil {
		return true, err
	}
	scripts := []string{}
	for _, script := range doc.FindElements("//camunda:script") {
		scripts = append(scripts, script.Text())
	}
	for _, script := range doc.FindElements("//bpmn:script") {
		scripts = append(scripts, script.Text())
	}
	for _, script := range scripts {
		if strings.Contains(script, "execution.") {
			return true, nil
		}
	}
	return false, nil
}

func (this Element) Validate(strict bool) error {
	if err := this.Task.Validate(strict); err != nil {
		return err
	}
	if err := this.MultiTask.Validate(strict); err != nil {
		return err
	}
	if err := this.TimeEvent.Validate(strict); err != nil {
		return err
	}
	if err := this.MsgEvent.Validate(strict); err != nil {
		return err
	}
	if err := this.ReceiveTaskEvent.Validate(strict); err != nil {
		return err
	}
	return nil
}

func (this *Task) Validate(strict bool) error {
	if this == nil {
		return nil
	}
	if this.BpmnElementId == "" {
		return errors.New("missing task bpmn id")
	}
	if err := this.Selection.Validate(strict); err != nil {
		return err
	}
	return nil
}

func (this *MultiTask) Validate(strict bool) error {
	if this == nil {
		return nil
	}
	if this.BpmnElementId == "" {
		return errors.New("missing multi task bpmn id")
	}
	if len(this.Selections) == 0 {
		return errors.New("expect at least one device selected for multi task")
	}
	for _, selection := range this.Selections {
		if err := selection.Validate(strict); err != nil {
			return err
		}
	}
	return nil
}

func (this *TimeEvent) Validate(strict bool) error {
	if this == nil {
		return nil
	}
	if this.BpmnElementId == "" {
		return errors.New("missing time event bpmn id")
	}
	if this.Kind == "" {
		return errors.New("missing time event kind")
	}
	if this.Time == "" {
		return errors.New("missing time event time")
	}
	return nil
}

func (this *MsgEvent) Validate(strict bool) error {
	if this == nil {
		return nil
	}
	if this.BpmnElementId == "" {
		return errors.New("missing msg event bpmn id")
	}
	if this.Device == nil {
		return errors.New("missing msg event device selection")
	}
	if this.Service == nil {
		return errors.New("missing msg event service selection")
	}
	if this.Device.Id == "" {
		return errors.New("missing msg event device id in selection")
	}
	if this.Service.Id == "" {
		return errors.New("missing msg event service id in selection")
	}
	if strict && this.EventId == "" {
		return errors.New("missing msg event id")
	}
	if this.Path == "" {
		return errors.New("missing msg event path")
	}
	if this.Operation == "" {
		return errors.New("missing msg event operation")
	}
	if this.TriggerConversion != nil && strict {
		if this.TriggerConversion.From == "" {
			return errors.New("missing msg event cast from")
		}
		if this.TriggerConversion.To == "" {
			return errors.New("missing msg event cast to")
		}
	}
	return nil
}

func (this LaneElement) Validate(strict bool) error {
	if err := this.Lane.Validate(strict); err != nil {
		return err
	}
	if err := this.MultiLane.Validate(strict); err != nil {
		return err
	}
	return nil
}

func (this *Lane) Validate(strict bool) error {
	if this == nil {
		return nil
	}
	if this.BpmnElementId == "" {
		return errors.New("missing lane bpmn id")
	}
	if this.Selection == nil || this.Selection.Id == "" {
		for _, element := range this.Elements {
			if element.LaneTask != nil {
				return errors.New("missing lane device selection")
			}
		}
	}
	for _, element := range this.Elements {
		if err := element.Validate(strict); err != nil {
			return err
		}
	}
	return nil
}

func (this *MultiLane) Validate(strict bool) error {
	if this == nil {
		return nil
	}
	if this.BpmnElementId == "" {
		return errors.New("missing lane bpmn id")
	}
	if len(this.Selections) == 0 {
		return errors.New("expect at least one device selected for multi lane")
	}
	for _, selection := range this.Selections {
		if selection.Id == "" {
			return errors.New("missing multi lane selection device id")
		}
	}
	for _, element := range this.Elements {
		if err := element.Validate(strict); err != nil {
			return err
		}
	}
	return nil
}

func (this LaneSubElement) Validate(strict bool) error {
	if err := this.LaneTask.Validate(strict); err != nil {
		return err
	}
	if err := this.TimeEvent.Validate(strict); err != nil {
		return err
	}
	if err := this.MsgEvent.Validate(strict); err != nil {
		return err
	}
	if err := this.ReceiveTaskEvent.Validate(strict); err != nil {
		return err
	}
	return nil
}

func (this *LaneTask) Validate(strict bool) error {
	if this == nil {
		return nil
	}
	if this.BpmnElementId == "" {
		return errors.New("missing task bpmn id")
	}
	if this.SelectedService == nil || this.SelectedService.Id == "" {
		return errors.New("missing lane task service selection id")
	}
	return nil
}

func (this Selection) Validate(strict bool) error {
	if this.Device == nil {
		return errors.New("missing device selection")
	}
	if this.Service == nil {
		return errors.New("missing service selection")
	}
	if this.Device.Id == "" {
		return errors.New("missing device id in selection")
	}
	if this.Service.Id == "" {
		return errors.New("missing service id in selection")
	}
	return nil
}
