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
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"log"
)

func (this *Ctrl) HandleUsersCommand(userMsg messages.UserCommandMsg) error {
	if userMsg.Command != "DELETE" {
		return nil
	}
	deployments, err := this.db.GetDeploymentIds(userMsg.Id)
	if err != nil {
		log.Println("ERROR: unable to handle users command (this.db.GetDeploymentIds)", err)
		return err
	}
	for _, id := range deployments {
		err = this.deleteDeployment(id)
		if err != nil {
			log.Println("ERROR: unable to handle users command (this.deleteDeployment)", err)
			return err
		}
	}
	return nil
}
