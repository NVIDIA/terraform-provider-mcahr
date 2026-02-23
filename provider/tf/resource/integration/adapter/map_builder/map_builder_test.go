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

package mapbuilder

import (
	"testing"

	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/common/version"

	"github.com/stretchr/testify/assert"
)

func TestNewMapBuilder(t *testing.T) {
	tests := []struct {
		name                 string
		backendVersion       *version.BackendVersion
		compatibilityOptions map[string]attribute.CompatibilityOptions
	}{
		{
			name:           "with valid version and options",
			backendVersion: version.NewBackendVersion("release-1.2.3"),
			compatibilityOptions: map[string]attribute.CompatibilityOptions{
				"field1": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
			},
		},
		{
			name:                 "with nil version",
			backendVersion:       nil,
			compatibilityOptions: map[string]attribute.CompatibilityOptions{},
		},
		{
			name:                 "with nil options",
			backendVersion:       version.NewBackendVersion("release-1.0.0"),
			compatibilityOptions: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewMapBuilder(tt.backendVersion, tt.compatibilityOptions)

			assert.NotNil(t, builder)
			assert.NotNil(t, builder.attrCompatibilityChecker)
			assert.NotNil(t, builder.fields)
			assert.Equal(t, 0, len(builder.fields))
		})
	}
}

func TestMapBuilder_SetField_CompatibleField(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.5.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{
		"field1": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
	}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)
	result := builder.SetField("mapField1", "field1", "testValue")

	// Should return the builder for chaining
	assert.Equal(t, builder, result)

	// Field should be set
	resultMap := builder.Build()
	assert.Equal(t, 1, len(resultMap))
	assert.Equal(t, "testValue", resultMap["mapField1"])
}

func TestMapBuilder_SetField_IncompatibleField_BelowMinVersion(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{
		"field1": {MinVersion: "release-1.5.0", MaxVersion: "release-2.0.0"},
	}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)
	result := builder.SetField("mapField1", "field1", "testValue")

	// Should return the builder for chaining
	assert.Equal(t, builder, result)

	// Field should NOT be set (backend version too old)
	resultMap := builder.Build()
	assert.Equal(t, 0, len(resultMap))
}

func TestMapBuilder_SetField_IncompatibleField_AboveMaxVersion(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-2.5.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{
		"field1": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
	}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)
	result := builder.SetField("mapField1", "field1", "testValue")

	// Should return the builder for chaining
	assert.Equal(t, builder, result)

	// Field should NOT be set (backend version too new, field deprecated)
	resultMap := builder.Build()
	assert.Equal(t, 0, len(resultMap))
}

func TestMapBuilder_SetField_NoCompatibilityOptions(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.5.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)
	builder.SetField("mapField1", "field1", "testValue")

	// Field should be set (no compatibility restrictions)
	resultMap := builder.Build()
	assert.Equal(t, 1, len(resultMap))
	assert.Equal(t, "testValue", resultMap["mapField1"])
}

func TestMapBuilder_SetField_NilBackendVersion(t *testing.T) {
	compatibilityOptions := map[string]attribute.CompatibilityOptions{
		"field1": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
	}

	builder := NewMapBuilder(nil, compatibilityOptions)
	builder.SetField("mapField1", "field1", "testValue")

	// Field should be set (nil version means assume latest)
	resultMap := builder.Build()
	assert.Equal(t, 1, len(resultMap))
	assert.Equal(t, "testValue", resultMap["mapField1"])
}

func TestMapBuilder_SetField_MultipleFields(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.5.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{
		"field1": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		"field2": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		"field3": {MinVersion: "release-2.5.0", MaxVersion: "release-3.0.0"}, // Incompatible (too old)
	}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)
	builder.SetField("mapField1", "field1", "value1")
	builder.SetField("mapField2", "field2", 42)
	builder.SetField("mapField3", "field3", "value3") // Should not be set

	resultMap := builder.Build()
	assert.Equal(t, 2, len(resultMap))
	assert.Equal(t, "value1", resultMap["mapField1"])
	assert.Equal(t, 42, resultMap["mapField2"])
	assert.NotContains(t, resultMap, "mapField3")
}

func TestMapBuilder_SetField_MethodChaining(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.5.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{
		"field1": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		"field2": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
	}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)

	// Test method chaining
	resultMap := builder.
		SetField("mapField1", "field1", "value1").
		SetField("mapField2", "field2", 42).
		Build()

	assert.Equal(t, 2, len(resultMap))
	assert.Equal(t, "value1", resultMap["mapField1"])
	assert.Equal(t, 42, resultMap["mapField2"])
}

func TestMapBuilder_SetField_DifferentTypes(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.5.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)

	// Test different value types
	stringVal := "string value"
	intVal := 123
	boolVal := true
	floatVal := 45.67
	sliceVal := []string{"a", "b", "c"}
	mapVal := map[string]interface{}{"nested": "value"}

	builder.
		SetField("stringField", "field1", stringVal).
		SetField("intField", "field2", intVal).
		SetField("boolField", "field3", boolVal).
		SetField("floatField", "field4", floatVal).
		SetField("sliceField", "field5", sliceVal).
		SetField("mapField", "field6", mapVal)

	resultMap := builder.Build()
	assert.Equal(t, 6, len(resultMap))
	assert.Equal(t, stringVal, resultMap["stringField"])
	assert.Equal(t, intVal, resultMap["intField"])
	assert.Equal(t, boolVal, resultMap["boolField"])
	assert.Equal(t, floatVal, resultMap["floatField"])
	assert.Equal(t, sliceVal, resultMap["sliceField"])
	assert.Equal(t, mapVal, resultMap["mapField"])
}

func TestMapBuilder_Build_EmptyMap(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)
	resultMap := builder.Build()

	assert.NotNil(t, resultMap)
	assert.Equal(t, 0, len(resultMap))
}

func TestMapBuilder_Build_MultipleCalls(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)
	builder.SetField("field1", "tf_field1", "value1")

	// Build should return the same map each time
	result1 := builder.Build()
	result2 := builder.Build()

	assert.Equal(t, result1, result2)
	assert.Equal(t, "value1", result1["field1"])
	assert.Equal(t, "value1", result2["field1"])
}

func TestMapBuilder_SetField_OverwriteValue(t *testing.T) {
	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityOptions := map[string]attribute.CompatibilityOptions{}

	builder := NewMapBuilder(backendVersion, compatibilityOptions)
	builder.SetField("field1", "tf_field1", "value1")
	builder.SetField("field1", "tf_field1", "value2") // Overwrite

	resultMap := builder.Build()
	assert.Equal(t, 1, len(resultMap))
	assert.Equal(t, "value2", resultMap["field1"])
}
