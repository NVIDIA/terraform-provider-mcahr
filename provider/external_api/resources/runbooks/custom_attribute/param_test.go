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

package customattribute

import (
	"encoding/json"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParamJson_SetGetConfig(t *testing.T) {
	// given
	param := &ParamJson{}
	config := common.JsonConfig{
		BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
	}

	// when
	param.SetConfig(config)
	result := param.GetConfig()

	// then
	assert.Equal(t, config.BackendVersion.Version, result.BackendVersion.Version)
}

func TestParamJson_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       *ParamJson
		expectError bool
	}{
		{
			name: "Complete param with description (version 28.4.0+)",
			input: &ParamJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
				},
				Name:        "test_param",
				Value:       "test_value",
				Required:    true,
				Export:      true,
				ParamType:   "PARAM",
				Description: "Test parameter",
			},
			expectError: false,
		},
		{
			name: "Param without description (older version)",
			input: &ParamJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-28.3.0", Major: 28, Minor: 3, Patch: 0},
				},
				Name:        "test_param",
				Value:       "test_value",
				Required:    false,
				Export:      false,
				ParamType:   "PARAM",
				Description: "Should be excluded",
			},
			expectError: false,
		},
		{
			name: "Param with empty values",
			input: &ParamJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
				},
				Name:        "",
				Value:       "",
				Required:    false,
				Export:      false,
				ParamType:   "PARAM",
				Description: "",
			},
			expectError: false,
		},
		{
			name: "Param with nil backend version",
			input: &ParamJson{
				Config: common.JsonConfig{
					BackendVersion: nil,
				},
				Name:      "test_param",
				Value:     "test_value",
				ParamType: "PARAM",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := tt.input.MarshalJSON()

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)

				var unmarshaled map[string]interface{}
				err = json.Unmarshal(result, &unmarshaled)
				require.NoError(t, err)

				// Basic field verification
				assert.Contains(t, unmarshaled, "name")
				assert.Contains(t, unmarshaled, "value")
				assert.Contains(t, unmarshaled, "required")
				assert.Contains(t, unmarshaled, "export")
				assert.Contains(t, unmarshaled, "param_type")

				// Check version filtering for description
				if tt.input.Config.BackendVersion != nil &&
					tt.input.Config.BackendVersion.Major == 28 &&
					tt.input.Config.BackendVersion.Minor < 4 {
					assert.NotContains(t, unmarshaled, "description")
				} else if tt.input.Config.BackendVersion != nil &&
					tt.input.Config.BackendVersion.Major >= 28 &&
					tt.input.Config.BackendVersion.Minor >= 4 {
					assert.Contains(t, unmarshaled, "description")
				}
			}
		})
	}
}

func TestParamJson_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		config      common.JsonConfig
		expected    *ParamJson
		expectError bool
	}{
		{
			name: "Complete param JSON",
			input: `{
				"name": "test_param",
				"value": "test_value",
				"required": true,
				"export": true,
				"param_type": "PARAM",
				"description": "Test parameter"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			expected: &ParamJson{
				Name:        "test_param",
				Value:       "test_value",
				Required:    true,
				Export:      true,
				ParamType:   "PARAM",
				Description: "Test parameter",
			},
			expectError: false,
		},
		{
			name: "Partial param JSON with defaults",
			input: `{
				"name": "minimal_param",
				"param_type": "CUSTOM"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			expected: &ParamJson{
				Name:        "minimal_param",
				Value:       DefaultParamValue,
				Required:    DefaultParamRequired,
				Export:      DefaultParamExport,
				ParamType:   "CUSTOM",
				Description: DefaultParamDescription,
			},
			expectError: false,
		},
		{
			name:  "Empty JSON object",
			input: `{}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			expected: &ParamJson{
				Name:        DefaultParamName,
				Value:       DefaultParamValue,
				Required:    DefaultParamRequired,
				Export:      DefaultParamExport,
				ParamType:   DefaultParamType,
				Description: DefaultParamDescription,
			},
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			input:       `{invalid json}`,
			config:      common.JsonConfig{},
			expectError: true,
		},
		{
			name: "Param with version filtering",
			input: `{
				"name": "version_test",
				"param_type": "PARAM",
				"description": "Should be filtered"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.3.0", Major: 28, Minor: 3, Patch: 0},
			},
			expected: &ParamJson{
				Name:        "version_test",
				Value:       DefaultParamValue,
				Required:    DefaultParamRequired,
				Export:      DefaultParamExport,
				ParamType:   "PARAM",
				Description: DefaultParamDescription, // Description should not be set due to version
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			param := &ParamJson{Config: tt.config}

			// when
			err := param.UnmarshalJSON([]byte(tt.input))

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Name, param.Name)
				assert.Equal(t, tt.expected.Value, param.Value)
				assert.Equal(t, tt.expected.Required, param.Required)
				assert.Equal(t, tt.expected.Export, param.Export)
				assert.Equal(t, tt.expected.ParamType, param.ParamType)
				assert.Equal(t, tt.expected.Description, param.Description)
			}
		})
	}
}

func TestParamJson_UnmarshalJSON_FieldPriority(t *testing.T) {
	// given
	input := `{
		"name": "override_test",
		"value": "override_value",
		"required": true,
		"export": false,
		"param_type": "CUSTOM",
		"description": "Override description"
	}`

	config := common.JsonConfig{
		BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
	}

	param := &ParamJson{Config: config}

	// when
	err := param.UnmarshalJSON([]byte(input))

	// then
	require.NoError(t, err)
	assert.Equal(t, "override_test", param.Name)
	assert.Equal(t, "override_value", param.Value)
	assert.Equal(t, true, param.Required)
	assert.Equal(t, false, param.Export)
	assert.Equal(t, "CUSTOM", param.ParamType)
	assert.Equal(t, "Override description", param.Description)
}

func TestParamJson_MarshalUnmarshal_RoundTrip(t *testing.T) {
	// given
	original := &ParamJson{
		Config: common.JsonConfig{
			BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
		},
		Name:        "roundtrip_param",
		Value:       "roundtrip_value",
		Required:    true,
		Export:      true,
		ParamType:   "PARAM",
		Description: "Roundtrip test",
	}

	// when
	marshaled, err := original.MarshalJSON()
	require.NoError(t, err)

	result := &ParamJson{Config: original.Config}
	err = result.UnmarshalJSON(marshaled)
	require.NoError(t, err)

	// then
	assert.Equal(t, original.Name, result.Name)
	assert.Equal(t, original.Value, result.Value)
	assert.Equal(t, original.Required, result.Required)
	assert.Equal(t, original.Export, result.Export)
	assert.Equal(t, original.ParamType, result.ParamType)
	assert.Equal(t, original.Description, result.Description)
}

func TestParamJson_DefaultValues(t *testing.T) {
	// Verify default constants
	assert.Equal(t, "", DefaultParamName)
	assert.Equal(t, "", DefaultParamValue)
	assert.Equal(t, false, DefaultParamRequired)
	assert.Equal(t, false, DefaultParamExport)
	assert.Equal(t, "PARAM", DefaultParamType)
	assert.Equal(t, "", DefaultParamDescription)
}
