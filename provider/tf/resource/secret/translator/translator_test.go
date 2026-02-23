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

	secretapi "terraform/terraform-provider/provider/external_api/resources/secrets"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createNVaultSecretResponseV2(name, integrationName, vaultPath, vaultKey string) *secretapi.NVaultSecretResponseAPIModel {
	return &secretapi.NVaultSecretResponseAPIModel{
		Output: secretapi.NVaultSecretOutput{
			Configurations: secretapi.NVaultSecretConfigurations{
				Items: []secretapi.NVaultConfigurationItem{
					{
						Config: secretapi.NVaultSecretConfig{
							ExternalValue: secretapi.NVaultSecretExternalValue{
								IntegrationName: integrationName,
								VaultSecretPath: vaultPath,
								VaultSecretKey:  vaultKey,
							},
						},
						EntityMetadata: secretapi.NVaultSecretEntityMetadata{
							Name: name,
						},
					},
				},
			},
		},
		Summary: secretapi.NVaultSecretSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}

func TestNVaultSecretTranslator_ToTFModel_Success(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslator{}
	apiModel := createNVaultSecretResponseV2("test_secret", "vault_integration", "path/to/secret", "username")
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_secret", result.Name.ValueString())
	assert.Equal(t, "path/to/secret", result.VaultSecretPath.ValueString())
	assert.Equal(t, "username", result.VaultSecretKey.ValueString())
	assert.Equal(t, "vault_integration", result.IntegrationName.ValueString())
}

func TestNVaultSecretTranslator_ToTFModel_NilInput(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslator{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestNVaultSecretTranslator_ToTFModel_NoConfigurations(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslator{}
	apiModel := &secretapi.NVaultSecretResponseAPIModel{
		Output: secretapi.NVaultSecretOutput{
			Configurations: secretapi.NVaultSecretConfigurations{
				Items: []secretapi.NVaultConfigurationItem{},
			},
		},
		Summary: secretapi.NVaultSecretSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no nvault secret configurations found")
}
