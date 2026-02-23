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

package commonstruct

import (
	"reflect"
	"strings"
	"terraform/terraform-provider/provider/common/version"
)

func ApplyCustomStructTags(structData any, opts map[string]any) map[string]any {

	result := make(map[string]any)

	v := reflect.ValueOf(structData)
	t := reflect.TypeOf(structData)

	// XTODO : improvement - make this function use the utility functions from the struct iteration file

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if field.Tag.Get("skip") == "true" {
			continue
		}

		fieldName := getFieldName(field)

		// Skip fields with empty names
		if fieldName == "" {
			continue
		}

		// Check if the field is supported in the current backend version
		minVersionStr := field.Tag.Get("min_version")
		maxVersionStr := field.Tag.Get("max_version")

		skipVersionPrefixes := field.Tag.Get("skip_version_prefixes")
		var skipVersionPrefixesList []string
		if skipVersionPrefixes != "" {
			skipVersionPrefixesList = strings.Split(skipVersionPrefixes, ",")
		}

		backendVer, ok := opts["backend_version"].(*version.BackendVersion)
		if ok && !version.IsFieldSupported(backendVer, minVersionStr, maxVersionStr, skipVersionPrefixesList) {
			continue
		}

		// ... Add other custom tags here if needed

		result[fieldName] = value.Interface()
	}

	return result
}

func getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return ""
	}
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}
