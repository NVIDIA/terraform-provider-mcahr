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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestAlertmanagerDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"external_url":  "https://alertmanager.example.com",
				"payload_paths": []interface{}{"alerts.receiver", "alerts.status"},
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "https://alertmanager.example.com", tfModel.ExternalUrl.ValueString())
				assert.False(t, tfModel.PayloadPaths.IsNull())
				assert.Len(t, tfModel.PayloadPaths.Elements(), 2)
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.ExternalUrl.ValueString())
				assert.True(t, tfModel.PayloadPaths.IsNull() || len(tfModel.PayloadPaths.Elements()) == 0)
			},
		},
		{
			name: "Partial data - only external_url",
			integrationData: map[string]interface{}{
				"external_url": "https://test.example.com",
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "https://test.example.com", tfModel.ExternalUrl.ValueString())
				assert.True(t, tfModel.PayloadPaths.IsNull() || len(tfModel.PayloadPaths.Elements()) == 0)
			},
		},
		{
			name: "Partial data - only payload_paths",
			integrationData: map[string]interface{}{
				"payload_paths": []interface{}{"single.path"},
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.ExternalUrl.ValueString())
				assert.False(t, tfModel.PayloadPaths.IsNull())
				assert.Len(t, tfModel.PayloadPaths.Elements(), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &AlertmanagerDataAdapter{}
			tfModel := &model.IntegrationTFModel{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			adapter.MapToTFModel(common.NewRequestContext(context.Background()), adapterOptions, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestAlertmanagerDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		expected map[string]interface{}
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				ExternalUrl: types.StringValue("https://alertmanager.example.com"),
				PayloadPaths: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("alerts.receiver"),
					types.StringValue("alerts.status"),
				}),
			},
			expected: map[string]interface{}{
				"external_url":  "https://alertmanager.example.com",
				"payload_paths": []string{"alerts.receiver", "alerts.status"},
			},
		},
		{
			name: "Empty TF model",
			tfModel: &model.IntegrationTFModel{
				ExternalUrl:  types.StringValue(""),
				PayloadPaths: types.ListValueMust(types.StringType, []attr.Value{}),
			},
			expected: map[string]interface{}{
				"external_url":  "",
				"payload_paths": []string{},
			},
		},
		{
			name: "TF model with null values",
			tfModel: &model.IntegrationTFModel{
				ExternalUrl:  types.StringNull(),
				PayloadPaths: types.ListNull(types.StringType),
			},
			expected: map[string]interface{}{
				"external_url":  "",
				"payload_paths": []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &AlertmanagerDataAdapter{}

			compatibilityOptions := make(map[string]attribute.CompatibilityOptions)
			adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{CompatibilityOptions: compatibilityOptions}
			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), adapterOptions, tt.tfModel)

			assert.Equal(t, tt.expected, result)
		})
	}
}
