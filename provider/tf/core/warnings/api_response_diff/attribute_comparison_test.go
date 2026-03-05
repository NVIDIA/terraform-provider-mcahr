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

package apiresponsediff

import (
	"reflect"
	"strings"
	coremodel "terraform/terraform-provider/provider/tf/core/schema"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CompleteTestModel struct {
	StringField  types.String  `tfsdk:"string_field"`
	BoolField    types.Bool    `tfsdk:"bool_field"`
	Int64Field   types.Int64   `tfsdk:"int64_field"`
	Float64Field types.Float64 `tfsdk:"float64_field"`
	ListField    types.List    `tfsdk:"list_field"`
	SetField     types.Set     `tfsdk:"set_field"`
	MapField     types.Map     `tfsdk:"map_field"`
	ObjectField  types.Object  `tfsdk:"object_field"`
	SkipField    types.String  `tfsdk:"-"`
	NoTagField   types.String
}

func (c CompleteTestModel) GetName() string {
	if !c.StringField.IsNull() && !c.StringField.IsUnknown() {
		return c.StringField.ValueString()
	}
	return ""
}

// Test models for testing computed field behavior
type runbookTestModel struct {
	Name        types.String `tfsdk:"name"`
	Cells       types.String `tfsdk:"cells"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func (r runbookTestModel) GetName() string {
	if !r.Name.IsNull() && !r.Name.IsUnknown() {
		return r.Name.ValueString()
	}
	return ""
}

type modelWithComputedFields struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	// Simulate computed-only fields (they would be null in config)
	CellsFull          types.String `tfsdk:"cells_full"`
	ParamsFull         types.String `tfsdk:"params_full"`
	ExternalParamsFull types.String `tfsdk:"external_params_full"`
}

func (m modelWithComputedFields) GetName() string {
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		return m.Name.ValueString()
	}
	return ""
}

func TestCompareModels_IdenticalModels(t *testing.T) {
	model1 := CompleteTestModel{
		StringField: types.StringValue("test"),
		BoolField:   types.BoolValue(true),
		Int64Field:  types.Int64Value(42),
	}

	model2 := CompleteTestModel{
		StringField: types.StringValue("test"),
		BoolField:   types.BoolValue(true),
		Int64Field:  types.Int64Value(42),
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	assert.Empty(t, differences)
}

func TestCompareModels_DifferentStringValues(t *testing.T) {
	model1 := CompleteTestModel{
		StringField: types.StringValue("original"),
	}

	model2 := CompleteTestModel{
		StringField: types.StringValue("modified"),
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	require.Len(t, differences, 1)
	assert.Equal(t, "string_field", differences[0].FieldName)
	assert.Equal(t, "original", differences[0].PlanValue)
	assert.Equal(t, "modified", differences[0].ResponseValue)
}

func TestCompareModels_NilPointers(t *testing.T) {
	var model1 *CompleteTestModel
	var model2 *CompleteTestModel

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	assert.Empty(t, differences)
}

func TestCompareModels_OneNilPointer(t *testing.T) {
	model1 := &CompleteTestModel{
		StringField: types.StringValue("test"),
	}
	var model2 *CompleteTestModel

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	assert.Empty(t, differences)
}

func TestCompareModels_PointerModels(t *testing.T) {
	model1 := &CompleteTestModel{
		StringField: types.StringValue("test1"),
	}

	model2 := &CompleteTestModel{
		StringField: types.StringValue("test2"),
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	require.Len(t, differences, 1)
	assert.Equal(t, "string_field", differences[0].FieldName)
}

func TestCompareModels_SkipFieldsWithoutTfsdkTag(t *testing.T) {
	model1 := CompleteTestModel{
		NoTagField: types.StringValue("value1"),
	}

	model2 := CompleteTestModel{
		NoTagField: types.StringValue("value2"),
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	assert.Empty(t, differences)
}

func TestCompareModels_SkipFieldsWithDashTag(t *testing.T) {
	model1 := CompleteTestModel{
		SkipField: types.StringValue("value1"),
	}

	model2 := CompleteTestModel{
		SkipField: types.StringValue("value2"),
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	assert.Empty(t, differences)
}

func TestCompareField_BothEmpty(t *testing.T) {
	configField := reflect.ValueOf(types.StringNull())
	responseField := reflect.ValueOf(types.StringNull())

	diff := compareField("test_field", configField, responseField)

	assert.Nil(t, diff)
}

func TestCompareField_DifferentValues(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("config"))
	responseField := reflect.ValueOf(types.StringValue("response"))

	diff := compareField("test_field", configField, responseField)

	require.NotNil(t, diff)
	assert.Equal(t, "test_field", diff.FieldName)
	assert.Equal(t, "config", diff.PlanValue)
	assert.Equal(t, "response", diff.ResponseValue)
}

func TestCompareField_SameValues(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("same"))
	responseField := reflect.ValueOf(types.StringValue("same"))

	diff := compareField("test_field", configField, responseField)

	assert.Nil(t, diff)
}

func TestExtractFieldValue_StringValue(t *testing.T) {
	field := reflect.ValueOf(types.StringValue("test value"))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "test value", result)
	assert.True(t, wasSet)
}

func TestExtractFieldValue_StringNull(t *testing.T) {
	field := reflect.ValueOf(types.StringNull())
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_StringUnknown(t *testing.T) {
	field := reflect.ValueOf(types.StringUnknown())
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_BoolValue(t *testing.T) {
	field := reflect.ValueOf(types.BoolValue(true))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "true", result)
	assert.True(t, wasSet)
}

func TestExtractFieldValue_BoolFalse(t *testing.T) {
	field := reflect.ValueOf(types.BoolValue(false))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "false", result)
	assert.True(t, wasSet)
}

func TestExtractFieldValue_BoolNull(t *testing.T) {
	field := reflect.ValueOf(types.BoolNull())
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_BoolUnknown(t *testing.T) {
	field := reflect.ValueOf(types.BoolUnknown())
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_Int64Value(t *testing.T) {
	field := reflect.ValueOf(types.Int64Value(42))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "42", result)
	assert.True(t, wasSet)
}

func TestExtractFieldValue_Int64Null(t *testing.T) {
	field := reflect.ValueOf(types.Int64Null())
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_Int64Unknown(t *testing.T) {
	field := reflect.ValueOf(types.Int64Unknown())
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_Float64Value(t *testing.T) {
	field := reflect.ValueOf(types.Float64Value(3.14))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "3.140000", result)
	assert.True(t, wasSet)
}

func TestExtractFieldValue_Float64Null(t *testing.T) {
	field := reflect.ValueOf(types.Float64Null())
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_Float64Unknown(t *testing.T) {
	field := reflect.ValueOf(types.Float64Unknown())
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_ListValue(t *testing.T) {
	elements := []attr.Value{
		types.StringValue("item1"),
		types.StringValue("item2"),
	}
	listValue, _ := types.ListValue(types.StringType, elements)
	field := reflect.ValueOf(listValue)

	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Contains(t, result, "item1")
	assert.Contains(t, result, "item2")
	assert.True(t, wasSet)
}

func TestExtractFieldValue_ListNull(t *testing.T) {
	field := reflect.ValueOf(types.ListNull(types.StringType))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_ListUnknown(t *testing.T) {
	field := reflect.ValueOf(types.ListUnknown(types.StringType))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_SetValue(t *testing.T) {
	elements := []attr.Value{
		types.StringValue("item1"),
		types.StringValue("item2"),
	}
	setValue, _ := types.SetValue(types.StringType, elements)
	field := reflect.ValueOf(setValue)

	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Contains(t, result, "item1")
	assert.Contains(t, result, "item2")
	assert.True(t, wasSet)
}

func TestExtractFieldValue_SetNull(t *testing.T) {
	field := reflect.ValueOf(types.SetNull(types.StringType))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_SetUnknown(t *testing.T) {
	field := reflect.ValueOf(types.SetUnknown(types.StringType))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_MapValue(t *testing.T) {
	elements := map[string]attr.Value{
		"key1": types.StringValue("value1"),
		"key2": types.StringValue("value2"),
	}
	mapValue, _ := types.MapValue(types.StringType, elements)
	field := reflect.ValueOf(mapValue)

	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Contains(t, result, "key1")
	assert.Contains(t, result, "value1")
	assert.True(t, wasSet)
}

func TestExtractFieldValue_MapNull(t *testing.T) {
	field := reflect.ValueOf(types.MapNull(types.StringType))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_MapUnknown(t *testing.T) {
	field := reflect.ValueOf(types.MapUnknown(types.StringType))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_ObjectValue(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
	}
	attrValues := map[string]attr.Value{
		"name": types.StringValue("test"),
	}
	objValue, _ := types.ObjectValue(attrTypes, attrValues)
	field := reflect.ValueOf(objValue)

	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Contains(t, result, "name")
	assert.True(t, wasSet)
}

func TestExtractFieldValue_ObjectNull(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
	}
	field := reflect.ValueOf(types.ObjectNull(attrTypes))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_ObjectUnknown(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
	}
	field := reflect.ValueOf(types.ObjectUnknown(attrTypes))
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_InvalidField(t *testing.T) {
	field := reflect.Value{}
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestExtractFieldValue_UnsupportedType(t *testing.T) {
	field := reflect.ValueOf("plain string")
	result, wasSet := extractFieldValueWithSetFlag(field)

	assert.Equal(t, "", result)
	assert.False(t, wasSet)
}

func TestCompareModels_MultipleDifferences(t *testing.T) {
	model1 := CompleteTestModel{
		StringField: types.StringValue("string1"),
		BoolField:   types.BoolValue(true),
		Int64Field:  types.Int64Value(10),
	}

	model2 := CompleteTestModel{
		StringField: types.StringValue("string2"),
		BoolField:   types.BoolValue(false),
		Int64Field:  types.Int64Value(20),
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	require.Len(t, differences, 3)

	fieldNames := make(map[string]bool)
	for _, diff := range differences {
		fieldNames[diff.FieldName] = true
	}

	assert.True(t, fieldNames["string_field"])
	assert.True(t, fieldNames["bool_field"])
	assert.True(t, fieldNames["int64_field"])
}

func TestCompareModels_MixedNullAndValues(t *testing.T) {
	model1 := CompleteTestModel{
		StringField: types.StringNull(),    // Not set by user - will be skipped
		BoolField:   types.BoolValue(true), // Set by user
	}

	model2 := CompleteTestModel{
		StringField: types.StringValue("value"), // Set by user (different from null)
		BoolField:   types.BoolNull(),           // Not set by user - will be skipped
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	// Only BoolField should be detected as different (model1 has it set, model2 doesn't)
	// StringField won't show as different because model1 didn't have it set
	require.Len(t, differences, 1)
	assert.Equal(t, "bool_field", differences[0].FieldName)
}

func TestCompareModels_ComputedFieldsAreSkipped(t *testing.T) {
	// Test that computed-only fields (null in config) are automatically skipped
	model1 := modelWithComputedFields{
		Name:               types.StringValue("test"),
		Description:        types.StringValue("desc1"),
		CellsFull:          types.StringNull(), // Computed field - not set in config
		ParamsFull:         types.StringNull(), // Computed field - not set in config
		ExternalParamsFull: types.StringNull(), // Computed field - not set in config
	}

	model2 := modelWithComputedFields{
		Name:               types.StringValue("test"),
		Description:        types.StringValue("desc2"),            // Different - should be detected
		CellsFull:          types.StringValue(`[{"op":"echo"}]`),  // Different but null in config - skipped
		ParamsFull:         types.StringValue(`[{"name":"p1"}]`),  // Different but null in config - skipped
		ExternalParamsFull: types.StringValue(`[{"name":"ep1"}]`), // Different but null in config - skipped
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	// Only description should show as different, computed fields should be skipped
	require.Len(t, differences, 1)
	assert.Equal(t, "description", differences[0].FieldName)
}

func TestCompareModels_BothUserAndComputedFields(t *testing.T) {
	// Verify that user-configurable fields are detected when computed fields also differ
	model1 := runbookTestModel{
		Name:        types.StringValue("test"),
		Cells:       types.StringValue(`base64data1`),
		Description: types.StringValue("desc1"),
		Enabled:     types.BoolValue(true),
	}

	model2 := runbookTestModel{
		Name:        types.StringValue("test2"),
		Cells:       types.StringValue(`base64data2`),
		Description: types.StringValue("desc2"),
		Enabled:     types.BoolValue(false),
	}

	differences := compareModels(model1, model2, map[string]coremodel.FieldComparisonRule{})

	// Should detect all user-configurable field differences
	require.Len(t, differences, 4)

	fieldNames := make(map[string]bool)
	for _, diff := range differences {
		fieldNames[diff.FieldName] = true
	}

	assert.True(t, fieldNames["name"])
	assert.True(t, fieldNames["cells"])
	assert.True(t, fieldNames["description"])
	assert.True(t, fieldNames["enabled"])
}

// Test model with fields that have _full variants
type modelWithFullFields struct {
	Name               types.String `tfsdk:"name"`
	Cells              types.String `tfsdk:"cells"`
	CellsFull          types.String `tfsdk:"cells_full"`
	Params             types.String `tfsdk:"params"`
	ParamsFull         types.String `tfsdk:"params_full"`
	ExternalParams     types.String `tfsdk:"external_params"`
	ExternalParamsFull types.String `tfsdk:"external_params_full"`
	Description        types.String `tfsdk:"description"`
}

func (m modelWithFullFields) GetName() string {
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		return m.Name.ValueString()
	}
	return ""
}

func TestCompareModels_SkipsFieldsWithFullVariant(t *testing.T) {
	config := modelWithFullFields{
		Name:           types.StringValue("test"),
		Cells:          types.StringValue(`[{"op":"echo"}]`),  // Minimal user config
		Params:         types.StringValue(`[{"name":"p1"}]`),  // Minimal user config
		ExternalParams: types.StringValue(`[{"name":"ep1"}]`), // Minimal user config
		Description:    types.StringValue("user description"),
	}

	apiResponse := modelWithFullFields{
		Name:               types.StringValue("test"),
		Cells:              types.StringValue(`[{"op":"echo","enabled":true}]`),           // API added defaults
		CellsFull:          types.StringValue(`[{"op":"echo","enabled":true}]`),           // Full representation
		Params:             types.StringValue(`[{"name":"p1","required":false}]`),         // API added defaults
		ParamsFull:         types.StringValue(`[{"name":"p1","required":false}]`),         // Full representation
		ExternalParams:     types.StringValue(`[{"name":"ep1","source":"alertmanager"}]`), // API added defaults
		ExternalParamsFull: types.StringValue(`[{"name":"ep1","source":"alertmanager"}]`), // Full representation
		Description:        types.StringValue("api modified description"),                 // Actually modified by API
	}

	// Explicitly exclude fields with _full variants using comparison rules
	comparisonRules := map[string]coremodel.FieldComparisonRule{
		"cells":           {Behavior: coremodel.SkipComparison},
		"params":          {Behavior: coremodel.SkipComparison},
		"external_params": {Behavior: coremodel.SkipComparison},
	}

	differences := compareModels(config, apiResponse, comparisonRules)

	// Should only detect description difference
	// cells, params, external_params are explicitly excluded via comparison rules
	require.Len(t, differences, 1, "Should only detect fields without exclusion rules")
	assert.Equal(t, "description", differences[0].FieldName)
	assert.Equal(t, "user description", differences[0].PlanValue)
	assert.Equal(t, "api modified description", differences[0].ResponseValue)
}

func TestCompareFieldWithRule_DefaultComparison(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("config value"))
	responseField := reflect.ValueOf(types.StringValue("response value"))
	rules := map[string]coremodel.FieldComparisonRule{}

	diff := compareFieldWithRule("test_field", configField, responseField, rules)

	require.NotNil(t, diff)
	assert.Equal(t, "test_field", diff.FieldName)
	assert.Equal(t, "config value", diff.PlanValue)
	assert.Equal(t, "response value", diff.ResponseValue)
}

func TestCompareFieldWithRule_CompareNormallyBehavior(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("value1"))
	responseField := reflect.ValueOf(types.StringValue("value2"))
	rules := map[string]coremodel.FieldComparisonRule{
		"test_field": {Behavior: coremodel.CompareNormally},
	}

	diff := compareFieldWithRule("test_field", configField, responseField, rules)

	require.NotNil(t, diff)
	assert.Equal(t, "test_field", diff.FieldName)
}

func TestCompareFieldWithRule_SkipComparison(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("config value"))
	responseField := reflect.ValueOf(types.StringValue("different value"))
	rules := map[string]coremodel.FieldComparisonRule{
		"test_field": {Behavior: coremodel.SkipComparison},
	}

	diff := compareFieldWithRule("test_field", configField, responseField, rules)

	assert.Nil(t, diff, "Should return nil for skipped fields")
}

func TestCompareFieldWithRule_CustomComparison_Equal(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("  value  "))
	responseField := reflect.ValueOf(types.StringValue("value"))
	rules := map[string]coremodel.FieldComparisonRule{
		"test_field": {
			Behavior: coremodel.CustomComparison,
			CustomCompare: func(fieldName, configValue, apiValue string) bool {
				// Custom comparison: trim whitespace
				return strings.TrimSpace(configValue) == strings.TrimSpace(apiValue)
			},
		},
	}

	diff := compareFieldWithRule("test_field", configField, responseField, rules)

	assert.Nil(t, diff, "Custom comparison should consider trimmed values equal")
}

func TestCompareFieldWithRule_CustomComparison_Different(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("value1"))
	responseField := reflect.ValueOf(types.StringValue("value2"))
	rules := map[string]coremodel.FieldComparisonRule{
		"test_field": {
			Behavior: coremodel.CustomComparison,
			CustomCompare: func(fieldName, configValue, apiValue string) bool {
				return configValue == apiValue
			},
		},
	}

	diff := compareFieldWithRule("test_field", configField, responseField, rules)

	require.NotNil(t, diff)
	assert.Equal(t, "test_field", diff.FieldName)
	assert.Equal(t, "value1", diff.PlanValue)
	assert.Equal(t, "value2", diff.ResponseValue)
}

func TestCompareFieldWithRule_CustomComparison_NilFunction(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("value1"))
	responseField := reflect.ValueOf(types.StringValue("value2"))
	rules := map[string]coremodel.FieldComparisonRule{
		"test_field": {
			Behavior:      coremodel.CustomComparison,
			CustomCompare: nil, // Invalid: nil function
		},
	}

	// Should fallback to default comparison when custom function is nil
	diff := compareFieldWithRule("test_field", configField, responseField, rules)

	require.NotNil(t, diff, "Should fallback to default comparison")
	assert.Equal(t, "test_field", diff.FieldName)
}

func TestCompareFieldWithCustomLogic_ValuesEqual(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("test"))
	responseField := reflect.ValueOf(types.StringValue("test"))
	customCompare := func(fieldName, configValue, apiValue string) bool {
		return configValue == apiValue
	}

	diff := compareFieldWithCustomLogic("test_field", configField, responseField, customCompare)

	assert.Nil(t, diff, "Should return nil when custom function returns true")
}

func TestCompareFieldWithCustomLogic_ValuesDifferent(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("config"))
	responseField := reflect.ValueOf(types.StringValue("response"))
	customCompare := func(fieldName, configValue, apiValue string) bool {
		return configValue == apiValue
	}

	diff := compareFieldWithCustomLogic("test_field", configField, responseField, customCompare)

	require.NotNil(t, diff)
	assert.Equal(t, "test_field", diff.FieldName)
	assert.Equal(t, "config", diff.PlanValue)
	assert.Equal(t, "response", diff.ResponseValue)
}

func TestCompareFieldWithCustomLogic_ConfigNotSet(t *testing.T) {
	configField := reflect.ValueOf(types.StringNull()) // Not set by user
	responseField := reflect.ValueOf(types.StringValue("response"))
	customCompare := func(fieldName, configValue, apiValue string) bool {
		return false // Always different
	}

	diff := compareFieldWithCustomLogic("test_field", configField, responseField, customCompare)

	assert.Nil(t, diff, "Should return nil when config field was not set by user")
}

func TestCompareFieldWithCustomLogic_CaseInsensitiveComparison(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("Value"))
	responseField := reflect.ValueOf(types.StringValue("value"))
	customCompare := func(fieldName, configValue, apiValue string) bool {
		return strings.EqualFold(configValue, apiValue)
	}

	diff := compareFieldWithCustomLogic("test_field", configField, responseField, customCompare)

	assert.Nil(t, diff, "Case-insensitive comparison should consider values equal")
}

func TestCompareFieldWithCustomLogic_PartialMatch(t *testing.T) {
	configField := reflect.ValueOf(types.StringValue("test"))
	responseField := reflect.ValueOf(types.StringValue("test with extra data"))
	customCompare := func(fieldName, configValue, apiValue string) bool {
		// Check if API value contains config value
		return strings.Contains(apiValue, configValue)
	}

	diff := compareFieldWithCustomLogic("test_field", configField, responseField, customCompare)

	assert.Nil(t, diff, "Partial match comparison should work")
}

// testModelMultiField is a test model for multiple custom rules test
type testModelMultiField struct {
	Field1 types.String `tfsdk:"field1"`
	Field2 types.String `tfsdk:"field2"`
	Field3 types.String `tfsdk:"field3"`
}

func (m testModelMultiField) GetName() string {
	return "test"
}

func TestCompareModels_MultipleCustomRules(t *testing.T) {
	config := testModelMultiField{
		Field1: types.StringValue("  value1  "),
		Field2: types.StringValue("VALUE2"),
		Field3: types.StringValue("value3"),
	}

	apiResponse := testModelMultiField{
		Field1: types.StringValue("value1"),    // Different but trimmed equal
		Field2: types.StringValue("value2"),    // Different but case-insensitive equal
		Field3: types.StringValue("different"), // Actually different
	}

	rules := map[string]coremodel.FieldComparisonRule{
		"field1": {
			Behavior: coremodel.CustomComparison,
			CustomCompare: func(fieldName, configValue, apiValue string) bool {
				return strings.TrimSpace(configValue) == strings.TrimSpace(apiValue)
			},
		},
		"field2": {
			Behavior: coremodel.CustomComparison,
			CustomCompare: func(fieldName, configValue, apiValue string) bool {
				return strings.EqualFold(configValue, apiValue)
			},
		},
		// field3 uses default comparison
	}

	differences := compareModels(config, apiResponse, rules)

	// Only field3 should be detected as different
	require.Len(t, differences, 1)
	assert.Equal(t, "field3", differences[0].FieldName)
}
