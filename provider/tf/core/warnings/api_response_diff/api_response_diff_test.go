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
	"context"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/tf/core/process"
	coremodel "terraform/terraform-provider/provider/tf/core/schema"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Count       types.Int64  `tfsdk:"count"`
}

func (t TestModel) GetName() string {
	if !t.Name.IsNull() && !t.Name.IsUnknown() {
		return t.Name.ValueString()
	}
	return ""
}

func (t TestModel) GetAttributeKeys() []string {
	return []string{"name", "description", "count"}
}

func (t TestModel) GetResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":        schema.StringAttribute{},
			"description": schema.StringAttribute{},
			"count":       schema.Int64Attribute{},
		},
	}
}

// mockSchema implements ResourceSchema for testing
type mockSchema struct {
	schema          schema.Schema
	comparisonRules map[string]coremodel.FieldComparisonRule
}

func (m *mockSchema) GetSchema() schema.Schema {
	return m.schema
}

func (m *mockSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{}
}

func (m *mockSchema) GetFieldComparisonRules() map[string]coremodel.FieldComparisonRule {
	if m.comparisonRules == nil {
		return map[string]coremodel.FieldComparisonRule{}
	}
	return m.comparisonRules
}

// createTestConfig creates a tfsdk.Config for testing purposes
func createTestConfig(sch schema.Schema) tfsdk.Config {
	attrTypes := make(map[string]tftypes.Type)
	values := make(map[string]tftypes.Value)

	for name, attr := range sch.Attributes {
		switch attr.(type) {
		case schema.StringAttribute:
			attrTypes[name] = tftypes.String
			values[name] = tftypes.NewValue(tftypes.String, nil)
		case schema.Int64Attribute:
			attrTypes[name] = tftypes.Number
			values[name] = tftypes.NewValue(tftypes.Number, nil)
		case schema.BoolAttribute:
			attrTypes[name] = tftypes.Bool
			values[name] = tftypes.NewValue(tftypes.Bool, nil)
		case schema.Float64Attribute:
			attrTypes[name] = tftypes.Number
			values[name] = tftypes.NewValue(tftypes.Number, nil)
		default:
			attrTypes[name] = tftypes.String
			values[name] = tftypes.NewValue(tftypes.String, nil)
		}
	}

	configData := tftypes.NewValue(tftypes.Object{
		AttributeTypes: attrTypes,
	}, values)

	return tfsdk.Config{
		Raw:    configData,
		Schema: sch,
	}
}

// createTestPlan creates a tfsdk.Plan for testing purposes
func createTestPlan(sch schema.Schema) tfsdk.Plan {
	attrTypes := make(map[string]tftypes.Type)
	values := make(map[string]tftypes.Value)

	for name, attr := range sch.Attributes {
		switch attr.(type) {
		case schema.StringAttribute:
			attrTypes[name] = tftypes.String
			values[name] = tftypes.NewValue(tftypes.String, nil)
		case schema.Int64Attribute:
			attrTypes[name] = tftypes.Number
			values[name] = tftypes.NewValue(tftypes.Number, nil)
		case schema.BoolAttribute:
			attrTypes[name] = tftypes.Bool
			values[name] = tftypes.NewValue(tftypes.Bool, nil)
		case schema.Float64Attribute:
			attrTypes[name] = tftypes.Number
			values[name] = tftypes.NewValue(tftypes.Number, nil)
		default:
			attrTypes[name] = tftypes.String
			values[name] = tftypes.NewValue(tftypes.String, nil)
		}
	}

	planData := tftypes.NewValue(tftypes.Object{
		AttributeTypes: attrTypes,
	}, values)

	return tfsdk.Plan{
		Raw:    planData,
		Schema: sch,
	}
}

// createTestConfigWithValues creates a tfsdk.Config with actual values for testing
func createTestConfigWithValues(model TestModel) tfsdk.Config {
	sch := model.GetResourceSchema()
	attrTypes := map[string]tftypes.Type{
		"name":        tftypes.String,
		"description": tftypes.String,
		"count":       tftypes.Number,
	}

	values := map[string]tftypes.Value{
		"name":        tftypes.NewValue(tftypes.String, model.Name.ValueStringPointer()),
		"description": tftypes.NewValue(tftypes.String, model.Description.ValueStringPointer()),
		"count":       tftypes.NewValue(tftypes.Number, model.Count.ValueInt64Pointer()),
	}

	configData := tftypes.NewValue(tftypes.Object{
		AttributeTypes: attrTypes,
	}, values)

	return tfsdk.Config{
		Raw:    configData,
		Schema: sch,
	}
}

// createTestPlanWithValues creates a tfsdk.Plan with actual values for testing
func createTestPlanWithValues(model TestModel) tfsdk.Plan {
	sch := model.GetResourceSchema()
	attrTypes := map[string]tftypes.Type{
		"name":        tftypes.String,
		"description": tftypes.String,
		"count":       tftypes.Number,
	}

	values := map[string]tftypes.Value{
		"name":        tftypes.NewValue(tftypes.String, model.Name.ValueStringPointer()),
		"description": tftypes.NewValue(tftypes.String, model.Description.ValueStringPointer()),
		"count":       tftypes.NewValue(tftypes.Number, model.Count.ValueInt64Pointer()),
	}

	planData := tftypes.NewValue(tftypes.Object{
		AttributeTypes: attrTypes,
	}, values)

	return tfsdk.Plan{
		Raw:    planData,
		Schema: sch,
	}
}

func TestCheckPlanVsApiResponseDelta_NoDifferences(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Create,
	}

	planModel := TestModel{
		Name:        types.StringValue("test-name"),
		Description: types.StringValue("test-description"),
		Count:       types.Int64Value(42),
	}

	apiResponseModel := &TestModel{
		Name:        types.StringValue("test-name"),
		Description: types.StringValue("test-description"),
		Count:       types.Int64Value(42),
	}

	plan := createTestPlanWithValues(planModel)

	processData := &process.ProcessData{
		CreateRequest: &resource.CreateRequest{
			Plan: plan,
		},
		CreateResponse: &resource.CreateResponse{},
	}

	mockSch := &mockSchema{schema: schema.Schema{}}
	err := CheckPlanVsApiResponseDelta(requestContext, processData, mockSch, apiResponseModel)

	assert.NoError(t, err)
	assert.Empty(t, processData.CreateResponse.Diagnostics)
}

func TestCheckPlanVsApiResponseDelta_WithDifferences(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Create,
	}

	apiResponseModel := &TestModel{
		Name:        types.StringValue("modified-name"),
		Description: types.StringValue("test-description"),
		Count:       types.Int64Value(42),
	}

	var tempModel TestModel
	plan := createTestPlan(tempModel.GetResourceSchema())

	processData := &process.ProcessData{
		CreateRequest:  &resource.CreateRequest{Plan: plan},
		CreateResponse: &resource.CreateResponse{},
	}

	mockSch := &mockSchema{schema: schema.Schema{}}
	err := CheckPlanVsApiResponseDelta(requestContext, processData, mockSch, apiResponseModel)

	assert.NoError(t, err)
}

func TestCheckPlanVsApiResponseDelta_NilPlan(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Create,
	}

	apiResponseModel := &TestModel{
		Name:        types.StringValue("test-name"),
		Description: types.StringValue("test-description"),
		Count:       types.Int64Value(42),
	}

	planModel := TestModel{
		Name:        types.StringValue("test-name"),
		Description: types.StringValue("test-description"),
		Count:       types.Int64Value(42),
	}
	plan := createTestPlanWithValues(planModel)

	processData := &process.ProcessData{
		CreateRequest:  &resource.CreateRequest{Plan: plan},
		CreateResponse: &resource.CreateResponse{},
	}

	mockSch := &mockSchema{schema: schema.Schema{}}
	err := CheckPlanVsApiResponseDelta(requestContext, processData, mockSch, apiResponseModel)

	assert.NoError(t, err)
	assert.Empty(t, processData.CreateResponse.Diagnostics)
}

func TestCheckPlanVsApiResponseDelta_NilApiResponse(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Create,
	}

	var nilApiResponse *TestModel

	var tempModel TestModel
	plan := createTestPlan(tempModel.GetResourceSchema())

	processData := &process.ProcessData{
		CreateRequest:  &resource.CreateRequest{Plan: plan},
		CreateResponse: &resource.CreateResponse{},
	}

	mockSch := &mockSchema{schema: schema.Schema{}}
	err := CheckPlanVsApiResponseDelta(requestContext, processData, mockSch, nilApiResponse)

	assert.NoError(t, err)
	assert.Empty(t, processData.CreateResponse.Diagnostics)
}

func TestCheckPlanVsApiResponseDelta_UpdateOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Update,
	}

	apiResponseModel := &TestModel{
		Name:        types.StringValue("test-name"),
		Description: types.StringValue("test-description"),
		Count:       types.Int64Value(42),
	}

	plan := createTestPlan(apiResponseModel.GetResourceSchema())

	processData := &process.ProcessData{
		UpdateRequest:  &resource.UpdateRequest{Plan: plan},
		UpdateResponse: &resource.UpdateResponse{},
	}

	mockSch := &mockSchema{schema: schema.Schema{}}
	err := CheckPlanVsApiResponseDelta(requestContext, processData, mockSch, apiResponseModel)

	assert.NoError(t, err)
}

func TestGetPlanTfModel_CreateOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Create,
	}

	var tempModel TestModel
	plan := createTestPlan(tempModel.GetResourceSchema())

	processData := &process.ProcessData{
		CreateRequest: &resource.CreateRequest{
			Plan: plan,
		},
	}

	var emptyModel TestModel
	result, err := getPlanTfModel[*TestModel](requestContext, processData)

	assert.NoError(t, err)
	assert.Equal(t, emptyModel.Name.IsNull(), result.Name.IsNull())
}

func TestGetPlanTfModel_UpdateOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Update,
	}

	var tempModel TestModel
	plan := createTestPlan(tempModel.GetResourceSchema())

	processData := &process.ProcessData{
		UpdateRequest: &resource.UpdateRequest{
			Plan: plan,
		},
	}

	var emptyModel TestModel
	result, err := getPlanTfModel[*TestModel](requestContext, processData)

	assert.NoError(t, err)
	assert.Equal(t, emptyModel.Name.IsNull(), result.Name.IsNull())
}

func TestGetPlanTfModel_ReadOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Read,
	}

	processData := &process.ProcessData{}

	result, err := getPlanTfModel[*TestModel](requestContext, processData)

	assert.NoError(t, err)
	assert.True(t, result == nil || result.Name.IsNull())
}

func TestGetPlanTfModel_DeleteOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Context:   context.Background(),
		Operation: common.Delete,
	}

	processData := &process.ProcessData{}

	result, err := getPlanTfModel[*TestModel](requestContext, processData)

	assert.NoError(t, err)
	assert.True(t, result == nil || result.Name.IsNull())
}
