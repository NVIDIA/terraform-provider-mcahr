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

package process

import (
	"context"
	"encoding/json"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"
	converters "terraform/terraform-provider/provider/tf/resource/runbook/translator/object_converters"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runbookModelGetter is a test helper that implements corecommon.Getter for RunbookTFModel.
type runbookModelGetter struct {
	model runbooktf.RunbookTFModel
}

func (g *runbookModelGetter) Get(_ context.Context, target interface{}) diag.Diagnostics {
	*(target.(*runbooktf.RunbookTFModel)) = g.model
	return nil
}

func (g *runbookModelGetter) GetAttribute(_ context.Context, _ path.Path, _ interface{}) diag.Diagnostics {
	return nil
}

// Mock implementation of JsonConfigurable for testing
type MockJsonConfigurable struct {
	Config common.JsonConfig
	Name   string `json:"name"`
	Type   string `json:"type"`
}

var _ common.JsonConfigurable = &MockJsonConfigurable{} // check that MockJsonConfigurable implements JsonConfigurable

func (m *MockJsonConfigurable) SetConfig(config common.JsonConfig) {
	m.Config = config
}

func (m *MockJsonConfigurable) GetConfig() common.JsonConfig {
	return m.Config
}

// TestPostProcessJsonFields_WithNullFields tests JSON field processing with null values
func TestPostProcessJsonFields_WithNullFields(t *testing.T) {
	t.Parallel()
	// given
	tfModel := &runbooktf.RunbookTFModel{
		Name:               types.StringValue("test_runbook"),
		ParamsFull:         types.StringNull(),
		ExternalParamsFull: types.StringNull(),
		CellsFull:          types.StringNull(),
	}
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	// when
	err := postProcessJsonFields(requestContext, tfModel)

	// then
	require.NoError(t, err)
	assert.True(t, tfModel.ParamsFull.IsNull())
	assert.True(t, tfModel.ExternalParamsFull.IsNull())
	assert.True(t, tfModel.CellsFull.IsNull())
}

// TestPostProcessJsonFields_WithUnknownFields tests JSON field processing with unknown values
func TestPostProcessJsonFields_WithUnknownFields(t *testing.T) {
	t.Parallel()
	// given
	tfModel := &runbooktf.RunbookTFModel{
		Name:               types.StringValue("test_runbook"),
		ParamsFull:         types.StringUnknown(),
		ExternalParamsFull: types.StringUnknown(),
		CellsFull:          types.StringUnknown(),
	}
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	// when
	err := postProcessJsonFields(requestContext, tfModel)

	// then
	require.NoError(t, err)
	assert.True(t, tfModel.ParamsFull.IsUnknown())
	assert.True(t, tfModel.ExternalParamsFull.IsUnknown())
	assert.True(t, tfModel.CellsFull.IsUnknown())
}

// TestPostProcessJsonFields_WithValidJSON tests JSON field processing with valid JSON
func TestPostProcessJsonFields_WithValidJSON(t *testing.T) {
	t.Parallel()
	// given
	cellsJSON := `[{"name":"cell1","type":"code","op":"print('hello')"}]`
	paramsJSON := `[{"name":"param1","type":"string","value":"value1"}]`
	externalParamsJSON := `[{"name":"ext1","source":"alertmanager","json_path":"$.<path>"}]`

	tfModel := &runbooktf.RunbookTFModel{
		Name:               types.StringValue("test_runbook"),
		ParamsFull:         types.StringValue(paramsJSON),
		ExternalParamsFull: types.StringValue(externalParamsJSON),
		CellsFull:          types.StringValue(cellsJSON),
	}
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	// when
	err := postProcessJsonFields(requestContext, tfModel)

	// then
	require.NoError(t, err)
	// Values should be processed and remain valid JSON
	assert.Equal(t, tfModel.ParamsFull, types.StringValue("[{\"export\":false,\"name\":\"param1\",\"param_type\":\"PARAM\",\"required\":false,\"value\":\"value1\"}]"))
	assert.Equal(t, tfModel.ExternalParamsFull, types.StringValue("[{\"export\":false,\"json_path\":\"$.\\u003cpath\\u003e\",\"name\":\"ext1\",\"param_type\":\"EXTERNAL\",\"source\":\"alertmanager\",\"value\":\"\"}]"))
	assert.Equal(t, tfModel.CellsFull, types.StringValue("[{\"enabled\":true,\"name\":\"cell1\",\"op\":\"print('hello')\",\"secret_aware\":false}]"))

	// Verify the JSON is still valid
	var params []interface{}
	err = json.Unmarshal([]byte(tfModel.ParamsFull.ValueString()), &params)
	assert.NoError(t, err)
}

// TestPostProcessJsonFields_WithInvalidJSON tests JSON field processing with invalid JSON
func TestPostProcessJsonFields_WithInvalidJSON(t *testing.T) {
	t.Parallel()
	// given
	tfModel := &runbooktf.RunbookTFModel{
		Name:               types.StringValue("test_runbook"),
		ParamsFull:         types.StringValue("invalid json"),
		ExternalParamsFull: types.StringValue("[valid json]"),
		CellsFull:          types.StringValue("[{\"name\":\"cell1\"}]"),
	}
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	// when
	err := postProcessJsonFields(requestContext, tfModel)

	// then
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid character 'i' looking for beginning of value")
}

// TestSetParamsGroups tests the setParamsGroups helper directly.
func TestSetParamsGroups(t *testing.T) {
	t.Parallel()

	knownParamsGroups, _ := types.ObjectValue(
		converters.ParamsGroupsAttrTypes,
		map[string]attr.Value{
			"required": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("p1")}),
			"optional": types.ListValueMust(types.StringType, []attr.Value{}),
			"exported": types.ListValueMust(types.StringType, []attr.Value{}),
			"external": types.ListValueMust(types.StringType, []attr.Value{}),
		},
	)

	tests := []struct {
		name           string
		source         types.Object
		expectNull     bool
		expectSameAsIn bool
	}{
		{
			name:       "null source is copied as-is (null in plan = intentional)",
			source:     types.ObjectNull(converters.ParamsGroupsAttrTypes),
			expectNull: true,
		},
		{
			name:       "unknown source is normalized to null",
			source:     types.ObjectUnknown(converters.ParamsGroupsAttrTypes),
			expectNull: true,
		},
		{
			name:           "known source is copied as-is",
			source:         knownParamsGroups,
			expectNull:     false,
			expectSameAsIn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			apiModel := &runbooktf.RunbookTFModel{
				ParamsGroups: knownParamsGroups,
			}

			setParamsGroups(tt.source, apiModel)

			if tt.expectNull {
				assert.True(t, apiModel.ParamsGroups.IsNull())
			} else if tt.expectSameAsIn {
				assert.Equal(t, tt.source, apiModel.ParamsGroups)
			}
		})
	}
}

// TestRestoreParamsGroupsFromPlan tests that params_groups is correctly restored from the plan,
// overriding whatever the API returned.
func TestRestoreParamsGroupsFromPlan(t *testing.T) {
	t.Parallel()

	apiComputedParamsGroups, _ := types.ObjectValue(
		converters.ParamsGroupsAttrTypes,
		map[string]attr.Value{
			"required": types.ListValueMust(types.StringType, []attr.Value{}),
			"optional": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("p1")}),
			"exported": types.ListValueMust(types.StringType, []attr.Value{}),
			"external": types.ListValueMust(types.StringType, []attr.Value{}),
		},
	)

	tests := []struct {
		name                  string
		planParamsGroups      types.Object
		expectNullInModel     bool
		expectAPIValueInModel bool
	}{
		{
			name:              "null in plan (user did not configure) → null in state, API value discarded",
			planParamsGroups:  types.ObjectNull(converters.ParamsGroupsAttrTypes),
			expectNullInModel: true,
		},
		{
			name:              "unknown in plan → normalized to null, API value discarded",
			planParamsGroups:  types.ObjectUnknown(converters.ParamsGroupsAttrTypes),
			expectNullInModel: true,
		},
		{
			name: "configured in plan → plan value wins over API value",
			planParamsGroups: func() types.Object {
				v, _ := types.ObjectValue(
					converters.ParamsGroupsAttrTypes,
					map[string]attr.Value{
						"required": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("p1")}),
						"optional": types.ListValueMust(types.StringType, []attr.Value{}),
						"exported": types.ListValueMust(types.StringType, []attr.Value{}),
						"external": types.ListValueMust(types.StringType, []attr.Value{}),
					},
				)
				return v
			}(),
			expectNullInModel: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			requestContext := common.NewRequestContext(ctx).WithOperation(common.Create)

			planGetter := &runbookModelGetter{
				model: runbooktf.RunbookTFModel{ParamsGroups: tt.planParamsGroups},
			}
			apiModel := &runbooktf.RunbookTFModel{
				ParamsGroups: apiComputedParamsGroups,
			}

			err := restoreFieldsFromPlan(requestContext, planGetter, apiModel)

			require.NoError(t, err)
			if tt.expectNullInModel {
				assert.True(t, apiModel.ParamsGroups.IsNull(), "expected params_groups to be null")
			} else {
				assert.Equal(t, tt.planParamsGroups, apiModel.ParamsGroups)
			}
		})
	}
}

// TestPostProcessJsonFullField_GenericType tests the generic postProcessJsonFullField function
func TestPostProcessJsonFullField_GenericType(t *testing.T) {
	t.Parallel()
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}

	tests := []struct {
		name        string
		input       types.String
		expectError bool
		expectNull  bool
	}{
		{
			name:        "Null value",
			input:       types.StringNull(),
			expectError: false,
			expectNull:  false,
		},
		{
			name:        "Unknown value",
			input:       types.StringUnknown(),
			expectError: false,
			expectNull:  false,
		},
		{
			name:        "Valid JSON array",
			input:       types.StringValue(`[{"name":"test","type":"string"}]`),
			expectError: false,
			expectNull:  false,
		},
		{
			name:        "Invalid JSON",
			input:       types.StringValue("not json"),
			expectError: true,
			expectNull:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := postProcessJsonFullField[*MockJsonConfigurable](&tt.input, backendVersion)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectNull {
					assert.True(t, result.IsNull())
				} else if tt.input.IsNull() {
					assert.True(t, result.IsNull())
				} else if tt.input.IsUnknown() {
					assert.True(t, result.IsUnknown())
				}
			}
		})
	}
}

func TestRestoreBaseFieldsFromState_CellsFullPreservedWhenCellsListActive(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	requestContext := common.NewRequestContext(ctx).WithOperation(common.Read)

	cellObj, _ := types.ObjectValue(
		converters.CellsListAttrTypes,
		map[string]attr.Value{
			"op": types.StringValue("host"), "md": types.StringNull(),
			"name": types.StringValue("unnamed"), "enabled": types.BoolValue(true),
			"secret_aware": types.BoolValue(false), "description": types.StringValue(""),
		},
	)
	cellsList, _ := types.ListValue(converters.CellsListObjectType, []attr.Value{cellObj})

	stateGetter := &runbookModelGetter{
		model: runbooktf.RunbookTFModel{
			Cells:          types.StringNull(),
			CellsFull:      types.StringNull(),
			CellsList:      cellsList,
			Params:         types.StringValue("[]"),
			ExternalParams: types.StringValue("[]"),
		},
	}

	apiModel := &runbooktf.RunbookTFModel{
		Cells:     types.StringValue(`[{"op":"host","name":"unnamed","enabled":true}]`),
		CellsFull: types.StringValue(`[{"op":"host","name":"unnamed","enabled":true}]`),
		CellsList: cellsList,
	}

	err := restoreBaseFieldsFromState(requestContext, stateGetter, apiModel)

	require.NoError(t, err)
	assert.True(t, apiModel.Cells.IsNull(), "cells should be null when cells_list is active")
	assert.True(t, apiModel.CellsFull.IsNull(), "cells_full should be null when cells_list is active")
	assert.Equal(t, cellsList, apiModel.CellsList, "cells_list should remain unchanged")
}

func TestRestoreBaseFieldsFromState_CellsFullNotRestoredWhenCellsPath(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	requestContext := common.NewRequestContext(ctx).WithOperation(common.Read)

	stateGetter := &runbookModelGetter{
		model: runbooktf.RunbookTFModel{
			Cells:          types.StringValue(`[{"op":"host"}]`),
			CellsFull:      types.StringValue(`[{"op":"host","enabled":true,"name":"unnamed"}]`),
			Params:         types.StringValue("[]"),
			ExternalParams: types.StringValue("[]"),
		},
	}

	apiCellsFull := `[{"op":"host","enabled":true,"name":"unnamed","secret_aware":false}]`
	apiModel := &runbooktf.RunbookTFModel{
		Cells:     types.StringValue(`[{"op":"host","enabled":true,"name":"unnamed","secret_aware":false}]`),
		CellsFull: types.StringValue(apiCellsFull),
	}

	err := restoreBaseFieldsFromState(requestContext, stateGetter, apiModel)

	require.NoError(t, err)
	assert.Equal(t, `[{"op":"host"}]`, apiModel.Cells.ValueString(), "cells should be restored from state")
	assert.Equal(t, apiCellsFull, apiModel.CellsFull.ValueString(), "cells_full should NOT be restored when cells_list is not active (used for drift)")
}
