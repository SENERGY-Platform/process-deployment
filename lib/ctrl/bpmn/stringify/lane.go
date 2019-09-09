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

package stringify

import (
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
)

func LaneElement(doc *etree.Document, lane model.LaneElement, selectionAsRef bool, deviceRepo interfaces.Devices) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if err := Lane(doc, lane.Lane, selectionAsRef, deviceRepo); err != nil {
		return err
	}
	if err := MultiLane(doc, lane.MultiLane, selectionAsRef, deviceRepo); err != nil {
		return err
	}
	return nil
}

func Lane(doc *etree.Document, lane *model.Lane, selectionAsRef bool, deviceRepo interfaces.Devices) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if lane == nil {
		return nil
	}

	for _, element := range lane.Elements {
		if err := LaneTask(doc, element.LaneTask, lane.Selection, selectionAsRef, deviceRepo); err != nil {
			return err
		}
		if err := MsgEvent(doc, element.MsgEvent); err != nil {
			return err
		}

		if err := ReceiverTask(doc, element.ReceiveTaskEvent); err != nil {
			return err
		}

		if err := TimeEvent(doc, element.TimeEvent); err != nil {
			return err
		}
	}
	return nil
}

func MultiLane(doc *etree.Document, lane *model.MultiLane, selectionAsRef bool, deviceRepo interfaces.Devices) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if lane == nil {
		return nil
	}

	for _, element := range lane.Elements {
		if err := LaneMultiTask(doc, element.LaneTask, lane.Selections, selectionAsRef, deviceRepo); err != nil {
			return err
		}
		if err := MsgEvent(doc, element.MsgEvent); err != nil {
			return err
		}

		if err := ReceiverTask(doc, element.ReceiveTaskEvent); err != nil {
			return err
		}

		if err := TimeEvent(doc, element.TimeEvent); err != nil {
			return err
		}
	}
	return nil
}
