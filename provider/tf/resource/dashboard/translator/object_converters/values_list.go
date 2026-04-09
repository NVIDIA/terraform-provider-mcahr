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
	customattribute "terraform/terraform-provider/provider/external_api/resources/dashboards/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ValueTFModel is the Terraform-side representation of a single value in values_list.
type ValueTFModel struct {
	Color  types.String `tfsdk:"color"`
	Values types.List   `tfsdk:"values"`
}

var ValuesListAttrTypes = map[string]attr.Type{
	"color":  types.StringType,
	"values": types.ListType{ElemType: types.StringType},
}

var ValuesListObjectType = types.ObjectType{AttrTypes: ValuesListAttrTypes}

// ValuesListToInternal converts a Terraform values_list (types.List) into []ValueJson.
func ValuesListToInternal(ctx context.Context, tfList types.List) ([]customattribute.ValueJson, error) {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil, nil
	}

	var models []ValueTFModel
	diags := tfList.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract values from list: %s", diags.Errors())
	}

	values := make([]customattribute.ValueJson, len(models))
	for i, m := range models {
		var vals []string
		m.Values.ElementsAs(ctx, &vals, false)
		values[i] = customattribute.ValueJson{
			Color:  m.Color.ValueString(),
			Values: vals,
		}
	}
	return values, nil
}

// ValuesListFromAPI converts API response values into a Terraform types.List.
func ValuesListFromAPI(apiValues []customattribute.ValueJson) (types.List, diag.Diagnostics) {
	if apiValues == nil {
		return types.ListNull(ValuesListObjectType), nil
	}

	objects := make([]attr.Value, len(apiValues))
	for i, v := range apiValues {
		valsVal, diags := types.ListValueFrom(context.Background(), types.StringType, v.Values)
		if diags.HasError() {
			return types.ListNull(ValuesListObjectType), diags
		}

		obj, diags := types.ObjectValue(ValuesListAttrTypes, map[string]attr.Value{
			"color":  types.StringValue(v.Color),
			"values": valsVal,
		})
		if diags.HasError() {
			return types.ListNull(ValuesListObjectType), diags
		}
		objects[i] = obj
	}

	result, diags := types.ListValue(ValuesListObjectType, objects)
	if diags.HasError() {
		return types.ListNull(ValuesListObjectType), diags
	}
	return result, nil
}
