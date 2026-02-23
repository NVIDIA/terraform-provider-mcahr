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

func TestDatadogDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"api_key":      "dd-api-key-123",
				"api_url":      "https://api.datadoghq.com",
				"site_url":     "https://app.datadoghq.com",
				"app_key":      "dd-app-key-456",
				"webhook_name": "my-webhook",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "dd-api-key-123", tfModel.APIKey.ValueString())
				assert.Equal(t, "https://api.datadoghq.com", tfModel.APIUrl.ValueString())
				assert.Equal(t, "https://app.datadoghq.com", tfModel.SiteUrl.ValueString())
				assert.Equal(t, "dd-app-key-456", tfModel.AppKey.ValueString())
				assert.Equal(t, "my-webhook", tfModel.WebhookName.ValueString())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.APIKey.ValueString())
				assert.Equal(t, "", tfModel.APIUrl.ValueString())
				assert.Equal(t, "", tfModel.SiteUrl.ValueString())
				assert.Equal(t, "", tfModel.AppKey.ValueString())
				assert.Equal(t, "", tfModel.WebhookName.ValueString())
			},
		},
		{
			name: "Partial data - required fields only",
			integrationData: map[string]interface{}{
				"api_key": "required-key",
				"api_url": "required-url",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "required-key", tfModel.APIKey.ValueString())
				assert.Equal(t, "required-url", tfModel.APIUrl.ValueString())
				assert.Equal(t, "", tfModel.SiteUrl.ValueString())
				assert.Equal(t, "", tfModel.AppKey.ValueString())
				assert.Equal(t, "", tfModel.WebhookName.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DatadogDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestDatadogDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				APIKey:      types.StringValue("prod-api-key"),
				APIUrl:      types.StringValue("https://api.datadoghq.eu"),
				SiteUrl:     types.StringValue("https://app.datadoghq.eu"),
				AppKey:      types.StringValue("prod-app-key"),
				WebhookName: types.StringValue("prod-webhook"),
			},
			expected: map[string]interface{}{
				"api_key":      "prod-api-key",
				"api_url":      "https://api.datadoghq.eu",
				"site_url":     "https://app.datadoghq.eu",
				"app_key":      "prod-app-key",
				"webhook_name": "prod-webhook",
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				APIKey:      types.StringValue(""),
				APIUrl:      types.StringValue(""),
				SiteUrl:     types.StringValue(""),
				AppKey:      types.StringValue(""),
				WebhookName: types.StringValue(""),
			},
			expected: map[string]interface{}{
				"api_key":      "",
				"api_url":      "",
				"site_url":     "",
				"app_key":      "",
				"webhook_name": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DatadogDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
