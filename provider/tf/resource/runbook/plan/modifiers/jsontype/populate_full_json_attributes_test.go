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

package jsonmodifier

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPopulateFullJsonAttributes(t *testing.T) {
	tests := []struct {
		name           string
		resultValues   *model.RunbookTFModel
		plan           *model.RunbookTFModel
		state          *model.RunbookTFModel
		backendVersion *version.BackendVersion
		expectError    bool
		validate       func(t *testing.T, result *model.RunbookTFModel)
	}{
		{
			name: "Process cells field",
			resultValues: &model.RunbookTFModel{
				Cells: types.StringValue(`[{"op": "print('test')", "name": "cell1"}]`),
			},
			plan: &model.RunbookTFModel{
				CellsFull: types.StringValue("not_null"),
			},
			state:          &model.RunbookTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.False(t, result.CellsFull.IsNull())
				assert.Contains(t, result.CellsFull.ValueString(), "cell1")
			},
		},
		{
			name: "Process params field",
			resultValues: &model.RunbookTFModel{
				Params: types.StringValue(`[{"name": "param1", "value": "value1"}]`),
			},
			plan: &model.RunbookTFModel{
				ParamsFull: types.StringValue("not_null"),
			},
			state:          &model.RunbookTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.False(t, result.ParamsFull.IsNull())
				assert.Contains(t, result.ParamsFull.ValueString(), "param1")
			},
		},
		{
			name: "Process external_params field",
			resultValues: &model.RunbookTFModel{
				ExternalParams: types.StringValue(`[{"name": "ext1", "source": "api"}]`),
			},
			plan: &model.RunbookTFModel{
				ExternalParamsFull: types.StringValue("not_null"),
			},
			state:          &model.RunbookTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.False(t, result.ExternalParamsFull.IsNull())
				assert.Contains(t, result.ExternalParamsFull.ValueString(), "ext1")
			},
		},
		{
			name: "Skip delete operation - null plan",
			resultValues: &model.RunbookTFModel{
				Cells: types.StringValue(`[{"op": "test"}]`),
			},
			plan: &model.RunbookTFModel{
				CellsFull: types.StringNull(), // Indicates delete
			},
			state:          &model.RunbookTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.True(t, result.CellsFull.IsNull())
			},
		},
		{
			name: "Handle invalid JSON in cells",
			resultValues: &model.RunbookTFModel{
				Cells: types.StringValue(`{invalid json}`),
			},
			plan: &model.RunbookTFModel{
				CellsFull: types.StringValue("not_null"),
			},
			state:          &model.RunbookTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    true,
		},
		{
			name: "Process all fields",
			resultValues: &model.RunbookTFModel{
				Cells:          types.StringValue(`[{"op": "code", "name": "c1"}]`),
				Params:         types.StringValue(`[{"name": "p1", "value": "v1"}]`),
				ExternalParams: types.StringValue(`[{"name": "ep1", "source": "src"}]`),
			},
			plan: &model.RunbookTFModel{
				CellsFull:          types.StringValue("not_null"),
				ParamsFull:         types.StringValue("not_null"),
				ExternalParamsFull: types.StringValue("not_null"),
			},
			state:          &model.RunbookTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.False(t, result.CellsFull.IsNull())
				assert.False(t, result.ParamsFull.IsNull())
				assert.False(t, result.ExternalParamsFull.IsNull())
			},
		},
		{
			name: "Version filtering in remarshal",
			resultValues: &model.RunbookTFModel{
				Params: types.StringValue(`[{"name": "p1", "description": "desc"}]`),
			},
			plan: &model.RunbookTFModel{
				ParamsFull: types.StringValue("not_null"),
			},
			state:          &model.RunbookTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-28.3.0", Major: 28, Minor: 3, Patch: 0}, // Before description support
			expectError:    false,
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.False(t, result.ParamsFull.IsNull())
				// Description should be filtered out due to version
				assert.NotContains(t, result.ParamsFull.ValueString(), "description")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()

			// when
			err := PopulateFullJsonAttributes(ctx, tt.resultValues, tt.resultValues, tt.plan, tt.state, tt.backendVersion)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, tt.resultValues)
				}
			}
		})
	}
}

func TestIsDeleteOperation(t *testing.T) {
	tests := []struct {
		name       string
		plan       *model.RunbookTFModel
		attrConfig JsonAttributeConfig
		expected   bool
	}{
		{
			name: "Delete operation - null cells_full",
			plan: &model.RunbookTFModel{
				CellsFull: types.StringNull(),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["cells"],
			expected:   true,
		},
		{
			name: "Not delete operation - non-null cells_full",
			plan: &model.RunbookTFModel{
				CellsFull: types.StringValue("some_value"),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["cells"],
			expected:   false,
		},
		{
			name: "Delete operation - null params_full",
			plan: &model.RunbookTFModel{
				ParamsFull: types.StringNull(),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["params"],
			expected:   true,
		},
		{
			name: "Not delete operation - empty string",
			plan: &model.RunbookTFModel{
				CellsFull: types.StringValue(""),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["cells"],
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := isDeleteOperation(tt.plan, tt.attrConfig)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJsonAttributeConfig_Functions(t *testing.T) {
	tests := []struct {
		name        string
		attribute   string
		setupModel  func() *model.RunbookTFModel
		testGetters bool
		testSetters bool
	}{
		{
			name:      "Cells attribute config",
			attribute: "cells",
			setupModel: func() *model.RunbookTFModel {
				return &model.RunbookTFModel{
					Cells:     types.StringValue("cells_value"),
					CellsFull: types.StringValue("cells_full_value"),
				}
			},
			testGetters: true,
			testSetters: true,
		},
		{
			name:      "Params attribute config",
			attribute: "params",
			setupModel: func() *model.RunbookTFModel {
				return &model.RunbookTFModel{
					Params:     types.StringValue("params_value"),
					ParamsFull: types.StringValue("params_full_value"),
				}
			},
			testGetters: true,
			testSetters: true,
		},
		{
			name:      "External params attribute config",
			attribute: "external_params",
			setupModel: func() *model.RunbookTFModel {
				return &model.RunbookTFModel{
					ExternalParams:     types.StringValue("ext_params_value"),
					ExternalParamsFull: types.StringValue("ext_params_full_value"),
				}
			},
			testGetters: true,
			testSetters: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			attrConfig := JSON_ATTRIBUTES_TO_POPULATE[tt.attribute]
			model := tt.setupModel()

			if tt.testGetters {
				// Test GetAttr and GetFullAttr
				regularValue := attrConfig.GetAttr(model)
				fullValue := attrConfig.GetFullAttr(model)

				switch tt.attribute {
				case "cells":
					assert.Equal(t, "cells_value", regularValue.ValueString())
					assert.Equal(t, "cells_full_value", fullValue.ValueString())
				case "params":
					assert.Equal(t, "params_value", regularValue.ValueString())
					assert.Equal(t, "params_full_value", fullValue.ValueString())
				case "external_params":
					assert.Equal(t, "ext_params_value", regularValue.ValueString())
					assert.Equal(t, "ext_params_full_value", fullValue.ValueString())
				}
			}

			if tt.testSetters {
				// Test SetFullAttr
				newFullValue := types.StringValue("new_full_value")
				attrConfig.SetFullAttr(model, newFullValue)

				fullValue := attrConfig.GetFullAttr(model)
				assert.Equal(t, "new_full_value", fullValue.ValueString())
			}
		})
	}
}

func TestJSON_ATTRIBUTES_TO_POPULATE_Configuration(t *testing.T) {
	// Verify all expected attributes are configured
	expectedAttributes := []string{"cells", "params", "external_params"}

	for _, attr := range expectedAttributes {
		t.Run("Has_"+attr, func(t *testing.T) {
			attrConfig, exists := JSON_ATTRIBUTES_TO_POPULATE[attr]
			assert.True(t, exists, "Expected attribute %s to be in JSON_ATTRIBUTES_TO_POPULATE", attr)

			// Verify full attribute name
			switch attr {
			case "cells":
				assert.Equal(t, "cells_full", attrConfig.FullAttrName)
			case "params":
				assert.Equal(t, "params_full", attrConfig.FullAttrName)
			case "external_params":
				assert.Equal(t, "external_params_full", attrConfig.FullAttrName)
			}

			// Verify functions are set
			assert.NotNil(t, attrConfig.RemarshalFunc)
			assert.NotNil(t, attrConfig.GetAttr)
			assert.NotNil(t, attrConfig.GetFullAttr)
			assert.NotNil(t, attrConfig.SetFullAttr)
		})
	}

	// Verify no extra attributes
	assert.Equal(t, len(expectedAttributes), len(JSON_ATTRIBUTES_TO_POPULATE))
}

func TestPopulateFullJsonAttributes_Integration(t *testing.T) {
	// Integration test with real JSON processing
	ctx := context.Background()
	backendVersion := &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1}

	resultValues := &model.RunbookTFModel{
		Cells: types.StringValue(`[
			{"op": "print('hello')", "name": "cell1", "enabled": true},
			{"md": "# Header", "name": "cell2", "enabled": false}
		]`),
		Params: types.StringValue(`[
			{"name": "param1", "value": "value1", "required": true},
			{"name": "param2", "value": "value2", "required": false}
		]`),
		ExternalParams: types.StringValue(`[
			{"name": "ext1", "source": "api", "json_path": "$.data"},
			{"name": "ext2", "source": "db", "json_path": "$.result"}
		]`),
	}

	plan := &model.RunbookTFModel{
		CellsFull:          types.StringValue("not_null"),
		ParamsFull:         types.StringValue("not_null"),
		ExternalParamsFull: types.StringValue("not_null"),
	}

	state := &model.RunbookTFModel{}

	// when
	err := PopulateFullJsonAttributes(ctx, resultValues, resultValues, plan, state, backendVersion)

	// then
	require.NoError(t, err)

	// Verify cells_full was populated
	assert.False(t, resultValues.CellsFull.IsNull())
	cellsFullStr := resultValues.CellsFull.ValueString()
	assert.Contains(t, cellsFullStr, "cell1")
	assert.Contains(t, cellsFullStr, "cell2")

	// Verify params_full was populated
	assert.False(t, resultValues.ParamsFull.IsNull())
	paramsFullStr := resultValues.ParamsFull.ValueString()
	assert.Contains(t, paramsFullStr, "param1")
	assert.Contains(t, paramsFullStr, "param2")

	// Verify external_params_full was populated
	assert.False(t, resultValues.ExternalParamsFull.IsNull())
	extParamsFullStr := resultValues.ExternalParamsFull.ValueString()
	assert.Contains(t, extParamsFullStr, "ext1")
	assert.Contains(t, extParamsFullStr, "ext2")
}

func TestShouldSkipForReplacement(t *testing.T) {
	cellsConfig := JSON_ATTRIBUTES_TO_POPULATE["cells"]
	paramsConfig := JSON_ATTRIBUTES_TO_POPULATE["params"]
	extParamsConfig := JSON_ATTRIBUTES_TO_POPULATE["external_params"]

	dummyList := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("item")})

	t.Run("cells_list active and cells not explicitly set → skip and null", func(t *testing.T) {
		resultValues := &model.RunbookTFModel{
			CellsList: dummyList,
			Cells:     types.StringValue("[]"),
			CellsFull: types.StringValue("[]"),
		}
		resultWithoutDefaults := &model.RunbookTFModel{
			Cells: types.StringNull(),
		}
		plan := &model.RunbookTFModel{
			CellsFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(cellsConfig, resultValues, resultWithoutDefaults, plan)

		assert.True(t, skipped)
		assert.True(t, resultValues.Cells.IsNull(), "cells should be nulled")
		assert.True(t, resultValues.CellsFull.IsNull(), "cells_full should be nulled")
		assert.True(t, plan.CellsFull.IsNull(), "plan cells_full should be nulled")
		assert.Equal(t, dummyList, resultValues.CellsList, "cells_list should remain unchanged")
	})

	t.Run("cells_list active but cells explicitly set → don't skip (conflict case)", func(t *testing.T) {
		resultValues := &model.RunbookTFModel{
			CellsList: dummyList,
			Cells:     types.StringValue(`[{"op":"test"}]`),
			CellsFull: types.StringValue("something"),
		}
		resultWithoutDefaults := &model.RunbookTFModel{
			Cells: types.StringValue(`[{"op":"test"}]`),
		}
		plan := &model.RunbookTFModel{
			CellsFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(cellsConfig, resultValues, resultWithoutDefaults, plan)

		assert.False(t, skipped)
		assert.Equal(t, `[{"op":"test"}]`, resultValues.Cells.ValueString())
	})

	t.Run("cells_list not active → don't skip", func(t *testing.T) {
		resultValues := &model.RunbookTFModel{
			CellsList: types.ListNull(types.StringType),
			Cells:     types.StringValue(`[{"op":"test"}]`),
		}
		resultWithoutDefaults := &model.RunbookTFModel{
			Cells: types.StringValue(`[{"op":"test"}]`),
		}
		plan := &model.RunbookTFModel{}

		skipped := shouldSkipForReplacement(cellsConfig, resultValues, resultWithoutDefaults, plan)

		assert.False(t, skipped)
	})

	t.Run("params_list active and params not explicitly set → skip and null", func(t *testing.T) {
		resultValues := &model.RunbookTFModel{
			ParamsList: dummyList,
			Params:     types.StringValue("[]"),
			ParamsFull: types.StringValue("[]"),
		}
		resultWithoutDefaults := &model.RunbookTFModel{
			Params: types.StringNull(),
		}
		plan := &model.RunbookTFModel{
			ParamsFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(paramsConfig, resultValues, resultWithoutDefaults, plan)

		assert.True(t, skipped)
		assert.True(t, resultValues.Params.IsNull(), "params should be nulled")
		assert.True(t, resultValues.ParamsFull.IsNull(), "params_full should be nulled")
		assert.True(t, plan.ParamsFull.IsNull(), "plan params_full should be nulled")
		assert.Equal(t, dummyList, resultValues.ParamsList, "params_list should remain unchanged")
	})

	t.Run("params_list active but params explicitly set → don't skip (conflict case)", func(t *testing.T) {
		resultValues := &model.RunbookTFModel{
			ParamsList: dummyList,
			Params:     types.StringValue(`[{"name":"p1","value":"v1"}]`),
			ParamsFull: types.StringValue("something"),
		}
		resultWithoutDefaults := &model.RunbookTFModel{
			Params: types.StringValue(`[{"name":"p1","value":"v1"}]`),
		}
		plan := &model.RunbookTFModel{
			ParamsFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(paramsConfig, resultValues, resultWithoutDefaults, plan)

		assert.False(t, skipped)
		assert.Equal(t, `[{"name":"p1","value":"v1"}]`, resultValues.Params.ValueString())
	})

	t.Run("external_params_list active and external_params not explicitly set → skip and null", func(t *testing.T) {
		resultValues := &model.RunbookTFModel{
			ExternalParamsList: dummyList,
			ExternalParams:     types.StringValue("[]"),
			ExternalParamsFull: types.StringValue("[]"),
		}
		resultWithoutDefaults := &model.RunbookTFModel{
			ExternalParams: types.StringNull(),
		}
		plan := &model.RunbookTFModel{
			ExternalParamsFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(extParamsConfig, resultValues, resultWithoutDefaults, plan)

		assert.True(t, skipped)
		assert.True(t, resultValues.ExternalParams.IsNull(), "external_params should be nulled")
		assert.True(t, resultValues.ExternalParamsFull.IsNull(), "external_params_full should be nulled")
		assert.True(t, plan.ExternalParamsFull.IsNull(), "plan external_params_full should be nulled")
		assert.Equal(t, dummyList, resultValues.ExternalParamsList, "external_params_list should remain unchanged")
	})

	t.Run("external_params_list active but external_params explicitly set → don't skip (conflict case)", func(t *testing.T) {
		resultValues := &model.RunbookTFModel{
			ExternalParamsList: dummyList,
			ExternalParams:     types.StringValue(`[{"name":"ep1","source":"api"}]`),
			ExternalParamsFull: types.StringValue("something"),
		}
		resultWithoutDefaults := &model.RunbookTFModel{
			ExternalParams: types.StringValue(`[{"name":"ep1","source":"api"}]`),
		}
		plan := &model.RunbookTFModel{
			ExternalParamsFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(extParamsConfig, resultValues, resultWithoutDefaults, plan)

		assert.False(t, skipped)
		assert.Equal(t, `[{"name":"ep1","source":"api"}]`, resultValues.ExternalParams.ValueString())
	})

	t.Run("replacement field at zero value → don't skip", func(t *testing.T) {
		resultValues := &model.RunbookTFModel{
			Params: types.StringValue(`[{"name":"p1"}]`),
		}
		resultWithoutDefaults := &model.RunbookTFModel{}
		plan := &model.RunbookTFModel{}

		skipped := shouldSkipForReplacement(paramsConfig, resultValues, resultWithoutDefaults, plan)

		assert.False(t, skipped)
	})
}
