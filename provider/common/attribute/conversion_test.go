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

package attribute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToInt64(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		expectError bool
		expected    int64
	}{
		{
			name:        "From float32",
			value:       float32(123.0),
			expectError: false,
			expected:    123,
		},
		{
			name:        "From float64",
			value:       float64(456.0),
			expectError: false,
			expected:    456,
		},
		{
			name:        "From int",
			value:       int(789),
			expectError: false,
			expected:    789,
		},
		{
			name:        "From int32",
			value:       int32(1000),
			expectError: false,
			expected:    1000,
		},
		{
			name:        "From int64",
			value:       int64(2000),
			expectError: false,
			expected:    2000,
		},
		{
			name:        "From string should error",
			value:       "123",
			expectError: true,
		},
		{
			name:        "From bool should error",
			value:       true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := ConvertToInt64(tt.value)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result.ValueInt64())
			}
		})
	}
}

func TestConvertToStringList(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		expectError bool
		validate    func(t *testing.T, result types.List)
	}{
		{
			name:        "String array",
			value:       []interface{}{"item1", "item2", "item3"},
			expectError: false,
			validate: func(t *testing.T, result types.List) {
				assert.Equal(t, 3, len(result.Elements()))
			},
		},
		{
			name:        "Empty array returns empty list not null",
			value:       []interface{}{},
			expectError: false,
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull(), "empty array should return empty list, not null")
				assert.Equal(t, 0, len(result.Elements()))
			},
		},
		{
			name:        "Mixed types array should error",
			value:       []interface{}{"string1", 123, "string2", true},
			expectError: true,
		},
		{
			name:        "Non-array value should error",
			value:       "not_an_array",
			expectError: true,
		},
		{
			name:        "Nil value should error",
			value:       nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := ConvertToStringList(tt.value)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}
