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
	"fmt"
	"terraform/terraform-provider/provider/common"
	runbookapi "terraform/terraform-provider/provider/external_api/resources/runbooks"
	utils "terraform/terraform-provider/provider/tf/core/translator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ParamsGroupsAttrTypes = map[string]attr.Type{
	"required": types.ListType{ElemType: types.StringType},
	"optional": types.ListType{ElemType: types.StringType},
	"exported": types.ListType{ElemType: types.StringType},
	"external": types.ListType{ElemType: types.StringType},
}

func ParamsGroupsFromTFModel(requestContext *common.RequestContext, tfParamsGroups types.Object) (runbookapi.ParamsGroups, error) {

	attributes := tfParamsGroups.Attributes()

	required, err := getStringListFromAttributes(requestContext, attributes, "required")
	if err != nil {
		return runbookapi.ParamsGroups{}, err
	}

	optional, err := getStringListFromAttributes(requestContext, attributes, "optional")
	if err != nil {
		return runbookapi.ParamsGroups{}, err
	}

	exported, err := getStringListFromAttributes(requestContext, attributes, "exported")
	if err != nil {
		return runbookapi.ParamsGroups{}, err
	}

	external, err := getStringListFromAttributes(requestContext, attributes, "external")
	if err != nil {
		return runbookapi.ParamsGroups{}, err
	}

	apiParamsGroups := runbookapi.ParamsGroups{
		Required: required,
		Optional: optional,
		Exported: exported,
		External: external,
	}

	return apiParamsGroups, nil
}

// getStringListFromAttributes retrieves a string list from the attributes map by key name
// Returns an empty list if the key is missing
// Returns an error if the key exists but has the wrong type
func getStringListFromAttributes(requestContext *common.RequestContext, attributes map[string]attr.Value, key string) ([]string, error) {
	attrValue, exists := attributes[key]
	if !exists {
		return []string{}, nil
	}

	list, ok := attrValue.(types.List)
	if !ok {
		return nil, fmt.Errorf("params_groups key '%s' is not a list, got type %T", key, attrValue)
	}

	return utils.ListSliceFromTFModel(requestContext.Context, list), nil
}

func ParamsGroupsToTFModel(requestContext *common.RequestContext, apiParamsGroups runbookapi.ParamsGroups) (types.Object, diag.Diagnostics) {

	tfParamsGroups, diags := types.ObjectValueFrom(requestContext.Context, ParamsGroupsAttrTypes, apiParamsGroups)
	if diags.HasError() {
		return types.ObjectNull(ParamsGroupsAttrTypes), diags
	}

	return tfParamsGroups, diags
}
