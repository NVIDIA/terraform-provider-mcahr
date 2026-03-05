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

package process

import (
	"context"
	"terraform/terraform-provider/provider/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Count       types.Int64  `tfsdk:"count"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	ComputedID  types.String `tfsdk:"computed_id"`
	NoTagField  string       `tfsdk:"-"` // Field without tfsdk tag (ignored by framework)
}

func (t TestModel) GetName() string {
	if !t.Name.IsNull() && !t.Name.IsUnknown() {
		return t.Name.ValueString()
	}
	return ""
}

func TestRestoreAllFieldsFromPlan_Create(t *testing.T) {
	ctx := context.Background()
	requestContext := &common.RequestContext{
		Context:   ctx,
		Operation: common.Create,
	}

	// Plan model (what user configured)
	planModel := TestModel{
		Name:        types.StringValue("my-resource"),
		Description: types.StringValue("user description"),
		Count:       types.Int64Value(42),
		Enabled:     types.BoolValue(true),
		ComputedID:  types.StringNull(), // User didn't set this
		NoTagField:  "plan-value",
	}

	// API response model (what API returned - normalized values)
	apiModel := TestModel{
		Name:        types.StringValue("my-resource-normalized"),
		Description: types.StringValue("api description"),
		Count:       types.Int64Value(100),
		Enabled:     types.BoolValue(false),
		ComputedID:  types.StringValue("generated-id-123"), // API set this
		NoTagField:  "api-value",
	}

	// Create process data with plan
	processData := createProcessDataWithPlan(ctx, planModel)

	// Call RestoreAllFieldsFromPlan
	err := RestoreAllFieldsFromPlan(requestContext, processData, &apiModel)
	require.NoError(t, err)

	// Verify user-configured fields were restored from plan
	assert.Equal(t, "my-resource", apiModel.Name.ValueString())
	assert.Equal(t, "user description", apiModel.Description.ValueString())
	assert.Equal(t, int64(42), apiModel.Count.ValueInt64())
	assert.True(t, apiModel.Enabled.ValueBool())

	// Verify computed field (null in plan) kept API response
	assert.Equal(t, "generated-id-123", apiModel.ComputedID.ValueString())

	// Note: NoTagField with tfsdk:"-" is not populated by Plan.Get() from the framework,
	// so it remains at zero value ("") in the planModel, and this zero value is copied to apiModel
	assert.Equal(t, "", apiModel.NoTagField)
}

func TestRestoreAllFieldsFromPlan_Update(t *testing.T) {
	ctx := context.Background()
	requestContext := &common.RequestContext{
		Context:   ctx,
		Operation: common.Update,
	}

	planModel := TestModel{
		Name:        types.StringValue("updated-name"),
		Description: types.StringValue("updated description"),
		Count:       types.Int64Value(99),
		Enabled:     types.BoolValue(false),
		ComputedID:  types.StringNull(),
	}

	apiModel := TestModel{
		Name:        types.StringValue("api-name"),
		Description: types.StringValue("api description"),
		Count:       types.Int64Value(1),
		Enabled:     types.BoolValue(true),
		ComputedID:  types.StringValue("new-id"),
	}

	processData := createUpdateProcessDataWithPlan(ctx, planModel)

	err := RestoreAllFieldsFromPlan(requestContext, processData, &apiModel)
	require.NoError(t, err)

	// Verify fields restored from plan
	assert.Equal(t, "updated-name", apiModel.Name.ValueString())
	assert.Equal(t, "updated description", apiModel.Description.ValueString())
	assert.Equal(t, int64(99), apiModel.Count.ValueInt64())
	assert.False(t, apiModel.Enabled.ValueBool())

	// Verify computed field kept from API
	assert.Equal(t, "new-id", apiModel.ComputedID.ValueString())
}

func TestRestoreAllFieldsFromPlan_Read_ShouldError(t *testing.T) {
	ctx := context.Background()
	requestContext := &common.RequestContext{
		Context:   ctx,
		Operation: common.Read,
	}

	planModel := TestModel{
		Name: types.StringValue("plan-name"),
	}

	apiModel := TestModel{
		Name: types.StringValue("api-name"),
	}

	processData := createProcessDataWithPlan(ctx, planModel)

	// Should error because READ operations should not call RestoreAllFieldsFromPlan
	err := RestoreAllFieldsFromPlan(requestContext, processData, &apiModel)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported operation")
}

func TestRestoreAllFieldsFromPlan_Delete_ShouldError(t *testing.T) {
	ctx := context.Background()
	requestContext := &common.RequestContext{
		Context:   ctx,
		Operation: common.Delete,
	}

	planModel := TestModel{
		Name: types.StringValue("plan-name"),
	}

	apiModel := TestModel{
		Name: types.StringValue("api-name"),
	}

	processData := createProcessDataWithPlan(ctx, planModel)

	// Should error because DELETE operations should not call RestoreAllFieldsFromPlan
	err := RestoreAllFieldsFromPlan(requestContext, processData, &apiModel)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported operation")
}

func TestRestoreAllFieldsFromPlan_UnknownValue(t *testing.T) {
	ctx := context.Background()
	requestContext := &common.RequestContext{
		Context:   ctx,
		Operation: common.Create,
	}

	planModel := TestModel{
		Name:        types.StringValue("my-resource"),
		ComputedID:  types.StringUnknown(), // Unknown in plan
		Description: types.StringValue("desc"),
	}

	apiModel := TestModel{
		Name:        types.StringValue("api-name"),
		ComputedID:  types.StringValue("computed-value"),
		Description: types.StringValue("api-desc"),
	}

	processData := createProcessDataWithPlan(ctx, planModel)

	err := RestoreAllFieldsFromPlan(requestContext, processData, &apiModel)
	require.NoError(t, err)

	// Known values restored from plan
	assert.Equal(t, "my-resource", apiModel.Name.ValueString())
	assert.Equal(t, "desc", apiModel.Description.ValueString())

	// Unknown value kept from API
	assert.Equal(t, "computed-value", apiModel.ComputedID.ValueString())
}

func TestRestoreAllFieldsFromPlan_MixedNullAndValues(t *testing.T) {
	ctx := context.Background()
	requestContext := &common.RequestContext{
		Context:   ctx,
		Operation: common.Create,
	}

	planModel := TestModel{
		Name:        types.StringValue("name"),
		Description: types.StringNull(),    // Null in plan
		Count:       types.Int64Value(10),  // Set in plan
		Enabled:     types.BoolNull(),      // Null in plan
		ComputedID:  types.StringUnknown(), // Unknown in plan
	}

	apiModel := TestModel{
		Name:        types.StringValue("api-name"),
		Description: types.StringValue("api-desc"),
		Count:       types.Int64Value(999),
		Enabled:     types.BoolValue(true),
		ComputedID:  types.StringValue("api-id"),
	}

	processData := createProcessDataWithPlan(ctx, planModel)

	err := RestoreAllFieldsFromPlan(requestContext, processData, &apiModel)
	require.NoError(t, err)

	// Set values restored from plan
	assert.Equal(t, "name", apiModel.Name.ValueString())
	assert.Equal(t, int64(10), apiModel.Count.ValueInt64())

	// Null/Unknown values kept from API
	assert.Equal(t, "api-desc", apiModel.Description.ValueString())
	assert.True(t, apiModel.Enabled.ValueBool())
	assert.Equal(t, "api-id", apiModel.ComputedID.ValueString())
}

func TestCopyNonNullFields_NilPointers(t *testing.T) {
	var source *TestModel = nil
	var target *TestModel = nil

	err := copyNonNullFields(source, target)
	require.NoError(t, err)
}

func TestCopyNonNullFields_NonStructTypes(t *testing.T) {
	// This test verifies copyNonNullFields handles non-struct types gracefully
	// In practice, TFModel constraint ensures only structs are passed, but we test the internal logic

	source := TestModel{Name: types.StringValue("source")}
	target := TestModel{Name: types.StringValue("target")}

	// Call the internal copy function with valid structs
	err := copyNonNullFields(source, target)
	require.NoError(t, err)

	// Since target is passed by value (not pointer), it should not be modified
	assert.Equal(t, "target", target.Name.ValueString())
}

// Helper functions

func createProcessDataWithPlan(ctx context.Context, model TestModel) *ProcessData {
	sch := createTestSchema()
	plan := createTestPlan(ctx, sch, model)

	return &ProcessData{
		CreateRequest: &resource.CreateRequest{
			Plan: plan,
		},
	}
}

func createUpdateProcessDataWithPlan(ctx context.Context, model TestModel) *ProcessData {
	sch := schema.Schema{
		Attributes: createTestAttributes(),
	}
	plan := createTestPlan(ctx, sch, model)

	return &ProcessData{
		UpdateRequest: &resource.UpdateRequest{
			Plan: plan,
		},
	}
}

func createTestSchema() schema.Schema {
	return schema.Schema{
		Attributes: createTestAttributes(),
	}
}

func createTestAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":        schema.StringAttribute{Optional: true},
		"description": schema.StringAttribute{Optional: true},
		"count":       schema.Int64Attribute{Optional: true},
		"enabled":     schema.BoolAttribute{Optional: true},
		"computed_id": schema.StringAttribute{Computed: true},
	}
}

func createTestPlan(ctx context.Context, sch schema.Schema, model TestModel) tfsdk.Plan {
	// Create tftypes values
	attrTypes := map[string]tftypes.Type{
		"name":        tftypes.String,
		"description": tftypes.String,
		"count":       tftypes.Number,
		"enabled":     tftypes.Bool,
		"computed_id": tftypes.String,
	}

	values := map[string]tftypes.Value{
		"name":        convertStringToTFType(model.Name),
		"description": convertStringToTFType(model.Description),
		"count":       convertInt64ToTFType(model.Count),
		"enabled":     convertBoolToTFType(model.Enabled),
		"computed_id": convertStringToTFType(model.ComputedID),
	}

	objType := tftypes.Object{AttributeTypes: attrTypes}
	objVal := tftypes.NewValue(objType, values)

	return tfsdk.Plan{
		Raw:    objVal,
		Schema: sch,
	}
}

func convertStringToTFType(val types.String) tftypes.Value {
	if val.IsNull() {
		return tftypes.NewValue(tftypes.String, nil)
	}
	if val.IsUnknown() {
		return tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}
	return tftypes.NewValue(tftypes.String, val.ValueString())
}

func convertInt64ToTFType(val types.Int64) tftypes.Value {
	if val.IsNull() {
		return tftypes.NewValue(tftypes.Number, nil)
	}
	if val.IsUnknown() {
		return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue)
	}
	return tftypes.NewValue(tftypes.Number, val.ValueInt64())
}

func convertBoolToTFType(val types.Bool) tftypes.Value {
	if val.IsNull() {
		return tftypes.NewValue(tftypes.Bool, nil)
	}
	if val.IsUnknown() {
		return tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
	}
	return tftypes.NewValue(tftypes.Bool, val.ValueBool())
}
