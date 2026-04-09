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

// GroupTFModel is the Terraform-side representation of a single group in groups_list.
type GroupTFModel struct {
	Name types.String `tfsdk:"name"`
	Tags types.List   `tfsdk:"tags"`
}

var GroupsListAttrTypes = map[string]attr.Type{
	"name": types.StringType,
	"tags": types.ListType{ElemType: types.StringType},
}

var GroupsListObjectType = types.ObjectType{AttrTypes: GroupsListAttrTypes}

// GroupsListToInternal converts a Terraform groups_list (types.List) into []GroupJson.
func GroupsListToInternal(ctx context.Context, tfList types.List) ([]customattribute.GroupJson, error) {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil, nil
	}

	var models []GroupTFModel
	diags := tfList.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract groups from list: %s", diags.Errors())
	}

	groups := make([]customattribute.GroupJson, len(models))
	for i, m := range models {
		var tags []string
		m.Tags.ElementsAs(ctx, &tags, false)
		groups[i] = customattribute.GroupJson{
			Name: m.Name.ValueString(),
			Tags: tags,
		}
	}
	return groups, nil
}

// GroupsListFromAPI converts API response groups into a Terraform types.List.
func GroupsListFromAPI(apiGroups []customattribute.GroupJson) (types.List, diag.Diagnostics) {
	if apiGroups == nil {
		return types.ListNull(GroupsListObjectType), nil
	}

	objects := make([]attr.Value, len(apiGroups))
	for i, g := range apiGroups {
		tagsVal, diags := types.ListValueFrom(context.Background(), types.StringType, g.Tags)
		if diags.HasError() {
			return types.ListNull(GroupsListObjectType), diags
		}

		obj, diags := types.ObjectValue(GroupsListAttrTypes, map[string]attr.Value{
			"name": types.StringValue(g.Name),
			"tags": tagsVal,
		})
		if diags.HasError() {
			return types.ListNull(GroupsListObjectType), diags
		}
		objects[i] = obj
	}

	result, diags := types.ListValue(GroupsListObjectType, objects)
	if diags.HasError() {
		return types.ListNull(GroupsListObjectType), diags
	}
	return result, nil
}
