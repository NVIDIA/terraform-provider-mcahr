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

package datavalidator

import (
	"context"
	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	data "terraform/terraform-provider/provider/tf/resource/runbook/data_attribute"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func ApplyDataValidators(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// If destroy then do nothing
	if req.Config.Raw.IsNull() {
		return
	}

	var rootModel model.RunbookTFModel
	diags := req.Config.Get(ctx, &rootModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataMap, err := data.ParseDataJSONToMap(rootModel.Data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Data JSON Parse Error",
			err.Error(),
		)
		return
	} else if dataMap == nil {
		return
	}

	// Validate that fields are not set in both root TF model and data JSON
	if err := validateNoFieldConflicts(ctx, &rootModel, dataMap); err != nil {
		resp.Diagnostics.AddError(
			"Data JSON Field Conflict Error",
			err.Error(),
		)
	}

	// Validate that required fields are set
	if err := validateRequiredFields(&rootModel, dataMap); err != nil {
		resp.Diagnostics.AddError(
			"Missing required argument",
			err.Error(),
		)
	}

	// Validate data JSON cell structure (type field is required for conversion to op/md)
	if err := validateDataCells(dataMap); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Cell in Data JSON",
			err.Error(),
		)
	}
}
