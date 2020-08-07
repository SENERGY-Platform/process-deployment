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
	"bytes"
	"encoding/json"
	"errors"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"runtime/debug"
)

func (this *Repository) CheckAccess(token jwt_http_router.JwtImpersonate, ids []string) (result map[string]bool, err error) {
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(ids)
	if err != nil {
		return result, err
	}
	req, err := http.NewRequest("POST", this.config.PermSearchUrl+"/ids/check/devices/x", buff)
	if err != nil {
		debug.PrintStack()
		return result, err
	}
	req.Header.Set("Authorization", string(token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		debug.PrintStack()
		return result, errors.New(buf.String())
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		debug.PrintStack()
		return result, err
	}
	return result, nil
}
