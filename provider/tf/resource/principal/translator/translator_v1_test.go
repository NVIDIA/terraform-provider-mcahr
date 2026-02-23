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

	principalapi "terraform/terraform-provider/provider/external_api/resources/principals"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a full V1 API response for testing
func createFullPrincipalResponseV1() *principalapi.PrincipalResponseAPIModelV1 {
	return &principalapi.PrincipalResponseAPIModelV1{
		GetPrincipalClass: &principalapi.PrincipalContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			PrincipalClasses: []principalapi.PrincipalClassV1{
				{
					Enabled:              false,
					Name:                 "test_principal",
					Deleted:              false,
					Identity:             "test@example.com",
					ActionLimit:          100,
					ExecuteLimit:         50,
					ConfigurePermission:  1,
					AdministerPermission: 1,
					IDPName:              "azure",
				},
			},
		},
	}
}

// Helper function to create a minimal V1 API response for testing
func createMinimalPrincipalResponseV1() *principalapi.PrincipalResponseAPIModelV1 {
	return &principalapi.PrincipalResponseAPIModelV1{
		GetPrincipalClass: &principalapi.PrincipalContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			PrincipalClasses: []principalapi.PrincipalClassV1{
				{
					Enabled:              false,
					Name:                 "minimal_principal",
					Deleted:              false,
					Identity:             "minimal@example.com",
					ActionLimit:          0,
					ExecuteLimit:         0,
					ConfigurePermission:  0,
					AdministerPermission: 0,
					IDPName:              "",
				},
			},
		},
	}
}

func TestPrincipalTranslatorV1_ToTFModel_Success(t *testing.T) {
	t.Parallel()

	// Given
	translator := &PrincipalTranslatorV1{}
	apiModel := createFullPrincipalResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_principal", result.Name.ValueString())
	assert.Equal(t, "test@example.com", result.Identity.ValueString())
	assert.Equal(t, int64(100), result.ActionLimit.ValueInt64())
	assert.Equal(t, int64(50), result.ExecuteLimit.ValueInt64())
	assert.True(t, result.ConfigurePermission.ValueBool())
	assert.True(t, result.AdministerPermission.ValueBool())
	assert.Equal(t, "azure", result.IDPName.ValueString())
	assert.True(t, result.ViewLimit.IsNull()) // ViewLimit is not set by translator, will be handled by post-processor
}

func TestPrincipalTranslatorV1_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &PrincipalTranslatorV1{}
	apiModel := createMinimalPrincipalResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify basic fields
	assert.Equal(t, "minimal_principal", result.Name.ValueString())
	assert.Equal(t, "minimal@example.com", result.Identity.ValueString())
	assert.Equal(t, int64(0), result.ActionLimit.ValueInt64())
	assert.Equal(t, int64(0), result.ExecuteLimit.ValueInt64())
	assert.False(t, result.ConfigurePermission.ValueBool())
	assert.False(t, result.AdministerPermission.ValueBool())

	// Verify that empty optional fields are empty strings or defaults
	assert.Equal(t, "", result.IDPName.ValueString())
	assert.True(t, result.ViewLimit.IsNull()) // ViewLimit is not set by translator, will be handled by post-processor
}

func TestPrincipalTranslatorV1_ToTFModel_EmptyPrincipalClasses(t *testing.T) {
	t.Parallel()

	// Given - V1 API response with empty principal classes (translator-level validation)
	apiModel := &principalapi.PrincipalResponseAPIModelV1{
		GetPrincipalClass: &principalapi.PrincipalContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			PrincipalClasses: []principalapi.PrincipalClassV1{}, // Empty list
		},
	}

	translator := &PrincipalTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	tfModel, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	require.Nil(t, tfModel)
	assert.Contains(t, err.Error(), "no principal classes found")
}

func TestPrincipalTranslatorV1_ToTFModel_NilInput(t *testing.T) {
	t.Parallel()

	// Given
	translator := &PrincipalTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}
