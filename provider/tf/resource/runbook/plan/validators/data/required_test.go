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

	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		tfModel     *model.RunbookTFModel
		dataMap     map[string]any
		expectError bool
		errorMsg    string
	}{
		{
			name: "name set in TF model",
			tfModel: &model.RunbookTFModel{
				Name: types.StringValue("my_runbook"),
			},
			dataMap:     map[string]any{},
			expectError: false,
		},
		{
			name: "name set in data JSON",
			tfModel: &model.RunbookTFModel{
				Name: types.StringNull(),
			},
			dataMap:     map[string]any{"name": "my_runbook"},
			expectError: false,
		},
		{
			name: "name set in both",
			tfModel: &model.RunbookTFModel{
				Name: types.StringValue("tf_name"),
			},
			dataMap:     map[string]any{"name": "data_name"},
			expectError: false,
		},
		{
			name: "name not set anywhere",
			tfModel: &model.RunbookTFModel{
				Name: types.StringNull(),
			},
			dataMap:     map[string]any{},
			expectError: true,
			errorMsg:    "\"name\" is required",
		},
		{
			name: "name unknown in TF model and not in data",
			tfModel: &model.RunbookTFModel{
				Name: types.StringUnknown(),
			},
			dataMap:     map[string]any{},
			expectError: true,
			errorMsg:    "\"name\" is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequiredFields(tt.tfModel, tt.dataMap)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsFieldSet(t *testing.T) {
	tests := []struct {
		name       string
		modelValue types.String
		dataMap    map[string]any
		fieldName  string
		expected   bool
	}{
		{
			name:       "set in model only",
			modelValue: types.StringValue("value"),
			dataMap:    map[string]any{},
			fieldName:  "name",
			expected:   true,
		},
		{
			name:       "set in data only",
			modelValue: types.StringNull(),
			dataMap:    map[string]any{"name": "value"},
			fieldName:  "name",
			expected:   true,
		},
		{
			name:       "set in both",
			modelValue: types.StringValue("tf"),
			dataMap:    map[string]any{"name": "data"},
			fieldName:  "name",
			expected:   true,
		},
		{
			name:       "set in neither",
			modelValue: types.StringNull(),
			dataMap:    map[string]any{},
			fieldName:  "name",
			expected:   false,
		},
		{
			name:       "model unknown counts as not set",
			modelValue: types.StringUnknown(),
			dataMap:    map[string]any{},
			fieldName:  "name",
			expected:   false,
		},
		{
			name:       "model unknown but data has it",
			modelValue: types.StringUnknown(),
			dataMap:    map[string]any{"name": "from_data"},
			fieldName:  "name",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFieldSet(tt.modelValue, tt.fieldName, tt.dataMap)
			assert.Equal(t, tt.expected, result)
		})
	}
}
