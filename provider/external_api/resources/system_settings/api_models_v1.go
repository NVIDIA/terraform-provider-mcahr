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

import (
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
)

// SystemSettingsV1 represents the system_settings configuration in V1 API responses
type SystemSettingsV1 struct {
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

// SystemSettingsContainerV1 represents the container for system_settings in V1 API responses
type SystemSettingsContainerV1 struct {
	SystemSettings SystemSettingsV1  `json:"system_settings"`
	Error          apicommon.ErrorV1 `json:"error"`
	Errors         []string          `json:"errors,omitempty"` // Direct errors array for some error cases
}

// GetNestedError returns the nested error structure
func (c *SystemSettingsContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *SystemSettingsContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// SystemSettingsResponseAPIModelV1 represents the complete V1 API response for system_settings operations
type SystemSettingsResponseAPIModelV1 struct {
	ConfigurationName     string                     `json:"configuration_name,omitempty"`
	GetConfigurationClass *SystemSettingsContainerV1 `json:"get_configuration_class,omitempty"`
	UpdateConfiguration   *SystemSettingsContainerV1 `json:"update_configuration,omitempty"`
	Errors                *apicommon.SyntaxErrorsV1  `json:"errors,omitempty"` // Top-level syntax errors
}

// GetContainer returns the appropriate container based on the operation type
func (s SystemSettingsResponseAPIModelV1) GetContainer() *SystemSettingsContainerV1 {
	if s.GetConfigurationClass != nil {
		return s.GetConfigurationClass
	}
	if s.UpdateConfiguration != nil {
		return s.UpdateConfiguration
	}
	return nil
}

// GetErrors returns a string representation of the API response errors
func (s SystemSettingsResponseAPIModelV1) GetErrors() string {
	container := s.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(s.Errors, container)
}
