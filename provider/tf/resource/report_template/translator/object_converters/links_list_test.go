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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinksListToInternal_Success(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	obj, _ := types.ObjectValue(LinksListAttrTypes, map[string]attr.Value{
		"label":                types.StringValue("View Details"),
		"report_template_name": types.StringValue("detail_template"),
	})
	tfList, _ := types.ListValue(LinksListObjectType, []attr.Value{obj})

	result, err := LinksListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "View Details", result[0].Label)
	assert.Equal(t, "detail_template", result[0].ReportTemplateName)
}

func TestLinksListToInternal_MultipleLinks(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	obj1, _ := types.ObjectValue(LinksListAttrTypes, map[string]attr.Value{
		"label":                types.StringValue("Link A"),
		"report_template_name": types.StringValue("template_a"),
	})
	obj2, _ := types.ObjectValue(LinksListAttrTypes, map[string]attr.Value{
		"label":                types.StringValue("Link B"),
		"report_template_name": types.StringValue("template_b"),
	})
	tfList, _ := types.ListValue(LinksListObjectType, []attr.Value{obj1, obj2})

	result, err := LinksListToInternal(ctx, tfList)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Link A", result[0].Label)
	assert.Equal(t, "template_a", result[0].ReportTemplateName)
	assert.Equal(t, "Link B", result[1].Label)
	assert.Equal(t, "template_b", result[1].ReportTemplateName)
}

func TestLinksListToInternal_EmptyList(t *testing.T) {
	t.Parallel()
	tfList, _ := types.ListValue(LinksListObjectType, []attr.Value{})

	result, err := LinksListToInternal(context.Background(), tfList)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestLinksListToInternal_NullList(t *testing.T) {
	t.Parallel()
	result, err := LinksListToInternal(context.Background(), types.ListNull(LinksListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestLinksListToInternal_UnknownList(t *testing.T) {
	t.Parallel()
	result, err := LinksListToInternal(context.Background(), types.ListUnknown(LinksListObjectType))

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestLinksListFromAPI_Success(t *testing.T) {
	t.Parallel()
	apiLinks := []customattribute.LinkJson{
		{Label: "View", ReportTemplateName: "other_template"},
	}

	result, diags := LinksListFromAPI(apiLinks)

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	require.Equal(t, 1, len(result.Elements()))

	obj := result.Elements()[0].(types.Object)
	assert.Equal(t, "View", obj.Attributes()["label"].(types.String).ValueString())
	assert.Equal(t, "other_template", obj.Attributes()["report_template_name"].(types.String).ValueString())
}

func TestLinksListFromAPI_MultipleLinks(t *testing.T) {
	t.Parallel()
	apiLinks := []customattribute.LinkJson{
		{Label: "A", ReportTemplateName: "tmpl_a"},
		{Label: "B", ReportTemplateName: "tmpl_b"},
	}

	result, diags := LinksListFromAPI(apiLinks)

	assert.False(t, diags.HasError())
	require.Equal(t, 2, len(result.Elements()))

	obj0 := result.Elements()[0].(types.Object)
	assert.Equal(t, "A", obj0.Attributes()["label"].(types.String).ValueString())
	assert.Equal(t, "tmpl_a", obj0.Attributes()["report_template_name"].(types.String).ValueString())

	obj1 := result.Elements()[1].(types.Object)
	assert.Equal(t, "B", obj1.Attributes()["label"].(types.String).ValueString())
	assert.Equal(t, "tmpl_b", obj1.Attributes()["report_template_name"].(types.String).ValueString())
}

func TestLinksListFromAPI_EmptySlice(t *testing.T) {
	t.Parallel()
	result, diags := LinksListFromAPI([]customattribute.LinkJson{})

	assert.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Empty(t, result.Elements())
}

func TestLinksListFromAPI_NilSlice(t *testing.T) {
	t.Parallel()
	result, diags := LinksListFromAPI(nil)

	assert.False(t, diags.HasError())
	assert.True(t, result.IsNull())
}

func TestLinksListRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	original := []customattribute.LinkJson{
		{Label: "Dashboard", ReportTemplateName: "dashboard_report"},
		{Label: "Summary", ReportTemplateName: "summary_report"},
	}

	tfList, diags := LinksListFromAPI(original)
	require.False(t, diags.HasError())

	roundTripped, err := LinksListToInternal(ctx, tfList)
	require.NoError(t, err)
	require.Len(t, roundTripped, len(original))

	for i := range original {
		assert.Equal(t, original[i].Label, roundTripped[i].Label)
		assert.Equal(t, original[i].ReportTemplateName, roundTripped[i].ReportTemplateName)
	}
}

func TestLinksListAttrTypesStructure(t *testing.T) {
	t.Parallel()
	expected := []string{"label", "report_template_name"}
	assert.Equal(t, len(expected), len(LinksListAttrTypes))
	for _, key := range expected {
		_, ok := LinksListAttrTypes[key]
		assert.True(t, ok, "missing key %s", key)
	}
}
