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
	permv2 "github.com/SENERGY-Platform/permissions-v2/pkg/client"
	"github.com/SENERGY-Platform/process-deployment/lib/auth"
)

func (this *Repository) CheckAccess(token auth.Token, kind string, ids []string) (result map[string]bool, err error) {
	if len(ids) == 0 {
		return map[string]bool{}, nil
	}
	result, err, _ = this.permv2.CheckMultiplePermissions(token.Jwt(), kind, ids, permv2.Execute)
	return result, err
}
