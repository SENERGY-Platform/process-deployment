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
	"errors"
	"log"
	"sync"

	"github.com/SmartEnergyPlatform/process-deployment/lib/model"
	"github.com/SmartEnergyPlatform/process-deployment/lib/util"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var instance *mgo.Session
var once sync.Once

func getDb() *mgo.Session {
	once.Do(func() {
		session, err := mgo.Dial(util.Config.MongoUrl)
		if err != nil {
			log.Fatal("error on connection to mongodb: ", err)
		}
		session.SetMode(mgo.Monotonic, true)
		instance = session
	})
	return instance.Copy()
}

func getMetadataCollection() (session *mgo.Session, collection *mgo.Collection) {
	session = getDb()
	collection = session.DB(util.Config.MongoTable).C(util.Config.ProcessMetadataCollection)
	err := collection.EnsureIndexKey("device")
	if err != nil {
		log.Fatal("error on db device index: ", err)
	}
	return
}

type Metadata struct {
	Process  string                `json:"process" bson:"process"`
	Abstract model.AbstractProcess `json:"abstract" bson:"abstract"`
	Online   bool                  `json:"online" bson:"-"`
	Owner    string                `json:"owner" bson:"owner"`
}

func GetMetadata(id string, owner string) (result Metadata, err error) {
	session, collection := getMetadataCollection()
	defer session.Close()
	err = collection.Find(bson.M{"process": id, "owner": owner}).One(&result)
	return
}

func RemoveMetadata(id string) (err error) {
	session, collection := getMetadataCollection()
	defer session.Close()
	_, err = collection.RemoveAll(bson.M{"process": id})
	return
}

func GetAllMetadata(owner string) (result []Metadata, err error) {
	session, collection := getMetadataCollection()
	defer session.Close()
	err = collection.Find(bson.M{"owner": owner}).All(&result)
	return
}

func SetMetadata(id string, deployment model.DeploymentRequest, owner string) (err error) {
	data := sanitizeDeployment(deployment)
	metadata := Metadata{Process: id, Abstract: data, Owner: owner}
	session, collection := getMetadataCollection()
	defer session.Close()
	_, err = collection.Upsert(bson.M{"process": id, "owner": owner}, metadata)
	if err != nil {
		b, _ := json.Marshal(metadata)
		log.Println("ERROR: metadata ", string(b))
	}
	return
}

func sanitizeDeployment(deployment model.DeploymentRequest) (result model.AbstractProcess) {
	result.Name = deployment.Process.Name
	result.Xml = deployment.Process.Xml
	result.MsgEvents = deployment.Process.MsgEvents
	result.ReceiveTasks = deployment.Process.ReceiveTasks
	result.AbstractTasks = sanitizeDeploymentParameter(deployment.Process.AbstractTasks)
	return
}

func sanitizeDeploymentParameter(parameters []model.AbstractTask) (result []model.AbstractTask) {
	for _, param := range parameters {
		result = append(result, model.AbstractTask{Selected: param.Selected, Tasks: param.Tasks, DeviceTypeId: param.DeviceTypeId})
	}
	return
}

func CheckAccess(id string, owner string) (err error) {
	session, collection := getMetadataCollection()
	defer session.Close()
	metadata := []Metadata{}
	err = collection.Find(bson.M{"process": id}).All(&metadata)
	if err != nil {
		return err
	}
	if len(metadata) == 0 {
		return nil	//allow deletion of inconsistent data
	}
	if metadata[0].Owner != owner {
		return errors.New("access denied")
	}
	return
}


func MetadataExists(id string) (exists bool, err error) {
	session, collection := getMetadataCollection()
	defer session.Close()
	metadata := []Metadata{}
	err = collection.Find(bson.M{"process": id}).All(&metadata)
	return len(metadata) > 0, err
}
