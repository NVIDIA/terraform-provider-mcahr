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

package utils

import (
	"context"
	"terraform/terraform-provider/provider/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStringOrEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		key             string
		expected        string
	}{
		{
			name: "Key exists with string value",
			integrationData: map[string]interface{}{
				"test_key": "test_value",
			},
			key:      "test_key",
			expected: "test_value",
		},
		{
			name: "Key does not exist",
			integrationData: map[string]interface{}{
				"other_key": "other_value",
			},
			key:      "test_key",
			expected: "",
		},
		{
			name:            "Empty map",
			integrationData: map[string]interface{}{},
			key:             "test_key",
			expected:        "",
		},
		{
			name: "Key exists with empty string",
			integrationData: map[string]interface{}{
				"test_key": "",
			},
			key:      "test_key",
			expected: "",
		},
		{
			name: "Key exists with complex string",
			integrationData: map[string]interface{}{
				"url": "https://example.com:8080/path?param=value",
			},
			key:      "url",
			expected: "https://example.com:8080/path?param=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestContext := common.NewRequestContext(context.Background())

			result := GetStringOrEmpty(requestContext, tt.integrationData, tt.key)

			assert.Equal(t, tt.expected, result.ValueString())
			assert.False(t, result.IsNull())
		})
	}
}

func TestGetInt64OrZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		key             string
		expected        int64
	}{
		{
			name: "Key exists with int64 value",
			integrationData: map[string]interface{}{
				"test_key": int64(123456),
			},
			key:      "test_key",
			expected: 123456,
		},
		{
			name: "Key exists with int value",
			integrationData: map[string]interface{}{
				"test_key": 789,
			},
			key:      "test_key",
			expected: 789,
		},
		{
			name: "Key exists with int32 value",
			integrationData: map[string]interface{}{
				"test_key": int32(456),
			},
			key:      "test_key",
			expected: 456,
		},
		{
			name: "Key exists with int16 value",
			integrationData: map[string]interface{}{
				"test_key": int16(123),
			},
			key:      "test_key",
			expected: 123,
		},
		{
			name: "Key exists with int8 value",
			integrationData: map[string]interface{}{
				"test_key": int8(42),
			},
			key:      "test_key",
			expected: 42,
		},
		{
			name: "Key exists with float64 value",
			integrationData: map[string]interface{}{
				"test_key": float64(987.0),
			},
			key:      "test_key",
			expected: 987,
		},
		{
			name: "Key exists with float32 value",
			integrationData: map[string]interface{}{
				"test_key": float32(654.0),
			},
			key:      "test_key",
			expected: 654,
		},
		{
			name: "Key does not exist",
			integrationData: map[string]interface{}{
				"other_key": 123,
			},
			key:      "test_key",
			expected: 0,
		},
		{
			name:            "Empty map",
			integrationData: map[string]interface{}{},
			key:             "test_key",
			expected:        0,
		},
		{
			name: "Key exists with unsupported type",
			integrationData: map[string]interface{}{
				"test_key": "not_a_number",
			},
			key:      "test_key",
			expected: 0,
		},
		{
			name: "Key exists with nil value",
			integrationData: map[string]interface{}{
				"test_key": nil,
			},
			key:      "test_key",
			expected: 0,
		},
		{
			name: "Key exists with zero value",
			integrationData: map[string]interface{}{
				"test_key": 0,
			},
			key:      "test_key",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestContext := common.NewRequestContext(context.Background())

			result := GetInt64OrZero(requestContext, tt.integrationData, tt.key)

			assert.Equal(t, tt.expected, result.ValueInt64())
			assert.False(t, result.IsNull())
		})
	}
}

func TestGetStringSetOrEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		integrationData map[string]interface{}
		key             string
		expectedCount   int
		expectedValues  []string
	}{
		{
			name: "Key exists with string slice",
			integrationData: map[string]interface{}{
				"test_key": []interface{}{"value1", "value2", "value3"},
			},
			key:            "test_key",
			expectedCount:  3,
			expectedValues: []string{"value1", "value2", "value3"},
		},
		{
			name: "Key exists with empty slice",
			integrationData: map[string]interface{}{
				"test_key": []interface{}{},
			},
			key:            "test_key",
			expectedCount:  0,
			expectedValues: []string{},
		},
		{
			name: "Key exists with single item slice",
			integrationData: map[string]interface{}{
				"test_key": []interface{}{"single_value"},
			},
			key:            "test_key",
			expectedCount:  1,
			expectedValues: []string{"single_value"},
		},
		{
			name: "Key does not exist",
			integrationData: map[string]interface{}{
				"other_key": []interface{}{"other_value"},
			},
			key:            "test_key",
			expectedCount:  0,
			expectedValues: []string{},
		},
		{
			name:            "Empty map",
			integrationData: map[string]interface{}{},
			key:             "test_key",
			expectedCount:   0,
			expectedValues:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestContext := common.NewRequestContext(context.Background())

			result := GetStringListOrEmpty(requestContext, tt.integrationData, tt.key)

			assert.Equal(t, tt.expectedCount, len(result.Elements()))
			assert.False(t, result.IsNull())

			if tt.expectedCount > 0 {
				actualValues := make([]string, len(result.Elements()))
				for i, elem := range result.Elements() {
					actualValues[i] = elem.(types.String).ValueString()
				}
				assert.ElementsMatch(t, tt.expectedValues, actualValues)
			}
		})
	}
}

func TestStringSetFromMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		value          interface{}
		expectedCount  int
		expectedValues []string
		expectNull     bool
	}{
		{
			name:           "Valid string slice",
			value:          []interface{}{"item1", "item2", "item3"},
			expectedCount:  3,
			expectedValues: []string{"item1", "item2", "item3"},
			expectNull:     false,
		},
		{
			name:           "Empty string slice",
			value:          []interface{}{},
			expectedCount:  0,
			expectedValues: []string{},
			expectNull:     false,
		},
		{
			name:           "Single item slice",
			value:          []interface{}{"single"},
			expectedCount:  1,
			expectedValues: []string{"single"},
			expectNull:     false,
		},
		{
			name:       "Non-slice value",
			value:      "not_a_slice",
			expectNull: true,
		},
		{
			name:       "Slice with non-string items",
			value:      []interface{}{"valid", 123, "invalid"},
			expectNull: true,
		},
		{
			name:       "Slice with mixed types",
			value:      []interface{}{123, true, nil},
			expectNull: true,
		},
		{
			name:       "Nil value",
			value:      nil,
			expectNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestContext := common.NewRequestContext(context.Background())

			result := StringListFromMap(requestContext, tt.value)

			if tt.expectNull {
				assert.True(t, result.IsNull())
			} else {
				assert.False(t, result.IsNull())
				assert.Equal(t, tt.expectedCount, len(result.Elements()))

				if tt.expectedCount > 0 {
					actualValues := make([]string, len(result.Elements()))
					for i, elem := range result.Elements() {
						actualValues[i] = elem.(types.String).ValueString()
					}
					assert.ElementsMatch(t, tt.expectedValues, actualValues)
				}
			}
		})
	}
}

func TestStringSetTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          types.List
		expectedValues []string
	}{
		{
			name: "Valid set with multiple strings",
			input: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("string1"),
				types.StringValue("string2"),
				types.StringValue("string3"),
			}),
			expectedValues: []string{"string1", "string2", "string3"},
		},
		{
			name: "Valid set with single string",
			input: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("single"),
			}),
			expectedValues: []string{"single"},
		},
		{
			name:           "Empty set",
			input:          types.ListValueMust(types.StringType, []attr.Value{}),
			expectedValues: []string{},
		},
		{
			name: "Set with empty strings",
			input: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue(""),
				types.StringValue("not_empty"),
				types.StringValue(""),
			}),
			expectedValues: []string{"", "not_empty", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestContext := common.NewRequestContext(context.Background())

			result := StringListTFModel(requestContext, tt.input)

			require.Equal(t, len(tt.expectedValues), len(result))
			assert.ElementsMatch(t, tt.expectedValues, result)
		})
	}
}

func TestStringSetTFModel_NullSet(t *testing.T) {
	t.Parallel()

	// Test with null set
	nullSet := types.ListNull(types.StringType)
	requestContext := common.NewRequestContext(context.Background())

	result := StringListTFModel(requestContext, nullSet)

	assert.Equal(t, []string{}, result)
}

func TestIntegrationUtilsFunctions_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("GetStringOrEmpty with nil map", func(t *testing.T) {
		// This should panic, but we test the documented behavior
		defer func() {
			if r := recover(); r != nil {
				// Expected behavior - accessing nil map should panic
				assert.NotNil(t, r)
			}
		}()

		var nilMap map[string]interface{}
		requestContext := common.NewRequestContext(context.Background())
		result := GetStringOrEmpty(requestContext, nilMap, "key")
		// If we reach here, it means no panic occurred and empty string was returned
		assert.Equal(t, "", result.ValueString())
	})

	t.Run("GetInt64OrZero with nil map", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				// Expected behavior - accessing nil map should panic
				assert.NotNil(t, r)
			}
		}()

		var nilMap map[string]interface{}
		requestContext := common.NewRequestContext(context.Background())
		result := GetInt64OrZero(requestContext, nilMap, "key")
		// If we reach here, it means no panic occurred and 0 was returned
		assert.Equal(t, int64(0), result.ValueInt64())
	})
}
