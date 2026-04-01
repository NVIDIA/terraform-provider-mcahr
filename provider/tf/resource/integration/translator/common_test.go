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
	"terraform/terraform-provider/provider/external_api/resources/statement"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"
	"terraform/terraform-provider/provider/tf/resource/integration/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationTranslatorCommon_ToAPIModelWithVersion(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslatorCommon{}

	tests := []struct {
		name           string
		tfModel        *integrationtf.IntegrationTFModel
		operation      common.CrudOperation
		backendVersion common.APIVersion
		expectError    bool
		validateResult func(t *testing.T, apiModel *statement.StatementInputAPIModel)
	}{
		{
			name: "Create operation with V2 backend",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:            types.StringValue("test-datadog-integration"),
				ServiceName:     types.StringValue("datadog"),
				SerialNumber:    types.StringValue("SN123456"),
				Enabled:         types.BoolValue(true),
				PermissionsUser: types.StringValue("admin@example.com"),
				APIKey:          types.StringValue("dd-api-key"),
				APIUrl:          types.StringValue("https://api.datadoghq.com"),
			},
			operation:      common.Create,
			backendVersion: common.V2,
			expectError:    false,
			validateResult: func(t *testing.T, apiModel *statement.StatementInputAPIModel) {
				expected := `define_integration(integration_name="test-datadog-integration", serial_number="SN123456", enabled=true, permissions_user="admin@example.com", params={"api_key":"dd-api-key","api_url":"https://api.datadoghq.com","app_key":"","site_url":"","webhook_name":""}, service_name="datadog")`
				assert.Equal(t, expected, apiModel.Statement)
				assert.Equal(t, common.V2, apiModel.APIVersion)
			},
		},
		{
			name: "Read operation with V1 backend",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue("test-okta-integration"),
			},
			operation:      common.Read,
			backendVersion: common.V1,
			expectError:    false,
			validateResult: func(t *testing.T, apiModel *statement.StatementInputAPIModel) {
				expectedStatement := `get_integration_class(integration_name="test-okta-integration")`
				assert.Equal(t, expectedStatement, apiModel.Statement)
				assert.Equal(t, common.V1, apiModel.APIVersion)
			},
		},
		{
			name: "Update operation with Azure AD",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:            types.StringValue("azure-ad-prod"),
				ServiceName:     types.StringValue("azure_active_directory"),
				SerialNumber:    types.StringValue("SN789012"),
				Enabled:         types.BoolValue(false),
				PermissionsUser: types.StringValue("service@domain.com"),
				IDPName:         types.StringValue("azure-ad-idp"),
				TenantID:        types.StringValue("tenant-123"),
				ClientID:        types.StringValue("client-456"),
				ClientSecret:    types.StringValue("secret-789"),
			},
			operation:      common.Update,
			backendVersion: common.V2,
			expectError:    false,
			validateResult: func(t *testing.T, apiModel *statement.StatementInputAPIModel) {
				expected := `update_integration(integration_name="azure-ad-prod", serial_number="SN789012", enabled=false, permissions_user="service@domain.com", params={"api_rate_limit":0,"cache_ttl_ms":0,"client_id":"client-456","client_secret":"secret-789","idp_name":"azure-ad-idp","tenant_id":"tenant-123"})`
				assert.Equal(t, expected, apiModel.Statement)
				assert.Equal(t, common.V2, apiModel.APIVersion)
			},
		},
		{
			name: "Delete operation",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue("integration-to-delete"),
			},
			operation:      common.Delete,
			backendVersion: common.V1,
			expectError:    false,
			validateResult: func(t *testing.T, apiModel *statement.StatementInputAPIModel) {
				expectedStatement := `delete_integration(integration_name="integration-to-delete")`
				assert.Equal(t, expectedStatement, apiModel.Statement)
				assert.Equal(t, common.V1, apiModel.APIVersion)
			},
		},
		{
			name: "Unsupported operation",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue("test-integration"),
			},
			operation:      common.CrudOperation(999), // Invalid operation
			backendVersion: common.V2,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestContext := common.NewRequestContext(context.Background()).WithOperation(tt.operation).WithAPIVersion(tt.backendVersion)
			translationData := &coretranslator.TranslationData{}
			result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tt.tfModel)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				tt.validateResult(t, result)
			}
		})
	}
}

func TestIntegrationTranslatorCommon_BuildCreateStatement(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslatorCommon{}

	integrationSchema := schema.IntegrationSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: integrationSchema.GetCompatibilityOptions()}

	tests := []struct {
		name        string
		tfModel     *integrationtf.IntegrationTFModel
		expectError bool
		validate    func(t *testing.T, statement string)
	}{
		{
			name: "Valid Alertmanager integration",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:            types.StringValue("alertmanager-test"),
				ServiceName:     types.StringValue("alertmanager"),
				SerialNumber:    types.StringValue("AM001"),
				Enabled:         types.BoolValue(true),
				PermissionsUser: types.StringValue("alerts@company.com"),
				ExternalUrl:     types.StringValue("https://alertmanager.company.com"),
			},
			expectError: false,
			validate: func(t *testing.T, statement string) {
				expected := `define_integration(integration_name="alertmanager-test", serial_number="AM001", enabled=true, permissions_user="alerts@company.com", params={"external_url":"https://alertmanager.company.com","payload_paths":[]}, service_name="alertmanager")`
				assert.Equal(t, expected, statement)
			},
		},
		{
			name: "Valid NVault integration",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:            types.StringValue("vault-prod"),
				ServiceName:     types.StringValue("nvault"),
				SerialNumber:    types.StringValue("VAULT001"),
				Enabled:         types.BoolValue(false),
				PermissionsUser: types.StringValue("vault@company.com"),
				Address:         types.StringValue("https://vault.company.com:8200"),
				Namespace:       types.StringValue("prod/team1"),
				RoleName:        types.StringValue("terraform-role"),
				JWTAuthPath:     types.StringValue("auth/jwt-prod"),
			},
			expectError: false,
			validate: func(t *testing.T, statement string) {
				expected := `define_integration(integration_name="vault-prod", serial_number="VAULT001", enabled=false, permissions_user="vault@company.com", params={"address":"https://vault.company.com:8200","jwt_auth_path":"auth/jwt-prod","namespace":"prod/team1","role_name":"terraform-role"}, service_name="nvault")`
				assert.Equal(t, expected, statement)
			},
		},
		{
			name: "Integration with unsupported service name",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:        types.StringValue("unsupported-integration"),
				ServiceName: types.StringValue("unsupported_service"),
			},
			expectError: true,
		},
		{
			name: "Integration with empty service name",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:        types.StringValue("empty-service-integration"),
				ServiceName: types.StringValue(""),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)

			// When
			result, err := translator.buildCreateStatement(requestContext, translationData, tt.tfModel)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, result)
				tt.validate(t, result)
			}
		})
	}
}

func TestIntegrationTranslatorCommon_BuildReadStatement(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslatorCommon{}

	tests := []struct {
		name     string
		tfModel  *integrationtf.IntegrationTFModel
		expected string
	}{
		{
			name: "Simple integration name",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue("test-integration"),
			},
			expected: `get_integration_class(integration_name="test-integration")`,
		},
		{
			name: "Integration name with special characters",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue("my-integration_2024"),
			},
			expected: `get_integration_class(integration_name="my-integration_2024")`,
		},
		{
			name: "Integration name with spaces",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue("integration with spaces"),
			},
			expected: `get_integration_class(integration_name="integration with spaces")`,
		},
		{
			name: "Empty integration name",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue(""),
			},
			expected: `get_integration_class(integration_name="")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translator.buildReadStatement(tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntegrationTranslatorCommon_BuildUpdateStatement(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslatorCommon{}

	integrationSchema := schema.IntegrationSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: integrationSchema.GetCompatibilityOptions()}

	tests := []struct {
		name        string
		tfModel     *integrationtf.IntegrationTFModel
		expectError bool
		validate    func(t *testing.T, statement string)
	}{
		{
			name: "Valid BCM update",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:            types.StringValue("bcm-update-test"),
				ServiceName:     types.StringValue("bcm"),
				SerialNumber:    types.StringValue("BCM002"),
				Enabled:         types.BoolValue(true),
				PermissionsUser: types.StringValue("bcm@company.com"),
				IDPName:         types.StringValue("bcm-idp-updated"),
				CacheTTLMs:      types.Int64Value(600000),
				APIRateLimit:    types.Int64Value(150),
			},
			expectError: false,
			validate: func(t *testing.T, statement string) {
				expected := `update_integration(integration_name="bcm-update-test", serial_number="BCM002", enabled=true, permissions_user="bcm@company.com", params={"api_rate_limit":150,"cache_ttl_ms":600000,"idp_name":"bcm-idp-updated"})`
				assert.Equal(t, expected, statement)
			},
		},
		{
			name: "Update with unsupported service",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:        types.StringValue("unsupported-update"),
				ServiceName: types.StringValue("unknown_service"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Update).WithAPIVersion(common.V2)

			// When
			result, err := translator.buildUpdateStatement(requestContext, translationData, tt.tfModel)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, result)
				tt.validate(t, result)
			}
		})
	}
}

func TestIntegrationTranslatorCommon_BuildDeleteStatement(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslatorCommon{}

	tests := []struct {
		name     string
		tfModel  *integrationtf.IntegrationTFModel
		expected string
	}{
		{
			name: "Standard delete",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue("delete-me"),
			},
			expected: `delete_integration(integration_name="delete-me")`,
		},
		{
			name: "Delete with complex name",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue("complex-integration_name-2024"),
			},
			expected: `delete_integration(integration_name="complex-integration_name-2024")`,
		},
		{
			name: "Delete with empty name",
			tfModel: &integrationtf.IntegrationTFModel{
				Name: types.StringValue(""),
			},
			expected: `delete_integration(integration_name="")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translator.buildDeleteStatement(tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntegrationTranslatorCommon_BuildIntegrationStatement(t *testing.T) {
	t.Parallel()

	translator := &IntegrationTranslatorCommon{}

	integrationSchema := schema.IntegrationSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: integrationSchema.GetCompatibilityOptions()}

	tests := []struct {
		name          string
		statementName string
		tfModel       *integrationtf.IntegrationTFModel
		expectError   bool
		validate      func(t *testing.T, statement string)
	}{
		{
			name:          "Define integration with Google Cloud Identity",
			statementName: "define_integration",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:            types.StringValue("gci-test"),
				ServiceName:     types.StringValue("google_cloud_identity"),
				SerialNumber:    types.StringValue("GCI001"),
				Enabled:         types.BoolValue(true),
				PermissionsUser: types.StringValue("gci@company.com"),
				IDPName:         types.StringValue("google-idp"),
				CacheTTLMs:      types.Int64Value(300000),
				APIRateLimit:    types.Int64Value(100),
				Subject:         types.StringValue("service@project.iam.gserviceaccount.com"),
				Credentials:     types.StringValue(`{"type": "service_account"}`),
			},
			expectError: false,
			validate: func(t *testing.T, statement string) {
				expected := `define_integration(integration_name="gci-test", serial_number="GCI001", enabled=true, permissions_user="gci@company.com", params={"api_rate_limit":100,"cache_ttl_ms":300000,"credentials":"{\"type\": \"service_account\"}","idp_name":"google-idp","subject":"service@project.iam.gserviceaccount.com"}, service_name="google_cloud_identity")`
				assert.Equal(t, expected, statement)
			},
		},
		{
			name:          "Update integration without service_name",
			statementName: "update_integration",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:            types.StringValue("update-test"),
				ServiceName:     types.StringValue("alertmanager"),
				SerialNumber:    types.StringValue("AM001"),
				Enabled:         types.BoolValue(false),
				PermissionsUser: types.StringValue("alerts@company.com"),
				ExternalUrl:     types.StringValue("https://alertmanager.company.com"),
			},
			expectError: false,
			validate: func(t *testing.T, statement string) {
				expected := `update_integration(integration_name="update-test", serial_number="AM001", enabled=false, permissions_user="alerts@company.com", params={"external_url":"https://alertmanager.company.com","payload_paths":[]})`
				assert.Equal(t, expected, statement)
			},
		},
		{
			name:          "Integration with unsupported service",
			statementName: "define_integration",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:        types.StringValue("bad-service"),
				ServiceName: types.StringValue("bad_service"),
			},
			expectError: true,
		},
		{
			name:          "Custom statement name",
			statementName: "custom_integration_operation",
			tfModel: &integrationtf.IntegrationTFModel{
				Name:            types.StringValue("custom-test"),
				ServiceName:     types.StringValue("datadog"),
				SerialNumber:    types.StringValue("CUSTOM001"),
				Enabled:         types.BoolValue(true),
				PermissionsUser: types.StringValue("custom@company.com"),
				APIKey:          types.StringValue("custom-key"),
			},
			expectError: false,
			validate: func(t *testing.T, statement string) {
				expected := `custom_integration_operation(integration_name="custom-test", serial_number="CUSTOM001", enabled=true, permissions_user="custom@company.com", params={"api_key":"custom-key","api_url":"","app_key":"","site_url":"","webhook_name":""})`
				assert.Equal(t, expected, statement)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)

			// When
			result, err := translator.buildIntegrationStatement(requestContext, translationData, tt.statementName, tt.tfModel)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, result)
				tt.validate(t, result)
			}
		})
	}
}
