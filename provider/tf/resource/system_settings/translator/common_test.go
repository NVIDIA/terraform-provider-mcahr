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
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	systemsettingstf "terraform/terraform-provider/provider/tf/resource/system_settings/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemSettingsTranslatorCommon_ToAPIModel(t *testing.T) {
	tests := []struct {
		name      string
		operation common.CrudOperation
	}{
		{"Create operation", common.Create},
		{"Read operation", common.Read},
		{"Update operation", common.Update},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			translator := &SystemSettingsTranslatorCommon{}
			tfModel := createTestSystemSettingsTFModel()
			requestContext := common.NewRequestContext(context.Background()).WithOperation(tt.operation).WithAPIVersion(common.V1)
			translationData := &coretranslator.TranslationData{}

			// When
			result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

			// Then
			assert.NoError(t, err)
			require.NotNil(t, result)

			// Expected statement for both Create and Update operations (same for system_settings)
			expectedUpdateStatement := "update_configuration(" +
				"configuration_name=\"system_settings\", " +
				"administrator_grants_create_user=true, " +
				"administrator_grants_create_user_token=false, " +
				"administrator_grants_regenerate_user_token=true, " +
				"administrator_grants_read_user_token=false, " +
				"approval_feature_enabled=false, " +
				"runbook_ad_hoc_approval_request_enabled=true, " +
				"runbook_approval_request_expiry_time=30, " +
				"run_approval_expiry_time=45, " +
				"approval_editable_allowed_resource_query_enabled=true, " +
				"approval_allow_individual_notification=false, " +
				"approval_optional_request_ticket_url=true, " +
				"time_trigger_permissions_user=\"test_user\", " +
				"external_audit_storage_enabled=true, " +
				"external_audit_storage_type=\"SPLUNK\", " +
				"external_audit_storage_batch_period_sec=15, " +
				"environment_name=\"test_env\", " +
				"environment_name_background=\"#FF0000\", " +
				"param_value_max_length=5000, " +
				"parallel_runs_fired_by_time_triggers=3, " +
				"maintenance_mode_enabled=true, " +
				"managed_secrets=\"EXTERNAL\", " +
				"allowed_tags=[], " +
				"skipped_tags=[])"

			switch tt.operation {
			case common.Create, common.Update:
				assert.Equal(t, expectedUpdateStatement, result.Statement)
			case common.Read:
				assert.Equal(t, "get_configuration_class(configuration_name=\"system_settings\")", result.Statement)
			}
		})
	}
}

func TestSystemSettingsTranslatorCommon_ToAPIModel_WithStringArrays(t *testing.T) {
	// Given
	translator := &SystemSettingsTranslatorCommon{}
	tfModel := createTestSystemSettingsTFModelWithArrays()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Update).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	assert.NoError(t, err)
	require.NotNil(t, result)

	expectedStatement := "update_configuration(" +
		"configuration_name=\"system_settings\", " +
		"administrator_grants_create_user=true, " +
		"administrator_grants_create_user_token=true, " +
		"administrator_grants_regenerate_user_token=true, " +
		"administrator_grants_read_user_token=true, " +
		"approval_feature_enabled=true, " +
		"runbook_ad_hoc_approval_request_enabled=true, " +
		"runbook_approval_request_expiry_time=60, " +
		"run_approval_expiry_time=60, " +
		"approval_editable_allowed_resource_query_enabled=true, " +
		"approval_allow_individual_notification=true, " +
		"approval_optional_request_ticket_url=false, " +
		"time_trigger_permissions_user=\"Shoreline\", " +
		"external_audit_storage_enabled=false, " +
		"external_audit_storage_type=\"ELASTIC\", " +
		"external_audit_storage_batch_period_sec=5, " +
		"environment_name=\"\", " +
		"environment_name_background=\"#EF5350\", " +
		"param_value_max_length=10000, " +
		"parallel_runs_fired_by_time_triggers=10, " +
		"maintenance_mode_enabled=false, " +
		"managed_secrets=\"LOCAL\", " +
		"allowed_tags=[\"tag1\", \"tag2\"], " +
		"skipped_tags=[\"skip1\", \"skip2\"])"
	assert.Equal(t, expectedStatement, result.Statement)
}

func TestSystemSettingsTranslatorCommon_ToAPIModel_DeleteOperation(t *testing.T) {
	// Given
	translator := &SystemSettingsTranslatorCommon{}
	tfModel := createTestSystemSettingsTFModel()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Delete).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "delete operation is not supported for system_settings resource", err.Error())
}

func createTestSystemSettingsTFModel() *systemsettingstf.SystemSettingsTFModel {
	return &systemsettingstf.SystemSettingsTFModel{
		Name:                                        types.StringValue("system_settings"),
		AdministratorGrantsCreateUser:               types.BoolValue(true),
		AdministratorGrantsCreateUserToken:          types.BoolValue(false),
		AdministratorGrantsRegenerateUserToken:      types.BoolValue(true),
		AdministratorGrantsReadUserToken:            types.BoolValue(false),
		ApprovalFeatureEnabled:                      types.BoolValue(false),
		RunbookAdHocApprovalRequestEnabled:          types.BoolValue(true),
		RunbookApprovalRequestExpiryTime:            types.Int64Value(30),
		RunApprovalExpiryTime:                       types.Int64Value(45),
		ApprovalEditableAllowedResourceQueryEnabled: types.BoolValue(true),
		ApprovalAllowIndividualNotification:         types.BoolValue(false),
		ApprovalOptionalRequestTicketURL:            types.BoolValue(true),
		TimeTriggerPermissionsUser:                  types.StringValue("test_user"),
		ExternalAuditStorageEnabled:                 types.BoolValue(true),
		ExternalAuditStorageType:                    types.StringValue("SPLUNK"),
		ExternalAuditStorageBatchPeriodSec:          types.Int64Value(15),
		EnvironmentName:                             types.StringValue("test_env"),
		EnvironmentNameBackground:                   types.StringValue("#FF0000"),
		ParamValueMaxLength:                         types.Int64Value(5000),
		ParallelRunsFiredByTimeTriggers:             types.Int64Value(3),
		MaintenanceModeEnabled:                      types.BoolValue(true),
		ManagedSecrets:                              types.StringValue("EXTERNAL"),
		AllowedTags:                                 types.ListValueMust(types.StringType, []attr.Value{}),
		SkippedTags:                                 types.ListValueMust(types.StringType, []attr.Value{}),
	}
}

func createTestSystemSettingsTFModelWithArrays() *systemsettingstf.SystemSettingsTFModel {
	allowedTagsElems := []types.String{
		types.StringValue("tag1"),
		types.StringValue("tag2"),
	}

	skippedTagsElems := []types.String{
		types.StringValue("skip1"),
		types.StringValue("skip2"),
	}

	return &systemsettingstf.SystemSettingsTFModel{
		Name:                                        types.StringValue("system_settings"),
		AdministratorGrantsCreateUser:               types.BoolValue(true),
		AdministratorGrantsCreateUserToken:          types.BoolValue(true),
		AdministratorGrantsRegenerateUserToken:      types.BoolValue(true),
		AdministratorGrantsReadUserToken:            types.BoolValue(true),
		ApprovalFeatureEnabled:                      types.BoolValue(true),
		RunbookAdHocApprovalRequestEnabled:          types.BoolValue(true),
		RunbookApprovalRequestExpiryTime:            types.Int64Value(60),
		RunApprovalExpiryTime:                       types.Int64Value(60),
		ApprovalEditableAllowedResourceQueryEnabled: types.BoolValue(true),
		ApprovalAllowIndividualNotification:         types.BoolValue(true),
		ApprovalOptionalRequestTicketURL:            types.BoolValue(false),
		TimeTriggerPermissionsUser:                  types.StringValue("Shoreline"),
		ExternalAuditStorageEnabled:                 types.BoolValue(false),
		ExternalAuditStorageType:                    types.StringValue("ELASTIC"),
		ExternalAuditStorageBatchPeriodSec:          types.Int64Value(5),
		EnvironmentName:                             types.StringValue(""),
		EnvironmentNameBackground:                   types.StringValue("#EF5350"),
		ParamValueMaxLength:                         types.Int64Value(10000),
		ParallelRunsFiredByTimeTriggers:             types.Int64Value(10),
		MaintenanceModeEnabled:                      types.BoolValue(false),
		ManagedSecrets:                              types.StringValue("LOCAL"),
		AllowedTags:                                 types.ListValueMust(types.StringType, []attr.Value{allowedTagsElems[0], allowedTagsElems[1]}),
		SkippedTags:                                 types.ListValueMust(types.StringType, []attr.Value{skippedTagsElems[0], skippedTagsElems[1]}),
	}
}
