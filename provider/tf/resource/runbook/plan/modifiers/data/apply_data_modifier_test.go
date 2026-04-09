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
	"testing"

	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyDataModifier(t *testing.T) {
	tests := []struct {
		name        string
		config      *model.RunbookTFModel
		expectError bool
		validate    func(t *testing.T, result *model.RunbookTFModel)
	}{
		{
			name: "Apply data JSON to empty fields",
			config: &model.RunbookTFModel{
				Name: types.StringNull(),
				Data: types.StringValue(`{
					"name": "from_data",
					"enabled": true,
					"description": "Applied from data"
				}`),
			},
			expectError: false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				// Result should have data values applied
				assert.Equal(t, "from_data", result.Name.ValueString())
				assert.Equal(t, true, result.Enabled.ValueBool())
				assert.Equal(t, "Applied from data", result.Description.ValueString())
				// Original config should be unchanged
			},
		},
		{
			name: "Don't override existing values",
			config: &model.RunbookTFModel{
				Name:    types.StringValue("original"),
				Enabled: types.BoolValue(false),
				Data: types.StringValue(`{
					"name": "from_data",
					"enabled": true
				}`),
			},
			expectError: false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				// Should keep original values
				assert.Equal(t, "original", result.Name.ValueString())
				assert.Equal(t, false, result.Enabled.ValueBool())
			},
		},
		{
			name: "Apply complex fields from data",
			config: &model.RunbookTFModel{
				Cells:           types.StringNull(),
				Params:          types.StringNull(),
				AllowedEntities: types.ListNull(types.StringType),
				Data: types.StringValue(`{
					"cells": [{"name": "cell1", "type": "OP_LANG", "content": "code"}],
					"params": [{"name": "p1", "value": "v1"}],
					"allowed_entities": ["entity1", "entity2"]
				}`),
			},
			expectError: false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.False(t, result.CellsList.IsNull())
				require.Equal(t, 1, len(result.CellsList.Elements()))
				cellObj := result.CellsList.Elements()[0].(types.Object)
				assert.Equal(t, "code", cellObj.Attributes()["op"].(types.String).ValueString())
				assert.Equal(t, "cell1", cellObj.Attributes()["name"].(types.String).ValueString())
				assert.Equal(t, true, cellObj.Attributes()["enabled"].(types.Bool).ValueBool())
				assert.Equal(t, false, cellObj.Attributes()["secret_aware"].(types.Bool).ValueBool())

				assert.False(t, result.ParamsList.IsNull())
				require.Equal(t, 1, len(result.ParamsList.Elements()))
				paramObj := result.ParamsList.Elements()[0].(types.Object)
				assert.Equal(t, "p1", paramObj.Attributes()["name"].(types.String).ValueString())
				assert.Equal(t, "v1", paramObj.Attributes()["value"].(types.String).ValueString())
				assert.True(t, result.Params.IsNull(), "deprecated params string should remain null")

				assert.Equal(t, result.AllowedEntities.Elements()[0], types.StringValue("entity1"))
				assert.Equal(t, result.AllowedEntities.Elements()[1], types.StringValue("entity2"))
				assert.Equal(t, 2, len(result.AllowedEntities.Elements()))
			},
		},
		{
			name: "Handle camelCase fields",
			config: &model.RunbookTFModel{
				TimeoutMs:              types.Int64Null(),
				IsRunOutputPersisted:   types.BoolNull(),
				FilterResourceToAction: types.BoolNull(),
				Data: types.StringValue(`{
					"timeoutMs": 3000,
					"isRunOutputPersisted": true,
					"filterResourceToAction": false
				}`),
			},
			expectError: false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.Equal(t, int64(3000), result.TimeoutMs.ValueInt64())
				assert.Equal(t, true, result.IsRunOutputPersisted.ValueBool())
				assert.Equal(t, false, result.FilterResourceToAction.ValueBool())
			},
		},
		{
			name: "Empty data JSON",
			config: &model.RunbookTFModel{
				Name: types.StringNull(),
				Data: types.StringValue(`{}`),
			},
			expectError: false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.True(t, result.Name.IsNull())
			},
		},
		{
			name: "Null data JSON",
			config: &model.RunbookTFModel{
				Name: types.StringNull(),
				Data: types.StringNull(),
			},
			expectError: false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.True(t, result.Name.IsNull())
			},
		},
		{
			name: "Invalid data JSON",
			config: &model.RunbookTFModel{
				Data: types.StringValue(`{invalid json}`),
			},
			expectError: true,
		},
		{
			name: "Original config unchanged",
			config: &model.RunbookTFModel{
				Name: types.StringValue("original"),
				Data: types.StringValue(`{"name": "modified"}`),
			},
			expectError: false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				// Result should have original value (not overridden)
				assert.Equal(t, "original", result.Name.ValueString())
			},
		},
		{
			name: "Mixed null and set values",
			config: &model.RunbookTFModel{
				Name:        types.StringValue("keep_this"),
				Description: types.StringNull(),
				Enabled:     types.BoolNull(),
				TimeoutMs:   types.Int64Value(1000),
				Data: types.StringValue(`{
					"name": "override_attempt",
					"description": "apply_this",
					"enabled": true,
					"timeout_ms": 5000
				}`),
			},
			expectError: false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.Equal(t, "keep_this", result.Name.ValueString())
				assert.Equal(t, "apply_this", result.Description.ValueString())
				assert.Equal(t, true, result.Enabled.ValueBool())
				assert.Equal(t, int64(1000), result.TimeoutMs.ValueInt64())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()

			// Create a copy of the original config to verify it's unchanged
			originalName := tt.config.Name
			originalEnabled := tt.config.Enabled
			originalData := tt.config.Data

			// when
			result, err := ApplyDataModifier(ctx, tt.config)

			// then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				// Verify original config is unchanged
				assert.Equal(t, originalName, tt.config.Name)
				assert.Equal(t, originalEnabled, tt.config.Enabled)
				assert.Equal(t, originalData, tt.config.Data)

				// Verify result is a different instance
				assert.NotSame(t, tt.config, result)

				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestApplyDataModifier_ShallowCopy(t *testing.T) {
	// given
	ctx := context.Background()
	config := &model.RunbookTFModel{
		Name:    types.StringValue("original"),
		Enabled: types.BoolValue(true),
		Data:    types.StringValue(`{"description": "from_data"}`),
		Labels: func() types.List {
			s, _ := types.ListValueFrom(ctx, types.StringType, []string{"label1", "label2"})
			return s
		}(),
	}

	// when
	result, err := ApplyDataModifier(ctx, config)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify it's a shallow copy
	assert.NotSame(t, config, result)
	assert.Equal(t, config.Name, result.Name)
	assert.Equal(t, config.Enabled, result.Enabled)
	assert.Equal(t, config.Labels, result.Labels)

	// Data JSON should be applied
	assert.Equal(t, "from_data", result.Description.ValueString())
}

func TestApplyDataModifier_AllFieldTypes(t *testing.T) {
	// given
	ctx := context.Background()
	config := &model.RunbookTFModel{
		// All fields null to test application
		Name:                                types.StringNull(),
		Enabled:                             types.BoolNull(),
		Description:                         types.StringNull(),
		TimeoutMs:                           types.Int64Null(),
		AllowedResourcesQuery:               types.StringNull(),
		CommunicationWorkspace:              types.StringNull(),
		CommunicationChannel:                types.StringNull(),
		IsRunOutputPersisted:                types.BoolNull(),
		FilterResourceToAction:              types.BoolNull(),
		CommunicationCudNotifications:       types.BoolNull(),
		CommunicationApprovalNotifications:  types.BoolNull(),
		CommunicationExecutionNotifications: types.BoolNull(),
		AllowedEntities:                     types.ListNull(types.StringType),
		Approvers:                           types.ListNull(types.StringType),
		Labels:                              types.ListNull(types.StringType),
		Editors:                             types.ListNull(types.StringType),
		SecretNames:                         types.ListNull(types.StringType),
		Cells:                               types.StringNull(),
		Params:                              types.StringNull(),
		ExternalParams:                      types.StringNull(),
		Data: types.StringValue(`{
			"name": "test_runbook",
			"enabled": true,
			"description": "Test description",
			"timeout_ms": 5000,
			"allowed_resources_query": "query",
			"communication_workspace": "workspace",
			"communication_channel": "channel",
			"is_run_output_persisted": true,
			"filter_resource_to_action": false,
			"communication_cud_notifications": true,
			"communication_approval_notifications": false,
			"communication_execution_notifications": true,
			"allowed_entities": ["e1", "e2"],
			"approvers": ["a1", "a2"],
			"labels": ["l1", "l2"],
			"editors": ["ed1", "ed2"],
			"secret_names": ["s1", "s2"],
			"params_groups": {
				"required": ["p1", "p2"],
				"optional": ["p3", "p4"],
				"exported": ["p5", "p6"],
				"external": ["p7", "p8"]
			},
			"cells": [{"name": "c1", "type": "OP_LANG", "content": "code"}],
			"params": [{"name": "p1", "value": "v1"}],
			"external_params": [{"name": "ep1", "source": "src"}]
		}`),
	}

	// when
	result, err := ApplyDataModifier(ctx, config)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify all fields were applied
	assert.Equal(t, "test_runbook", result.Name.ValueString())
	assert.Equal(t, true, result.Enabled.ValueBool())
	assert.Equal(t, "Test description", result.Description.ValueString())
	assert.Equal(t, int64(5000), result.TimeoutMs.ValueInt64())
	assert.Equal(t, "query", result.AllowedResourcesQuery.ValueString())
	assert.Equal(t, "workspace", result.CommunicationWorkspace.ValueString())
	assert.Equal(t, "channel", result.CommunicationChannel.ValueString())
	assert.Equal(t, true, result.IsRunOutputPersisted.ValueBool())
	assert.Equal(t, false, result.FilterResourceToAction.ValueBool())
	assert.Equal(t, true, result.CommunicationCudNotifications.ValueBool())
	assert.Equal(t, false, result.CommunicationApprovalNotifications.ValueBool())
	assert.Equal(t, true, result.CommunicationExecutionNotifications.ValueBool())

	// Check sets
	assert.Equal(t, 2, len(result.AllowedEntities.Elements()))
	assert.Equal(t, result.AllowedEntities.Elements()[0], types.StringValue("e1"))
	assert.Equal(t, result.AllowedEntities.Elements()[1], types.StringValue("e2"))

	assert.Equal(t, 2, len(result.Approvers.Elements()))
	assert.Equal(t, result.Approvers.Elements()[0], types.StringValue("a1"))
	assert.Equal(t, result.Approvers.Elements()[1], types.StringValue("a2"))

	assert.Equal(t, 2, len(result.Labels.Elements()))
	assert.Equal(t, result.Labels.Elements()[0], types.StringValue("l1"))
	assert.Equal(t, result.Labels.Elements()[1], types.StringValue("l2"))

	assert.Equal(t, 2, len(result.Editors.Elements()))
	assert.Equal(t, result.Editors.Elements()[0], types.StringValue("ed1"))
	assert.Equal(t, result.Editors.Elements()[1], types.StringValue("ed2"))

	assert.Equal(t, 2, len(result.SecretNames.Elements()))
	assert.Equal(t, result.SecretNames.Elements()[0], types.StringValue("s1"))
	assert.Equal(t, result.SecretNames.Elements()[1], types.StringValue("s2"))

	// Check JSON fields - cells from data go into cells_list (not deprecated cells string)
	assert.False(t, result.CellsList.IsNull())
	require.Equal(t, 1, len(result.CellsList.Elements()))
	cellObj := result.CellsList.Elements()[0].(types.Object)
	assert.Equal(t, "code", cellObj.Attributes()["op"].(types.String).ValueString())
	assert.Equal(t, "c1", cellObj.Attributes()["name"].(types.String).ValueString())
	assert.Equal(t, true, cellObj.Attributes()["enabled"].(types.Bool).ValueBool())
	assert.Equal(t, false, cellObj.Attributes()["secret_aware"].(types.Bool).ValueBool())
	assert.True(t, result.Cells.IsNull(), "cells string should remain null when cells come from data JSON")

	// params from data go into params_list (not deprecated params string)
	assert.False(t, result.ParamsList.IsNull())
	require.Equal(t, 1, len(result.ParamsList.Elements()))
	paramObj := result.ParamsList.Elements()[0].(types.Object)
	assert.Equal(t, "p1", paramObj.Attributes()["name"].(types.String).ValueString())
	assert.Equal(t, "v1", paramObj.Attributes()["value"].(types.String).ValueString())
	assert.True(t, result.Params.IsNull(), "params string should remain null when params come from data JSON")

	// external_params from data go into external_params_list (not deprecated external_params string)
	assert.False(t, result.ExternalParamsList.IsNull())
	require.Equal(t, 1, len(result.ExternalParamsList.Elements()))
	epObj := result.ExternalParamsList.Elements()[0].(types.Object)
	assert.Equal(t, "ep1", epObj.Attributes()["name"].(types.String).ValueString())
	assert.Equal(t, "src", epObj.Attributes()["source"].(types.String).ValueString())
	assert.True(t, result.ExternalParams.IsNull(), "external_params string should remain null when external_params come from data JSON")
}
