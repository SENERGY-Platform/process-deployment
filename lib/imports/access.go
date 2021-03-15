/*
 * Copyright 2021 InfAI (CC SES)
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

package imports

import (
	"github.com/SENERGY-Platform/process-deployment/lib/model/importmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/util"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

func (check *Check) CheckAccess(token jwt_http_router.JwtImpersonate, ids []string, alsoCheckTypes bool) (b bool, err error) {
	typeIds := make([]string, len(ids))
	var imports []importmodel.Import
	err = token.GetJSON(check.config.ImportDeployUrl+"/instances", &imports)
	if err != nil {
		return
	}
IDLOOP:
	for i, id := range ids {
		for _, imp := range imports {
			if id == imp.Id {
				typeIds[i] = imp.ImportTypeId
				continue IDLOOP
			}
		}
		return false, nil
	}
	if !alsoCheckTypes {
		return true, nil
	}
	typesAccess, err := util.CheckAccess(check.config.PermSearchUrl, token, "import-types", typeIds)
	if err != nil {
		return
	}
	for _, typeAccess := range typesAccess {
		if !typeAccess {
			return false, nil
		}
	}
	return true, nil
}
