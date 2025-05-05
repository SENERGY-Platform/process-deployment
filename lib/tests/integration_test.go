/*
 * Copyright 2025 InfAI (CC SES)
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

package tests

import "testing"

func TestIntegration(t *testing.T) {
	t.Skip("TODO")
	//TODO:
	//	- create full process stack (kafka, events, incidents, engine, engine-wrapper, task-worker?)
	//		- those services should communicate via http
	//  - check event deployment:
	//  	- deploy process with event handling
	//  	- start process
	//  	- trigger event
	//  	- check process instance finished
	//  - check incident create/delete:
	//  	- deploy and start process
	//		- trigger incident (via task-worker?)
	//		- check incident existence
	//		- stop-process and/or delete deployment
	//		- check incident is deleted
	//	- check incident handling:
	//		- deploy and start process with on-incident handler
	//		- trigger incident (via task-worker?)
	//		- check handler has worked
}
