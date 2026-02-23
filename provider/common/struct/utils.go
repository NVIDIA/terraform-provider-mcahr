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

package commonstruct

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetPlanAttributeValueToNil(ctx context.Context, attributeName string, resp *resource.ModifyPlanResponse) error {
	// Get the current attribute value from the plan to determine its type
	var currentValue attr.Value
	diags := resp.Plan.GetAttribute(ctx, path.Root(attributeName), &currentValue)
	if diags.HasError() {
		return fmt.Errorf("failed to get attribute %s from plan: %v", attributeName, diags.Errors())
	}

	nullValue, err := SetAttrValueToNil(ctx, currentValue)
	if err != nil {
		return fmt.Errorf("failed to set attribute value to null: %w", err)
	}

	// Set the null value in the response plan
	diags = resp.Plan.SetAttribute(ctx, path.Root(attributeName), nullValue)
	if diags.HasError() {
		return fmt.Errorf("failed to set attribute %s to null: %v", attributeName, diags.Errors())
	}
	return nil
}

func SetAttrValueToNil(ctx context.Context, attrValue attr.Value) (attr.Value, error) {

	var nullValue attr.Value

	// Create a null value based on the current value's type
	switch v := attrValue.(type) {
	case types.String:
		nullValue = types.StringNull()
	case types.Bool:
		nullValue = types.BoolNull()
	case types.Int64:
		nullValue = types.Int64Null()
	case types.Int32:
		nullValue = types.Int32Null()
	case types.Float64:
		nullValue = types.Float64Null()
	case types.Set:
		// For sets, preserve the element type
		nullValue = types.SetNull(v.ElementType(ctx))
	case types.List:
		// For lists, preserve the element type
		nullValue = types.ListNull(v.ElementType(ctx))
	case types.Map:
		// For maps, preserve the element type
		nullValue = types.MapNull(v.ElementType(ctx))
	case types.Object:
		// For objects, preserve the attribute types
		nullValue = types.ObjectNull(v.AttributeTypes(ctx))
	default:
		return nil, fmt.Errorf("unsupported attribute type %T", attrValue)
	}
	return nullValue, nil
}

// setFieldToNil sets a Terraform framework type field to null based on its type
func SetFieldValueToNil(ctx context.Context, fieldValue *reflect.Value) error {
	attrValue, ok := fieldValue.Interface().(attr.Value)
	if !ok {
		return fmt.Errorf("failed to convert field value to attr.Value")
	}
	nullValue, err := SetAttrValueToNil(ctx, attrValue)
	if err != nil {
		return fmt.Errorf("failed to set attribute value to null: %w", err)
	}
	fieldValue.Set(reflect.ValueOf(nullValue))

	return nil
}
