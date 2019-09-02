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
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model"
	"github.com/beevik/etree"
)

func Deployment(deployment model.Deployment, selectionAsRef bool, deviceRepo interfaces.DeviceRepository) (xml string, err error) {
	doc := etree.NewDocument()
	err = doc.ReadFromString(deployment.XmlRaw)
	if err != nil {
		return
	}

	for _, element := range deployment.Elements {
		err = Element(doc, element, selectionAsRef, deviceRepo)
		if err != nil {
			return "", err
		}
	}

	for _, lane := range deployment.Lanes {
		err = LaneElement(doc, lane, selectionAsRef, deviceRepo)
		if err != nil {
			return "", err
		}
	}

	xml, err = doc.WriteToString()
	return
}
