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

package mocks

import (
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
	"github.com/SENERGY-Platform/process-deployment/lib/model/importmodel"

	"sync"
)

type ImportsMock struct {
	imports       []importmodel.Import
	importTypeIds []string
	mux           sync.Mutex
}

var Imports = &ImportsMock{}

func (this *ImportsMock) New(_ config.Config) (interfaces.Imports, error) {
	return Imports, nil
}

func (this *ImportsMock) CheckAccess(_ string, ids []string, alsoCheckTypes bool) (b bool, err error) {
	typeIds := make([]string, len(ids))
IDLOOP:
	for i, id := range ids {
		for _, imp := range this.imports {
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

TYPEIDLOOP:
	for _, typeId := range typeIds {
		for _, imp := range this.imports {
			if typeId == imp.ImportTypeId {
				continue TYPEIDLOOP
			}
		}
		return false, nil
	}

	return true, nil
}

func (this *ImportsMock) SetImports(imports []importmodel.Import) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.imports = imports
}

func (this *ImportsMock) SetImportTypeIds(importTypeIds []string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.importTypeIds = importTypeIds
}
