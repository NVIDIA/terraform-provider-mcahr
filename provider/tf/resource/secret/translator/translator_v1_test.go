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

func createNVaultSecretResponseV1(name, integrationName, vaultPath, vaultKey string) *secretapi.NVaultSecretResponseAPIModelV1 {
	return &secretapi.NVaultSecretResponseAPIModelV1{
		GetSecret: &secretapi.NVaultSecretContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			Secrets: []secretapi.NVaultSecretV1{
				{
					Name: name,
					SecretInfo: secretapi.NVaultSecretInfoV1{
						IntegrationName: integrationName,
						VaultSecretPath: vaultPath,
						VaultSecretKey:  vaultKey,
					},
				},
			},
		},
	}
}

func TestNVaultSecretTranslatorV1_ToTFModel_Success(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorV1{}
	apiModel := createNVaultSecretResponseV1("test_secret", "vault_integration", "path/to/secret", "username")
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
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

func TestNVaultSecretTranslatorV1_ToTFModel_DifferentValues(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorV1{}
	apiModel := createNVaultSecretResponseV1("minimal_secret", "minimal_vault", "simple/path", "key")
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal_secret", result.Name.ValueString())
	assert.Equal(t, "simple/path", result.VaultSecretPath.ValueString())
	assert.Equal(t, "key", result.VaultSecretKey.ValueString())
	assert.Equal(t, "minimal_vault", result.IntegrationName.ValueString())
}

func TestNVaultSecretTranslatorV1_ToTFModel_NilInput(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestNVaultSecretTranslatorV1_ToTFModel_NoContainer(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorV1{}
	apiModel := &secretapi.NVaultSecretResponseAPIModelV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no nvault secret container found")
}

func TestNVaultSecretTranslatorV1_ToTFModel_NoSecrets(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorV1{}
	apiModel := &secretapi.NVaultSecretResponseAPIModelV1{
		GetSecret: &secretapi.NVaultSecretContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			Secrets: []secretapi.NVaultSecretV1{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no nvault secrets found")
}
