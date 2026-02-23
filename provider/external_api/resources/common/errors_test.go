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

func TestFormatV2Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   string
		errors   []Error
		expected string
	}{
		{
			name:     "no errors returns empty string",
			status:   "OP_COMPLETED",
			errors:   []Error{},
			expected: "",
		},
		{
			name:     "nil errors returns empty string",
			status:   "OP_COMPLETED",
			errors:   nil,
			expected: "",
		},
		{
			name:   "single error",
			status: "OP_FAILED",
			errors: []Error{
				{Type: "VALIDATION_ERROR", Message: "Field is required"},
			},
			expected: "Status: OP_FAILED; Errors: VALIDATION_ERROR: Field is required",
		},
		{
			name:   "multiple errors",
			status: "OP_FAILED",
			errors: []Error{
				{Type: "VALIDATION_ERROR", Message: "name is required"},
				{Type: "VALIDATION_ERROR", Message: "blocks must not be empty"},
			},
			expected: "Status: OP_FAILED; Errors: VALIDATION_ERROR: name is required, VALIDATION_ERROR: blocks must not be empty",
		},
		{
			name:   "errors with different types",
			status: "OP_FAILED",
			errors: []Error{
				{Type: "VALIDATION_ERROR", Message: "Invalid input"},
				{Type: "AUTHORIZATION_ERROR", Message: "Insufficient permissions"},
			},
			expected: "Status: OP_FAILED; Errors: VALIDATION_ERROR: Invalid input, AUTHORIZATION_ERROR: Insufficient permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := FormatErrors(tt.status, tt.errors)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}
