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

package datavalidator

import (
	"fmt"

	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"
)

// validateArrayField validates each element in data.<fieldName> using the provided callback.
// One loop per field: the callback runs all checks for a single element.
func validateArrayField(fieldName string, dataMap map[string]any, validate func(index int, obj map[string]interface{}) error) error {
	raw, ok := dataMap[fieldName]
	if !ok {
		return nil
	}

	items, ok := raw.([]interface{})
	if !ok {
		return fmt.Errorf("data.%s must be an array, got %T", fieldName, raw)
	}

	for index, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("data.%s[%d] must be an object, got %T", fieldName, index, item)
		}
		if err := validate(index, obj); err != nil {
			return err
		}
	}

	return nil
}

// --- Per-field validators ---

func validateDataCells(dataMap map[string]any) error {
	return validateArrayField("cells", dataMap, func(index int, element map[string]interface{}) error {

		return validateCellType(index, element)
	})
}

// --- Shared element checks ---

func validateCellType(index int, element map[string]interface{}) error {
	cellType, _ := element["type"].(string)
	if !isValidCellType(cellType) {
		return fmt.Errorf("data.cells[%d] has invalid or missing \"type\" field %q (must be %q or %q)",
			index, cellType, customattribute.OP_LANG_TYPE, customattribute.MARKDOWN_TYPE)
	}
	return nil
}

func isValidCellType(cellType string) bool {
	switch cellType {
	case customattribute.OP_LANG_TYPE, customattribute.MARKDOWN_TYPE:
		return true
	default:
		return false
	}
}
