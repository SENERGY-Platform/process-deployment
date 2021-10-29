/*
 * Copyright 2020 InfAI (CC SES)
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

package deploymentmodel

import "testing"

func TestElement_Validate(t *testing.T) {
	type fields = TimeEvent
	type args struct {
		kind ValidationKind
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		//P1MT1M
		{name: "PT1M", fields: fields{Type: "timeDuration", Time: "PT1M"}, wantErr: false},
		{name: "P1MT1M", fields: fields{Type: "timeDuration", Time: "P1MT1M"}, wantErr: false},
		{name: "P1MT1M2S", fields: fields{Type: "timeDuration", Time: "P1MT1M2S"}, wantErr: false},
		{name: "PT8S", fields: fields{Type: "timeDuration", Time: "PT8S"}, wantErr: false},
		{name: "PT2S", fields: fields{Type: "timeDuration", Time: "PT2S"}, wantErr: true},
		{name: "PTS", fields: fields{Type: "timeDuration", Time: "PTS"}, wantErr: true},
		{name: "PT", fields: fields{Type: "timeDuration", Time: "PT"}, wantErr: true},
		{name: "P", fields: fields{Type: "timeDuration", Time: "P"}, wantErr: true},
		{name: "", fields: fields{Type: "timeDuration", Time: ""}, wantErr: true},
		{name: "T", fields: fields{Type: "timeDuration", Time: "T"}, wantErr: true},
		{name: "2S", fields: fields{Type: "timeDuration", Time: "2S"}, wantErr: true},
		{name: "2", fields: fields{Type: "timeDuration", Time: "2"}, wantErr: true},
		{name: "8", fields: fields{Type: "timeDuration", Time: "8"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := Element{
				BpmnId:    "foo",
				TimeEvent: &tt.fields,
			}
			if err := this.Validate(ValidateRequest); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
