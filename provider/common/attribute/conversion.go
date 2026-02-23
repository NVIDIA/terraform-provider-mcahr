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
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConvertToInt64 converts a value to types.Int64
func ConvertToInt64(value interface{}) (types.Int64, error) {
	switch v := value.(type) {
	case float32:
		return types.Int64Value(int64(v)), nil
	case float64:
		return types.Int64Value(int64(v)), nil
	case int:
		return types.Int64Value(int64(v)), nil
	case int32:
		return types.Int64Value(int64(v)), nil
	case int64:
		return types.Int64Value(v), nil
	default:
		return types.Int64Null(), fmt.Errorf("cannot convert %T to int64", v)
	}
}

// ConvertToStringList converts a value to a types.List of strings
// Always returns an empty list for empty arrays (never null)
// Returns an error if the value is not an array or contains non-string elements
func ConvertToStringList(value interface{}) (types.List, error) {
	arr, ok := value.([]interface{})
	if !ok {
		return types.ListNull(types.StringType), fmt.Errorf("cannot convert %T to list", value)
	}

	elements := make([]attr.Value, 0, len(arr))
	for _, item := range arr {
		str, ok := item.(string)
		if !ok {
			return types.ListNull(types.StringType), fmt.Errorf("expected string list element, got %T", item)
		}
		elements = append(elements, types.StringValue(str))
	}

	listVal, diags := types.ListValue(types.StringType, elements)
	if diags.HasError() {
		return types.ListNull(types.StringType), fmt.Errorf("failed to create list value: %s", diags.Errors()[0].Summary())
	}
	return listVal, nil
}
