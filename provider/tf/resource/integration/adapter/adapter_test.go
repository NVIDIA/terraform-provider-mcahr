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

package adapter

import (
	"context"
	"encoding/json"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTFDataToMap(t *testing.T) {
	t.Parallel()

	compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
	adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}

	tests := []struct {
		name        string
		tfModel     *integrationtf.IntegrationTFModel
		expectError bool
		validate    func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "Valid Datadog model",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName: types.StringValue("datadog"),
				APIKey:      types.StringValue("dd-api-key"),
				APIUrl:      types.StringValue("https://api.datadoghq.com"),
				SiteUrl:     types.StringValue("https://app.datadoghq.com"),
				AppKey:      types.StringValue("dd-app-key"),
				WebhookName: types.StringValue("my-webhook"),
			},
			expectError: false,
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "dd-api-key", result["api_key"])
				assert.Equal(t, "https://api.datadoghq.com", result["api_url"])
				assert.Equal(t, "https://app.datadoghq.com", result["site_url"])
				assert.Equal(t, "dd-app-key", result["app_key"])
				assert.Equal(t, "my-webhook", result["webhook_name"])
			},
		},
		{
			name: "Valid Alertmanager model",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName: types.StringValue("alertmanager"),
				ExternalUrl: types.StringValue("https://alertmanager.example.com"),
				PayloadPaths: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("alerts.receiver"),
					types.StringValue("alerts.status"),
				}),
			},
			expectError: false,
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "https://alertmanager.example.com", result["external_url"])
				payloadPaths, ok := result["payload_paths"].([]string)
				require.True(t, ok)
				assert.ElementsMatch(t, []string{"alerts.receiver", "alerts.status"}, payloadPaths)
			},
		},
		{
			name: "Valid Azure AD model",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName:  types.StringValue("azure_active_directory"),
				IDPName:      types.StringValue("azure-ad-test"),
				CacheTTLMs:   types.Int64Value(300000),
				APIRateLimit: types.Int64Value(100),
				TenantID:     types.StringValue("tenant-123"),
				ClientID:     types.StringValue("client-456"),
				ClientSecret: types.StringValue("secret-789"),
			},
			expectError: false,
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "azure-ad-test", result["idp_name"])
				assert.Equal(t, int64(300000), result["cache_ttl_ms"])
				assert.Equal(t, int64(100), result["api_rate_limit"])
				assert.Equal(t, "tenant-123", result["tenant_id"])
				assert.Equal(t, "client-456", result["client_id"])
				assert.Equal(t, "secret-789", result["client_secret"])
			},
		},
		{
			name: "Unsupported service name",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName: types.StringValue("unsupported_service"),
			},
			expectError: true,
		},
		{
			name: "Empty service name",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName: types.StringValue(""),
			},
			expectError: true,
		},
		{
			name: "Null service name",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName: types.StringNull(),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// When
			result, err := TFDataToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			// Then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				tt.validate(t, result)
			}
		})
	}
}

func TestMapToTFData(t *testing.T) {
	t.Parallel()

	compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
	adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		serviceName     string
		expectError     bool
		validate        func(t *testing.T, tfModel *integrationtf.IntegrationTFModel)
	}{
		{
			name: "Valid Datadog data",
			integrationData: map[string]interface{}{
				"api_key":      "dd-key-123",
				"api_url":      "https://api.datadoghq.eu",
				"site_url":     "https://app.datadoghq.eu",
				"app_key":      "dd-app-456",
				"webhook_name": "test-webhook",
			},
			serviceName: "datadog",
			expectError: false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "dd-key-123", tfModel.APIKey.ValueString())
				assert.Equal(t, "https://api.datadoghq.eu", tfModel.APIUrl.ValueString())
				assert.Equal(t, "https://app.datadoghq.eu", tfModel.SiteUrl.ValueString())
				assert.Equal(t, "dd-app-456", tfModel.AppKey.ValueString())
				assert.Equal(t, "test-webhook", tfModel.WebhookName.ValueString())
			},
		},
		{
			name: "Valid Elastic data with field mapping",
			integrationData: map[string]interface{}{
				"api_token": "elastic-token", // Data field name
				"url":       "https://elastic.example.com",
			},
			serviceName: "elastic",
			expectError: false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				// Should map to TF field names
				assert.Equal(t, "elastic-token", tfModel.APIKey.ValueString())
				assert.Equal(t, "https://elastic.example.com", tfModel.APIUrl.ValueString())
			},
		},
		{
			name: "Valid NVault data",
			integrationData: map[string]interface{}{
				"address":       "https://vault.example.com:8200",
				"namespace":     "admin/test",
				"role_name":     "test-role",
				"jwt_auth_path": "auth/jwt",
			},
			serviceName: "nvault",
			expectError: false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "https://vault.example.com:8200", tfModel.Address.ValueString())
				assert.Equal(t, "admin/test", tfModel.Namespace.ValueString())
				assert.Equal(t, "test-role", tfModel.RoleName.ValueString())
				assert.Equal(t, "auth/jwt", tfModel.JWTAuthPath.ValueString())
			},
		},
		{
			name: "Unsupported service name",
			integrationData: map[string]interface{}{
				"api_key": "test-key",
			},
			serviceName: "unsupported_service",
			expectError: true,
		},
		{
			name: "Empty service name",
			integrationData: map[string]interface{}{
				"api_key": "test-key",
			},
			serviceName: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Given
			tfModel := &integrationtf.IntegrationTFModel{
				ServiceName: types.StringValue(tt.serviceName),
			}

			// When
			err := MapToTFData(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			// Then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.validate(t, tfModel)
			}
		})
	}
}

func TestTFDataToJSON(t *testing.T) {
	t.Parallel()

	compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
	adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}

	tests := []struct {
		name        string
		tfModel     *integrationtf.IntegrationTFModel
		expectError bool
		validate    func(t *testing.T, jsonStr string)
	}{
		{
			name: "Valid BCM model",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName:  types.StringValue("bcm"),
				IDPName:      types.StringValue("bcm-test"),
				CacheTTLMs:   types.Int64Value(180000),
				APIRateLimit: types.Int64Value(50),
			},
			expectError: false,
			validate: func(t *testing.T, jsonStr string) {
				var data map[string]interface{}
				err := json.Unmarshal([]byte(jsonStr), &data)
				require.NoError(t, err)

				assert.Equal(t, "bcm-test", data["idp_name"])
				assert.Equal(t, float64(180000), data["cache_ttl_ms"]) // JSON unmarshals numbers as float64
				assert.Equal(t, float64(50), data["api_rate_limit"])
			},
		},
		{
			name: "Valid BCM Connectivity model",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName:    types.StringValue("bcm_connectivity"),
				APIKey:         types.StringValue("bcm-key"),
				APICertificate: types.StringValue("cert-data"),
			},
			expectError: false,
			validate: func(t *testing.T, jsonStr string) {
				var data map[string]interface{}
				err := json.Unmarshal([]byte(jsonStr), &data)
				require.NoError(t, err)

				assert.Equal(t, "bcm-key", data["api_key"])
				assert.Equal(t, "cert-data", data["api_certificate"])
			},
		},
		{
			name: "Unsupported service name",
			tfModel: &integrationtf.IntegrationTFModel{
				ServiceName: types.StringValue("unsupported"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// When
			result, err := TFDataToJSON(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

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

func TestJSONToTFData(t *testing.T) {
	t.Parallel()

	compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
	adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}

	tests := []struct {
		name            string
		integrationData string
		serviceName     string
		expectError     bool
		validate        func(t *testing.T, tfModel *integrationtf.IntegrationTFModel)
	}{
		{
			name: "Valid Google Cloud Identity JSON",
			integrationData: `{
				"idp_name": "google-test",
				"cache_ttl_ms": 600000,
				"api_rate_limit": 100,
				"subject": "service@project.iam.gserviceaccount.com",
				"credentials": "{\"type\": \"service_account\"}"
			}`,
			serviceName: "google_cloud_identity",
			expectError: false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "google-test", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(600000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(100), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "service@project.iam.gserviceaccount.com", tfModel.Subject.ValueString())
				assert.Equal(t, `{"type": "service_account"}`, tfModel.Credentials.ValueString())
			},
		},
		{
			name: "Valid Okta JSON with field mapping",
			integrationData: `{
				"idp_name": "okta-prod",
				"cache_ttl_ms": 300000,
				"api_rate_limit": 60,
				"api_token": "okta-token-xyz",
				"url": "https://prod.okta.com"
			}`,
			serviceName: "okta",
			expectError: false,
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				assert.Equal(t, "okta-prod", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(300000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(60), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "okta-token-xyz", tfModel.APIKey.ValueString())        // Maps to TF field
				assert.Equal(t, "https://prod.okta.com", tfModel.APIUrl.ValueString()) // Maps to TF field
			},
		},
		{
			name:            "Invalid JSON",
			integrationData: `{invalid json}`,
			serviceName:     "datadog",
			expectError:     true,
		},
		{
			name:            "Empty JSON",
			integrationData: "",
			serviceName:     "datadog",
			expectError:     true,
		},
		{
			name:            "Null JSON",
			integrationData: "null",
			serviceName:     "datadog",
			expectError:     false, // null JSON is valid, should result in empty data
			validate: func(t *testing.T, tfModel *integrationtf.IntegrationTFModel) {
				// Should have empty/default values since JSON was null
				assert.Equal(t, "", tfModel.APIKey.ValueString())
			},
		},
		{
			name: "Valid JSON with unsupported service",
			integrationData: `{
				"api_key": "test-key"
			}`,
			serviceName: "unsupported_service",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Given
			tfModel := &integrationtf.IntegrationTFModel{
				ServiceName: types.StringValue(tt.serviceName),
			}

			// When
			err := JSONToTFData(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			// Then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, tfModel)
				}
			}
		})
	}
}

func TestAdapter_RoundTrip(t *testing.T) {
	t.Parallel()

	compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
	adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}

	// Test round-trip conversion: TF Model -> Map -> JSON -> Map -> TF Model
	t.Run("Fluentbit Elastic round trip", func(t *testing.T) {
		original := &integrationtf.IntegrationTFModel{
			ServiceName: types.StringValue("fluentbit_elastic"),
			APIUrl:      types.StringValue("https://fluentbit.elastic.co:9200"),
		}

		// TF -> Map
		mapData, err := TFDataToMap(common.NewRequestContext(context.Background()), adapterOptions, original)
		require.NoError(t, err)

		// TF -> JSON
		jsonData, err := TFDataToJSON(common.NewRequestContext(context.Background()), adapterOptions, original)
		require.NoError(t, err)

		// JSON -> TF
		recovered := &integrationtf.IntegrationTFModel{
			ServiceName: types.StringValue("fluentbit_elastic"),
		}
		err = JSONToTFData(common.NewRequestContext(context.Background()), adapterOptions, jsonData, recovered)
		require.NoError(t, err)

		// Verify round-trip integrity
		assert.Equal(t, original.APIUrl.ValueString(), recovered.APIUrl.ValueString())

		// Map -> TF
		recovered2 := &integrationtf.IntegrationTFModel{
			ServiceName: types.StringValue("fluentbit_elastic"),
		}
		err = MapToTFData(common.NewRequestContext(context.Background()), adapterOptions, mapData, recovered2)
		require.NoError(t, err)

		assert.Equal(t, original.APIUrl.ValueString(), recovered2.APIUrl.ValueString())
	})
}
