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
	reporttemplateapi "terraform/terraform-provider/provider/external_api/resources/report_templates"
	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	reporttemplatetf "terraform/terraform-provider/provider/tf/resource/report_template/model"
	converters "terraform/terraform-provider/provider/tf/resource/report_template/translator/object_converters"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReportTemplateTranslator handles translation between TF models and V2 API models for report template resources
type ReportTemplateTranslator struct {
	ReportTemplateTranslatorCommon
}

var _ translator.Translator[*reporttemplatetf.ReportTemplateTFModel, *reporttemplateapi.ReportTemplateResponseAPIModel] = &ReportTemplateTranslator{}

// ToTFModel converts a V2 API model to a TF model
func (t *ReportTemplateTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *reporttemplateapi.ReportTemplateResponseAPIModel) (*reporttemplatetf.ReportTemplateTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no report template configurations found in V2 API response")
	}

	// Get the first configuration item, current implementation only supports one report template to be returned by the API
	configItem := apiModel.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name: types.StringValue(metadata.Name),
	}

	if err := toTFModelJsonFields(tfModel, config.Blocks, config.Links); err != nil {
		return nil, err
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V2 backend
func (t *ReportTemplateTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *reporttemplatetf.ReportTemplateTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}

// toTFModelJsonFields populates JSON and list fields on the TF model from parsed API structs.
// Called by both V1 and V2 translators.
func toTFModelJsonFields(tfModel *reporttemplatetf.ReportTemplateTFModel, blocks []customattribute.BlockJson, links []customattribute.LinkJson) error {
	// Blocks
	blocksJSON, err := json.Marshal(blocks)
	if err != nil {
		return fmt.Errorf("failed to marshal blocks: %w", err)
	}
	blocksValue := types.StringValue(string(blocksJSON))
	tfModel.Blocks = blocksValue
	tfModel.BlocksFull = blocksValue

	blocksList, bDiags := converters.BlocksListFromAPI(blocks)
	if bDiags.HasError() {
		return fmt.Errorf("failed to convert blocks to blocks_list: %s", bDiags.Errors())
	}
	tfModel.BlocksList = blocksList

	// Links
	linksJSON, err := json.Marshal(links)
	if err != nil {
		return fmt.Errorf("failed to marshal links: %w", err)
	}
	linksValue := types.StringValue(string(linksJSON))
	if len(links) == 0 {
		linksValue = types.StringValue("[]")
	}
	tfModel.Links = linksValue
	tfModel.LinksFull = linksValue

	linksList, lDiags := converters.LinksListFromAPI(links)
	if lDiags.HasError() {
		return fmt.Errorf("failed to convert links to links_list: %s", lDiags.Errors())
	}
	tfModel.LinksList = linksList

	return nil
}
