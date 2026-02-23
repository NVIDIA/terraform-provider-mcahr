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
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	"terraform/terraform-provider/provider/tf/resource/dashboard/model"
)

// No need for custom structs, using customattribute package instead

// DashboardTranslatorCommon contains shared logic for dashboard translators
type DashboardTranslatorCommon struct{}

// ToAPIModelWithVersion creates a statement API model for the given TF model and API version
func (t *DashboardTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *model.DashboardTFModel) (*statement.StatementInputAPIModel, error) {
	var stmt string

	switch requestContext.Operation {
	case common.Create:
		stmt = t.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = t.buildReadStatement(tfModel)
	case common.Update:
		stmt = t.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt = t.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}

	apiModel := &statement.StatementInputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (t *DashboardTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *model.DashboardTFModel) string {
	return t.buildDashboardStatement(requestContext, translationData, "define_dashboard", tfModel)
}

func (t *DashboardTranslatorCommon) buildReadStatement(tfModel *model.DashboardTFModel) string {
	return fmt.Sprintf("get_dashboard_class(dashboard_name=\"%s\")", tfModel.Name.ValueString())
}

func (t *DashboardTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *model.DashboardTFModel) string {
	return t.buildDashboardStatement(requestContext, translationData, "update_dashboard", tfModel)
}

func (t *DashboardTranslatorCommon) buildDeleteStatement(tfModel *model.DashboardTFModel) string {
	return fmt.Sprintf("delete_dashboard(dashboard_name=\"%s\")", tfModel.Name.ValueString())
}

func (t *DashboardTranslatorCommon) buildDashboardStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *model.DashboardTFModel) string {

	configJSON, _ := t.buildDashboardConfigurationJSON(requestContext, tfModel)

	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("dashboard_name", tfModel.Name.ValueString(), "name").
		SetStringField("dashboard_type", tfModel.DashboardType.ValueString(), "dashboard_type").
		SetField("dashboard_configuration", utils.EncodeBase64(configJSON), "")

	return builder.Build()
}

func (t *DashboardTranslatorCommon) buildDashboardConfigurationJSON(requestContext *common.RequestContext, tfModel *model.DashboardTFModel) (string, error) {
	// Parse groups and values from JSON strings
	var groups, values any
	json.Unmarshal([]byte(tfModel.GroupsFull.ValueString()), &groups)
	json.Unmarshal([]byte(tfModel.ValuesFull.ValueString()), &values)

	// Build configuration directly from TF model
	config := map[string]any{
		"resource_query": tfModel.ResourceQuery.ValueString(),
		"groups":         groups,
		"values":         values,
		"other_tags":     utils.ListSliceFromTFModel(requestContext.Context, tfModel.OtherTags),
		"identifiers":    utils.ListSliceFromTFModel(requestContext.Context, tfModel.Identifiers),
	}

	configBytes, _ := json.Marshal(config)
	return string(configBytes), nil
}
