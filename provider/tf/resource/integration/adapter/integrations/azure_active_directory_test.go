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

package integrations

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	"terraform/terraform-provider/provider/tf/resource/integration/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestAzureActiveDirectoryDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"idp_name":       "azure-ad-test",
				"cache_ttl_ms":   300000,
				"api_rate_limit": 100,
				"tenant_id":      "12345678-1234-1234-1234-123456789abc",
				"client_id":      "87654321-4321-4321-4321-cba987654321",
				"client_secret":  "super-secret-key",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "azure-ad-test", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(300000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(100), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "12345678-1234-1234-1234-123456789abc", tfModel.TenantID.ValueString())
				assert.Equal(t, "87654321-4321-4321-4321-cba987654321", tfModel.ClientID.ValueString())
				assert.Equal(t, "super-secret-key", tfModel.ClientSecret.ValueString())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(0), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(0), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "", tfModel.TenantID.ValueString())
				assert.Equal(t, "", tfModel.ClientID.ValueString())
				assert.Equal(t, "", tfModel.ClientSecret.ValueString())
			},
		},
		{
			name: "Partial data - required fields only",
			integrationData: map[string]interface{}{
				"tenant_id": "test-tenant",
				"client_id": "test-client",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(0), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(0), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "test-tenant", tfModel.TenantID.ValueString())
				assert.Equal(t, "test-client", tfModel.ClientID.ValueString())
				assert.Equal(t, "", tfModel.ClientSecret.ValueString())
			},
		},
		{
			name: "Data with different numeric types",
			integrationData: map[string]interface{}{
				"cache_ttl_ms":   float64(600000),
				"api_rate_limit": int32(50),
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, int64(600000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(50), tfModel.APIRateLimit.ValueInt64())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &AzureActiveDirectoryDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestAzureActiveDirectoryDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringValue("azure-ad-prod"),
				CacheTTLMs:   types.Int64Value(600000),
				APIRateLimit: types.Int64Value(200),
				TenantID:     types.StringValue("prod-tenant-id"),
				ClientID:     types.StringValue("prod-client-id"),
				ClientSecret: types.StringValue("prod-secret"),
			},
			expected: map[string]interface{}{
				"idp_name":       "azure-ad-prod",
				"cache_ttl_ms":   int64(600000),
				"api_rate_limit": int64(200),
				"tenant_id":      "prod-tenant-id",
				"client_id":      "prod-client-id",
				"client_secret":  "prod-secret",
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringValue(""),
				CacheTTLMs:   types.Int64Value(0),
				APIRateLimit: types.Int64Value(0),
				TenantID:     types.StringValue(""),
				ClientID:     types.StringValue(""),
				ClientSecret: types.StringValue(""),
			},
			expected: map[string]interface{}{
				"idp_name":       "",
				"cache_ttl_ms":   int64(0),
				"api_rate_limit": int64(0),
				"tenant_id":      "",
				"client_id":      "",
				"client_secret":  "",
			},
		},
		{
			name: "TF model with null values",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringNull(),
				CacheTTLMs:   types.Int64Null(),
				APIRateLimit: types.Int64Null(),
				TenantID:     types.StringNull(),
				ClientID:     types.StringNull(),
				ClientSecret: types.StringNull(),
			},
			expected: map[string]interface{}{
				"idp_name":       "",
				"cache_ttl_ms":   int64(0),
				"api_rate_limit": int64(0),
				"tenant_id":      "",
				"client_id":      "",
				"client_secret":  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &AzureActiveDirectoryDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
