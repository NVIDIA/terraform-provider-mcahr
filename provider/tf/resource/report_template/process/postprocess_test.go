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

package process

import (
	"context"
	"encoding/json"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	reporttemplatetf "terraform/terraform-provider/provider/tf/resource/report_template/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostProcessJsonFullFields_WithNullFields tests JSON field processing with null values
func TestPostProcessJsonFullFields_WithNullFields(t *testing.T) {
	t.Parallel()
	// given
	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name:       types.StringValue("test_report_template"),
		Blocks:     types.StringValue(`[{"title":"Block"}]`),
		BlocksFull: types.StringNull(),
		Links:      types.StringValue(`[{"label":"Link"}]`),
		LinksFull:  types.StringNull(),
	}
	backendVersion := &version.BackendVersion{Version: "release-2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	// when
	err := postProcessJsonFullFields(requestContext, tfModel)

	// then
	require.NoError(t, err)
	assert.True(t, tfModel.BlocksFull.IsNull())
	assert.True(t, tfModel.LinksFull.IsNull())
}

// TestPostProcessJsonFullFields_WithUnknownFields tests JSON field processing with unknown values
func TestPostProcessJsonFullFields_WithUnknownFields(t *testing.T) {
	t.Parallel()
	// given
	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name:       types.StringValue("test_report_template"),
		Blocks:     types.StringValue(`[{"title":"Block"}]`),
		BlocksFull: types.StringUnknown(),
		Links:      types.StringValue(`[{"label":"Link"}]`),
		LinksFull:  types.StringUnknown(),
	}
	backendVersion := &version.BackendVersion{Version: "release-2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	// when
	err := postProcessJsonFullFields(requestContext, tfModel)

	// then
	require.NoError(t, err)
	assert.True(t, tfModel.BlocksFull.IsUnknown())
	assert.True(t, tfModel.LinksFull.IsUnknown())
}

// TestPostProcessJsonFullFields_WithValidJSON tests JSON field processing with valid JSON
func TestPostProcessJsonFullFields_WithValidJSON(t *testing.T) {
	t.Parallel()
	// given
	blocksJSON := `[{"title":"Block Name","resource_query":"host","group_by_tag":"tag_0"}]`
	linksJSON := `[{"label":"Test Link","report_template_name":"other_template"}]`

	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name:       types.StringValue("test_report_template"),
		Blocks:     types.StringValue(blocksJSON),
		BlocksFull: types.StringValue(blocksJSON),
		Links:      types.StringValue(linksJSON),
		LinksFull:  types.StringValue(linksJSON),
	}
	backendVersion := &version.BackendVersion{Version: "release-2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	// when
	err := postProcessJsonFullFields(requestContext, tfModel)

	// then
	require.NoError(t, err)
	// Values should be processed and remain valid JSON with defaults applied
	expectedBlocksJSON := `[{` +
		`"breakdown_by_tag":"",` +
		`"breakdown_tags_values":[],` +
		`"group_by_tag":"tag_0",` +
		`"group_by_tag_order":{"type":"DEFAULT","values":[]},` +
		`"include_other_breakdown_tag_values":false,` +
		`"include_resources_without_group_tag":false,` +
		`"other_tags_to_export":[],` +
		`"resource_query":"host",` +
		`"resources_breakdown":[],` +
		`"title":"Block Name",` +
		`"view_mode":"COUNT"` +
		`}]`
	assert.Equal(t, expectedBlocksJSON, tfModel.BlocksFull.ValueString())

	expectedLinksJSON := `[{` +
		`"label":"Test Link",` +
		`"report_template_name":"other_template"` +
		`}]`
	assert.Equal(t, expectedLinksJSON, tfModel.LinksFull.ValueString())

	// Verify the JSON is still valid
	var blocks []interface{}
	err = json.Unmarshal([]byte(tfModel.BlocksFull.ValueString()), &blocks)
	assert.NoError(t, err)

	var links []interface{}
	err = json.Unmarshal([]byte(tfModel.LinksFull.ValueString()), &links)
	assert.NoError(t, err)
}

// TestPostProcessJsonFullFields_WithInvalidJSON tests JSON field processing with invalid JSON
func TestPostProcessJsonFullFields_WithInvalidJSON(t *testing.T) {
	t.Parallel()
	// given
	invalidJSON := `[{"title":"Block","invalid":}]` // Invalid JSON

	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name:       types.StringValue("test_report_template"),
		Blocks:     types.StringValue(invalidJSON),
		BlocksFull: types.StringValue(invalidJSON),
		Links:      types.StringValue(`[]`),
		LinksFull:  types.StringValue(`[]`),
	}
	backendVersion := &version.BackendVersion{Version: "release-2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	// when
	err := postProcessJsonFullFields(requestContext, tfModel)

	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

// TestPostProcessJsonFullField_GenericType tests the generic postProcessJsonFullField function
func TestPostProcessJsonFullField_GenericType(t *testing.T) {
	t.Parallel()
	// given
	backendVersion := &version.BackendVersion{Version: "release-2.0.0"}

	tests := []struct {
		name        string
		input       types.String
		expectError bool
		expectNull  bool
	}{
		{
			name:        "Null value",
			input:       types.StringNull(),
			expectError: false,
			expectNull:  true,
		},
		{
			name:        "Unknown value",
			input:       types.StringUnknown(),
			expectError: false,
			expectNull:  false,
		},
		{
			name:        "Valid JSON",
			input:       types.StringValue(`[{"title":"Test"}]`),
			expectError: false,
			expectNull:  false,
		},
		{
			name:        "Invalid JSON",
			input:       types.StringValue(`[{"invalid":}]`),
			expectError: true,
			expectNull:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := postProcessJsonFullField[*MockJsonConfigurable](&tt.input, backendVersion)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectNull {
					assert.True(t, result.IsNull())
				} else {
					assert.False(t, result.IsNull())
				}
			}
		})
	}
}

// Mock implementation of JsonConfigurable for testing
type MockJsonConfigurable struct {
	Config common.JsonConfig
	Name   string `json:"name"`
	Type   string `json:"type"`
}

func (m *MockJsonConfigurable) SetConfig(config common.JsonConfig) {
	m.Config = config
}

func (m *MockJsonConfigurable) GetConfig() common.JsonConfig {
	return m.Config
}
