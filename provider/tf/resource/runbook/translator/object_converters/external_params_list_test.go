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

// --- ExternalParamsListToInternal ---

func TestExternalParamsListToInternal_Success(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	obj, _ := types.ObjectValue(ExternalParamsListAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("ep1"),
		"value":       types.StringValue("default_val"),
		"source":      types.StringValue("alertmanager"),
		"json_path":   types.StringValue("$.data.token"),
		"export":      types.BoolValue(true),
		"description": types.StringValue("auth token"),
	})
	tfList, _ := types.ListValue(ExternalParamsListObjectType, []attr.Value{obj})

	result, err := ExternalParamsListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "ep1", result[0].Name)
	assert.Equal(t, "default_val", result[0].Value)
	assert.Equal(t, "alertmanager", result[0].Source)
	assert.Equal(t, "$.data.token", result[0].JsonPath)
	assert.True(t, result[0].Export)
	assert.Equal(t, customattribute.DefaultExternalParamType, result[0].ParamType)
	assert.Equal(t, "auth token", result[0].Description)
}

func TestExternalParamsListToInternal_MultipleParams(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	obj1, _ := types.ObjectValue(ExternalParamsListAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("ep1"),
		"value":       types.StringValue(""),
		"source":      types.StringValue("alertmanager"),
		"json_path":   types.StringValue("$.alerts"),
		"export":      types.BoolValue(false),
		"description": types.StringValue(""),
	})
	obj2, _ := types.ObjectValue(ExternalParamsListAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("ep2"),
		"value":       types.StringValue("fallback"),
		"source":      types.StringValue("alertmanager"),
		"json_path":   types.StringValue("$.status"),
		"export":      types.BoolValue(true),
		"description": types.StringValue("status param"),
	})
	tfList, _ := types.ListValue(ExternalParamsListObjectType, []attr.Value{obj1, obj2})

	result, err := ExternalParamsListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "ep1", result[0].Name)
	assert.Equal(t, "$.alerts", result[0].JsonPath)
	assert.Equal(t, "ep2", result[1].Name)
	assert.True(t, result[1].Export)
}

func TestExternalParamsListToInternal_EmptyList(t *testing.T) {
	t.Parallel()
	tfList, _ := types.ListValue(ExternalParamsListObjectType, []attr.Value{})

	result, err := ExternalParamsListToInternal(context.Background(), tfList)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestExternalParamsListToInternal_NullList(t *testing.T) {
	t.Parallel()
	result, err := ExternalParamsListToInternal(context.Background(), types.ListNull(ExternalParamsListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExternalParamsListToInternal_UnknownList(t *testing.T) {
	t.Parallel()
	result, err := ExternalParamsListToInternal(context.Background(), types.ListUnknown(ExternalParamsListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExternalParamsListToInternal_SetsDefaultParamType(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	obj, _ := types.ObjectValue(ExternalParamsListAttrTypes, map[string]attr.Value{
		"name":        types.StringValue("ep"),
		"value":       types.StringValue(""),
		"source":      types.StringValue("alertmanager"),
		"json_path":   types.StringValue("$.x"),
		"export":      types.BoolValue(false),
		"description": types.StringValue(""),
	})
	tfList, _ := types.ListValue(ExternalParamsListObjectType, []attr.Value{obj})

	result, err := ExternalParamsListToInternal(ctx, tfList)

	require.NoError(t, err)
	assert.Equal(t, "EXTERNAL", result[0].ParamType)
}

// --- ExternalParamsListFromAPI ---

func TestExternalParamsListFromAPI_Success(t *testing.T) {
	t.Parallel()
	apiParams := []customattribute.ExternalParamJson{
		{Name: "ep1", Value: "v", Source: "alertmanager", JsonPath: "$.data", Export: true, Description: "desc"},
	}

	result, diags := ExternalParamsListFromAPI(apiParams)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	require.Equal(t, 1, len(result.Elements()))

	obj := result.Elements()[0].(types.Object)
	assert.Equal(t, "ep1", obj.Attributes()["name"].(types.String).ValueString())
	assert.Equal(t, "v", obj.Attributes()["value"].(types.String).ValueString())
	assert.Equal(t, "alertmanager", obj.Attributes()["source"].(types.String).ValueString())
	assert.Equal(t, "$.data", obj.Attributes()["json_path"].(types.String).ValueString())
	assert.True(t, obj.Attributes()["export"].(types.Bool).ValueBool())
	assert.Equal(t, "desc", obj.Attributes()["description"].(types.String).ValueString())
}

func TestExternalParamsListFromAPI_MultipleParams(t *testing.T) {
	t.Parallel()
	apiParams := []customattribute.ExternalParamJson{
		{Name: "ep1", Source: "alertmanager", JsonPath: "$.a"},
		{Name: "ep2", Source: "alertmanager", JsonPath: "$.b"},
	}

	result, diags := ExternalParamsListFromAPI(apiParams)

	assert.False(t, diags.HasError())
	assert.Equal(t, 2, len(result.Elements()))
}

func TestExternalParamsListFromAPI_EmptySlice(t *testing.T) {
	t.Parallel()
	result, diags := ExternalParamsListFromAPI([]customattribute.ExternalParamJson{})

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Empty(t, result.Elements())
}

func TestExternalParamsListFromAPI_NilSlice(t *testing.T) {
	t.Parallel()
	result, diags := ExternalParamsListFromAPI(nil)

	assert.False(t, diags.HasError())
	assert.True(t, result.IsNull())
}

// --- Round-trip ---

func TestExternalParamsListRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	original := []customattribute.ExternalParamJson{
		{Name: "ep1", Value: "default", Source: "alertmanager", JsonPath: "$.alerts[0].name", Export: true, Description: "alert name"},
		{Name: "ep2", Value: "", Source: "alertmanager", JsonPath: "$.status", Export: false, Description: ""},
	}

	tfList, diags := ExternalParamsListFromAPI(original)
	require.False(t, diags.HasError())

	roundTripped, err := ExternalParamsListToInternal(ctx, tfList)
	require.NoError(t, err)
	require.Len(t, roundTripped, len(original))

	for i := range original {
		assert.Equal(t, original[i].Name, roundTripped[i].Name)
		assert.Equal(t, original[i].Value, roundTripped[i].Value)
		assert.Equal(t, original[i].Source, roundTripped[i].Source)
		assert.Equal(t, original[i].JsonPath, roundTripped[i].JsonPath)
		assert.Equal(t, original[i].Export, roundTripped[i].Export)
		assert.Equal(t, original[i].Description, roundTripped[i].Description)
	}
}

// --- Attr types ---

func TestExternalParamsListAttrTypesStructure(t *testing.T) {
	t.Parallel()
	expected := []string{"name", "value", "source", "json_path", "export", "description"}
	assert.Equal(t, len(expected), len(ExternalParamsListAttrTypes))
	for _, key := range expected {
		_, ok := ExternalParamsListAttrTypes[key]
		assert.True(t, ok, "missing key %s", key)
	}
}
