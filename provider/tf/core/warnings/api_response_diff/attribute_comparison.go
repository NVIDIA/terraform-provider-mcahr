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

package apiresponsediff

import (
	"fmt"
	"reflect"
	model "terraform/terraform-provider/provider/tf/core/model"
	coremodel "terraform/terraform-provider/provider/tf/core/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func compareModels[TF model.TFModel](plan TF, apiResponse TF, comparisonRules map[string]coremodel.FieldComparisonRule) []fieldDifference {
	var differences []fieldDifference

	planValue := reflect.ValueOf(plan)
	responseValue := reflect.ValueOf(apiResponse)

	if planValue.Kind() == reflect.Ptr {
		if planValue.IsNil() || responseValue.IsNil() {
			return differences
		}
		planValue = planValue.Elem()
		responseValue = responseValue.Elem()
	}

	if planValue.Kind() != reflect.Struct {
		return differences
	}

	planType := planValue.Type()

	for i := 0; i < planValue.NumField(); i++ {
		field := planType.Field(i)
		tfsdkTag := field.Tag.Get("tfsdk")
		if tfsdkTag == "" || tfsdkTag == "-" {
			continue
		}

		planField := planValue.Field(i)
		responseField := responseValue.Field(i)

		diff := compareFieldWithRule(tfsdkTag, planField, responseField, comparisonRules)
		if diff != nil {
			differences = append(differences, *diff)
		}
	}

	return differences
}

// compareFieldWithRule compares a single field using the appropriate comparison rule
func compareFieldWithRule(fieldName string, planField, responseField reflect.Value, comparisonRules map[string]coremodel.FieldComparisonRule) *fieldDifference {
	rule, hasRule := comparisonRules[fieldName]

	if !hasRule || rule.Behavior == coremodel.CompareNormally {
		// Default comparison
		return compareField(fieldName, planField, responseField)
	}

	if rule.Behavior == coremodel.SkipComparison {
		// Skip this field entirely
		return nil
	}

	if rule.Behavior == coremodel.CustomComparison && rule.CustomCompare != nil {
		// Use custom comparison logic
		return compareFieldWithCustomLogic(fieldName, planField, responseField, rule.CustomCompare)
	}

	// Fallback to default comparison if rule is invalid
	return compareField(fieldName, planField, responseField)
}

// compareFieldWithCustomLogic applies a custom comparison function to a field
func compareFieldWithCustomLogic(fieldName string, planField, responseField reflect.Value, customCompare func(string, string, string) bool) *fieldDifference {
	planStr, planWasSet := extractFieldValueWithSetFlag(planField)
	responseStr, _ := extractFieldValueWithSetFlag(responseField)

	// Skip if plan field was not set by the user
	if !planWasSet {
		return nil
	}

	// Use custom comparison function
	// Returns true if values are considered equal, false if different
	if customCompare(fieldName, planStr, responseStr) {
		// Custom function says they're equal - no difference
		return nil
	}

	// Custom function says they're different
	return &fieldDifference{
		FieldName:     fieldName,
		PlanValue:     planStr,
		ResponseValue: responseStr,
	}
}

func compareField(fieldName string, planField, responseField reflect.Value) *fieldDifference {
	planStr, planWasSet := extractFieldValueWithSetFlag(planField)
	responseStr, _ := extractFieldValueWithSetFlag(responseField)

	// Skip if plan field was not set by the user
	if !planWasSet {
		return nil
	}

	// Only report if values differ
	if planStr != responseStr {
		return &fieldDifference{
			FieldName:     fieldName,
			PlanValue:     planStr,
			ResponseValue: responseStr,
		}
	}

	return nil
}

// extractFieldValueWithSetFlag returns the string value and whether the field was explicitly set
func extractFieldValueWithSetFlag(field reflect.Value) (string, bool) {
	if !field.IsValid() {
		return "", false
	}

	fieldType := field.Type()

	switch {
	case fieldType == reflect.TypeOf(types.String{}):
		tfString := field.Interface().(types.String)
		if tfString.IsNull() || tfString.IsUnknown() {
			return "", false
		}
		return tfString.ValueString(), true

	case fieldType == reflect.TypeOf(types.Bool{}):
		tfBool := field.Interface().(types.Bool)
		if tfBool.IsNull() || tfBool.IsUnknown() {
			return "", false
		}
		return fmt.Sprintf("%v", tfBool.ValueBool()), true

	case fieldType == reflect.TypeOf(types.Int64{}):
		tfInt := field.Interface().(types.Int64)
		if tfInt.IsNull() || tfInt.IsUnknown() {
			return "", false
		}
		return fmt.Sprintf("%d", tfInt.ValueInt64()), true

	case fieldType == reflect.TypeOf(types.Float64{}):
		tfFloat := field.Interface().(types.Float64)
		if tfFloat.IsNull() || tfFloat.IsUnknown() {
			return "", false
		}
		return fmt.Sprintf("%f", tfFloat.ValueFloat64()), true

	case fieldType == reflect.TypeOf(types.List{}):
		tfList := field.Interface().(types.List)
		if tfList.IsNull() || tfList.IsUnknown() {
			return "", false
		}
		return fmt.Sprintf("%v", tfList.Elements()), true

	case fieldType == reflect.TypeOf(types.Set{}):
		tfSet := field.Interface().(types.Set)
		if tfSet.IsNull() || tfSet.IsUnknown() {
			return "", false
		}
		return fmt.Sprintf("%v", tfSet.Elements()), true

	case fieldType == reflect.TypeOf(types.Map{}):
		tfMap := field.Interface().(types.Map)
		if tfMap.IsNull() || tfMap.IsUnknown() {
			return "", false
		}
		return fmt.Sprintf("%v", tfMap.Elements()), true

	case fieldType == reflect.TypeOf(types.Object{}):
		tfObj := field.Interface().(types.Object)
		if tfObj.IsNull() || tfObj.IsUnknown() {
			return "", false
		}
		return fmt.Sprintf("%v", tfObj.Attributes()), true

	default:
		return "", false
	}
}
