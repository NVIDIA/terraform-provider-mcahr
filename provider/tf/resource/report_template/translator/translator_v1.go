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
	reporttemplateapi "terraform/terraform-provider/provider/external_api/resources/report_templates"
	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	reporttemplatetf "terraform/terraform-provider/provider/tf/resource/report_template/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReportTemplateTranslatorV1 handles translation between TF models and V1 API models for report template resources
type ReportTemplateTranslatorV1 struct {
	ReportTemplateTranslatorCommon
}

var _ translator.Translator[*reporttemplatetf.ReportTemplateTFModel, *reporttemplateapi.ReportTemplateResponseAPIModelV1] = &ReportTemplateTranslatorV1{}

// ToTFModel converts a V1 API model to a TF model
func (t *ReportTemplateTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *reporttemplateapi.ReportTemplateResponseAPIModelV1) (*reporttemplatetf.ReportTemplateTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the report template container regardless of operation type
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no report template container found in V1 API response")
	}

	if len(container.ReportTemplateClasses) == 0 {
		return nil, fmt.Errorf("no report template classes found in V1 API response")
	}

	// Get the first report template class, current implementation only supports one report template to be returned by the API
	reportTemplateClass := container.ReportTemplateClasses[0]
	jsonConfig := common.JsonConfig{
		BackendVersion: requestContext.BackendVersion,
	}

	blocksJson, err := common.RemarshalListWithConfig[*customattribute.BlockJson](reportTemplateClass.Blocks, jsonConfig)
	if err != nil {
		return nil, err
	}

	linksJson, err := common.RemarshalListWithConfig[*customattribute.LinkJson](reportTemplateClass.Links, jsonConfig)
	if err != nil {
		return nil, err
	}

	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name:       types.StringValue(reportTemplateClass.Name),
		Blocks:     types.StringValue(blocksJson),
		BlocksFull: types.StringValue(blocksJson),
		Links:      types.StringValue(linksJson),
		LinksFull:  types.StringValue(linksJson),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (t *ReportTemplateTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *reporttemplatetf.ReportTemplateTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
