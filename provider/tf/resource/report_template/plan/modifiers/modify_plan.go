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
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/tf/core/plan"
	"terraform/terraform-provider/provider/tf/resource/report_template/model"
	jsonmodifier "terraform/terraform-provider/provider/tf/resource/report_template/plan/modifiers/jsontype"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, backendVersion *version.BackendVersion) {

	var err error
	var resultValues *model.ReportTemplateTFModel

	doReturn, planValues, configValues, stateValues := plan.GetValues[model.ReportTemplateTFModel](ctx, req, resp)
	if doReturn {
		return
	}

	// Make a shallow copy to avoid modifying the original config
	resultValuesCopy := configValues
	resultValues = &resultValuesCopy

	// Apply defaults from plan (important for fields with schema defaults)
	plan.AddDefaultsFromPlan(resultValues, &planValues)

	// Populate the full JSON attributes with normalized values and defaults.
	// Deprecated JSON fields with active replacements (e.g. blocks when blocks_list is set)
	// are automatically skipped and nulled inside PopulateFullJsonAttributes.
	// configValues (pre-defaults) is passed to distinguish user-set fields from plan defaults.
	err = jsonmodifier.PopulateFullJsonAttributes(ctx, resultValues, &configValues, &planValues, &stateValues, backendVersion)
	if err != nil {
		resp.Diagnostics.AddError("Error populating full JSON attributes", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, resultValues)...)
}
