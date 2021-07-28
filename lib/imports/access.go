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
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
	"github.com/SENERGY-Platform/process-deployment/lib/model/importmodel"
	"github.com/SENERGY-Platform/process-deployment/lib/util"
	"net/http"
)

func (check *Check) CheckAccess(token auth.Token, ids []string, alsoCheckTypes bool) (b bool, err error) {
	typeIds := make([]string, len(ids))
	var imports []importmodel.Import
	req, err := http.NewRequest("GET", check.config.ImportDeployUrl+"/instances", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return false, errors.New(buf.String())
	}
	err = json.NewDecoder(resp.Body).Decode(&imports)
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
