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
	"encoding/json"
	"fmt"
	"reflect"
	commonattribute "terraform/terraform-provider/provider/common/attribute"
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"
	converters "terraform/terraform-provider/provider/tf/resource/runbook/translator/object_converters"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// convertDataValueToTerraformValue converts a data value to the appropriate Terraform type
func convertDataValueToTerraformValue(fieldType reflect.Type, dataValue interface{}, fieldName string) (attr.Value, error) {
	// Create a zero value of the field type to check its type
	zeroValue := reflect.Zero(fieldType).Interface()

	switch zeroValue.(type) {
	case types.String:
		strValue, err := convertToString(dataValue, fieldName)
		if err != nil {
			return nil, err
		}
		return strValue, nil

	case types.Bool:
		if boolVal, ok := dataValue.(bool); ok {
			return types.BoolValue(boolVal), nil
		}
		return nil, fmt.Errorf("cannot convert %T to bool", dataValue)

	case types.Int64:
		int64Val, err := commonattribute.ConvertToInt64(dataValue)
		if err != nil {
			return nil, err
		}
		return int64Val, nil

	case types.List:
		setValue, err := commonattribute.ConvertToStringList(dataValue)
		if err != nil {
			return nil, err
		}
		return setValue, nil

	case types.Object:
		switch fieldName {
		case "params_groups":
			return convertToParamsGroupsObject(dataValue)

		// Add other object fields here if needed

		default:
			return nil, fmt.Errorf("unsupported object field: %s", fieldName)
		}

	default:
		return nil, fmt.Errorf("unsupported field type: %T", zeroValue)
	}
}

// convertToString converts a value to types.String
func convertToString(value interface{}, fieldName string) (types.String, error) {
	switch v := value.(type) {
	case string:
		return types.StringValue(v), nil
	case map[string]interface{}, []interface{}:
		if fieldName == "cells" {
			return convertCellsToInternalModel(v)
		}
		if IsJSONField(fieldName) {
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return types.StringNull(), err
			}
			return types.StringValue(string(jsonBytes)), nil
		}
		return types.StringNull(), fmt.Errorf("unexpected non-string value for non-JSON field")
	default:
		return types.StringNull(), fmt.Errorf("cannot convert %T to string", v)
	}
}

func convertCellsToInternalModel(inputCells interface{}) (types.String, error) {
	// cells in the data JSON use the API model of a cell
	// but in the root terraform TF model

	apiCells, ok := inputCells.([]interface{})
	if !ok {
		return types.StringNull(), fmt.Errorf("cells field is not an array of objects, got %T", inputCells)
	}

	internalCells := make([]customattribute.CellJson, len(apiCells))
	for i, cell := range apiCells {

		var cellMap map[string]interface{}
		if cellMap, ok = cell.(map[string]interface{}); !ok {
			return types.StringNull(), fmt.Errorf("cell at index %d is not a valid object, got %T", i, cell)
		}
		apiCell := customattribute.CellJsonAPI{}
		apiCell.SetFromMap(cellMap)

		internalCells[i] = *apiCell.ToInternalModel()
	}

	marshaledCells, err := json.Marshal(internalCells)
	if err != nil {
		return types.StringNull(), err
	}
	return types.StringValue(string(marshaledCells)), nil
}

//
// Object converters
//

func convertToParamsGroupsObject(value interface{}) (types.Object, error) {
	obj, ok := value.(map[string]interface{})
	if !ok {
		return types.ObjectNull(converters.ParamsGroupsAttrTypes), fmt.Errorf("cannot convert %T to object", value)
	}

	// Iterate over the keys defined in ParamsGroupsAttrTypes to ensure consistency
	attrs := make(map[string]attr.Value, len(converters.ParamsGroupsAttrTypes))
	for key := range converters.ParamsGroupsAttrTypes {
		raw, exists := obj[key]
		if !exists {
			return types.ObjectNull(converters.ParamsGroupsAttrTypes), fmt.Errorf("params_groups is missing required key \"%q\"", key)
		}

		listVal, err := commonattribute.ConvertToStringList(raw)
		if err != nil {
			return types.ObjectNull(converters.ParamsGroupsAttrTypes), fmt.Errorf("params_groups.%s: %w", key, err)
		}
		attrs[key] = listVal
	}

	val, diags := types.ObjectValue(converters.ParamsGroupsAttrTypes, attrs)
	if diags.HasError() {
		return types.ObjectNull(converters.ParamsGroupsAttrTypes), fmt.Errorf("failed to create params_groups object: %s", diags.Errors()[0].Summary())
	}
	return val, nil
}
