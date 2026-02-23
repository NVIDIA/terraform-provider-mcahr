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
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCellJsonAPI_ToInternalModel(t *testing.T) {
	tests := []struct {
		name     string
		input    *CellJsonAPI
		expected *CellJson
	}{
		{
			name: "OP_LANG type cell",
			input: &CellJsonAPI{
				Type:        OP_LANG_TYPE,
				Content:     "print('hello')",
				Name:        "test_cell",
				Enabled:     true,
				SecretAware: true,
				Description: "Test cell",
			},
			expected: &CellJson{
				Op:          common.NewOptional("print('hello')"),
				Name:        "test_cell",
				Enabled:     true,
				SecretAware: true,
				Description: "Test cell",
			},
		},
		{
			name: "MARKDOWN type cell",
			input: &CellJsonAPI{
				Type:        MARKDOWN_TYPE,
				Content:     "# Hello World",
				Name:        "md_cell",
				Enabled:     false,
				SecretAware: false,
				Description: "Markdown cell",
			},
			expected: &CellJson{
				Md:          common.NewOptional("# Hello World"),
				Name:        "md_cell",
				Enabled:     false,
				SecretAware: false,
				Description: "Markdown cell",
			},
		},
		{
			name: "Cell with CellType field",
			input: &CellJsonAPI{
				CellType:    OP_LANG_TYPE,
				Content:     "code content",
				Name:        "cell_type_test",
				Enabled:     true,
				SecretAware: false,
				Description: "",
			},
			expected: &CellJson{
				Op:          common.NewOptional("code content"),
				Name:        "cell_type_test",
				Enabled:     true,
				SecretAware: false,
				Description: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := tt.input.ToInternalModel()

			// then
			assert.Equal(t, tt.expected.Op, result.Op)
			assert.Equal(t, tt.expected.Md, result.Md)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.SecretAware, result.SecretAware)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}

func TestCellJsonAPI_SetFromMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected *CellJsonAPI
	}{
		{
			name: "Complete cell data with type field",
			input: map[string]interface{}{
				"name":         "test_cell",
				"content":      "print('test')",
				"enabled":      true,
				"type":         OP_LANG_TYPE,
				"secret_aware": true,
				"description":  "Test description",
			},
			expected: &CellJsonAPI{
				Name:        "test_cell",
				Content:     "print('test')",
				Enabled:     true,
				CellType:    OP_LANG_TYPE,
				SecretAware: true,
				Description: "Test description",
			},
		},
		{
			name: "Cell data with cell_type field",
			input: map[string]interface{}{
				"name":         "md_cell",
				"content":      "# Markdown",
				"enabled":      false,
				"cell_type":    MARKDOWN_TYPE,
				"secret_aware": false,
			},
			expected: &CellJsonAPI{
				Name:        "md_cell",
				Content:     "# Markdown",
				Enabled:     false,
				CellType:    MARKDOWN_TYPE,
				SecretAware: false,
				Description: "",
			},
		},
		{
			name:  "Empty map",
			input: map[string]interface{}{},
			expected: &CellJsonAPI{
				Name:        "",
				Content:     "",
				Enabled:     false,
				CellType:    "",
				SecretAware: false,
				Description: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			cell := &CellJsonAPI{}

			// when
			cell.SetFromMap(tt.input)

			// then
			assert.Equal(t, tt.expected.Name, cell.Name)
			assert.Equal(t, tt.expected.Content, cell.Content)
			assert.Equal(t, tt.expected.Enabled, cell.Enabled)
			assert.Equal(t, tt.expected.CellType, cell.CellType)
			assert.Equal(t, tt.expected.SecretAware, cell.SecretAware)
			assert.Equal(t, tt.expected.Description, cell.Description)
		})
	}
}

func TestCellJson_ToAPIModel(t *testing.T) {
	tests := []struct {
		name     string
		input    *CellJson
		expected *CellJsonAPI
	}{
		{
			name: "Op cell",
			input: &CellJson{
				Op:          common.NewOptional("print('hello')"),
				Name:        "op_cell",
				Enabled:     true,
				SecretAware: true,
				Description: "Op cell",
			},
			expected: &CellJsonAPI{
				Content:     "print('hello')",
				Type:        OP_LANG_TYPE,
				Name:        "op_cell",
				Enabled:     true,
				SecretAware: true,
				Description: "Op cell",
			},
		},
		{
			name: "Md cell",
			input: &CellJson{
				Md:          common.NewOptional("# Markdown"),
				Name:        "md_cell",
				Enabled:     false,
				SecretAware: false,
				Description: "Md cell",
			},
			expected: &CellJsonAPI{
				Content:     "# Markdown",
				Type:        MARKDOWN_TYPE,
				Name:        "md_cell",
				Enabled:     false,
				SecretAware: false,
				Description: "Md cell",
			},
		},
		{
			name: "Empty cell defaults to OP_LANG",
			input: &CellJson{
				Name:        "empty_cell",
				Enabled:     true,
				SecretAware: false,
				Description: "",
			},
			expected: &CellJsonAPI{
				Content:     DefaultCellContent,
				Type:        DefaultCellType,
				Name:        "empty_cell",
				Enabled:     true,
				SecretAware: false,
				Description: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := tt.input.ToAPIModel()

			// then
			assert.Equal(t, tt.expected.Content, result.Content)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.SecretAware, result.SecretAware)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}

func TestCellJson_MarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         *CellJson
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid Op cell",
			input: &CellJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
				},
				Op:          common.NewOptional("print('test')"),
				Name:        "test_cell",
				Enabled:     true,
				SecretAware: true,
				Description: "Test",
			},
			expectError: false,
		},
		{
			name: "Valid Md cell",
			input: &CellJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
				},
				Md:          common.NewOptional("# Test"),
				Name:        "md_cell",
				Enabled:     false,
				SecretAware: false,
				Description: "Markdown",
			},
			expectError: false,
		},
		{
			name: "Both Op and Md set",
			input: &CellJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
				},
				Op:   common.NewOptional("print('test')"),
				Md:   common.NewOptional("# Test"),
				Name: "invalid_cell",
			},
			expectError:   true,
			errorContains: "cannot have both op and md",
		},
		{
			name: "Neither Op nor Md set",
			input: &CellJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
				},
				Name: "invalid_cell",
			},
			expectError:   true,
			errorContains: "must have either op or md",
		},
		{
			name: "Version filtering - description excluded",
			input: &CellJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "release-28.0.0", Major: 28, Minor: 0, Patch: 0},
				},
				Op:          common.NewOptional("print('test')"),
				Name:        "test_cell",
				Description: "Should be excluded",
			},
			expectError: false,
		},
		{
			name: "Cell with default enabled value should marshal to true",
			input: &CellJson{
				Config: common.JsonConfig{
					BackendVersion: &version.BackendVersion{Version: "2.0.0"},
				},
				Op:          common.NewOptional("print('hello')"),
				Name:        "test_cell",
				Enabled:     true,
				SecretAware: false,
				Description: "",
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
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)

				var unmarshaled map[string]interface{}
				err = json.Unmarshal(result, &unmarshaled)
				require.NoError(t, err)

				// Check that unset fields are omitted
				if !tt.input.Op.IsSet {
					assert.NotContains(t, unmarshaled, "op")
				}
				if !tt.input.Md.IsSet {
					assert.NotContains(t, unmarshaled, "md")
				}

				// Check version filtering
				if tt.input.Config.BackendVersion != nil && tt.input.Config.BackendVersion.Major < 29 {
					assert.NotContains(t, unmarshaled, "description")
				}
			}
		})
	}
}

func TestCellJson_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		config        common.JsonConfig
		expected      *CellJson
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid Op cell JSON",
			input: `{
				"op": "print('hello')",
				"name": "test_cell",
				"enabled": true,
				"secret_aware": true,
				"description": "Test cell"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
			},
			expected: &CellJson{
				Op:          common.NewOptional("print('hello')"),
				Name:        "test_cell",
				Enabled:     true,
				SecretAware: true,
				Description: "Test cell",
			},
			expectError: false,
		},
		{
			name: "Valid Md cell JSON",
			input: `{
				"md": "# Markdown",
				"name": "md_cell",
				"enabled": false,
				"secret_aware": false
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
			},
			expected: &CellJson{
				Md:          common.NewOptional("# Markdown"),
				Name:        "md_cell",
				Enabled:     false,
				SecretAware: false,
				Description: DefaultCellDescription,
			},
			expectError: false,
		},
		{
			name: "Op cell JSON without enabled field - should default to true",
			input: `{
				"op": "print('hello')",
				"name": "test_cell"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "2.0.0"},
			},
			expected: &CellJson{
				Op:          common.NewOptional("print('hello')"),
				Name:        "test_cell",
				Enabled:     true,
				SecretAware: false,
				Description: "",
			},
			expectError: false,
		},
		{
			name: "Both op and md in JSON",
			input: `{
				"op": "print('hello')",
				"md": "# Markdown",
				"name": "invalid"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
			},
			expectError:   true,
			errorContains: "cannot have both op and md",
		},
		{
			name: "Neither op nor md in JSON",
			input: `{
				"name": "invalid"
			}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
			},
			expectError:   true,
			errorContains: "must have either op or md",
		},
		{
			name:  "Invalid JSON",
			input: `{invalid json}`,
			config: common.JsonConfig{
				BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			cell := &CellJson{Config: tt.config}

			// when
			err := cell.UnmarshalJSON([]byte(tt.input))

			// then
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Op, cell.Op)
				assert.Equal(t, tt.expected.Md, cell.Md)
				assert.Equal(t, tt.expected.Name, cell.Name)
				assert.Equal(t, tt.expected.Enabled, cell.Enabled)
				assert.Equal(t, tt.expected.SecretAware, cell.SecretAware)
				assert.Equal(t, tt.expected.Description, cell.Description)
			}
		})
	}
}

func TestValidateOpAndMd(t *testing.T) {
	tests := []struct {
		name          string
		op            common.Optional[string]
		md            common.Optional[string]
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid op only",
			op:          common.NewOptional("print('test')"),
			md:          common.NewOptionalUnset[string](),
			expectError: false,
		},
		{
			name:        "Valid md only",
			op:          common.NewOptionalUnset[string](),
			md:          common.NewOptional("# Markdown"),
			expectError: false,
		},
		{
			name:          "Both op and md",
			op:            common.NewOptional("print('test')"),
			md:            common.NewOptional("# Markdown"),
			expectError:   true,
			errorContains: "cannot have both op and md",
		},
		{
			name:          "Neither op nor md",
			op:            common.NewOptionalUnset[string](),
			md:            common.NewOptionalUnset[string](),
			expectError:   true,
			errorContains: "must have either op or md",
		},
		{
			name:        "Empty op content",
			op:          common.NewOptional(""),
			md:          common.NewOptionalUnset[string](),
			expectError: false,
		},
		{
			name:        "Empty md content",
			op:          common.NewOptionalUnset[string](),
			md:          common.NewOptional(""),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			err := validateOpAndMd(tt.op, tt.md)

			// then
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMapCellsToAPIModel(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid cells JSON",
			input: `[
				{"op": "print('hello')", "name": "cell1", "enabled": true},
				{"md": "# Markdown", "name": "cell2", "enabled": false}
			]`,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			input:       `{invalid json}`,
			expectError: true,
		},
		{
			name:        "Empty cells",
			input:       `[]`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := MapCellsToAPIModel(tt.input)

			// then
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, result) // Result should be base64 encoded with quotes

				// Remove surrounding quotes
				assert.True(t, strings.HasPrefix(result, `"`))
				assert.True(t, strings.HasSuffix(result, `"`))
				unquoted := result[1 : len(result)-1]

				// Verify it's valid base64
				decoded, err := base64.StdEncoding.DecodeString(unquoted)
				require.NoError(t, err)
				assert.NotEmpty(t, string(decoded))

				// Verify the decoded JSON is valid
				var cells []CellJsonAPI
				err = json.Unmarshal(decoded, &cells)
				require.NoError(t, err)
			}
		})
	}
}

func TestMapCellsToInternalModel(t *testing.T) {
	tests := []struct {
		name        string
		input       []CellJsonAPI
		expectError bool
	}{
		{
			name: "Valid API cells",
			input: []CellJsonAPI{
				{
					Type:    OP_LANG_TYPE,
					Content: "print('test')",
					Name:    "cell1",
					Enabled: true,
				},
				{
					Type:    MARKDOWN_TYPE,
					Content: "# Header",
					Name:    "cell2",
					Enabled: false,
				},
			},
			expectError: false,
		},
		{
			name:        "Empty cells",
			input:       []CellJsonAPI{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := MapCellsToInternalModel(tt.input)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify the result is valid JSON
				var cells []CellJson
				err = json.Unmarshal([]byte(result), &cells)
				assert.NoError(t, err)
				assert.Len(t, cells, len(tt.input))
			}
		})
	}
}

func TestCellJson_SetGetConfig(t *testing.T) {
	// given
	cell := &CellJson{}
	config := common.JsonConfig{
		BackendVersion: &version.BackendVersion{Version: "release-29.0.1", Major: 29, Minor: 0, Patch: 1},
	}

	// when
	cell.SetConfig(config)
	result := cell.GetConfig()

	// then
	assert.Equal(t, config.BackendVersion.Version, result.BackendVersion.Version)
}
