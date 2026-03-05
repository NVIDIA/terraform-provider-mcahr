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

package process

import (
	"fmt"
	"reflect"
	"terraform/terraform-provider/provider/common"
	model "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// RestoreAllFieldsFromPlan copies non-null/non-unknown fields from the plan to the tfModel.
// This is the default behavior to ensure consistency and prevent "inconsistent result after apply" errors.
// Resource-specific PostProcessors can then override specific fields that need special handling.
//
// This function should only be called for Create and Update operations (enforced by caller).
func RestoreAllFieldsFromPlan[TF model.TFModel](requestContext *common.RequestContext, processData *ProcessData, tfModel TF) error {
	var planModel TF
	var diags diag.Diagnostics

	switch requestContext.Operation {
	case common.Create:
		diags = processData.CreateRequest.Plan.Get(requestContext.Context, &planModel)
	case common.Update:
		diags = processData.UpdateRequest.Plan.Get(requestContext.Context, &planModel)
	default:
		// This should not happen if caller enforces operation check
		return fmt.Errorf("Restore plan was called with unsupported operation: %v", requestContext.Operation)
	}

	if diags.HasError() {
		return fmt.Errorf("failed to get plan model: %s", diags.Errors())
	}

	// Copy only non-null/non-unknown fields from plan to tfModel
	// This allows computed fields set by the API to remain
	return copyNonNullFields(planModel, tfModel)
}

// copyNonNullFields uses reflection to copy only non-null and non-unknown fields from source to target.
// This ensures user-configured values match the plan while allowing API-set computed fields to remain.
func copyNonNullFields[TF model.TFModel](source, target TF) error {
	sourceVal := reflect.ValueOf(source)
	targetVal := reflect.ValueOf(target)

	// Handle pointer types
	if sourceVal.Kind() == reflect.Ptr {
		if sourceVal.IsNil() || targetVal.IsNil() {
			return nil
		}
		sourceVal = sourceVal.Elem()
		targetVal = targetVal.Elem()
	}

	if sourceVal.Kind() != reflect.Struct || targetVal.Kind() != reflect.Struct {
		return nil
	}

	sourceType := sourceVal.Type()

	// Copy fields selectively
	for i := 0; i < sourceVal.NumField(); i++ {
		sourceField := sourceVal.Field(i)
		targetField := targetVal.Field(i)

		if shouldCopyField(sourceType.Field(i), sourceField, targetField) {
			targetField.Set(sourceField)
		}
	}

	return nil
}

// shouldCopyField determines if a field should be copied from source to target
func shouldCopyField(fieldType reflect.StructField, sourceField, targetField reflect.Value) bool {
	// Skip if target field is not settable
	if !targetField.CanSet() || !sourceField.IsValid() {
		return false
	}

	tfsdkTag := fieldType.Tag.Get("tfsdk")

	// For fields with tfsdk tag, check if the value is known
	if tfsdkTag != "" && tfsdkTag != "-" {
		// If field has tfsdk tag, it must be an attr.Value type
		// Only copy if the value is known (not null and not unknown)
		return common.IsAttrKnown(sourceField.Interface().(attr.Value))
	}

	// Copy all non-tfsdk fields
	return true
}
