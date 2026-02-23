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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListSliceFromTFModel extracts a string slice from types.List, returning empty slice for null/unknown
func ListSliceFromTFModel(ctx context.Context, tfList types.List) []string {
	var result []string
	tfList.ElementsAs(ctx, &result, false)
	return result
}

// SetSliceFromTFModel extracts a string slice from types.Set, returning empty slice for null/unknown
func SetSliceFromTFModel(ctx context.Context, tfSet types.Set) []string {
	var result []string
	tfSet.ElementsAs(ctx, &result, false)
	return result
}

// EscapeString escapes strings for op lang format using strconv.Quote
func EscapeString(s string) string {
	if IsEmpty(s) {
		return "\"\""
	}
	return strconv.Quote(s)
}

// arrayToOpLang converts string slice to op lang array format
func ArrayToOpLang(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	var escaped []string
	for _, item := range arr {
		escaped = append(escaped, EscapeString(item))
	}
	return "[" + strings.Join(escaped, ", ") + "]"
}

// parseStringArray parses JSON string arrays commonly found in API responses
func ParseStringArray(jsonStr string) []string {
	if IsEmpty(jsonStr) || jsonStr == "[]" {
		return []string{}
	}
	var result []string
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return s == ""
}

// BoolToInt converts a boolean to an integer (1 for true, 0 for false)
func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// IntToBool converts an integer to a boolean (0 for false, any other value for true)
// This matches the API response format for permission fields
func IntToBool(value int) bool {
	return value != 0
}

// EncodeBase64 encodes a JSON string to base64 with quotes for API calls
func EncodeBase64(jsonStr string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(jsonStr))
	return fmt.Sprintf("\"%s\"", encoded)
}

// ListValueFromStringSlice converts a slice of strings to a Terraform list.
// It returns an empty list when the input slice is nil to avoid null values.
func ListValueFromStringSlice(ctx context.Context, values []string) types.List {
	if values == nil {
		values = []string{}
	}

	list, diags := types.ListValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}

	return list
}
