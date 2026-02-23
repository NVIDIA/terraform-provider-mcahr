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

package system_settings

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// SystemSettingsConfigurationV2 represents the system_settings configuration in V2 API responses
type SystemSettingsConfigurationV2 struct {
	AdministratorGrantsCreateUser               bool     `json:"administrator_grants_create_user"`
	AdministratorGrantsCreateUserToken          bool     `json:"administrator_grants_create_user_token"`
	AdministratorGrantsRegenerateUserToken      bool     `json:"administrator_grants_regenerate_user_token"`
	AdministratorGrantsReadUserToken            bool     `json:"administrator_grants_read_user_token"`
	ApprovalFeatureEnabled                      bool     `json:"approval_feature_enabled"`
	RunbookAdHocApprovalRequestEnabled          bool     `json:"runbook_ad_hoc_approval_request_enabled"`
	RunbookApprovalRequestExpiryTime            int64    `json:"runbook_approval_request_expiry_time"`
	RunApprovalExpiryTime                       int64    `json:"run_approval_expiry_time"`
	ApprovalEditableAllowedResourceQueryEnabled bool     `json:"approval_editable_allowed_resource_query_enabled"`
	ApprovalAllowIndividualNotification         bool     `json:"approval_allow_individual_notification"`
	ApprovalOptionalRequestTicketURL            bool     `json:"approval_optional_request_ticket_url"`
	TimeTriggerPermissionsUser                  string   `json:"time_trigger_permissions_user"`
	ExternalAuditStorageEnabled                 bool     `json:"external_audit_storage_enabled"`
	ExternalAuditStorageType                    string   `json:"external_audit_storage_type"`
	ExternalAuditStorageBatchPeriodSec          int64    `json:"external_audit_storage_batch_period_sec"`
	EnvironmentName                             string   `json:"environment_name"`
	EnvironmentNameBackground                   string   `json:"environment_name_background"`
	ParamValueMaxLength                         int64    `json:"param_value_max_length"`
	ParallelRunsFiredByTimeTriggers             int64    `json:"parallel_runs_fired_by_time_triggers"`
	MaintenanceModeEnabled                      bool     `json:"maintenance_mode_enabled"`
	AllowedTags                                 []string `json:"allowed_tags"`
	SkippedTags                                 []string `json:"skipped_tags"`
	ManagedSecrets                              string   `json:"managed_secrets"`
}

// SystemSettingsItemV2 represents a single system_settings item in V2 API responses
type SystemSettingsItemV2 struct {
	Name          string                        `json:"name"`
	Configuration SystemSettingsConfigurationV2 `json:"configuration"`
}

// SystemConfigurationsV2 represents the system_configurations section in V2 API responses
type SystemConfigurationsV2 struct {
	Count int                    `json:"count"`
	Items []SystemSettingsItemV2 `json:"items"`
}

// OutputV2 represents the output section in V2 API responses
type OutputV2 struct {
	SystemConfigurations SystemConfigurationsV2 `json:"system_configurations"`
}

// SummaryV2 represents the summary section in V2 API responses
type SummaryV2 struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// SystemSettingsResponseAPIModel represents the complete V2 API response for system_settings operations
type SystemSettingsResponseAPIModel struct {
	Output  OutputV2  `json:"output"`
	Summary SummaryV2 `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (s SystemSettingsResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(s.Summary.Status, s.Summary.Errors)
}
