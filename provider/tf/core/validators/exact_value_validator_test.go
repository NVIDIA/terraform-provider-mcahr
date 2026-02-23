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

package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Helper function to validate error messages in validator responses
func assertValidatorError(t *testing.T, response *validator.StringResponse, expectError bool, expectedErrorMsg string) {
	hasError := response.Diagnostics.HasError()
	if expectError != hasError {
		t.Errorf("Expected error: %v, got error: %v", expectError, hasError)
	}

	if expectError && expectedErrorMsg != "" {
		if !hasError {
			t.Errorf("Expected error message but got no error")
			return
		}

		// Check if any of the error messages contains the expected message
		errorMessages := response.Diagnostics.Errors()
		found := false
		for _, err := range errorMessages {
			if err.Detail() == expectedErrorMsg {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected error message: '%s', but got: %v", expectedErrorMsg, errorMessages[0].Detail())
		}
	}
}

func TestExactValueValidator_ValidateString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		// Given
		requiredValue string
		inputValue    string
		// When/Then
		expectError      bool
		expectedErrorMsg string
	}{
		"exact_match_system_settings": {
			requiredValue: "system_settings",
			inputValue:    "system_settings",
			expectError:   false,
		},
		"exact_match_different_value": {
			requiredValue: "production",
			inputValue:    "production",
			expectError:   false,
		},
		"mismatch_wrong_value": {
			requiredValue:    "system_settings",
			inputValue:       "wrong_name",
			expectError:      true,
			expectedErrorMsg: "The value must be exactly 'system_settings', got: 'wrong_name'",
		},
		"mismatch_empty_string": {
			requiredValue:    "system_settings",
			inputValue:       "",
			expectError:      true,
			expectedErrorMsg: "The value must be exactly 'system_settings', got: ''",
		},
		"mismatch_similar_value": {
			requiredValue:    "system_settings",
			inputValue:       "system_setting", // missing 's'
			expectError:      true,
			expectedErrorMsg: "The value must be exactly 'system_settings', got: 'system_setting'",
		},
		"mismatch_case_sensitive": {
			requiredValue:    "system_settings",
			inputValue:       "System_Settings", // different case
			expectError:      true,
			expectedErrorMsg: "The value must be exactly 'system_settings', got: 'System_Settings'",
		},
		"mismatch_extra_characters": {
			requiredValue:    "system_settings",
			inputValue:       "system_settings_extra",
			expectError:      true,
			expectedErrorMsg: "The value must be exactly 'system_settings', got: 'system_settings_extra'",
		},
		"empty_string_validator": {
			requiredValue: "",
			inputValue:    "",
			expectError:   false,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			ctx := context.Background()
			request := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(test.inputValue),
			}
			response := &validator.StringResponse{}

			validator := ExactValueValidator(test.requiredValue)

			// When
			validator.ValidateString(ctx, request, response)

			// Then - Test validation logic
			assertValidatorError(t, response, test.expectError, test.expectedErrorMsg)
		})
	}
}
func TestExactValueValidator_ValidateString_NullAndUnknown(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		// Given
		value types.String
		// Then
		expectError      bool
		expectedErrorMsg string
	}{
		"null_value": {
			// Given
			value: types.StringNull(),
			// Then
			expectError:      true, // null values should trigger validation errors for singleton resources
			expectedErrorMsg: "The value must be exactly 'system_settings', got: ''",
		},
		"unknown_value": {
			// Given
			value: types.StringUnknown(),
			// Then
			expectError:      true, // unknown values should trigger validation errors for singleton resources
			expectedErrorMsg: "The value must be exactly 'system_settings', got: ''",
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			ctx := context.Background()
			request := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: test.value,
			}
			response := &validator.StringResponse{}

			validator := ExactValueValidator("system_settings")

			// When
			validator.ValidateString(ctx, request, response)

			// Then - Test validation logic
			assertValidatorError(t, response, test.expectError, test.expectedErrorMsg)
		})
	}
}
