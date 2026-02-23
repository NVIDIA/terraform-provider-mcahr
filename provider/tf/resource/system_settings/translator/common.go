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

package translator

import (
	"fmt"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	systemsettingstf "terraform/terraform-provider/provider/tf/resource/system_settings/model"
)

// SystemSettingsTranslatorCommon provides common functionality for system_settings translators across API versions
type SystemSettingsTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (t *SystemSettingsTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *systemsettingstf.SystemSettingsTFModel) (*statement.StatementInputAPIModel, error) {
	var stmt string

	switch requestContext.Operation {
	case common.Create:
		// For system_settings, create uses update statement
		stmt = t.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = t.buildReadStatement(tfModel)
	case common.Update:
		stmt = t.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		// system_settings doesn't support delete
		return nil, fmt.Errorf("delete operation is not supported for system_settings resource")
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}

	apiModel := &statement.StatementInputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (t *SystemSettingsTranslatorCommon) buildReadStatement(tfModel *systemsettingstf.SystemSettingsTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_configuration_class(configuration_name=\"%s\")", name)
}

func (t *SystemSettingsTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *systemsettingstf.SystemSettingsTFModel) string {
	// Build the update_configuration statement from the TF model using the builder pattern
	builder := utils.NewStatementBuilder("update_configuration", requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("configuration_name", tfModel.Name.ValueString(), "name")

	// Set all configuration fields
	builder = builder.
		SetField("administrator_grants_create_user", tfModel.AdministratorGrantsCreateUser.ValueBool(), "administrator_grants_create_user").
		SetField("administrator_grants_create_user_token", tfModel.AdministratorGrantsCreateUserToken.ValueBool(), "administrator_grants_create_user_token").
		SetField("administrator_grants_regenerate_user_token", tfModel.AdministratorGrantsRegenerateUserToken.ValueBool(), "administrator_grants_regenerate_user_token").
		SetField("administrator_grants_read_user_token", tfModel.AdministratorGrantsReadUserToken.ValueBool(), "administrator_grants_read_user_token").
		SetField("approval_feature_enabled", tfModel.ApprovalFeatureEnabled.ValueBool(), "approval_feature_enabled").
		SetField("runbook_ad_hoc_approval_request_enabled", tfModel.RunbookAdHocApprovalRequestEnabled.ValueBool(), "runbook_ad_hoc_approval_request_enabled").
		SetField("runbook_approval_request_expiry_time", tfModel.RunbookApprovalRequestExpiryTime.ValueInt64(), "runbook_approval_request_expiry_time").
		SetField("run_approval_expiry_time", tfModel.RunApprovalExpiryTime.ValueInt64(), "run_approval_expiry_time").
		SetField("approval_editable_allowed_resource_query_enabled", tfModel.ApprovalEditableAllowedResourceQueryEnabled.ValueBool(), "approval_editable_allowed_resource_query_enabled").
		SetField("approval_allow_individual_notification", tfModel.ApprovalAllowIndividualNotification.ValueBool(), "approval_allow_individual_notification").
		SetField("approval_optional_request_ticket_url", tfModel.ApprovalOptionalRequestTicketURL.ValueBool(), "approval_optional_request_ticket_url").
		SetStringField("time_trigger_permissions_user", tfModel.TimeTriggerPermissionsUser.ValueString(), "time_trigger_permissions_user").
		SetField("external_audit_storage_enabled", tfModel.ExternalAuditStorageEnabled.ValueBool(), "external_audit_storage_enabled").
		SetStringField("external_audit_storage_type", tfModel.ExternalAuditStorageType.ValueString(), "external_audit_storage_type").
		SetField("external_audit_storage_batch_period_sec", tfModel.ExternalAuditStorageBatchPeriodSec.ValueInt64(), "external_audit_storage_batch_period_sec").
		SetStringField("environment_name", tfModel.EnvironmentName.ValueString(), "environment_name").
		SetStringField("environment_name_background", tfModel.EnvironmentNameBackground.ValueString(), "environment_name_background").
		SetField("param_value_max_length", tfModel.ParamValueMaxLength.ValueInt64(), "param_value_max_length").
		SetField("parallel_runs_fired_by_time_triggers", tfModel.ParallelRunsFiredByTimeTriggers.ValueInt64(), "parallel_runs_fired_by_time_triggers").
		SetField("maintenance_mode_enabled", tfModel.MaintenanceModeEnabled.ValueBool(), "maintenance_mode_enabled").
		SetStringField("managed_secrets", tfModel.ManagedSecrets.ValueString(), "managed_secrets")

	// Handle string sets - only add to statement if explicitly set (not null/unknown)
	if common.IsAttrKnown(tfModel.AllowedTags) {
		builder = builder.SetArrayField("allowed_tags", utils.ListSliceFromTFModel(requestContext.Context, tfModel.AllowedTags), "allowed_tags")
	}
	if common.IsAttrKnown(tfModel.SkippedTags) {
		builder = builder.SetArrayField("skipped_tags", utils.ListSliceFromTFModel(requestContext.Context, tfModel.SkippedTags), "skipped_tags")
	}

	return builder.Build()
}
