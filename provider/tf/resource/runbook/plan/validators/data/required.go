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
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

func validateRequiredFields(originalTFModel *model.RunbookTFModel, dataMap map[string]any) error {

	if err := validateReqField("name", "name", originalTFModel.Name, dataMap); err != nil {
		return err
	}
	// Any other required fields should be added here

	return nil
}

func validateReqField(tfFieldName string, dataFieldName string, modelValue attr.Value, dataMap map[string]any) error {

	if isFieldSet(modelValue, dataFieldName, dataMap) {
		return nil
	}
	return fmt.Errorf("The argument \"%s\" is required, but no definition was found.", tfFieldName)
}

func isFieldSet(modelValue attr.Value, dataFieldName string, dataMap map[string]any) bool {
	modelAttrHasValue := common.IsAttrKnown(modelValue)
	_, dataAttrHasValue := dataMap[dataFieldName]
	return modelAttrHasValue || dataAttrHasValue
}
