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
	secrettf "terraform/terraform-provider/provider/tf/resource/secret/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNVaultSecretTranslatorCommon_ToAPIModel(t *testing.T) {
	tests := []struct {
		name      string
		operation common.CrudOperation
		expected  string
	}{
		{
			name:      "Create operation",
			operation: common.Create,
			expected: "define_secret(" +
				"secret_name=\"test_secret\", " +
				"external_value={\"integration_name\":\"vault_integration\",\"vault_secret_path\":\"path/to/secret\",\"vault_secret_key\":\"username\"})",
		},
		{
			name:      "Read operation",
			operation: common.Read,
			expected:  "get_secret(secret_name=\"test_secret\")",
		},
		{
			name:      "Update operation",
			operation: common.Update,
			expected: "update_secret(" +
				"secret_name=\"test_secret\", " +
				"external_value={\"integration_name\":\"vault_integration\",\"vault_secret_path\":\"path/to/secret\",\"vault_secret_key\":\"username\"})",
		},
		{
			name:      "Delete operation",
			operation: common.Delete,
			expected:  "delete_secret(secret_name=\"test_secret\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			translator := &NVaultSecretTranslatorCommon{}
			tfModel := &secrettf.NVaultSecretTFModel{
				Name:            types.StringValue("test_secret"),
				VaultSecretPath: types.StringValue("path/to/secret"),
				VaultSecretKey:  types.StringValue("username"),
				IntegrationName: types.StringValue("vault_integration"),
			}
			requestContext := common.NewRequestContext(context.Background()).WithOperation(tt.operation).WithAPIVersion(common.V1)
			translationData := &coretranslator.TranslationData{}

			// When
			result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

			// Then
			assert.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Statement)
		})
	}
}

func TestNVaultSecretTranslatorCommon_ToAPIModel_UnsupportedOperation(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorCommon{}
	tfModel := &secrettf.NVaultSecretTFModel{
		Name: types.StringValue("test_secret"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.CrudOperation(999)).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}

func TestNVaultSecretTranslatorCommon_BuildSecretStatement(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorCommon{}
	tfModel := &secrettf.NVaultSecretTFModel{
		Name:            types.StringValue("my_secret"),
		VaultSecretPath: types.StringValue("secrets/app/db"),
		VaultSecretKey:  types.StringValue("password"),
		IntegrationName: types.StringValue("vault_prod"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)

	// When
	result := translator.buildSecretStatement(requestContext, &coretranslator.TranslationData{}, "define_secret", tfModel)

	// Then
	expected := "define_secret(" +
		"secret_name=\"my_secret\", " +
		"external_value={\"integration_name\":\"vault_prod\",\"vault_secret_path\":\"secrets/app/db\",\"vault_secret_key\":\"password\"})"
	assert.Equal(t, expected, result)
}

func TestNVaultSecretTranslatorCommon_BuildReadStatement(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorCommon{}
	tfModel := &secrettf.NVaultSecretTFModel{
		Name: types.StringValue("test_secret"),
	}

	// When
	result := translator.buildReadStatement(tfModel)

	// Then
	assert.Equal(t, "get_secret(secret_name=\"test_secret\")", result)
}

func TestNVaultSecretTranslatorCommon_BuildDeleteStatement(t *testing.T) {
	// Given
	translator := &NVaultSecretTranslatorCommon{}
	tfModel := &secrettf.NVaultSecretTFModel{
		Name: types.StringValue("test_secret"),
	}

	// When
	result := translator.buildDeleteStatement(tfModel)

	// Then
	assert.Equal(t, "delete_secret(secret_name=\"test_secret\")", result)
}
