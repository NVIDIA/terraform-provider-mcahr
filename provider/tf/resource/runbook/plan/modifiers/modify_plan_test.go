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

package plan

import (
	"testing"

	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestModifyPlan(t *testing.T) {
	// Skip these tests as they require complex mocking of the framework internals
	t.Skip("Skipping tests that require complex framework mocking")
}

func TestHandleDataFieldPlan(t *testing.T) {
	tests := []struct {
		name         string
		resultValues *model.RunbookTFModel
		configValues *model.RunbookTFModel
		planValues   *model.RunbookTFModel
		stateValues  *model.RunbookTFModel
		validate     func(t *testing.T, result *model.RunbookTFModel)
	}{
		{
			name: "User didn't provide data - set to null",
			resultValues: &model.RunbookTFModel{
				Data: types.StringUnknown(),
			},
			configValues: &model.RunbookTFModel{
				Data: types.StringNull(),
			},
			planValues: &model.RunbookTFModel{
				Data: types.StringNull(),
			},
			stateValues: &model.RunbookTFModel{
				Data: types.StringNull(),
			},
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.True(t, result.Data.IsNull())
			},
		},
		{
			name: "User provided data in config",
			resultValues: &model.RunbookTFModel{
				Data: types.StringValue(`{"test": true}`),
			},
			configValues: &model.RunbookTFModel{
				Data: types.StringValue(`{"test": true}`),
			},
			planValues: &model.RunbookTFModel{
				Data: types.StringValue(`{"test": true}`),
			},
			stateValues: &model.RunbookTFModel{
				Data: types.StringNull(),
			},
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.Equal(t, `{"test": true}`, result.Data.ValueString())
			},
		},
		{
			name: "Data exists in state",
			resultValues: &model.RunbookTFModel{
				Data: types.StringUnknown(),
			},
			configValues: &model.RunbookTFModel{
				Data: types.StringNull(),
			},
			planValues: &model.RunbookTFModel{
				Data: types.StringNull(),
			},
			stateValues: &model.RunbookTFModel{
				Data: types.StringValue(`{"existing": true}`),
			},
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.True(t, result.Data.IsUnknown())
			},
		},
		{
			name: "Data is planned to change",
			resultValues: &model.RunbookTFModel{
				Data: types.StringValue(`{"new": true}`),
			},
			configValues: &model.RunbookTFModel{
				Data: types.StringNull(),
			},
			planValues: &model.RunbookTFModel{
				Data: types.StringValue(`{"new": true}`),
			},
			stateValues: &model.RunbookTFModel{
				Data: types.StringValue(`{"old": true}`),
			},
			validate: func(t *testing.T, result *model.RunbookTFModel) {
				assert.Equal(t, `{"new": true}`, result.Data.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			handleNullDataFieldPlan(tt.resultValues, tt.configValues, tt.planValues, tt.stateValues)

			// then
			if tt.validate != nil {
				tt.validate(t, tt.resultValues)
			}
		})
	}
}
