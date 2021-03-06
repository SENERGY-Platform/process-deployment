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
	"github.com/SENERGY-Platform/process-deployment/lib/util"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

func (this *Repository) CheckAccess(token jwt_http_router.JwtImpersonate, kind string, ids []string) (result map[string]bool, err error) {
	if len(ids) == 0 {
		return map[string]bool{}, nil
	}
	return util.CheckAccess(this.config.PermSearchUrl, token, kind, ids)
}
