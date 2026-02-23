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

package translator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSetSliceFromTFModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    types.List
		expected []string
	}{
		{
			name:     "valid set with multiple elements",
			input:    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("item1"), types.StringValue("item2"), types.StringValue("item3")}),
			expected: []string{"item1", "item2", "item3"},
		},
		{
			name:     "valid set with single element",
			input:    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("single")}),
			expected: []string{"single"},
		},
		{
			name:     "empty set",
			input:    types.ListValueMust(types.StringType, []attr.Value{}),
			expected: []string{},
		},
		{
			name:     "null set",
			input:    types.ListNull(types.StringType),
			expected: nil,
		},
		{
			name:     "unknown set",
			input:    types.ListUnknown(types.StringType),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := ListSliceFromTFModel(context.Background(), input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEscapeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: `""`,
		},
		{
			name:     "simple string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "string with quotes",
			input:    `hello "world"`,
			expected: `"hello \"world\""`,
		},
		{
			name:     "string with backslashes",
			input:    `C:\Program Files\app`,
			expected: `"C:\\Program Files\\app"`,
		},
		{
			name:     "string with newlines and tabs",
			input:    "hello\nworld\ttab",
			expected: `"hello\nworld\ttab"`,
		},
		{
			name:     "string with single quotes",
			input:    "hello 'world'",
			expected: `"hello 'world'"`,
		},
		{
			name:     "string with mixed special characters",
			input:    `Line 1\nLine 2\t"quoted"\backslash`,
			expected: `"Line 1\\nLine 2\\t\"quoted\"\\backslash"`,
		},
		{
			name:     "unicode characters",
			input:    "Hello 世界 🌍",
			expected: `"Hello 世界 🌍"`,
		},
		{
			name:     "control characters",
			input:    "hello\x00\x01world",
			expected: `"hello\x00\x01world"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := EscapeString(input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestArrayToOpLang(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty array",
			input:    []string{},
			expected: "[]",
		},
		{
			name:     "nil array",
			input:    nil,
			expected: "[]",
		},
		{
			name:     "single element",
			input:    []string{"item1"},
			expected: `["item1"]`,
		},
		{
			name:     "multiple elements",
			input:    []string{"item1", "item2", "item3"},
			expected: `["item1", "item2", "item3"]`,
		},
		{
			name:     "elements with special characters",
			input:    []string{`hello "world"`, "item\nwith\nnewlines", `C:\path`},
			expected: `["hello \"world\"", "item\nwith\nnewlines", "C:\\path"]`,
		},
		{
			name:     "empty strings in array",
			input:    []string{"", "item2", ""},
			expected: `["", "item2", ""]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := ArrayToOpLang(input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseStringArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "empty array string",
			input:    "[]",
			expected: []string{},
		},
		{
			name:     "valid single element array",
			input:    `["item1"]`,
			expected: []string{"item1"},
		},
		{
			name:     "valid multiple element array",
			input:    `["item1", "item2", "item3"]`,
			expected: []string{"item1", "item2", "item3"},
		},
		{
			name:     "array with special characters",
			input:    `["hello \"world\"", "item\nwith\nnewlines"]`,
			expected: []string{`hello "world"`, "item\nwith\nnewlines"},
		},
		{
			name:     "array with empty strings",
			input:    `["", "item2", ""]`,
			expected: []string{"", "item2", ""},
		},
		{
			name:     "invalid JSON returns nil",
			input:    `[invalid json`,
			expected: nil,
		},
		{
			name:     "non-array JSON returns nil",
			input:    `{"key": "value"}`,
			expected: nil,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := ParseStringArray(input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "non-empty string",
			input:    "hello",
			expected: false,
		},
		{
			name:     "whitespace string",
			input:    " ",
			expected: false,
		},
		{
			name:     "tab string",
			input:    "\t",
			expected: false,
		},
		{
			name:     "newline string",
			input:    "\n",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input

			// when
			result := IsEmpty(input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Edge case and integration tests
func TestSetSliceFromTFModel_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("set with mixed content", func(t *testing.T) {
		// given
		// Test with strings that have special characters
		input := types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("normal"),
			types.StringValue(""),
			types.StringValue("with\nnewlines"),
			types.StringValue(`with "quotes"`),
		})
		expected := []string{"normal", "", "with\nnewlines", `with "quotes"`}

		// when
		result := ListSliceFromTFModel(context.Background(), input)

		// then
		assert.Equal(t, expected, result)
	})
}

// Integration test: round-trip conversion
func TestArrayToOpLang_ParseStringArray_RoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("simple array structure", func(t *testing.T) {
		// given
		original := []string{"item1", "item2", "item3"}

		// when
		opLangResult := ArrayToOpLang(original)

		// then
		// OpLang format should be exact array structure
		expected := `["item1", "item2", "item3"]`
		assert.Equal(t, expected, opLangResult)
	})
}

func TestBoolToInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    bool
		expected int
	}{
		{"true converts to 1", true, 1},
		{"false converts to 0", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolToInt(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntToBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"0 converts to false", 0, false},
		{"1 converts to true", 1, true},
		{"positive number converts to true", 5, true},
		{"negative number converts to true", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntToBool(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
