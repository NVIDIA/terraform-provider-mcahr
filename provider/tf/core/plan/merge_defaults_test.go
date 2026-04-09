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

package plan

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Test models and helpers for AddDefaultsFromPlan ---

type addDefaultsTestModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	TimeoutMs   types.Int64  `tfsdk:"timeout_ms"`
	Labels      types.List   `tfsdk:"labels"`
}

type nestedTestModel struct {
	Name    types.String `tfsdk:"name"`
	Enabled types.Bool   `tfsdk:"enabled"`
	Items   types.List   `tfsdk:"items"`
}

var testItemAttrTypes = map[string]attr.Type{
	"name":    types.StringType,
	"enabled": types.BoolType,
}

var testItemObjectType = types.ObjectType{AttrTypes: testItemAttrTypes}

func testNestedItemObject(t *testing.T, name attr.Value, enabled attr.Value) types.Object {
	t.Helper()
	o, diags := types.ObjectValue(testItemAttrTypes, map[string]attr.Value{
		"name":    name,
		"enabled": enabled,
	})
	require.False(t, diags.HasError())
	return o
}

func testNestedItemList(t *testing.T, elems ...types.Object) types.List {
	t.Helper()
	vals := make([]attr.Value, len(elems))
	for i := range elems {
		vals[i] = elems[i]
	}
	l, diags := types.ListValue(testItemObjectType, vals)
	require.False(t, diags.HasError())
	return l
}

// --- AddDefaultsFromPlan tests ---

func TestAddDefaultsFromPlan(t *testing.T) {
	tests := []struct {
		name         string
		resultValues *addDefaultsTestModel
		planValues   *addDefaultsTestModel
		validate     func(t *testing.T, result *addDefaultsTestModel)
	}{
		{
			name: "Copy null fields from plan",
			resultValues: &addDefaultsTestModel{
				Name:    types.StringValue("test"),
				Enabled: types.BoolNull(),
			},
			planValues: &addDefaultsTestModel{
				Name:    types.StringValue("test"),
				Enabled: types.BoolValue(true),
			},
			validate: func(t *testing.T, result *addDefaultsTestModel) {
				assert.Equal(t, "test", result.Name.ValueString())
				assert.Equal(t, true, result.Enabled.ValueBool())
			},
		},
		{
			name: "Copy unknown fields from plan",
			resultValues: &addDefaultsTestModel{
				Name:        types.StringValue("test"),
				Description: types.StringUnknown(),
			},
			planValues: &addDefaultsTestModel{
				Name:        types.StringValue("test"),
				Description: types.StringValue("from plan"),
			},
			validate: func(t *testing.T, result *addDefaultsTestModel) {
				assert.Equal(t, "test", result.Name.ValueString())
				assert.Equal(t, "from plan", result.Description.ValueString())
			},
		},
		{
			name: "Don't override non-null/unknown fields",
			resultValues: &addDefaultsTestModel{
				Name:        types.StringValue("result"),
				Description: types.StringValue("result desc"),
			},
			planValues: &addDefaultsTestModel{
				Name:        types.StringValue("plan"),
				Description: types.StringValue("plan desc"),
			},
			validate: func(t *testing.T, result *addDefaultsTestModel) {
				assert.Equal(t, "result", result.Name.ValueString())
				assert.Equal(t, "result desc", result.Description.ValueString())
			},
		},
		{
			name: "Handle all field types",
			resultValues: &addDefaultsTestModel{
				Name:      types.StringNull(),
				Enabled:   types.BoolNull(),
				TimeoutMs: types.Int64Null(),
				Labels:    types.ListNull(types.StringType),
			},
			planValues: &addDefaultsTestModel{
				Name:      types.StringValue("from plan"),
				Enabled:   types.BoolValue(true),
				TimeoutMs: types.Int64Value(5000),
				Labels: func() types.List {
					s, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"label1"})
					return s
				}(),
			},
			validate: func(t *testing.T, result *addDefaultsTestModel) {
				assert.Equal(t, "from plan", result.Name.ValueString())
				assert.Equal(t, true, result.Enabled.ValueBool())
				assert.Equal(t, int64(5000), result.TimeoutMs.ValueInt64())
				assert.False(t, result.Labels.IsNull())
			},
		},
		{
			name: "Mixed null, unknown, and known fields",
			resultValues: &addDefaultsTestModel{
				Name:        types.StringNull(),
				Description: types.StringUnknown(),
				Enabled:     types.BoolValue(false),
				TimeoutMs:   types.Int64Value(1000),
			},
			planValues: &addDefaultsTestModel{
				Name:        types.StringValue("plan name"),
				Description: types.StringValue("plan desc"),
				Enabled:     types.BoolValue(true),
				TimeoutMs:   types.Int64Value(2000),
			},
			validate: func(t *testing.T, result *addDefaultsTestModel) {
				assert.Equal(t, "plan name", result.Name.ValueString())
				assert.Equal(t, "plan desc", result.Description.ValueString())
				assert.Equal(t, false, result.Enabled.ValueBool())
				assert.Equal(t, int64(1000), result.TimeoutMs.ValueInt64())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			AddDefaultsFromPlan(tt.resultValues, tt.planValues)
			if tt.validate != nil {
				tt.validate(t, tt.resultValues)
			}
		})
	}
}

func TestAddDefaultsFromPlan_NestedListDefaults(t *testing.T) {
	t.Parallel()

	result := &nestedTestModel{
		Name:    types.StringValue("root"),
		Enabled: types.BoolNull(),
		Items: testNestedItemList(t,
			testNestedItemObject(t, types.StringValue("user-set"), types.BoolNull()),
		),
	}
	plan := &nestedTestModel{
		Name:    types.StringValue("plan-root"),
		Enabled: types.BoolValue(true),
		Items: testNestedItemList(t,
			testNestedItemObject(t, types.StringValue("plan-default"), types.BoolValue(true)),
		),
	}

	AddDefaultsFromPlan(result, plan)

	assert.Equal(t, "root", result.Name.ValueString())
	assert.True(t, result.Enabled.ValueBool())
	elems := result.Items.Elements()
	require.Len(t, elems, 1)
	obj := elems[0].(types.Object)
	attrs := obj.Attributes()
	assert.Equal(t, "user-set", attrs["name"].(types.String).ValueString())
	assert.True(t, attrs["enabled"].(types.Bool).ValueBool())
}

func TestAddDefaultsFromPlan_NonListFieldsUnchanged(t *testing.T) {
	t.Parallel()

	result := &nestedTestModel{
		Name:    types.StringNull(),
		Enabled: types.BoolValue(false),
		Items: testNestedItemList(t,
			testNestedItemObject(t, types.StringValue("keep"), types.BoolValue(true)),
		),
	}
	plan := &nestedTestModel{
		Name:    types.StringValue("from-plan"),
		Enabled: types.BoolValue(true),
		Items: testNestedItemList(t,
			testNestedItemObject(t, types.StringValue("plan-item"), types.BoolNull()),
		),
	}

	AddDefaultsFromPlan(result, plan)

	assert.Equal(t, "from-plan", result.Name.ValueString())
	assert.False(t, result.Enabled.ValueBool())
	elems := result.Items.Elements()
	require.Len(t, elems, 1)
	obj := elems[0].(types.Object)
	attrs := obj.Attributes()
	assert.Equal(t, "keep", attrs["name"].(types.String).ValueString())
	assert.True(t, attrs["enabled"].(types.Bool).ValueBool())
}

func TestAddDefaultsFromPlan_EmptyOrNullListsNoElementMerge(t *testing.T) {
	t.Parallel()

	t.Run("both empty known lists unchanged", func(t *testing.T) {
		t.Parallel()
		empty, diags := types.ListValue(testItemObjectType, []attr.Value{})
		require.False(t, diags.HasError())

		result := &nestedTestModel{Name: types.StringValue("n"), Enabled: types.BoolValue(true), Items: empty}
		plan := &nestedTestModel{Name: types.StringValue("n"), Enabled: types.BoolValue(true), Items: empty}

		AddDefaultsFromPlan(result, plan)
		assert.True(t, result.Items.IsNull() || len(result.Items.Elements()) == 0)
	})

	t.Run("result list known plan list null no merge", func(t *testing.T) {
		t.Parallel()
		before := testNestedItemList(t, testNestedItemObject(t, types.StringValue("only"), types.BoolNull()))
		result := &nestedTestModel{Name: types.StringValue("n"), Enabled: types.BoolValue(true), Items: before}
		plan := &nestedTestModel{Name: types.StringValue("n"), Enabled: types.BoolValue(true), Items: types.ListNull(testItemObjectType)}

		AddDefaultsFromPlan(result, plan)
		assert.True(t, result.Items.Equal(before))
	})

	t.Run("result list null copies whole field from plan", func(t *testing.T) {
		t.Parallel()
		planList := testNestedItemList(t, testNestedItemObject(t, types.StringValue("p"), types.BoolValue(true)))
		result := &nestedTestModel{Name: types.StringValue("n"), Enabled: types.BoolValue(true), Items: types.ListNull(testItemObjectType)}
		plan := &nestedTestModel{Name: types.StringValue("n"), Enabled: types.BoolValue(true), Items: planList}

		AddDefaultsFromPlan(result, plan)
		assert.True(t, result.Items.Equal(planList))
	})
}

func TestAddDefaultsFromPlan_DifferentLengthListsNoElementMerge(t *testing.T) {
	t.Parallel()

	before := testNestedItemList(t, testNestedItemObject(t, types.StringValue("one"), types.BoolNull()))
	result := &nestedTestModel{Name: types.StringValue("n"), Enabled: types.BoolValue(true), Items: before}
	plan := &nestedTestModel{
		Name: types.StringValue("n"), Enabled: types.BoolValue(true),
		Items: testNestedItemList(t,
			testNestedItemObject(t, types.StringValue("a"), types.BoolNull()),
			testNestedItemObject(t, types.StringValue("b"), types.BoolNull()),
		),
	}

	AddDefaultsFromPlan(result, plan)
	assert.True(t, result.Items.Equal(before))
}

// --- mergeDefaults tests ---

func TestMergeDefaults_Primitives(t *testing.T) {
	tests := []struct {
		name       string
		result     attr.Value
		plan       attr.Value
		wantValue  attr.Value
		wantChange bool
	}{
		{
			name:       "null result, known plan → use plan",
			result:     types.StringNull(),
			plan:       types.StringValue("default"),
			wantValue:  types.StringValue("default"),
			wantChange: true,
		},
		{
			name:       "unknown result, known plan → use plan",
			result:     types.StringUnknown(),
			plan:       types.StringValue("default"),
			wantValue:  types.StringValue("default"),
			wantChange: true,
		},
		{
			name:       "known result, known plan → keep result",
			result:     types.StringValue("user"),
			plan:       types.StringValue("default"),
			wantValue:  types.StringValue("user"),
			wantChange: false,
		},
		{
			name:       "null result, null plan → keep null",
			result:     types.StringNull(),
			plan:       types.StringNull(),
			wantValue:  types.StringNull(),
			wantChange: false,
		},
		{
			name:       "known result, null plan → keep result",
			result:     types.BoolValue(true),
			plan:       types.BoolNull(),
			wantValue:  types.BoolValue(true),
			wantChange: false,
		},
		{
			name:       "null bool, known plan → use plan",
			result:     types.BoolNull(),
			plan:       types.BoolValue(false),
			wantValue:  types.BoolValue(false),
			wantChange: true,
		},
		{
			name:       "null result, unknown plan → propagate unknown (UseStateForUnknown on first create)",
			result:     types.StringNull(),
			plan:       types.StringUnknown(),
			wantValue:  types.StringUnknown(),
			wantChange: true,
		},
		{
			name:       "unknown result, unknown plan → no change",
			result:     types.StringUnknown(),
			plan:       types.StringUnknown(),
			wantValue:  types.StringUnknown(),
			wantChange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed := mergeDefaults(tt.result, tt.plan)
			assert.Equal(t, tt.wantChange, changed)
			assert.Equal(t, tt.wantValue, got)
		})
	}
}

func TestMergeDefaults_Object(t *testing.T) {
	objType := map[string]attr.Type{
		"name":    types.StringType,
		"enabled": types.BoolType,
	}

	t.Run("fills null attrs from plan, keeps user-set attrs", func(t *testing.T) {
		result, _ := types.ObjectValue(objType, map[string]attr.Value{
			"name":    types.StringValue("user_name"),
			"enabled": types.BoolNull(),
		})
		plan, _ := types.ObjectValue(objType, map[string]attr.Value{
			"name":    types.StringValue("default_name"),
			"enabled": types.BoolValue(true),
		})

		got, changed := mergeDefaults(result, plan)
		require.True(t, changed)

		obj := got.(types.Object)
		assert.Equal(t, "user_name", obj.Attributes()["name"].(types.String).ValueString())
		assert.Equal(t, true, obj.Attributes()["enabled"].(types.Bool).ValueBool())
	})

	t.Run("no change when all attrs set", func(t *testing.T) {
		result, _ := types.ObjectValue(objType, map[string]attr.Value{
			"name":    types.StringValue("a"),
			"enabled": types.BoolValue(false),
		})
		plan, _ := types.ObjectValue(objType, map[string]attr.Value{
			"name":    types.StringValue("b"),
			"enabled": types.BoolValue(true),
		})

		_, changed := mergeDefaults(result, plan)
		assert.False(t, changed)
	})

	t.Run("null result object → use plan", func(t *testing.T) {
		result := types.ObjectNull(objType)
		plan, _ := types.ObjectValue(objType, map[string]attr.Value{
			"name":    types.StringValue("default"),
			"enabled": types.BoolValue(true),
		})

		got, changed := mergeDefaults(result, plan)
		require.True(t, changed)
		assert.Equal(t, plan, got)
	})
}

func TestMergeDefaults_List(t *testing.T) {
	t.Run("fills null elements from plan", func(t *testing.T) {
		result, _ := types.ListValue(types.StringType, []attr.Value{
			types.StringValue("a"),
			types.StringNull(),
			types.StringValue("c"),
		})
		plan, _ := types.ListValue(types.StringType, []attr.Value{
			types.StringValue("x"),
			types.StringValue("y"),
			types.StringValue("z"),
		})

		got, changed := mergeDefaults(result, plan)
		require.True(t, changed)

		list := got.(types.List)
		assert.Equal(t, "a", list.Elements()[0].(types.String).ValueString())
		assert.Equal(t, "y", list.Elements()[1].(types.String).ValueString())
		assert.Equal(t, "c", list.Elements()[2].(types.String).ValueString())
	})

	t.Run("no change when all elements set", func(t *testing.T) {
		result, _ := types.ListValue(types.StringType, []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
		})
		plan, _ := types.ListValue(types.StringType, []attr.Value{
			types.StringValue("x"),
			types.StringValue("y"),
		})

		_, changed := mergeDefaults(result, plan)
		assert.False(t, changed)
	})

	t.Run("different length lists → no merge", func(t *testing.T) {
		result, _ := types.ListValue(types.StringType, []attr.Value{
			types.StringNull(),
		})
		plan, _ := types.ListValue(types.StringType, []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
		})

		_, changed := mergeDefaults(result, plan)
		assert.False(t, changed)
	})

	t.Run("null list → use plan", func(t *testing.T) {
		result := types.ListNull(types.StringType)
		plan, _ := types.ListValue(types.StringType, []attr.Value{types.StringValue("a")})

		got, changed := mergeDefaults(result, plan)
		require.True(t, changed)
		assert.Equal(t, plan, got)
	})
}

func TestMergeDefaults_NestedListOfObjects(t *testing.T) {
	cellType := map[string]attr.Type{
		"op":      types.StringType,
		"md":      types.StringType,
		"name":    types.StringType,
		"enabled": types.BoolType,
	}
	cellObjType := types.ObjectType{AttrTypes: cellType}

	t.Run("fills null nested attrs in list elements", func(t *testing.T) {
		resultCell, _ := types.ObjectValue(cellType, map[string]attr.Value{
			"op":      types.StringValue("host"),
			"md":      types.StringNull(),
			"name":    types.StringNull(),
			"enabled": types.BoolNull(),
		})
		planCell, _ := types.ObjectValue(cellType, map[string]attr.Value{
			"op":      types.StringValue("host"),
			"md":      types.StringNull(),
			"name":    types.StringValue("unnamed"),
			"enabled": types.BoolValue(true),
		})

		resultList, _ := types.ListValue(cellObjType, []attr.Value{resultCell})
		planList, _ := types.ListValue(cellObjType, []attr.Value{planCell})

		got, changed := mergeDefaults(resultList, planList)
		require.True(t, changed)

		list := got.(types.List)
		cell := list.Elements()[0].(types.Object)
		assert.Equal(t, "host", cell.Attributes()["op"].(types.String).ValueString())
		assert.True(t, cell.Attributes()["md"].(types.String).IsNull())
		assert.Equal(t, "unnamed", cell.Attributes()["name"].(types.String).ValueString())
		assert.Equal(t, true, cell.Attributes()["enabled"].(types.Bool).ValueBool())
	})

	t.Run("preserves user-set values in nested objects", func(t *testing.T) {
		resultCell, _ := types.ObjectValue(cellType, map[string]attr.Value{
			"op":      types.StringNull(),
			"md":      types.StringValue("# My Title"),
			"name":    types.StringValue("custom_name"),
			"enabled": types.BoolValue(false),
		})
		planCell, _ := types.ObjectValue(cellType, map[string]attr.Value{
			"op":      types.StringNull(),
			"md":      types.StringValue("# My Title"),
			"name":    types.StringValue("unnamed"),
			"enabled": types.BoolValue(true),
		})

		resultList, _ := types.ListValue(cellObjType, []attr.Value{resultCell})
		planList, _ := types.ListValue(cellObjType, []attr.Value{planCell})

		got, changed := mergeDefaults(resultList, planList)
		assert.False(t, changed)

		list := got.(types.List)
		cell := list.Elements()[0].(types.Object)
		assert.Equal(t, "# My Title", cell.Attributes()["md"].(types.String).ValueString())
		assert.Equal(t, "custom_name", cell.Attributes()["name"].(types.String).ValueString())
		assert.Equal(t, false, cell.Attributes()["enabled"].(types.Bool).ValueBool())
	})
}

func TestMergeDefaults_DeeplyNested(t *testing.T) {
	innerType := map[string]attr.Type{
		"value": types.StringType,
	}
	innerObjType := types.ObjectType{AttrTypes: innerType}

	outerType := map[string]attr.Type{
		"name":  types.StringType,
		"items": types.ListType{ElemType: innerObjType},
	}
	outerObjType := types.ObjectType{AttrTypes: outerType}

	innerResult, _ := types.ObjectValue(innerType, map[string]attr.Value{
		"value": types.StringNull(),
	})
	innerPlan, _ := types.ObjectValue(innerType, map[string]attr.Value{
		"value": types.StringValue("deep_default"),
	})

	innerListResult, _ := types.ListValue(innerObjType, []attr.Value{innerResult})
	innerListPlan, _ := types.ListValue(innerObjType, []attr.Value{innerPlan})

	outerResult, _ := types.ObjectValue(outerType, map[string]attr.Value{
		"name":  types.StringValue("outer"),
		"items": innerListResult,
	})
	outerPlan, _ := types.ObjectValue(outerType, map[string]attr.Value{
		"name":  types.StringValue("outer"),
		"items": innerListPlan,
	})

	resultList, _ := types.ListValue(outerObjType, []attr.Value{outerResult})
	planList, _ := types.ListValue(outerObjType, []attr.Value{outerPlan})

	got, changed := mergeDefaults(resultList, planList)
	require.True(t, changed)

	list := got.(types.List)
	outer := list.Elements()[0].(types.Object)
	assert.Equal(t, "outer", outer.Attributes()["name"].(types.String).ValueString())

	items := outer.Attributes()["items"].(types.List)
	inner := items.Elements()[0].(types.Object)
	assert.Equal(t, "deep_default", inner.Attributes()["value"].(types.String).ValueString())
}
