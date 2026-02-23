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

func TestExternalParamJson_SetGetConfig(t *testing.T) {
	// given
	param := &ExternalParamJson{}
	config := common.JsonConfig{
		BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
	}

	// when
	param.SetConfig(config)
	result := param.GetConfig()

	// then
	assert.Equal(t, config.BackendVersion.Version, result.BackendVersion.Version)
}

func TestExternalParamJson_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       *ExternalParamJson
		expectError bool
	}{
		{
			name: "Complete external param with description (version 28.4.0+)",
			input: &ExternalParamJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
				},
				Name:        "test_ext_param",
				Value:       "test_value",
				Source:      "test_source",
				JsonPath:    "$.data.value",
				Export:      true,
				ParamType:   "EXTERNAL",
				Description: "Test external parameter",
			},
			expectError: false,
		},
		{
			name: "External param without description (older version)",
			input: &ExternalParamJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-28.3.0", Major: 28, Minor: 3, Patch: 0},
				},
				Name:        "test_ext_param",
				Value:       "test_value",
				Source:      "test_source",
				JsonPath:    "$.data",
				Export:      false,
				ParamType:   "EXTERNAL",
				Description: "Should be excluded",
			},
			expectError: false,
		},
		{
			name: "External param with empty values",
			input: &ExternalParamJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
				},
				Name:        "",
				Value:       "",
				Source:      "",
				JsonPath:    "",
				Export:      false,
				ParamType:   "EXTERNAL",
				Description: "",
			},
			expectError: false,
		},
		{
			name: "External param with nil backend version",
			input: &ExternalParamJson{
				Config: common.JsonConfig{
					BackendVersion: nil,
				},
				Name:      "test_param",
				Value:     "test_value",
				Source:    "test_source",
				ParamType: "EXTERNAL",
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
				assert.Contains(t, unmarshaled, "source")
				assert.Contains(t, unmarshaled, "json_path")
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

func TestExternalParamJson_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		config      common.JsonConfig
		expected    *ExternalParamJson
		expectError bool
	}{
		{
			name: "Complete external param JSON",
			input: `{
				"name": "test_ext_param",
				"value": "test_value",
				"source": "test_source",
				"json_path": "$.data.value",
				"export": true,
				"param_type": "EXTERNAL",
				"description": "Test external parameter"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			expected: &ExternalParamJson{
				Name:        "test_ext_param",
				Value:       "test_value",
				Source:      "test_source",
				JsonPath:    "$.data.value",
				Export:      true,
				ParamType:   "EXTERNAL",
				Description: "Test external parameter",
			},
			expectError: false,
		},
		{
			name: "Partial external param JSON with defaults",
			input: `{
				"name": "minimal_param",
				"source": "api_endpoint"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			expected: &ExternalParamJson{
				Name:        "minimal_param",
				Value:       DefaultExternalParamValue,
				Source:      "api_endpoint",
				JsonPath:    DefaultExternalParamJsonPath,
				Export:      DefaultExternalParamExport,
				ParamType:   DefaultExternalParamType,
				Description: DefaultExternalParamDescription,
			},
			expectError: false,
		},
		{
			name:  "Empty JSON object",
			input: `{}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			expected: &ExternalParamJson{
				Name:        DefaultExternalParamName,
				Value:       DefaultExternalParamValue,
				Source:      DefaultExternalParamSource,
				JsonPath:    DefaultExternalParamJsonPath,
				Export:      DefaultExternalParamExport,
				ParamType:   DefaultExternalParamType,
				Description: DefaultExternalParamDescription,
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
			name: "External param with version filtering",
			input: `{
				"name": "version_test",
				"source": "test_source",
				"param_type": "EXTERNAL",
				"description": "Should be filtered"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.3.0", Major: 28, Minor: 3, Patch: 0},
			},
			expected: &ExternalParamJson{
				Name:        "version_test",
				Value:       DefaultExternalParamValue,
				Source:      "test_source",
				JsonPath:    DefaultExternalParamJsonPath,
				Export:      DefaultExternalParamExport,
				ParamType:   "EXTERNAL",
				Description: DefaultExternalParamDescription, // Description should not be set due to version
			},
			expectError: false,
		},
		{
			name: "External param with complex JSON path",
			input: `{
				"name": "complex_path",
				"source": "api_response",
				"json_path": "$.results[?(@.status=='active')].data",
				"param_type": "EXTERNAL"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			expected: &ExternalParamJson{
				Name:        "complex_path",
				Value:       DefaultExternalParamValue,
				Source:      "api_response",
				JsonPath:    "$.results[?(@.status=='active')].data",
				Export:      DefaultExternalParamExport,
				ParamType:   "EXTERNAL",
				Description: DefaultExternalParamDescription,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			param := &ExternalParamJson{Config: tt.config}

			// when
			err := param.UnmarshalJSON([]byte(tt.input))

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Name, param.Name)
				assert.Equal(t, tt.expected.Value, param.Value)
				assert.Equal(t, tt.expected.Source, param.Source)
				assert.Equal(t, tt.expected.JsonPath, param.JsonPath)
				assert.Equal(t, tt.expected.Export, param.Export)
				assert.Equal(t, tt.expected.ParamType, param.ParamType)
				assert.Equal(t, tt.expected.Description, param.Description)
			}
		})
	}
}

func TestExternalParamJson_UnmarshalJSON_FieldPriority(t *testing.T) {
	// given
	input := `{
		"name": "override_test",
		"value": "override_value",
		"source": "override_source",
		"json_path": "$.override.path",
		"export": true,
		"param_type": "CUSTOM_EXTERNAL",
		"description": "Override description"
	}`

	config := common.JsonConfig{
		BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
	}

	param := &ExternalParamJson{Config: config}

	// when
	err := param.UnmarshalJSON([]byte(input))

	// then
	require.NoError(t, err)
	assert.Equal(t, "override_test", param.Name)
	assert.Equal(t, "override_value", param.Value)
	assert.Equal(t, "override_source", param.Source)
	assert.Equal(t, "$.override.path", param.JsonPath)
	assert.Equal(t, true, param.Export)
	assert.Equal(t, "CUSTOM_EXTERNAL", param.ParamType)
	assert.Equal(t, "Override description", param.Description)
}

func TestExternalParamJson_MarshalUnmarshal_RoundTrip(t *testing.T) {
	// given
	original := &ExternalParamJson{
		Config: common.JsonConfig{
			BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
		},
		Name:        "roundtrip_param",
		Value:       "roundtrip_value",
		Source:      "roundtrip_source",
		JsonPath:    "$.roundtrip.data",
		Export:      true,
		ParamType:   "EXTERNAL",
		Description: "Roundtrip test",
	}

	// when
	marshaled, err := original.MarshalJSON()
	require.NoError(t, err)

	result := &ExternalParamJson{Config: original.Config}
	err = result.UnmarshalJSON(marshaled)
	require.NoError(t, err)

	// then
	assert.Equal(t, original.Name, result.Name)
	assert.Equal(t, original.Value, result.Value)
	assert.Equal(t, original.Source, result.Source)
	assert.Equal(t, original.JsonPath, result.JsonPath)
	assert.Equal(t, original.Export, result.Export)
	assert.Equal(t, original.ParamType, result.ParamType)
	assert.Equal(t, original.Description, result.Description)
}

func TestExternalParamJson_DefaultValues(t *testing.T) {
	// Verify default constants
	assert.Equal(t, "", DefaultExternalParamName)
	assert.Equal(t, "", DefaultExternalParamValue)
	assert.Equal(t, "", DefaultExternalParamSource)
	assert.Equal(t, "", DefaultExternalParamJsonPath)
	assert.Equal(t, false, DefaultExternalParamExport)
	assert.Equal(t, "EXTERNAL", DefaultExternalParamType)
	assert.Equal(t, "", DefaultExternalParamDescription)
}

func TestExternalParamJson_SpecialCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		config   common.JsonConfig
		validate func(t *testing.T, param *ExternalParamJson)
	}{
		{
			name: "Empty source and json_path",
			input: `{
				"name": "no_source",
				"value": "static_value",
				"param_type": "EXTERNAL"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			validate: func(t *testing.T, param *ExternalParamJson) {
				assert.Equal(t, "no_source", param.Name)
				assert.Equal(t, "static_value", param.Value)
				assert.Equal(t, "", param.Source)
				assert.Equal(t, "", param.JsonPath)
			},
		},
		{
			name: "Source without json_path",
			input: `{
				"name": "source_only",
				"source": "data_source",
				"param_type": "EXTERNAL"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-28.4.0", Major: 28, Minor: 4, Patch: 0},
			},
			validate: func(t *testing.T, param *ExternalParamJson) {
				assert.Equal(t, "source_only", param.Name)
				assert.Equal(t, "data_source", param.Source)
				assert.Equal(t, "", param.JsonPath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			param := &ExternalParamJson{Config: tt.config}

			// when
			err := param.UnmarshalJSON([]byte(tt.input))

			// then
			require.NoError(t, err)
			tt.validate(t, param)
		})
	}
}
