/*
 * Copyright 2018 InfAI (CC SES)
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

package lib

import (
	"errors"
	"github.com/SmartEnergyPlatform/process-deployment/lib/etree"
	"github.com/SmartEnergyPlatform/process-deployment/lib/model"
	"github.com/cbroglie/mustache"
	"log"
	"strings"
)

func getPlaceholder(element *etree.Element) (result []model.Placeholder, err error) {
	placeholderNames, err := GetTemplateParameterList(element.Text())
	if err != nil {
		return result, err
	}
	for _, name := range placeholderNames {
		placeholder := model.Placeholder{Name: name}
		placeholder.Label, placeholder.Value, err = parsePlaceholderName(name)
		if err != nil {
			log.Println("WARNING: error in getPlaceholder()", err)
			//continue work but ignore placeholder
			err = nil
		} else {
			result = append(result, placeholder)
		}
	}
	return
}

func parsePlaceholderName(name string) (label string, value string, err error) {
	if name == "" {
		err = errors.New("missing placeholder name to parse")
		return
	}
	parts := strings.Split(name, "=")
	switch len(parts) {
	case 0:
		err = errors.New("wtfh")
		return
	case 1:
		label = strings.TrimSpace(name)
		return
	default:
		label = strings.TrimSpace(parts[0])
		value = strings.TrimSpace(strings.Join(parts[1:], "="))
		return
	}
}

func GetTemplateParameterList(text string) (result []string, err error) {
	templ, err := mustache.ParseString(text)
	if err != nil {
		return result, err
	}
	tags := templ.Tags()
	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		tagType := tag.Type()
		tagName := tag.Name()
		if tagType != mustache.Variable {
			subTags := tag.Tags()
			if len(subTags) > 0 {
				tags = append(tags, subTags...)
			}
		}
		if tagType == mustache.Variable {
			result = append(result, tagName)
		}
	}
	return result, err
}

func renderPlaceholder(text string, placeholders []model.Placeholder) (result string, err error) {
	context := map[string]interface{}{}
	for _, placeholder := range placeholders {
		contextName, contextValue, err := getPlaceholderRenderingContext(strings.Split(placeholder.Name, "."), placeholder.Value)
		if err != nil {
			return result, err
		}
		context[contextName] = contextValue
	}
	return mustache.Render(text, context)
}

func getPlaceholderRenderingContext(path []string, value string) (name string, resultvalue interface{}, err error) {
	if len(path) == 0 {
		err = errors.New("expect path with min length of 1")
		return
	}
	if len(path) == 1 {
		return path[0], value, nil
	}
	name = path[0]

	subName, subValue, err := getPlaceholderRenderingContext(path[1:], value)
	if err != nil {
		return name, resultvalue, err
	}
	resultvalue = map[string]interface{}{subName: subValue}
	return
}
