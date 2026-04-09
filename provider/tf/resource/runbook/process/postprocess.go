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
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"
	corecommon "terraform/terraform-provider/provider/tf/core/common"
	"terraform/terraform-provider/provider/tf/core/process"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"
	converters "terraform/terraform-provider/provider/tf/resource/runbook/translator/object_converters"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RunbookPostProcessor struct{}

var _ process.PostProcessor[*runbooktf.RunbookTFModel] = &RunbookPostProcessor{}

func (p *RunbookPostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *runbooktf.RunbookTFModel) (err error) {
	return restoreFieldsFromPlan(requestContext, data.CreateRequest.Plan, tfModel)
}

func (p *RunbookPostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, tfModel *runbooktf.RunbookTFModel) (err error) {

	// Process JSON fields to populate _full attributes from API response for drift detection
	err = postProcessJsonFields(requestContext, tfModel)
	if err != nil {
		return err
	}

	// For READ, restore base fields from state (skip _full variants for drift detection)
	// Keep _full fields from API response to enable drift detection
	err = restoreBaseFieldsFromState(requestContext, data.ReadRequest.State, tfModel)
	if err != nil {
		return err
	}

	return nil
}

func (p *RunbookPostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *runbooktf.RunbookTFModel) (err error) {
	return restoreFieldsFromPlan(requestContext, data.UpdateRequest.Plan, tfModel)
}

func (p *RunbookPostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, tfModel *runbooktf.RunbookTFModel) error {
	return nil
}

//
// Custom postprocess functions
//

func postProcessJsonFields(requestContext *common.RequestContext, tfModel *runbooktf.RunbookTFModel) (err error) {

	tfModel.ParamsFull, err = postProcessJsonFullField[*customattribute.ParamJson](&tfModel.ParamsFull, requestContext.BackendVersion)
	if err != nil {
		return err
	}

	tfModel.ExternalParamsFull, err = postProcessJsonFullField[*customattribute.ExternalParamJson](&tfModel.ExternalParamsFull, requestContext.BackendVersion)
	if err != nil {
		return err
	}

	tfModel.CellsFull, err = postProcessJsonFullField[*customattribute.CellJson](&tfModel.CellsFull, requestContext.BackendVersion)
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

// restoreFieldsFromPlan restores fields that need explicit handling from the plan.
// Called for Create/Update operations (after RestoreAllFieldsFromPlan in the orchestrator).
func restoreFieldsFromPlan(requestContext *common.RequestContext, source corecommon.Getter, tfModel *runbooktf.RunbookTFModel) error {

	var sourceModel runbooktf.RunbookTFModel
	diags := source.Get(requestContext.Context, &sourceModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get source model: %s", diags.Errors())
	}

	// params_groups is set to null by NullObjectIfUnknownModifier when not configured by the user,
	// but the API always computes and returns it. Without explicit restoration, the API value leaks
	// into state, causing "inconsistent result after apply" errors.
	setParamsGroups(sourceModel.ParamsGroups, tfModel)

	// Enforce cells mode: the translator populates cells, cells_full, and cells_list from every
	// API response. Null out the fields that don't belong to the active mode so the plan matches.
	enforceCellsMode(&sourceModel, tfModel)

	return nil
}

// setParamsGroups copies params_groups from a source value to the model, normalizing unknown to null.
func setParamsGroups(source types.Object, tfModel *runbooktf.RunbookTFModel) {
	if source.IsUnknown() {
		tfModel.ParamsGroups = types.ObjectNull(converters.ParamsGroupsAttrTypes)
	} else {
		tfModel.ParamsGroups = source
	}
}

// restoreBaseFieldsFromState restores base fields and special fields from state (excludes _full computed variants)
// Used during READ to preserve user input while allowing _full fields to reflect API state for drift detection
func restoreBaseFieldsFromState(requestContext *common.RequestContext, source corecommon.Getter, tfModel *runbooktf.RunbookTFModel) error {

	var originalModel runbooktf.RunbookTFModel
	diags := source.Get(requestContext.Context, &originalModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original model: %s", diags.Errors())
	}

	// Restore base fields from state (skip _full variants to enable drift detection)
	tfModel.Cells = originalModel.Cells
	tfModel.Params = originalModel.Params
	tfModel.ExternalParams = originalModel.ExternalParams

	// Enforce cells mode: the translator populates all three (cells, cells_full, cells_list)
	// from every API response. Null out the fields that don't belong to the active mode.
	enforceCellsMode(&originalModel, tfModel)

	// Restore special feature field
	tfModel.Data = originalModel.Data

	setParamsGroups(originalModel.ParamsGroups, tfModel)

	return nil
}

// enforceCellsMode ensures only one cells representation is active in the model.
// The translator always populates cells, cells_full, and cells_list from the API response.
// This function nulls out the fields that don't belong to the active mode:
//   - cells_list mode (source has cells_list set): cells and cells_full stay at "[]"
//   - cells mode (source has cells_list null): cells_list is set to null
func enforceCellsMode(source *runbooktf.RunbookTFModel, tfModel *runbooktf.RunbookTFModel) {
	if common.IsAttrKnown(source.CellsList) {
		// cells_list mode: cells and cells_full must be null
		tfModel.Cells = types.StringNull()
		tfModel.CellsFull = types.StringNull()
	} else {
		// cells mode: cells_list must be null
		tfModel.CellsList = types.ListNull(converters.CellsListObjectType)
	}
}
