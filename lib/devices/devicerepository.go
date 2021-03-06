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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/model/devicemodel"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/coocood/freecache"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)

func (this *Repository) GetDevice(token jwt_http_router.JwtImpersonate, id string) (result devicemodel.Device, err error, code int) {
	err, code = this.get(string(token), "devices", id, &result)
	return
}

func (this *Repository) GetService(token jwt_http_router.JwtImpersonate, id string) (result devicemodel.Service, err error, code int) {
	err, code = this.get(string(token), "services", id, &result)
	return
}

func (this *Repository) GetDeviceGroup(token jwt_http_router.JwtImpersonate, id string) (result devicemodel.DeviceGroup, err error, code int) {
	err, code = this.get(string(token), "device-groups", id, &result)
	return
}

func (this *Repository) get(token string, resource string, id string, result interface{}) (error, int) {
	temp, err := this.l1.Get([]byte(resource + "." + id))
	if err != nil && this.config.Debug {
		if err == freecache.ErrNotFound {
			log.Println("DEBUG: "+resource+" not in cache", id)
		} else {
			log.Println("ERROR: "+resource+" cache retrieval error", id, err)
		}
	}
	if err == nil {
		err = json.Unmarshal(temp, result)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		} else {
			return nil, 200
		}
	} else {
		client := http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := http.NewRequest(
			"GET",
			this.config.DeviceRepoUrl+"/"+resource+"/"+url.PathEscape(id),
			nil,
		)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		req.Header.Set("Authorization", token)

		resp, err := client.Do(req)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			debug.PrintStack()
			return errors.New("unexpected statuscode"), resp.StatusCode
		}

		temp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		err = json.Unmarshal(temp, result)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		err = this.l1.Set([]byte(resource+"."+id), temp, L1Expiration)
		if err != nil {
			log.Println("WARNING: unable to save resource in cache", resource+"."+id, string(temp))
		}
		return nil, 200
	}
}

func (this *Repository) getUncachedList(token string, resource string, result interface{}) (error, int) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(
		"GET",
		this.config.DeviceRepoUrl+"/"+resource,
		nil,
	)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		debug.PrintStack()
		return errors.New("unexpected statuscode"), resp.StatusCode
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	return nil, 200
}
