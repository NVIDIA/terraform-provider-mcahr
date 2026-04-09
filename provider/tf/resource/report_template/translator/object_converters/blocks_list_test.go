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

	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fullBlock() customattribute.BlockJson {
	return customattribute.BlockJson{
		Title:                           "Block 1",
		ResourceQuery:                   "host",
		GroupByTag:                      "environment",
		BreakdownByTag:                  "service",
		ViewMode:                        "COUNT",
		IncludeResourcesWithoutGroupTag: true,
		IncludeOtherBreakdownTagValues:  false,
		OtherTagsToExport:               []string{"tag1", "tag2"},
		GroupByTagOrder:                 customattribute.GroupByTagOrder{Type: "CUSTOM", Values: []string{"prod", "staging"}},
		BreakdownTagsValues: []customattribute.BreakdownTagValue{
			{Color: "#FF0000", Label: "Production", Values: []string{"prod", "production"}},
		},
		ResourcesBreakdown: []customattribute.ResourcesBreakdown{
			{GroupByValue: "env", BreakdownValues: []customattribute.BreakdownValue{{Value: "prod", Count: 5}}},
		},
	}
}

func minimalBlock() customattribute.BlockJson {
	return customattribute.BlockJson{
		Title:         "Minimal",
		ResourceQuery: "host",
	}
}

// --- BlocksListFromAPI ---

func TestBlocksListFromAPI_FullBlock(t *testing.T) {
	t.Parallel()
	apiBlocks := []customattribute.BlockJson{fullBlock()}

	result, diags := BlocksListFromAPI(apiBlocks)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	require.Equal(t, 1, len(result.Elements()))

	obj := result.Elements()[0].(types.Object)
	assert.Equal(t, "Block 1", obj.Attributes()["title"].(types.String).ValueString())
	assert.Equal(t, "host", obj.Attributes()["resource_query"].(types.String).ValueString())
	assert.Equal(t, "environment", obj.Attributes()["group_by_tag"].(types.String).ValueString())
	assert.Equal(t, "service", obj.Attributes()["breakdown_by_tag"].(types.String).ValueString())
	assert.Equal(t, "COUNT", obj.Attributes()["view_mode"].(types.String).ValueString())
	assert.Equal(t, true, obj.Attributes()["include_resources_without_group_tag"].(types.Bool).ValueBool())
	assert.Equal(t, false, obj.Attributes()["include_other_breakdown_tag_values"].(types.Bool).ValueBool())
}

func TestBlocksListFromAPI_MinimalBlock(t *testing.T) {
	t.Parallel()
	apiBlocks := []customattribute.BlockJson{minimalBlock()}

	result, diags := BlocksListFromAPI(apiBlocks)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	require.Equal(t, 1, len(result.Elements()))

	obj := result.Elements()[0].(types.Object)
	assert.Equal(t, "Minimal", obj.Attributes()["title"].(types.String).ValueString())
	assert.Equal(t, "host", obj.Attributes()["resource_query"].(types.String).ValueString())
	assert.Equal(t, "", obj.Attributes()["group_by_tag"].(types.String).ValueString())
	assert.Equal(t, "", obj.Attributes()["view_mode"].(types.String).ValueString())
}

func TestBlocksListFromAPI_MultipleBlocks(t *testing.T) {
	t.Parallel()
	apiBlocks := []customattribute.BlockJson{fullBlock(), minimalBlock()}

	result, diags := BlocksListFromAPI(apiBlocks)

	assert.False(t, diags.HasError())
	require.Equal(t, 2, len(result.Elements()))

	obj0 := result.Elements()[0].(types.Object)
	assert.Equal(t, "Block 1", obj0.Attributes()["title"].(types.String).ValueString())
	assert.Equal(t, "COUNT", obj0.Attributes()["view_mode"].(types.String).ValueString())

	obj1 := result.Elements()[1].(types.Object)
	assert.Equal(t, "Minimal", obj1.Attributes()["title"].(types.String).ValueString())
	assert.Equal(t, "host", obj1.Attributes()["resource_query"].(types.String).ValueString())
}

func TestBlocksListFromAPI_EmptySlice(t *testing.T) {
	t.Parallel()
	result, diags := BlocksListFromAPI([]customattribute.BlockJson{})

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Empty(t, result.Elements())
}

func TestBlocksListFromAPI_NilSlice(t *testing.T) {
	t.Parallel()
	result, diags := BlocksListFromAPI(nil)

	assert.False(t, diags.HasError())
	assert.True(t, result.IsNull())
}

// --- BlocksListToInternal ---

func TestBlocksListToInternal_NullList(t *testing.T) {
	t.Parallel()
	result, err := BlocksListToInternal(context.Background(), BlocksListNullValue())

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestBlocksListToInternal_UnknownList(t *testing.T) {
	t.Parallel()
	result, err := BlocksListToInternal(context.Background(), BlocksListUnknownValue())

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestBlocksListToInternal_EmptyList(t *testing.T) {
	t.Parallel()
	tfList, diags := BlocksListFromAPI([]customattribute.BlockJson{})
	require.False(t, diags.HasError())

	result, err := BlocksListToInternal(context.Background(), tfList)

	require.NoError(t, err)
	assert.Empty(t, result)
}

// --- Round-trip ---

func TestBlocksListRoundTrip_FullBlock(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	original := []customattribute.BlockJson{fullBlock()}

	tfList, diags := BlocksListFromAPI(original)
	require.False(t, diags.HasError())

	roundTripped, err := BlocksListToInternal(ctx, tfList)
	require.NoError(t, err)
	require.Len(t, roundTripped, 1)

	b := roundTripped[0]
	assert.Equal(t, "Block 1", b.Title)
	assert.Equal(t, "host", b.ResourceQuery)
	assert.Equal(t, "environment", b.GroupByTag)
	assert.Equal(t, "service", b.BreakdownByTag)
	assert.Equal(t, "COUNT", b.ViewMode)
	assert.True(t, b.IncludeResourcesWithoutGroupTag)
	assert.False(t, b.IncludeOtherBreakdownTagValues)
	assert.Equal(t, []string{"tag1", "tag2"}, b.OtherTagsToExport)
	assert.Equal(t, "CUSTOM", b.GroupByTagOrder.Type)
	assert.Equal(t, []string{"prod", "staging"}, b.GroupByTagOrder.Values)
	require.Len(t, b.BreakdownTagsValues, 1)
	assert.Equal(t, "#FF0000", b.BreakdownTagsValues[0].Color)
	assert.Equal(t, "Production", b.BreakdownTagsValues[0].Label)
	assert.Equal(t, []string{"prod", "production"}, b.BreakdownTagsValues[0].Values)
	require.Len(t, b.ResourcesBreakdown, 1)
	assert.Equal(t, "env", b.ResourcesBreakdown[0].GroupByValue)
	require.Len(t, b.ResourcesBreakdown[0].BreakdownValues, 1)
	assert.Equal(t, "prod", b.ResourcesBreakdown[0].BreakdownValues[0].Value)
	assert.Equal(t, 5, b.ResourcesBreakdown[0].BreakdownValues[0].Count)
}

func TestBlocksListRoundTrip_MinimalBlock(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	original := []customattribute.BlockJson{minimalBlock()}

	tfList, diags := BlocksListFromAPI(original)
	require.False(t, diags.HasError())

	roundTripped, err := BlocksListToInternal(ctx, tfList)
	require.NoError(t, err)
	require.Len(t, roundTripped, 1)

	b := roundTripped[0]
	assert.Equal(t, "Minimal", b.Title)
	assert.Equal(t, "host", b.ResourceQuery)
	assert.Equal(t, "", b.GroupByTag)
	assert.Equal(t, "", b.BreakdownByTag)
	assert.False(t, b.IncludeResourcesWithoutGroupTag)
}

func TestBlocksListRoundTrip_MultipleBlocks(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	original := []customattribute.BlockJson{fullBlock(), minimalBlock()}

	tfList, diags := BlocksListFromAPI(original)
	require.False(t, diags.HasError())

	roundTripped, err := BlocksListToInternal(ctx, tfList)
	require.NoError(t, err)
	require.Len(t, roundTripped, 2)

	assert.Equal(t, "Block 1", roundTripped[0].Title)
	assert.Equal(t, "Minimal", roundTripped[1].Title)
}

// --- Nested type tests ---

func TestGroupByTagOrderToTF_RoundTrip(t *testing.T) {
	t.Parallel()
	original := customattribute.GroupByTagOrder{Type: "CUSTOM", Values: []string{"a", "b", "c"}}

	tfObj, diags := groupByTagOrderToTF(original)
	require.False(t, diags.HasError())

	roundTripped := groupByTagOrderFromAttr(context.Background(), tfObj)
	assert.Equal(t, original.Type, roundTripped.Type)
	assert.Equal(t, original.Values, roundTripped.Values)
}

func TestGroupByTagOrderFromAttr_NullObject(t *testing.T) {
	t.Parallel()
	result := groupByTagOrderFromAttr(context.Background(), nil)
	assert.Equal(t, "DEFAULT", result.Type)
	assert.Equal(t, []string{}, result.Values)
}

func TestBreakdownTagValuesRoundTrip(t *testing.T) {
	t.Parallel()
	original := []customattribute.BreakdownTagValue{
		{Color: "#FF0000", Label: "Prod", Values: []string{"prod", "production"}},
		{Color: "#00FF00", Label: "Dev", Values: []string{"dev"}},
	}

	tfList, diags := breakdownTagValuesToTF(original)
	require.False(t, diags.HasError())

	roundTripped := breakdownTagValuesFromAttr(context.Background(), tfList)
	require.Len(t, roundTripped, 2)
	assert.Equal(t, "#FF0000", roundTripped[0].Color)
	assert.Equal(t, "Prod", roundTripped[0].Label)
	assert.Equal(t, []string{"prod", "production"}, roundTripped[0].Values)
	assert.Equal(t, "#00FF00", roundTripped[1].Color)
}

func TestBreakdownTagValuesToTF_NilSlice(t *testing.T) {
	t.Parallel()
	result, diags := breakdownTagValuesToTF(nil)
	assert.False(t, diags.HasError())
	assert.Empty(t, result.Elements())
}

func TestResourcesBreakdownRoundTrip(t *testing.T) {
	t.Parallel()
	original := []customattribute.ResourcesBreakdown{
		{
			GroupByValue: "env",
			BreakdownValues: []customattribute.BreakdownValue{
				{Value: "prod", Count: 10},
				{Value: "staging", Count: 3},
			},
		},
	}

	tfList, diags := resourcesBreakdownToTF(original)
	require.False(t, diags.HasError())

	roundTripped := resourcesBreakdownFromAttr(context.Background(), tfList)
	require.Len(t, roundTripped, 1)
	assert.Equal(t, "env", roundTripped[0].GroupByValue)
	require.Len(t, roundTripped[0].BreakdownValues, 2)
	assert.Equal(t, "prod", roundTripped[0].BreakdownValues[0].Value)
	assert.Equal(t, 10, roundTripped[0].BreakdownValues[0].Count)
}

func TestResourcesBreakdownToTF_NilSlice(t *testing.T) {
	t.Parallel()
	result, diags := resourcesBreakdownToTF(nil)
	assert.False(t, diags.HasError())
	assert.Empty(t, result.Elements())
}

// --- Attr types structure ---

func TestBlocksListAttrTypesStructure(t *testing.T) {
	t.Parallel()
	expected := []string{
		"title", "resource_query", "group_by_tag", "breakdown_by_tag",
		"view_mode", "include_resources_without_group_tag", "include_other_breakdown_tag_values",
		"other_tags_to_export", "group_by_tag_order", "breakdown_tags_values", "resources_breakdown",
	}
	assert.Equal(t, len(expected), len(BlocksListAttrTypes))
	for _, key := range expected {
		_, ok := BlocksListAttrTypes[key]
		assert.True(t, ok, "missing key %s", key)
	}
}

// --- Helpers for null/unknown ---

func BlocksListNullValue() types.List {
	return types.ListNull(BlocksListObjectType)
}

func BlocksListUnknownValue() types.List {
	return types.ListUnknown(BlocksListObjectType)
}
