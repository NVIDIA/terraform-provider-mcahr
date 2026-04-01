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
	integrationapi "terraform/terraform-provider/provider/external_api/resources/integrations"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationTranslatorV1_ToTFModel(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslatorV1{}

	tests := []struct {
		name        string
		apiModel    *integrationapi.IntegrationResponseAPIModelV1
		expectError bool
		expectNil   bool
		validate    func(t *testing.T, tfModel *integrationtf.IntegrationTFModel)
	}{
		{
			name:        "Nil API model",
			apiModel:    nil,
			expectError: false,
			expectNil:   true,
		},
		{
			name: "Valid GetIntegrationClass response",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				GetIntegrationClass: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{
						{
							Name:            "okta-v1-test",
							ServiceName:     "okta",
							SerialNumber:    "OKTA001",
							Enabled:         true,
							PermissionsUser: "okta@company.com",
							Params:          `{"idp_name":"okta-idp","cache_ttl_ms":300000,"api_rate_limit":60,"api_token":"okta-token-123","url":"https://test.okta.com"}`,
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "okta-v1-test", tfModel.Name.ValueString())
				assert.Equal(t, "okta", tfModel.ServiceName.ValueString())
				assert.Equal(t, "OKTA001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "okta@company.com", tfModel.PermissionsUser.ValueString())

				// Check fields from JSON params (mapped through adapter)
				assert.Equal(t, "okta-idp", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(300000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(60), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "okta-token-123", tfModel.APIKey.ValueString())        // Maps from api_token
				assert.Equal(t, "https://test.okta.com", tfModel.APIUrl.ValueString()) // Maps from url
			},
		},
		{
			name: "Valid DefineIntegration response",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				DefineIntegration: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{
						{
							Name:            "alertmanager-v1-create",
							ServiceName:     "alertmanager",
							SerialNumber:    "AM001",
							Enabled:         false,
							PermissionsUser: "alertmanager@company.com",
							Params:          `{"external_url":"https://alertmanager.company.com","payload_paths":["alerts.receiver","alerts.status"]}`,
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "alertmanager-v1-create", tfModel.Name.ValueString())
				assert.Equal(t, "alertmanager", tfModel.ServiceName.ValueString())
				assert.Equal(t, "AM001", tfModel.SerialNumber.ValueString())
				assert.False(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "alertmanager@company.com", tfModel.PermissionsUser.ValueString())

				// Check Alertmanager-specific fields
				assert.Equal(t, "https://alertmanager.company.com", tfModel.ExternalUrl.ValueString())
				assert.False(t, tfModel.PayloadPaths.IsNull())
				assert.Len(t, tfModel.PayloadPaths.Elements(), 2)
			},
		},
		{
			name: "Valid UpdateIntegration response",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				UpdateIntegration: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{
						{
							Name:            "azure-ad-v1-update",
							ServiceName:     "azure_active_directory",
							SerialNumber:    "AAD002",
							Enabled:         true,
							PermissionsUser: "azuread-updated@company.com",
							Params:          `{"idp_name":"azure-idp-updated","cache_ttl_ms":600000,"api_rate_limit":120,"tenant_id":"updated-tenant","client_id":"updated-client","client_secret":"updated-secret"}`,
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "azure-ad-v1-update", tfModel.Name.ValueString())
				assert.Equal(t, "azure_active_directory", tfModel.ServiceName.ValueString())
				assert.Equal(t, "AAD002", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "azuread-updated@company.com", tfModel.PermissionsUser.ValueString())

				// Check Azure AD-specific fields
				assert.Equal(t, "azure-idp-updated", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(600000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(120), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "updated-tenant", tfModel.TenantID.ValueString())
				assert.Equal(t, "updated-client", tfModel.ClientID.ValueString())
				assert.Equal(t, "updated-secret", tfModel.ClientSecret.ValueString())
			},
		},
		{
			name: "Valid DeleteIntegration response",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				DeleteIntegration: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{
						{
							Name:            "deleted-integration",
							ServiceName:     "okta",
							SerialNumber:    "DEL001",
							Enabled:         false,
							PermissionsUser: "deleted@company.com",
							Params:          `{"api_token":"deleted-token","url":"https://deleted.okta.com"}`,
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "deleted-integration", tfModel.Name.ValueString())
				assert.Equal(t, "okta", tfModel.ServiceName.ValueString())
				assert.Equal(t, "DEL001", tfModel.SerialNumber.ValueString())
				assert.False(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "deleted@company.com", tfModel.PermissionsUser.ValueString())
			},
		},
		{
			name: "Empty integration classes",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				GetIntegrationClass: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{},
				},
			},
			expectError: true,
		},
		{
			name:     "No container found",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				// All containers are nil
			},
			expectError: true,
		},
		{
			name: "Multiple integration classes - should return first",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				GetIntegrationClass: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{
						{
							Name:            "first-v1-integration",
							ServiceName:     "bcm",
							SerialNumber:    "FIRST001",
							Enabled:         true,
							PermissionsUser: "first@company.com",
							Params:          `{"idp_name":"first-bcm","cache_ttl_ms":180000,"api_rate_limit":50}`,
						},
						{
							Name:            "second-v1-integration",
							ServiceName:     "nvault",
							SerialNumber:    "SECOND001",
							Enabled:         false,
							PermissionsUser: "second@company.com",
							Params:          `{"address":"https://second.vault.com","role_name":"second-role"}`,
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				// Should return the first integration class
				assert.Equal(t, "first-v1-integration", tfModel.Name.ValueString())
				assert.Equal(t, "bcm", tfModel.ServiceName.ValueString())
				assert.Equal(t, "FIRST001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "first@company.com", tfModel.PermissionsUser.ValueString())

				// Check BCM-specific fields
				assert.Equal(t, "first-bcm", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(180000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(50), tfModel.APIRateLimit.ValueInt64())
			},
		},
		{
			name: "Integration with complex JSON params",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				GetIntegrationClass: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{
						{
							Name:            "complex-v1-integration",
							ServiceName:     "google_cloud_identity",
							SerialNumber:    "COMPLEX001",
							Enabled:         true,
							PermissionsUser: "complex@company.com",
							Params:          `{"idp_name":"google-complex","cache_ttl_ms":900000,"api_rate_limit":200,"subject":"complex-service@project.iam.gserviceaccount.com","credentials":"{\"type\":\"service_account\",\"project_id\":\"complex-project\"}"}`,
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "complex-v1-integration", tfModel.Name.ValueString())
				assert.Equal(t, "google_cloud_identity", tfModel.ServiceName.ValueString())
				assert.Equal(t, "COMPLEX001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "complex@company.com", tfModel.PermissionsUser.ValueString())

				// Check Google Cloud Identity-specific fields
				assert.Equal(t, "google-complex", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(900000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(200), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "complex-service@project.iam.gserviceaccount.com", tfModel.Subject.ValueString())
				assert.Equal(t, `{"type":"service_account","project_id":"complex-project"}`, tfModel.Credentials.ValueString())
			},
		},
		{
			name: "Integration with empty/null params",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				GetIntegrationClass: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{
						{
							Name:            "empty-params-integration",
							ServiceName:     "alertmanager",
							SerialNumber:    "EMPTY001",
							Enabled:         false,
							PermissionsUser: "empty@company.com",
							Params:          `{}`, // Empty JSON
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "empty-params-integration", tfModel.Name.ValueString())
				assert.Equal(t, "alertmanager", tfModel.ServiceName.ValueString())
				assert.Equal(t, "EMPTY001", tfModel.SerialNumber.ValueString())
				assert.False(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "empty@company.com", tfModel.PermissionsUser.ValueString())

				// Fields should have default/empty values
				assert.Equal(t, "", tfModel.ExternalUrl.ValueString())
			},
		},
		{
			name: "Integration with invalid JSON params",
			apiModel: &integrationapi.IntegrationResponseAPIModelV1{
				GetIntegrationClass: &integrationapi.IntegrationContainerV1{
					IntegrationClasses: []integrationapi.IntegrationClassV1{
						{
							Name:            "invalid-json-integration",
							ServiceName:     "bcm_connectivity",
							SerialNumber:    "INVALID001",
							Enabled:         true,
							PermissionsUser: "invalid@company.com",
							Params:          `{invalid json}`, // Invalid JSON
						},
					},
				},
			},
			expectError: false, // Should not error, but fields will have default values
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "invalid-json-integration", tfModel.Name.ValueString())
				assert.Equal(t, "bcm_connectivity", tfModel.ServiceName.ValueString())
				assert.Equal(t, "INVALID001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "invalid@company.com", tfModel.PermissionsUser.ValueString())

				// Fields should have default/empty values due to invalid JSON
				assert.Equal(t, "", tfModel.APIKey.ValueString())
				assert.Equal(t, "", tfModel.APICertificate.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
			translationData := &coretranslator.TranslationData{}

			// When
			result, err := translator.ToTFModel(requestContext, translationData, tt.apiModel)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else if tt.expectNil {
				assert.NoError(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				tt.validate(t, result)
			}
		})
	}
}

func TestIntegrationTranslatorV1_EdgeCases(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslatorV1{}

	t.Run("Null JSON params", func(t *testing.T) {
		// Given
		apiModel := &integrationapi.IntegrationResponseAPIModelV1{
			GetIntegrationClass: &integrationapi.IntegrationContainerV1{
				IntegrationClasses: []integrationapi.IntegrationClassV1{
					{
						Name:            "null-params",
						ServiceName:     "okta",
						SerialNumber:    "NULL001",
						Enabled:         true,
						PermissionsUser: "null@company.com",
						Params:          "null", // JSON null
					},
				},
			},
		}
		requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
		translationData := &coretranslator.TranslationData{}

		// When
		result, err := translator.ToTFModel(requestContext, translationData, apiModel)

		// Then
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "null-params", result.Name.ValueString())
		assert.Equal(t, "okta", result.ServiceName.ValueString())
		// Integration fields should have default values
		assert.Equal(t, "", result.APIKey.ValueString())
	})

	t.Run("Complex nested JSON in params", func(t *testing.T) {
		// Given
		apiModel := &integrationapi.IntegrationResponseAPIModelV1{
			GetIntegrationClass: &integrationapi.IntegrationContainerV1{
				IntegrationClasses: []integrationapi.IntegrationClassV1{
					{
						Name:            "nested-json",
						ServiceName:     "google_cloud_identity",
						SerialNumber:    "NESTED001",
						Enabled:         true,
						PermissionsUser: "nested@company.com",
						Params:          `{"credentials":"{\"type\":\"service_account\",\"project_id\":\"nested-project\",\"private_key\":\"-----BEGIN PRIVATE KEY-----\\nMIIC...\\n-----END PRIVATE KEY-----\\n\"}","subject":"nested@project.iam.gserviceaccount.com"}`,
					},
				},
			},
		}
		requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
		translationData := &coretranslator.TranslationData{}

		// When
		result, err := translator.ToTFModel(requestContext, translationData, apiModel)

		// Then
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "nested-json", result.Name.ValueString())
		assert.Equal(t, "google_cloud_identity", result.ServiceName.ValueString())
		assert.Contains(t, result.Credentials.ValueString(), "service_account")
		assert.Contains(t, result.Credentials.ValueString(), "nested-project")
		assert.Equal(t, "nested@project.iam.gserviceaccount.com", result.Subject.ValueString())
	})

	t.Run("V1 with null TF model values", func(t *testing.T) {
		// Given
		tfModel := &integrationtf.IntegrationTFModel{
			Name:            types.StringValue("null-tf-values"),
			ServiceName:     types.StringValue("alertmanager"),
			SerialNumber:    types.StringNull(),
			Enabled:         types.BoolNull(),
			PermissionsUser: types.StringNull(),
			ExternalUrl:     types.StringNull(),
		}
		requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
		translationData := &coretranslator.TranslationData{}

		// When
		result, err := translator.ToAPIModel(requestContext, translationData, tfModel)

		// Then
		require.NoError(t, err)
		expected := `define_integration(integration_name="null-tf-values", serial_number="", enabled=false, permissions_user="", params={"external_url":"","payload_paths":[]}, service_name="alertmanager")`
		assert.Equal(t, expected, result.Statement)
		assert.Equal(t, common.V1, result.APIVersion)
	})
}
