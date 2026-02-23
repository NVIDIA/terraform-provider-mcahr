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
	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JsonAttributeConfig struct {
	FullAttrName  string
	RemarshalFunc func(string, common.JsonConfig) (string, error)
	GetAttr       func(*model.RunbookTFModel) types.String
	GetFullAttr   func(*model.RunbookTFModel) types.String
	SetFullAttr   func(*model.RunbookTFModel, types.String)
}

var (
	JSON_ATTRIBUTES_TO_POPULATE = map[string]JsonAttributeConfig{
		"cells": {
			FullAttrName:  "cells_full",
			RemarshalFunc: common.RemarshalListWithConfig[*customattribute.CellJson],
			GetAttr:       func(model *model.RunbookTFModel) types.String { return model.Cells },
			GetFullAttr:   func(model *model.RunbookTFModel) types.String { return model.CellsFull },
			SetFullAttr:   func(model *model.RunbookTFModel, value types.String) { model.CellsFull = value },
		},
		"params": {
			FullAttrName:  "params_full",
			RemarshalFunc: common.RemarshalListWithConfig[*customattribute.ParamJson],
			GetAttr:       func(model *model.RunbookTFModel) types.String { return model.Params },
			GetFullAttr:   func(model *model.RunbookTFModel) types.String { return model.ParamsFull },
			SetFullAttr:   func(model *model.RunbookTFModel, value types.String) { model.ParamsFull = value },
		},
		"external_params": {
			FullAttrName:  "external_params_full",
			RemarshalFunc: common.RemarshalListWithConfig[*customattribute.ExternalParamJson],
			GetAttr:       func(model *model.RunbookTFModel) types.String { return model.ExternalParams },
			GetFullAttr:   func(model *model.RunbookTFModel) types.String { return model.ExternalParamsFull },
			SetFullAttr:   func(model *model.RunbookTFModel, value types.String) { model.ExternalParamsFull = value },
		},
	}
)

// PopulateFullJsonAttributes normalizes JSON attributes from user input and populates the corresponding _full fields.
// It applies version-aware defaults and struct tag rules (min_version, max_version) during normalization.
func PopulateFullJsonAttributes(ctx context.Context, resultValues, plan, state *model.RunbookTFModel, backendVersion *version.BackendVersion) error {

	for _, attrConfig := range JSON_ATTRIBUTES_TO_POPULATE {

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

func isDeleteOperation(plan *model.RunbookTFModel, attrConfig JsonAttributeConfig) bool {
	fullPlanValue := attrConfig.GetFullAttr(plan)
	if fullPlanValue.IsNull() {
		// It's a delete operation, do nothing.
		return true
	}
	return false
}
