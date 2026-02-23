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
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestModel is a generic test model for testing plan utilities
type TestModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	TimeoutMs   types.Int64  `tfsdk:"timeout_ms"`
	Labels      types.List   `tfsdk:"labels"`
}

func TestAddDefaultsFromPlan(t *testing.T) {
	tests := []struct {
		name         string
		resultValues *TestModel
		planValues   *TestModel
		validate     func(t *testing.T, result *TestModel)
	}{
		{
			name: "Copy null fields from plan",
			resultValues: &TestModel{
				Name:    types.StringValue("test"),
				Enabled: types.BoolNull(),
			},
			planValues: &TestModel{
				Name:    types.StringValue("test"),
				Enabled: types.BoolValue(true),
			},
			validate: func(t *testing.T, result *TestModel) {
				assert.Equal(t, "test", result.Name.ValueString())
				assert.Equal(t, true, result.Enabled.ValueBool())
			},
		},
		{
			name: "Copy unknown fields from plan",
			resultValues: &TestModel{
				Name:        types.StringValue("test"),
				Description: types.StringUnknown(),
			},
			planValues: &TestModel{
				Name:        types.StringValue("test"),
				Description: types.StringValue("from plan"),
			},
			validate: func(t *testing.T, result *TestModel) {
				assert.Equal(t, "test", result.Name.ValueString())
				assert.Equal(t, "from plan", result.Description.ValueString())
			},
		},
		{
			name: "Don't override non-null/unknown fields",
			resultValues: &TestModel{
				Name:        types.StringValue("result"),
				Description: types.StringValue("result desc"),
			},
			planValues: &TestModel{
				Name:        types.StringValue("plan"),
				Description: types.StringValue("plan desc"),
			},
			validate: func(t *testing.T, result *TestModel) {
				assert.Equal(t, "result", result.Name.ValueString())
				assert.Equal(t, "result desc", result.Description.ValueString())
			},
		},
		{
			name: "Handle all field types",
			resultValues: &TestModel{
				Name:      types.StringNull(),
				Enabled:   types.BoolNull(),
				TimeoutMs: types.Int64Null(),
				Labels:    types.ListNull(types.StringType),
			},
			planValues: &TestModel{
				Name:      types.StringValue("from plan"),
				Enabled:   types.BoolValue(true),
				TimeoutMs: types.Int64Value(5000),
				Labels: func() types.List {
					s, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"label1"})
					return s
				}(),
			},
			validate: func(t *testing.T, result *TestModel) {
				assert.Equal(t, "from plan", result.Name.ValueString())
				assert.Equal(t, true, result.Enabled.ValueBool())
				assert.Equal(t, int64(5000), result.TimeoutMs.ValueInt64())
				assert.False(t, result.Labels.IsNull())
			},
		},
		{
			name: "Mixed null, unknown, and known fields",
			resultValues: &TestModel{
				Name:        types.StringNull(),     // Should be copied from plan
				Description: types.StringUnknown(),  // Should be copied from plan
				Enabled:     types.BoolValue(false), // Should stay as-is
				TimeoutMs:   types.Int64Value(1000), // Should stay as-is
			},
			planValues: &TestModel{
				Name:        types.StringValue("plan name"),
				Description: types.StringValue("plan desc"),
				Enabled:     types.BoolValue(true),
				TimeoutMs:   types.Int64Value(2000),
			},
			validate: func(t *testing.T, result *TestModel) {
				assert.Equal(t, "plan name", result.Name.ValueString())
				assert.Equal(t, "plan desc", result.Description.ValueString())
				assert.Equal(t, false, result.Enabled.ValueBool())          // Not copied
				assert.Equal(t, int64(1000), result.TimeoutMs.ValueInt64()) // Not copied
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// when
			AddDefaultsFromPlan(tt.resultValues, tt.planValues)

			// then
			if tt.validate != nil {
				tt.validate(t, tt.resultValues)
			}
		})
	}
}

func TestCheckIfFieldIsNullOrUnknown(t *testing.T) {
	tests := []struct {
		name     string
		field    interface{}
		expected bool
	}{
		{
			name:     "Null field should return true",
			field:    types.StringNull(),
			expected: true,
		},
		{
			name:     "Unknown field should return true",
			field:    types.StringUnknown(),
			expected: true,
		},
		{
			name:     "Known field should return false",
			field:    types.StringValue("test"),
			expected: false,
		},
		{
			name:     "Bool null field should return true",
			field:    types.BoolNull(),
			expected: true,
		},
		{
			name:     "Bool unknown field should return true",
			field:    types.BoolUnknown(),
			expected: true,
		},
		{
			name:     "Bool known field should return false",
			field:    types.BoolValue(true),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Use reflection to simulate the internal call
			fieldValue := reflect.ValueOf(tt.field)

			// when
			result := CheckIfFieldIsNullOrUnknown(fieldValue)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}
