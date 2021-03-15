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

package interfaces

import (
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
)

type ImportsFactory interface {
	New(config config.Config) (Imports, error)
}

type Imports interface {
	CheckAccess(token jwt_http_router.JwtImpersonate, ids []string, alsoCheckTypes bool) (bool, error)
}
