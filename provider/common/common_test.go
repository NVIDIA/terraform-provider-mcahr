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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestEncodeBase64_Success tests successful base64 encoding
func TestEncodeBase64_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string",
			input:    "hello world",
			expected: "\"aGVsbG8gd29ybGQ=\"",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "\"\"",
		},
		{
			name:     "Special characters",
			input:    "Hello, 世界! @#$%",
			expected: "\"SGVsbG8sIOS4lueVjCEgQCMkJQ==\"",
		},
		{
			name:     "Multi-line string",
			input:    "line1\nline2\nline3",
			expected: "\"bGluZTEKbGluZTIKbGluZTM=\"",
		},
		{
			name:     "JSON string",
			input:    `{"key": "value", "number": 123}`,
			expected: "\"eyJrZXkiOiAidmFsdWUiLCAibnVtYmVyIjogMTIzfQ==\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := EncodeBase64(tt.input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestWrapInQuotes_Success tests successful quote wrapping
func TestWrapInQuotes_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string",
			input:    "hello",
			expected: "\"hello\"",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "\"\"",
		},
		{
			name:     "String with quotes",
			input:    "hello \"world\"",
			expected: "\"hello \"world\"\"",
		},
		{
			name:     "String with special characters",
			input:    "test\nline\ttab",
			expected: "\"test\nline\ttab\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := WrapInQuotes(tt.input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsAttrKnown_AllCases tests attribute known status checking
func TestIsAttrKnown_AllCases(t *testing.T) {
	tests := []struct {
		name     string
		value    attr.Value
		expected bool
	}{
		{
			name:     "Known string value",
			value:    types.StringValue("test"),
			expected: true,
		},
		{
			name:     "Known bool value",
			value:    types.BoolValue(true),
			expected: true,
		},
		{
			name:     "Known int64 value",
			value:    types.Int64Value(42),
			expected: true,
		},
		{
			name:     "Null string value",
			value:    types.StringNull(),
			expected: false,
		},
		{
			name:     "Unknown string value",
			value:    types.StringUnknown(),
			expected: false,
		},
		{
			name:     "Null bool value",
			value:    types.BoolNull(),
			expected: false,
		},
		{
			name:     "Unknown bool value",
			value:    types.BoolUnknown(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := IsAttrKnown(tt.value)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSnakeToCamelCase_Success tests snake_case to camelCase conversion
func TestSnakeToCamelCase_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple snake case",
			input:    "hello_world",
			expected: "helloWorld",
		},
		{
			name:     "Multiple underscores",
			input:    "this_is_a_test",
			expected: "thisIsATest",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "No underscores",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "Leading underscore",
			input:    "_private_field",
			expected: "PrivateField",
		},
		{
			name:     "Trailing underscore",
			input:    "field_name_",
			expected: "fieldName",
		},
		{
			name:     "Multiple consecutive underscores",
			input:    "field__name",
			expected: "fieldName",
		},
		{
			name:     "All uppercase",
			input:    "FIELD_NAME",
			expected: "FIELDNAME",
		},
		{
			name:     "Mixed case input",
			input:    "Field_Name",
			expected: "FieldName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := SnakeToCamelCase(tt.input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCamelToSnakeCase_Success tests camelCase to snake_case conversion
func TestCamelToSnakeCase_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple camel case",
			input:    "helloWorld",
			expected: "hello_world",
		},
		{
			name:     "Multiple words",
			input:    "thisIsATest",
			expected: "this_is_a_test",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "No uppercase",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "Starting with uppercase",
			input:    "HelloWorld",
			expected: "hello_world",
		},
		{
			name:     "Consecutive uppercase",
			input:    "XMLParser",
			expected: "x_m_l_parser",
		},
		{
			name:     "Single letter",
			input:    "A",
			expected: "a",
		},
		{
			name:     "Numbers in string",
			input:    "version2API",
			expected: "version2_a_p_i",
		},
		{
			name:     "Already snake case",
			input:    "already_snake_case",
			expected: "already_snake_case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := CamelToSnakeCase(tt.input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSnakeToCamelCase_EdgeCases tests edge cases for snake to camel conversion
func TestSnakeToCamelCase_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Underscore only",
			input:    "_",
			expected: "",
		},
		{
			name:     "Multiple underscores only",
			input:    "___",
			expected: "",
		},
		{
			name:     "Underscore at beginning",
			input:    "_test_case",
			expected: "TestCase",
		},
		{
			name:     "Numbers with underscores",
			input:    "test_123_case",
			expected: "test123Case",
		},
		{
			name:     "Special characters preserved",
			input:    "test-case_name",
			expected: "test-caseName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := SnakeToCamelCase(tt.input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCamelToSnakeCase_EdgeCases tests edge cases for camel to snake conversion
func TestCamelToSnakeCase_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "All uppercase",
			input:    "ALLUPPERCASE",
			expected: "a_l_l_u_p_p_e_r_c_a_s_e",
		},
		{
			name:     "Acronym at start",
			input:    "XMLHttpRequest",
			expected: "x_m_l_http_request",
		},
		{
			name:     "Acronym at end",
			input:    "parseXML",
			expected: "parse_x_m_l",
		},
		{
			name:     "Single uppercase letter",
			input:    "A",
			expected: "a",
		},
		{
			name:     "Mixed with numbers",
			input:    "base64Encode",
			expected: "base64_encode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := CamelToSnakeCase(tt.input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestUtilityFunctionsCombined tests combining multiple utility functions
func TestUtilityFunctionsCombined(t *testing.T) {
	// given
	originalSnake := "test_field_name"

	// when
	camelCase := SnakeToCamelCase(originalSnake)
	backToSnake := CamelToSnakeCase(camelCase)

	// then
	assert.Equal(t, "testFieldName", camelCase)
	assert.Equal(t, "test_field_name", backToSnake)
	assert.Equal(t, originalSnake, backToSnake)
}

// TestEncodeBase64_WithWrapInQuotes tests combining EncodeBase64 functionality
func TestEncodeBase64_WithWrapInQuotes(t *testing.T) {
	// given
	input := "test data"

	// when
	encoded := EncodeBase64(input)

	// then
	// EncodeBase64 already wraps in quotes
	assert.True(t, len(encoded) > 2)
	assert.Equal(t, byte('"'), encoded[0])
	assert.Equal(t, byte('"'), encoded[len(encoded)-1])

	// Verify it's valid base64 between quotes
	base64Part := encoded[1 : len(encoded)-1]
	assert.NotEmpty(t, base64Part)
}
