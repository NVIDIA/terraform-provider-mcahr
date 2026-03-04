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
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertDataValueToTerraformValue(t *testing.T) {
	tests := []struct {
		name        string
		fieldType   reflect.Type
		dataValue   interface{}
		fieldName   string
		expectError bool
		validate    func(t *testing.T, result interface{})
	}{
		{
			name:        "Convert string value",
			fieldType:   reflect.TypeOf(types.StringNull()),
			dataValue:   "test_string",
			fieldName:   "name",
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				strVal, ok := result.(types.String)
				require.True(t, ok)
				assert.Equal(t, "test_string", strVal.ValueString())
			},
		},
		{
			name:        "Convert bool value",
			fieldType:   reflect.TypeOf(types.BoolNull()),
			dataValue:   true,
			fieldName:   "enabled",
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				boolVal, ok := result.(types.Bool)
				require.True(t, ok)
				assert.Equal(t, true, boolVal.ValueBool())
			},
		},
		{
			name:        "Convert int64 from float64",
			fieldType:   reflect.TypeOf(types.Int64Null()),
			dataValue:   float64(12345),
			fieldName:   "timeout_ms",
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				intVal, ok := result.(types.Int64)
				require.True(t, ok)
				assert.Equal(t, int64(12345), intVal.ValueInt64())
			},
		},
		{
			name:        "Convert int64 from int",
			fieldType:   reflect.TypeOf(types.Int64Null()),
			dataValue:   int(100),
			fieldName:   "count",
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				intVal, ok := result.(types.Int64)
				require.True(t, ok)
				assert.Equal(t, int64(100), intVal.ValueInt64())
			},
		},
		{
			name:        "Convert set from array",
			fieldType:   reflect.TypeOf(types.ListNull(types.StringType)),
			dataValue:   []interface{}{"item1", "item2", "item3"},
			fieldName:   "labels",
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				setVal, ok := result.(types.List)
				require.True(t, ok)
				assert.Equal(t, 3, len(setVal.Elements()))
			},
		},
		{
			name:        "Convert JSON field from map",
			fieldType:   reflect.TypeOf(types.StringNull()),
			dataValue:   map[string]interface{}{"key": "value"},
			fieldName:   "params",
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				strVal, ok := result.(types.String)
				require.True(t, ok)
				assert.Contains(t, strVal.ValueString(), `"key":"value"`)
			},
		},
		{
			name:        "Convert cells field",
			fieldType:   reflect.TypeOf(types.StringNull()),
			dataValue:   []interface{}{map[string]interface{}{"name": "cell1", "type": "OP_LANG", "content": "code"}},
			fieldName:   "cells",
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				strVal, ok := result.(types.String)
				require.True(t, ok)
				assert.Contains(t, strVal.ValueString(), `"op":"code"`)
			},
		},
		{
			name: "Convert params_groups object",
			fieldType: reflect.TypeOf(types.ObjectNull(map[string]attr.Type{
				"required": types.ListType{ElemType: types.StringType},
				"optional": types.ListType{ElemType: types.StringType},
				"exported": types.ListType{ElemType: types.StringType},
				"external": types.ListType{ElemType: types.StringType},
			})),
			dataValue: map[string]interface{}{
				"required": []interface{}{"p1"},
				"optional": []interface{}{},
				"exported": []interface{}{"p2", "p3"},
				"external": []interface{}{"p4"},
			},
			fieldName:   "params_groups",
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				objVal, ok := result.(types.Object)
				require.True(t, ok)
				require.False(t, objVal.IsNull())
				attrs := objVal.Attributes()
				// required
				reqList, ok := attrs["required"].(types.List)
				require.True(t, ok)
				assert.Equal(t, 1, len(reqList.Elements()))
				// optional should be present (empty list, not null)
				optList, ok := attrs["optional"].(types.List)
				require.True(t, ok)
				assert.Equal(t, 0, len(optList.Elements()))
			},
		},
		{
			name:        "Invalid bool conversion",
			fieldType:   reflect.TypeOf(types.BoolNull()),
			dataValue:   "not_a_bool",
			fieldName:   "enabled",
			expectError: true,
		},
		{
			name:        "Invalid int64 conversion",
			fieldType:   reflect.TypeOf(types.Int64Null()),
			dataValue:   "not_a_number",
			fieldName:   "timeout_ms",
			expectError: true,
		},
		{
			name:        "Unsupported field type",
			fieldType:   reflect.TypeOf(struct{}{}),
			dataValue:   "any_value",
			fieldName:   "unsupported",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := convertDataValueToTerraformValue(tt.fieldType, tt.dataValue, tt.fieldName)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestConvertToString(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		fieldName   string
		expectError bool
		expected    string
	}{
		{
			name:        "Simple string",
			value:       "hello",
			fieldName:   "name",
			expectError: false,
			expected:    "hello",
		},
		{
			name:        "Map to JSON string",
			value:       map[string]interface{}{"key": "value", "num": 42},
			fieldName:   "params",
			expectError: false,
			expected:    `{"key":"value","num":42}`,
		},
		{
			name:        "Array to JSON string",
			value:       []interface{}{"item1", "item2"},
			fieldName:   "external_params",
			expectError: false,
			expected:    `["item1","item2"]`,
		},
		{
			name: "Cells conversion",
			value: []interface{}{
				map[string]interface{}{
					"name":    "test_cell",
					"type":    "OP_LANG",
					"content": "print('hello')",
				},
			},
			fieldName:   "cells",
			expectError: false,
			// Should contain internal model format with enabled defaulting to true
			expected: `[{"description":"","enabled":true,"name":"test_cell","op":"print('hello')","secret_aware":false}]`,
		},
		{
			name:        "Non-JSON field with map should error",
			value:       map[string]interface{}{"key": "value"},
			fieldName:   "regular_field",
			expectError: true,
		},
		{
			name:        "Invalid type",
			value:       123,
			fieldName:   "name",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := convertToString(tt.value, tt.fieldName)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if strings.Contains(tt.expected, "{") || strings.Contains(tt.expected, "[") {
					// For JSON, check contains instead of exact match
					assert.Contains(t, result.ValueString(), tt.expected)
				} else {
					assert.Equal(t, tt.expected, result.ValueString())
				}
			}
		})
	}
}

func TestConvertCellsToInternalModel(t *testing.T) {
	tests := []struct {
		name        string
		inputCells  interface{}
		expectError bool
		validate    func(t *testing.T, result types.String)
	}{
		{
			name: "Valid OP_LANG cells",
			inputCells: []interface{}{
				map[string]interface{}{
					"name":    "cell1",
					"type":    "OP_LANG",
					"content": "print('test')",
					"enabled": true,
				},
				map[string]interface{}{
					"name":    "cell2",
					"type":    "OP_LANG",
					"content": "print('hello')",
					"enabled": false,
				},
			},
			expectError: false,
			validate: func(t *testing.T, result types.String) {
				assert.Contains(t, result.ValueString(), `"op":"print('test')"`)
				assert.Contains(t, result.ValueString(), `"op":"print('hello')"`)
				assert.Contains(t, result.ValueString(), `"name":"cell1"`)
				assert.Contains(t, result.ValueString(), `"name":"cell2"`)
			},
		},
		{
			name: "Valid MARKDOWN cells",
			inputCells: []interface{}{
				map[string]interface{}{
					"name":    "md_cell",
					"type":    "MARKDOWN",
					"content": "# Header",
					"enabled": true,
				},
			},
			expectError: false,
			validate: func(t *testing.T, result types.String) {
				assert.Contains(t, result.ValueString(), `"md":"# Header"`)
				assert.Contains(t, result.ValueString(), `"name":"md_cell"`)
			},
		},
		{
			name: "Mixed cell types",
			inputCells: []interface{}{
				map[string]interface{}{
					"name":    "op_cell",
					"type":    "OP_LANG",
					"content": "code",
				},
				map[string]interface{}{
					"name":    "md_cell",
					"type":    "MARKDOWN",
					"content": "text",
				},
			},
			expectError: false,
			validate: func(t *testing.T, result types.String) {
				assert.Contains(t, result.ValueString(), `"op":"code"`)
				assert.Contains(t, result.ValueString(), `"md":"text"`)
			},
		},
		{
			name:        "Not an array",
			inputCells:  "not_an_array",
			expectError: true,
		},
		{
			name: "Invalid cell object",
			inputCells: []interface{}{
				"not_a_map",
			},
			expectError: true,
		},
		{
			name:        "Empty cells array",
			inputCells:  []interface{}{},
			expectError: false,
			validate: func(t *testing.T, result types.String) {
				assert.Equal(t, "[]", result.ValueString())
			},
		},
		{
			name: "Cell with cell_type field",
			inputCells: []interface{}{
				map[string]interface{}{
					"name":      "cell1",
					"cell_type": "OP_LANG",
					"content":   "print('test')",
				},
			},
			expectError: false,
			validate: func(t *testing.T, result types.String) {
				assert.Contains(t, result.ValueString(), `"op":"print('test')"`)
			},
		},
		{
			name: "Cell without enabled field should default to true",
			inputCells: []interface{}{
				map[string]interface{}{
					"name":    "cell1",
					"type":    "OP_LANG",
					"content": "echo hello",
				},
				map[string]interface{}{
					"name":    "cell2",
					"type":    "MARKDOWN",
					"content": "# Title",
				},
			},
			expectError: false,
			validate: func(t *testing.T, result types.String) {
				// Verify both cells have enabled=true by default
				assert.Contains(t, result.ValueString(), `"enabled":true`)
				// Count occurrences - should appear twice (once for each cell)
				count := strings.Count(result.ValueString(), `"enabled":true`)
				assert.Equal(t, 2, count, "Both cells should have enabled=true")
			},
		},
		{
			name: "Cell with explicit enabled=false should remain false",
			inputCells: []interface{}{
				map[string]interface{}{
					"name":    "cell1",
					"type":    "OP_LANG",
					"content": "echo hello",
					"enabled": false,
				},
			},
			expectError: false,
			validate: func(t *testing.T, result types.String) {
				assert.Contains(t, result.ValueString(), `"enabled":false`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := convertCellsToInternalModel(tt.inputCells)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}
