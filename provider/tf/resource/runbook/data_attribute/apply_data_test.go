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
	"strings"
	"testing"

	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyDataJSONValues(t *testing.T) {
	tests := []struct {
		name        string
		input       *model.RunbookTFModel
		expectError bool
		validate    func(t *testing.T, tfModel *model.RunbookTFModel)
	}{
		{
			name: "Apply simple fields from data JSON",
			input: &model.RunbookTFModel{
				Data: types.StringValue(`{
					"name": "test_runbook",
					"enabled": true,
					"description": "Test description",
					"timeout_ms": 5000
				}`),
				// These fields are not set in TF model
				Name:        types.StringNull(),
				Enabled:     types.BoolNull(),
				Description: types.StringNull(),
				TimeoutMs:   types.Int64Null(),
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.Equal(t, "test_runbook", tfModel.Name.ValueString())
				assert.Equal(t, true, tfModel.Enabled.ValueBool())
				assert.Equal(t, "Test description", tfModel.Description.ValueString())
				assert.Equal(t, int64(5000), tfModel.TimeoutMs.ValueInt64())
			},
		},
		{
			name: "Don't override existing TF model values",
			input: &model.RunbookTFModel{
				Data: types.StringValue(`{
					"name": "from_data",
					"enabled": false
				}`),
				// These are already set in TF model
				Name:    types.StringValue("from_tf"),
				Enabled: types.BoolValue(true),
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				// Should keep TF values
				assert.Equal(t, "from_tf", tfModel.Name.ValueString())
				assert.Equal(t, true, tfModel.Enabled.ValueBool())
			},
		},
		{
			name: "Apply JSON fields from data",
			input: &model.RunbookTFModel{
				Data: types.StringValue(`{
					"params": [{"name": "param1", "value": "value1"}],
					"cells": [{"name": "cell1", "type": "OP_LANG", "content": "print('hello')"}]
				}`),
				Params: types.StringNull(),
				Cells:  types.StringNull(),
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.True(t, tfModel.Params.ValueString() != "")
				assert.True(t, tfModel.Cells.IsNull(), "cells string should remain null, data cells go to cells_list")
				assert.False(t, tfModel.CellsList.IsNull(), "cells from data should populate cells_list")
				require.Equal(t, 1, len(tfModel.CellsList.Elements()))
				cellObj := tfModel.CellsList.Elements()[0].(types.Object)
				assert.Equal(t, "print('hello')", cellObj.Attributes()["op"].(types.String).ValueString())
			},
		},
		{
			name: "Handle camelCase fields in data JSON",
			input: &model.RunbookTFModel{
				Data: types.StringValue(`{
					"timeoutMs": 3000,
					"isRunOutputPersisted": true,
					"filterResourceToAction": false
				}`),
				TimeoutMs:              types.Int64Null(),
				IsRunOutputPersisted:   types.BoolNull(),
				FilterResourceToAction: types.BoolNull(),
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.Equal(t, int64(3000), tfModel.TimeoutMs.ValueInt64())
				assert.Equal(t, true, tfModel.IsRunOutputPersisted.ValueBool())
				assert.Equal(t, false, tfModel.FilterResourceToAction.ValueBool())
			},
		},
		{
			name: "Skip _full fields and data field",
			input: &model.RunbookTFModel{
				Data: types.StringValue(`{
					"cells_full": "[{invalid}]",
					"params_full": "[{invalid}]",
					"data": "should_be_skipped"
				}`),
				CellsFull:  types.StringNull(),
				ParamsFull: types.StringNull(),
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.True(t, tfModel.CellsFull.IsNull())
				assert.True(t, tfModel.ParamsFull.IsNull())
			},
		},
		{
			name: "Invalid data JSON",
			input: &model.RunbookTFModel{
				Data: types.StringValue(`{invalid json}`),
			},
			expectError: true,
		},
		{
			name: "Empty data JSON",
			input: &model.RunbookTFModel{
				Data: types.StringValue(`{}`),
				Name: types.StringNull(),
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.True(t, tfModel.Name.IsNull())
			},
		},
		{
			name: "Null data JSON",
			input: &model.RunbookTFModel{
				Data: types.StringNull(),
				Name: types.StringNull(),
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.True(t, tfModel.Name.IsNull())
			},
		},
		{
			name: "Apply set fields from data JSON",
			input: &model.RunbookTFModel{
				Data: types.StringValue(`{
					"allowed_entities": ["entity1", "entity2"],
					"approvers": ["user1", "user2"],
					"labels": ["label1", "label2"]
				}`),
				AllowedEntities: types.ListNull(types.StringType),
				Approvers:       types.ListNull(types.StringType),
				Labels:          types.ListNull(types.StringType),
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.False(t, tfModel.AllowedEntities.IsNull())
				assert.Equal(t, 2, len(tfModel.AllowedEntities.Elements()))
				assert.False(t, tfModel.Approvers.IsNull())
				assert.Equal(t, 2, len(tfModel.Approvers.Elements()))
				assert.False(t, tfModel.Labels.IsNull())
				assert.Equal(t, 2, len(tfModel.Labels.Elements()))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()

			// when
			err := ApplyDataJSONValues(ctx, tt.input)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, tt.input)
				}
			}
		})
	}
}

func TestDataJSONToTFModel(t *testing.T) {
	tests := []struct {
		name        string
		dataJSON    types.String
		expectError bool
		validate    func(t *testing.T, tfModel *model.RunbookTFModel)
	}{
		{
			name: "Create TF model from data JSON",
			dataJSON: types.StringValue(`{
				"name": "test_runbook",
				"enabled": true,
				"description": "Created from data",
				"timeout_ms": 2000,
				"params": [{"name": "p1", "value": "v1"}]
			}`),
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.Equal(t, "test_runbook", tfModel.Name.ValueString())
				assert.Equal(t, true, tfModel.Enabled.ValueBool())
				assert.Equal(t, "Created from data", tfModel.Description.ValueString())
				assert.Equal(t, int64(2000), tfModel.TimeoutMs.ValueInt64())
				assert.False(t, tfModel.Params.IsNull())
			},
		},
		{
			name:        "Invalid JSON",
			dataJSON:    types.StringValue(`{invalid}`),
			expectError: true,
		},
		{
			name:        "Empty JSON",
			dataJSON:    types.StringValue(`{}`),
			expectError: false,
			validate: func(t *testing.T, tfModel *model.RunbookTFModel) {
				assert.NotNil(t, tfModel)
				assert.Equal(t, `{}`, tfModel.Data.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()

			// when
			result, err := DataJSONToTFModel(ctx, tt.dataJSON)

			// then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.dataJSON.ValueString(), result.Data.ValueString())
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestSetFieldFromDataJSON(t *testing.T) {
	tests := []struct {
		name        string
		fieldName   string
		dataMap     map[string]interface{}
		expectError bool
	}{
		{
			name:      "Set string field",
			fieldName: "name",
			dataMap: map[string]interface{}{
				"name": "test_value",
			},
			expectError: false,
		},
		{
			name:      "Set bool field",
			fieldName: "enabled",
			dataMap: map[string]interface{}{
				"enabled": true,
			},
			expectError: false,
		},
		{
			name:      "Set int64 field",
			fieldName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeout_ms": float64(5000), // JSON numbers are float64
			},
			expectError: false,
		},
		{
			name:        "Field not in data map",
			fieldName:   "missing_field",
			dataMap:     map[string]interface{}{},
			expectError: false, // Should not error, just skip
		},
		{
			name:      "Invalid type conversion",
			fieldName: "timeout_ms",
			dataMap: map[string]interface{}{
				"timeout_ms": "not_a_number",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tfModel := &model.RunbookTFModel{
				Name:      types.StringNull(),
				Enabled:   types.BoolNull(),
				TimeoutMs: types.Int64Null(),
			}

			// Use reflection to get the field
			modelValue := reflect.ValueOf(tfModel).Elem()
			modelType := modelValue.Type()

			var fieldValue reflect.Value
			for i := 0; i < modelType.NumField(); i++ {
				field := modelType.Field(i)
				jsonTag := field.Tag.Get("json")
				if jsonTag != "" {
					fieldName := strings.Split(jsonTag, ",")[0]
					if fieldName == tt.fieldName {
						fieldValue = modelValue.Field(i)
						break
					}
				}
			}

			if !fieldValue.IsValid() {
				t.Skip("Field not found in model")
			}

			// when
			err := setFieldFromDataJSON(tt.fieldName, fieldValue, tt.dataMap)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
