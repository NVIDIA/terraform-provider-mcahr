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
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ParamTFModel is the Terraform-side representation of a single param in params_list.
// param_type is omitted because it is always "PARAM" and is a backend-internal detail.
type ParamTFModel struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Required    types.Bool   `tfsdk:"required"`
	Export      types.Bool   `tfsdk:"export"`
	Description types.String `tfsdk:"description"`
}

// ParamsListAttrTypes defines the attribute type map for the nested param object.
var ParamsListAttrTypes = map[string]attr.Type{
	"name":        types.StringType,
	"value":       types.StringType,
	"required":    types.BoolType,
	"export":      types.BoolType,
	"description": types.StringType,
}

// ParamsListObjectType is the types.ObjectType for a single param element in params_list.
var ParamsListObjectType = types.ObjectType{AttrTypes: ParamsListAttrTypes}

// ParamsListToInternal converts a Terraform params_list (types.List) into []ParamJson.
// Returns nil, nil when the list is null or unknown.
func ParamsListToInternal(ctx context.Context, tfList types.List) ([]customattribute.ParamJson, error) {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil, nil
	}

	var models []ParamTFModel
	diags := tfList.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract params from list: %s", diags.Errors())
	}

	params := make([]customattribute.ParamJson, len(models))
	for i, m := range models {
		params[i] = customattribute.ParamJson{
			Name:        m.Name.ValueString(),
			Value:       m.Value.ValueString(),
			Required:    m.Required.ValueBool(),
			Export:      m.Export.ValueBool(),
			ParamType:   customattribute.DefaultParamType,
			Description: m.Description.ValueString(),
		}
	}
	return params, nil
}

// ParamsListFromAPI converts API response params ([]ParamJson) into a Terraform types.List.
// Returns a null list when apiParams is nil.
func ParamsListFromAPI(apiParams []customattribute.ParamJson) (types.List, diag.Diagnostics) {
	if apiParams == nil {
		return types.ListNull(ParamsListObjectType), nil
	}

	objects := make([]attr.Value, len(apiParams))
	for i, p := range apiParams {
		obj, diags := types.ObjectValue(ParamsListAttrTypes, map[string]attr.Value{
			"name":        types.StringValue(p.Name),
			"value":       types.StringValue(p.Value),
			"required":    types.BoolValue(p.Required),
			"export":      types.BoolValue(p.Export),
			"description": types.StringValue(p.Description),
		})
		if diags.HasError() {
			return types.ListNull(ParamsListObjectType), diags
		}
		objects[i] = obj
	}

	result, diags := types.ListValue(ParamsListObjectType, objects)
	if diags.HasError() {
		return types.ListNull(ParamsListObjectType), diags
	}
	return result, nil
}
