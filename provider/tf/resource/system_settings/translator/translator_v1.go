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
	systemsettingsapi "terraform/terraform-provider/provider/external_api/resources/system_settings"
	"terraform/terraform-provider/provider/tf/core/translator"
	systemsettingstf "terraform/terraform-provider/provider/tf/resource/system_settings/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SystemSettingsTranslatorV1 handles translation between TF models and V1 API models for system_settings resources
type SystemSettingsTranslatorV1 struct {
	SystemSettingsTranslatorCommon
}

var _ translator.Translator[*systemsettingstf.SystemSettingsTFModel, *systemsettingsapi.SystemSettingsResponseAPIModelV1] = &SystemSettingsTranslatorV1{}

// ToAPIModel converts a TF model to an API input model for V1
func (t *SystemSettingsTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *systemsettingstf.SystemSettingsTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}

// ToTFModel converts a V1 API model to a TF model
func (t *SystemSettingsTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *systemsettingsapi.SystemSettingsResponseAPIModelV1) (*systemsettingstf.SystemSettingsTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the system_settings container
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no system_settings container found in V1 API response")
	}

	// Build TF model from V1 system_settings
	tfModel := &systemsettingstf.SystemSettingsTFModel{
		Name:                                        types.StringValue(apiModel.ConfigurationName),
		AdministratorGrantsCreateUser:               types.BoolValue(container.SystemSettings.AdministratorGrantsCreateUser),
		AdministratorGrantsCreateUserToken:          types.BoolValue(container.SystemSettings.AdministratorGrantsCreateUserToken),
		AdministratorGrantsRegenerateUserToken:      types.BoolValue(container.SystemSettings.AdministratorGrantsRegenerateUserToken),
		AdministratorGrantsReadUserToken:            types.BoolValue(container.SystemSettings.AdministratorGrantsReadUserToken),
		ApprovalFeatureEnabled:                      types.BoolValue(container.SystemSettings.ApprovalFeatureEnabled),
		RunbookAdHocApprovalRequestEnabled:          types.BoolValue(container.SystemSettings.RunbookAdHocApprovalRequestEnabled),
		RunbookApprovalRequestExpiryTime:            types.Int64Value(container.SystemSettings.RunbookApprovalRequestExpiryTime),
		RunApprovalExpiryTime:                       types.Int64Value(container.SystemSettings.RunApprovalExpiryTime),
		ApprovalEditableAllowedResourceQueryEnabled: types.BoolValue(container.SystemSettings.ApprovalEditableAllowedResourceQueryEnabled),
		ApprovalAllowIndividualNotification:         types.BoolValue(container.SystemSettings.ApprovalAllowIndividualNotification),
		ApprovalOptionalRequestTicketURL:            types.BoolValue(container.SystemSettings.ApprovalOptionalRequestTicketURL),
		TimeTriggerPermissionsUser:                  types.StringValue(container.SystemSettings.TimeTriggerPermissionsUser),
		ExternalAuditStorageEnabled:                 types.BoolValue(container.SystemSettings.ExternalAuditStorageEnabled),
		ExternalAuditStorageType:                    types.StringValue(container.SystemSettings.ExternalAuditStorageType),
		ExternalAuditStorageBatchPeriodSec:          types.Int64Value(container.SystemSettings.ExternalAuditStorageBatchPeriodSec),
		EnvironmentName:                             types.StringValue(container.SystemSettings.EnvironmentName),
		EnvironmentNameBackground:                   types.StringValue(container.SystemSettings.EnvironmentNameBackground),
		ParamValueMaxLength:                         types.Int64Value(container.SystemSettings.ParamValueMaxLength),
		ParallelRunsFiredByTimeTriggers:             types.Int64Value(container.SystemSettings.ParallelRunsFiredByTimeTriggers),
		MaintenanceModeEnabled:                      types.BoolValue(container.SystemSettings.MaintenanceModeEnabled),
		ManagedSecrets:                              types.StringValue(container.SystemSettings.ManagedSecrets),
	}

	// Convert string arrays to sets
	tfModel.AllowedTags, _ = types.ListValueFrom(requestContext.Context, types.StringType, container.SystemSettings.AllowedTags)
	tfModel.SkippedTags, _ = types.ListValueFrom(requestContext.Context, types.StringType, container.SystemSettings.SkippedTags)

	return tfModel, nil
}
