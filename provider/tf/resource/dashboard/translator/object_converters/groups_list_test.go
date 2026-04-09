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

	customattribute "terraform/terraform-provider/provider/external_api/resources/dashboards/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupsListToInternal_Success(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tagsVal, _ := types.ListValueFrom(ctx, types.StringType, []string{"tag1", "tag2"})
	obj, _ := types.ObjectValue(GroupsListAttrTypes, map[string]attr.Value{
		"name": types.StringValue("group1"),
		"tags": tagsVal,
	})
	tfList, _ := types.ListValue(GroupsListObjectType, []attr.Value{obj})

	result, err := GroupsListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "group1", result[0].Name)
	assert.Equal(t, []string{"tag1", "tag2"}, result[0].Tags)
}

func TestGroupsListToInternal_MultipleGroups(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tags1, _ := types.ListValueFrom(ctx, types.StringType, []string{"cloud_provider"})
	obj1, _ := types.ObjectValue(GroupsListAttrTypes, map[string]attr.Value{
		"name": types.StringValue("g1"),
		"tags": tags1,
	})

	tags2, _ := types.ListValueFrom(ctx, types.StringType, []string{"release_tag", "env"})
	obj2, _ := types.ObjectValue(GroupsListAttrTypes, map[string]attr.Value{
		"name": types.StringValue("g2"),
		"tags": tags2,
	})

	tfList, _ := types.ListValue(GroupsListObjectType, []attr.Value{obj1, obj2})

	result, err := GroupsListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "g1", result[0].Name)
	assert.Equal(t, []string{"cloud_provider"}, result[0].Tags)
	assert.Equal(t, "g2", result[1].Name)
	assert.Equal(t, []string{"release_tag", "env"}, result[1].Tags)
}

func TestGroupsListToInternal_EmptyList(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tfList, _ := types.ListValue(GroupsListObjectType, []attr.Value{})

	result, err := GroupsListToInternal(ctx, tfList)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestGroupsListToInternal_NullList(t *testing.T) {
	t.Parallel()

	result, err := GroupsListToInternal(context.Background(), types.ListNull(GroupsListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestGroupsListToInternal_UnknownList(t *testing.T) {
	t.Parallel()

	result, err := GroupsListToInternal(context.Background(), types.ListUnknown(GroupsListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestGroupsListToInternal_EmptyTags(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	emptyTags, _ := types.ListValueFrom(ctx, types.StringType, []string{})
	obj, _ := types.ObjectValue(GroupsListAttrTypes, map[string]attr.Value{
		"name": types.StringValue("no-tags-group"),
		"tags": emptyTags,
	})
	tfList, _ := types.ListValue(GroupsListObjectType, []attr.Value{obj})

	result, err := GroupsListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "no-tags-group", result[0].Name)
	assert.Empty(t, result[0].Tags)
}

func TestGroupsListFromAPI_Success(t *testing.T) {
	t.Parallel()

	apiGroups := []customattribute.GroupJson{
		{Name: "group1", Tags: []string{"tag1", "tag2"}},
	}

	result, diags := GroupsListFromAPI(apiGroups)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Equal(t, 1, len(result.Elements()))
}

func TestGroupsListFromAPI_MultipleGroups(t *testing.T) {
	t.Parallel()

	apiGroups := []customattribute.GroupJson{
		{Name: "g1", Tags: []string{"cloud_provider"}},
		{Name: "g2", Tags: []string{"release_tag", "env"}},
	}

	result, diags := GroupsListFromAPI(apiGroups)

	assert.False(t, diags.HasError())
	assert.Equal(t, 2, len(result.Elements()))
}

func TestGroupsListFromAPI_EmptyGroups(t *testing.T) {
	t.Parallel()

	result, diags := GroupsListFromAPI([]customattribute.GroupJson{})

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Empty(t, result.Elements())
}

func TestGroupsListFromAPI_NilGroups(t *testing.T) {
	t.Parallel()

	result, diags := GroupsListFromAPI(nil)

	assert.False(t, diags.HasError())
	assert.True(t, result.IsNull())
}

func TestGroupsListFromAPI_EmptyTags(t *testing.T) {
	t.Parallel()

	apiGroups := []customattribute.GroupJson{
		{Name: "group1", Tags: []string{}},
	}

	result, diags := GroupsListFromAPI(apiGroups)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Equal(t, 1, len(result.Elements()))
}

func TestGroupsListRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	original := []customattribute.GroupJson{
		{Name: "g1", Tags: []string{"tag1", "tag2"}},
		{Name: "g2", Tags: []string{"tag3"}},
	}

	tfList, diags := GroupsListFromAPI(original)
	require.False(t, diags.HasError())

	roundTripped, err := GroupsListToInternal(ctx, tfList)
	require.NoError(t, err)

	require.Len(t, roundTripped, len(original))
	for i := range original {
		assert.Equal(t, original[i].Name, roundTripped[i].Name)
		assert.Equal(t, original[i].Tags, roundTripped[i].Tags)
	}
}
