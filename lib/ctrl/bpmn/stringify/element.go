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
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/beevik/etree"
	"log"
	"runtime/debug"
)

func Element(doc *etree.Document, element deploymentmodel.Element, deviceRepo interfaces.Devices) (err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			log.Printf("%s: %s", r, debug.Stack())
			err = errors.New(fmt.Sprint("Recovered Error: ", r))
		}
	}()
	if err := Task(doc, element.Task, deviceRepo); err != nil {
		return err
	}

	if err := MultiTask(doc, element.MultiTask, deviceRepo); err != nil {
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
	return nil
}
