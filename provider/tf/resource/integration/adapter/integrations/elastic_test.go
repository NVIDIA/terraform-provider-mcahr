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

func TestElasticDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"api_token": "elastic-token-123",
				"url":       "https://my-cluster.es.io:9243",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "elastic-token-123", tfModel.APIKey.ValueString())
				assert.Equal(t, "https://my-cluster.es.io:9243", tfModel.APIUrl.ValueString())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.APIKey.ValueString())
				assert.Equal(t, "", tfModel.APIUrl.ValueString())
			},
		},
		{
			name: "Partial data - only token",
			integrationData: map[string]interface{}{
				"api_token": "partial-token",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "partial-token", tfModel.APIKey.ValueString())
				assert.Equal(t, "", tfModel.APIUrl.ValueString())
			},
		},
		{
			name: "Partial data - only url",
			integrationData: map[string]interface{}{
				"url": "https://localhost:9200",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.APIKey.ValueString())
				assert.Equal(t, "https://localhost:9200", tfModel.APIUrl.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &ElasticDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestElasticDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				APIKey: types.StringValue("prod-elastic-token"),
				APIUrl: types.StringValue("https://prod-cluster.elastic.co:9243"),
			},
			expected: map[string]interface{}{
				"api_token": "prod-elastic-token",
				"url":       "https://prod-cluster.elastic.co:9243",
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				APIKey: types.StringValue(""),
				APIUrl: types.StringValue(""),
			},
			expected: map[string]interface{}{
				"api_token": "",
				"url":       "",
			},
		},
		{
			name: "TF model with null values",
			tfModel: &model.IntegrationTFModel{
				APIKey: types.StringNull(),
				APIUrl: types.StringNull(),
			},
			expected: map[string]interface{}{
				"api_token": "",
				"url":       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &ElasticDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
