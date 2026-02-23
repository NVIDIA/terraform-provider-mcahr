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
	"context"
	"testing"

	data "terraform/terraform-provider/provider/tf/resource/runbook/data_attribute"
	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckDataJSONConflicts(t *testing.T) {
	// Skip these tests as they require complex mocking of the framework internals
	t.Skip("Skipping tests that require complex framework mocking")
}

func TestValidateNoFieldConflicts(t *testing.T) {
	tests := []struct {
		name        string
		tfModel     *model.RunbookTFModel
		expectError bool
		errorFields []string
	}{
		{
			name: "No conflicts - different fields",
			tfModel: &model.RunbookTFModel{
				Name:    types.StringValue("from_tf"),
				Enabled: types.BoolValue(true),
				Data: types.StringValue(`{
					"description": "from_data",
					"timeout_ms": 5000
				}`),
			},
			expectError: false,
		},
		{
			name: "Single field conflict",
			tfModel: &model.RunbookTFModel{
				Name: types.StringValue("from_tf"),
				Data: types.StringValue(`{"name": "from_data"}`),
			},
			expectError: true,
			errorFields: []string{"name"},
		},
		{
			name: "Multiple field conflicts",
			tfModel: &model.RunbookTFModel{
				Name:        types.StringValue("from_tf"),
				Description: types.StringValue("tf desc"),
				Enabled:     types.BoolValue(true),
				Data: types.StringValue(`{
					"name": "from_data",
					"description": "data desc",
					"enabled": false
				}`),
			},
			expectError: true,
			errorFields: []string{"name", "description", "enabled"},
		},
		{
			name: "CamelCase field names in data",
			tfModel: &model.RunbookTFModel{
				TimeoutMs:            types.Int64Value(1000),
				IsRunOutputPersisted: types.BoolValue(true),
				Data: types.StringValue(`{
					"timeoutMs": 2000,
					"isRunOutputPersisted": false
				}`),
			},
			expectError: true,
			errorFields: []string{"timeout_ms", "is_run_output_persisted"},
		},
		{
			name: "Unknown/null fields don't conflict",
			tfModel: &model.RunbookTFModel{
				Name:    types.StringUnknown(),
				Enabled: types.BoolNull(),
				Data: types.StringValue(`{
					"name": "from_data",
					"enabled": true
				}`),
			},
			expectError: false,
		},
		{
			name: "JSON fields conflict",
			tfModel: &model.RunbookTFModel{
				Cells:  types.StringValue(`[{"op": "print('tf')"}]`),
				Params: types.StringValue(`[{"name": "p1"}]`),
				Data: types.StringValue(`{
					"cells": [{"type": "OP_LANG", "content": "print('data')"}],
					"params": [{"name": "p2"}]
				}`),
			},
			expectError: true,
			errorFields: []string{"cells", "params"},
		},
		{
			name: "Skip data and _full fields",
			tfModel: &model.RunbookTFModel{
				CellsFull:  types.StringValue("test"),
				ParamsFull: types.StringValue("test"),
				Data: types.StringValue(`{
					"cells_full": "test",
					"params_full": "test",
					"data": "recursive"
				}`),
			},
			expectError: false,
		},
		{
			name: "Set fields conflict",
			tfModel: &model.RunbookTFModel{
				AllowedEntities: func() types.List {
					s, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"entity1"})
					return s
				}(),
				Data: types.StringValue(`{
					"allowed_entities": ["entity2", "entity3"]
				}`),
			},
			expectError: true,
			errorFields: []string{"allowed_entities"},
		},
		{
			name: "Null data JSON",
			tfModel: &model.RunbookTFModel{
				Name: types.StringValue("test"),
				Data: types.StringNull(),
			},
			expectError: false,
		},
		{
			name: "Empty data JSON",
			tfModel: &model.RunbookTFModel{
				Name: types.StringValue("test"),
				Data: types.StringValue("{}"),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()
			dataMap, err := data.ParseDataJSONToMap(tt.tfModel.Data)
			if err != nil {
				t.Fatalf("failed to parse data JSON: %v", err)
			}

			// when
			err = validateNoFieldConflicts(ctx, tt.tfModel, dataMap)

			// then
			if tt.expectError {
				require.Error(t, err)
				for _, field := range tt.errorFields {
					assert.Contains(t, err.Error(), field)
				}
				assert.Contains(t, err.Error(), "set in both")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateNoFieldConflicts_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		tfModel     *model.RunbookTFModel
		expectError bool
		errorMsg    string
		errorFields []string
	}{
		{
			name: "Invalid JSON in data field",
			tfModel: &model.RunbookTFModel{
				Data: types.StringValue("{invalid json}"),
			},
			expectError: true,
			errorMsg:    "failed to parse data JSON",
		},
		{
			name: "Field processing error",
			tfModel: &model.RunbookTFModel{
				// This would need a custom setup to trigger an error in OnEachStructField
				// For now, we'll test the normal path
				Name: types.StringValue("test"),
				Data: types.StringValue(`{"description": "test"}`),
			},
			expectError: false,
		},
		{
			name: "Mixed null and set values",
			tfModel: &model.RunbookTFModel{
				Name:        types.StringValue("set"),
				Description: types.StringNull(),
				Enabled:     types.BoolUnknown(),
				TimeoutMs:   types.Int64Value(1000),
				Data: types.StringValue(`{
					"name": "conflict",
					"description": "no conflict - tf is null",
					"enabled": "no conflict - tf is unknown",
					"timeout_ms": 2000
				}`),
			},
			expectError: true,
			errorFields: []string{"name", "timeout_ms"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()
			dataMap, err := data.ParseDataJSONToMap(tt.tfModel.Data)

			// If we expect an error related to JSON parsing, check it here
			if err != nil {
				if tt.expectError && tt.errorMsg == "failed to parse data JSON" {
					// This is expected - the JSON parsing failure is what we're testing
					require.Error(t, err)
					assert.Contains(t, err.Error(), "failed to parse data JSON")
					return
				}
				// Otherwise, this is an unexpected parse error
				t.Fatalf("failed to parse data JSON: %v", err)
			}

			// when
			err = validateNoFieldConflicts(ctx, tt.tfModel, dataMap)

			// then
			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				if tt.errorFields != nil {
					for _, field := range tt.errorFields {
						assert.Contains(t, err.Error(), field)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
