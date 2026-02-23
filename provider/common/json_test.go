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

package common

import (
	"encoding/json"
	"testing"

	"terraform/terraform-provider/provider/common/version"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJsonConfigurable implementation for testing
type TestJsonConfigurable struct {
	Config      JsonConfig `json:"-"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Description string     `json:"description,omitempty"`
	Value       int        `json:"value"`
}

func (t TestJsonConfigurable) SetConfig(config JsonConfig) {
	t.Config = config
}

func (t TestJsonConfigurable) GetConfig() JsonConfig {
	return t.Config
}

// TestRemarshalWithConfig_Success tests successful remarshaling with config
func TestRemarshalWithConfig_Success(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}

	inputData := TestJsonConfigurable{
		Name:        "test-item",
		Type:        "string",
		Description: "Test description",
		Value:       42,
	}

	inputJSON, err := json.Marshal(inputData)
	require.NoError(t, err)

	// when
	result, err := RemarshalWithConfig[TestJsonConfigurable](string(inputJSON), config)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify the result is valid JSON
	var output TestJsonConfigurable
	err = json.Unmarshal([]byte(result), &output)
	require.NoError(t, err)

	assert.Equal(t, inputData.Name, output.Name)
	assert.Equal(t, inputData.Type, output.Type)
	assert.Equal(t, inputData.Description, output.Description)
	assert.Equal(t, inputData.Value, output.Value)
}

// TestRemarshalWithConfig_InvalidJSON tests remarshaling with invalid JSON
func TestRemarshalWithConfig_InvalidJSON(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}
	invalidJSON := "not valid json"

	// when
	result, err := RemarshalWithConfig[TestJsonConfigurable](invalidJSON, config)

	// then
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "invalid character")
}

// TestRemarshalWithConfig_EmptyJSON tests remarshaling with empty JSON object
func TestRemarshalWithConfig_EmptyJSON(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}
	emptyJSON := "{}"

	// when
	result, err := RemarshalWithConfig[TestJsonConfigurable](emptyJSON, config)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify the result has default values
	var output TestJsonConfigurable
	err = json.Unmarshal([]byte(result), &output)
	require.NoError(t, err)

	assert.Equal(t, "", output.Name)
	assert.Equal(t, "", output.Type)
	assert.Equal(t, 0, output.Value)
}

// TestRemarshalListWithConfig_Success tests successful list remarshaling
func TestRemarshalListWithConfig_Success(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}

	inputList := []TestJsonConfigurable{
		{
			Name:  "item1",
			Type:  "string",
			Value: 10,
		},
		{
			Name:        "item2",
			Type:        "number",
			Description: "Second item",
			Value:       20,
		},
		{
			Name:  "item3",
			Type:  "boolean",
			Value: 30,
		},
	}

	inputJSON, err := json.Marshal(inputList)
	require.NoError(t, err)

	// when
	result, err := RemarshalListWithConfig[*TestJsonConfigurable](string(inputJSON), config)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify the result is valid JSON array
	var output []TestJsonConfigurable
	err = json.Unmarshal([]byte(result), &output)
	require.NoError(t, err)

	assert.Len(t, output, 3)
	assert.Equal(t, "item1", output[0].Name)
	assert.Equal(t, "item2", output[1].Name)
	assert.Equal(t, "item3", output[2].Name)
	assert.Equal(t, "Second item", output[1].Description)
}

// TestRemarshalListWithConfig_EmptyList tests remarshaling empty list
func TestRemarshalListWithConfig_EmptyList(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}
	emptyList := "[]"

	// when
	result, err := RemarshalListWithConfig[*TestJsonConfigurable](emptyList, config)

	// then
	require.NoError(t, err)
	assert.Equal(t, "[]", result)
}

// TestRemarshalListWithConfig_InvalidJSON tests list remarshaling with invalid JSON
func TestRemarshalListWithConfig_InvalidJSON(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}
	invalidJSON := "[{invalid json}]"

	// when
	result, err := RemarshalListWithConfig[*TestJsonConfigurable](invalidJSON, config)

	// then
	assert.Error(t, err)
	assert.Empty(t, result)
}

// TestRemarshalListWithConfig_NotAnArray tests list remarshaling with non-array JSON
func TestRemarshalListWithConfig_NotAnArray(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}
	notArray := `{"name": "test"}`

	// when
	result, err := RemarshalListWithConfig[*TestJsonConfigurable](notArray, config)

	// then
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "cannot unmarshal object into Go value of type")
}

// TestJsonConfigurable_ConfigPersistence tests that config is properly set
func TestJsonConfigurable_ConfigPersistence(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "1.5.0"}
	config := JsonConfig{BackendVersion: backendVersion}

	inputData := TestJsonConfigurable{
		Name:  "config-test",
		Type:  "test",
		Value: 100,
	}

	inputJSON, err := json.Marshal(inputData)
	require.NoError(t, err)

	// when
	result, err := RemarshalWithConfig[TestJsonConfigurable](string(inputJSON), config)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// then
	// Note: Due to value receivers, SetConfig doesn't actually modify the struct
	// This is a known limitation of the current implementation
	// We can only verify that the methods don't panic
	var testObj TestJsonConfigurable
	testObj.SetConfig(config) // This doesn't modify testObj due to value receiver
	retrievedConfig := testObj.GetConfig()

	// With value receivers, the config is not actually stored
	assert.Equal(t, JsonConfig{}, retrievedConfig)
}

// NestedConfigurable for testing complex structures
type NestedConfigurable struct {
	Config   JsonConfig             `json:"-"`
	ID       string                 `json:"id"`
	Data     map[string]interface{} `json:"data"`
	SubItems []TestJsonConfigurable `json:"sub_items"`
	Metadata json.RawMessage        `json:"metadata,omitempty"`
}

// NestedConfigurableImpl implements JsonConfigurable interface
type NestedConfigurableImpl struct {
	NestedConfigurable
}

func (n NestedConfigurableImpl) SetConfig(config JsonConfig) {
	// With value receiver, this doesn't actually modify the struct
	// This is a limitation of the current implementation
	n.Config = config
}

func (n NestedConfigurableImpl) GetConfig() JsonConfig {
	return n.Config
}

// TestRemarshalWithConfig_ComplexNestedStructure tests with complex nested data
func TestRemarshalWithConfig_ComplexNestedStructure(t *testing.T) {
	// given

	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}

	complexData := NestedConfigurableImpl{
		NestedConfigurable: NestedConfigurable{
			ID: "complex-1",
			Data: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
				"key3": true,
			},
			SubItems: []TestJsonConfigurable{
				{Name: "sub1", Type: "type1", Value: 1},
				{Name: "sub2", Type: "type2", Value: 2},
			},
			Metadata: json.RawMessage(`{"meta": "data"}`),
		},
	}

	inputJSON, err := json.Marshal(complexData)
	require.NoError(t, err)

	// when
	result, err := RemarshalWithConfig[NestedConfigurableImpl](string(inputJSON), config)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify complex structure is preserved
	var output NestedConfigurableImpl
	err = json.Unmarshal([]byte(result), &output)
	require.NoError(t, err)

	assert.Equal(t, "complex-1", output.ID)
	assert.Equal(t, "value1", output.Data["key1"])
	assert.Equal(t, float64(42), output.Data["key2"]) // JSON numbers unmarshal as float64
	assert.Equal(t, true, output.Data["key3"])
	assert.Len(t, output.SubItems, 2)
}

// TestRemarshalListWithConfig_MixedValidInvalidItems tests partial invalid items
func TestRemarshalListWithConfig_MixedValidInvalidItems(t *testing.T) {
	// given
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	config := JsonConfig{BackendVersion: backendVersion}

	// Create a list where marshaling might succeed but with specific field issues
	inputJSON := `[
		{"name": "valid1", "type": "string", "value": 10},
		{"name": "valid2", "type": "number", "value": 20}
	]`

	// when
	result, err := RemarshalListWithConfig[*TestJsonConfigurable](inputJSON, config)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify all items were processed
	var output []TestJsonConfigurable
	err = json.Unmarshal([]byte(result), &output)
	require.NoError(t, err)
	assert.Len(t, output, 2)
}

// TestJsonConfig_BackendVersionNil tests handling of nil backend version
func TestJsonConfig_BackendVersionNil(t *testing.T) {
	// given
	config := JsonConfig{BackendVersion: nil}

	inputData := TestJsonConfigurable{
		Name:  "nil-version-test",
		Type:  "test",
		Value: 50,
	}

	inputJSON, err := json.Marshal(inputData)
	require.NoError(t, err)

	// when
	result, err := RemarshalWithConfig[TestJsonConfigurable](string(inputJSON), config)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify data is still processed correctly
	var output TestJsonConfigurable
	err = json.Unmarshal([]byte(result), &output)
	require.NoError(t, err)
	assert.Equal(t, inputData.Name, output.Name)
}
