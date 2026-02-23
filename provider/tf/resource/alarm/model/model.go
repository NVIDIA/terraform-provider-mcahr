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

var _ core.TFModel = &AlarmTFModel{}

type AlarmTFModel struct {
	// Required fields
	Name      types.String `tfsdk:"name" json:"name"`
	FireQuery types.String `tfsdk:"fire_query" json:"fire_query"`

	// Optional fields
	ClearQuery       types.String `tfsdk:"clear_query" json:"clear_query,omitempty"`
	MuteQuery        types.String `tfsdk:"mute_query" json:"mute_query,omitempty"`
	Description      types.String `tfsdk:"description" json:"description,omitempty"`
	ResourceQuery    types.String `tfsdk:"resource_query" json:"resource_query,omitempty"`
	ResourceType     types.String `tfsdk:"resource_type" json:"resource_type,omitempty"`
	CheckIntervalSec types.Int64  `tfsdk:"check_interval_sec" json:"check_interval_sec,omitempty"`
	ConditionType    types.String `tfsdk:"condition_type" json:"condition_type,omitempty"`
	ConditionValue   types.String `tfsdk:"condition_value" json:"condition_value,omitempty"`
	MetricName       types.String `tfsdk:"metric_name" json:"metric_name,omitempty"`
	RaiseFor         types.String `tfsdk:"raise_for" json:"raise_for,omitempty"`
	Family           types.String `tfsdk:"family" json:"family,omitempty"`
	Enabled          types.Bool   `tfsdk:"enabled" json:"enabled,omitempty"`

	// Template fields
	FireTitleTemplate    types.String `tfsdk:"fire_title_template" json:"fire_title_template,omitempty"`
	FireLongTemplate     types.String `tfsdk:"fire_long_template" json:"fire_long_template,omitempty"`
	FireShortTemplate    types.String `tfsdk:"fire_short_template" json:"fire_short_template,omitempty"`
	ResolveTitleTemplate types.String `tfsdk:"resolve_title_template" json:"resolve_title_template,omitempty"`
	ResolveLongTemplate  types.String `tfsdk:"resolve_long_template" json:"resolve_long_template,omitempty"`
	ResolveShortTemplate types.String `tfsdk:"resolve_short_template" json:"resolve_short_template,omitempty"`
}

func (a *AlarmTFModel) GetName() string {
	return a.Name.ValueString()
}
