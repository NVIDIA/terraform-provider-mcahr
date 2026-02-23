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
	"reflect"
	"testing"

	"terraform/terraform-provider/provider/tf/core/plan"
	"terraform/terraform-provider/provider/tf/resource/report_template/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestAddDefaultsFromPlan(t *testing.T) {
	tests := []struct {
		name         string
		resultValues *model.ReportTemplateTFModel
		planValues   *model.ReportTemplateTFModel
		validate     func(t *testing.T, result *model.ReportTemplateTFModel)
	}{
		{
			name: "Copy null fields from plan",
			resultValues: &model.ReportTemplateTFModel{
				Name:   types.StringValue("test"),
				Blocks: types.StringNull(),
			},
			planValues: &model.ReportTemplateTFModel{
				Name:   types.StringValue("test"),
				Blocks: types.StringValue(`[{"title":"Block"}]`),
			},
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.Equal(t, "test", result.Name.ValueString())
				assert.Equal(t, `[{"title":"Block"}]`, result.Blocks.ValueString())
			},
		},
		{
			name: "Copy unknown fields from plan",
			resultValues: &model.ReportTemplateTFModel{
				Name:  types.StringValue("test"),
				Links: types.StringUnknown(),
			},
			planValues: &model.ReportTemplateTFModel{
				Name:  types.StringValue("test"),
				Links: types.StringValue(`[{"label":"link"}]`),
			},
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.Equal(t, "test", result.Name.ValueString())
				assert.Equal(t, `[{"label":"link"}]`, result.Links.ValueString())
			},
		},
		{
			name: "Don't override non-null/unknown fields",
			resultValues: &model.ReportTemplateTFModel{
				Name:   types.StringValue("result"),
				Blocks: types.StringValue(`[{"title":"Result"}]`),
			},
			planValues: &model.ReportTemplateTFModel{
				Name:   types.StringValue("plan"),
				Blocks: types.StringValue(`[{"title":"Plan"}]`),
			},
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.Equal(t, "result", result.Name.ValueString())
				assert.Equal(t, `[{"title":"Result"}]`, result.Blocks.ValueString())
			},
		},
		{
			name: "Handle all report template field types",
			resultValues: &model.ReportTemplateTFModel{
				Name:       types.StringNull(),
				Blocks:     types.StringNull(),
				BlocksFull: types.StringNull(),
				Links:      types.StringNull(),
				LinksFull:  types.StringNull(),
			},
			planValues: &model.ReportTemplateTFModel{
				Name:       types.StringValue("from plan"),
				Blocks:     types.StringValue(`[{"title":"Plan Block"}]`),
				BlocksFull: types.StringValue(`[{"title":"Plan Block Full"}]`),
				Links:      types.StringValue(`[{"label":"Plan Link"}]`),
				LinksFull:  types.StringValue(`[{"label":"Plan Link Full"}]`),
			},
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.Equal(t, "from plan", result.Name.ValueString())
				assert.Equal(t, `[{"title":"Plan Block"}]`, result.Blocks.ValueString())
				assert.Equal(t, `[{"title":"Plan Block Full"}]`, result.BlocksFull.ValueString())
				assert.Equal(t, `[{"label":"Plan Link"}]`, result.Links.ValueString())
				assert.Equal(t, `[{"label":"Plan Link Full"}]`, result.LinksFull.ValueString())
			},
		},
		{
			name: "Mixed null, unknown, and known fields",
			resultValues: &model.ReportTemplateTFModel{
				Name:       types.StringValue("existing"), // Known - should not be overridden
				Blocks:     types.StringNull(),            // Null - should be copied from plan
				BlocksFull: types.StringUnknown(),         // Unknown - should be copied from plan
				Links:      types.StringValue(`[]`),       // Known - should not be overridden
				LinksFull:  types.StringNull(),            // Null - should be copied from plan
			},
			planValues: &model.ReportTemplateTFModel{
				Name:       types.StringValue("plan name"),
				Blocks:     types.StringValue(`[{"title":"Plan"}]`),
				BlocksFull: types.StringValue(`[{"title":"Plan Full"}]`),
				Links:      types.StringValue(`[{"label":"Plan"}]`),
				LinksFull:  types.StringValue(`[{"label":"Plan Full"}]`),
			},
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.Equal(t, "existing", result.Name.ValueString())                      // Not overridden
				assert.Equal(t, `[{"title":"Plan"}]`, result.Blocks.ValueString())          // Copied from plan
				assert.Equal(t, `[{"title":"Plan Full"}]`, result.BlocksFull.ValueString()) // Copied from plan
				assert.Equal(t, `[]`, result.Links.ValueString())                           // Not overridden
				assert.Equal(t, `[{"label":"Plan Full"}]`, result.LinksFull.ValueString())  // Copied from plan
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// when
			plan.AddDefaultsFromPlan(tt.resultValues, tt.planValues)

			// then
			tt.validate(t, tt.resultValues)
		})
	}
}

func TestCheckIfFieldIsNullOrUnknown(t *testing.T) {
	tests := []struct {
		name     string
		field    types.String
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// when
			result := plan.CheckIfFieldIsNullOrUnknown(reflect.ValueOf(tt.field))

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}
