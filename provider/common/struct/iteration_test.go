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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test structs for iteration
type SimpleStruct struct {
	Name    string
	Age     int
	Active  bool
	Score   float64
	Private string // unexported field
}

type TaggedStruct struct {
	Field1 string `json:"field1" xml:"f1"`
	Field2 int    `json:"field2,omitempty" validate:"required"`
	Field3 bool   `json:"-" custom:"skip"`
	Field4 string `json:"field4" min_version:"1.0.0"`
}

type EmbeddedStruct struct {
	Outer string
	SimpleStruct
	Additional int
}

func TestIterateStruct_BasicIteration(t *testing.T) {
	// given
	s := SimpleStruct{
		Name:    "John",
		Age:     30,
		Active:  true,
		Score:   95.5,
		Private: "private",
	}

	// when
	type FieldData struct {
		Name  string
		Value interface{}
	}
	var fields []FieldData
	IterateStruct(s, func(field *reflect.StructField, value *reflect.Value) bool {
		fields = append(fields, FieldData{
			Name:  field.Name,
			Value: value.Interface(),
		})
		return false
	})

	// then
	assert.Equal(t, 5, len(fields))

	// Verify each field name corresponds to its correct value
	fieldMap := make(map[string]interface{})
	for _, f := range fields {
		fieldMap[f.Name] = f.Value
	}

	assert.Equal(t, "John", fieldMap["Name"])
	assert.Equal(t, 30, fieldMap["Age"])
	assert.Equal(t, true, fieldMap["Active"])
	assert.Equal(t, 95.5, fieldMap["Score"])
	assert.Equal(t, "private", fieldMap["Private"])
}

func TestIterateStruct_WithPointer(t *testing.T) {
	// given
	s := &SimpleStruct{
		Name:   "Alice",
		Age:    25,
		Active: false,
	}

	// when
	var fieldNames []string
	IterateStruct(s, func(field *reflect.StructField, value *reflect.Value) bool {
		fieldNames = append(fieldNames, field.Name)
		return false
	})

	// then
	assert.Equal(t, 5, len(fieldNames))
	assert.ElementsMatch(t, []string{"Name", "Age", "Active", "Score", "Private"}, fieldNames)
}

func TestIterateStruct_EarlyBreak(t *testing.T) {
	// given
	s := SimpleStruct{
		Name:   "Bob",
		Age:    40,
		Active: true,
		Score:  88.0,
	}

	// when
	var processedFields []string
	IterateStruct(s, func(field *reflect.StructField, value *reflect.Value) bool {
		processedFields = append(processedFields, field.Name)
		// Break after processing 2 fields
		return len(processedFields) >= 2
	})

	// then
	assert.Equal(t, 2, len(processedFields))
}

func TestIterateStruct_EmptyStruct(t *testing.T) {
	// given
	type EmptyStruct struct{}
	s := EmptyStruct{}

	// when
	var count int
	IterateStruct(s, func(field *reflect.StructField, value *reflect.Value) bool {
		count++
		return false
	})

	// then
	assert.Equal(t, 0, count)
}

func TestIterateStruct_FieldModification(t *testing.T) {
	// given
	s := SimpleStruct{
		Name: "Original",
		Age:  20,
	}

	// when
	IterateStruct(&s, func(field *reflect.StructField, value *reflect.Value) bool {
		if field.Name == "Name" && value.CanSet() {
			value.SetString("Modified")
		}
		if field.Name == "Age" && value.CanSet() {
			value.SetInt(99)
		}
		return false
	})

	// then
	assert.Equal(t, "Modified", s.Name)
	assert.Equal(t, 99, s.Age)
}

func TestIterateStruct_ReadOnly(t *testing.T) {
	// given
	s := SimpleStruct{Name: "Test"}

	// when
	IterateStruct(s, func(field *reflect.StructField, value *reflect.Value) bool {
		// Try to set a value on non-pointer struct
		if field.Name == "Name" {
			assert.False(t, value.CanSet(), "Should not be able to set value on non-pointer struct")
		}
		return false
	})
}

func TestIterateStructWithTags_BasicTags(t *testing.T) {
	// given
	s := TaggedStruct{
		Field1: "value1",
		Field2: 42,
		Field3: true,
		Field4: "value4",
	}

	// when
	var results []struct {
		name    string
		jsonTag string
		xmlTag  string
	}

	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		results = append(results, struct {
			name    string
			jsonTag string
			xmlTag  string
		}{
			name:    fieldName,
			jsonTag: fieldTags.Get("json"),
			xmlTag:  fieldTags.Get("xml"),
		})
		return false
	})

	// then
	assert.Equal(t, 4, len(results))

	// Check Field1
	field1 := results[0]
	assert.Equal(t, "Field1", field1.name)
	assert.Equal(t, "field1", field1.jsonTag)
	assert.Equal(t, "f1", field1.xmlTag)

	// Check Field2
	field2 := results[1]
	assert.Equal(t, "Field2", field2.name)
	assert.Equal(t, "field2,omitempty", field2.jsonTag)

	// Check Field3
	field3 := results[2]
	assert.Equal(t, "Field3", field3.name)
	assert.Equal(t, "-", field3.jsonTag)
}

func TestIterateStructWithTags_FilterByTag(t *testing.T) {
	// given
	s := TaggedStruct{
		Field1: "value1",
		Field2: 42,
		Field3: true,
		Field4: "value4",
	}

	// when - filter fields that have validate tag
	var fieldsWithValidate []string
	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		if validateTag := fieldTags.Get("validate"); validateTag != "" {
			fieldsWithValidate = append(fieldsWithValidate, fieldName)
		}
		return false
	})

	// then
	assert.Equal(t, 1, len(fieldsWithValidate))
	assert.Equal(t, []string{"Field2"}, fieldsWithValidate)
}

func TestIterateStructWithTags_SkipDashTag(t *testing.T) {
	// given
	s := TaggedStruct{
		Field1: "value1",
		Field2: 42,
		Field3: true,
		Field4: "value4",
	}

	// when - skip fields with json:"-"
	var nonSkippedFields []string
	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		jsonTag := fieldTags.Get("json")
		if jsonTag != "-" && jsonTag != "" {
			nonSkippedFields = append(nonSkippedFields, fieldName)
		}
		return false
	})

	// then
	assert.Equal(t, 3, len(nonSkippedFields))
	assert.ElementsMatch(t, []string{"Field1", "Field2", "Field4"}, nonSkippedFields)
}

func TestIterateStructWithTags_EarlyBreak(t *testing.T) {
	// given
	s := TaggedStruct{
		Field1: "value1",
		Field2: 42,
		Field3: true,
		Field4: "value4",
	}

	// when - break after finding first field with omitempty
	var foundField string
	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		jsonTag := fieldTags.Get("json")
		if contains(jsonTag, "omitempty") {
			foundField = fieldName
			return true // break
		}
		return false
	})

	// then
	assert.Equal(t, "Field2", foundField)
}

func TestIterateStructWithTags_WithPointer(t *testing.T) {
	// given
	s := &TaggedStruct{
		Field1: "pointer_value",
		Field2: 100,
	}

	// when
	var fieldNames []string
	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		fieldNames = append(fieldNames, fieldName)
		return false
	})

	// then
	assert.Equal(t, 4, len(fieldNames))
}

func TestIterateStructWithTags_AccessFieldValues(t *testing.T) {
	// given
	s := TaggedStruct{
		Field1: "test_string",
		Field2: 777,
		Field3: false,
		Field4: "versioned",
	}

	// when
	type FieldData struct {
		Name  string
		Value interface{}
		Tag   string
	}
	var fieldData []FieldData

	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		fieldData = append(fieldData, FieldData{
			Name:  fieldName,
			Value: fieldValue.Interface(),
			Tag:   fieldTags.Get("json"),
		})
		return false
	})

	// then
	assert.Equal(t, 4, len(fieldData))

	// Build a map to verify each field name corresponds to its correct value and tag
	fieldMap := make(map[string]FieldData)
	for _, f := range fieldData {
		fieldMap[f.Name] = f
	}

	// Verify Field1 data
	field1 := fieldMap["Field1"]
	assert.Equal(t, "test_string", field1.Value)
	assert.Equal(t, "field1", field1.Tag)

	// Verify Field2 data
	field2 := fieldMap["Field2"]
	assert.Equal(t, 777, field2.Value)
	assert.Equal(t, "field2,omitempty", field2.Tag)

	// Verify Field3 data
	field3 := fieldMap["Field3"]
	assert.Equal(t, false, field3.Value)
	assert.Equal(t, "-", field3.Tag)

	// Verify Field4 data
	field4 := fieldMap["Field4"]
	assert.Equal(t, "versioned", field4.Value)
	assert.Equal(t, "field4", field4.Tag)
}

func TestIterateStructWithTags_EmptyStruct(t *testing.T) {
	// given
	type EmptyStruct struct{}
	s := EmptyStruct{}

	// when
	var count int
	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		count++
		return false
	})

	// then
	assert.Equal(t, 0, count)
}

func TestIterateStructWithTags_FieldTypes(t *testing.T) {
	// given
	type TypedStruct struct {
		StringField  string         `json:"string_field"`
		IntField     int            `json:"int_field"`
		BoolField    bool           `json:"bool_field"`
		SliceField   []string       `json:"slice_field"`
		MapField     map[string]int `json:"map_field"`
		PointerField *string        `json:"pointer_field"`
	}

	s := TypedStruct{}

	// when
	var fieldTypes []reflect.Kind

	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		fieldTypes = append(fieldTypes, fieldType.Kind())
		return false
	})

	// then
	assert.Equal(t, 6, len(fieldTypes))
	assert.ElementsMatch(t, []reflect.Kind{reflect.String, reflect.Int, reflect.Bool, reflect.Slice, reflect.Map, reflect.Ptr}, fieldTypes)
}

func TestIterateStructWithTags_ModifyValues(t *testing.T) {
	// given
	s := &TaggedStruct{
		Field1: "original",
		Field2: 10,
	}

	// when - modify values through iteration
	IterateStructWithTags(s, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
		if fieldName == "Field1" && fieldValue.CanSet() {
			fieldValue.SetString("modified")
		}
		if fieldName == "Field2" && fieldValue.CanSet() {
			fieldValue.SetInt(999)
		}
		return false
	})

	// then
	assert.Equal(t, "modified", s.Field1)
	assert.Equal(t, 999, s.Field2)
}

func TestIterateStruct_NestedStructs(t *testing.T) {
	// given
	s := EmbeddedStruct{
		Outer:      "outer_value",
		Additional: 42,
		SimpleStruct: SimpleStruct{
			Name: "embedded",
			Age:  25,
		},
	}

	// when
	var fieldNames []string
	IterateStruct(s, func(field *reflect.StructField, value *reflect.Value) bool {
		fieldNames = append(fieldNames, field.Name)
		return false
	})

	// then
	// Should iterate over the fields of EmbeddedStruct, including the embedded SimpleStruct
	assert.Equal(t, 3, len(fieldNames))
	assert.ElementsMatch(t, []string{"Outer", "SimpleStruct", "Additional"}, fieldNames)
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(s)] != "" && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
