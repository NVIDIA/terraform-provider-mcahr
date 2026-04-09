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
	"fmt"
	"reflect"
	"terraform/terraform-provider/provider/common"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApplyDataJSONValues applies values from data JSON to root fields
func ApplyDataJSONValues(ctx context.Context, tfModel *runbooktf.RunbookTFModel) error {

	// Parse data JSON
	dataMap, err := ParseDataJSONToMap(tfModel.Data)
	if err != nil || dataMap == nil {
		return err
	}

	err = OnEachStructField(ctx, tfModel,
		func(snakeCaseFieldName string, fieldValue *reflect.Value) error {

			if IsJSONSkipField(snakeCaseFieldName) || IsDeprecatedAliasTarget(snakeCaseFieldName) {
				return nil
			}

			if common.IsAttrKnown(fieldValue.Interface().(attr.Value)) {
				// Value is set in TF model, so do not apply from data JSON
				return nil
			}

			if err := setFieldFromDataJSON(snakeCaseFieldName, *fieldValue, dataMap); err != nil {
				return fmt.Errorf("failed to apply field %s from data JSON: %w", snakeCaseFieldName, err)
			}

			return nil
		})
	if err != nil {
		return fmt.Errorf("failed to process fields from data JSON: %w", err)
	}

	return nil

}

func DataJSONToTFModel(ctx context.Context, dataJSON types.String) (*runbooktf.RunbookTFModel, error) {

	tfModel := runbooktf.RunbookTFModel{}
	tfModel.Data = dataJSON

	err := ApplyDataJSONValues(ctx, &tfModel)
	if err != nil {
		return nil, fmt.Errorf("failed to apply data JSON to TF model: %w", err)
	}

	return &tfModel, nil
}

// setFieldFromDataJSON applies a single field value from data JSON to the model
func setFieldFromDataJSON(fieldName string, tfModelValue reflect.Value, dataMap map[string]interface{}) error {

	// Get the value from data JSON (using alias for migrated fields like cells_list → cells)
	dataFieldName := ResolveDataFieldName(fieldName)
	dataValueRaw := findValueInMap(dataFieldName, dataMap)
	if dataValueRaw == nil {
		// Do nothing if field is not present in data JSON
		return nil
	}

	// Convert the data JSON value to the appropriate Terraform type
	dataValue, err := convertDataValueToTerraformValue(tfModelValue.Type(), dataValueRaw, fieldName)
	if err != nil {
		return fmt.Errorf("failed to convert value for field %s: %w", fieldName, err)
	}

	tfModelValue.Set(reflect.ValueOf(dataValue))

	return nil
}
