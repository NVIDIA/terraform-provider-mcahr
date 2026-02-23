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

package converters

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	runbookapi "terraform/terraform-provider/provider/external_api/resources/runbooks"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParamsGroupsFromTFModel(t *testing.T) {
	tests := []struct {
		name        string
		tfModel     types.Object
		expected    runbookapi.ParamsGroups
		description string
	}{
		{
			name: "Complete params groups with all fields populated",
			tfModel: createTFParamsGroups(t, map[string][]string{
				"required": {"param1", "param2", "param3"},
				"optional": {"param4", "param5"},
				"exported": {"param6"},
				"external": {"param7", "param8"},
			}),
			expected: runbookapi.ParamsGroups{
				Required: []string{"param1", "param2", "param3"},
				Optional: []string{"param4", "param5"},
				Exported: []string{"param6"},
				External: []string{"param7", "param8"},
			},
			description: "Should correctly convert all populated lists",
		},
		{
			name: "Params groups with empty lists",
			tfModel: createTFParamsGroups(t, map[string][]string{
				"required": {},
				"optional": {},
				"exported": {},
				"external": {},
			}),
			expected: runbookapi.ParamsGroups{
				Required: []string{},
				Optional: []string{},
				Exported: []string{},
				External: []string{},
			},
			description: "Should handle empty lists correctly",
		},
		{
			name: "Params groups with only required fields",
			tfModel: createTFParamsGroups(t, map[string][]string{
				"required": {"required_param1", "required_param2"},
				"optional": {},
				"exported": {},
				"external": {},
			}),
			expected: runbookapi.ParamsGroups{
				Required: []string{"required_param1", "required_param2"},
				Optional: []string{},
				Exported: []string{},
				External: []string{},
			},
			description: "Should handle only required params populated",
		},
		{
			name: "Params groups with mixed fields",
			tfModel: createTFParamsGroups(t, map[string][]string{
				"required": {"req1"},
				"optional": {},
				"exported": {"exp1", "exp2", "exp3"},
				"external": {"ext1"},
			}),
			expected: runbookapi.ParamsGroups{
				Required: []string{"req1"},
				Optional: []string{},
				Exported: []string{"exp1", "exp2", "exp3"},
				External: []string{"ext1"},
			},
			description: "Should handle mixed populated and empty lists",
		},
		{
			name: "Params groups with single items",
			tfModel: createTFParamsGroups(t, map[string][]string{
				"required": {"single_req"},
				"optional": {"single_opt"},
				"exported": {"single_exp"},
				"external": {"single_ext"},
			}),
			expected: runbookapi.ParamsGroups{
				Required: []string{"single_req"},
				Optional: []string{"single_opt"},
				Exported: []string{"single_exp"},
				External: []string{"single_ext"},
			},
			description: "Should handle single items in each list",
		},
		{
			name:    "Params groups with null lists",
			tfModel: createTFParamsGroupsWithNulls(t),
			expected: runbookapi.ParamsGroups{
				Required: nil,
				Optional: nil,
				Exported: nil,
				External: nil,
			},
			description: "Should handle null lists and convert to nil slices",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			requestContext := common.NewRequestContext(context.Background())

			// when
			result, err := ParamsGroupsFromTFModel(requestContext, tt.tfModel)

			// then
			require.NoError(t, err, "ParamsGroupsFromTFModel should not error")
			assert.Equal(t, tt.expected.Required, result.Required, "Required params should match")
			assert.Equal(t, tt.expected.Optional, result.Optional, "Optional params should match")
			assert.Equal(t, tt.expected.Exported, result.Exported, "Exported params should match")
			assert.Equal(t, tt.expected.External, result.External, "External params should match")
		})
	}
}

func TestParamsGroupsToTFModel(t *testing.T) {
	tests := []struct {
		name        string
		apiModel    runbookapi.ParamsGroups
		validate    func(t *testing.T, result types.Object)
		expectError bool
		description string
	}{
		{
			name: "Complete params groups with all fields",
			apiModel: runbookapi.ParamsGroups{
				Required: []string{"param1", "param2"},
				Optional: []string{"param3"},
				Exported: []string{"param4", "param5", "param6"},
				External: []string{"param7"},
			},
			expectError: false,
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull(), "Result should not be null")
				assert.False(t, result.IsUnknown(), "Result should not be unknown")

				attrs := result.Attributes()
				require.NotNil(t, attrs, "Attributes should not be nil")

				// Validate required
				reqList, ok := attrs["required"].(types.List)
				require.True(t, ok, "Required should be a List")
				require.Equal(t, 2, len(reqList.Elements()), "Required should have 2 elements")
				reqElems := reqList.Elements()
				assert.Equal(t, "param1", reqElems[0].(types.String).ValueString())
				assert.Equal(t, "param2", reqElems[1].(types.String).ValueString())

				// Validate optional
				optList, ok := attrs["optional"].(types.List)
				require.True(t, ok, "Optional should be a List")
				require.Equal(t, 1, len(optList.Elements()), "Optional should have 1 element")
				optElems := optList.Elements()
				assert.Equal(t, "param3", optElems[0].(types.String).ValueString())

				// Validate exported
				expList, ok := attrs["exported"].(types.List)
				require.True(t, ok, "Exported should be a List")
				require.Equal(t, 3, len(expList.Elements()), "Exported should have 3 elements")
				expElems := expList.Elements()
				assert.Equal(t, "param4", expElems[0].(types.String).ValueString())
				assert.Equal(t, "param5", expElems[1].(types.String).ValueString())
				assert.Equal(t, "param6", expElems[2].(types.String).ValueString())

				// Validate external
				extList, ok := attrs["external"].(types.List)
				require.True(t, ok, "External should be a List")
				require.Equal(t, 1, len(extList.Elements()), "External should have 1 element")
				extElems := extList.Elements()
				assert.Equal(t, "param7", extElems[0].(types.String).ValueString())
			},
			description: "Should correctly convert all populated lists to TF model",
		},
		{
			name: "Empty params groups",
			apiModel: runbookapi.ParamsGroups{
				Required: []string{},
				Optional: []string{},
				Exported: []string{},
				External: []string{},
			},
			expectError: false,
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull(), "Result should not be null")
				attrs := result.Attributes()

				// All lists should be empty but not null
				reqList, ok := attrs["required"].(types.List)
				require.True(t, ok, "Required should be a List")
				assert.Equal(t, 0, len(reqList.Elements()), "Required should be empty")

				optList, ok := attrs["optional"].(types.List)
				require.True(t, ok, "Optional should be a List")
				assert.Equal(t, 0, len(optList.Elements()), "Optional should be empty")

				expList, ok := attrs["exported"].(types.List)
				require.True(t, ok, "Exported should be a List")
				assert.Equal(t, 0, len(expList.Elements()), "Exported should be empty")

				extList, ok := attrs["external"].(types.List)
				require.True(t, ok, "External should be a List")
				assert.Equal(t, 0, len(extList.Elements()), "External should be empty")
			},
			description: "Should handle empty lists correctly",
		},
		{
			name: "Only required params",
			apiModel: runbookapi.ParamsGroups{
				Required: []string{"only_required"},
				Optional: []string{},
				Exported: []string{},
				External: []string{},
			},
			expectError: false,
			validate: func(t *testing.T, result types.Object) {
				attrs := result.Attributes()

				// Validate required
				reqList, ok := attrs["required"].(types.List)
				require.True(t, ok, "Required should be a List")
				require.Equal(t, 1, len(reqList.Elements()), "Required should have 1 element")
				reqElems := reqList.Elements()
				assert.Equal(t, "only_required", reqElems[0].(types.String).ValueString())

				// Validate optional
				optList, ok := attrs["optional"].(types.List)
				require.True(t, ok, "Optional should be a List")
				assert.Equal(t, 0, len(optList.Elements()), "Optional should be empty")

				// Validate exported
				expList, ok := attrs["exported"].(types.List)
				require.True(t, ok, "Exported should be a List")
				assert.Equal(t, 0, len(expList.Elements()), "Exported should be empty")

				// Validate external
				extList, ok := attrs["external"].(types.List)
				require.True(t, ok, "External should be a List")
				assert.Equal(t, 0, len(extList.Elements()), "External should be empty")
			},
			description: "Should handle only required params populated",
		},
		{
			name: "Nil slices",
			apiModel: runbookapi.ParamsGroups{
				Required: nil,
				Optional: nil,
				Exported: nil,
				External: nil,
			},
			expectError: false,
			validate: func(t *testing.T, result types.Object) {
				assert.False(t, result.IsNull(), "Result should not be null")
				// types.ObjectValueFrom should handle nil slices gracefully
			},
			description: "Should handle nil slices without error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			requestContext := common.NewRequestContext(context.Background())

			// when
			result, diags := ParamsGroupsToTFModel(requestContext, tt.apiModel)

			// then
			if tt.expectError {
				assert.True(t, diags.HasError(), "Should have errors")
			} else {
				assert.False(t, diags.HasError(), "Should not have errors: %v", diags)
				require.NotNil(t, result, "Result should not be nil")
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestParamsGroupsRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		original runbookapi.ParamsGroups
	}{
		{
			name: "Complete round trip with all fields",
			original: runbookapi.ParamsGroups{
				Required: []string{"req1", "req2"},
				Optional: []string{"opt1"},
				Exported: []string{"exp1", "exp2", "exp3"},
				External: []string{"ext1"},
			},
		},
		{
			name: "Round trip with empty lists",
			original: runbookapi.ParamsGroups{
				Required: []string{},
				Optional: []string{},
				Exported: []string{},
				External: []string{},
			},
		},
		{
			name: "Round trip with single values",
			original: runbookapi.ParamsGroups{
				Required: []string{"single"},
				Optional: []string{},
				Exported: []string{},
				External: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			requestContext := common.NewRequestContext(context.Background())

			// when - convert API -> TF
			tfModel, diags := ParamsGroupsToTFModel(requestContext, tt.original)
			require.False(t, diags.HasError(), "API to TF conversion should not error")

			// then - convert TF -> API
			result, err := ParamsGroupsFromTFModel(requestContext, tfModel)
			require.NoError(t, err, "ParamsGroupsFromTFModel should not error")

			// verify the round trip preserves data
			assert.Equal(t, tt.original.Required, result.Required, "Required should match after round trip")
			assert.Equal(t, tt.original.Optional, result.Optional, "Optional should match after round trip")
			assert.Equal(t, tt.original.Exported, result.Exported, "Exported should match after round trip")
			assert.Equal(t, tt.original.External, result.External, "External should match after round trip")
		})
	}
}

func TestParamsGroupsAttrTypes(t *testing.T) {
	// Verify the attribute types are correctly defined
	t.Run("Verify attribute types structure", func(t *testing.T) {
		require.NotNil(t, ParamsGroupsAttrTypes, "ParamsGroupsAttrTypes should be defined")
		assert.Equal(t, 4, len(ParamsGroupsAttrTypes), "Should have 4 attributes")

		// Check each attribute type
		assert.Contains(t, ParamsGroupsAttrTypes, "required", "Should have required attribute")
		assert.Contains(t, ParamsGroupsAttrTypes, "optional", "Should have optional attribute")
		assert.Contains(t, ParamsGroupsAttrTypes, "exported", "Should have exported attribute")
		assert.Contains(t, ParamsGroupsAttrTypes, "external", "Should have external attribute")

		// Verify all are ListType with StringType elements
		for key, attrType := range ParamsGroupsAttrTypes {
			listType, ok := attrType.(types.ListType)
			assert.True(t, ok, "Attribute %s should be a ListType", key)
			assert.Equal(t, types.StringType, listType.ElemType, "Attribute %s should have StringType elements", key)
		}
	})
}

// Helper functions

func createTFParamsGroups(t *testing.T, params map[string][]string) types.Object {
	ctx := context.Background()

	required, diags := types.ListValueFrom(ctx, types.StringType, params["required"])
	require.False(t, diags.HasError())

	optional, diags := types.ListValueFrom(ctx, types.StringType, params["optional"])
	require.False(t, diags.HasError())

	exported, diags := types.ListValueFrom(ctx, types.StringType, params["exported"])
	require.False(t, diags.HasError())

	external, diags := types.ListValueFrom(ctx, types.StringType, params["external"])
	require.False(t, diags.HasError())

	obj, diags := types.ObjectValue(
		ParamsGroupsAttrTypes,
		map[string]attr.Value{
			"required": required,
			"optional": optional,
			"exported": exported,
			"external": external,
		},
	)
	require.False(t, diags.HasError())

	return obj
}

func createTFParamsGroupsWithNulls(t *testing.T) types.Object {
	obj, diags := types.ObjectValue(
		ParamsGroupsAttrTypes,
		map[string]attr.Value{
			"required": types.ListNull(types.StringType),
			"optional": types.ListNull(types.StringType),
			"exported": types.ListNull(types.StringType),
			"external": types.ListNull(types.StringType),
		},
	)
	require.False(t, diags.HasError(), "Should create object with null lists: %v", diags)

	return obj
}
