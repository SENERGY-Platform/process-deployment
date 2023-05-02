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

package devices

import (
	"github.com/SENERGY-Platform/permission-search/lib/client"
	"github.com/SENERGY-Platform/permission-search/lib/model"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
)

func (this *Repository) CheckAccess(token auth.Token, kind string, ids []string) (result map[string]bool, err error) {
	if len(ids) == 0 {
		return map[string]bool{}, nil
	}
	result, _, err = client.Query[map[string]bool](this.permissionsearch, token.String(), model.QueryMessage{
		Resource: kind,
		CheckIds: &model.QueryCheckIds{
			Ids:    ids,
			Rights: "x",
		},
	})
	return result, err
}
