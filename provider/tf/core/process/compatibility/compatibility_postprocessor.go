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

package compatibility

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	commonstruct "terraform/terraform-provider/provider/common/struct"
	model "terraform/terraform-provider/provider/tf/core/model"
)

func PostProcess[T model.TFModel](requestContext *common.RequestContext, tfModel T, options map[string]attribute.CompatibilityOptions) error {
	compatibilityChecker := attribute.NewCompatibilityChecker(requestContext.BackendVersion, options)

	// Get the underlying struct from the interface
	tfModelValue := reflect.ValueOf(tfModel)
	if tfModelValue.Kind() == reflect.Ptr {
		tfModelValue = tfModelValue.Elem()
	}

	// Iterate over all fields of the struct
	return iterateAndNullifyIncompatibleFields(requestContext.Context, tfModelValue, compatibilityChecker)
}

// iterateAndNullifyIncompatibleFields iterates over struct fields and nullifies incompatible ones
func iterateAndNullifyIncompatibleFields(ctx context.Context, structValue reflect.Value, checker *attribute.CompatibilityChecker) error {
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Skip unexported fields or fields that can't be set
		if !fieldValue.CanSet() {
			continue
		}

		// Get the tfsdk tag to match against compatibility options
		tfsdk := field.Tag.Get("tfsdk")
		if tfsdk == "" || tfsdk == "-" {
			continue
		}

		// Extract the attribute name from the tfsdk tag
		attributeName := strings.Split(tfsdk, ",")[0]

		// Check if the attribute is compatible
		if !checker.IsAttributeCompatible(attributeName) {
			if err := commonstruct.SetFieldValueToNil(ctx, &fieldValue); err != nil {
				return fmt.Errorf("failed to set field %s to null in post process: %w", attributeName, err)
			}
		}
	}

	return nil
}
