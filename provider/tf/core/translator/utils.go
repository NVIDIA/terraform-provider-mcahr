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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListSliceFromTFModel extracts a string slice from types.List, returning a
// nil slice for null/unknown.
//
// Conversion failures are returned rather than discarded. These helpers feed
// SetArrayField for attributes including allowed_entities, approvers, editors
// and allowed_tags: swallowing the diagnostics turned a failed conversion into
// an empty list, which sends allowed_entities=[] and wipes the entity's
// authorization list instead of failing the apply.
func ListSliceFromTFModel(ctx context.Context, tfList types.List) ([]string, error) {
	// nil, not an empty slice: ElementsAs left the result nil for null/unknown
	// lists, and callers marshal these straight to JSON where nil renders as
	// null and []string{} renders as []. Preserving nil keeps the wire format
	// identical -- only the error handling changes here.
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil, nil
	}

	var result []string
	if diags := tfList.ElementsAs(ctx, &result, false); diags.HasError() {
		return nil, diagnosticsError("list", diags)
	}

	return result, nil
}

// SetSliceFromTFModel extracts a string slice from types.Set, returning a nil
// slice for null/unknown. See ListSliceFromTFModel on error handling.
func SetSliceFromTFModel(ctx context.Context, tfSet types.Set) ([]string, error) {
	if tfSet.IsNull() || tfSet.IsUnknown() {
		return nil, nil
	}

	var result []string
	if diags := tfSet.ElementsAs(ctx, &result, false); diags.HasError() {
		return nil, diagnosticsError("set", diags)
	}

	return result, nil
}

// ListSlicesFromTFModel extracts several string slices in one call, keyed by
// attribute name so a failure says which attribute could not be converted and
// the builder chain at the call site stays readable.
func ListSlicesFromTFModel(ctx context.Context, lists map[string]types.List) (map[string][]string, error) {
	result := make(map[string][]string, len(lists))

	for name, tfList := range lists {
		slice, err := ListSliceFromTFModel(ctx, tfList)
		if err != nil {
			return nil, fmt.Errorf("attribute %q: %w", name, err)
		}
		result[name] = slice
	}

	return result, nil
}

// diagnosticsError renders element-conversion diagnostics as an error.
func diagnosticsError(kind string, diags diag.Diagnostics) error {
	messages := make([]string, 0, len(diags.Errors()))
	for _, d := range diags.Errors() {
		messages = append(messages, fmt.Sprintf("%s: %s", d.Summary(), d.Detail()))
	}

	return fmt.Errorf("failed to convert %s elements: %s", kind, strings.Join(messages, "; "))
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
