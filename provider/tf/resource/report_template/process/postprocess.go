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

package process

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"
	corecommon "terraform/terraform-provider/provider/tf/core/common"
	"terraform/terraform-provider/provider/tf/core/process"
	reporttemplatetf "terraform/terraform-provider/provider/tf/resource/report_template/model"
	converters "terraform/terraform-provider/provider/tf/resource/report_template/translator/object_converters"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ReportTemplatePostProcessor struct{}

var _ process.PostProcessor[*reporttemplatetf.ReportTemplateTFModel] = &ReportTemplatePostProcessor{}

func (p *ReportTemplatePostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *reporttemplatetf.ReportTemplateTFModel) error {
	return restoreFieldsFromPlan(requestContext, data.CreateRequest.Plan, tfModel)
}

func (p *ReportTemplatePostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, tfModel *reporttemplatetf.ReportTemplateTFModel) error {
	// Process JSON fields first to populate _full attributes from API response for drift detection
	err := postProcessJsonFullFields(requestContext, tfModel)
	if err != nil {
		return err
	}

	// For READ, restore base fields from state (skip _full variants for drift detection)
	// Keep _full fields from API response to enable drift detection
	return restoreBaseFieldsFromState(requestContext, data.ReadRequest.State, tfModel)
}

func (p *ReportTemplatePostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *reporttemplatetf.ReportTemplateTFModel) error {
	return restoreFieldsFromPlan(requestContext, data.UpdateRequest.Plan, tfModel)
}

func (p *ReportTemplatePostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, tfModel *reporttemplatetf.ReportTemplateTFModel) error {
	// No special post-processing needed for delete operations
	return nil
}

// restoreFieldsFromPlan restores fields that need explicit handling from the plan.
// Called for Create/Update operations (after RestoreAllFieldsFromPlan in the orchestrator).
// The translator always populates both deprecated and _list fields from every API response.
// This nulls out the fields that don't belong to the active mode so the result matches the plan.
func restoreFieldsFromPlan(requestContext *common.RequestContext, source corecommon.Getter, tfModel *reporttemplatetf.ReportTemplateTFModel) error {
	var sourceModel reporttemplatetf.ReportTemplateTFModel
	diags := source.Get(requestContext.Context, &sourceModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get plan model: %s", diags.Errors())
	}

	enforceBlocksMode(&sourceModel, tfModel)
	enforceLinksMode(&sourceModel, tfModel)

	return nil
}

// postProcessJsonFullFields processes all _full JSON fields with version-aware logic
func postProcessJsonFullFields(requestContext *common.RequestContext, tfModel *reporttemplatetf.ReportTemplateTFModel) error {

	var err error

	tfModel.BlocksFull, err = postProcessJsonFullField[*customattribute.BlockJson](&tfModel.BlocksFull, requestContext.BackendVersion)
	if err != nil {
		return err
	}

	tfModel.LinksFull, err = postProcessJsonFullField[*customattribute.LinkJson](&tfModel.LinksFull, requestContext.BackendVersion)
	if err != nil {
		return err
	}

	return nil
}

func postProcessJsonFullField[T common.JsonConfigurable](fullField *types.String, backendVersion *version.BackendVersion) (types.String, error) {

	if fullField.IsNull() || fullField.IsUnknown() {
		return *fullField, nil
	}

	fullFieldString := fullField.ValueString()

	// Remarshal the full field to apply the custom struct tags (like min_version, max_version, etc.)
	// and set the default values for the fields that are not present in the JSON
	// See customattribute structs for more details
	// This is necessary to avoid any TF errors in case the backend returns values that are not supported by it (not common, but possible in case of backend bugs)
	overriddenValues, err := common.RemarshalListWithConfig[T](fullFieldString, common.JsonConfig{BackendVersion: backendVersion})
	if err != nil {
		return types.StringNull(), err
	}

	return types.StringValue(overriddenValues), nil
}

// restoreBaseFieldsFromState restores base fields from state (excludes _full computed variants)
// Used during READ to preserve user input while allowing _full fields to reflect API state for drift detection
func restoreBaseFieldsFromState(requestContext *common.RequestContext, source corecommon.Getter, tfModel *reporttemplatetf.ReportTemplateTFModel) error {

	var originalModel reporttemplatetf.ReportTemplateTFModel
	diags := source.Get(requestContext.Context, &originalModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original model: %s", diags.Errors())
	}

	// Restore base fields from state (skip _full variants to enable drift detection)
	tfModel.Blocks = originalModel.Blocks
	tfModel.Links = originalModel.Links

	enforceBlocksMode(&originalModel, tfModel)
	enforceLinksMode(&originalModel, tfModel)

	return nil
}

func enforceBlocksMode(source *reporttemplatetf.ReportTemplateTFModel, tfModel *reporttemplatetf.ReportTemplateTFModel) {
	if common.IsAttrKnown(source.BlocksList) {
		tfModel.Blocks = types.StringNull()
		tfModel.BlocksFull = types.StringNull()
	} else {
		tfModel.BlocksList = types.ListNull(converters.BlocksListObjectType)
	}
}

func enforceLinksMode(source *reporttemplatetf.ReportTemplateTFModel, tfModel *reporttemplatetf.ReportTemplateTFModel) {
	if common.IsAttrKnown(source.LinksList) {
		tfModel.Links = types.StringNull()
		tfModel.LinksFull = types.StringNull()
	} else {
		tfModel.LinksList = types.ListNull(converters.LinksListObjectType)
	}
}
