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
	"fmt"
	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LinkTFModel struct {
	Label              types.String `tfsdk:"label"`
	ReportTemplateName types.String `tfsdk:"report_template_name"`
}

var LinksListAttrTypes = map[string]attr.Type{
	"label":                types.StringType,
	"report_template_name": types.StringType,
}

var LinksListObjectType = types.ObjectType{AttrTypes: LinksListAttrTypes}

func LinksListToInternal(ctx context.Context, tfList types.List) ([]customattribute.LinkJson, error) {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil, nil
	}

	var models []LinkTFModel
	diags := tfList.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract links from list: %s", diags.Errors())
	}

	links := make([]customattribute.LinkJson, len(models))
	for i, m := range models {
		links[i] = customattribute.LinkJson{
			Label:              m.Label.ValueString(),
			ReportTemplateName: m.ReportTemplateName.ValueString(),
		}
	}
	return links, nil
}

func LinksListFromAPI(apiLinks []customattribute.LinkJson) (types.List, diag.Diagnostics) {
	if apiLinks == nil {
		return types.ListNull(LinksListObjectType), nil
	}

	objects := make([]attr.Value, len(apiLinks))
	for i, l := range apiLinks {
		obj, diags := types.ObjectValue(LinksListAttrTypes, map[string]attr.Value{
			"label":                types.StringValue(l.Label),
			"report_template_name": types.StringValue(l.ReportTemplateName),
		})
		if diags.HasError() {
			return types.ListNull(LinksListObjectType), diags
		}
		objects[i] = obj
	}

	result, diags := types.ListValue(LinksListObjectType, objects)
	if diags.HasError() {
		return types.ListNull(LinksListObjectType), diags
	}
	return result, nil
}
