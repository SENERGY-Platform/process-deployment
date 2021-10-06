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

package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ConfigStruct struct {
	ApiPort                     string `json:"api_port"`
	KafkaUrl                    string `json:"kafka_url"`
	DeviceRepoUrl               string `json:"device_repo_url"`
	ProcessRepoUrl              string `json:"process_repo_url"`
	PermSearchUrl               string `json:"perm_search_url"`
	DeviceSelectionUrl          string `json:"device_selection_url"`
	MongoUrl                    string `json:"mongo_url"`
	MongoTable                  string `json:"mongo_table"`
	MongoDeploymentCollection   string `json:"mongo_deployment_collection"`
	MongoDependenciesCollection string `json:"mongo_dependencies_collection"`
	ConsumerGroup               string `json:"consumer_group"`
	Debug                       bool   `json:"debug"`
	ConnectivityTest            bool   `json:"connectivity_test"`
	DeploymentTopic             string `json:"deployment_topic"`
	UsersTopic                  string `json:"users_topic"`
	NotificationUrl             string `json:"notification_url"`
	ImportDeployUrl             string `json:"import_deploy_url"`

	EnableDeviceGroupsForTasks  bool `json:"enable_device_groups_for_tasks"`
	EnableDeviceGroupsForEvents bool `json:"enable_device_groups_for_events"`
	EnableImportsForEvents      bool `json:"enable_imports_for_events"`

	HttpClientTimeout     string `json:"http_client_timeout"`
	HttpServerTimeout     string `json:"http_server_timeout"`
	HttpServerReadTimeout string `json:"http_server_read_timeout"`
}

type Config = *ConfigStruct

func LoadConfig(location string) (config Config, err error) {
	file, err := os.Open(location)
	if err != nil {
		log.Println("error on config load: ", err)
		return config, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Println("invalid config json: ", err)
		return config, err
	}
	HandleEnvironmentVars(config)
	setDefaultHttpClient(config)
	return config, nil
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func fieldNameToEnvName(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToUpper(strings.Join(a, "_"))
}

// preparations for docker
func HandleEnvironmentVars(config Config) {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	configType := configValue.Type()
	for index := 0; index < configType.NumField(); index++ {
		fieldName := configType.Field(index).Name
		envName := fieldNameToEnvName(fieldName)
		envValue := os.Getenv(envName)
		if envValue != "" {
			fmt.Println("use environment variable: ", envName, " = ", envValue)
			if configValue.FieldByName(fieldName).Kind() == reflect.Int64 {
				i, _ := strconv.ParseInt(envValue, 10, 64)
				configValue.FieldByName(fieldName).SetInt(i)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Float64 {
				f, _ := strconv.ParseFloat(envValue, 64)
				configValue.FieldByName(fieldName).SetFloat(f)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.String {
				configValue.FieldByName(fieldName).SetString(envValue)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Bool {
				b, _ := strconv.ParseBool(envValue)
				configValue.FieldByName(fieldName).SetBool(b)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Slice {
				val := []string{}
				for _, element := range strings.Split(envValue, ",") {
					val = append(val, strings.TrimSpace(element))
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(val))
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Map {
				value := map[string]string{}
				for _, element := range strings.Split(envValue, ",") {
					keyVal := strings.Split(element, ":")
					key := strings.TrimSpace(keyVal[0])
					val := strings.TrimSpace(keyVal[1])
					value[key] = val
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(value))
			}
		}
	}
}

func setDefaultHttpClient(config Config) {
	var err error
	http.DefaultClient.Timeout, err = time.ParseDuration(config.HttpClientTimeout)
	if err != nil {
		log.Println("WARNING: invalid http timeout --> no timeouts\n", err)
	}
}
