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

func TestBcmConnectivityDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"api_key":         "test-api-key-123",
				"api_certificate": "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "test-api-key-123", tfModel.APIKey.ValueString())
				assert.Equal(t, "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----", tfModel.APICertificate.ValueString())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.APIKey.ValueString())
				assert.Equal(t, "", tfModel.APICertificate.ValueString())
			},
		},
		{
			name: "Partial data - only api_key",
			integrationData: map[string]interface{}{
				"api_key": "partial-key",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "partial-key", tfModel.APIKey.ValueString())
				assert.Equal(t, "", tfModel.APICertificate.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &BcmConnectivityDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestBcmConnectivityDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				APIKey:         types.StringValue("prod-api-key"),
				APICertificate: types.StringValue("prod-certificate"),
			},
			expected: map[string]interface{}{
				"api_key":         "prod-api-key",
				"api_certificate": "prod-certificate",
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				APIKey:         types.StringValue(""),
				APICertificate: types.StringValue(""),
			},
			expected: map[string]interface{}{
				"api_key":         "",
				"api_certificate": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &BcmConnectivityDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
