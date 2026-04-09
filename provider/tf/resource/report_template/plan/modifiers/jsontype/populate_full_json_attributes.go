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
	"terraform/terraform-provider/provider/tf/resource/report_template/model"

	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JsonAttributeConfig struct {
	FullAttrName  string
	RemarshalFunc func(string, common.JsonConfig) (string, error)
	GetAttr       func(*model.ReportTemplateTFModel) types.String
	SetAttr       func(*model.ReportTemplateTFModel, types.String)
	GetFullAttr   func(*model.ReportTemplateTFModel) types.String
	SetFullAttr   func(*model.ReportTemplateTFModel, types.String)

	// GetReplacementAttr returns the replacement field value (e.g. blocks_list for blocks).
	// Leave nil for attributes that have no replacement yet.
	GetReplacementAttr func(*model.ReportTemplateTFModel) attr.Value
}

var (
	JSON_ATTRIBUTES_TO_POPULATE = map[string]JsonAttributeConfig{
		"blocks": {
			FullAttrName:       "blocks_full",
			RemarshalFunc:      common.RemarshalListWithConfig[*customattribute.BlockJson],
			GetAttr:            func(m *model.ReportTemplateTFModel) types.String { return m.Blocks },
			SetAttr:            func(m *model.ReportTemplateTFModel, v types.String) { m.Blocks = v },
			GetFullAttr:        func(m *model.ReportTemplateTFModel) types.String { return m.BlocksFull },
			SetFullAttr:        func(m *model.ReportTemplateTFModel, v types.String) { m.BlocksFull = v },
			GetReplacementAttr: func(m *model.ReportTemplateTFModel) attr.Value { return m.BlocksList },
		},
		"links": {
			FullAttrName:       "links_full",
			RemarshalFunc:      common.RemarshalListWithConfig[*customattribute.LinkJson],
			GetAttr:            func(m *model.ReportTemplateTFModel) types.String { return m.Links },
			SetAttr:            func(m *model.ReportTemplateTFModel, v types.String) { m.Links = v },
			GetFullAttr:        func(m *model.ReportTemplateTFModel) types.String { return m.LinksFull },
			SetFullAttr:        func(m *model.ReportTemplateTFModel, v types.String) { m.LinksFull = v },
			GetReplacementAttr: func(m *model.ReportTemplateTFModel) attr.Value { return m.LinksList },
		},
	}
)

// PopulateFullJsonAttributes normalizes JSON attributes from user input and populates the corresponding _full fields.
// It applies version-aware defaults and struct tag rules (min_version, max_version) during normalization.
//
// When a deprecated JSON attribute has a replacement (e.g. blocks → blocks_list), and the replacement
// is active while the base field was not explicitly set, both the base and _full fields are nulled
// and the attribute is skipped. The resultValuesWithoutDefaults parameter is used to distinguish
// "user explicitly set the field" from "field has a plan default".
func PopulateFullJsonAttributes(ctx context.Context, resultValues, resultValuesWithoutDefaults, plan, state *model.ReportTemplateTFModel, backendVersion *version.BackendVersion) error {

	for _, attrConfig := range JSON_ATTRIBUTES_TO_POPULATE {

		// Skip deprecated JSON fields when their replacement list field is active
		if shouldSkipForReplacement(attrConfig, resultValues, resultValuesWithoutDefaults, plan) {
			continue
		}

		if isDeleteOperation(plan, attrConfig) {
			continue
		}

		configValue := attrConfig.GetAttr(resultValues)

		// Remarshal the json field to apply the custom struct tags (like min_version, max_version, etc.)
		// and set the default values for the fields that are not present in the JSON
		// See customattribute structs for more details
		normalizedValue, err := attrConfig.RemarshalFunc(configValue.ValueString(), common.JsonConfig{BackendVersion: backendVersion})
		if err != nil {
			return fmt.Errorf("error populating full JSON attributes: %s", err.Error())
		}

		attrConfig.SetFullAttr(resultValues, types.StringValue(normalizedValue))
	}

	return nil
}

// shouldSkipForReplacement checks if a deprecated JSON attribute should be skipped because its
// replacement field is active. When skipping, nulls the base and _full fields in both
// resultValues and planValues (plan needs _full nulled so isDeleteOperation works on next call).
func shouldSkipForReplacement(attrConfig JsonAttributeConfig, resultValues, resultValuesWithoutDefaults, plan *model.ReportTemplateTFModel) bool {
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

func isDeleteOperation(plan *model.ReportTemplateTFModel, attrConfig JsonAttributeConfig) bool {
	fullPlanValue := attrConfig.GetFullAttr(plan)
	if fullPlanValue.IsNull() {
		// It's a delete operation, do nothing.
		return true
	}
	return false
}
