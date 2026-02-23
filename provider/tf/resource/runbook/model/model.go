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

var _ core.TFModel = &RunbookTFModel{}

type RunbookTFModel struct {
	// Required fields
	Name types.String `tfsdk:"name" json:"name"`

	// Optional fields with complex types
	Cells     types.String `tfsdk:"cells" json:"cells,omitempty"`
	CellsFull types.String `tfsdk:"cells_full" json:"cells_full,omitempty"`

	Params     types.String `tfsdk:"params" json:"params,omitempty"`
	ParamsFull types.String `tfsdk:"params_full" json:"params_full,omitempty"`

	ExternalParams     types.String `tfsdk:"external_params" json:"external_params,omitempty"`
	ExternalParamsFull types.String `tfsdk:"external_params_full" json:"external_params_full,omitempty"`

	// Optional boolean fields
	Enabled types.Bool `tfsdk:"enabled" json:"enabled,omitempty"`

	// Optional string fields
	Description            types.String `tfsdk:"description" json:"description,omitempty"`
	AllowedResourcesQuery  types.String `tfsdk:"allowed_resources_query" json:"allowed_resources_query,omitempty"`
	CommunicationWorkspace types.String `tfsdk:"communication_workspace" json:"communication_workspace,omitempty"`
	CommunicationChannel   types.String `tfsdk:"communication_channel" json:"communication_channel,omitempty"`
	Category               types.String `tfsdk:"category" json:"category,omitempty"`

	// Optional numeric fields
	TimeoutMs types.Int64 `tfsdk:"timeout_ms" json:"timeout_ms,omitempty"`

	// Optional boolean fields for communication
	CommunicationCudNotifications       types.Bool `tfsdk:"communication_cud_notifications" json:"communication_cud_notifications,omitempty"`
	CommunicationApprovalNotifications  types.Bool `tfsdk:"communication_approval_notifications" json:"communication_approval_notifications,omitempty"`
	CommunicationExecutionNotifications types.Bool `tfsdk:"communication_execution_notifications" json:"communication_execution_notifications,omitempty"`
	IsRunOutputPersisted                types.Bool `tfsdk:"is_run_output_persisted" json:"is_run_output_persisted,omitempty"`
	FilterResourceToAction              types.Bool `tfsdk:"filter_resource_to_action" json:"filter_resource_to_action,omitempty"`

	// Set attributes
	AllowedEntities types.List `tfsdk:"allowed_entities" json:"allowed_entities,omitempty"`
	Approvers       types.List `tfsdk:"approvers" json:"approvers,omitempty"`
	Labels          types.List `tfsdk:"labels" json:"labels,omitempty"`
	Editors         types.List `tfsdk:"editors" json:"editors,omitempty"`
	SecretNames     types.List `tfsdk:"secret_names" json:"secret_names,omitempty"`

	// Nested attributes
	ParamsGroups types.Object `tfsdk:"params_groups" json:"params_groups,omitempty"`

	// Data fields
	Data types.String `tfsdk:"data" json:"data,omitempty"`
}

func (r *RunbookTFModel) GetName() string {
	return r.Name.ValueString()
}

// Copy creates a copy of the RunbookTFModel
// This preserves the internal state of Terraform types (null, unknown, known)
func (r *RunbookTFModel) Copy() *RunbookTFModel {
	if r == nil {
		return nil
	}
	// SDK types are immutable, so this acts like a deep copy
	copy := *r
	return &copy
}
