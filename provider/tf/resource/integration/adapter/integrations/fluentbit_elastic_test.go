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

func TestFluentbitElasticDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with URL",
			integrationData: map[string]interface{}{
				"url": "https://fluentbit-elastic.example.com:9200",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "https://fluentbit-elastic.example.com:9200", tfModel.APIUrl.ValueString())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.APIUrl.ValueString())
			},
		},
		{
			name: "Data with localhost URL",
			integrationData: map[string]interface{}{
				"url": "http://localhost:9200",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "http://localhost:9200", tfModel.APIUrl.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &FluentbitElasticDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestFluentbitElasticDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with URL",
			tfModel: &model.IntegrationTFModel{
				APIUrl: types.StringValue("https://prod-fluentbit.elastic.co:9200"),
			},
			expected: map[string]interface{}{
				"url": "https://prod-fluentbit.elastic.co:9200",
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				APIUrl: types.StringValue(""),
			},
			expected: map[string]interface{}{
				"url": "",
			},
		},
		{
			name: "TF model with null URL",
			tfModel: &model.IntegrationTFModel{
				APIUrl: types.StringNull(),
			},
			expected: map[string]interface{}{
				"url": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &FluentbitElasticDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
