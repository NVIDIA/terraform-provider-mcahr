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
	timetriggertf "terraform/terraform-provider/provider/tf/resource/time_trigger/model"
)

// TimeTriggerTranslatorCommon contains shared functionality between V1 and V2 translators
type TimeTriggerTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (t *TimeTriggerTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *timetriggertf.TimeTriggerTFModel) (*statement.StatementInputAPIModel, error) {
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

func (t *TimeTriggerTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *timetriggertf.TimeTriggerTFModel) string {
	return t.buildTimeTriggerStatement(requestContext, translationData, "define_time_trigger", tfModel)
}

func (t *TimeTriggerTranslatorCommon) buildReadStatement(tfModel *timetriggertf.TimeTriggerTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_time_trigger_class(time_trigger_name=%s)", utils.EscapeString(name))
}

func (t *TimeTriggerTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *timetriggertf.TimeTriggerTFModel) string {
	return t.buildTimeTriggerStatement(requestContext, translationData, "update_time_trigger", tfModel)
}

func (t *TimeTriggerTranslatorCommon) buildDeleteStatement(tfModel *timetriggertf.TimeTriggerTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_time_trigger(time_trigger_name=%s)", utils.EscapeString(name))
}

func (t *TimeTriggerTranslatorCommon) buildTimeTriggerStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *timetriggertf.TimeTriggerTFModel) string {
	return utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("time_trigger_name", tfModel.Name.ValueString(), "name").
		SetCommandField("fire_query", tfModel.FireQuery.ValueString(), "fire_query").
		SetStringField("start_date", tfModel.StartDate.ValueString(), "start_date").
		SetStringField("end_date", tfModel.EndDate.ValueString(), "end_date").
		SetField("enabled", tfModel.Enabled.ValueBool(), "enabled").
		Build()
}
