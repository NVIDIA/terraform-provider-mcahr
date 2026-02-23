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

func TestIntegrationTranslator_ToTFModel(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslator{}

	tests := []struct {
		name        string
		apiModel    *integrationapi.IntegrationResponseAPIModel
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
			name: "Empty integrations count",
			apiModel: &integrationapi.IntegrationResponseAPIModel{
				Output: integrationapi.IntegrationOutput{
					Integrations: integrationapi.IntegrationConfigurations{
						Items: []integrationapi.IntegrationItem{},
					},
				},
			},
			expectError: true,
		},
		{
			name: "Valid Datadog integration",
			apiModel: &integrationapi.IntegrationResponseAPIModel{
				Output: integrationapi.IntegrationOutput{
					Integrations: integrationapi.IntegrationConfigurations{
						Items: []integrationapi.IntegrationItem{
							{
								Name:            "datadog-prod",
								Enabled:         true,
								SerialNumber:    "DD001",
								PermissionsUser: "datadog@company.com",
								IntegrationType: "DATADOG",
								IntegrationData: map[string]interface{}{
									"api_key":      "dd-api-key-123",
									"api_url":      "https://api.datadoghq.com",
									"site_url":     "https://app.datadoghq.com",
									"app_key":      "dd-app-key-456",
									"webhook_name": "prod-webhook",
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "datadog-prod", tfModel.Name.ValueString())
				assert.Equal(t, "datadog", tfModel.ServiceName.ValueString()) // Should be lowercase
				assert.Equal(t, "DD001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "datadog@company.com", tfModel.PermissionsUser.ValueString())

				// Check integration-specific fields
				assert.Equal(t, "dd-api-key-123", tfModel.APIKey.ValueString())
				assert.Equal(t, "https://api.datadoghq.com", tfModel.APIUrl.ValueString())
				assert.Equal(t, "https://app.datadoghq.com", tfModel.SiteUrl.ValueString())
				assert.Equal(t, "dd-app-key-456", tfModel.AppKey.ValueString())
				assert.Equal(t, "prod-webhook", tfModel.WebhookName.ValueString())

				// PayloadPaths should be null since it's not set for Datadog
				assert.True(t, tfModel.PayloadPaths.IsNull())
			},
		},
		{
			name: "Valid Alertmanager integration",
			apiModel: &integrationapi.IntegrationResponseAPIModel{
				Output: integrationapi.IntegrationOutput{
					Integrations: integrationapi.IntegrationConfigurations{
						Items: []integrationapi.IntegrationItem{
							{
								Name:            "alertmanager-test",
								Enabled:         false,
								SerialNumber:    "AM001",
								PermissionsUser: "alerts@company.com",
								IntegrationType: "ALERTMANAGER",
								IntegrationData: map[string]interface{}{
									"external_url":  "https://alertmanager.company.com",
									"payload_paths": []interface{}{"alerts.receiver", "alerts.status", "alerts.labels"},
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "alertmanager-test", tfModel.Name.ValueString())
				assert.Equal(t, "alertmanager", tfModel.ServiceName.ValueString())
				assert.Equal(t, "AM001", tfModel.SerialNumber.ValueString())
				assert.False(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "alerts@company.com", tfModel.PermissionsUser.ValueString())

				// Check Alertmanager-specific fields
				assert.Equal(t, "https://alertmanager.company.com", tfModel.ExternalUrl.ValueString())
				assert.False(t, tfModel.PayloadPaths.IsNull())
				assert.Len(t, tfModel.PayloadPaths.Elements(), 3)
			},
		},
		{
			name: "Valid Azure Active Directory integration",
			apiModel: &integrationapi.IntegrationResponseAPIModel{
				Output: integrationapi.IntegrationOutput{
					Integrations: integrationapi.IntegrationConfigurations{
						Items: []integrationapi.IntegrationItem{
							{
								Name:            "azure-ad-prod",
								Enabled:         true,
								SerialNumber:    "AAD001",
								PermissionsUser: "azuread@company.com",
								IntegrationType: "Azure_Active_Directory",
								IntegrationData: map[string]interface{}{
									"idp_name":       "azure-idp",
									"cache_ttl_ms":   600000,
									"api_rate_limit": 100,
									"tenant_id":      "12345678-1234-1234-1234-123456789abc",
									"client_id":      "87654321-4321-4321-4321-cba987654321",
									"client_secret":  "super-secret-value",
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "azure-ad-prod", tfModel.Name.ValueString())
				assert.Equal(t, "azure_active_directory", tfModel.ServiceName.ValueString())
				assert.Equal(t, "AAD001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "azuread@company.com", tfModel.PermissionsUser.ValueString())

				// Check Azure AD-specific fields
				assert.Equal(t, "azure-idp", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(600000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(100), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "12345678-1234-1234-1234-123456789abc", tfModel.TenantID.ValueString())
				assert.Equal(t, "87654321-4321-4321-4321-cba987654321", tfModel.ClientID.ValueString())
				assert.Equal(t, "super-secret-value", tfModel.ClientSecret.ValueString())
			},
		},
		{
			name: "Valid NVault integration",
			apiModel: &integrationapi.IntegrationResponseAPIModel{
				Output: integrationapi.IntegrationOutput{
					Integrations: integrationapi.IntegrationConfigurations{
						Items: []integrationapi.IntegrationItem{
							{
								Name:            "vault-prod",
								Enabled:         true,
								SerialNumber:    "VAULT001",
								PermissionsUser: "vault@company.com",
								IntegrationType: "NVAULT",
								IntegrationData: map[string]interface{}{
									"address":       "https://vault.company.com:8200",
									"namespace":     "prod/team1",
									"role_name":     "terraform-role",
									"jwt_auth_path": "auth/jwt-prod",
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "vault-prod", tfModel.Name.ValueString())
				assert.Equal(t, "nvault", tfModel.ServiceName.ValueString())
				assert.Equal(t, "VAULT001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "vault@company.com", tfModel.PermissionsUser.ValueString())

				// Check NVault-specific fields
				assert.Equal(t, "https://vault.company.com:8200", tfModel.Address.ValueString())
				assert.Equal(t, "prod/team1", tfModel.Namespace.ValueString())
				assert.Equal(t, "terraform-role", tfModel.RoleName.ValueString())
				assert.Equal(t, "auth/jwt-prod", tfModel.JWTAuthPath.ValueString())
			},
		},
		{
			name: "Multiple integrations - should return first",
			apiModel: &integrationapi.IntegrationResponseAPIModel{
				Output: integrationapi.IntegrationOutput{
					Integrations: integrationapi.IntegrationConfigurations{
						Items: []integrationapi.IntegrationItem{
							{
								Name:            "first-integration",
								Enabled:         true,
								SerialNumber:    "FIRST001",
								PermissionsUser: "first@company.com",
								IntegrationType: "ELASTIC",
								IntegrationData: map[string]interface{}{
									"api_token": "first-token",
									"url":       "https://first.elastic.co",
								},
							},
							{
								Name:            "second-integration",
								Enabled:         false,
								SerialNumber:    "SECOND001",
								PermissionsUser: "second@company.com",
								IntegrationType: "OKTA",
								IntegrationData: map[string]interface{}{
									"api_token": "second-token",
									"url":       "https://second.okta.com",
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				// Should return the first integration
				assert.Equal(t, "first-integration", tfModel.Name.ValueString())
				assert.Equal(t, "elastic", tfModel.ServiceName.ValueString())
				assert.Equal(t, "FIRST001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
			},
		},
		{
			name: "Integration with unknown type",
			apiModel: &integrationapi.IntegrationResponseAPIModel{
				Output: integrationapi.IntegrationOutput{
					Integrations: integrationapi.IntegrationConfigurations{
						Items: []integrationapi.IntegrationItem{
							{
								Name:            "unknown-integration",
								Enabled:         true,
								SerialNumber:    "UNK001",
								PermissionsUser: "unknown@company.com",
								IntegrationType: "UNKNOWN_TYPE",
								IntegrationData: map[string]interface{}{
									"some_field": "some_value",
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "unknown-integration", tfModel.Name.ValueString())
				assert.Equal(t, "unknown_type", tfModel.ServiceName.ValueString()) // Should be lowercase
				assert.Equal(t, "UNK001", tfModel.SerialNumber.ValueString())
				assert.True(t, tfModel.Enabled.ValueBool())
				assert.Equal(t, "unknown@company.com", tfModel.PermissionsUser.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
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
func TestIntegrationTranslator_EdgeCases(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslator{}

	t.Run("Empty integration data", func(t *testing.T) {
		// Given
		apiModel := &integrationapi.IntegrationResponseAPIModel{
			Output: integrationapi.IntegrationOutput{
				Integrations: integrationapi.IntegrationConfigurations{
					Items: []integrationapi.IntegrationItem{
						{
							Name:            "empty-data-integration",
							Enabled:         true,
							SerialNumber:    "EMPTY001",
							PermissionsUser: "empty@company.com",
							IntegrationType: "DATADOG",
							IntegrationData: map[string]interface{}{}, // Empty data
						},
					},
				},
			},
		}
		requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
		translationData := &coretranslator.TranslationData{}

		// When
		result, err := translator.ToTFModel(requestContext, translationData, apiModel)

		// Then
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "empty-data-integration", result.Name.ValueString())
		assert.Equal(t, "datadog", result.ServiceName.ValueString())
		// Integration-specific fields should have default/empty values
		assert.Equal(t, "", result.APIKey.ValueString())
		assert.Equal(t, "", result.APIUrl.ValueString())
	})

	t.Run("Null values in TF model", func(t *testing.T) {
		// Given
		tfModel := &integrationtf.IntegrationTFModel{
			Name:            types.StringValue("null-values-test"),
			ServiceName:     types.StringValue("okta"),
			SerialNumber:    types.StringNull(),
			Enabled:         types.BoolNull(),
			PermissionsUser: types.StringNull(),
			APIKey:          types.StringNull(),
			APIUrl:          types.StringNull(),
		}
		requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
		translationData := &coretranslator.TranslationData{}

		// When
		result, err := translator.ToAPIModel(requestContext, translationData, tfModel)

		// Then
		require.NoError(t, err)
		expected := `define_integration(integration_name="null-values-test", serial_number="", enabled=false, permissions_user="", params={"api_rate_limit":0,"api_token":"","cache_ttl_ms":0,"idp_name":"","url":""}, service_name="okta")`
		assert.Equal(t, expected, result.Statement)
	})

	t.Run("Complex payload paths", func(t *testing.T) {
		// Given
		apiModel := &integrationapi.IntegrationResponseAPIModel{
			Output: integrationapi.IntegrationOutput{
				Integrations: integrationapi.IntegrationConfigurations{
					Items: []integrationapi.IntegrationItem{
						{
							Name:            "complex-paths",
							Enabled:         true,
							SerialNumber:    "COMPLEX001",
							PermissionsUser: "complex@company.com",
							IntegrationType: "ALERTMANAGER",
							IntegrationData: map[string]interface{}{
								"external_url": "https://alertmanager.complex.com",
								"payload_paths": []interface{}{
									"alerts[0].receiver",
									"alerts[*].labels.severity",
									"commonLabels.alertname",
									"groupLabels.instance",
								},
							},
						},
					},
				},
			},
		}
		requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
		translationData := &coretranslator.TranslationData{}

		// When
		result, err := translator.ToTFModel(requestContext, translationData, apiModel)

		// Then
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "complex-paths", result.Name.ValueString())
		assert.False(t, result.PayloadPaths.IsNull())
		assert.Len(t, result.PayloadPaths.Elements(), 4)
	})
}
