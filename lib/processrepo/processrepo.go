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

package processrepo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/processmodel"

	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)

type RepositoryFactory struct{}

type Repository struct {
	config config.Config
}

func (this *RepositoryFactory) New(ctx context.Context, config config.Config) (interfaces.ProcessRepo, error) {
	return &Repository{
		config: config,
	}, nil
}

var Factory = &RepositoryFactory{}

func (this *Repository) GetProcessModel(token string, id string) (result processmodel.ProcessModel, err error, errCode int) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(
		"GET",
		this.config.ProcessRepoUrl+"/processes/"+url.PathEscape(id),
		nil,
	)
	if err != nil {
		debug.PrintStack()
		return result, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		return result, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		debug.PrintStack()
		return result, errors.New("unexpected statuscode"), resp.StatusCode
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		debug.PrintStack()
		return result, err, http.StatusInternalServerError
	}
	return result, nil, 200
}
