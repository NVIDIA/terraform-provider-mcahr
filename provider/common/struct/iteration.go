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

import "reflect"

// IterateStruct iterates over all fields of a struct, invoking the provided callback function for each field.
//
// Parameters:
//   - structData: The struct to iterate over. Can be a struct value or a pointer to a struct.
//   - callback: Function called for each field with:
//   - field: Pointer to the reflect.StructField containing field metadata (name, type, tags, etc.)
//   - value: Pointer to the reflect.Value containing the field's actual value
//     Returns shouldBreak (bool): Return true to stop iteration early, false to continue
//
// The function automatically handles pointer dereferencing. If structData is a pointer to a struct,
// it will dereference it before iteration.
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//	p := Person{Name: "John", Age: 30}
//	IterateStruct(p, func(field *reflect.StructField, value *reflect.Value) bool {
//	    fmt.Printf("%s: %v\n", field.Name, value.Interface())
//	    return false // continue iteration
//	})
func IterateStruct(structData any, callback func(field *reflect.StructField, value *reflect.Value) (shouldBreak bool)) {

	t := reflect.TypeOf(structData)
	v := reflect.ValueOf(structData)

	// If the struct is a pointer, get the underlying struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		shouldBreak := callback(&field, &value)
		if shouldBreak {
			return
		}
	}
}

// IterateStructWithTags iterates over all fields of a struct with convenient access to field tags.
//
// This is a convenience wrapper around IterateStruct that provides direct access to commonly needed
// field properties: type, name, tags, and value.
//
// Parameters:
//   - structData: The struct to iterate over. Can be a struct value or a pointer to a struct.
//   - callback: Function called for each field with:
//   - fieldType: The reflect.Type of the field
//   - fieldName: The name of the field as a string
//   - fieldTags: The reflect.StructTag containing all struct tags (json, tfsdk, xml, etc.)
//   - fieldValue: Pointer to the reflect.Value containing the field's actual value
//     Returns shouldBreak (bool): Return true to stop iteration early, false to continue
//
// This function is particularly useful when you need to work with struct tags, as it provides
// direct access to the tags without needing to extract them from the StructField.
//
// Example:
//
//	type User struct {
//	    Name  string `json:"name" tfsdk:"name"`
//	    Email string `json:"email,omitempty" tfsdk:"email"`
//	}
//	u := User{Name: "Alice", Email: "alice@example.com"}
//	IterateStructWithTags(u, func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) bool {
//	    tfsdkTag := fieldTags.Get("tfsdk")
//	    jsonTag := fieldTags.Get("json")
//	    fmt.Printf("%s: tfsdk=%s, json=%s, value=%v\n", fieldName, tfsdkTag, jsonTag, fieldValue.Interface())
//	    return false // continue iteration
//	})
func IterateStructWithTags(structData any, callback func(fieldType reflect.Type, fieldName string, fieldTags reflect.StructTag, fieldValue *reflect.Value) (shouldBreak bool)) {
	IterateStruct(
		structData,
		func(field *reflect.StructField, value *reflect.Value) (shouldBreak bool) {
			return callback(field.Type, field.Name, field.Tag, value)
		},
	)
}
