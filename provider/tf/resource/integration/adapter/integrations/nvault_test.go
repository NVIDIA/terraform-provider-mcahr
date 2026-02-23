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

func TestNvaultDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"address":       "https://vault.example.com:8200",
				"namespace":     "admin/tenant1",
				"role_name":     "my-vault-role",
				"jwt_auth_path": "auth/jwt",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "https://vault.example.com:8200", tfModel.Address.ValueString())
				assert.Equal(t, "admin/tenant1", tfModel.Namespace.ValueString())
				assert.Equal(t, "my-vault-role", tfModel.RoleName.ValueString())
				assert.Equal(t, "auth/jwt", tfModel.JWTAuthPath.ValueString())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.Address.ValueString())
				assert.Equal(t, "", tfModel.Namespace.ValueString())
				assert.Equal(t, "", tfModel.RoleName.ValueString())
				assert.Equal(t, "", tfModel.JWTAuthPath.ValueString())
			},
		},
		{
			name: "Partial data - minimal configuration",
			integrationData: map[string]interface{}{
				"address":   "http://localhost:8200",
				"role_name": "default-role",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "http://localhost:8200", tfModel.Address.ValueString())
				assert.Equal(t, "", tfModel.Namespace.ValueString())
				assert.Equal(t, "default-role", tfModel.RoleName.ValueString())
				assert.Equal(t, "", tfModel.JWTAuthPath.ValueString())
			},
		},
		{
			name: "Data with custom auth path",
			integrationData: map[string]interface{}{
				"address":       "https://vault.prod.com",
				"jwt_auth_path": "auth/kubernetes",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "https://vault.prod.com", tfModel.Address.ValueString())
				assert.Equal(t, "", tfModel.Namespace.ValueString())
				assert.Equal(t, "", tfModel.RoleName.ValueString())
				assert.Equal(t, "auth/kubernetes", tfModel.JWTAuthPath.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &NvaultDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestNvaultDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				Address:     types.StringValue("https://vault.production.com:8200"),
				Namespace:   types.StringValue("prod/team1"),
				RoleName:    types.StringValue("production-role"),
				JWTAuthPath: types.StringValue("auth/jwt-prod"),
			},
			expected: map[string]interface{}{
				"address":       "https://vault.production.com:8200",
				"namespace":     "prod/team1",
				"role_name":     "production-role",
				"jwt_auth_path": "auth/jwt-prod",
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				Address:     types.StringValue(""),
				Namespace:   types.StringValue(""),
				RoleName:    types.StringValue(""),
				JWTAuthPath: types.StringValue(""),
			},
			expected: map[string]interface{}{
				"address":       "",
				"namespace":     "",
				"role_name":     "",
				"jwt_auth_path": "",
			},
		},
		{
			name: "TF model with null values",
			tfModel: &model.IntegrationTFModel{
				Address:     types.StringNull(),
				Namespace:   types.StringNull(),
				RoleName:    types.StringNull(),
				JWTAuthPath: types.StringNull(),
			},
			expected: map[string]interface{}{
				"address":       "",
				"namespace":     "",
				"role_name":     "",
				"jwt_auth_path": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &NvaultDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
