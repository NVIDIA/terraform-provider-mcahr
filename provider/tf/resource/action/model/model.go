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

var _ core.TFModel = &ActionTFModel{} // check that ActionTFModel implements TFModel

type ActionTFModel struct {
	// Required fields
	Name    types.String `tfsdk:"name" json:"name"`
	Command types.String `tfsdk:"command" json:"command"`

	// Optional fields
	Description            types.String `tfsdk:"description" json:"description,omitempty"`
	ResEnvVar              types.String `tfsdk:"res_env_var" json:"res_env_var,omitempty"`
	ResourceQuery          types.String `tfsdk:"resource_query" json:"resource_query,omitempty"`
	Shell                  types.String `tfsdk:"shell" json:"shell,omitempty"`
	StartShortTemplate     types.String `tfsdk:"start_short_template" json:"start_short_template,omitempty"`
	StartLongTemplate      types.String `tfsdk:"start_long_template" json:"start_long_template,omitempty"`
	StartTitleTemplate     types.String `tfsdk:"start_title_template" json:"start_title_template,omitempty"`
	ErrorShortTemplate     types.String `tfsdk:"error_short_template" json:"error_short_template,omitempty"`
	ErrorLongTemplate      types.String `tfsdk:"error_long_template" json:"error_long_template,omitempty"`
	ErrorTitleTemplate     types.String `tfsdk:"error_title_template" json:"error_title_template,omitempty"`
	CompleteShortTemplate  types.String `tfsdk:"complete_short_template" json:"complete_short_template,omitempty"`
	CompleteLongTemplate   types.String `tfsdk:"complete_long_template" json:"complete_long_template,omitempty"`
	CompleteTitleTemplate  types.String `tfsdk:"complete_title_template" json:"complete_title_template,omitempty"`
	CommunicationWorkspace types.String `tfsdk:"communication_workspace" json:"communication_workspace,omitempty"`
	CommunicationChannel   types.String `tfsdk:"communication_channel" json:"communication_channel,omitempty"`
	AllowedResourcesQuery  types.String `tfsdk:"allowed_resources_query" json:"allowed_resources_query,omitempty"`
	Enabled                types.Bool   `tfsdk:"enabled" json:"enabled,omitempty"`
	Timeout                types.Int64  `tfsdk:"timeout" json:"timeout,omitempty"`
	Params                 types.List   `tfsdk:"params" json:"params,omitempty"`
	ResourceTagsToExport   types.List   `tfsdk:"resource_tags_to_export" json:"resource_tags_to_export,omitempty"`
	FileDeps               types.List   `tfsdk:"file_deps" json:"file_deps,omitempty"`
	AllowedEntities        types.List   `tfsdk:"allowed_entities" json:"allowed_entities,omitempty"`
	Editors                types.List   `tfsdk:"editors" json:"editors,omitempty"`
}

func (a *ActionTFModel) GetName() string {
	return a.Name.ValueString()
}
