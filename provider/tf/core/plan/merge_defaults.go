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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AddDefaultsFromPlan copies field values from planValues to resultValues for fields that are null or unknown.
// For complex fields (lists, objects), it recurses into nested structures to fill null/unknown
// values from the plan while preserving user-set values at any depth.
//
// This is needed because resultValues comes from config (which has null for unset nested fields),
// while planValues has schema defaults applied by the framework.
//
// T should be a struct type that contains Terraform framework types.
func AddDefaultsFromPlan[T any](resultValues *T, planValues *T) {
	// Get the underlying struct values via reflection
	resultValue := reflect.ValueOf(resultValues).Elem()
	planValue := reflect.ValueOf(planValues).Elem()

	// Iterate over each field in the top-level TF model struct
	for i := 0; i < resultValue.NumField(); i++ {
		resultField := resultValue.Field(i)
		planField := planValue.Field(i)

		// Skip unexported or unsettable fields
		if !resultField.CanSet() || !planField.CanInterface() {
			continue
		}

		// Both fields must be attr.Value (TF framework types: String, Bool, List, Object, etc.)
		resultAttr, ok := resultField.Interface().(attr.Value)
		if !ok {
			continue
		}
		planAttr, ok := planField.Interface().(attr.Value)
		if !ok {
			continue
		}

		// Recursively merge: fills null/unknown values from plan at any depth
		merged, didChange := mergeDefaults(resultAttr, planAttr)
		if didChange {
			resultField.Set(reflect.ValueOf(merged))
		}
	}
}

// mergeDefaults recursively merges plan defaults into a result value at any depth.
// If result is null/unknown and plan is known, uses the plan value.
// For lists, merges element-by-element. For objects, merges attribute-by-attribute.
// Primitives are returned as-is when already set.
// Returns the merged value and whether any change was made.
func mergeDefaults(result, plan attr.Value) (attr.Value, bool) {
	// If result is null/unknown, use the plan value. This includes the case where plan is
	// also unknown (e.g., category with UseStateForUnknown on first create) -- the unknown
	// must propagate so the framework shows "(known after apply)" instead of null.
	if result.IsNull() || result.IsUnknown() {
		// No change if both are in the same state (both null or both unknown)
		if (result.IsNull() && plan.IsNull()) || (result.IsUnknown() && plan.IsUnknown()) {
			return result, false
		}
		return plan, true
	}

	// If plan is null/unknown, nothing to merge from
	if plan.IsNull() || plan.IsUnknown() {
		return result, false
	}

	// Recurse into lists
	if resultList, ok := result.(types.List); ok {
		if planList, planOk := plan.(types.List); planOk {
			return mergeListDefaults(resultList, planList)
		}
	}

	// Recurse into objects
	if resultObj, ok := result.(types.Object); ok {
		if planObj, planOk := plan.(types.Object); planOk {
			return mergeObjectDefaults(resultObj, planObj)
		}
	}

	// Primitive types (String, Bool, Int64, etc.) -- already known, keep as-is
	return result, false
}

// mergeListDefaults merges plan defaults into list elements by index.
func mergeListDefaults(resultList, planList types.List) (attr.Value, bool) {
	resultElements := resultList.Elements()
	planElements := planList.Elements()

	if len(resultElements) != len(planElements) || len(resultElements) == 0 {
		return resultList, false
	}

	merged := make([]attr.Value, len(resultElements))
	changed := false

	for i := range resultElements {
		mergedElem, didChange := mergeDefaults(resultElements[i], planElements[i])
		merged[i] = mergedElem
		if didChange {
			changed = true
		}
	}

	if changed {
		ctx := context.Background()
		newList, diags := types.ListValue(resultList.ElementType(ctx), merged)
		if !diags.HasError() {
			return newList, true
		}
	}

	return resultList, false
}

// mergeObjectDefaults merges plan defaults into object attributes.
func mergeObjectDefaults(resultObj, planObj types.Object) (attr.Value, bool) {
	ctx := context.Background()
	resultAttrs := resultObj.Attributes()
	planAttrs := planObj.Attributes()
	changed := false

	for key, resultAttr := range resultAttrs {
		planAttr, exists := planAttrs[key]
		if !exists {
			continue
		}

		mergedAttr, didChange := mergeDefaults(resultAttr, planAttr)
		if didChange {
			resultAttrs[key] = mergedAttr
			changed = true
		}
	}

	if changed {
		newObj, diags := types.ObjectValue(resultObj.AttributeTypes(ctx), resultAttrs)
		if !diags.HasError() {
			return newObj, true
		}
	}

	return resultObj, false
}
