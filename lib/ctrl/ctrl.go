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

package ctrl

import (
	"context"
	"github.com/SENERGY-Platform/process-deployment/lib/config"
	"github.com/SENERGY-Platform/process-deployment/lib/interfaces"
)

type Ctrl struct {
	config        config.Config
	db            interfaces.Database
	connectionLog interfaces.Connectionlog
	semanticRepo  interfaces.SemanticRepository
}

func New(ctx context.Context, config config.Config, sourcing interfaces.SourcingFactory, db interfaces.Database, connlog interfaces.Connectionlog, repo interfaces.SemanticRepository) (result *Ctrl, err error) {
	result = &Ctrl{
		config:        config,
		db:            db,
		connectionLog: connlog,
		semanticRepo:  repo,
	}
	return result, nil
}
