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

package model

import (
	core "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ core.TFModel = &TimeTriggerTFModel{}

type TimeTriggerTFModel struct {
	// Required fields
	Name      types.String `tfsdk:"name" json:"name"`
	FireQuery types.String `tfsdk:"fire_query" json:"fire_query"`

	// Optional fields
	StartDate types.String `tfsdk:"start_date" json:"start_date,omitempty"`
	EndDate   types.String `tfsdk:"end_date" json:"end_date,omitempty"`
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled,omitempty"`
}

func (t *TimeTriggerTFModel) GetName() string {
	return t.Name.ValueString()
}
