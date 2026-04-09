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

package plan

import (
	"context"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/tf/core/plan"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	"terraform/terraform-provider/provider/tf/resource/runbook/model"
	datamodifier "terraform/terraform-provider/provider/tf/resource/runbook/plan/modifiers/data"
	jsonmodifier "terraform/terraform-provider/provider/tf/resource/runbook/plan/modifiers/jsontype"
	"terraform/terraform-provider/provider/tf/resource/runbook/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema *schema.RunbookSchema, backendVersion *version.BackendVersion) {

	var err error
	var resultValues *model.RunbookTFModel

	doReturn, planValues, configValues, stateValues := plan.GetValues[model.RunbookTFModel](ctx, req, resp)
	if doReturn {
		return
	}

	// Apply data JSON to the config struct
	resultValues, err = datamodifier.ApplyDataModifier(ctx, &configValues)
	if err != nil {
		resp.Diagnostics.AddError("Error applying data JSON", err.Error())
		return
	}

	// Create a copy of the result values without the defaults
	resultValuesWithoutDefaults := resultValues.Copy()

	plan.AddDefaultsFromPlan(resultValues, &planValues)

	handleNullDataFieldPlan(resultValues, &configValues, &planValues, &stateValues)

	// Populate the full JSON attributes with normalized values and defaults.
	// Deprecated JSON fields with active replacements (e.g. cells when cells_list is set)
	// are automatically skipped and nulled inside PopulateFullJsonAttributes.
	err = jsonmodifier.PopulateFullJsonAttributes(ctx, resultValues, resultValuesWithoutDefaults, &planValues, &stateValues, backendVersion)
	if err != nil {
		resp.Diagnostics.AddError("Error populating full JSON attributes", err.Error())
		return
	}

	// Set the result values to the plan
	resp.Diagnostics.Append(resp.Plan.Set(ctx, resultValues)...)

	// For incompatible attributes, this will:
	// - set the plan values to null if they are not provided by the user
	// - raise a validation error if they are provided by the user (in the config)
	// This uses the resultValues to apply the compatibility modifiers to the merged values from root TF config and "data"
	// It is required to not have the defaults in the resultValues to avoid applying them to incompatible attributes (which would raise a validation error all the time)
	compatibility.ApplyCompatibilityModifiers(ctx, &req, resp, schema, backendVersion, resultValuesWithoutDefaults)
}

// handleNullDataFieldPlan ensures that the data field is not marked as computed if the user didn't set it
func handleNullDataFieldPlan(resultValues *model.RunbookTFModel, configValues *model.RunbookTFModel, planValues *model.RunbookTFModel, stateValues *model.RunbookTFModel) {
	// If user didn't provide data in config, it's not planned to change and it's not in state, don't mark it as computed
	if !common.IsAttrKnown(configValues.Data) && !common.IsAttrKnown(planValues.Data) && !common.IsAttrKnown(stateValues.Data) {
		// Don't make it show as "known after apply" in the plan
		resultValues.Data = types.StringNull()
	}
}
