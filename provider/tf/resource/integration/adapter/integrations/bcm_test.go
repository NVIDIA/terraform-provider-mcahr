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

func TestBcmDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"idp_name":       "bcm-test",
				"cache_ttl_ms":   180000,
				"api_rate_limit": 50,
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "bcm-test", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(180000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(50), tfModel.APIRateLimit.ValueInt64())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(0), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(0), tfModel.APIRateLimit.ValueInt64())
			},
		},
		{
			name: "Partial data",
			integrationData: map[string]interface{}{
				"idp_name": "partial-bcm",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "partial-bcm", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(0), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(0), tfModel.APIRateLimit.ValueInt64())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &BcmDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestBcmDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringValue("bcm-prod"),
				CacheTTLMs:   types.Int64Value(240000),
				APIRateLimit: types.Int64Value(75),
			},
			expected: map[string]interface{}{
				"idp_name":       "bcm-prod",
				"cache_ttl_ms":   int64(240000),
				"api_rate_limit": int64(75),
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringValue(""),
				CacheTTLMs:   types.Int64Value(0),
				APIRateLimit: types.Int64Value(0),
			},
			expected: map[string]interface{}{
				"idp_name":       "",
				"cache_ttl_ms":   int64(0),
				"api_rate_limit": int64(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &BcmDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
