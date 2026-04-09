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

package data

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"terraform/terraform-provider/provider/common"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IsDataJSONEmpty checks if data JSON is empty
func IsDataJSONEmpty(data *types.String) bool {
	if data == nil || data.IsNull() || data.IsUnknown() || data.ValueString() == "" || data.ValueString() == "{}" {
		return true
	}
	return false
}

// IsFieldInDataJSON checks if a field is present in the data JSON
func IsFieldInDataJSON(fieldName string, dataMap map[string]interface{}) bool {
	return findValueInMap(fieldName, dataMap) != nil
}

// GetFieldValue tries to get value from data fields using both snake_case and camelCase
func GetFieldValue(fieldName string, dataMap map[string]interface{}) interface{} {
	return findValueInMap(fieldName, dataMap)
}

// IsJSONField determines if a field should be treated as JSON (exported for other packages)
func IsJSONField(fieldName string) bool {
	jsonFields := map[string]bool{
		"params":          true,
		"cells":           true,
		"external_params": true,
	}
	return jsonFields[fieldName]
}

func IsJSONSkipField(fieldName string) bool {
	return strings.HasSuffix(fieldName, "_full") || fieldName == "data"
}

// IsDeprecatedAliasTarget returns true if the field is a deprecated field that has been
// replaced by an aliased field (e.g., cells → cells_list). These fields should be skipped
// during data apply (the replacement reads from the same data JSON key) but NOT during
// conflict validation (to detect root vs data conflicts on the deprecated field).
func IsDeprecatedAliasTarget(fieldName string) bool {
	_, isDeprecated := deprecatedAliasTargets[fieldName]
	return isDeprecated
}

// dataFieldAliases maps struct field names to their data JSON equivalents
// when the names differ (e.g., cells_list in struct → cells in data JSON).
var dataFieldAliases = map[string]string{
	"cells_list":           "cells",
	"params_list":          "params",
	"external_params_list": "external_params",
}

// deprecatedAliasTargets is the reverse of dataFieldAliases — maps deprecated data JSON
// field names back to their replacement struct field names. Pre-computed to avoid looping.
var deprecatedAliasTargets = func() map[string]string {
	m := make(map[string]string, len(dataFieldAliases))
	for replacement, deprecated := range dataFieldAliases {
		m[deprecated] = replacement
	}
	return m
}()

// ResolveDataFieldName returns the data JSON field name for a given struct field name,
// applying aliases for migrated fields (e.g., cells_list → cells).
func ResolveDataFieldName(fieldName string) string {
	if alias, ok := dataFieldAliases[fieldName]; ok {
		return alias
	}
	return fieldName
}

// OnEachStructField iterates over all fields of the model and calls the fieldFunc for each field
func OnEachStructField(ctx context.Context, tfModel *runbooktf.RunbookTFModel, fieldFunc func(fieldName string, fieldValue *reflect.Value) error) error {
	// Use reflection to process all model fields
	modelValue := reflect.ValueOf(tfModel).Elem()
	modelType := modelValue.Type()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// Get the JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Extract field name from JSON tag
		fieldName := strings.Split(jsonTag, ",")[0]

		if err := fieldFunc(fieldName, &fieldValue); err != nil {
			return fmt.Errorf("failed to process field %s from data JSON: %w", fieldName, err)
		}

	}

	return nil
}

// findValueInMap looks for a value in the map using both snake_case and camelCase
func findValueInMap(snakeCaseName string, dataMap map[string]interface{}) interface{} {
	// Try exact match first (snake_case)
	if value, ok := dataMap[snakeCaseName]; ok {
		return value
	}

	// Try camelCase
	camelCaseName := common.SnakeToCamelCase(snakeCaseName)
	if camelCaseName != snakeCaseName {
		if value, ok := dataMap[camelCaseName]; ok {
			return value
		}
	}

	return nil
}

func ParseDataJSONToMap(dataJSON types.String) (map[string]interface{}, error) {

	if IsDataJSONEmpty(&dataJSON) {
		return nil, nil
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(dataJSON.ValueString()), &dataMap); err != nil {
		return nil, fmt.Errorf("failed to parse data JSON: %w", err)
	}
	return dataMap, nil
}
