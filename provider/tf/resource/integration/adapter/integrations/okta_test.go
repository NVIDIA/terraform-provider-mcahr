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

func TestOktaDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"idp_name":       "okta-test",
				"cache_ttl_ms":   300000,
				"api_rate_limit": 60,
				"api_token":      "okta-token-123",
				"url":            "https://dev-123456.okta.com",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "okta-test", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(300000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(60), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "okta-token-123", tfModel.APIKey.ValueString())
				assert.Equal(t, "https://dev-123456.okta.com", tfModel.APIUrl.ValueString())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(0), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(0), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "", tfModel.APIKey.ValueString())
				assert.Equal(t, "", tfModel.APIUrl.ValueString())
			},
		},
		{
			name: "Partial data - required fields only",
			integrationData: map[string]interface{}{
				"api_token": "required-token",
				"url":       "https://required.okta.com",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(0), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(0), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "required-token", tfModel.APIKey.ValueString())
				assert.Equal(t, "https://required.okta.com", tfModel.APIUrl.ValueString())
			},
		},
		{
			name: "Data with IDP configuration",
			integrationData: map[string]interface{}{
				"idp_name":       "production-okta",
				"cache_ttl_ms":   float64(600000),
				"api_rate_limit": int32(100),
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "production-okta", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(600000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(100), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "", tfModel.APIKey.ValueString())
				assert.Equal(t, "", tfModel.APIUrl.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &OktaDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestOktaDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringValue("okta-prod"),
				CacheTTLMs:   types.Int64Value(900000),
				APIRateLimit: types.Int64Value(120),
				APIKey:       types.StringValue("prod-okta-token"),
				APIUrl:       types.StringValue("https://prod-123456.okta.com"),
			},
			expected: map[string]interface{}{
				"idp_name":       "okta-prod",
				"cache_ttl_ms":   int64(900000),
				"api_rate_limit": int64(120),
				"api_token":      "prod-okta-token",
				"url":            "https://prod-123456.okta.com",
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringValue(""),
				CacheTTLMs:   types.Int64Value(0),
				APIRateLimit: types.Int64Value(0),
				APIKey:       types.StringValue(""),
				APIUrl:       types.StringValue(""),
			},
			expected: map[string]interface{}{
				"idp_name":       "",
				"cache_ttl_ms":   int64(0),
				"api_rate_limit": int64(0),
				"api_token":      "",
				"url":            "",
			},
		},
		{
			name: "TF model with null values",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringNull(),
				CacheTTLMs:   types.Int64Null(),
				APIRateLimit: types.Int64Null(),
				APIKey:       types.StringNull(),
				APIUrl:       types.StringNull(),
			},
			expected: map[string]interface{}{
				"idp_name":       "",
				"cache_ttl_ms":   int64(0),
				"api_rate_limit": int64(0),
				"api_token":      "",
				"url":            "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &OktaDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
