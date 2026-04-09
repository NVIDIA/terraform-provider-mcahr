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
	"terraform/terraform-provider/provider/tf/resource/dashboard/model"
	converters "terraform/terraform-provider/provider/tf/resource/dashboard/translator/object_converters"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildGroupsList(groups ...map[string]attr.Value) types.List {
	objects := make([]attr.Value, len(groups))
	for i, g := range groups {
		obj, _ := types.ObjectValue(converters.GroupsListAttrTypes, g)
		objects[i] = obj
	}
	list, _ := types.ListValue(converters.GroupsListObjectType, objects)
	return list
}

func buildValuesList(values ...map[string]attr.Value) types.List {
	objects := make([]attr.Value, len(values))
	for i, v := range values {
		obj, _ := types.ObjectValue(converters.ValuesListAttrTypes, v)
		objects[i] = obj
	}
	list, _ := types.ListValue(converters.ValuesListObjectType, objects)
	return list
}

func tagsListValue(tags ...string) types.List {
	v, _ := types.ListValueFrom(context.Background(), types.StringType, tags)
	return v
}

func valsListValue(vals ...string) types.List {
	v, _ := types.ListValueFrom(context.Background(), types.StringType, vals)
	return v
}

func defaultBackendVersion() *version.BackendVersion {
	return &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0}
}

func TestPopulateFullJsonAttributes(t *testing.T) {
	tests := []struct {
		name           string
		resultValues   *model.DashboardTFModel
		plan           *model.DashboardTFModel
		state          *model.DashboardTFModel
		backendVersion *version.BackendVersion
		expectError    bool
		validate       func(t *testing.T, result *model.DashboardTFModel)
	}{
		{
			name: "Process groups field",
			resultValues: &model.DashboardTFModel{
				Groups: types.StringValue(`[{"name":"g1","tags":["tag1"]}]`),
			},
			plan: &model.DashboardTFModel{
				GroupsFull: types.StringValue("not_null"),
			},
			state:          &model.DashboardTFModel{},
			backendVersion: defaultBackendVersion(),
			validate: func(t *testing.T, result *model.DashboardTFModel) {
				assert.Equal(t, `[{"name":"g1","tags":["tag1"]}]`, result.GroupsFull.ValueString())
			},
		},
		{
			name: "Process values field",
			resultValues: &model.DashboardTFModel{
				Values: types.StringValue(`[{"color":"#78909c","values":["aws"]}]`),
			},
			plan: &model.DashboardTFModel{
				ValuesFull: types.StringValue("not_null"),
			},
			state:          &model.DashboardTFModel{},
			backendVersion: defaultBackendVersion(),
			validate: func(t *testing.T, result *model.DashboardTFModel) {
				assert.Equal(t, `[{"color":"#78909c","values":["aws"]}]`, result.ValuesFull.ValueString())
			},
		},
		{
			name: "Process both fields together",
			resultValues: &model.DashboardTFModel{
				Groups: types.StringValue(`[{"name":"g1","tags":["tag1"]}]`),
				Values: types.StringValue(`[{"color":"red","values":["v1"]}]`),
			},
			plan: &model.DashboardTFModel{
				GroupsFull: types.StringValue("not_null"),
				ValuesFull: types.StringValue("not_null"),
			},
			state:          &model.DashboardTFModel{},
			backendVersion: defaultBackendVersion(),
			validate: func(t *testing.T, result *model.DashboardTFModel) {
				assert.Equal(t, `[{"name":"g1","tags":["tag1"]}]`, result.GroupsFull.ValueString())
				assert.Equal(t, `[{"color":"red","values":["v1"]}]`, result.ValuesFull.ValueString())
			},
		},
		{
			name: "Skip delete operation - null groups_full in plan",
			resultValues: &model.DashboardTFModel{
				Groups: types.StringValue(`[{"name":"g1","tags":[]}]`),
			},
			plan: &model.DashboardTFModel{
				GroupsFull: types.StringNull(),
			},
			state:          &model.DashboardTFModel{},
			backendVersion: defaultBackendVersion(),
			validate: func(t *testing.T, result *model.DashboardTFModel) {
				assert.True(t, result.GroupsFull.IsNull())
			},
		},
		{
			name: "Skip delete operation - null values_full in plan",
			resultValues: &model.DashboardTFModel{
				Values: types.StringValue(`[{"color":"red","values":[]}]`),
			},
			plan: &model.DashboardTFModel{
				ValuesFull: types.StringNull(),
			},
			state:          &model.DashboardTFModel{},
			backendVersion: defaultBackendVersion(),
			validate: func(t *testing.T, result *model.DashboardTFModel) {
				assert.True(t, result.ValuesFull.IsNull())
			},
		},
		{
			name: "Handle invalid JSON in groups",
			resultValues: &model.DashboardTFModel{
				Groups: types.StringValue(`{invalid json}`),
			},
			plan: &model.DashboardTFModel{
				GroupsFull: types.StringValue("not_null"),
			},
			state:          &model.DashboardTFModel{},
			backendVersion: defaultBackendVersion(),
			expectError:    true,
		},
		{
			name: "Handle invalid JSON in values",
			resultValues: &model.DashboardTFModel{
				Values: types.StringValue(`{invalid json}`),
			},
			plan: &model.DashboardTFModel{
				ValuesFull: types.StringValue("not_null"),
			},
			state:          &model.DashboardTFModel{},
			backendVersion: defaultBackendVersion(),
			expectError:    true,
		},
		{
			name: "Empty array JSON",
			resultValues: &model.DashboardTFModel{
				Groups: types.StringValue(`[]`),
				Values: types.StringValue(`[]`),
			},
			plan: &model.DashboardTFModel{
				GroupsFull: types.StringValue("not_null"),
				ValuesFull: types.StringValue("not_null"),
			},
			state:          &model.DashboardTFModel{},
			backendVersion: defaultBackendVersion(),
			validate: func(t *testing.T, result *model.DashboardTFModel) {
				assert.False(t, result.GroupsFull.IsNull())
				assert.False(t, result.ValuesFull.IsNull())
				assert.Equal(t, "[]", result.GroupsFull.ValueString())
				assert.Equal(t, "[]", result.ValuesFull.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			err := PopulateFullJsonAttributes(ctx, tt.resultValues, tt.resultValues, tt.plan, tt.state, tt.backendVersion)

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

func TestShouldSkipForReplacement(t *testing.T) {
	groupsConfig := JSON_ATTRIBUTES_TO_POPULATE["groups"]
	valuesConfig := JSON_ATTRIBUTES_TO_POPULATE["values"]

	groupsList := buildGroupsList(map[string]attr.Value{
		"name": types.StringValue("g1"),
		"tags": tagsListValue("tag1"),
	})
	valuesList := buildValuesList(map[string]attr.Value{
		"color":  types.StringValue("#78909c"),
		"values": valsListValue("aws"),
	})

	t.Run("groups_list active and groups not explicitly set → skip and null", func(t *testing.T) {
		resultValues := &model.DashboardTFModel{
			GroupsList: groupsList,
			Groups:     types.StringValue("[]"),
			GroupsFull: types.StringValue("[]"),
		}
		resultWithoutDefaults := &model.DashboardTFModel{
			Groups: types.StringNull(),
		}
		plan := &model.DashboardTFModel{
			GroupsFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(groupsConfig, resultValues, resultWithoutDefaults, plan)

		assert.True(t, skipped)
		assert.True(t, resultValues.Groups.IsNull())
		assert.True(t, resultValues.GroupsFull.IsNull())
		assert.True(t, plan.GroupsFull.IsNull())
		assert.Equal(t, groupsList, resultValues.GroupsList, "groups_list should remain unchanged")
	})

	t.Run("values_list active and values not explicitly set → skip and null", func(t *testing.T) {
		resultValues := &model.DashboardTFModel{
			ValuesList: valuesList,
			Values:     types.StringValue("[]"),
			ValuesFull: types.StringValue("[]"),
		}
		resultWithoutDefaults := &model.DashboardTFModel{
			Values: types.StringNull(),
		}
		plan := &model.DashboardTFModel{
			ValuesFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(valuesConfig, resultValues, resultWithoutDefaults, plan)

		assert.True(t, skipped)
		assert.True(t, resultValues.Values.IsNull())
		assert.True(t, resultValues.ValuesFull.IsNull())
		assert.True(t, plan.ValuesFull.IsNull())
		assert.Equal(t, valuesList, resultValues.ValuesList, "values_list should remain unchanged")
	})

	t.Run("groups_list active but groups explicitly set → don't skip (conflict case)", func(t *testing.T) {
		resultValues := &model.DashboardTFModel{
			GroupsList: groupsList,
			Groups:     types.StringValue(`[{"name":"g1","tags":[]}]`),
			GroupsFull: types.StringValue("something"),
		}
		resultWithoutDefaults := &model.DashboardTFModel{
			Groups: types.StringValue(`[{"name":"g1","tags":[]}]`),
		}
		plan := &model.DashboardTFModel{
			GroupsFull: types.StringValue("not_null"),
		}

		skipped := shouldSkipForReplacement(groupsConfig, resultValues, resultWithoutDefaults, plan)

		assert.False(t, skipped)
		assert.Equal(t, `[{"name":"g1","tags":[]}]`, resultValues.Groups.ValueString())
	})

	t.Run("groups_list not active → don't skip", func(t *testing.T) {
		resultValues := &model.DashboardTFModel{
			GroupsList: types.ListNull(converters.GroupsListObjectType),
			Groups:     types.StringValue(`[{"name":"g1","tags":[]}]`),
		}
		resultWithoutDefaults := &model.DashboardTFModel{
			Groups: types.StringValue(`[{"name":"g1","tags":[]}]`),
		}
		plan := &model.DashboardTFModel{}

		skipped := shouldSkipForReplacement(groupsConfig, resultValues, resultWithoutDefaults, plan)

		assert.False(t, skipped)
	})
}

func TestIsDeleteOperation(t *testing.T) {
	tests := []struct {
		name       string
		plan       *model.DashboardTFModel
		attrConfig JsonAttributeConfig
		expected   bool
	}{
		{
			name: "Delete operation - null groups_full",
			plan: &model.DashboardTFModel{
				GroupsFull: types.StringNull(),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["groups"],
			expected:   true,
		},
		{
			name: "Not delete operation - non-null groups_full",
			plan: &model.DashboardTFModel{
				GroupsFull: types.StringValue("some_value"),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["groups"],
			expected:   false,
		},
		{
			name: "Delete operation - null values_full",
			plan: &model.DashboardTFModel{
				ValuesFull: types.StringNull(),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["values"],
			expected:   true,
		},
		{
			name: "Not delete operation - empty string values_full",
			plan: &model.DashboardTFModel{
				ValuesFull: types.StringValue(""),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["values"],
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDeleteOperation(tt.plan, tt.attrConfig)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJsonAttributeConfig_Functions(t *testing.T) {
	tests := []struct {
		name       string
		attribute  string
		setupModel func() *model.DashboardTFModel
	}{
		{
			name:      "Groups attribute config",
			attribute: "groups",
			setupModel: func() *model.DashboardTFModel {
				return &model.DashboardTFModel{
					Groups:     types.StringValue("groups_value"),
					GroupsFull: types.StringValue("groups_full_value"),
				}
			},
		},
		{
			name:      "Values attribute config",
			attribute: "values",
			setupModel: func() *model.DashboardTFModel {
				return &model.DashboardTFModel{
					Values:     types.StringValue("values_value"),
					ValuesFull: types.StringValue("values_full_value"),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrConfig := JSON_ATTRIBUTES_TO_POPULATE[tt.attribute]
			m := tt.setupModel()

			regularValue := attrConfig.GetAttr(m)
			fullValue := attrConfig.GetFullAttr(m)

			switch tt.attribute {
			case "groups":
				assert.Equal(t, "groups_value", regularValue.ValueString())
				assert.Equal(t, "groups_full_value", fullValue.ValueString())
			case "values":
				assert.Equal(t, "values_value", regularValue.ValueString())
				assert.Equal(t, "values_full_value", fullValue.ValueString())
			}

			newFullValue := types.StringValue("new_full_value")
			attrConfig.SetFullAttr(m, newFullValue)
			assert.Equal(t, "new_full_value", attrConfig.GetFullAttr(m).ValueString())

			newBaseValue := types.StringValue("new_base_value")
			attrConfig.SetAttr(m, newBaseValue)
			assert.Equal(t, "new_base_value", attrConfig.GetAttr(m).ValueString())
		})
	}
}

func TestJSON_ATTRIBUTES_TO_POPULATE_Configuration(t *testing.T) {
	expectedAttributes := []string{"groups", "values"}

	for _, a := range expectedAttributes {
		t.Run("Has_"+a, func(t *testing.T) {
			attrConfig, exists := JSON_ATTRIBUTES_TO_POPULATE[a]
			assert.True(t, exists, "Expected attribute %s to be configured", a)

			switch a {
			case "groups":
				assert.Equal(t, "groups_full", attrConfig.FullAttrName)
			case "values":
				assert.Equal(t, "values_full", attrConfig.FullAttrName)
			}

			assert.NotNil(t, attrConfig.RemarshalFunc)
			assert.NotNil(t, attrConfig.GetAttr)
			assert.NotNil(t, attrConfig.SetAttr)
			assert.NotNil(t, attrConfig.GetFullAttr)
			assert.NotNil(t, attrConfig.SetFullAttr)
			assert.NotNil(t, attrConfig.GetReplacementAttr, "dashboard attributes should have replacement attrs")
		})
	}

	assert.Equal(t, len(expectedAttributes), len(JSON_ATTRIBUTES_TO_POPULATE))
}
