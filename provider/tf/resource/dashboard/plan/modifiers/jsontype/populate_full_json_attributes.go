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

package jsonmodifier

import (
	"context"
	"fmt"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/tf/resource/dashboard/model"

	customattribute "terraform/terraform-provider/provider/external_api/resources/dashboards/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JsonAttributeConfig struct {
	FullAttrName  string
	RemarshalFunc func(string, common.JsonConfig) (string, error)
	GetAttr       func(*model.DashboardTFModel) types.String
	SetAttr       func(*model.DashboardTFModel, types.String)
	GetFullAttr   func(*model.DashboardTFModel) types.String
	SetFullAttr   func(*model.DashboardTFModel, types.String)

	// GetReplacementAttr returns the replacement field value (e.g. groups_list for groups).
	// Leave nil for attributes that have no replacement yet.
	GetReplacementAttr func(*model.DashboardTFModel) attr.Value
}

var (
	JSON_ATTRIBUTES_TO_POPULATE = map[string]JsonAttributeConfig{
		"groups": {
			FullAttrName:       "groups_full",
			RemarshalFunc:      common.RemarshalListWithConfig[*customattribute.GroupJson],
			GetAttr:            func(m *model.DashboardTFModel) types.String { return m.Groups },
			SetAttr:            func(m *model.DashboardTFModel, v types.String) { m.Groups = v },
			GetFullAttr:        func(m *model.DashboardTFModel) types.String { return m.GroupsFull },
			SetFullAttr:        func(m *model.DashboardTFModel, v types.String) { m.GroupsFull = v },
			GetReplacementAttr: func(m *model.DashboardTFModel) attr.Value { return m.GroupsList },
		},
		"values": {
			FullAttrName:       "values_full",
			RemarshalFunc:      common.RemarshalListWithConfig[*customattribute.ValueJson],
			GetAttr:            func(m *model.DashboardTFModel) types.String { return m.Values },
			SetAttr:            func(m *model.DashboardTFModel, v types.String) { m.Values = v },
			GetFullAttr:        func(m *model.DashboardTFModel) types.String { return m.ValuesFull },
			SetFullAttr:        func(m *model.DashboardTFModel, v types.String) { m.ValuesFull = v },
			GetReplacementAttr: func(m *model.DashboardTFModel) attr.Value { return m.ValuesList },
		},
	}
)

// PopulateFullJsonAttributes normalizes JSON attributes from user input and populates the corresponding _full fields.
// It applies version-aware defaults and struct tag rules (min_version, max_version) during normalization.
func PopulateFullJsonAttributes(ctx context.Context, resultValues, resultValuesWithoutDefaults, plan, state *model.DashboardTFModel, backendVersion *version.BackendVersion) error {

	for _, attrConfig := range JSON_ATTRIBUTES_TO_POPULATE {

		if shouldSkipForReplacement(attrConfig, resultValues, resultValuesWithoutDefaults, plan) {
			continue
		}

		if isDeleteOperation(plan, attrConfig) {
			continue
		}

		configValue := attrConfig.GetAttr(resultValues)

		normalizedValue, err := attrConfig.RemarshalFunc(configValue.ValueString(), common.JsonConfig{BackendVersion: backendVersion})
		if err != nil {
			return fmt.Errorf("error populating full JSON attributes: %s", err.Error())
		}

		attrConfig.SetFullAttr(resultValues, types.StringValue(normalizedValue))
	}

	return nil
}

func shouldSkipForReplacement(attrConfig JsonAttributeConfig, resultValues, resultValuesWithoutDefaults, plan *model.DashboardTFModel) bool {
	if attrConfig.GetReplacementAttr == nil {
		return false
	}

	replacementActive := common.IsAttrKnown(attrConfig.GetReplacementAttr(resultValues))
	baseExplicitlySet := common.IsAttrKnown(attrConfig.GetAttr(resultValuesWithoutDefaults))

	if replacementActive && !baseExplicitlySet {
		attrConfig.SetAttr(resultValues, types.StringNull())
		attrConfig.SetFullAttr(resultValues, types.StringNull())
		attrConfig.SetFullAttr(plan, types.StringNull())
		return true
	}

	return false
}

func isDeleteOperation(plan *model.DashboardTFModel, attrConfig JsonAttributeConfig) bool {
	fullPlanValue := attrConfig.GetFullAttr(plan)
	if fullPlanValue.IsNull() {
		return true
	}
	return false
}
