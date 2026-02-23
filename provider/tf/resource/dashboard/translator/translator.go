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
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	dashboardtf "terraform/terraform-provider/provider/tf/resource/dashboard/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DashboardTranslator handles V2 API translation for dashboard resources
type DashboardTranslator struct {
	DashboardTranslatorCommon
}

var _ translator.Translator[*dashboardtf.DashboardTFModel, *dashboardapi.DashboardResponseAPIModel] = &DashboardTranslator{}

// NewDashboardTranslator creates a new V2 dashboard translator
func NewDashboardTranslator() *DashboardTranslator {
	return &DashboardTranslator{}
}

// ToTFModel converts API response to TF model using container approach
func (t *DashboardTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, response *dashboardapi.DashboardResponseAPIModel) (*dashboardtf.DashboardTFModel, error) {
	if response == nil {
		return nil, nil
	}

	// Extract configurations from response
	if len(response.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no configurations found in API response")
	}

	// Get the first configuration
	configItem := response.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	// Convert groups to JSON string
	groupsJSON, err := json.Marshal(config.Configuration.Groups)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal groups: %w", err)
	}
	groupsValue := types.StringValue(string(groupsJSON))

	// Convert values to JSON string
	valuesJSON, err := json.Marshal(config.Configuration.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal values: %w", err)
	}
	valuesValue := types.StringValue(string(valuesJSON))

	// Convert other_tags and identifiers to sets
	otherTags := translator.ListValueFromStringSlice(requestContext.Context, config.Configuration.OtherTags)
	identifiers := translator.ListValueFromStringSlice(requestContext.Context, config.Configuration.Identifiers)

	tfModel := &dashboardtf.DashboardTFModel{
		Name:          types.StringValue(metadata.Name),
		DashboardType: types.StringValue(config.DashboardType),
		ResourceQuery: types.StringValue(config.Configuration.ResourceQuery),
		Groups:        groupsValue,
		GroupsFull:    groupsValue,
		Values:        valuesValue,
		ValuesFull:    valuesValue,
		OtherTags:     otherTags,
		Identifiers:   identifiers,
	}

	return tfModel, nil
}

// ToAPIModel converts the TF model to API model for create operations
func (t *DashboardTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *dashboardtf.DashboardTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
