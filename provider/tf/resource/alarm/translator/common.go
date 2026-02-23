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
	"fmt"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	alarmtf "terraform/terraform-provider/provider/tf/resource/alarm/model"
)

// AlarmTranslatorCommon contains shared functionality between V1 and V2 translators
type AlarmTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (a *AlarmTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *alarmtf.AlarmTFModel) (*statement.StatementInputAPIModel, error) {
	var stmt string

	switch requestContext.Operation {
	case common.Create:
		stmt = a.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = a.buildReadStatement(tfModel)
	case common.Update:
		stmt = a.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt = a.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}

	apiModel := &statement.StatementInputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (a *AlarmTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *alarmtf.AlarmTFModel) string {
	return a.buildAlarmStatement(requestContext, translationData, "define_alarm", tfModel)
}

func (a *AlarmTranslatorCommon) buildReadStatement(tfModel *alarmtf.AlarmTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_alarm_class(alarm_name=\"%s\")", name)
}

func (a *AlarmTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *alarmtf.AlarmTFModel) string {
	return a.buildAlarmStatement(requestContext, translationData, "update_alarm", tfModel)
}

func (a *AlarmTranslatorCommon) buildDeleteStatement(tfModel *alarmtf.AlarmTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_alarm(alarm_name=\"%s\")", name)
}

func (a *AlarmTranslatorCommon) buildAlarmStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *alarmtf.AlarmTFModel) string {
	// Build the alarm statement from the TF model using the builder pattern
	// V1 and V2 use the same fields (deprecated fields excluded from both versions)

	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("alarm_name", tfModel.Name.ValueString(), "name").
		SetStringField("fire_query", tfModel.FireQuery.ValueString(), "fire_query").
		SetStringField("clear_query", tfModel.ClearQuery.ValueString(), "clear_query").
		SetStringField("description", tfModel.Description.ValueString(), "description").
		SetStringField("resource_query", tfModel.ResourceQuery.ValueString(), "resource_query").
		SetStringField("resource_type", tfModel.ResourceType.ValueString(), "resource_type").
		SetField("check_interval_sec", tfModel.CheckIntervalSec.ValueInt64(), "check_interval_sec").
		SetStringField("family", tfModel.Family.ValueString(), "family").
		SetField("enabled", tfModel.Enabled.ValueBool(), "enabled").
		SetStringField("fire_title_template", tfModel.FireTitleTemplate.ValueString(), "fire_title_template").
		SetStringField("fire_short_template", tfModel.FireShortTemplate.ValueString(), "fire_short_template").
		SetStringField("resolve_title_template", tfModel.ResolveTitleTemplate.ValueString(), "resolve_title_template").
		SetStringField("resolve_short_template", tfModel.ResolveShortTemplate.ValueString(), "resolve_short_template")

	return builder.Build()
}
