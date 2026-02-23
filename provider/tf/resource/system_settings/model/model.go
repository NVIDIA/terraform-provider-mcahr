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

var _ core.TFModel = &SystemSettingsTFModel{} // check that SystemSettingsTFModel implements TFModel

// SystemSettingsTFModel represents the Terraform model for system_settings resources
type SystemSettingsTFModel struct {
	Name                                        types.String `tfsdk:"name"`
	AdministratorGrantsCreateUser               types.Bool   `tfsdk:"administrator_grants_create_user"`
	AdministratorGrantsCreateUserToken          types.Bool   `tfsdk:"administrator_grants_create_user_token"`
	AdministratorGrantsRegenerateUserToken      types.Bool   `tfsdk:"administrator_grants_regenerate_user_token"`
	AdministratorGrantsReadUserToken            types.Bool   `tfsdk:"administrator_grants_read_user_token"`
	ApprovalFeatureEnabled                      types.Bool   `tfsdk:"approval_feature_enabled"`
	RunbookAdHocApprovalRequestEnabled          types.Bool   `tfsdk:"runbook_ad_hoc_approval_request_enabled"`
	RunbookApprovalRequestExpiryTime            types.Int64  `tfsdk:"runbook_approval_request_expiry_time"`
	RunApprovalExpiryTime                       types.Int64  `tfsdk:"run_approval_expiry_time"`
	ApprovalEditableAllowedResourceQueryEnabled types.Bool   `tfsdk:"approval_editable_allowed_resource_query_enabled"`
	ApprovalAllowIndividualNotification         types.Bool   `tfsdk:"approval_allow_individual_notification"`
	ApprovalOptionalRequestTicketURL            types.Bool   `tfsdk:"approval_optional_request_ticket_url"`
	TimeTriggerPermissionsUser                  types.String `tfsdk:"time_trigger_permissions_user"`
	ExternalAuditStorageEnabled                 types.Bool   `tfsdk:"external_audit_storage_enabled"`
	ExternalAuditStorageType                    types.String `tfsdk:"external_audit_storage_type"`
	ExternalAuditStorageBatchPeriodSec          types.Int64  `tfsdk:"external_audit_storage_batch_period_sec"`
	EnvironmentName                             types.String `tfsdk:"environment_name"`
	EnvironmentNameBackground                   types.String `tfsdk:"environment_name_background"`
	ParamValueMaxLength                         types.Int64  `tfsdk:"param_value_max_length"`
	ParallelRunsFiredByTimeTriggers             types.Int64  `tfsdk:"parallel_runs_fired_by_time_triggers"`
	MaintenanceModeEnabled                      types.Bool   `tfsdk:"maintenance_mode_enabled"`
	AllowedTags                                 types.List   `tfsdk:"allowed_tags"`
	SkippedTags                                 types.List   `tfsdk:"skipped_tags"`
	ManagedSecrets                              types.String `tfsdk:"managed_secrets"`
}

// GetName returns the name of the system_settings resource
func (s SystemSettingsTFModel) GetName() string {
	return s.Name.ValueString()
}
