/*
 * Copyright 2025 InfAI (CC SES)
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
	eventdeployment "github.com/SENERGY-Platform/event-deployment/lib/client"
	"log"
)

func (this *Ctrl) HandleDeviceGroupCommand(dgId string) error {
	err, _ := this.eventdeployment.UpdateDeploymentsOfDeviceGroup(eventdeployment.InternalAdminToken, dgId)
	if err != nil {
		log.Printf("ERROR: unable to handle device-group update err=%v, dg=%#v", err, dgId)
		return err
	}
	return nil
}
