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

	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- ParamsListToInternal ---

func TestParamsListToInternal_Success(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	obj, _ := types.ObjectValue(ParamsListAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("param1"),
		"value":       types.StringValue("val1"),
		"required":    types.BoolValue(true),
		"export":      types.BoolValue(false),
		"description": types.StringValue("first param"),
	})
	tfList, _ := types.ListValue(ParamsListObjectType, []attr.Value{obj})

	result, err := ParamsListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "param1", result[0].Name)
	assert.Equal(t, "val1", result[0].Value)
	assert.True(t, result[0].Required)
	assert.False(t, result[0].Export)
	assert.Equal(t, customattribute.DefaultParamType, result[0].ParamType)
	assert.Equal(t, "first param", result[0].Description)
}

func TestParamsListToInternal_MultipleParams(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	obj1, _ := types.ObjectValue(ParamsListAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("p1"),
		"value":       types.StringValue("v1"),
		"required":    types.BoolValue(true),
		"export":      types.BoolValue(true),
		"description": types.StringValue(""),
	})
	obj2, _ := types.ObjectValue(ParamsListAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("p2"),
		"value":       types.StringValue("v2"),
		"required":    types.BoolValue(false),
		"export":      types.BoolValue(false),
		"description": types.StringValue("second"),
	})
	tfList, _ := types.ListValue(ParamsListObjectType, []attr.Value{obj1, obj2})

	result, err := ParamsListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "p1", result[0].Name)
	assert.True(t, result[0].Export)
	assert.Equal(t, "p2", result[1].Name)
	assert.Equal(t, "second", result[1].Description)
}

func TestParamsListToInternal_EmptyList(t *testing.T) {
	t.Parallel()
	tfList, _ := types.ListValue(ParamsListObjectType, []attr.Value{})

	result, err := ParamsListToInternal(context.Background(), tfList)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestParamsListToInternal_NullList(t *testing.T) {
	t.Parallel()
	result, err := ParamsListToInternal(context.Background(), types.ListNull(ParamsListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParamsListToInternal_UnknownList(t *testing.T) {
	t.Parallel()
	result, err := ParamsListToInternal(context.Background(), types.ListUnknown(ParamsListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParamsListToInternal_SetsDefaultParamType(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	obj, _ := types.ObjectValue(ParamsListAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("p"),
		"value":       types.StringValue(""),
		"required":    types.BoolValue(false),
		"export":      types.BoolValue(false),
		"description": types.StringValue(""),
	})
	tfList, _ := types.ListValue(ParamsListObjectType, []attr.Value{obj})

	result, err := ParamsListToInternal(ctx, tfList)

	require.NoError(t, err)
	assert.Equal(t, "PARAM", result[0].ParamType)
}

// --- ParamsListFromAPI ---

func TestParamsListFromAPI_Success(t *testing.T) {
	t.Parallel()
	apiParams := []customattribute.ParamJson{
		{Name: "p1", Value: "v1", Required: true, Export: true, Description: "desc"},
	}

	result, diags := ParamsListFromAPI(apiParams)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	require.Equal(t, 1, len(result.Elements()))

	obj := result.Elements()[0].(types.Object)
	assert.Equal(t, "p1", obj.Attributes()["name"].(types.String).ValueString())
	assert.Equal(t, "v1", obj.Attributes()["value"].(types.String).ValueString())
	assert.True(t, obj.Attributes()["required"].(types.Bool).ValueBool())
	assert.True(t, obj.Attributes()["export"].(types.Bool).ValueBool())
	assert.Equal(t, "desc", obj.Attributes()["description"].(types.String).ValueString())
}

func TestParamsListFromAPI_MultipleParams(t *testing.T) {
	t.Parallel()
	apiParams := []customattribute.ParamJson{
		{Name: "p1", Value: "v1"},
		{Name: "p2", Value: "v2", Required: true},
	}

	result, diags := ParamsListFromAPI(apiParams)

	assert.False(t, diags.HasError())
	assert.Equal(t, 2, len(result.Elements()))
}

func TestParamsListFromAPI_EmptySlice(t *testing.T) {
	t.Parallel()
	result, diags := ParamsListFromAPI([]customattribute.ParamJson{})

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Empty(t, result.Elements())
}

func TestParamsListFromAPI_NilSlice(t *testing.T) {
	t.Parallel()
	result, diags := ParamsListFromAPI(nil)

	assert.False(t, diags.HasError())
	assert.True(t, result.IsNull())
}

// --- Round-trip ---

func TestParamsListRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	original := []customattribute.ParamJson{
		{Name: "api_key", Value: "secret", Required: true, Export: true, Description: "API key"},
		{Name: "timeout", Value: "30", Required: false, Export: false, Description: ""},
	}

	tfList, diags := ParamsListFromAPI(original)
	require.False(t, diags.HasError())

	roundTripped, err := ParamsListToInternal(ctx, tfList)
	require.NoError(t, err)
	require.Len(t, roundTripped, len(original))

	for i := range original {
		assert.Equal(t, original[i].Name, roundTripped[i].Name)
		assert.Equal(t, original[i].Value, roundTripped[i].Value)
		assert.Equal(t, original[i].Required, roundTripped[i].Required)
		assert.Equal(t, original[i].Export, roundTripped[i].Export)
		assert.Equal(t, original[i].Description, roundTripped[i].Description)
	}
}

// --- Attr types ---

func TestParamsListAttrTypesStructure(t *testing.T) {
	t.Parallel()
	expected := []string{"name", "value", "required", "export", "description"}
	assert.Equal(t, len(expected), len(ParamsListAttrTypes))
	for _, key := range expected {
		_, ok := ParamsListAttrTypes[key]
		assert.True(t, ok, "missing key %s", key)
	}
}
