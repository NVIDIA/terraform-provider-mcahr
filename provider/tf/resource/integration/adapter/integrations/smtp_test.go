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

func TestSmtpDataAdapter_MapToTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		validate        func(t *testing.T, tfModel *model.IntegrationTFModel)
	}{
		{
			name: "Valid data with all fields",
			integrationData: map[string]interface{}{
				"smtp_host":             "smtp.example.com",
				"smtp_port":             float64(587),
				"username":              "foo@bar.com",
				"password":              "password",
				"sender":                "foo@bar.com",
				"max_emails_per_day":    float64(1000),
				"max_emails_per_second": float64(10),
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "smtp.example.com", tfModel.SmtpHost.ValueString())
				assert.Equal(t, int64(587), tfModel.SmtpPort.ValueInt64())
				assert.Equal(t, "foo@bar.com", tfModel.Username.ValueString())
				assert.Equal(t, "password", tfModel.Password.ValueString())
				assert.Equal(t, "foo@bar.com", tfModel.Sender.ValueString())
				assert.Equal(t, int64(1000), tfModel.MaxEmailsPerDay.ValueInt64())
				assert.Equal(t, int64(10), tfModel.MaxEmailsPerSecond.ValueInt64())
			},
		},
		{
			name: "Valid data with minimal fields",
			integrationData: map[string]interface{}{
				"smtp_host": "smtp.minimal.com",
				"smtp_port": float64(25),
			},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "smtp.minimal.com", tfModel.SmtpHost.ValueString())
				assert.Equal(t, int64(25), tfModel.SmtpPort.ValueInt64())
				assert.Equal(t, "", tfModel.Username.ValueString())
				assert.Equal(t, "", tfModel.Password.ValueString())
				assert.Equal(t, "", tfModel.Sender.ValueString())
				assert.Equal(t, int64(0), tfModel.MaxEmailsPerDay.ValueInt64())
				assert.Equal(t, int64(0), tfModel.MaxEmailsPerSecond.ValueInt64())
			},
		},
		{
			name:            "Empty data",
			integrationData: map[string]interface{}{},
			validate: func(t *testing.T, tfModel *model.IntegrationTFModel) {
				assert.Equal(t, "", tfModel.SmtpHost.ValueString())
				assert.Equal(t, int64(0), tfModel.SmtpPort.ValueInt64())
				assert.Equal(t, "", tfModel.Username.ValueString())
				assert.Equal(t, "", tfModel.Password.ValueString())
				assert.Equal(t, "", tfModel.Sender.ValueString())
				assert.Equal(t, int64(0), tfModel.MaxEmailsPerDay.ValueInt64())
				assert.Equal(t, int64(0), tfModel.MaxEmailsPerSecond.ValueInt64())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			adapter := &SmtpDataAdapter{}
			tfModel := &model.IntegrationTFModel{}
			options := &adapterinterface.IntegrationDataAdapterOptions{
				CompatibilityOptions: make(map[string]attribute.CompatibilityOptions),
			}

			adapter.MapToTFModel(common.NewRequestContext(context.Background()), options, tt.integrationData, tfModel)

			tt.validate(t, tfModel)
		})
	}
}

func TestSmtpDataAdapter_TFModelToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tfModel  *model.IntegrationTFModel
		validate func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "Valid TF model with all fields",
			tfModel: &model.IntegrationTFModel{
				SmtpHost:           types.StringValue("smtp.example.com"),
				SmtpPort:           types.Int64Value(587),
				Username:           types.StringValue("foo@bar.com"),
				Password:           types.StringValue("password"),
				Sender:             types.StringValue("foo@bar.com"),
				MaxEmailsPerDay:    types.Int64Value(1000),
				MaxEmailsPerSecond: types.Int64Value(10),
			},
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "smtp.example.com", result["smtp_host"])
				assert.Equal(t, int64(587), result["smtp_port"])
				assert.Equal(t, "foo@bar.com", result["username"])
				assert.Equal(t, "password", result["password"])
				assert.Equal(t, "foo@bar.com", result["sender"])
				assert.Equal(t, int64(1000), result["max_emails_per_day"])
				assert.Equal(t, int64(10), result["max_emails_per_second"])
			},
		},
		{
			name: "Valid TF model with minimal fields",
			tfModel: &model.IntegrationTFModel{
				SmtpHost: types.StringValue("smtp.minimal.com"),
				SmtpPort: types.Int64Value(25),
			},
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "smtp.minimal.com", result["smtp_host"])
				assert.Equal(t, int64(25), result["smtp_port"])
				assert.Equal(t, "", result["username"])
				assert.Equal(t, "", result["password"])
				assert.Equal(t, "", result["sender"])
				assert.Equal(t, int64(0), result["max_emails_per_day"])
				assert.Equal(t, int64(0), result["max_emails_per_second"])
			},
		},
		{
			name:    "Empty TF model",
			tfModel: &model.IntegrationTFModel{},
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "", result["smtp_host"])
				assert.Equal(t, int64(0), result["smtp_port"])
				assert.Equal(t, "", result["username"])
				assert.Equal(t, "", result["password"])
				assert.Equal(t, "", result["sender"])
				assert.Equal(t, int64(0), result["max_emails_per_day"])
				assert.Equal(t, int64(0), result["max_emails_per_second"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			adapter := &SmtpDataAdapter{}
			options := &adapterinterface.IntegrationDataAdapterOptions{
				CompatibilityOptions: make(map[string]attribute.CompatibilityOptions),
			}

			result := adapter.TFModelToMap(common.NewRequestContext(context.Background()), options, tt.tfModel)

			tt.validate(t, result)
		})
	}
}

func TestSmtpDataAdapter_DataFieldNames(t *testing.T) {
	t.Parallel()

	adapter := &SmtpDataAdapter{}
	fieldNames := adapter.DataFieldNames()

	expectedFields := []string{"smtp_host", "smtp_port", "username", "password", "sender", "max_emails_per_day", "max_emails_per_second"}
	assert.Equal(t, expectedFields, fieldNames)
}

func TestSmtpDataAdapter_TFModelFieldNames(t *testing.T) {
	t.Parallel()

	adapter := &SmtpDataAdapter{}
	fieldNames := adapter.TFModelFieldNames()

	expectedFields := []string{"smtp_host", "smtp_port", "username", "password", "sender", "max_emails_per_day", "max_emails_per_second"}
	assert.Equal(t, expectedFields, fieldNames)
}
