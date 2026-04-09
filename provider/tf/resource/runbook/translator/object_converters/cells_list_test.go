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
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCellsListToInternalCells(t *testing.T) {
	tests := []struct {
		name        string
		tfList      types.List
		expected    []customattribute.CellJson
		expectError bool
		errContains string
	}{
		{
			name:   "Single op cell with all fields",
			tfList: createCellsList(t, []cellInput{{op: strPtr("host | limit 1"), name: "my_cell", enabled: true, secretAware: true, description: "A test cell"}}),
			expected: []customattribute.CellJson{
				{
					Op:          common.NewOptional("host | limit 1"),
					Md:          common.NewOptionalUnset[string](),
					Name:        "my_cell",
					Enabled:     true,
					SecretAware: true,
					Description: "A test cell",
				},
			},
		},
		{
			name:   "Single md cell with defaults",
			tfList: createCellsList(t, []cellInput{{md: strPtr("# Hello"), name: "unnamed", enabled: true, secretAware: false, description: ""}}),
			expected: []customattribute.CellJson{
				{
					Op:          common.NewOptionalUnset[string](),
					Md:          common.NewOptional("# Hello"),
					Name:        "unnamed",
					Enabled:     true,
					SecretAware: false,
					Description: "",
				},
			},
		},
		{
			name: "Multiple cells mixed op and md",
			tfList: createCellsList(t, []cellInput{
				{md: strPtr("CREATE"), name: "unnamed", enabled: true, secretAware: false, description: ""},
				{op: strPtr("action success = `echo SUCCESS`"), name: "unnamed", enabled: true, secretAware: false, description: "Creates action"},
				{op: strPtr("success"), name: "run_cell", enabled: false, secretAware: false, description: "Runs the action"},
			}),
			expected: []customattribute.CellJson{
				{Op: common.NewOptionalUnset[string](), Md: common.NewOptional("CREATE"), Name: "unnamed", Enabled: true, SecretAware: false, Description: ""},
				{Op: common.NewOptional("action success = `echo SUCCESS`"), Md: common.NewOptionalUnset[string](), Name: "unnamed", Enabled: true, SecretAware: false, Description: "Creates action"},
				{Op: common.NewOptional("success"), Md: common.NewOptionalUnset[string](), Name: "run_cell", Enabled: false, SecretAware: false, Description: "Runs the action"},
			},
		},
		{
			name:     "Empty list",
			tfList:   createCellsList(t, []cellInput{}),
			expected: []customattribute.CellJson{},
		},
		{
			name:     "Null list returns nil",
			tfList:   types.ListNull(CellsListObjectType),
			expected: nil,
		},
		{
			name:        "Cell with both op and md errors",
			tfList:      createCellsListRaw(t, []map[string]attr.Value{{"op": types.StringValue("cmd"), "md": types.StringValue("text"), "name": types.StringValue("unnamed"), "enabled": types.BoolValue(true), "secret_aware": types.BoolValue(false), "description": types.StringValue("")}}),
			expectError: true,
			errContains: "both op and md set",
		},
		{
			name:        "Cell with neither op nor md errors",
			tfList:      createCellsListRaw(t, []map[string]attr.Value{{"op": types.StringNull(), "md": types.StringNull(), "name": types.StringValue("unnamed"), "enabled": types.BoolValue(true), "secret_aware": types.BoolValue(false), "description": types.StringValue("")}}),
			expectError: true,
			errContains: "must have either op or md set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CellsListToInternalCells(context.Background(), tt.tfList)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.Equal(t, len(tt.expected), len(result))
			for i := range tt.expected {
				assert.Equal(t, tt.expected[i].Op.IsSet, result[i].Op.IsSet, "cell %d Op.IsSet", i)
				if tt.expected[i].Op.IsSet {
					assert.Equal(t, tt.expected[i].Op.Get(), result[i].Op.Get(), "cell %d Op value", i)
				}
				assert.Equal(t, tt.expected[i].Md.IsSet, result[i].Md.IsSet, "cell %d Md.IsSet", i)
				if tt.expected[i].Md.IsSet {
					assert.Equal(t, tt.expected[i].Md.Get(), result[i].Md.Get(), "cell %d Md value", i)
				}
				assert.Equal(t, tt.expected[i].Name, result[i].Name, "cell %d Name", i)
				assert.Equal(t, tt.expected[i].Enabled, result[i].Enabled, "cell %d Enabled", i)
				assert.Equal(t, tt.expected[i].SecretAware, result[i].SecretAware, "cell %d SecretAware", i)
				assert.Equal(t, tt.expected[i].Description, result[i].Description, "cell %d Description", i)
			}
		})
	}
}

func TestCellsListFromAPICells(t *testing.T) {
	tests := []struct {
		name        string
		apiCells    []customattribute.CellJsonAPI
		validate    func(t *testing.T, result types.List)
		expectError bool
	}{
		{
			name: "Op cell with all fields",
			apiCells: []customattribute.CellJsonAPI{
				{Type: "OP_LANG", Content: "host | limit 1", Name: "my_cell", Enabled: true, SecretAware: true, Description: "A test cell"},
			},
			validate: func(t *testing.T, result types.List) {
				require.Equal(t, 1, len(result.Elements()))
				obj := result.Elements()[0].(types.Object)
				attrs := obj.Attributes()
				assert.Equal(t, "host | limit 1", attrs["op"].(types.String).ValueString())
				assert.True(t, attrs["md"].(types.String).IsNull())
				assert.Equal(t, "my_cell", attrs["name"].(types.String).ValueString())
				assert.Equal(t, true, attrs["enabled"].(types.Bool).ValueBool())
				assert.Equal(t, true, attrs["secret_aware"].(types.Bool).ValueBool())
				assert.Equal(t, "A test cell", attrs["description"].(types.String).ValueString())
			},
		},
		{
			name: "Md cell",
			apiCells: []customattribute.CellJsonAPI{
				{Type: "MARKDOWN", Content: "# Hello World", Name: "unnamed", Enabled: true, SecretAware: false, Description: ""},
			},
			validate: func(t *testing.T, result types.List) {
				require.Equal(t, 1, len(result.Elements()))
				obj := result.Elements()[0].(types.Object)
				attrs := obj.Attributes()
				assert.True(t, attrs["op"].(types.String).IsNull())
				assert.Equal(t, "# Hello World", attrs["md"].(types.String).ValueString())
			},
		},
		{
			name: "Multiple cells",
			apiCells: []customattribute.CellJsonAPI{
				{Type: "MARKDOWN", Content: "CREATE", Name: "unnamed", Enabled: true},
				{Type: "OP_LANG", Content: "action success = `echo SUCCESS`", Name: "unnamed", Enabled: true, Description: "Creates action"},
				{Type: "OP_LANG", Content: "success", Name: "run_cell", Enabled: false},
			},
			validate: func(t *testing.T, result types.List) {
				require.Equal(t, 3, len(result.Elements()))
				obj0 := result.Elements()[0].(types.Object)
				assert.Equal(t, "CREATE", obj0.Attributes()["md"].(types.String).ValueString())
				obj1 := result.Elements()[1].(types.Object)
				assert.Equal(t, "action success = `echo SUCCESS`", obj1.Attributes()["op"].(types.String).ValueString())
				obj2 := result.Elements()[2].(types.Object)
				assert.Equal(t, false, obj2.Attributes()["enabled"].(types.Bool).ValueBool())
			},
		},
		{
			name:     "Empty cells list",
			apiCells: []customattribute.CellJsonAPI{},
			validate: func(t *testing.T, result types.List) {
				assert.Equal(t, 0, len(result.Elements()))
				assert.False(t, result.IsNull())
			},
		},
		{
			name:     "Nil cells returns null list",
			apiCells: nil,
			validate: func(t *testing.T, result types.List) {
				assert.True(t, result.IsNull())
			},
		},
		{
			name: "Cell using cell_type field instead of type",
			apiCells: []customattribute.CellJsonAPI{
				{CellType: "OP_LANG", Content: "host", Name: "unnamed", Enabled: true},
			},
			validate: func(t *testing.T, result types.List) {
				require.Equal(t, 1, len(result.Elements()))
				obj := result.Elements()[0].(types.Object)
				attrs := obj.Attributes()
				assert.Equal(t, "host", attrs["op"].(types.String).ValueString())
				assert.True(t, attrs["md"].(types.String).IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := CellsListFromAPICells(tt.apiCells)

			if tt.expectError {
				assert.True(t, diags.HasError())
				return
			}

			assert.False(t, diags.HasError(), "unexpected diags: %v", diags)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestCellsListRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		original []customattribute.CellJsonAPI
	}{
		{
			name: "Round trip with op and md cells",
			original: []customattribute.CellJsonAPI{
				{Type: "MARKDOWN", Content: "# Setup", Name: "setup_md", Enabled: true, SecretAware: false, Description: "Setup section"},
				{Type: "OP_LANG", Content: "host | limit 5", Name: "op_cell", Enabled: true, SecretAware: true, Description: "Op cell"},
				{Type: "OP_LANG", Content: "success", Name: "disabled_cell", Enabled: false, SecretAware: false, Description: "Disabled"},
			},
		},
		{
			name:     "Round trip with empty list",
			original: []customattribute.CellJsonAPI{},
		},
		{
			name: "Round trip with single cell",
			original: []customattribute.CellJsonAPI{
				{Type: "OP_LANG", Content: "host", Name: "unnamed", Enabled: true, SecretAware: false, Description: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// API -> TF
			tfList, diags := CellsListFromAPICells(tt.original)
			require.False(t, diags.HasError(), "API to TF should not error: %v", diags)

			// TF -> internal
			internalCells, err := CellsListToInternalCells(context.Background(), tfList)
			require.NoError(t, err, "TF to internal should not error")

			require.Equal(t, len(tt.original), len(internalCells))
			for i, orig := range tt.original {
				cell := internalCells[i]
				assert.Equal(t, orig.Name, cell.Name, "cell %d Name", i)
				assert.Equal(t, orig.Enabled, cell.Enabled, "cell %d Enabled", i)
				assert.Equal(t, orig.SecretAware, cell.SecretAware, "cell %d SecretAware", i)
				assert.Equal(t, orig.Description, cell.Description, "cell %d Description", i)

				if orig.Type == "OP_LANG" {
					assert.True(t, cell.Op.IsSet, "cell %d should have Op set", i)
					assert.Equal(t, orig.Content, cell.Op.Get(), "cell %d Op content", i)
				} else {
					assert.True(t, cell.Md.IsSet, "cell %d should have Md set", i)
					assert.Equal(t, orig.Content, cell.Md.Get(), "cell %d Md content", i)
				}
			}
		})
	}
}

func TestCellsListAttrTypesStructure(t *testing.T) {
	require.NotNil(t, CellsListAttrTypes)
	assert.Equal(t, 6, len(CellsListAttrTypes))

	expectedKeys := []string{"op", "md", "name", "enabled", "secret_aware", "description"}
	for _, key := range expectedKeys {
		assert.Contains(t, CellsListAttrTypes, key, "Should have %s attribute", key)
	}

	assert.Equal(t, types.StringType, CellsListAttrTypes["op"])
	assert.Equal(t, types.StringType, CellsListAttrTypes["md"])
	assert.Equal(t, types.StringType, CellsListAttrTypes["name"])
	assert.Equal(t, types.BoolType, CellsListAttrTypes["enabled"])
	assert.Equal(t, types.BoolType, CellsListAttrTypes["secret_aware"])
	assert.Equal(t, types.StringType, CellsListAttrTypes["description"])
}

func TestCellTFModelToInternal(t *testing.T) {
	tests := []struct {
		name        string
		model       CellTFModel
		expectError bool
		errContains string
		validate    func(t *testing.T, cell *customattribute.CellJson)
	}{
		{
			name: "Op cell",
			model: CellTFModel{
				Op: types.StringValue("host | limit 1"), Md: types.StringNull(),
				Name: types.StringValue("my_cell"), Enabled: types.BoolValue(true),
				SecretAware: types.BoolValue(true), Description: types.StringValue("desc"),
			},
			validate: func(t *testing.T, cell *customattribute.CellJson) {
				assert.True(t, cell.Op.IsSet)
				assert.Equal(t, "host | limit 1", cell.Op.Get())
				assert.False(t, cell.Md.IsSet)
				assert.Equal(t, "my_cell", cell.Name)
				assert.Equal(t, true, cell.Enabled)
				assert.Equal(t, true, cell.SecretAware)
				assert.Equal(t, "desc", cell.Description)
			},
		},
		{
			name: "Md cell",
			model: CellTFModel{
				Op: types.StringNull(), Md: types.StringValue("# Hello"),
				Name: types.StringValue("unnamed"), Enabled: types.BoolValue(false),
				SecretAware: types.BoolValue(false), Description: types.StringValue(""),
			},
			validate: func(t *testing.T, cell *customattribute.CellJson) {
				assert.False(t, cell.Op.IsSet)
				assert.True(t, cell.Md.IsSet)
				assert.Equal(t, "# Hello", cell.Md.Get())
				assert.Equal(t, false, cell.Enabled)
			},
		},
		{
			name: "Both op and md set",
			model: CellTFModel{
				Op: types.StringValue("cmd"), Md: types.StringValue("text"),
				Name: types.StringValue("x"), Enabled: types.BoolValue(true),
				SecretAware: types.BoolValue(false), Description: types.StringValue(""),
			},
			expectError: true,
			errContains: "both op and md set",
		},
		{
			name: "Neither op nor md set",
			model: CellTFModel{
				Op: types.StringNull(), Md: types.StringNull(),
				Name: types.StringValue("x"), Enabled: types.BoolValue(true),
				SecretAware: types.BoolValue(false), Description: types.StringValue(""),
			},
			expectError: true,
			errContains: "must have either op or md set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell, err := cellTFModelToInternal(&tt.model, 0)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cell)
				tt.validate(t, cell)
			}
		})
	}
}

func TestInternalCellToTFObject(t *testing.T) {
	tests := []struct {
		name     string
		internal *customattribute.CellJson
		validate func(t *testing.T, obj types.Object)
	}{
		{
			name: "Op cell",
			internal: &customattribute.CellJson{
				Op: common.NewOptional("host"), Md: common.NewOptionalUnset[string](),
				Name: "my_cell", Enabled: true, SecretAware: true, Description: "desc",
			},
			validate: func(t *testing.T, obj types.Object) {
				attrs := obj.Attributes()
				assert.Equal(t, "host", attrs["op"].(types.String).ValueString())
				assert.True(t, attrs["md"].(types.String).IsNull())
				assert.Equal(t, "my_cell", attrs["name"].(types.String).ValueString())
				assert.Equal(t, true, attrs["enabled"].(types.Bool).ValueBool())
				assert.Equal(t, true, attrs["secret_aware"].(types.Bool).ValueBool())
				assert.Equal(t, "desc", attrs["description"].(types.String).ValueString())
			},
		},
		{
			name: "Md cell",
			internal: &customattribute.CellJson{
				Op: common.NewOptionalUnset[string](), Md: common.NewOptional("# Title"),
				Name: "unnamed", Enabled: false, SecretAware: false, Description: "",
			},
			validate: func(t *testing.T, obj types.Object) {
				attrs := obj.Attributes()
				assert.True(t, attrs["op"].(types.String).IsNull())
				assert.Equal(t, "# Title", attrs["md"].(types.String).ValueString())
				assert.Equal(t, false, attrs["enabled"].(types.Bool).ValueBool())
			},
		},
		{
			name: "Neither op nor md set",
			internal: &customattribute.CellJson{
				Op: common.NewOptionalUnset[string](), Md: common.NewOptionalUnset[string](),
				Name: "empty", Enabled: true, SecretAware: false, Description: "",
			},
			validate: func(t *testing.T, obj types.Object) {
				attrs := obj.Attributes()
				assert.True(t, attrs["op"].(types.String).IsNull())
				assert.True(t, attrs["md"].(types.String).IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, diags := internalCellToTFObject(tt.internal)
			require.False(t, diags.HasError(), "unexpected diags: %v", diags)
			tt.validate(t, obj)
		})
	}
}

// Helper types and functions

type cellInput struct {
	op          *string
	md          *string
	name        string
	enabled     bool
	secretAware bool
	description string
}

func strPtr(s string) *string {
	return &s
}

func createCellsList(t *testing.T, cells []cellInput) types.List {
	t.Helper()
	rawCells := make([]map[string]attr.Value, len(cells))
	for i, c := range cells {
		attrs := map[string]attr.Value{
			"name":         types.StringValue(c.name),
			"enabled":      types.BoolValue(c.enabled),
			"secret_aware": types.BoolValue(c.secretAware),
			"description":  types.StringValue(c.description),
		}
		if c.op != nil {
			attrs["op"] = types.StringValue(*c.op)
		} else {
			attrs["op"] = types.StringNull()
		}
		if c.md != nil {
			attrs["md"] = types.StringValue(*c.md)
		} else {
			attrs["md"] = types.StringNull()
		}
		rawCells[i] = attrs
	}
	return createCellsListRaw(t, rawCells)
}

func createCellsListRaw(t *testing.T, cells []map[string]attr.Value) types.List {
	t.Helper()
	objects := make([]attr.Value, len(cells))
	for i, c := range cells {
		obj, diags := types.ObjectValue(CellsListAttrTypes, c)
		require.False(t, diags.HasError(), "creating cell object: %v", diags)
		objects[i] = obj
	}
	list, diags := types.ListValue(CellsListObjectType, objects)
	require.False(t, diags.HasError(), "creating cells list: %v", diags)
	return list
}
