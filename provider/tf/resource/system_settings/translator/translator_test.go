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
	"terraform/terraform-provider/provider/common"
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	"testing"

	systemsettingsapi "terraform/terraform-provider/provider/external_api/resources/system_settings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemSettingsTranslator_ToTFModel_Success(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SystemSettingsTranslator{}
	apiModel := createFullSystemSettingsResponseV2()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "system_settings", result.Name.ValueString())
	assert.Equal(t, true, result.AdministratorGrantsCreateUser.ValueBool())
	assert.Equal(t, true, result.AdministratorGrantsCreateUserToken.ValueBool())
	assert.Equal(t, false, result.AdministratorGrantsRegenerateUserToken.ValueBool())
	assert.Equal(t, true, result.AdministratorGrantsReadUserToken.ValueBool())
	assert.Equal(t, true, result.ApprovalFeatureEnabled.ValueBool())
	assert.Equal(t, true, result.RunbookAdHocApprovalRequestEnabled.ValueBool())
	assert.Equal(t, int64(6), result.RunbookApprovalRequestExpiryTime.ValueInt64())
	assert.Equal(t, int64(5), result.RunApprovalExpiryTime.ValueInt64())
	assert.Equal(t, true, result.ApprovalEditableAllowedResourceQueryEnabled.ValueBool())
	assert.Equal(t, true, result.ApprovalAllowIndividualNotification.ValueBool())
	assert.Equal(t, false, result.ApprovalOptionalRequestTicketURL.ValueBool())
	assert.Equal(t, "test_user", result.TimeTriggerPermissionsUser.ValueString())
	assert.Equal(t, false, result.ExternalAuditStorageEnabled.ValueBool())
	assert.Equal(t, "ELASTIC", result.ExternalAuditStorageType.ValueString())
	assert.Equal(t, int64(10), result.ExternalAuditStorageBatchPeriodSec.ValueInt64())
	assert.Equal(t, "Env_Name via TF", result.EnvironmentName.ValueString())
	assert.Equal(t, "#673ab7", result.EnvironmentNameBackground.ValueString())
	assert.Equal(t, int64(10000), result.ParamValueMaxLength.ValueInt64())
	assert.Equal(t, int64(5), result.ParallelRunsFiredByTimeTriggers.ValueInt64())
	assert.Equal(t, false, result.MaintenanceModeEnabled.ValueBool())
	assert.Equal(t, "LOCAL", result.ManagedSecrets.ValueString())

	// Verify sets
	var allowedTags []string
	result.AllowedTags.ElementsAs(context.Background(), &allowedTags, false)
	assert.Equal(t, []string{"production", "staging", "development", "testing"}, allowedTags)

	var skippedTags []string
	result.SkippedTags.ElementsAs(context.Background(), &skippedTags, false)
	assert.Equal(t, []string{"deprecated", "internal", "temp"}, skippedTags)
}

func TestSystemSettingsTranslator_ToTFModel_NilInput(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SystemSettingsTranslator{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestSystemSettingsTranslator_ToTFModel_EmptyItems(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SystemSettingsTranslator{}
	apiModel := &systemsettingsapi.SystemSettingsResponseAPIModel{
		Output: systemsettingsapi.OutputV2{
			SystemConfigurations: systemsettingsapi.SystemConfigurationsV2{
				Count: 0,
				Items: []systemsettingsapi.SystemSettingsItemV2{},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no system_settings configurations found in API response")
}

func createFullSystemSettingsResponseV2() *systemsettingsapi.SystemSettingsResponseAPIModel {
	return &systemsettingsapi.SystemSettingsResponseAPIModel{
		Output: systemsettingsapi.OutputV2{
			SystemConfigurations: systemsettingsapi.SystemConfigurationsV2{
				Count: 1,
				Items: []systemsettingsapi.SystemSettingsItemV2{
					{
						Name: "system_settings",
						Configuration: systemsettingsapi.SystemSettingsConfigurationV2{
							AdministratorGrantsCreateUser:               true,
							AdministratorGrantsCreateUserToken:          true,
							AdministratorGrantsRegenerateUserToken:      false,
							AdministratorGrantsReadUserToken:            true,
							ApprovalFeatureEnabled:                      true,
							RunbookAdHocApprovalRequestEnabled:          true,
							RunbookApprovalRequestExpiryTime:            6,
							RunApprovalExpiryTime:                       5,
							ApprovalEditableAllowedResourceQueryEnabled: true,
							ApprovalAllowIndividualNotification:         true,
							ApprovalOptionalRequestTicketURL:            false,
							TimeTriggerPermissionsUser:                  "test_user",
							ExternalAuditStorageEnabled:                 false,
							ExternalAuditStorageType:                    "ELASTIC",
							ExternalAuditStorageBatchPeriodSec:          10,
							EnvironmentName:                             "Env_Name via TF",
							EnvironmentNameBackground:                   "#673ab7",
							ParamValueMaxLength:                         10000,
							ParallelRunsFiredByTimeTriggers:             5,
							MaintenanceModeEnabled:                      false,
							AllowedTags:                                 []string{"production", "staging", "development", "testing"},
							SkippedTags:                                 []string{"deprecated", "internal", "temp"},
							ManagedSecrets:                              "LOCAL",
						},
					},
				},
			},
		},
		Summary: systemsettingsapi.SummaryV2{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}
