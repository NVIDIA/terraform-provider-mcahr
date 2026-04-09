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

package datavalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- validateArrayField ---

func TestValidateArrayField_FieldNotPresent(t *testing.T) {
	t.Parallel()
	err := validateArrayField("cells", map[string]any{}, func(i int, obj map[string]interface{}) error { return nil })
	assert.NoError(t, err)
}

func TestValidateArrayField_NotAnArray(t *testing.T) {
	t.Parallel()
	err := validateArrayField("cells", map[string]any{"cells": "not-an-array"}, func(i int, obj map[string]interface{}) error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be an array")
}

func TestValidateArrayField_ElementNotAnObject(t *testing.T) {
	t.Parallel()
	dataMap := map[string]any{"cells": []interface{}{"not-an-object"}}
	err := validateArrayField("cells", dataMap, func(i int, obj map[string]interface{}) error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be an object")
}

func TestValidateArrayField_EmptyArray(t *testing.T) {
	t.Parallel()
	dataMap := map[string]any{"cells": []interface{}{}}
	err := validateArrayField("cells", dataMap, func(i int, obj map[string]interface{}) error { return nil })
	assert.NoError(t, err)
}

// --- validateDataCells ---

func TestValidateDataCells(t *testing.T) {
	tests := []struct {
		name        string
		dataMap     map[string]any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no cells field",
			dataMap:     map[string]any{},
			expectError: false,
		},
		{
			name: "valid OP_LANG cell",
			dataMap: map[string]any{
				"cells": []interface{}{
					map[string]interface{}{"type": "OP_LANG", "content": "print('hello')"},
				},
			},
			expectError: false,
		},
		{
			name: "valid MARKDOWN cell",
			dataMap: map[string]any{
				"cells": []interface{}{
					map[string]interface{}{"type": "MARKDOWN", "content": "# Header"},
				},
			},
			expectError: false,
		},
		{
			name: "multiple valid cells",
			dataMap: map[string]any{
				"cells": []interface{}{
					map[string]interface{}{"type": "OP_LANG", "content": "code"},
					map[string]interface{}{"type": "MARKDOWN", "content": "text"},
				},
			},
			expectError: false,
		},
		{
			name: "invalid cell type",
			dataMap: map[string]any{
				"cells": []interface{}{
					map[string]interface{}{"type": "INVALID_TYPE", "content": "code"},
				},
			},
			expectError: true,
			errorMsg:    "invalid or missing \"type\"",
		},
		{
			name: "missing cell type",
			dataMap: map[string]any{
				"cells": []interface{}{
					map[string]interface{}{"content": "code"},
				},
			},
			expectError: true,
			errorMsg:    "invalid or missing \"type\"",
		},
		{
			name: "empty cell type",
			dataMap: map[string]any{
				"cells": []interface{}{
					map[string]interface{}{"type": "", "content": "code"},
				},
			},
			expectError: true,
			errorMsg:    "invalid or missing \"type\"",
		},
		{
			name: "first valid second invalid",
			dataMap: map[string]any{
				"cells": []interface{}{
					map[string]interface{}{"type": "OP_LANG", "content": "ok"},
					map[string]interface{}{"type": "BAD"},
				},
			},
			expectError: true,
			errorMsg:    "cells[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDataCells(tt.dataMap)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- validateDataParams ---

func TestValidateDataParams(t *testing.T) {
	tests := []struct {
		name        string
		dataMap     map[string]any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no params field",
			dataMap:     map[string]any{},
			expectError: false,
		},
		{
			name: "valid param",
			dataMap: map[string]any{
				"params": []interface{}{
					map[string]interface{}{"name": "param_1", "value": "val"},
				},
			},
			expectError: false,
		},
		{
			name: "multiple valid params",
			dataMap: map[string]any{
				"params": []interface{}{
					map[string]interface{}{"name": "param_1", "value": "v1"},
					map[string]interface{}{"name": "param_2", "value": "v2"},
				},
			},
			expectError: false,
		},
		{
			name: "missing name",
			dataMap: map[string]any{
				"params": []interface{}{
					map[string]interface{}{"value": "val"},
				},
			},
			expectError: true,
			errorMsg:    "missing the required \"name\"",
		},
		{
			name: "empty name",
			dataMap: map[string]any{
				"params": []interface{}{
					map[string]interface{}{"name": "", "value": "val"},
				},
			},
			expectError: true,
			errorMsg:    "missing the required \"name\"",
		},
		{
			name: "invalid name — starts with digit",
			dataMap: map[string]any{
				"params": []interface{}{
					map[string]interface{}{"name": "1param", "value": "val"},
				},
			},
			expectError: true,
			errorMsg:    "alphanumeric",
		},
		{
			name: "invalid name — special characters",
			dataMap: map[string]any{
				"params": []interface{}{
					map[string]interface{}{"name": "my-param", "value": "val"},
				},
			},
			expectError: true,
			errorMsg:    "alphanumeric",
		},
		{
			name: "second param invalid",
			dataMap: map[string]any{
				"params": []interface{}{
					map[string]interface{}{"name": "good_name", "value": "v1"},
					map[string]interface{}{"value": "v2"},
				},
			},
			expectError: true,
			errorMsg:    "params[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDataParams(tt.dataMap)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- validateDataExternalParams ---

func TestValidateDataExternalParams(t *testing.T) {
	tests := []struct {
		name        string
		dataMap     map[string]any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no external_params field",
			dataMap:     map[string]any{},
			expectError: false,
		},
		{
			name: "valid external param",
			dataMap: map[string]any{
				"external_params": []interface{}{
					map[string]interface{}{"name": "ep_1", "source": "alertmanager", "json_path": "$.data"},
				},
			},
			expectError: false,
		},
		{
			name: "missing name",
			dataMap: map[string]any{
				"external_params": []interface{}{
					map[string]interface{}{"source": "alertmanager", "json_path": "$.data"},
				},
			},
			expectError: true,
			errorMsg:    "missing the required \"name\"",
		},
		{
			name: "invalid name",
			dataMap: map[string]any{
				"external_params": []interface{}{
					map[string]interface{}{"name": "bad-name", "source": "alertmanager"},
				},
			},
			expectError: true,
			errorMsg:    "alphanumeric",
		},
		{
			name: "missing source",
			dataMap: map[string]any{
				"external_params": []interface{}{
					map[string]interface{}{"name": "ep_1", "json_path": "$.data"},
				},
			},
			expectError: true,
			errorMsg:    "missing the required \"source\"",
		},
		{
			name: "empty source",
			dataMap: map[string]any{
				"external_params": []interface{}{
					map[string]interface{}{"name": "ep_1", "source": ""},
				},
			},
			expectError: true,
			errorMsg:    "missing the required \"source\"",
		},
		{
			name: "invalid source",
			dataMap: map[string]any{
				"external_params": []interface{}{
					map[string]interface{}{"name": "ep_1", "source": "invalid_source"},
				},
			},
			expectError: true,
			errorMsg:    "invalid \"source\"",
		},
		{
			name: "valid name but invalid source",
			dataMap: map[string]any{
				"external_params": []interface{}{
					map[string]interface{}{"name": "good_name", "source": "bad"},
				},
			},
			expectError: true,
			errorMsg:    "invalid \"source\"",
		},
		{
			name: "multiple — first valid second missing source",
			dataMap: map[string]any{
				"external_params": []interface{}{
					map[string]interface{}{"name": "ep_1", "source": "alertmanager"},
					map[string]interface{}{"name": "ep_2"},
				},
			},
			expectError: true,
			errorMsg:    "external_params[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDataExternalParams(tt.dataMap)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- helper function tests ---

func TestIsValidCellType(t *testing.T) {
	assert.True(t, isValidCellType("OP_LANG"))
	assert.True(t, isValidCellType("MARKDOWN"))
	assert.False(t, isValidCellType("INVALID"))
	assert.False(t, isValidCellType(""))
	assert.False(t, isValidCellType("op_lang"))
}

func TestIsValidExternalParamSource(t *testing.T) {
	assert.True(t, isValidExternalParamSource("alertmanager"))
	assert.False(t, isValidExternalParamSource("invalid"))
	assert.False(t, isValidExternalParamSource(""))
	assert.False(t, isValidExternalParamSource("ALERTMANAGER"))
}
