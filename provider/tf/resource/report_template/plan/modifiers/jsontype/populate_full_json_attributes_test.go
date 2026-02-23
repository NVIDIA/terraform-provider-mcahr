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

package jsonmodifier

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/tf/resource/report_template/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// normalizeJSON properly compacts JSON while preserving string content
func normalizeJSON(jsonStr string) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(jsonStr)); err != nil {
		// If compact fails, return original string
		return jsonStr
	}
	return buf.String()
}

func TestPopulateFullJsonAttributes(t *testing.T) {
	tests := []struct {
		name           string
		resultValues   *model.ReportTemplateTFModel
		plan           *model.ReportTemplateTFModel
		state          *model.ReportTemplateTFModel
		backendVersion *version.BackendVersion
		expectError    bool
		validate       func(t *testing.T, result *model.ReportTemplateTFModel)
	}{
		{
			name: "Process blocks field",
			resultValues: &model.ReportTemplateTFModel{
				Blocks: types.StringValue(`[{"title": "Test Block", "resource_query": "resource.type = 'server'"}]`),
			},
			plan: &model.ReportTemplateTFModel{
				BlocksFull: types.StringValue("not_null"),
			},
			state:          &model.ReportTemplateTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.False(t, result.BlocksFull.IsNull())
				expected := `[{
					"breakdown_by_tag": "",
					"breakdown_tags_values": [],
					"group_by_tag": "",
					"group_by_tag_order": {
						"type": "DEFAULT",
						"values": []
					},
					"include_other_breakdown_tag_values": false,
					"include_resources_without_group_tag": false,
					"other_tags_to_export": [],
					"resource_query": "resource.type = 'server'",
					"resources_breakdown": [],
					"title": "Test Block",
					"view_mode": "COUNT"
				}]`
				assert.Equal(t, normalizeJSON(expected), result.BlocksFull.ValueString())
			},
		},
		{
			name: "Process links field",
			resultValues: &model.ReportTemplateTFModel{
				Links: types.StringValue(`[{"label": "Test Link", "report_template_name": "target_template"}]`),
			},
			plan: &model.ReportTemplateTFModel{
				LinksFull: types.StringValue("not_null"),
			},
			state:          &model.ReportTemplateTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.False(t, result.LinksFull.IsNull())
				expected := `[{
					"label": "Test Link",
					"report_template_name": "target_template"
				}]`
				assert.Equal(t, normalizeJSON(expected), result.LinksFull.ValueString())
			},
		},
		{
			name: "Skip delete operation - null blocks_full plan",
			resultValues: &model.ReportTemplateTFModel{
				Blocks: types.StringValue(`[{"title": "test"}]`),
			},
			plan: &model.ReportTemplateTFModel{
				BlocksFull: types.StringNull(), // Indicates delete
			},
			state:          &model.ReportTemplateTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.True(t, result.BlocksFull.IsNull())
			},
		},
		{
			name: "Skip delete operation - null links_full plan",
			resultValues: &model.ReportTemplateTFModel{
				Links: types.StringValue(`[{"label": "test"}]`),
			},
			plan: &model.ReportTemplateTFModel{
				LinksFull: types.StringNull(), // Indicates delete
			},
			state:          &model.ReportTemplateTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.True(t, result.LinksFull.IsNull())
			},
		},
		{
			name: "Handle invalid JSON in blocks",
			resultValues: &model.ReportTemplateTFModel{
				Blocks: types.StringValue(`{invalid json}`),
			},
			plan: &model.ReportTemplateTFModel{
				BlocksFull: types.StringValue("not_null"),
			},
			state:          &model.ReportTemplateTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    true,
		},
		{
			name: "Handle invalid JSON in links",
			resultValues: &model.ReportTemplateTFModel{
				Links: types.StringValue(`{invalid json}`),
			},
			plan: &model.ReportTemplateTFModel{
				LinksFull: types.StringValue("not_null"),
			},
			state:          &model.ReportTemplateTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    true,
		},
		{
			name: "Process both blocks and links fields",
			resultValues: &model.ReportTemplateTFModel{
				Blocks: types.StringValue(`[{"title": "Block 1", "resource_query": "query1"}]`),
				Links:  types.StringValue(`[{"label": "Link 1", "report_template_name": "template1"}]`),
			},
			plan: &model.ReportTemplateTFModel{
				BlocksFull: types.StringValue("not_null"),
				LinksFull:  types.StringValue("not_null"),
			},
			state:          &model.ReportTemplateTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.False(t, result.BlocksFull.IsNull())
				assert.False(t, result.LinksFull.IsNull())
				expectedBlocks := `[{
					"breakdown_by_tag": "",
					"breakdown_tags_values": [],
					"group_by_tag": "",
					"group_by_tag_order": {
						"type": "DEFAULT",
						"values": []
					},
					"include_other_breakdown_tag_values": false,
					"include_resources_without_group_tag": false,
					"other_tags_to_export": [],
					"resource_query": "query1",
					"resources_breakdown": [],
					"title": "Block 1",
					"view_mode": "COUNT"
				}]`
				expectedLinks := `[{
					"label": "Link 1",
					"report_template_name": "template1"
				}]`
				assert.Equal(t, normalizeJSON(expectedBlocks), result.BlocksFull.ValueString())
				assert.Equal(t, normalizeJSON(expectedLinks), result.LinksFull.ValueString())
			},
		},
		{
			name: "Process complex blocks with default values",
			resultValues: &model.ReportTemplateTFModel{
				Blocks: types.StringValue(`[{
					"title": "Complex Block",
					"resource_query": "resource.type = 'server'",
					"group_by_tag": "environment",
					"breakdown_by_tag": "region"
				}]`),
			},
			plan: &model.ReportTemplateTFModel{
				BlocksFull: types.StringValue("not_null"),
			},
			state:          &model.ReportTemplateTFModel{},
			backendVersion: &version.BackendVersion{Version: "release-29.0.0", Major: 29, Minor: 0, Patch: 0},
			expectError:    false,
			validate: func(t *testing.T, result *model.ReportTemplateTFModel) {
				assert.False(t, result.BlocksFull.IsNull())
				expected := `[{
					"breakdown_by_tag": "region",
					"breakdown_tags_values": [],
					"group_by_tag": "environment",
					"group_by_tag_order": {
						"type": "DEFAULT",
						"values": []
					},
					"include_other_breakdown_tag_values": false,
					"include_resources_without_group_tag": false,
					"other_tags_to_export": [],
					"resource_query": "resource.type = 'server'",
					"resources_breakdown": [],
					"title": "Complex Block",
					"view_mode": "COUNT"
				}]`
				assert.Equal(t, normalizeJSON(expected), result.BlocksFull.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// given
			ctx := context.Background()

			// when
			err := PopulateFullJsonAttributes(ctx, tt.resultValues, tt.plan, tt.state, tt.backendVersion)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, tt.resultValues)
				}
			}
		})
	}
}

func TestIsDeleteOperation(t *testing.T) {
	tests := []struct {
		name       string
		plan       *model.ReportTemplateTFModel
		attrConfig JsonAttributeConfig
		expected   bool
	}{
		{
			name: "Delete operation - null blocks_full",
			plan: &model.ReportTemplateTFModel{
				BlocksFull: types.StringNull(),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["blocks"],
			expected:   true,
		},
		{
			name: "Not delete operation - non-null blocks_full",
			plan: &model.ReportTemplateTFModel{
				BlocksFull: types.StringValue("some_value"),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["blocks"],
			expected:   false,
		},
		{
			name: "Delete operation - null links_full",
			plan: &model.ReportTemplateTFModel{
				LinksFull: types.StringNull(),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["links"],
			expected:   true,
		},
		{
			name: "Not delete operation - empty string",
			plan: &model.ReportTemplateTFModel{
				BlocksFull: types.StringValue(""),
			},
			attrConfig: JSON_ATTRIBUTES_TO_POPULATE["blocks"],
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := isDeleteOperation(tt.plan, tt.attrConfig)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJsonAttributeConfig_Functions(t *testing.T) {
	tests := []struct {
		name        string
		attribute   string
		setupModel  func() *model.ReportTemplateTFModel
		testGetters bool
		testSetters bool
	}{
		{
			name:      "Blocks attribute config",
			attribute: "blocks",
			setupModel: func() *model.ReportTemplateTFModel {
				return &model.ReportTemplateTFModel{
					Blocks:     types.StringValue("blocks_value"),
					BlocksFull: types.StringValue("blocks_full_value"),
				}
			},
			testGetters: true,
			testSetters: true,
		},
		{
			name:      "Links attribute config",
			attribute: "links",
			setupModel: func() *model.ReportTemplateTFModel {
				return &model.ReportTemplateTFModel{
					Links:     types.StringValue("links_value"),
					LinksFull: types.StringValue("links_full_value"),
				}
			},
			testGetters: true,
			testSetters: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// given
			attrConfig := JSON_ATTRIBUTES_TO_POPULATE[tt.attribute]
			model := tt.setupModel()

			if tt.testGetters {
				// Test GetAttr and GetFullAttr
				regularValue := attrConfig.GetAttr(model)
				fullValue := attrConfig.GetFullAttr(model)

				switch tt.attribute {
				case "blocks":
					assert.Equal(t, "blocks_value", regularValue.ValueString())
					assert.Equal(t, "blocks_full_value", fullValue.ValueString())
				case "links":
					assert.Equal(t, "links_value", regularValue.ValueString())
					assert.Equal(t, "links_full_value", fullValue.ValueString())
				}
			}

			if tt.testSetters {
				// Test SetFullAttr
				newValue := types.StringValue("new_full_value")
				attrConfig.SetFullAttr(model, newValue)

				fullValue := attrConfig.GetFullAttr(model)
				assert.Equal(t, "new_full_value", fullValue.ValueString())
			}
		})
	}
}

func TestJSON_ATTRIBUTES_TO_POPULATE_Configuration(t *testing.T) {
	// Verify all expected attributes are configured
	expectedAttributes := []string{"blocks", "links"}

	for _, attr := range expectedAttributes {
		t.Run("Has_"+attr, func(t *testing.T) {
			attrConfig, exists := JSON_ATTRIBUTES_TO_POPULATE[attr]
			assert.True(t, exists, "Expected attribute %s to be in JSON_ATTRIBUTES_TO_POPULATE", attr)

			// Verify full attribute name
			switch attr {
			case "blocks":
				assert.Equal(t, "blocks_full", attrConfig.FullAttrName)
			case "links":
				assert.Equal(t, "links_full", attrConfig.FullAttrName)
			}

			// Verify functions are set
			assert.NotNil(t, attrConfig.RemarshalFunc)
			assert.NotNil(t, attrConfig.GetAttr)
			assert.NotNil(t, attrConfig.GetFullAttr)
			assert.NotNil(t, attrConfig.SetFullAttr)
		})
	}

	// Verify no extra attributes
	assert.Equal(t, len(expectedAttributes), len(JSON_ATTRIBUTES_TO_POPULATE))
}

func TestPopulateFullJsonAttributes_Integration(t *testing.T) {
	// Integration test with real JSON processing
	ctx := context.Background()
	backendVersion := &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1}

	resultValues := &model.ReportTemplateTFModel{
		Blocks: types.StringValue(`[
			{
				"title": "Server Resources",
				"resource_query": "resource.type = 'server'",
				"group_by_tag": "environment",
				"breakdown_by_tag": "region",
				"view_mode": "COUNT"
			},
			{
				"title": "Database Resources", 
				"resource_query": "resource.type = 'database'",
				"group_by_tag": "tier"
			}
		]`),
		Links: types.StringValue(`[
			{
				"label": "Server Dashboard",
				"report_template_name": "server_report"
			},
			{
				"label": "Database Overview",
				"report_template_name": "db_report"
			}
		]`),
	}

	plan := &model.ReportTemplateTFModel{
		BlocksFull: types.StringValue("not_null"),
		LinksFull:  types.StringValue("not_null"),
	}

	state := &model.ReportTemplateTFModel{}

	// when
	err := PopulateFullJsonAttributes(ctx, resultValues, plan, state, backendVersion)

	// then
	require.NoError(t, err)

	// Verify blocks_full was populated
	assert.False(t, resultValues.BlocksFull.IsNull())
	expectedBlocks := `[{
		"breakdown_by_tag": "region",
		"breakdown_tags_values": [],
		"group_by_tag": "environment",
		"group_by_tag_order": {
			"type": "DEFAULT",
			"values": []
		},
		"include_other_breakdown_tag_values": false,
		"include_resources_without_group_tag": false,
		"other_tags_to_export": [],
		"resource_query": "resource.type = 'server'",
		"resources_breakdown": [],
		"title": "Server Resources",
		"view_mode": "COUNT"
	}, {
		"breakdown_by_tag": "",
		"breakdown_tags_values": [],
		"group_by_tag": "tier",
		"group_by_tag_order": {
			"type": "DEFAULT",
			"values": []
		},
		"include_other_breakdown_tag_values": false,
		"include_resources_without_group_tag": false,
		"other_tags_to_export": [],
		"resource_query": "resource.type = 'database'",
		"resources_breakdown": [],
		"title": "Database Resources",
		"view_mode": "COUNT"
	}]`
	assert.Equal(t, normalizeJSON(expectedBlocks), resultValues.BlocksFull.ValueString())

	// Verify links_full was populated
	assert.False(t, resultValues.LinksFull.IsNull())
	expectedLinks := `[{
		"label": "Server Dashboard",
		"report_template_name": "server_report"
	}, {
		"label": "Database Overview",
		"report_template_name": "db_report"
	}]`
	assert.Equal(t, normalizeJSON(expectedLinks), resultValues.LinksFull.ValueString())
}
