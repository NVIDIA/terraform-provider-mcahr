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
	"testing"

	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/common/version"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestModel for testing
type TestModel struct {
	Name            types.String `tfsdk:"name" json:"name"`
	CompatibleField types.String `tfsdk:"compatible_field" json:"compatible_field"`
	MinVersionField types.String `tfsdk:"min_version_field" json:"min_version_field"`
	MaxVersionField types.String `tfsdk:"max_version_field" json:"max_version_field"`
	NoTagField      types.String `json:"no_tag"`
	SkipField       types.String `tfsdk:"-" json:"skip"`
}

func (t *TestModel) GetName() string {
	if t.Name.IsNull() || t.Name.IsUnknown() {
		return ""
	}
	return t.Name.ValueString()
}

// MockResourceSchema for testing
type MockResourceSchema struct {
	mock.Mock
}

func (m *MockResourceSchema) GetSchema() schema.Schema {
	args := m.Called()
	return args.Get(0).(schema.Schema)
}

func (m *MockResourceSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	args := m.Called()
	return args.Get(0).(map[string]attribute.CompatibilityOptions)
}

// Helper to create test plan
func createTestPlan(schema schema.Schema, values map[string]tftypes.Value) tfsdk.Plan {
	attrTypes := make(map[string]tftypes.Type)
	for name := range schema.Attributes {
		attrTypes[name] = tftypes.String
	}

	planData := tftypes.NewValue(tftypes.Object{
		AttributeTypes: attrTypes,
	}, values)

	return tfsdk.Plan{
		Raw:    planData,
		Schema: schema,
	}
}

// Helper to create test config
func createTestConfig(schema schema.Schema, values map[string]tftypes.Value) tfsdk.Config {
	attrTypes := make(map[string]tftypes.Type)
	for name := range schema.Attributes {
		attrTypes[name] = tftypes.String
	}

	configData := tftypes.NewValue(tftypes.Object{
		AttributeTypes: attrTypes,
	}, values)

	return tfsdk.Config{
		Raw:    configData,
		Schema: schema,
	}
}

func TestConstructCompatibilityErrorMessage_WithMinVersion(t *testing.T) {
	// given
	backendVersion := version.NewBackendVersion("release-1.0.0")
	options := map[string]attribute.CompatibilityOptions{
		"test_field": {
			MinVersion: "release-2.0.0",
		},
	}
	checker := attribute.NewCompatibilityChecker(backendVersion, options)

	// when
	result := constructCompatibilityErrorMessage(checker, "test_field")

	// then
	expected := "test_field attribute is not supported by the current platform. Current version: release-1.0.0. Minimum version required: release-2.0.0."
	assert.Equal(t, expected, result)
}

func TestConstructCompatibilityErrorMessage_WithMaxVersion(t *testing.T) {
	// given
	backendVersion := version.NewBackendVersion("release-3.0.0")
	options := map[string]attribute.CompatibilityOptions{
		"test_field": {
			MaxVersion: "release-2.0.0",
		},
	}
	checker := attribute.NewCompatibilityChecker(backendVersion, options)

	// when
	result := constructCompatibilityErrorMessage(checker, "test_field")

	// then
	expected := "test_field attribute is not supported by the current platform. Current version: release-3.0.0. Maximum allowed version: release-2.0.0."
	assert.Equal(t, expected, result)
}

func TestConstructCompatibilityErrorMessage_WithBothVersions(t *testing.T) {
	// given
	backendVersion := version.NewBackendVersion("release-1.0.0")
	options := map[string]attribute.CompatibilityOptions{
		"test_field": {
			MinVersion: "release-2.0.0",
			MaxVersion: "release-5.0.0",
		},
	}
	checker := attribute.NewCompatibilityChecker(backendVersion, options)

	// when
	result := constructCompatibilityErrorMessage(checker, "test_field")

	// then
	expected := "test_field attribute is not supported by the current platform. Current version: release-1.0.0. Minimum version required: release-2.0.0. Maximum allowed version: release-5.0.0."
	assert.Equal(t, expected, result)
}

func TestConstructCompatibilityErrorMessage_NoOptions(t *testing.T) {
	// given
	backendVersion := version.NewBackendVersion("release-1.0.0")
	options := map[string]attribute.CompatibilityOptions{}
	checker := attribute.NewCompatibilityChecker(backendVersion, options)

	// when
	result := constructCompatibilityErrorMessage(checker, "test_field")

	// then
	expected := "test_field attribute is not supported by the current platform. Current version: release-1.0.0."
	assert.Equal(t, expected, result)
}

func TestApplyVersionValidationModifier_SkipDestroyOperation(t *testing.T) {
	ctx := context.Background()

	// given - create a destroy operation (plan is null)
	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	nullPlanData := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name": tftypes.String,
		},
	}, nil)

	req := &resource.ModifyPlanRequest{
		Plan: tfsdk.Plan{
			Raw:    nullPlanData,
			Schema: planSchema,
		},
	}
	resp := &resource.ModifyPlanResponse{
		Plan: req.Plan,
	}

	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{})
	configValues := &TestModel{}

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should not add any diagnostics
	assert.False(t, resp.Diagnostics.HasError())
}

func TestApplyVersionValidationModifier_CompatibleAttribute(t *testing.T) {
	ctx := context.Background()

	// given - all attributes are compatible
	configValues := &TestModel{
		Name:            types.StringValue("test"),
		CompatibleField: types.StringValue("value"),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"compatible_field": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := &resource.ModifyPlanRequest{
		Plan: createTestPlan(planSchema, map[string]tftypes.Value{
			"name":             tftypes.NewValue(tftypes.String, "test"),
			"compatible_field": tftypes.NewValue(tftypes.String, "value"),
		}),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	backendVersion := version.NewBackendVersion("release-2.0.0")
	// No compatibility restrictions - all fields compatible
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{})

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should not add any diagnostics
	assert.False(t, resp.Diagnostics.HasError())
}

func TestApplyVersionValidationModifier_IncompatibleAttributeWithValue(t *testing.T) {
	ctx := context.Background()

	// given - user provides value for incompatible attribute
	configValues := &TestModel{
		Name:            types.StringValue("test"),
		MinVersionField: types.StringValue("user_provided_value"),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"min_version_field": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := &resource.ModifyPlanRequest{
		Plan: createTestPlan(planSchema, map[string]tftypes.Value{
			"name":              tftypes.NewValue(tftypes.String, "test"),
			"min_version_field": tftypes.NewValue(tftypes.String, "user_provided_value"),
		}),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	backendVersion := version.NewBackendVersion("release-1.0.0")
	// min_version_field requires version 2.0.0 or higher
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{
		"min_version_field": {
			MinVersion: "release-2.0.0",
		},
	})

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should add error diagnostic
	assert.True(t, resp.Diagnostics.HasError())
	errors := resp.Diagnostics.Errors()
	require.Len(t, errors, 1)
	assert.Equal(t, "Unsupported attribute", errors[0].Summary())
	expectedDetail := "min_version_field attribute is not supported by the current platform. Current version: release-1.0.0. Minimum version required: release-2.0.0."
	assert.Equal(t, expectedDetail, errors[0].Detail())
}

func TestApplyVersionValidationModifier_IncompatibleAttributeUnknown(t *testing.T) {
	ctx := context.Background()

	// given - attribute is unknown (not provided by user)
	configValues := &TestModel{
		Name:            types.StringValue("test"),
		MinVersionField: types.StringUnknown(),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"min_version_field": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := &resource.ModifyPlanRequest{
		Plan: createTestPlan(planSchema, map[string]tftypes.Value{
			"name":              tftypes.NewValue(tftypes.String, "test"),
			"min_version_field": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		}),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{
		"min_version_field": {
			MinVersion: "release-2.0.0",
		},
	})

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should not add error, but should set to null
	assert.False(t, resp.Diagnostics.HasError())

	// Verify the attribute was set to null
	var resultValue attr.Value
	diags := resp.Plan.GetAttribute(ctx, path.Root("min_version_field"), &resultValue)
	require.False(t, diags.HasError())
	assert.True(t, resultValue.IsNull())
}

func TestApplyVersionValidationModifier_IncompatibleAttributeNull(t *testing.T) {
	ctx := context.Background()

	// given - attribute is null (not provided by user)
	configValues := &TestModel{
		Name:            types.StringValue("test"),
		MinVersionField: types.StringNull(),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"min_version_field": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := &resource.ModifyPlanRequest{
		Plan: createTestPlan(planSchema, map[string]tftypes.Value{
			"name":              tftypes.NewValue(tftypes.String, "test"),
			"min_version_field": tftypes.NewValue(tftypes.String, nil),
		}),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{
		"min_version_field": {
			MinVersion: "release-2.0.0",
		},
	})

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should not add error
	assert.False(t, resp.Diagnostics.HasError())
}

func TestApplyVersionValidationModifier_SkipFieldsWithoutTfsdkTag(t *testing.T) {
	ctx := context.Background()

	// given - field without tfsdk tag
	configValues := &TestModel{
		Name:       types.StringValue("test"),
		NoTagField: types.StringValue("should_be_ignored"),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := &resource.ModifyPlanRequest{
		Plan: createTestPlan(planSchema, map[string]tftypes.Value{
			"name": tftypes.NewValue(tftypes.String, "test"),
		}),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{})

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should not add any diagnostics
	assert.False(t, resp.Diagnostics.HasError())
}

func TestApplyVersionValidationModifier_SkipFieldsWithDashTag(t *testing.T) {
	ctx := context.Background()

	// given - field with tfsdk:"-" tag
	configValues := &TestModel{
		Name:      types.StringValue("test"),
		SkipField: types.StringValue("should_be_skipped"),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := &resource.ModifyPlanRequest{
		Plan: createTestPlan(planSchema, map[string]tftypes.Value{
			"name": tftypes.NewValue(tftypes.String, "test"),
		}),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	backendVersion := version.NewBackendVersion("release-1.0.0")
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{})

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should not add any diagnostics
	assert.False(t, resp.Diagnostics.HasError())
}

func TestApplyVersionValidationModifier_MaxVersionExceeded(t *testing.T) {
	ctx := context.Background()

	// given - backend version exceeds max version
	configValues := &TestModel{
		Name:            types.StringValue("test"),
		MaxVersionField: types.StringValue("user_provided_value"),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"max_version_field": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := &resource.ModifyPlanRequest{
		Plan: createTestPlan(planSchema, map[string]tftypes.Value{
			"name":              tftypes.NewValue(tftypes.String, "test"),
			"max_version_field": tftypes.NewValue(tftypes.String, "user_provided_value"),
		}),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	backendVersion := version.NewBackendVersion("release-3.0.0")
	// max_version_field only supported up to version 2.0.0
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{
		"max_version_field": {
			MaxVersion: "release-2.0.0",
		},
	})

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should add error diagnostic
	assert.True(t, resp.Diagnostics.HasError())
	errors := resp.Diagnostics.Errors()
	require.Len(t, errors, 1)
	assert.Equal(t, "Unsupported attribute", errors[0].Summary())
	expectedDetail := "max_version_field attribute is not supported by the current platform. Current version: release-3.0.0. Maximum allowed version: release-2.0.0."
	assert.Equal(t, expectedDetail, errors[0].Detail())
}

func TestApplyVersionValidationModifier_MultipleIncompatibleFields(t *testing.T) {
	ctx := context.Background()

	// given - multiple incompatible fields with user values
	configValues := &TestModel{
		Name:            types.StringValue("test"),
		MinVersionField: types.StringValue("value1"),
		MaxVersionField: types.StringValue("value2"),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"min_version_field": schema.StringAttribute{
				Optional: true,
			},
			"max_version_field": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := &resource.ModifyPlanRequest{
		Plan: createTestPlan(planSchema, map[string]tftypes.Value{
			"name":              tftypes.NewValue(tftypes.String, "test"),
			"min_version_field": tftypes.NewValue(tftypes.String, "value1"),
			"max_version_field": tftypes.NewValue(tftypes.String, "value2"),
		}),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	backendVersion := version.NewBackendVersion("release-1.5.0")
	compatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, map[string]attribute.CompatibilityOptions{
		"min_version_field": {
			MinVersion: "release-2.0.0",
		},
		"max_version_field": {
			MaxVersion: "release-1.0.0",
		},
	})

	// when
	applyVersionValidationModifier(ctx, req, resp, compatibilityChecker, configValues)

	// then - should add error diagnostics for both fields
	assert.True(t, resp.Diagnostics.HasError())
	errors := resp.Diagnostics.Errors()
	assert.Len(t, errors, 2)
}

func TestApplyCompatibilityModifiers_Integration(t *testing.T) {
	ctx := context.Background()

	// given
	configValues := &TestModel{
		Name:            types.StringValue("test"),
		MinVersionField: types.StringValue("incompatible_value"),
	}

	planSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"min_version_field": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	testValues := map[string]tftypes.Value{
		"name":              tftypes.NewValue(tftypes.String, "test"),
		"min_version_field": tftypes.NewValue(tftypes.String, "incompatible_value"),
	}

	req := &resource.ModifyPlanRequest{
		Plan:   createTestPlan(planSchema, testValues),
		Config: createTestConfig(planSchema, testValues),
	}
	resp := &resource.ModifyPlanResponse{
		Plan:        req.Plan,
		Diagnostics: diag.Diagnostics{},
	}

	mockSchema := &MockResourceSchema{}
	mockSchema.On("GetCompatibilityOptions").Return(map[string]attribute.CompatibilityOptions{
		"min_version_field": {
			MinVersion: "release-2.0.0",
		},
	})

	backendVersion := version.NewBackendVersion("release-1.0.0")

	// when
	ApplyCompatibilityModifiers(ctx, req, resp, mockSchema, backendVersion, configValues)

	// then
	assert.True(t, resp.Diagnostics.HasError())
	mockSchema.AssertExpectations(t)
}

// Note: ApplyResourceCompatibilityModifiers is a thin wrapper around ApplyCompatibilityModifiers
// that first extracts config. It's thoroughly tested through integration tests and the underlying
// ApplyCompatibilityModifiers function is comprehensively tested above.

func TestConstructCompatibilityErrorMessage_VariousScenarios(t *testing.T) {
	tests := []struct {
		name            string
		currentVersion  string
		minVersion      string
		maxVersion      string
		attributeName   string
		expectedMessage string
	}{
		{
			name:            "Only min version constraint",
			currentVersion:  "release-1.0.0",
			minVersion:      "release-2.0.0",
			maxVersion:      "",
			attributeName:   "feature_flag",
			expectedMessage: "feature_flag attribute is not supported by the current platform. Current version: release-1.0.0. Minimum version required: release-2.0.0.",
		},
		{
			name:            "Only max version constraint",
			currentVersion:  "release-5.0.0",
			minVersion:      "",
			maxVersion:      "release-3.0.0",
			attributeName:   "deprecated_field",
			expectedMessage: "deprecated_field attribute is not supported by the current platform. Current version: release-5.0.0. Maximum allowed version: release-3.0.0.",
		},
		{
			name:            "Both constraints",
			currentVersion:  "release-1.0.0",
			minVersion:      "release-2.0.0",
			maxVersion:      "release-4.0.0",
			attributeName:   "limited_field",
			expectedMessage: "limited_field attribute is not supported by the current platform. Current version: release-1.0.0. Minimum version required: release-2.0.0. Maximum allowed version: release-4.0.0.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			backendVersion := version.NewBackendVersion(tt.currentVersion)
			options := map[string]attribute.CompatibilityOptions{
				tt.attributeName: {
					MinVersion: tt.minVersion,
					MaxVersion: tt.maxVersion,
				},
			}
			checker := attribute.NewCompatibilityChecker(backendVersion, options)

			// when
			result := constructCompatibilityErrorMessage(checker, tt.attributeName)

			// then
			assert.Equal(t, tt.expectedMessage, result)
		})
	}
}
