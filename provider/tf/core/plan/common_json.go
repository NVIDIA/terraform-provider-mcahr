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
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// CheckIfFieldIsNullOrUnknown checks if a Terraform field is null or unknown using reflection
func CheckIfFieldIsNullOrUnknown(resultField reflect.Value) bool {
	// Check if the field has IsNull and IsUnknown methods (Terraform types)
	if resultField.CanInterface() {
		resultInterface := resultField.Interface()

		// Use reflection to check if the field has IsNull and IsUnknown methods
		resultFieldValue := reflect.ValueOf(resultInterface)
		if resultFieldValue.IsValid() {
			// Check if the field has IsNull method
			isNullMethod := resultFieldValue.MethodByName("IsNull")
			isUnknownMethod := resultFieldValue.MethodByName("IsUnknown")

			if isNullMethod.IsValid() && isUnknownMethod.IsValid() {
				// Call IsNull() and IsUnknown() methods
				isNullResult := isNullMethod.Call(nil)
				isUnknownResult := isUnknownMethod.Call(nil)

				// Check if the field is null or unknown
				if len(isNullResult) > 0 && len(isUnknownResult) > 0 {
					isNull := GetInvocationResult(isNullResult)
					isUnknown := GetInvocationResult(isUnknownResult)

					return isNull || isUnknown
				}
			}
		}
	}
	return false
}

// GetInvocationResult extracts a boolean result from reflection method invocation
func GetInvocationResult(result []reflect.Value) bool {
	return result[0].Bool()
}

// AddDefaultsFromPlan copies field values from planValues to resultValues for fields that are null or unknown
// T should be a struct type that contains Terraform framework types
func AddDefaultsFromPlan[T any](resultValues *T, planValues *T) {
	resultValue := reflect.ValueOf(resultValues).Elem()
	planValue := reflect.ValueOf(planValues).Elem()

	// Iterate through all fields in the struct
	for i := 0; i < resultValue.NumField(); i++ {
		resultField := resultValue.Field(i)
		planField := planValue.Field(i)

		if CheckIfFieldIsNullOrUnknown(resultField) && planField.CanInterface() && resultField.CanSet() {
			resultField.Set(planField)
		}
	}
}

// GetValues extracts plan, config, and state values from the ModifyPlanRequest
// T should be a struct type that contains Terraform framework types
// Returns (doReturn bool, planValues T, configValues T, stateValues T)
func GetValues[T any](ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) (bool, T, T, T) {
	var planValues, configValues, stateValues T

	if !req.Plan.Raw.IsNull() {
		diags := req.Plan.Get(ctx, &planValues)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return true, planValues, configValues, stateValues
		}
	} else {
		// It's a destroy operation, do nothing
		return true, planValues, configValues, stateValues
	}

	if !req.Config.Raw.IsNull() {
		diags := req.Config.Get(ctx, &configValues)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return true, planValues, configValues, stateValues
		}
	}

	if !req.State.Raw.IsNull() {
		diags := req.State.Get(ctx, &stateValues)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return true, planValues, configValues, stateValues
		}
	}

	return false, planValues, configValues, stateValues
}
