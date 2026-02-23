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
	secrettf "terraform/terraform-provider/provider/tf/resource/secret/model"
)

// NVaultSecretTranslatorCommon provides common functionality for nvault secret translators across API versions
type NVaultSecretTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (t *NVaultSecretTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *secrettf.NVaultSecretTFModel) (*statement.StatementInputAPIModel, error) {
	var stmt string

	switch requestContext.Operation {
	case common.Create:
		stmt = t.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = t.buildReadStatement(tfModel)
	case common.Update:
		stmt = t.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt = t.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}

	apiModel := &statement.StatementInputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (t *NVaultSecretTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *secrettf.NVaultSecretTFModel) string {
	return t.buildSecretStatement(requestContext, translationData, "define_secret", tfModel)
}

func (t *NVaultSecretTranslatorCommon) buildReadStatement(tfModel *secrettf.NVaultSecretTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_secret(secret_name=\"%s\")", name)
}

func (t *NVaultSecretTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *secrettf.NVaultSecretTFModel) string {
	return t.buildSecretStatement(requestContext, translationData, "update_secret", tfModel)
}

func (t *NVaultSecretTranslatorCommon) buildDeleteStatement(tfModel *secrettf.NVaultSecretTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_secret(secret_name=\"%s\")", name)
}

func (t *NVaultSecretTranslatorCommon) buildSecretStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *secrettf.NVaultSecretTFModel) string {
	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions)

	// Add required fields
	name := tfModel.Name.ValueString()
	builder.SetStringField("secret_name", name, "name")

	// Build external_value as a JSON object for the secret configuration
	externalValue := fmt.Sprintf("{\"integration_name\":\"%s\",\"vault_secret_path\":\"%s\",\"vault_secret_key\":\"%s\"}",
		tfModel.IntegrationName.ValueString(),
		tfModel.VaultSecretPath.ValueString(),
		tfModel.VaultSecretKey.ValueString())

	builder.SetField("external_value", externalValue, "")

	return builder.Build()
}
