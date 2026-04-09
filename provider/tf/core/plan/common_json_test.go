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

package plan

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestCheckIfFieldIsNullOrUnknown(t *testing.T) {
	tests := []struct {
		name     string
		field    interface{}
		expected bool
	}{
		{
			name:     "Null field should return true",
			field:    types.StringNull(),
			expected: true,
		},
		{
			name:     "Unknown field should return true",
			field:    types.StringUnknown(),
			expected: true,
		},
		{
			name:     "Known field should return false",
			field:    types.StringValue("test"),
			expected: false,
		},
		{
			name:     "Bool null field should return true",
			field:    types.BoolNull(),
			expected: true,
		},
		{
			name:     "Bool unknown field should return true",
			field:    types.BoolUnknown(),
			expected: true,
		},
		{
			name:     "Bool known field should return false",
			field:    types.BoolValue(true),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fieldValue := reflect.ValueOf(tt.field)

			result := CheckIfFieldIsNullOrUnknown(fieldValue)

			assert.Equal(t, tt.expected, result)
		})
	}
}
