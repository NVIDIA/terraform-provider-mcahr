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

func TestPrincipalTranslator_ToTFModel_Success(t *testing.T) {
	// Given
	translator := &PrincipalTranslator{}
	apiModel := createFullPrincipalResponseV2()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
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

func TestPrincipalTranslator_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &PrincipalTranslator{}
	apiModel := createMinimalPrincipalResponseV2()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
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

func TestPrincipalTranslator_ToTFModel_EmptyAccessControlItems(t *testing.T) {
	// Given
	translator := &PrincipalTranslator{}
	apiModel := &principalapi.PrincipalResponseAPIModel{
		Output: principalapi.PrincipalOutput{
			AccessControl: principalapi.AccessControl{
				Items: []principalapi.AccessControlItem{}, // Empty list
			},
		},
		Summary: principalapi.PrincipalSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "No principal access control items found")
}

func TestPrincipalTranslator_ToTFModel_NilInput(t *testing.T) {
	// Given
	translator := &PrincipalTranslator{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// Helper function to create a full V2 API response for testing
func createFullPrincipalResponseV2() *principalapi.PrincipalResponseAPIModel {
	return &principalapi.PrincipalResponseAPIModel{
		Output: principalapi.PrincipalOutput{
			AccessControl: principalapi.AccessControl{
				Items: []principalapi.AccessControlItem{
					{
						Data: principalapi.PrincipalData{
							Name:                 "test_principal",
							Identity:             "test@example.com",
							IDPName:              "azure",
							ActionLimit:          100,
							ExecuteLimit:         50,
							ConfigurePermission:  1,
							AdministerPermission: 1,
						},
					},
				},
			},
		},
		Summary: principalapi.PrincipalSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}

// Helper function to create a minimal V2 API response for testing
func createMinimalPrincipalResponseV2() *principalapi.PrincipalResponseAPIModel {
	return &principalapi.PrincipalResponseAPIModel{
		Output: principalapi.PrincipalOutput{
			AccessControl: principalapi.AccessControl{
				Items: []principalapi.AccessControlItem{
					{
						Data: principalapi.PrincipalData{
							Name:                 "minimal_principal",
							Identity:             "minimal@example.com",
							IDPName:              "",
							ActionLimit:          0,
							ExecuteLimit:         0,
							ConfigurePermission:  0,
							AdministerPermission: 0,
						},
					},
				},
			},
		},
		Summary: principalapi.PrincipalSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}
