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

func TestValuesListToInternal_Success(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	valsVal, _ := types.ListValueFrom(ctx, types.StringType, []string{"aws", "gcp"})
	obj, _ := types.ObjectValue(ValuesListAttrTypes, map[string]attr.Value{
		"color":  types.StringValue("#78909c"),
		"values": valsVal,
	})
	tfList, _ := types.ListValue(ValuesListObjectType, []attr.Value{obj})

	result, err := ValuesListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "#78909c", result[0].Color)
	assert.Equal(t, []string{"aws", "gcp"}, result[0].Values)
}

func TestValuesListToInternal_MultipleEntries(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	vals1, _ := types.ListValueFrom(ctx, types.StringType, []string{"aws"})
	obj1, _ := types.ObjectValue(ValuesListAttrTypes, map[string]attr.Value{
		"color":  types.StringValue("#78909c"),
		"values": vals1,
	})

	vals2, _ := types.ListValueFrom(ctx, types.StringType, []string{"release-X"})
	obj2, _ := types.ObjectValue(ValuesListAttrTypes, map[string]attr.Value{
		"color":  types.StringValue("#ffa726"),
		"values": vals2,
	})

	tfList, _ := types.ListValue(ValuesListObjectType, []attr.Value{obj1, obj2})

	result, err := ValuesListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "#78909c", result[0].Color)
	assert.Equal(t, []string{"aws"}, result[0].Values)
	assert.Equal(t, "#ffa726", result[1].Color)
	assert.Equal(t, []string{"release-X"}, result[1].Values)
}

func TestValuesListToInternal_EmptyList(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tfList, _ := types.ListValue(ValuesListObjectType, []attr.Value{})

	result, err := ValuesListToInternal(ctx, tfList)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestValuesListToInternal_NullList(t *testing.T) {
	t.Parallel()

	result, err := ValuesListToInternal(context.Background(), types.ListNull(ValuesListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestValuesListToInternal_UnknownList(t *testing.T) {
	t.Parallel()

	result, err := ValuesListToInternal(context.Background(), types.ListUnknown(ValuesListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestValuesListToInternal_EmptyValues(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	emptyVals, _ := types.ListValueFrom(ctx, types.StringType, []string{})
	obj, _ := types.ObjectValue(ValuesListAttrTypes, map[string]attr.Value{
		"color":  types.StringValue("#000000"),
		"values": emptyVals,
	})
	tfList, _ := types.ListValue(ValuesListObjectType, []attr.Value{obj})

	result, err := ValuesListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "#000000", result[0].Color)
	assert.Empty(t, result[0].Values)
}

func TestValuesListFromAPI_Success(t *testing.T) {
	t.Parallel()

	apiValues := []customattribute.ValueJson{
		{Color: "#78909c", Values: []string{"aws", "gcp"}},
	}

	result, diags := ValuesListFromAPI(apiValues)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Equal(t, 1, len(result.Elements()))
}

func TestValuesListFromAPI_MultipleEntries(t *testing.T) {
	t.Parallel()

	apiValues := []customattribute.ValueJson{
		{Color: "#78909c", Values: []string{"aws"}},
		{Color: "#ffa726", Values: []string{"release-X"}},
	}

	result, diags := ValuesListFromAPI(apiValues)

	assert.False(t, diags.HasError())
	assert.Equal(t, 2, len(result.Elements()))
}

func TestValuesListFromAPI_EmptyValues(t *testing.T) {
	t.Parallel()

	result, diags := ValuesListFromAPI([]customattribute.ValueJson{})

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Empty(t, result.Elements())
}

func TestValuesListFromAPI_NilValues(t *testing.T) {
	t.Parallel()

	result, diags := ValuesListFromAPI(nil)

	assert.False(t, diags.HasError())
	assert.True(t, result.IsNull())
}

func TestValuesListFromAPI_EmptyValuesField(t *testing.T) {
	t.Parallel()

	apiValues := []customattribute.ValueJson{
		{Color: "#78909c", Values: []string{}},
	}

	result, diags := ValuesListFromAPI(apiValues)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Equal(t, 1, len(result.Elements()))
}

func TestValuesListRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	original := []customattribute.ValueJson{
		{Color: "#78909c", Values: []string{"aws", "gcp"}},
		{Color: "#ffa726", Values: []string{"release-X"}},
	}

	tfList, diags := ValuesListFromAPI(original)
	require.False(t, diags.HasError())

	roundTripped, err := ValuesListToInternal(ctx, tfList)
	require.NoError(t, err)

	require.Len(t, roundTripped, len(original))
	for i := range original {
		assert.Equal(t, original[i].Color, roundTripped[i].Color)
		assert.Equal(t, original[i].Values, roundTripped[i].Values)
	}
}
