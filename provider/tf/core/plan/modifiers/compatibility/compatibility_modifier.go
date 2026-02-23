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
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	commonstruct "terraform/terraform-provider/provider/common/struct"
	"terraform/terraform-provider/provider/common/version"
	model "terraform/terraform-provider/provider/tf/core/model"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// T should always be a pointer to a struct
func ApplyResourceCompatibilityModifiers[T model.TFModel](ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema coreschema.ResourceSchema, backendVersion *version.BackendVersion, model T) {

	if isDestroyOperation(req) {
		return
	}

	diags := req.Config.Get(ctx, model)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	ApplyCompatibilityModifiers(ctx, req, resp, schema, backendVersion, model)
}

func ApplyCompatibilityModifiers[T model.TFModel](ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, schema coreschema.ResourceSchema, backendVersion *version.BackendVersion, configValues T) {

	if isDestroyOperation(req) {
		return
	}

	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, schema.GetCompatibilityOptions())

	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)
}

// applyVersionValidationModifier applies the version validation modifier to the resource
func applyVersionValidationModifier[T model.TFModel](ctx context.Context, req *resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, compatibilityChecker *attribute.CompatibilityChecker, configValues T) {

	// Iterate over all Config values of the struct attributes and check if they are compatible
	commonstruct.IterateStructWithTags(
		configValues,
		func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) (shouldBreak bool) {
			tfName := fieldTags.Get("tfsdk")

			if !hasTfSdkTag(tfName) {
				return false // skip if no tfsdk tag
			}
			if compatibilityChecker.IsAttributeCompatible(tfName) {
				return false // skip if compatible
			}

			// If the attribute is not compatible then:

			configValue, ok := fieldValue.Interface().(attr.Value)
			if !ok || configValue == nil {
				return false // skip if not an attr.Value type or is nil
			}

			if common.IsAttrKnown(configValue) {
				// user is providing a value for an incompatible attribute, raise a validation error
				resp.Diagnostics.AddAttributeError(
					path.Root(tfName),
					"Unsupported attribute",
					constructCompatibilityErrorMessage(compatibilityChecker, tfName),
				)

			} else {
				// set to null if not provided by the user (to avoid the default value being applied)
				err := commonstruct.SetPlanAttributeValueToNil(ctx, tfName, resp)
				if err != nil {
					resp.Diagnostics.AddError("Error setting plan attribute value to null", err.Error())
				}
			}

			return false
		},
	)
}

// constructCompatibilityErrorMessage creates a detailed error message for incompatible attributes
func constructCompatibilityErrorMessage(checker *attribute.CompatibilityChecker, attributeName string) string {
	errorMessage := fmt.Sprintf("%s attribute is not supported by the current platform. Current version: %s.", attributeName, checker.BackendVersion.Version)

	compatibilityOptions, ok := checker.Options[attributeName]
	if !ok {
		return errorMessage
	}

	if compatibilityOptions.MinVersion != "" {
		errorMessage += fmt.Sprintf(" Minimum version required: %s.", compatibilityOptions.MinVersion)
	}

	if compatibilityOptions.MaxVersion != "" {
		errorMessage += fmt.Sprintf(" Maximum allowed version: %s.", compatibilityOptions.MaxVersion)
	}

	return errorMessage

}

func hasTfSdkTag(tfName string) bool {
	return tfName != "" && tfName != "-"
}

func isDestroyOperation(req *resource.ModifyPlanRequest) bool {
	return req.Plan.Raw.IsNull() || req.Config.Raw.IsNull()
}
