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
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: SetPlanAttributeValueToNil is not unit-tested here because it requires
// the full Terraform plugin framework context. It is integration-tested through
// its usage in the compatibility modifier (see plan/modifiers/compatibility).

func TestSetAttrValueToNil_PrimitiveTypes(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		input          attr.Value
		expectedType   attr.Value
		expectedIsNull bool
	}{
		{
			name:           "String type",
			input:          types.StringValue("test"),
			expectedType:   types.StringNull(),
			expectedIsNull: true,
		},
		{
			name:           "Bool type",
			input:          types.BoolValue(true),
			expectedType:   types.BoolNull(),
			expectedIsNull: true,
		},
		{
			name:           "Int64 type",
			input:          types.Int64Value(42),
			expectedType:   types.Int64Null(),
			expectedIsNull: true,
		},
		{
			name:           "Int32 type",
			input:          types.Int32Value(32),
			expectedType:   types.Int32Null(),
			expectedIsNull: true,
		},
		{
			name:           "Float64 type",
			input:          types.Float64Value(3.14),
			expectedType:   types.Float64Null(),
			expectedIsNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := SetAttrValueToNil(ctx, tt.input)

			// then
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.True(t, result.IsNull())
			assert.Equal(t, reflect.TypeOf(tt.expectedType), reflect.TypeOf(result))
		})
	}
}

func TestSetAttrValueToNil_CollectionTypes(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		input        attr.Value
		expectedType func() attr.Value
	}{
		{
			name:  "Set type",
			input: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
			expectedType: func() attr.Value {
				return types.ListNull(types.StringType)
			},
		},
		{
			name:  "List type",
			input: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
			expectedType: func() attr.Value {
				return types.ListNull(types.StringType)
			},
		},
		{
			name:  "Map type",
			input: types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("value")}),
			expectedType: func() attr.Value {
				return types.MapNull(types.StringType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := SetAttrValueToNil(ctx, tt.input)

			// then
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.True(t, result.IsNull())

			expected := tt.expectedType()
			assert.Equal(t, reflect.TypeOf(expected), reflect.TypeOf(result))
		})
	}
}

func TestSetAttrValueToNil_ObjectType(t *testing.T) {
	ctx := context.Background()

	// given
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
		"age":  types.Int64Type,
	}
	attrValues := map[string]attr.Value{
		"name": types.StringValue("test"),
		"age":  types.Int64Value(30),
	}
	objValue := types.ObjectValueMust(attrTypes, attrValues)

	// when
	result, err := SetAttrValueToNil(ctx, objValue)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.IsNull())

	// Verify it's an object type with the same attribute types
	objResult, ok := result.(types.Object)
	require.True(t, ok)
	assert.Equal(t, attrTypes, objResult.AttributeTypes(ctx))
}

func TestSetAttrValueToNil_AlreadyNull(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		input attr.Value
	}{
		{
			name:  "Null string",
			input: types.StringNull(),
		},
		{
			name:  "Null int64",
			input: types.Int64Null(),
		},
		{
			name:  "Null bool",
			input: types.BoolNull(),
		},
		{
			name:  "Null set",
			input: types.ListNull(types.StringType),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := SetAttrValueToNil(ctx, tt.input)

			// then
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.True(t, result.IsNull())
		})
	}
}

func TestSetAttrValueToNil_UnsupportedType(t *testing.T) {
	ctx := context.Background()

	// given - create a custom type that's not supported
	type UnsupportedType struct {
		attr.Value
	}
	unsupported := UnsupportedType{}

	// when
	result, err := SetAttrValueToNil(ctx, unsupported)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported attribute type")
}

func TestSetFieldValueToNil_Success(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		initialValue attr.Value
	}{
		{
			name:         "String field",
			initialValue: types.StringValue("test"),
		},
		{
			name:         "Int64 field",
			initialValue: types.Int64Value(42),
		},
		{
			name:         "Bool field",
			initialValue: types.BoolValue(true),
		},
		{
			name:         "Set field",
			initialValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test")}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given - create a struct with the field to make it addressable
			type TestFieldStruct struct {
				Field attr.Value
			}
			testStruct := TestFieldStruct{Field: tt.initialValue}
			structValue := reflect.ValueOf(&testStruct).Elem()
			fieldValue := structValue.Field(0)

			// when
			err := SetFieldValueToNil(ctx, &fieldValue)

			// then
			require.NoError(t, err)
			resultAttr, ok := fieldValue.Interface().(attr.Value)
			require.True(t, ok)
			assert.True(t, resultAttr.IsNull())
		})
	}
}

func TestSetFieldValueToNil_NonAttrValue(t *testing.T) {
	ctx := context.Background()

	// given - a field that is not an attr.Value
	type TestStruct struct {
		Field string
	}
	testStruct := TestStruct{Field: "just a string"}
	structValue := reflect.ValueOf(&testStruct).Elem()
	fieldValue := structValue.Field(0)

	// when
	err := SetFieldValueToNil(ctx, &fieldValue)

	// then
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert field value to attr.Value")
}

func TestSetFieldValueToNil_UnsupportedAttrType(t *testing.T) {
	ctx := context.Background()

	// given - create a custom attr.Value type
	type CustomAttrValue struct {
		attr.Value
	}
	type TestStruct struct {
		Field attr.Value
	}
	testStruct := TestStruct{Field: CustomAttrValue{}}
	structValue := reflect.ValueOf(&testStruct).Elem()
	fieldValue := structValue.Field(0)

	// when
	err := SetFieldValueToNil(ctx, &fieldValue)

	// then
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported attribute type")
}

func TestSetAttrValueToNil_PreservesElementTypes(t *testing.T) {
	ctx := context.Background()

	t.Run("Set preserves element type", func(t *testing.T) {
		// given
		setVal := types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(1), types.Int64Value(2)})

		// when
		result, err := SetAttrValueToNil(ctx, setVal)

		// then
		require.NoError(t, err)
		require.True(t, result.IsNull())

		setResult, ok := result.(types.List)
		require.True(t, ok)
		assert.Equal(t, types.Int64Type, setResult.ElementType(ctx))
	})

	t.Run("List preserves element type", func(t *testing.T) {
		// given
		listVal := types.ListValueMust(types.BoolType, []attr.Value{types.BoolValue(true)})

		// when
		result, err := SetAttrValueToNil(ctx, listVal)

		// then
		require.NoError(t, err)
		require.True(t, result.IsNull())

		listResult, ok := result.(types.List)
		require.True(t, ok)
		assert.Equal(t, types.BoolType, listResult.ElementType(ctx))
	})

	t.Run("Map preserves element type", func(t *testing.T) {
		// given
		mapVal := types.MapValueMust(types.Float64Type, map[string]attr.Value{"key": types.Float64Value(3.14)})

		// when
		result, err := SetAttrValueToNil(ctx, mapVal)

		// then
		require.NoError(t, err)
		require.True(t, result.IsNull())

		mapResult, ok := result.(types.Map)
		require.True(t, ok)
		assert.Equal(t, types.Float64Type, mapResult.ElementType(ctx))
	})
}

func TestSetAttrValueToNil_NestedCollections(t *testing.T) {
	ctx := context.Background()

	t.Run("List of objects", func(t *testing.T) {
		// given
		objType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name": types.StringType,
			},
		}
		objVal := types.ObjectValueMust(objType.AttrTypes, map[string]attr.Value{
			"name": types.StringValue("test"),
		})
		listVal := types.ListValueMust(objType, []attr.Value{objVal})

		// when
		result, err := SetAttrValueToNil(ctx, listVal)

		// then
		require.NoError(t, err)
		require.True(t, result.IsNull())

		listResult, ok := result.(types.List)
		require.True(t, ok)
		assert.Equal(t, objType, listResult.ElementType(ctx))
	})

	t.Run("Map of sets", func(t *testing.T) {
		// given
		setType := types.ListType{ElemType: types.StringType}
		setVal := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("item")})
		mapVal := types.MapValueMust(setType, map[string]attr.Value{"key": setVal})

		// when
		result, err := SetAttrValueToNil(ctx, mapVal)

		// then
		require.NoError(t, err)
		require.True(t, result.IsNull())

		mapResult, ok := result.(types.Map)
		require.True(t, ok)
		assert.Equal(t, setType, mapResult.ElementType(ctx))
	})
}
