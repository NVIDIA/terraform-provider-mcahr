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
	"encoding/base64"
	"fmt"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	reporttemplatetf "terraform/terraform-provider/provider/tf/resource/report_template/model"
)

// ReportTemplateTranslatorCommon contains shared logic for report template translators
type ReportTemplateTranslatorCommon struct{}

// ToAPIModelWithVersion creates a statement API model for the given TF model and API version
func (t *ReportTemplateTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *reporttemplatetf.ReportTemplateTFModel) (*statement.StatementInputAPIModel, error) {
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

func (t *ReportTemplateTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *reporttemplatetf.ReportTemplateTFModel) string {
	return t.buildReportTemplateStatement(requestContext, translationData, "define_report_template", tfModel)
}

func (t *ReportTemplateTranslatorCommon) buildReadStatement(tfModel *reporttemplatetf.ReportTemplateTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_report_template_class(report_template_name=\"%s\")", name)
}

func (t *ReportTemplateTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *reporttemplatetf.ReportTemplateTFModel) string {
	return t.buildReportTemplateStatement(requestContext, translationData, "update_report_template", tfModel)
}

func (t *ReportTemplateTranslatorCommon) buildDeleteStatement(tfModel *reporttemplatetf.ReportTemplateTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_report_template(report_template_name=\"%s\")", name)
}

func (t *ReportTemplateTranslatorCommon) buildReportTemplateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *reporttemplatetf.ReportTemplateTFModel) string {

	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("report_template_name", tfModel.Name.ValueString(), "name")

	builder = builder.SetField("blocks", t.encodeBase64(tfModel.BlocksFull.ValueString()), "blocks")
	builder = builder.SetField("links", t.encodeBase64(tfModel.LinksFull.ValueString()), "links")

	return builder.Build()
}

// encodeBase64 encodes the JSON string to base64 with quotes
func (t ReportTemplateTranslatorCommon) encodeBase64(jsonStr string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(jsonStr))
	return fmt.Sprintf("\"%s\"", encoded)
}
