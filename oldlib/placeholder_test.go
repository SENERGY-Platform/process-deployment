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

package oldlib

import (
	"encoding/json"
	"github.com/SENERGY-Platform/process-deployment/oldlib/etree"
	"github.com/SENERGY-Platform/process-deployment/oldlib/model"
	"github.com/SENERGY-Platform/process-deployment/oldlib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGetTemplateParameterList(t *testing.T) {
	text_1 := `"abc" ${Foo} lasd {{Bar}} asdj ${Batz {asdd} fars {{asgn=42}}`
	expected_1 := []string{"Bar", "asgn=42"}
	result_1, err := GetTemplateParameterList(text_1)
	if err != nil {
		t.Fatal(err)
	}
	test_compareUnorderedSlice(t, expected_1, result_1)
}

func TestParsePlaceholderName(t *testing.T) {
	label, value, err := parsePlaceholderName("asgn=42")
	if err != nil {
		t.Fatal(err)
	}
	if label != "asgn" || value != "42" {
		t.Fatal("unexpected result", label, value)
	}

	label, value, err = parsePlaceholderName("asgn = 42 asd aa")
	if err != nil {
		t.Fatal(err)
	}
	if label != "asgn" || value != "42 asd aa" {
		t.Fatal("unexpected result", label, value)
	}

	label, value, err = parsePlaceholderName("asgn = 42 asd=aa")
	if err != nil {
		t.Fatal(err)
	}
	if label != "asgn" || value != "42 asd=aa" {
		t.Fatal("unexpected result", label, value)
	}

	label, value, err = parsePlaceholderName("asgn")
	if err != nil {
		t.Fatal(err)
	}
	if label != "asgn" || value != "" {
		t.Fatal("unexpected result", label, value)
	}

	label, value, err = parsePlaceholderName("")
	if err == nil || label != "" || value != "" {
		t.Fatal("unexpected result", label, value, err)
	}
}

func TestGetPlaceholder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]interface{}{})
	}))
	defer server.Close()
	util.Config = &util.ConfigStruct{IotRepoUrl: server.URL, DeprecatedTopic: "execute_in_dose"}

	process, err := GetBpmnAbstractPrepare(test_bpmn_example, jwt_http_router.JwtImpersonate(""))

	if err != nil {
		t.Error(err)
		return
	}

	if len(process.PlaceholderTasks) != 1 {
		t.Error("unexpected result: ", process.PlaceholderTasks)
		return
	}

	sort.Slice(process.PlaceholderTasks[0].Parameter, func(i, j int) bool {
		return strings.Compare(process.PlaceholderTasks[0].Parameter[i].Name, process.PlaceholderTasks[0].Parameter[j].Name) >= 0
	})

	if len(process.PlaceholderTasks[0].Parameter) != 2 {
		t.Error("unexpected result: ", process.PlaceholderTasks)
		return
	}

	if process.PlaceholderTasks[0].Parameter[0].Name != "to=foo@bar.com" ||
		process.PlaceholderTasks[0].Parameter[0].Label != "to" ||
		process.PlaceholderTasks[0].Parameter[0].Value != "foo@bar.com" {
		t.Error("unexpected result: ", process.PlaceholderTasks[0].Parameter[0])
		return
	}

	if process.PlaceholderTasks[0].Parameter[1].Name != "subj=Labor Alarm" ||
		process.PlaceholderTasks[0].Parameter[1].Label != "subj" ||
		process.PlaceholderTasks[0].Parameter[1].Value != "Labor Alarm" {
		t.Error("unexpected result: ", process.PlaceholderTasks[0].Parameter[1])
		return
	}
}

func TestRenderPlaceholder(t *testing.T) {
	result, err := renderPlaceholder("abc {{foo=bar@batz.com}} def", []model.Placeholder{{
		Name:  "foo=bar@batz.com",
		Value: "test",
	}})
	if err != nil {
		t.Error(err)
		return
	}
	if result != "abc test def" {
		t.Error(result)
		return
	}
}

func TestInstantiateAbstractProcessWithPlaceholder(t *testing.T) {
	sliceserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]interface{}{})
	}))
	defer sliceserver.Close()
	boolserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(true)
	}))
	defer sliceserver.Close()
	util.Config = &util.ConfigStruct{IotRepoUrl: sliceserver.URL, DeprecatedTopic: "execute_in_dose", PermissionsUrl: boolserver.URL}

	process, err := GetBpmnAbstractPrepare(test_bpmn_example, jwt_http_router.JwtImpersonate(""))

	if err != nil {
		t.Error(err)
		return
	}

	sort.Slice(process.PlaceholderTasks[0].Parameter, func(i, j int) bool {
		return strings.Compare(process.PlaceholderTasks[0].Parameter[i].Name, process.PlaceholderTasks[0].Parameter[j].Name) >= 0
	})

	process.PlaceholderTasks[0].Parameter[0].Value = "test@result.com"

	result, err := InstantiateAbstractProcess(process, jwt_http_router.JwtImpersonate(""), "")
	if err != nil {
		t.Error(err)
		return
	}

	doc := etree.NewDocument()
	err = doc.ReadFromString(result)
	if err != nil {
		t.Error(err)
		return
	}

	id_1 := "Task_163k0g0"
	subDoc_1 := etree.NewDocument()
	subDoc_1.AddChild(doc.FindElement("//bpmn:serviceTask[@id='" + id_1 + "']"))
	subResult_1, err := subDoc_1.WriteToString()
	if err != nil {
		t.Error(err)
		return
	}
	expected_1 := `<bpmn:serviceTask id="Task_163k0g0" name="E-Mail Benachrichtigung an Nutzer">
         <bpmn:extensionElements>
            <camunda:connector>
               <camunda:inputOutput>
                  <camunda:inputParameter name="to">test@result.com</camunda:inputParameter>
                  <camunda:inputParameter name="subject">Labor Alarm</camunda:inputParameter>
                  <camunda:inputParameter name="text">Einbrecher im Labor!</camunda:inputParameter>
               </camunda:inputOutput>
               <camunda:connectorId>mail-send</camunda:connectorId>
            </camunda:connector>
         </bpmn:extensionElements>
         <bpmn:incoming>SequenceFlow_1dwrlmp</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_0d8n4ls</bpmn:outgoing>
      </bpmn:serviceTask>`

	if strings.TrimSpace(subResult_1) != strings.TrimSpace(expected_1) {
		t.Error("unexpected result: \n", subResult_1, "\n\n", expected_1)
		return
	}
}

func test_compareUnorderedSlice(t *testing.T, expected []string, actual []string) {
	t.Helper()
	sort.Strings(expected)
	sort.Strings(actual)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("unexpected result: ", expected, actual)
	}
}

const test_bpmn_example = `<?xml version="1.0" encoding="UTF-8"?>
<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:camunda="http://camunda.org/schema/1.0/bpmn" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn">
   <bpmn:process id="Lab_Alarm_v2" isExecutable="true">
      <bpmn:startEvent id="StartEvent_1">
         <bpmn:outgoing>SequenceFlow_1by6v84</bpmn:outgoing>
      </bpmn:startEvent>
      <bpmn:sequenceFlow id="SequenceFlow_1bzicgx" sourceRef="IntermediateThrowEvent_0y2xm95" targetRef="SubProcess_1ezaxs6" />
      <bpmn:sequenceFlow id="SequenceFlow_1by6v84" sourceRef="StartEvent_1" targetRef="Task_1l82idm" />
      <bpmn:sequenceFlow id="SequenceFlow_0fnoj7y" sourceRef="Task_1l82idm" targetRef="IntermediateThrowEvent_0y2xm95" />
      <bpmn:sequenceFlow id="SequenceFlow_1uopd35" sourceRef="SubProcess_1ezaxs6" targetRef="ExclusiveGateway_1ktiw5p" />
      <bpmn:sequenceFlow id="SequenceFlow_0lfm9ae" sourceRef="ExclusiveGateway_1ktiw5p" targetRef="IntermediateCatchEvent_03l35pd" />
      <bpmn:sequenceFlow id="SequenceFlow_05cqdvm" sourceRef="ExclusiveGateway_1ktiw5p" targetRef="IntermediateCatchEvent_0isg4wc" />
      <bpmn:sequenceFlow id="SequenceFlow_09itq8j" sourceRef="IntermediateCatchEvent_03l35pd" targetRef="ExclusiveGateway_1xmuvj8" />
      <bpmn:sequenceFlow id="SequenceFlow_179ncwq" sourceRef="IntermediateCatchEvent_0isg4wc" targetRef="ExclusiveGateway_1xmuvj8" />
      <bpmn:sequenceFlow id="SequenceFlow_1vpz0nt" sourceRef="ExclusiveGateway_1ktiw5p" targetRef="IntermediateCatchEvent_0h26u1z" />
      <bpmn:sequenceFlow id="SequenceFlow_0qf4iq2" sourceRef="IntermediateCatchEvent_0h26u1z" targetRef="ExclusiveGateway_1xmuvj8" />
      <bpmn:sequenceFlow id="SequenceFlow_0v1xf2m" sourceRef="ExclusiveGateway_1xmuvj8" targetRef="ExclusiveGateway_1e5ticp" />
      <bpmn:sequenceFlow id="SequenceFlow_0edysgu" sourceRef="ExclusiveGateway_1e5ticp" targetRef="Task_062eahn" />
      <bpmn:sequenceFlow id="SequenceFlow_1dwrlmp" sourceRef="ExclusiveGateway_1e5ticp" targetRef="Task_163k0g0" />
      <bpmn:sequenceFlow id="SequenceFlow_0m2pxym" sourceRef="ExclusiveGateway_1e5ticp" targetRef="Task_0rn7nov" />
      <bpmn:sequenceFlow id="SequenceFlow_0lnlwgl" sourceRef="ExclusiveGateway_1e5ticp" targetRef="Task_1s7fede" />
      <bpmn:sequenceFlow id="SequenceFlow_00ld1yr" sourceRef="Task_1s7fede" targetRef="ExclusiveGateway_09ndxg0" />
      <bpmn:sequenceFlow id="SequenceFlow_1hz285l" sourceRef="Task_062eahn" targetRef="ExclusiveGateway_09ndxg0" />
      <bpmn:sequenceFlow id="SequenceFlow_0d8n4ls" sourceRef="Task_163k0g0" targetRef="ExclusiveGateway_09ndxg0" />
      <bpmn:sequenceFlow id="SequenceFlow_1l9t9n8" sourceRef="Task_0rn7nov" targetRef="ExclusiveGateway_09ndxg0" />
      <bpmn:sequenceFlow id="SequenceFlow_0isvsuc" sourceRef="ExclusiveGateway_09ndxg0" targetRef="EndEvent_0jc7ro9" />
      <bpmn:serviceTask id="Task_1l82idm" name="Philips-Hue-Light set_state" camunda:type="external" camunda:topic="execute_in_dose">
         <bpmn:extensionElements>
            <camunda:inputOutput>
               <camunda:inputParameter name="payload">{
    "label": "set_state",
    "device_type": "iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d",
    "service": "iot#fc6c0792-cf86-436b-9f6f-2a4e0db655ed",
    "values": {
        "inputs": {
            "value": {
                "b": 0,
                "bri": 0,
                "g": 0,
                "on": true,
                "r": 0
            }
        },
        "outputs": {
            "status": 0
        }
    }
}</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.b">255</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.bri">255</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.g">255</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.on">false</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.r">255</camunda:inputParameter>
               <camunda:outputParameter name="status">${result.outputs.status}</camunda:outputParameter>
            </camunda:inputOutput>
         </bpmn:extensionElements>
         <bpmn:incoming>SequenceFlow_1by6v84</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_0fnoj7y</bpmn:outgoing>
      </bpmn:serviceTask>
      <bpmn:serviceTask id="Task_062eahn" name="ZWay-SwitchBinary on" camunda:type="external" camunda:topic="execute_in_dose">
         <bpmn:extensionElements>
            <camunda:inputOutput>
               <camunda:inputParameter name="payload">{
    "label": "on",
    "device_type": "iot#7075b34a-23b5-49ba-9723-867613205e0c",
    "service": "iot#deb8c74e-d556-4ab6-b2d8-4f1a74d528a6",
    "values": {}
}</camunda:inputParameter>
            </camunda:inputOutput>
         </bpmn:extensionElements>
         <bpmn:incoming>SequenceFlow_0edysgu</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_1hz285l</bpmn:outgoing>
      </bpmn:serviceTask>
      <bpmn:serviceTask id="Task_163k0g0" name="E-Mail Benachrichtigung an Nutzer">
         <bpmn:extensionElements>
            <camunda:connector>
               <camunda:inputOutput>
                  <camunda:inputParameter name="to">{{to=foo@bar.com}}</camunda:inputParameter>
                  <camunda:inputParameter name="subject">{{subj=Labor Alarm}}</camunda:inputParameter>
                  <camunda:inputParameter name="text">Einbrecher im Labor!</camunda:inputParameter>
               </camunda:inputOutput>
               <camunda:connectorId>mail-send</camunda:connectorId>
            </camunda:connector>
         </bpmn:extensionElements>
         <bpmn:incoming>SequenceFlow_1dwrlmp</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_0d8n4ls</bpmn:outgoing>
      </bpmn:serviceTask>
      <bpmn:serviceTask id="Task_1s7fede" name="Philips-Hue-Light set_state" camunda:type="external" camunda:topic="execute_in_dose">
         <bpmn:extensionElements>
            <camunda:inputOutput>
               <camunda:inputParameter name="payload">{
    "label": "set_state",
    "device_type": "iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d",
    "service": "iot#fc6c0792-cf86-436b-9f6f-2a4e0db655ed",
    "values": {
        "inputs": {
            "value": {
                "b": 0,
                "bri": 0,
                "g": 0,
                "on": true,
                "r": 0
            }
        },
        "outputs": {
            "status": 0
        }
    }
}</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.b">0</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.bri">255</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.g">0</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.on">true</camunda:inputParameter>
               <camunda:inputParameter name="inputs.value.r">255</camunda:inputParameter>
               <camunda:outputParameter name="status">${result.outputs.status}</camunda:outputParameter>
            </camunda:inputOutput>
         </bpmn:extensionElements>
         <bpmn:incoming>SequenceFlow_0lnlwgl</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_00ld1yr</bpmn:outgoing>
      </bpmn:serviceTask>
      <bpmn:serviceTask id="Task_0rn7nov" name="E-Mail Benachrichtigung an Nutzer">
         <bpmn:extensionElements>
            <camunda:connector>
               <camunda:inputOutput>
                  <camunda:inputParameter name="to">batz@blub.org</camunda:inputParameter>
                  <camunda:inputParameter name="subject">Labor Alarm</camunda:inputParameter>
                  <camunda:inputParameter name="text">Einbrecher im Labor!</camunda:inputParameter>
               </camunda:inputOutput>
               <camunda:connectorId>mail-send</camunda:connectorId>
            </camunda:connector>
         </bpmn:extensionElements>
         <bpmn:incoming>SequenceFlow_0m2pxym</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_1l9t9n8</bpmn:outgoing>
      </bpmn:serviceTask>
      <bpmn:intermediateCatchEvent id="IntermediateThrowEvent_0y2xm95" name="10 Sekunden">
         <bpmn:incoming>SequenceFlow_0fnoj7y</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_1bzicgx</bpmn:outgoing>
         <bpmn:timerEventDefinition>
            <bpmn:timeDuration xsi:type="bpmn:tFormalExpression">PT10S</bpmn:timeDuration>
         </bpmn:timerEventDefinition>
      </bpmn:intermediateCatchEvent>
      <bpmn:intermediateCatchEvent id="IntermediateCatchEvent_03l35pd">
         <bpmn:incoming>SequenceFlow_0lfm9ae</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_09itq8j</bpmn:outgoing>
         <bpmn:messageEventDefinition />
      </bpmn:intermediateCatchEvent>
      <bpmn:intermediateCatchEvent id="IntermediateCatchEvent_0isg4wc">
         <bpmn:incoming>SequenceFlow_05cqdvm</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_179ncwq</bpmn:outgoing>
         <bpmn:messageEventDefinition />
      </bpmn:intermediateCatchEvent>
      <bpmn:intermediateCatchEvent id="IntermediateCatchEvent_0h26u1z">
         <bpmn:incoming>SequenceFlow_1vpz0nt</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_0qf4iq2</bpmn:outgoing>
         <bpmn:messageEventDefinition />
      </bpmn:intermediateCatchEvent>
      <bpmn:subProcess id="SubProcess_1ezaxs6" name="Flash Lights">
         <bpmn:incoming>SequenceFlow_1bzicgx</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_1uopd35</bpmn:outgoing>
         <bpmn:multiInstanceLoopCharacteristics isSequential="true">
            <bpmn:loopCardinality xsi:type="bpmn:tFormalExpression">2</bpmn:loopCardinality>
         </bpmn:multiInstanceLoopCharacteristics>
         <bpmn:startEvent id="StartEvent_06keazs">
            <bpmn:outgoing>SequenceFlow_0zwb28u</bpmn:outgoing>
         </bpmn:startEvent>
         <bpmn:sequenceFlow id="SequenceFlow_0zwb28u" sourceRef="StartEvent_06keazs" targetRef="Task_1w1172b" />
         <bpmn:sequenceFlow id="SequenceFlow_0c1q4yx" sourceRef="Task_10avate" targetRef="EndEvent_035t0wg" />
         <bpmn:sequenceFlow id="SequenceFlow_0nrtugv" sourceRef="Task_1w1172b" targetRef="Task_10avate" />
         <bpmn:endEvent id="EndEvent_035t0wg">
            <bpmn:incoming>SequenceFlow_0c1q4yx</bpmn:incoming>
         </bpmn:endEvent>
         <bpmn:serviceTask id="Task_1w1172b" name="Philips-Hue-Light set_state" camunda:type="external" camunda:topic="execute_in_dose">
            <bpmn:extensionElements>
               <camunda:inputOutput>
                  <camunda:inputParameter name="payload">{
    "label": "set_state",
    "device_type": "iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d",
    "service": "iot#fc6c0792-cf86-436b-9f6f-2a4e0db655ed",
    "values": {
        "inputs": {
            "value": {
                "b": 0,
                "bri": 0,
                "g": 0,
                "on": true,
                "r": 0
            }
        },
        "outputs": {
            "status": 0
        }
    }
}</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.b">0</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.bri">255</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.g">255</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.on">true</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.r">0</camunda:inputParameter>
                  <camunda:outputParameter name="status">${result.outputs.status}</camunda:outputParameter>
               </camunda:inputOutput>
            </bpmn:extensionElements>
            <bpmn:incoming>SequenceFlow_0zwb28u</bpmn:incoming>
            <bpmn:outgoing>SequenceFlow_0nrtugv</bpmn:outgoing>
         </bpmn:serviceTask>
         <bpmn:serviceTask id="Task_10avate" name="Philips-Hue-Light set_state" camunda:type="external" camunda:topic="execute_in_dose">
            <bpmn:extensionElements>
               <camunda:inputOutput>
                  <camunda:inputParameter name="payload">{
    "label": "set_state",
    "device_type": "iot#58b47359-a3fe-4e5a-84a2-08bf7b83395d",
    "service": "iot#fc6c0792-cf86-436b-9f6f-2a4e0db655ed",
    "values": {
        "inputs": {
            "value": {
                "b": 0,
                "bri": 0,
                "g": 0,
                "on": true,
                "r": 0
            }
        },
        "outputs": {
            "status": 0
        }
    }
}</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.b">0</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.bri">255</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.g">255</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.on">false</camunda:inputParameter>
                  <camunda:inputParameter name="inputs.value.r">0</camunda:inputParameter>
                  <camunda:outputParameter name="status">${result.outputs.status}</camunda:outputParameter>
               </camunda:inputOutput>
            </bpmn:extensionElements>
            <bpmn:incoming>SequenceFlow_0nrtugv</bpmn:incoming>
            <bpmn:outgoing>SequenceFlow_0c1q4yx</bpmn:outgoing>
         </bpmn:serviceTask>
      </bpmn:subProcess>
      <bpmn:eventBasedGateway id="ExclusiveGateway_1ktiw5p">
         <bpmn:incoming>SequenceFlow_1uopd35</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_0lfm9ae</bpmn:outgoing>
         <bpmn:outgoing>SequenceFlow_05cqdvm</bpmn:outgoing>
         <bpmn:outgoing>SequenceFlow_1vpz0nt</bpmn:outgoing>
      </bpmn:eventBasedGateway>
      <bpmn:exclusiveGateway id="ExclusiveGateway_1xmuvj8">
         <bpmn:incoming>SequenceFlow_09itq8j</bpmn:incoming>
         <bpmn:incoming>SequenceFlow_179ncwq</bpmn:incoming>
         <bpmn:incoming>SequenceFlow_0qf4iq2</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_0v1xf2m</bpmn:outgoing>
      </bpmn:exclusiveGateway>
      <bpmn:parallelGateway id="ExclusiveGateway_1e5ticp">
         <bpmn:incoming>SequenceFlow_0v1xf2m</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_0edysgu</bpmn:outgoing>
         <bpmn:outgoing>SequenceFlow_1dwrlmp</bpmn:outgoing>
         <bpmn:outgoing>SequenceFlow_0m2pxym</bpmn:outgoing>
         <bpmn:outgoing>SequenceFlow_0lnlwgl</bpmn:outgoing>
      </bpmn:parallelGateway>
      <bpmn:parallelGateway id="ExclusiveGateway_09ndxg0">
         <bpmn:incoming>SequenceFlow_00ld1yr</bpmn:incoming>
         <bpmn:incoming>SequenceFlow_1hz285l</bpmn:incoming>
         <bpmn:incoming>SequenceFlow_0d8n4ls</bpmn:incoming>
         <bpmn:incoming>SequenceFlow_1l9t9n8</bpmn:incoming>
         <bpmn:outgoing>SequenceFlow_0isvsuc</bpmn:outgoing>
      </bpmn:parallelGateway>
      <bpmn:endEvent id="EndEvent_0jc7ro9">
         <bpmn:incoming>SequenceFlow_0isvsuc</bpmn:incoming>
      </bpmn:endEvent>
   </bpmn:process>
   <bpmndi:BPMNDiagram id="BPMNDiagram_1">
      <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="Lab_Alarm_v2">
         <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
            <dc:Bounds x="-610" y="252" width="36" height="36" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-592" y="288" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="SubProcess_1ezaxs6_di" bpmnElement="SubProcess_1ezaxs6" isExpanded="true">
            <dc:Bounds x="-261" y="169" width="469" height="202" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="IntermediateCatchEvent_18y6tpe_di" bpmnElement="IntermediateThrowEvent_0y2xm95">
            <dc:Bounds x="-365" y="252" width="36" height="36" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-380" y="288" width="65" height="12" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="StartEvent_06keazs_di" bpmnElement="StartEvent_06keazs">
            <dc:Bounds x="-239" y="248" width="36" height="36" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-220" y="284" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="EndEvent_035t0wg_di" bpmnElement="EndEvent_035t0wg">
            <dc:Bounds x="150" y="248" width="36" height="36" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="169" y="284" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ServiceTask_0848dgc_di" bpmnElement="Task_1w1172b">
            <dc:Bounds x="-149" y="229" width="100" height="80" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ServiceTask_1h01wsj_di" bpmnElement="Task_10avate">
            <dc:Bounds x="4" y="227" width="100" height="80" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ServiceTask_14pg2rf_di" bpmnElement="Task_1l82idm">
            <dc:Bounds x="-523" y="230" width="100" height="80" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ServiceTask_0h54o6v_di" bpmnElement="Task_1s7fede">
            <dc:Bounds x="-264" y="899" width="100" height="80" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="EventBasedGateway_0sakl8k_di" bpmnElement="ExclusiveGateway_1ktiw5p">
            <dc:Bounds x="-51" y="484" width="50" height="50" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-26" y="533.6556064073227" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ExclusiveGateway_1xmuvj8_di" bpmnElement="ExclusiveGateway_1xmuvj8" isMarkerVisible="true">
            <dc:Bounds x="-50.57780320366135" y="639" width="50" height="50" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-25" y="689" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="IntermediateCatchEvent_03l35pd_di" bpmnElement="IntermediateCatchEvent_03l35pd">
            <dc:Bounds x="42.6967963386727" y="559.0995423340961" width="36" height="36" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="61" y="595.0995423340961" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="IntermediateCatchEvent_0isg4wc_di" bpmnElement="IntermediateCatchEvent_0isg4wc">
            <dc:Bounds x="-129.3032036613273" y="559" width="36" height="36" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-111" y="595" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="IntermediateCatchEvent_0h26u1z_di" bpmnElement="IntermediateCatchEvent_0h26u1z">
            <dc:Bounds x="-44" y="559" width="36" height="36" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-26" y="595" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ServiceTask_0owjyph_di" bpmnElement="Task_062eahn">
            <dc:Bounds x="-134" y="899" width="100" height="80" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ServiceTask_06onj0k_di" bpmnElement="Task_163k0g0">
            <dc:Bounds x="-13" y="899" width="100" height="80" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ServiceTask_1jhku9c_di" bpmnElement="Task_0rn7nov">
            <dc:Bounds x="108" y="899" width="100" height="80" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ParallelGateway_0jn95v6_di" bpmnElement="ExclusiveGateway_1e5ticp">
            <dc:Bounds x="-51" y="752" width="50" height="50" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="ParallelGateway_0popgyr_di" bpmnElement="ExclusiveGateway_09ndxg0">
            <dc:Bounds x="-51" y="1082" width="50" height="50" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNShape id="EndEvent_0jc7ro9_di" bpmnElement="EndEvent_0jc7ro9">
            <dc:Bounds x="-44" y="1169" width="36" height="36" />
         </bpmndi:BPMNShape>
         <bpmndi:BPMNEdge id="SequenceFlow_1bzicgx_di" bpmnElement="SequenceFlow_1bzicgx">
            <di:waypoint x="-329" y="270" />
            <di:waypoint x="-261" y="270" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-295" y="255" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0zwb28u_di" bpmnElement="SequenceFlow_0zwb28u">
            <di:waypoint x="-203" y="267" />
            <di:waypoint x="-149" y="267" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-176" y="252" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0c1q4yx_di" bpmnElement="SequenceFlow_0c1q4yx">
            <di:waypoint x="104" y="266" />
            <di:waypoint x="150" y="266" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="127" y="251" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0nrtugv_di" bpmnElement="SequenceFlow_0nrtugv">
            <di:waypoint x="-49" y="267" />
            <di:waypoint x="4" y="267" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-22" y="252" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_1by6v84_di" bpmnElement="SequenceFlow_1by6v84">
            <di:waypoint x="-574" y="270" />
            <di:waypoint x="-523" y="270" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-548" y="255" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0fnoj7y_di" bpmnElement="SequenceFlow_0fnoj7y">
            <di:waypoint x="-423" y="270" />
            <di:waypoint x="-365" y="270" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-394" y="255" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_1uopd35_di" bpmnElement="SequenceFlow_1uopd35">
            <di:waypoint x="-26" y="371" />
            <di:waypoint x="-26" y="484" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-11" y="427.5" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0lfm9ae_di" bpmnElement="SequenceFlow_0lfm9ae">
            <di:waypoint x="-1" y="509" />
            <di:waypoint x="61" y="509" />
            <di:waypoint x="61" y="559" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="30" y="494" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_05cqdvm_di" bpmnElement="SequenceFlow_05cqdvm">
            <di:waypoint x="-51" y="509" />
            <di:waypoint x="-111" y="509" />
            <di:waypoint x="-111" y="559" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-81" y="494" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_09itq8j_di" bpmnElement="SequenceFlow_09itq8j">
            <di:waypoint x="61" y="595" />
            <di:waypoint x="61" y="664" />
            <di:waypoint x="-1" y="664" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="76" y="629.5" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_179ncwq_di" bpmnElement="SequenceFlow_179ncwq">
            <di:waypoint x="-111" y="595" />
            <di:waypoint x="-111" y="664" />
            <di:waypoint x="-51" y="664" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-96" y="629.5" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_1vpz0nt_di" bpmnElement="SequenceFlow_1vpz0nt">
            <di:waypoint x="-26" y="534" />
            <di:waypoint x="-26" y="559" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-11" y="536.5" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0qf4iq2_di" bpmnElement="SequenceFlow_0qf4iq2">
            <di:waypoint x="-26" y="595" />
            <di:waypoint x="-26" y="639" />
            <bpmndi:BPMNLabel>
               <dc:Bounds x="-11" y="607" width="0" height="0" />
            </bpmndi:BPMNLabel>
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0v1xf2m_di" bpmnElement="SequenceFlow_0v1xf2m">
            <di:waypoint x="-26" y="689" />
            <di:waypoint x="-26" y="752" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0edysgu_di" bpmnElement="SequenceFlow_0edysgu">
            <di:waypoint x="-51" y="777" />
            <di:waypoint x="-84" y="777" />
            <di:waypoint x="-84" y="899" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_1dwrlmp_di" bpmnElement="SequenceFlow_1dwrlmp">
            <di:waypoint x="-1" y="777" />
            <di:waypoint x="37" y="777" />
            <di:waypoint x="37" y="899" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0m2pxym_di" bpmnElement="SequenceFlow_0m2pxym">
            <di:waypoint x="-1" y="777" />
            <di:waypoint x="158" y="777" />
            <di:waypoint x="158" y="899" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0lnlwgl_di" bpmnElement="SequenceFlow_0lnlwgl">
            <di:waypoint x="-51" y="777" />
            <di:waypoint x="-214" y="777" />
            <di:waypoint x="-214" y="899" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_00ld1yr_di" bpmnElement="SequenceFlow_00ld1yr">
            <di:waypoint x="-214" y="979" />
            <di:waypoint x="-214" y="1107" />
            <di:waypoint x="-51" y="1107" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_1hz285l_di" bpmnElement="SequenceFlow_1hz285l">
            <di:waypoint x="-84" y="979" />
            <di:waypoint x="-84" y="1107" />
            <di:waypoint x="-51" y="1107" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0d8n4ls_di" bpmnElement="SequenceFlow_0d8n4ls">
            <di:waypoint x="37" y="979" />
            <di:waypoint x="37" y="1107" />
            <di:waypoint x="-1" y="1107" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_1l9t9n8_di" bpmnElement="SequenceFlow_1l9t9n8">
            <di:waypoint x="158" y="979" />
            <di:waypoint x="158" y="1107" />
            <di:waypoint x="-1" y="1107" />
         </bpmndi:BPMNEdge>
         <bpmndi:BPMNEdge id="SequenceFlow_0isvsuc_di" bpmnElement="SequenceFlow_0isvsuc">
            <di:waypoint x="-26" y="1132" />
            <di:waypoint x="-26" y="1169" />
         </bpmndi:BPMNEdge>
      </bpmndi:BPMNPlane>
   </bpmndi:BPMNDiagram>
</bpmn:definitions>`
