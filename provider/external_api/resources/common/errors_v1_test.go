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

	"github.com/stretchr/testify/assert"
)

func TestFormatV1Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      ErrorV1
		expected string
	}{
		{
			name: "OK error type returns empty string",
			err: ErrorV1{
				Type:    "OK",
				Message: "some message",
			},
			expected: "",
		},
		{
			name: "empty error type returns empty string",
			err: ErrorV1{
				Type:    "",
				Message: "some message",
			},
			expected: "",
		},
		{
			name: "error with type only",
			err: ErrorV1{
				Type:    "VALIDATION_ERROR",
				Message: "",
			},
			expected: "Error Type: VALIDATION_ERROR",
		},
		{
			name: "error with type and message",
			err: ErrorV1{
				Type:    "VALIDATION_ERROR",
				Message: "Field is required",
			},
			expected: "Error Type: VALIDATION_ERROR; Message: Field is required",
		},
		{
			name: "error with validation errors",
			err: ErrorV1{
				Type:    "VALIDATION_ERROR",
				Message: "Invalid input",
				ValidationErrors: []ValidationError{
					{Field: "name", Message: "name is required"},
					{Field: "blocks", Message: "blocks must not be empty"},
				},
			},
			expected: "Error Type: VALIDATION_ERROR; Message: Invalid input; Validation Errors: name: name is required, blocks: blocks must not be empty",
		},
		{
			name: "error with validation errors without field",
			err: ErrorV1{
				Type:    "VALIDATION_ERROR",
				Message: "Invalid input",
				ValidationErrors: []ValidationError{
					{Message: "name is required"},
					{Message: "blocks must not be empty"},
				},
			},
			expected: "Error Type: VALIDATION_ERROR; Message: Invalid input; Validation Errors: name is required, blocks must not be empty",
		},
		{
			name: "error with empty validation errors",
			err: ErrorV1{
				Type:             "VALIDATION_ERROR",
				Message:          "Invalid input",
				ValidationErrors: []ValidationError{},
			},
			expected: "Error Type: VALIDATION_ERROR; Message: Invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := FormatV1Error(tt.err)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatSyntaxErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		errors   []string
		expected string
	}{
		{
			name:     "no errors returns empty string",
			errors:   []string{},
			expected: "",
		},
		{
			name:     "nil errors returns empty string",
			errors:   nil,
			expected: "",
		},
		{
			name:     "single syntax error",
			errors:   []string{"unexpected token at line 5"},
			expected: "Syntax Errors: unexpected token at line 5",
		},
		{
			name:     "multiple syntax errors",
			errors:   []string{"unexpected token at line 5", "missing closing bracket at line 10"},
			expected: "Syntax Errors: unexpected token at line 5; missing closing bracket at line 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := FormatSyntaxErrors(tt.errors)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatDirectValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		errors   []string
		expected string
	}{
		{
			name:     "no errors returns empty string",
			errors:   []string{},
			expected: "",
		},
		{
			name:     "nil errors returns empty string",
			errors:   nil,
			expected: "",
		},
		{
			name:     "single validation error",
			errors:   []string{"field 'name' is required"},
			expected: "Validation Errors: field 'name' is required",
		},
		{
			name:     "multiple validation errors",
			errors:   []string{"field 'name' is required", "field 'value' must be positive"},
			expected: "Validation Errors: field 'name' is required; field 'value' must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := FormatDirectValidationErrors(tt.errors)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatV1ErrorsWithPriority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		syntaxErrors *SyntaxErrorsV1
		directErrors []string
		nestedError  ErrorV1
		expected     string
	}{
		{
			name:         "all errors nil/empty returns empty string",
			syntaxErrors: nil,
			directErrors: nil,
			nestedError:  ErrorV1{Type: ""},
			expected:     "",
		},
		{
			name:         "syntax errors take priority over everything",
			syntaxErrors: &SyntaxErrorsV1{Root: []string{"syntax error at line 5"}},
			directErrors: []string{"validation error"},
			nestedError:  ErrorV1{Type: "NESTED_ERROR", Message: "nested message"},
			expected:     "Syntax Errors: syntax error at line 5",
		},
		{
			name:         "direct errors take priority over nested errors",
			syntaxErrors: nil,
			directErrors: []string{"field 'name' is required"},
			nestedError:  ErrorV1{Type: "NESTED_ERROR", Message: "nested message"},
			expected:     "Validation Errors: field 'name' is required",
		},
		{
			name:         "nested error used when no syntax or direct errors",
			syntaxErrors: nil,
			directErrors: nil,
			nestedError:  ErrorV1{Type: "VALIDATION_ERROR", Message: "Invalid input"},
			expected:     "Error Type: VALIDATION_ERROR; Message: Invalid input",
		},
		{
			name:         "empty syntax errors fall through to direct errors",
			syntaxErrors: &SyntaxErrorsV1{Root: []string{}},
			directErrors: []string{"field 'value' must be positive"},
			nestedError:  ErrorV1{Type: "NESTED_ERROR", Message: "nested message"},
			expected:     "Validation Errors: field 'value' must be positive",
		},
		{
			name:         "empty direct errors fall through to nested error",
			syntaxErrors: nil,
			directErrors: []string{},
			nestedError:  ErrorV1{Type: "ERROR", Message: "Something went wrong"},
			expected:     "Error Type: ERROR; Message: Something went wrong",
		},
		{
			name:         "OK error type returns empty string",
			syntaxErrors: nil,
			directErrors: nil,
			nestedError:  ErrorV1{Type: "OK", Message: "Success"},
			expected:     "",
		},
		{
			name:         "multiple syntax errors formatted correctly",
			syntaxErrors: &SyntaxErrorsV1{Root: []string{"error 1", "error 2", "error 3"}},
			directErrors: nil,
			nestedError:  ErrorV1{Type: ""},
			expected:     "Syntax Errors: error 1; error 2; error 3",
		},
		{
			name:         "multiple direct errors formatted correctly",
			syntaxErrors: nil,
			directErrors: []string{"error 1", "error 2"},
			nestedError:  ErrorV1{Type: ""},
			expected:     "Validation Errors: error 1; error 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := FormatV1ErrorsWithPriority(tt.syntaxErrors, tt.directErrors, tt.nestedError)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

// mockV1ErrorContainer is a mock implementation of V1ErrorContainer for testing
type mockV1ErrorContainer struct {
	nestedError  ErrorV1
	directErrors []string
}

func (m *mockV1ErrorContainer) GetNestedError() ErrorV1 {
	return m.nestedError
}

func (m *mockV1ErrorContainer) GetDirectErrors() []string {
	return m.directErrors
}

func TestFormatV1ErrorsFromContainer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		syntaxErrors *SyntaxErrorsV1
		container    V1ErrorContainer
		expected     string
	}{
		{
			name:         "syntax errors take priority",
			syntaxErrors: &SyntaxErrorsV1{Root: []string{"syntax error"}},
			container: &mockV1ErrorContainer{
				directErrors: []string{"direct error"},
				nestedError:  ErrorV1{Type: "NESTED", Message: "nested error"},
			},
			expected: "Syntax Errors: syntax error",
		},
		{
			name:         "direct errors from container take priority over nested",
			syntaxErrors: nil,
			container: &mockV1ErrorContainer{
				directErrors: []string{"validation failed"},
				nestedError:  ErrorV1{Type: "NESTED", Message: "nested error"},
			},
			expected: "Validation Errors: validation failed",
		},
		{
			name:         "nested error from container used when no others",
			syntaxErrors: nil,
			container: &mockV1ErrorContainer{
				directErrors: []string{},
				nestedError:  ErrorV1{Type: "ERROR", Message: "something went wrong"},
			},
			expected: "Error Type: ERROR; Message: something went wrong",
		},
		{
			name:         "empty container returns empty string",
			syntaxErrors: nil,
			container: &mockV1ErrorContainer{
				directErrors: []string{},
				nestedError:  ErrorV1{Type: ""},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := FormatV1ErrorsFromContainer(tt.syntaxErrors, tt.container)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}
