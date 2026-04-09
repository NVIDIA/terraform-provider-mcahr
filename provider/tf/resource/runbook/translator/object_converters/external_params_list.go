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

// ExternalParamTFModel is the Terraform-side representation of a single external param.
// param_type is omitted because it is always "EXTERNAL" and is a backend-internal detail.
type ExternalParamTFModel struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Source      types.String `tfsdk:"source"`
	JsonPath    types.String `tfsdk:"json_path"`
	Export      types.Bool   `tfsdk:"export"`
	Description types.String `tfsdk:"description"`
}

// ExternalParamsListAttrTypes defines the attribute type map for the nested external param object.
var ExternalParamsListAttrTypes = map[string]attr.Type{
	"name":        types.StringType,
	"value":       types.StringType,
	"source":      types.StringType,
	"json_path":   types.StringType,
	"export":      types.BoolType,
	"description": types.StringType,
}

// ExternalParamsListObjectType is the types.ObjectType for a single external param element.
var ExternalParamsListObjectType = types.ObjectType{AttrTypes: ExternalParamsListAttrTypes}

// ExternalParamsListToInternal converts a Terraform external_params_list (types.List) into []ExternalParamJson.
// Returns nil, nil when the list is null or unknown.
func ExternalParamsListToInternal(ctx context.Context, tfList types.List) ([]customattribute.ExternalParamJson, error) {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil, nil
	}

	var models []ExternalParamTFModel
	diags := tfList.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract external_params from list: %s", diags.Errors())
	}

	params := make([]customattribute.ExternalParamJson, len(models))
	for i, m := range models {
		params[i] = customattribute.ExternalParamJson{
			Name:        m.Name.ValueString(),
			Value:       m.Value.ValueString(),
			Source:      m.Source.ValueString(),
			JsonPath:    m.JsonPath.ValueString(),
			Export:      m.Export.ValueBool(),
			ParamType:   customattribute.DefaultExternalParamType,
			Description: m.Description.ValueString(),
		}
	}
	return params, nil
}

// ExternalParamsListFromAPI converts API response external params into a Terraform types.List.
// Returns a null list when apiParams is nil.
func ExternalParamsListFromAPI(apiParams []customattribute.ExternalParamJson) (types.List, diag.Diagnostics) {
	if apiParams == nil {
		return types.ListNull(ExternalParamsListObjectType), nil
	}

	objects := make([]attr.Value, len(apiParams))
	for i, p := range apiParams {
		obj, diags := types.ObjectValue(ExternalParamsListAttrTypes, map[string]attr.Value{
			"name":        types.StringValue(p.Name),
			"value":       types.StringValue(p.Value),
			"source":      types.StringValue(p.Source),
			"json_path":   types.StringValue(p.JsonPath),
			"export":      types.BoolValue(p.Export),
			"description": types.StringValue(p.Description),
		})
		if diags.HasError() {
			return types.ListNull(ExternalParamsListObjectType), diags
		}
		objects[i] = obj
	}

	result, diags := types.ListValue(ExternalParamsListObjectType, objects)
	if diags.HasError() {
		return types.ListNull(ExternalParamsListObjectType), diags
	}
	return result, nil
}
