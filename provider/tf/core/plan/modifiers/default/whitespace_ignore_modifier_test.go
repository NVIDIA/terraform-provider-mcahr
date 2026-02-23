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
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIgnoreWhitespaceModifier_PlanModifyString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		planValue         types.String
		stateValue        types.String
		configValue       types.String
		expectedPlanValue types.String
		shouldModify      bool
	}{
		{
			name:              "identical strings - no modification",
			planValue:         types.StringValue("hello world"),
			stateValue:        types.StringValue("hello world"),
			configValue:       types.StringValue("hello world"),
			expectedPlanValue: types.StringValue("hello world"),
			shouldModify:      false,
		},
		{
			name:              "plan has extra spaces - should use state value",
			planValue:         types.StringValue("hello  world"),
			stateValue:        types.StringValue("hello world"),
			configValue:       types.StringValue("hello  world"),
			expectedPlanValue: types.StringValue("hello world"),
			shouldModify:      true,
		},
		{
			name:              "state has extra spaces - should use state value",
			planValue:         types.StringValue("hello world"),
			stateValue:        types.StringValue("hello  world"),
			configValue:       types.StringValue("hello world"),
			expectedPlanValue: types.StringValue("hello  world"),
			shouldModify:      true,
		},
		{
			name:              "multiple spaces normalized - should use state value",
			planValue:         types.StringValue("hello    world    test"),
			stateValue:        types.StringValue("helloworld test"),
			configValue:       types.StringValue("hello    world    test"),
			expectedPlanValue: types.StringValue("helloworld test"),
			shouldModify:      true,
		},
		{
			name:              "completely different strings - no modification",
			planValue:         types.StringValue("hello world"),
			stateValue:        types.StringValue("bye world"),
			configValue:       types.StringValue("hello world"),
			expectedPlanValue: types.StringValue("hello world"),
			shouldModify:      false,
		},
		{
			name:              "plan value unknown - no modification",
			planValue:         types.StringUnknown(),
			stateValue:        types.StringValue("hello world"),
			configValue:       types.StringValue("hello world"),
			expectedPlanValue: types.StringUnknown(),
			shouldModify:      false,
		},
		{
			name:              "plan value null - no modification",
			planValue:         types.StringNull(),
			stateValue:        types.StringValue("hello world"),
			configValue:       types.StringValue("hello world"),
			expectedPlanValue: types.StringNull(),
			shouldModify:      false,
		},
		{
			name:              "state value unknown - no modification",
			planValue:         types.StringValue("hello world"),
			stateValue:        types.StringUnknown(),
			configValue:       types.StringValue("hello world"),
			expectedPlanValue: types.StringValue("hello world"),
			shouldModify:      false,
		},
		{
			name:              "state value null - no modification",
			planValue:         types.StringValue("hello world"),
			stateValue:        types.StringNull(),
			configValue:       types.StringValue("hello world"),
			expectedPlanValue: types.StringValue("hello world"),
			shouldModify:      false,
		},
		{
			name:              "both values null - no modification",
			planValue:         types.StringNull(),
			stateValue:        types.StringNull(),
			configValue:       types.StringNull(),
			expectedPlanValue: types.StringNull(),
			shouldModify:      false,
		},
		{
			name:              "both values unknown - no modification",
			planValue:         types.StringUnknown(),
			stateValue:        types.StringUnknown(),
			configValue:       types.StringUnknown(),
			expectedPlanValue: types.StringUnknown(),
			shouldModify:      false,
		},
		{
			name:              "empty strings - no modification needed",
			planValue:         types.StringValue(""),
			stateValue:        types.StringValue(""),
			configValue:       types.StringValue(""),
			expectedPlanValue: types.StringValue(""),
			shouldModify:      false,
		},
		{
			name:              "plan empty, state with spaces - should use state value",
			planValue:         types.StringValue(""),
			stateValue:        types.StringValue("   "),
			configValue:       types.StringValue(""),
			expectedPlanValue: types.StringValue("   "),
			shouldModify:      true,
		},
		{
			name:              "complex whitespace differences - should use state value",
			planValue:         types.StringValue("SELECT * FROM table WHERE id = 1"),
			stateValue:        types.StringValue("SELECT*FROMtableWHEREid=1"),
			configValue:       types.StringValue("SELECT * FROM table WHERE id = 1"),
			expectedPlanValue: types.StringValue("SELECT*FROMtableWHEREid=1"),
			shouldModify:      true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Given
			modifier := IgnoreWhitespaceModifier()
			ctx := context.Background()
			req := planmodifier.StringRequest{
				PlanValue:   tt.planValue,
				StateValue:  tt.stateValue,
				ConfigValue: tt.configValue,
			}
			resp := &planmodifier.StringResponse{
				PlanValue: tt.planValue, // Initialize with plan value
			}

			// When
			modifier.PlanModifyString(ctx, req, resp)

			// Then
			if !resp.PlanValue.Equal(tt.expectedPlanValue) {
				t.Errorf("PlanModifyString() planValue = %v, expected %v", resp.PlanValue, tt.expectedPlanValue)
			}

			// Verify modification behavior
			wasModified := !resp.PlanValue.Equal(tt.planValue)
			if wasModified != tt.shouldModify {
				t.Errorf("PlanModifyString() modification behavior = %v, expected %v", wasModified, tt.shouldModify)
			}
		})
	}
}

func TestRemoveWhitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no spaces",
			input:    "helloworld",
			expected: "helloworld",
		},
		{
			name:     "single space",
			input:    "hello world",
			expected: "helloworld",
		},
		{
			name:     "multiple spaces",
			input:    "hello   world",
			expected: "helloworld",
		},
		{
			name:     "leading space",
			input:    " hello",
			expected: "hello",
		},
		{
			name:     "trailing space",
			input:    "hello ",
			expected: "hello",
		},
		{
			name:     "spaces everywhere",
			input:    "  hello   world  ",
			expected: "helloworld",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "mixed content with spaces",
			input:    "SELECT * FROM table WHERE id = 1",
			expected: "SELECT*FROMtableWHEREid=1",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Given
			input := tt.input

			// When
			result := removeWhitespace(input)

			// Then
			if result != tt.expected {
				t.Errorf("removeWhitespace() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
