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

func TestGoogleCloudIdentityDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"idp_name":       "google-cloud-test",
				"cache_ttl_ms":   600000,
				"api_rate_limit": 100,
				"subject":        "service-account@project.iam.gserviceaccount.com",
				"credentials":    `{"type": "service_account", "project_id": "test-project"}`,
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "google-cloud-test", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(600000), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(100), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "service-account@project.iam.gserviceaccount.com", tfModel.Subject.ValueString())
				assert.Equal(t, `{"type": "service_account", "project_id": "test-project"}`, tfModel.Credentials.ValueString())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(0), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(0), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "", tfModel.Subject.ValueString())
				assert.Equal(t, "", tfModel.Credentials.ValueString())
			},
		},
		{
			name: "Partial data - required fields only",
			integrationData: map[string]interface{}{
				"subject":     "minimal@example.com",
				"credentials": "minimal-creds",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.IDPName.ValueString())
				assert.Equal(t, int64(0), tfModel.CacheTTLMs.ValueInt64())
				assert.Equal(t, int64(0), tfModel.APIRateLimit.ValueInt64())
				assert.Equal(t, "minimal@example.com", tfModel.Subject.ValueString())
				assert.Equal(t, "minimal-creds", tfModel.Credentials.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &GoogleCloudIdentityDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestGoogleCloudIdentityDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringValue("google-cloud-prod"),
				CacheTTLMs:   types.Int64Value(900000),
				APIRateLimit: types.Int64Value(200),
				Subject:      types.StringValue("prod-service@prod-project.iam.gserviceaccount.com"),
				Credentials:  types.StringValue(`{"type": "service_account", "project_id": "prod-project"}`),
			},
			expected: map[string]interface{}{
				"idp_name":       "google-cloud-prod",
				"cache_ttl_ms":   int64(900000),
				"api_rate_limit": int64(200),
				"subject":        "prod-service@prod-project.iam.gserviceaccount.com",
				"credentials":    `{"type": "service_account", "project_id": "prod-project"}`,
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				IDPName:      types.StringValue(""),
				CacheTTLMs:   types.Int64Value(0),
				APIRateLimit: types.Int64Value(0),
				Subject:      types.StringValue(""),
				Credentials:  types.StringValue(""),
			},
			expected: map[string]interface{}{
				"idp_name":       "",
				"cache_ttl_ms":   int64(0),
				"api_rate_limit": int64(0),
				"subject":        "",
				"credentials":    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &GoogleCloudIdentityDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
