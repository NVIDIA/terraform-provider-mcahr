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

package translator

import (
	"encoding/json"
	"fmt"

	"terraform/terraform-provider/provider/common"
	dashboardapi "terraform/terraform-provider/provider/external_api/resources/dashboards"
	customattribute "terraform/terraform-provider/provider/external_api/resources/dashboards/custom_attribute"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	dashboardtf "terraform/terraform-provider/provider/tf/resource/dashboard/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DashboardTranslatorV1 handles translation between TF models and V1 API models for dashboard resources
type DashboardTranslatorV1 struct {
	DashboardTranslatorCommon
}

var _ translator.Translator[*dashboardtf.DashboardTFModel, *dashboardapi.DashboardResponseAPIModelV1] = &DashboardTranslatorV1{}

// ToTFModel converts a V1 API model to a TF model
func (t *DashboardTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *dashboardapi.DashboardResponseAPIModelV1) (*dashboardtf.DashboardTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the dashboard container
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no dashboard container found in V1 API response")
	}

	if len(container.DashboardClasses) == 0 {
		return nil, fmt.Errorf("no dashboard classes found in V1 API response")
	}

	// Get the first dashboard class
	dashboardClass := container.DashboardClasses[0]

	// Parse the configuration JSON string
	var config customattribute.ConfigurationJson
	if err := json.Unmarshal([]byte(dashboardClass.Configuration), &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration JSON: %w", err)
	}

	// Convert groups to JSON string
	groupsJSON, err := json.Marshal(config.Groups)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal groups: %w", err)
	}
	groupsValue := types.StringValue(string(groupsJSON))

	// Convert values to JSON string
	valuesJSON, err := json.Marshal(config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal values: %w", err)
	}
	valuesValue := types.StringValue(string(valuesJSON))

	// Convert other_tags and identifiers to sets
	otherTags := translator.ListValueFromStringSlice(requestContext.Context, config.OtherTags)
	identifiers := translator.ListValueFromStringSlice(requestContext.Context, config.Identifiers)

	tfModel := &dashboardtf.DashboardTFModel{
		Name:          types.StringValue(dashboardClass.Name),
		DashboardType: types.StringValue(dashboardClass.DashboardType),
		ResourceQuery: types.StringValue(config.ResourceQuery),
		Groups:        groupsValue,
		GroupsFull:    groupsValue,
		Values:        valuesValue,
		ValuesFull:    valuesValue,
		OtherTags:     otherTags,
		Identifiers:   identifiers,
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (t *DashboardTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *dashboardtf.DashboardTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
