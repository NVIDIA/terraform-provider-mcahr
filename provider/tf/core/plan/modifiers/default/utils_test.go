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

package modifiers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsValueKnown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    TerraformValue
		expected bool
	}{
		{
			name:     "known string value",
			value:    types.StringValue("test"),
			expected: true,
		},
		{
			name:     "null value",
			value:    types.StringNull(),
			expected: false,
		},
		{
			name:     "unknown value",
			value:    types.StringUnknown(),
			expected: false,
		},
		{
			name:     "empty string value",
			value:    types.StringValue(""),
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Given
			value := tt.value

			// When
			result := isValueKnown(value)

			// Then
			if result != tt.expected {
				t.Errorf("isValueKnown() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsPlanOrConfigKnown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		planValue   TerraformValue
		configValue TerraformValue
		expected    bool
	}{
		{
			name:        "both known",
			planValue:   types.StringValue("plan"),
			configValue: types.StringValue("config"),
			expected:    true,
		},
		{
			name:        "plan known, config null",
			planValue:   types.StringValue("plan"),
			configValue: types.StringNull(),
			expected:    true,
		},
		{
			name:        "plan null, config known",
			planValue:   types.StringNull(),
			configValue: types.StringValue("config"),
			expected:    true,
		},
		{
			name:        "both null",
			planValue:   types.StringNull(),
			configValue: types.StringNull(),
			expected:    false,
		},
		{
			name:        "both unknown",
			planValue:   types.StringUnknown(),
			configValue: types.StringUnknown(),
			expected:    false,
		},
		{
			name:        "plan unknown, config null",
			planValue:   types.StringUnknown(),
			configValue: types.StringNull(),
			expected:    false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Given
			planValue := tt.planValue
			configValue := tt.configValue

			// When
			result := IsPlanOrConfigKnown(planValue, configValue)

			// Then
			if result != tt.expected {
				t.Errorf("IsPlanOrConfigKnown() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsPlanOrStateUnknown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		planValue  TerraformValue
		stateValue TerraformValue
		expected   bool
	}{
		{
			name:       "both known",
			planValue:  types.StringValue("plan"),
			stateValue: types.StringValue("state"),
			expected:   false,
		},
		{
			name:       "plan unknown",
			planValue:  types.StringUnknown(),
			stateValue: types.StringValue("state"),
			expected:   true,
		},
		{
			name:       "plan null",
			planValue:  types.StringNull(),
			stateValue: types.StringValue("state"),
			expected:   true,
		},
		{
			name:       "state unknown",
			planValue:  types.StringValue("plan"),
			stateValue: types.StringUnknown(),
			expected:   true,
		},
		{
			name:       "state null",
			planValue:  types.StringValue("plan"),
			stateValue: types.StringNull(),
			expected:   true,
		},
		{
			name:       "both unknown",
			planValue:  types.StringUnknown(),
			stateValue: types.StringUnknown(),
			expected:   true,
		},
		{
			name:       "both null",
			planValue:  types.StringNull(),
			stateValue: types.StringNull(),
			expected:   true,
		},
		{
			name:       "plan unknown, state null",
			planValue:  types.StringUnknown(),
			stateValue: types.StringNull(),
			expected:   true,
		},
		{
			name:       "plan null, state unknown",
			planValue:  types.StringNull(),
			stateValue: types.StringUnknown(),
			expected:   true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Given
			planValue := tt.planValue
			stateValue := tt.stateValue

			// When
			result := IsPlanOrStateUnknown(planValue, stateValue)

			// Then
			if result != tt.expected {
				t.Errorf("IsPlanOrStateUnknown() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
