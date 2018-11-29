/*
 * Copyright 2018 InfAI (CC SES)
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


package lib

import (
	"encoding/json"
	"github.com/SmartEnergyPlatform/amqp-wrapper-lib"
	"github.com/SmartEnergyPlatform/process-deployment/lib/model"
	"github.com/SmartEnergyPlatform/process-deployment/lib/util"
	"log"
)

var amqp *amqp_wrapper_lib.Connection

type DeploymentCommand struct {
	Command    string          		 	`json:"command"`
	Id         string           		`json:"id"`
	Owner      string           		`json:"owner"`
	Deployment model.DeploymentRequest	`json:"deployment"`
	DeploymentXml string				`json:"deployment_xml"`
}

func InitEventSourcing()(err error){
	amqp, err = amqp_wrapper_lib.Init(util.Config.AmqpUrl, []string{util.Config.AmqpDeploymentTopic}, util.Config.AmqpReconnectTimeout)
	if err != nil {
		return err
	}

	//metadata delete
	err = amqp.Consume(util.Config.AmqpConsumerName + "_" +util.Config.AmqpDeploymentTopic, util.Config.AmqpDeploymentTopic, func(delivery []byte) error {
		command := DeploymentCommand{}
		err = json.Unmarshal(delivery, &command)
		if err != nil {
			log.Println("ERROR: unable to parse amqp event as json \n", err, "\n --> ignore event \n", string(delivery))
			return nil
		}
		log.Println("amqp receive ", string(delivery))
		switch command.Command {
		case "POST":
			return nil
		case "PUT":
			return handleDeploymentMetadataUpdate(command)
		case "DELETE":
			return handleDeploymentMetadataDelete(command)
		default:
			log.Println("WARNING: unknown event type", string(delivery))
			return nil
		}
	})
	if err != nil {
		return err
	}

	return err
}


func handleDeploymentMetadataDelete(command DeploymentCommand) error {
	return RemoveMetadata(command.Id)
}


func handleDeploymentMetadataUpdate(command DeploymentCommand)error{
	err := SetMetadata(command.Id, command.Deployment, command.Owner)
	if err != nil {
		log.Println("WARNING: unable to update process metadata", err, command)
	}
	return err
}

func CloseEventSourcing(){
	amqp.Close()
}

func PublishDeployment(userId string, deployment model.DeploymentRequest, xml string)error{
	command := DeploymentCommand{Owner:userId, Deployment:deployment, DeploymentXml:xml,Command:"POST"}
	payload, err := json.Marshal(command)
	if err != nil {
		return err
	}
	return amqp.Publish(util.Config.AmqpDeploymentTopic, payload)
}

func PublishDeploymentDelete(id string)error{
	command := DeploymentCommand{Id:id,Command:"DELETE"}
	payload, err := json.Marshal(command)
	if err != nil {
		return err
	}
	return amqp.Publish(util.Config.AmqpDeploymentTopic, payload)
}