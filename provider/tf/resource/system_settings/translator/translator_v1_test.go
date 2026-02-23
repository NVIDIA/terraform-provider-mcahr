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

func TestSystemSettingsTranslatorV1_ToTFModel_Success(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SystemSettingsTranslatorV1{}
	apiModel := createFullSystemSettingsResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "system_settings", result.Name.ValueString())
	assert.Equal(t, true, result.AdministratorGrantsCreateUser.ValueBool())
	assert.Equal(t, true, result.AdministratorGrantsCreateUserToken.ValueBool())
	assert.Equal(t, true, result.AdministratorGrantsRegenerateUserToken.ValueBool())
	assert.Equal(t, true, result.AdministratorGrantsReadUserToken.ValueBool())
	assert.Equal(t, true, result.ApprovalFeatureEnabled.ValueBool())
	assert.Equal(t, true, result.RunbookAdHocApprovalRequestEnabled.ValueBool())
	assert.Equal(t, int64(60), result.RunbookApprovalRequestExpiryTime.ValueInt64())
	assert.Equal(t, int64(60), result.RunApprovalExpiryTime.ValueInt64())
	assert.Equal(t, true, result.ApprovalEditableAllowedResourceQueryEnabled.ValueBool())
	assert.Equal(t, true, result.ApprovalAllowIndividualNotification.ValueBool())
	assert.Equal(t, false, result.ApprovalOptionalRequestTicketURL.ValueBool())
	assert.Equal(t, "Shoreline", result.TimeTriggerPermissionsUser.ValueString())
	assert.Equal(t, false, result.ExternalAuditStorageEnabled.ValueBool())
	assert.Equal(t, "ELASTIC", result.ExternalAuditStorageType.ValueString())
	assert.Equal(t, int64(5), result.ExternalAuditStorageBatchPeriodSec.ValueInt64())
	assert.Equal(t, "", result.EnvironmentName.ValueString())
	assert.Equal(t, "#EF5350", result.EnvironmentNameBackground.ValueString())
	assert.Equal(t, int64(100000), result.ParamValueMaxLength.ValueInt64())
	assert.Equal(t, int64(10), result.ParallelRunsFiredByTimeTriggers.ValueInt64())
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

func TestSystemSettingsTranslatorV1_ToTFModel_NilInput(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SystemSettingsTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestSystemSettingsTranslatorV1_ToTFModel_NilContainer(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SystemSettingsTranslatorV1{}
	apiModel := &systemsettingsapi.SystemSettingsResponseAPIModelV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no system_settings container found in V1 API response")
}

func createFullSystemSettingsResponseV1() *systemsettingsapi.SystemSettingsResponseAPIModelV1 {
	return &systemsettingsapi.SystemSettingsResponseAPIModelV1{
		ConfigurationName: "system_settings",
		GetConfigurationClass: &systemsettingsapi.SystemSettingsContainerV1{
			Error: apicommon.ErrorV1{
				Message:          "",
				Type:             "OK",
				ValidationErrors: []apicommon.ValidationError{},
			},
			SystemSettings: systemsettingsapi.SystemSettingsV1{
				AdministratorGrantsCreateUser:               true,
				AdministratorGrantsCreateUserToken:          true,
				AdministratorGrantsRegenerateUserToken:      true,
				AdministratorGrantsReadUserToken:            true,
				ApprovalFeatureEnabled:                      true,
				RunbookAdHocApprovalRequestEnabled:          true,
				RunbookApprovalRequestExpiryTime:            60,
				RunApprovalExpiryTime:                       60,
				ApprovalEditableAllowedResourceQueryEnabled: true,
				ApprovalAllowIndividualNotification:         true,
				ApprovalOptionalRequestTicketURL:            false,
				TimeTriggerPermissionsUser:                  "Shoreline",
				ExternalAuditStorageEnabled:                 false,
				ExternalAuditStorageType:                    "ELASTIC",
				ExternalAuditStorageBatchPeriodSec:          5,
				EnvironmentName:                             "",
				EnvironmentNameBackground:                   "#EF5350",
				ParamValueMaxLength:                         100000,
				ParallelRunsFiredByTimeTriggers:             10,
				MaintenanceModeEnabled:                      false,
				AllowedTags:                                 []string{"production", "staging", "development", "testing"},
				SkippedTags:                                 []string{"deprecated", "internal", "temp"},
				ManagedSecrets:                              "LOCAL",
			},
		},
	}
}
