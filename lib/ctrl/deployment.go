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

package ctrl

import (
	"context"
	"errors"
	engine "github.com/SENERGY-Platform/camunda-engine-wrapper/lib/client"
	eventdeployment "github.com/SENERGY-Platform/event-deployment/lib/client"
	"github.com/SENERGY-Platform/process-deployment/lib/model/deploymentmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/model/messages"
	"log"
	"runtime/debug"
	"time"
)

func (this *Ctrl) StartSyncLoop(ctx context.Context, interval time.Duration, lockduration time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := this.Sync(lockduration)
				if err != nil {
					log.Printf("ERROR: while db sync run: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (this *Ctrl) Sync(lockduration time.Duration) interface{} {
	return this.db.RetryDeploymentSync(lockduration, this.syncDeploymentDelete, func(command messages.DeploymentCommand) error {
		return this.deleteDeployment(command.Id)
	})
}

func (this *Ctrl) setDeployment(owner string, source string, deployment deploymentmodel.Deployment) error {
	err := this.db.SetDeployment(
		messages.DeploymentCommand{
			Command:    "PUT",
			Id:         deployment.Id,
			Owner:      owner,
			Deployment: &deployment,
			Source:     source,
			Version:    deployment.Version,
		},
		this.syncDeployment,
	)
	if err != nil {
		_ = this.deleteDeployment(deployment.Id)
	}
	return err
}

func (this *Ctrl) syncDeployment(command messages.DeploymentCommand) error {
	if command.Deployment == nil {
		debug.PrintStack()
		return errors.New("deployment cannot be nil")
	}

	err := this.SaveDependencies(command)
	if err != nil {
		return err
	}

	err, _ = this.engine.Deploy(engine.InternalAdminToken, engine.DeploymentMessage{
		Deployment: engine.Deployment{
			Id:               command.Deployment.Id,
			Name:             command.Deployment.Name,
			Diagram:          command.Deployment.Diagram,
			IncidentHandling: command.Deployment.IncidentHandling,
		},
		UserId: command.Owner,
		Source: command.Source,
	})
	if err != nil {
		return err
	}

	err, _ = this.eventdeployment.Deploy(eventdeployment.InternalAdminToken, eventdeployment.Deployment{
		Deployment: *command.Deployment,
		UserId:     command.Owner,
	})
	if err != nil {
		return err
	}
	err = this.publishDeployment(command)
	if err != nil {
		return err
	}
	return nil
}

func (this *Ctrl) publishDeployment(command messages.DeploymentCommand) error {
	command.Command = "PUT"
	return this.deploymentPublisher.Produce(command)
}

func (this *Ctrl) deleteDeployment(id string) error {
	return this.db.DeleteDeployment(id, this.syncDeploymentDelete)
}

func (this *Ctrl) syncDeploymentDelete(command messages.DeploymentCommand) error {
	err := this.DeleteDependencies(command)
	if err != nil {
		return err
	}

	err, _ = this.eventdeployment.DeleteDeployment(eventdeployment.InternalAdminToken, command.Owner, command.Id)
	if err != nil {
		return err
	}

	err, _ = this.engine.DeleteDeployment(engine.InternalAdminToken, command.Owner, command.Id)
	if err != nil {
		return err
	}

	err = this.publishDeploymentDelete(command.Owner, command.Id)
	if err != nil {
		return err
	}
	return nil
}

func (this *Ctrl) publishDeploymentDelete(user string, id string) error {
	return this.deploymentPublisher.Produce(messages.DeploymentCommand{
		Command: "DELETE",
		Id:      id,
		Owner:   user,
		Version: deploymentmodel.CurrentVersion,
	})
}
