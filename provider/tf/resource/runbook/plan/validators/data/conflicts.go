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
	"context"
	"fmt"
	"reflect"
	"strings"
	"terraform/terraform-provider/provider/common"
	data "terraform/terraform-provider/provider/tf/resource/runbook/data_attribute"
	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// crossFieldConflicts maps root config fields to data JSON fields they conflict with
// (when the names differ). Used for migrated struct fields that replace deprecated JSON fields.
var crossFieldConflicts = map[string]string{
	"cells_list":           "cells",
	"params_list":          "params",
	"external_params_list": "external_params",
}

// validateNoFieldConflicts validates that fields are set in at most one of the original TF model or data JSON
func validateNoFieldConflicts(ctx context.Context, originalTFModel *model.RunbookTFModel, dataMap map[string]any) error {

	var conflictingFields []string

	err := data.OnEachStructField(ctx, originalTFModel,
		func(fieldName string, fieldValue *reflect.Value) error {

			if data.IsJSONSkipField(fieldName) {
				return nil
			}

			// Check if field is set in original model
			originalFieldSet := common.IsAttrKnown(fieldValue.Interface().(attr.Value))

			// Check if field is present in data JSON
			dataFieldSet := data.IsFieldInDataJSON(fieldName, dataMap)

			// If both are set, it's a conflict
			if originalFieldSet && dataFieldSet {
				conflictingFields = append(conflictingFields, fieldName)
			}

			return nil
		})

	if err != nil {
		return fmt.Errorf("failed to validate field conflicts: %w", err)
	}

	// Check cross-field conflicts (e.g. root cells_list vs data cells)
	validateCrossFieldConflicts(originalTFModel, dataMap, &conflictingFields)

	if len(conflictingFields) > 0 {
		return fmt.Errorf("the following fields are set in both the root TF configuration and the data JSON: %s. Each field must be set in at most one location", strings.Join(conflictingFields, ", "))
	}

	return nil
}

// validateCrossFieldConflicts checks for conflicts between root config fields and data JSON fields
// that have different names (e.g. root cells_list vs data cells).
func validateCrossFieldConflicts(originalTFModel *model.RunbookTFModel, dataMap map[string]any, conflictingFields *[]string) {
	for rootFieldName, dataFieldName := range crossFieldConflicts {
		rootValue := getRootFieldValue(originalTFModel, rootFieldName)
		if rootValue == nil || !common.IsAttrKnown(rootValue) {
			continue
		}
		if data.IsFieldInDataJSON(dataFieldName, dataMap) {
			*conflictingFields = append(*conflictingFields,
				fmt.Sprintf("%s (conflicts with %s in data JSON)", rootFieldName, dataFieldName))
		}
	}
}

func getRootFieldValue(m *model.RunbookTFModel, jsonFieldName string) attr.Value {
	switch jsonFieldName {
	case "cells_list":
		return m.CellsList
	case "params_list":
		return m.ParamsList
	case "external_params_list":
		return m.ExternalParamsList
	default:
		return nil
	}
}
