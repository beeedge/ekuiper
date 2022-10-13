// Copyright 2021 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processor

import (
	"github.com/lf-edge/ekuiper/pkg/api"
	"reflect"
	"testing"
)

func TestRuleActionParse_Apply(t *testing.T) {
	var tests = []struct {
		ruleStr string
		result  *api.Rule
	}{
		{
			ruleStr: `{
			  "id": "ruleTest",
			  "sql": "SELECT * from demo",
			  "actions": [
				{
				  	"funcName": "RFC_READ_TABLE",
					"ashost":   "192.168.1.100",
					"sysnr":    "02",
					"client":   "900",
					"user":     "SPERF",
					"passwd":   "PASSPASS",
					"params": {
						"QUERY_TABLE": "VBAP",
						"ROWCOUNT":    10,
						"FIELDS": [
							{"FIELDNAME": "MANDT"},
							{"FIELDNAME": "VBELN"},
							{"FIELDNAME": "POSNR"}
						]
					}
				}
			  ]
			}`,
			result: &api.Rule{
				Triggered: false,
				Id:        "ruleTest",
				Sql:       "SELECT * from demo",
				Actions: []map[string]interface{}{
					{
						"funcName": "RFC_READ_TABLE",
						"ashost":   "192.168.1.100",
						"sysnr":    "02",
						"client":   "900",
						"user":     "SPERF",
						"passwd":   "PASSPASS",
						"params": map[string]interface{}{
							"QUERY_TABLE": "VBAP",
							"ROWCOUNT":    float64(10),
							"FIELDS": []interface{}{
								map[string]interface{}{"FIELDNAME": "MANDT"},
								map[string]interface{}{"FIELDNAME": "VBELN"},
								map[string]interface{}{"FIELDNAME": "POSNR"},
							},
						},
					},
				},
				Options: &api.RuleOption{
					IsEventTime:        false,
					LateTol:            1000,
					Concurrency:        1,
					BufferLength:       1024,
					SendMetaToSink:     false,
					Qos:                api.AtMostOnce,
					CheckpointInterval: 300000,
					SendError:          true,
				},
			},
		}, {
			ruleStr: `{
				"id": "ruleTest2",
				"sql": "SELECT * from demo",
				"actions": [
					{
						"log": ""
					},
					{
						"sap": {
							"funcName": "RFC_READ_TABLE",
							"ashost": "192.168.100.10",
							"sysnr": "02",
							"client": "900",
							"user": "uuu",
							"passwd": "ppp."
						}
					}
				],
				"options": {
					"isEventTime": true,
					"lateTolerance": 1000,
					"bufferLength": 10240,
					"qos": 2,
					"checkpointInterval": 60000
				}
			}`,
			result: &api.Rule{
				Triggered: false,
				Id:        "ruleTest2",
				Sql:       "SELECT * from demo",
				Actions: []map[string]interface{}{
					{
						"log": "",
					}, {
						"sap": map[string]interface{}{
							"funcName": "RFC_READ_TABLE",
							"ashost":   "192.168.100.10",
							"sysnr":    "02",
							"client":   "900",
							"user":     "uuu",
							"passwd":   "ppp.",
						},
					},
				},
				Options: &api.RuleOption{
					IsEventTime:        true,
					LateTol:            1000,
					Concurrency:        1,
					BufferLength:       10240,
					SendMetaToSink:     false,
					Qos:                api.ExactlyOnce,
					CheckpointInterval: 60000,
					SendError:          true,
				},
			},
		},
	}

	p := NewRuleProcessor()
	for i, tt := range tests {
		r, err := p.getRuleByJson(tt.result.Id, tt.ruleStr)
		if err != nil {
			t.Errorf("get rule error: %s", err)
		}
		if !reflect.DeepEqual(tt.result, r) {
			t.Errorf("%d \tresult mismatch:\n\nexp=%+v\n\ngot=%+v\n\n", i, tt.result, r)
		}
	}

}
