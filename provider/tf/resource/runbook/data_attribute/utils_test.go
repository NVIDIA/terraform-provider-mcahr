// SPDX-FileCopyrightText: Copyright (c) 2025 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"context"
	"reflect"
	"testing"

	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsDataJSONEmpty(t *testing.T) {
	tests := []struct {
		name     string
		data     *types.String
		expected bool
	}{
		{
			name:     "Nil data",
			data:     nil,
			expected: true,
		},
		{
			name: "Null data",
			data: func() *types.String {
				s := types.StringNull()
				return &s
			}(),
			expected: true,
		},
		{
			name: "Unknown data",
			data: func() *types.String {
				s := types.StringUnknown()
				return &s
			}(),
			expected: true,
		},
		{
			name: "Empty string data",
			data: func() *types.String {
				s := types.StringValue("")
				return &s
			}(),
			expected: true,
		},
		{
			name: "Empty JSON object",
			data: func() *types.String {
				s := types.StringValue("{}")
				return &s
			}(),
			expected: true,
		},
		{
			name: "Non-empty JSON",
			data: func() *types.String {
				s := types.StringValue(`{"key": "value"}`)
				return &s
			}(),
			expected: false,
		},
		{
			name: "JSON with whitespace",
			data: func() *types.String {
				s := types.StringValue(`{ "key": "value" }`)
				return &s
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := IsDataJSONEmpty(tt.data)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsFieldInDataJSON(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		dataMap   map[string]interface{}
		expected  bool
	}{
		{
			name:      "Field exists in snake_case",
			fieldName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeout_ms": 5000,
			},
			expected: true,
		},
		{
			name:      "Field exists in camelCase",
			fieldName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeoutMs": 5000,
			},
			expected: true,
		},
		{
			name:      "Field does not exist",
			fieldName: "missing_field",
			dataMap: map[string]interface{}{
				"other_field": "value",
			},
			expected: false,
		},
		{
			name:      "Field with null value exists",
			fieldName: "null_field",
			dataMap: map[string]interface{}{
				"null_field": nil,
			},
			expected: false, // nil values are not considered as existing
		},
		{
			name:      "Empty map",
			fieldName: "any_field",
			dataMap:   map[string]interface{}{},
			expected:  false,
		},
		{
			name:      "Nil map",
			fieldName: "any_field",
			dataMap:   nil,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := IsFieldInDataJSON(tt.fieldName, tt.dataMap)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		dataMap   map[string]interface{}
		expected  interface{}
	}{
		{
			name:      "Get value in snake_case",
			fieldName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeout_ms": 5000,
			},
			expected: 5000,
		},
		{
			name:      "Get value in camelCase",
			fieldName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeoutMs": 3000,
			},
			expected: 3000,
		},
		{
			name:      "Prefer exact match over case conversion",
			fieldName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeout_ms": 1000,
				"timeoutMs":  2000,
			},
			expected: 1000, // Should return exact match
		},
		{
			name:      "Field not found",
			fieldName: "missing",
			dataMap: map[string]interface{}{
				"other": "value",
			},
			expected: nil,
		},
		{
			name:      "Complex value",
			fieldName: "params",
			dataMap: map[string]interface{}{
				"params": []interface{}{"p1", "p2"},
			},
			expected: []interface{}{"p1", "p2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := GetFieldValue(tt.fieldName, tt.dataMap)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsJSONField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		expected  bool
	}{
		{
			name:      "params is JSON field",
			fieldName: "params",
			expected:  true,
		},
		{
			name:      "cells is JSON field",
			fieldName: "cells",
			expected:  true,
		},
		{
			name:      "external_params is JSON field",
			fieldName: "external_params",
			expected:  true,
		},
		{
			name:      "Regular field",
			fieldName: "name",
			expected:  false,
		},
		{
			name:      "Another regular field",
			fieldName: "description",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := IsJSONField(tt.fieldName)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsJSONSkipField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		expected  bool
	}{
		{
			name:      "cells_full should be skipped",
			fieldName: "cells_full",
			expected:  true,
		},
		{
			name:      "params_full should be skipped",
			fieldName: "params_full",
			expected:  true,
		},
		{
			name:      "external_params_full should be skipped",
			fieldName: "external_params_full",
			expected:  true,
		},
		{
			name:      "data field should be skipped",
			fieldName: "data",
			expected:  true,
		},
		{
			name:      "Regular field not skipped",
			fieldName: "cells",
			expected:  false,
		},
		{
			name:      "Another regular field",
			fieldName: "name",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := IsJSONSkipField(tt.fieldName)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOnEachStructField(t *testing.T) {
	tests := []struct {
		name          string
		tfModel       *model.RunbookTFModel
		expectError   bool
		expectedCalls []string
	}{
		{
			name: "Iterate over all exportable fields",
			tfModel: &model.RunbookTFModel{
				Name:        types.StringValue("test"),
				Enabled:     types.BoolValue(true),
				Description: types.StringValue("desc"),
			},
			expectError: false,
			expectedCalls: []string{
				"name", "enabled", "description", "timeout_ms",
				"allowed_resources_query", "communication_workspace", "category",
				"communication_channel", "is_run_output_persisted",
				"filter_resource_to_action", "communication_cud_notifications",
				"communication_approval_notifications", "communication_execution_notifications",
				"allowed_entities", "approvers", "labels", "editors", "secret_names",
				"cells", "params", "external_params", "cells_full", "params_full", "params_groups",
				"external_params_full", "data",
			},
		},
		{
			name:        "Handle error from field function",
			tfModel:     &model.RunbookTFModel{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()
			calledFields := []string{}

			var fieldFunc func(string, *reflect.Value) error
			if tt.expectError {
				fieldFunc = func(fieldName string, fieldValue *reflect.Value) error {
					if fieldName == "name" {
						return assert.AnError
					}
					return nil
				}
			} else {
				fieldFunc = func(fieldName string, fieldValue *reflect.Value) error {
					calledFields = append(calledFields, fieldName)
					return nil
				}
			}

			// when
			err := OnEachStructField(ctx, tt.tfModel, fieldFunc)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				// Check that we iterated over expected fields
				for _, expectedField := range tt.expectedCalls {
					assert.Contains(t, calledFields, expectedField)
				}
			}
		})
	}
}

func TestFindValueInMap(t *testing.T) {
	tests := []struct {
		name          string
		snakeCaseName string
		dataMap       map[string]interface{}
		expected      interface{}
	}{
		{
			name:          "Find exact snake_case match",
			snakeCaseName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeout_ms": 5000,
			},
			expected: 5000,
		},
		{
			name:          "Find camelCase equivalent",
			snakeCaseName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeoutMs": 3000,
			},
			expected: 3000,
		},
		{
			name:          "Find with complex field name",
			snakeCaseName: "communication_cud_notifications",
			dataMap: map[string]interface{}{
				"communicationCudNotifications": true,
			},
			expected: true,
		},
		{
			name:          "Not found returns nil",
			snakeCaseName: "missing_field",
			dataMap: map[string]interface{}{
				"other_field": "value",
			},
			expected: nil,
		},
		{
			name:          "Nil value returns nil",
			snakeCaseName: "null_field",
			dataMap: map[string]interface{}{
				"null_field": nil,
			},
			expected: nil,
		},
		{
			name:          "Single word field",
			snakeCaseName: "name",
			dataMap: map[string]interface{}{
				"name": "test_name",
			},
			expected: "test_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := findValueInMap(tt.snakeCaseName, tt.dataMap)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDataJSONToMap(t *testing.T) {
	tests := []struct {
		name        string
		dataJSON    types.String
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "Valid JSON",
			dataJSON: types.StringValue(`{
				"name": "test",
				"enabled": true,
				"timeout_ms": 5000
			}`),
			expected: map[string]interface{}{
				"name":       "test",
				"enabled":    true,
				"timeout_ms": float64(5000),
			},
			expectError: false,
		},
		{
			name:        "Empty JSON",
			dataJSON:    types.StringValue("{}"),
			expected:    nil,
			expectError: false,
		},
		{
			name:        "Null JSON",
			dataJSON:    types.StringNull(),
			expected:    nil,
			expectError: false,
		},
		{
			name:        "Unknown JSON",
			dataJSON:    types.StringUnknown(),
			expected:    nil,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			dataJSON:    types.StringValue("{invalid json}"),
			expectError: true,
		},
		{
			name: "Nested JSON",
			dataJSON: types.StringValue(`{
				"name": "test",
				"config": {
					"key": "value",
					"nested": {
						"deep": true
					}
				}
			}`),
			expected: map[string]interface{}{
				"name": "test",
				"config": map[string]interface{}{
					"key": "value",
					"nested": map[string]interface{}{
						"deep": true,
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := ParseDataJSONToMap(tt.dataJSON)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
