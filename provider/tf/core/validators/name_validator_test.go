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

func TestNameValidator_ValidateString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		// Given
		value string
		// When/Then
		expectError bool
		errorCount  int
	}{
		"valid_simple_name": {
			// Given
			value: "my_action",
			// Then
			expectError: false,
		},
		"valid_starts_with_letter": {
			// Given
			value: "action123",
			// Then
			expectError: false,
		},
		"valid_starts_with_underscore": {
			// Given
			value: "_private_action",
			// Then
			expectError: false,
		},
		"valid_all_letters": {
			// Given
			value: "MyAction",
			// Then
			expectError: false,
		},
		"valid_all_uppercase": {
			// Given
			value: "CPU_ALARM",
			// Then
			expectError: false,
		},
		"valid_mixed_case": {
			// Given
			value: "CpuAlarm_v2",
			// Then
			expectError: false,
		},
		"invalid_starts_with_number": {
			// Given
			value: "123action",
			// Then
			expectError: true,
			errorCount:  1,
		},
		"invalid_contains_hyphen": {
			// Given
			value: "my-action",
			// Then
			expectError: true,
			errorCount:  1,
		},
		"invalid_contains_space": {
			// Given
			value: "my action",
			// Then
			expectError: true,
			errorCount:  1,
		},
		"invalid_contains_special_chars": {
			// Given
			value: "action@123",
			// Then
			expectError: true,
			errorCount:  1,
		},
		"invalid_starts_with_special_char": {
			// Given
			value: "@action",
			// Then
			expectError: true,
			errorCount:  1, // Only invalid character (returns early)
		},
		"empty_string": {
			// Given
			value: "",
			// Then
			expectError: true,
			errorCount:  1, // Invalid start (empty doesn't start with letter/underscore)
		},
		"only_numbers": {
			// Given
			value: "123",
			// Then
			expectError: true,
			errorCount:  1, // Invalid start
		},
		"only_underscores": {
			// Given
			value: "___",
			// Then
			expectError: false, // Valid - starts with underscore and contains only underscores
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			ctx := context.Background()
			request := validator.StringRequest{
				Path:        path.Root("name"),
				ConfigValue: types.StringValue(test.value),
			}
			response := &validator.StringResponse{}

			validator := NameValidator()

			// When
			validator.ValidateString(ctx, request, response)

			// Then
			hasError := response.Diagnostics.HasError()
			if test.expectError != hasError {
				t.Errorf("Expected error: %v, got error: %v", test.expectError, hasError)
			}

			if test.expectError && test.errorCount > 0 {
				actualCount := response.Diagnostics.ErrorsCount()
				if actualCount != test.errorCount {
					t.Errorf("Expected %d errors, got %d errors", test.errorCount, actualCount)
				}
			}
		})
	}
}

func TestNameValidator_ValidateString_NullAndUnknown(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		// Given
		value types.String
		// Then
		expectError bool
	}{
		"null_value": {
			// Given
			value: types.StringNull(),
			// Then
			expectError: false, // null values should be allowed (handled by Required/Optional)
		},
		"unknown_value": {
			// Given
			value: types.StringUnknown(),
			// Then
			expectError: false, // unknown values should be allowed during planning
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			ctx := context.Background()
			request := validator.StringRequest{
				Path:        path.Root("name"),
				ConfigValue: test.value,
			}
			response := &validator.StringResponse{}

			validator := NameValidator()

			// When
			validator.ValidateString(ctx, request, response)

			// Then
			hasError := response.Diagnostics.HasError()
			if test.expectError != hasError {
				t.Errorf("Expected error: %v, got error: %v", test.expectError, hasError)
			}
		})
	}
}
