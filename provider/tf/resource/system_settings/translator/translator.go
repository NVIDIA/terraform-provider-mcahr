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

type SystemSettingsTranslator struct {
	SystemSettingsTranslatorCommon
}

// ToAPIModel converts a TF model to an API input model for V2
func (t *SystemSettingsTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *systemsettingstf.SystemSettingsTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}

func (t *SystemSettingsTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *systemsettingsapi.SystemSettingsResponseAPIModel) (*systemsettingstf.SystemSettingsTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.SystemConfigurations.Items) == 0 {
		return nil, fmt.Errorf("no system_settings configurations found in API response")
	}

	// Get the first configuration item, current implementation only supports one system_settings to be returned by the API
	configItem := apiModel.Output.SystemConfigurations.Items[0]
	config := configItem.Configuration

	// Build TF model from V2 system_settings configuration
	tfModel := &systemsettingstf.SystemSettingsTFModel{
		Name:                                        types.StringValue(configItem.Name),
		AdministratorGrantsCreateUser:               types.BoolValue(config.AdministratorGrantsCreateUser),
		AdministratorGrantsCreateUserToken:          types.BoolValue(config.AdministratorGrantsCreateUserToken),
		AdministratorGrantsRegenerateUserToken:      types.BoolValue(config.AdministratorGrantsRegenerateUserToken),
		AdministratorGrantsReadUserToken:            types.BoolValue(config.AdministratorGrantsReadUserToken),
		ApprovalFeatureEnabled:                      types.BoolValue(config.ApprovalFeatureEnabled),
		RunbookAdHocApprovalRequestEnabled:          types.BoolValue(config.RunbookAdHocApprovalRequestEnabled),
		RunbookApprovalRequestExpiryTime:            types.Int64Value(config.RunbookApprovalRequestExpiryTime),
		RunApprovalExpiryTime:                       types.Int64Value(config.RunApprovalExpiryTime),
		ApprovalEditableAllowedResourceQueryEnabled: types.BoolValue(config.ApprovalEditableAllowedResourceQueryEnabled),
		ApprovalAllowIndividualNotification:         types.BoolValue(config.ApprovalAllowIndividualNotification),
		ApprovalOptionalRequestTicketURL:            types.BoolValue(config.ApprovalOptionalRequestTicketURL),
		TimeTriggerPermissionsUser:                  types.StringValue(config.TimeTriggerPermissionsUser),
		ExternalAuditStorageEnabled:                 types.BoolValue(config.ExternalAuditStorageEnabled),
		ExternalAuditStorageType:                    types.StringValue(config.ExternalAuditStorageType),
		ExternalAuditStorageBatchPeriodSec:          types.Int64Value(config.ExternalAuditStorageBatchPeriodSec),
		EnvironmentName:                             types.StringValue(config.EnvironmentName),
		EnvironmentNameBackground:                   types.StringValue(config.EnvironmentNameBackground),
		ParamValueMaxLength:                         types.Int64Value(config.ParamValueMaxLength),
		ParallelRunsFiredByTimeTriggers:             types.Int64Value(config.ParallelRunsFiredByTimeTriggers),
		MaintenanceModeEnabled:                      types.BoolValue(config.MaintenanceModeEnabled),
		ManagedSecrets:                              types.StringValue(config.ManagedSecrets),
	}

	// Convert string arrays to sets
	tfModel.AllowedTags, _ = types.ListValueFrom(requestContext.Context, types.StringType, config.AllowedTags)
	tfModel.SkippedTags, _ = types.ListValueFrom(requestContext.Context, types.StringType, config.SkippedTags)

	return tfModel, nil
}
