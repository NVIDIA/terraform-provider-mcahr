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

package orchestrator

import (
	"context"
	"fmt"
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	actionapi "terraform/terraform-provider/provider/external_api/resources/actions"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	tfmodel "terraform/terraform-provider/provider/tf/core/model"
	"terraform/terraform-provider/provider/tf/core/process"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Simple test model that implements the required interface
type TestTFModel struct {
	BackendVersion *version.BackendVersion `json:"-" tfsdk:"-"`

	Name    types.String `tfsdk:"name" json:"name"`
	Command types.String `tfsdk:"command" json:"command"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,omitempty"`
	Timeout types.Int64  `tfsdk:"timeout" json:"timeout,omitempty"`
}

func (t TestTFModel) GetName() string {
	return t.Name.ValueString()
}

func (t TestTFModel) SetBackendVersion(backendVersion *version.BackendVersion) {
	t.BackendVersion = backendVersion
}

// MockResourceSchema implements coreschema.ResourceSchema for testing
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

// Helper function to create a simple ResourceSchema for testing
func createMockResourceSchema() coreschema.ResourceSchema {
	mockSchema := &MockResourceSchema{}

	// Mock GetSchema to return a simple schema
	mockSchema.On("GetSchema").Return(schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"command": schema.StringAttribute{
				Required: true,
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
			},
			"timeout": schema.Int64Attribute{
				Optional: true,
			},
		},
	})

	// Mock GetCompatibilityOptions to return empty options
	mockSchema.On("GetCompatibilityOptions").Return(map[string]attribute.CompatibilityOptions{})

	return mockSchema
}

// MockPreProcessor using testify/mock
type MockPreProcessor[TF tfmodel.TFModel] struct {
	mock.Mock
}

func (m *MockPreProcessor[TF]) PreProcessCreate(requestContext *common.RequestContext, data *process.ProcessData) (TF, error) {
	args := m.Called(requestContext, data)
	if args.Get(0) == nil {
		var nilTF TF
		return nilTF, args.Error(1)
	}
	return args.Get(0).(TF), args.Error(1)
}

func (m *MockPreProcessor[TF]) PreProcessRead(requestContext *common.RequestContext, data *process.ProcessData) (TF, error) {
	args := m.Called(requestContext, data)
	if args.Get(0) == nil {
		var nilTF TF
		return nilTF, args.Error(1)
	}
	return args.Get(0).(TF), args.Error(1)
}

func (m *MockPreProcessor[TF]) PreProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData) (TF, error) {
	args := m.Called(requestContext, data)
	if args.Get(0) == nil {
		var nilTF TF
		return nilTF, args.Error(1)
	}
	return args.Get(0).(TF), args.Error(1)
}

func (m *MockPreProcessor[TF]) PreProcessDelete(requestContext *common.RequestContext, data *process.ProcessData) (TF, error) {
	args := m.Called(requestContext, data)
	if args.Get(0) == nil {
		var nilTF TF
		return nilTF, args.Error(1)
	}
	return args.Get(0).(TF), args.Error(1)
}

// MockPostProcessor using testify/mock
type MockPostProcessor[TF tfmodel.TFModel] struct {
	mock.Mock
}

func (m *MockPostProcessor[TF]) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, tf TF) error {
	args := m.Called(requestContext, data, tf)
	return args.Error(0)
}

func (m *MockPostProcessor[TF]) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, tf TF) error {
	args := m.Called(requestContext, data, tf)
	return args.Error(0)
}

func (m *MockPostProcessor[TF]) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, tf TF) error {
	args := m.Called(requestContext, data, tf)
	return args.Error(0)
}

func (m *MockPostProcessor[TF]) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, tf TF) error {
	args := m.Called(requestContext, data, tf)
	return args.Error(0)
}

// MockTranslator using testify/mock
type MockTranslator[TF tfmodel.TFModel] struct {
	mock.Mock
}

func (m *MockTranslator[TF]) ToTFModel(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, apiModel *actionapi.ActionResponseAPIModel) (TF, error) {
	args := m.Called(requestContext, translationData, apiModel)
	if args.Get(0) == nil {
		var nilTF TF
		return nilTF, args.Error(1)
	}
	return args.Get(0).(TF), args.Error(1)
}

func (m *MockTranslator[TF]) ToAPIModel(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, tfModel TF) (*statement.StatementInputAPIModel, error) {
	args := m.Called(requestContext, translationData, tfModel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*statement.StatementInputAPIModel), args.Error(1)
}

// Helper functions
func createTestProcessData() *process.ProcessData {
	return &process.ProcessData{
		CreateRequest:  &resource.CreateRequest{},
		CreateResponse: &resource.CreateResponse{},
	}
}

func createTestRequestContext(backendVersion *version.BackendVersion) *common.RequestContext {
	return common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)
}

func createTestTranslationData() *coretranslator.TranslationData {
	return &coretranslator.TranslationData{}
}

func createTestModel() *TestTFModel {
	return &TestTFModel{
		Name:    types.StringValue("test_action"),
		Command: types.StringValue("echo test"),
		Enabled: types.BoolValue(true),
		Timeout: types.Int64Value(30),
	}
}

func TestOrchestrate_Success(t *testing.T) {
	t.Parallel()

	// given
	processData := createTestProcessData()
	platformClient := client.NewPlatformClient("http://test.com", "test-key")

	mockPreProcessor := &MockPreProcessor[*TestTFModel]{}
	mockPostProcessor := &MockPostProcessor[*TestTFModel]{}
	mockTranslator := &MockTranslator[*TestTFModel]{}

	testModel := createTestModel()
	apiInput := &statement.StatementInputAPIModel{Statement: "CREATE ACTION test_action"}
	apiResponse := &actionapi.ActionResponseAPIModel{
		Output: actionapi.ActionOutput{
			Configurations: actionapi.ActionConfigurations{
				Count: 1,
				Items: []actionapi.ConfigurationItem{
					{
						Config: actionapi.ActionConfig{
							Timeout:     30,
							CommandText: "echo test",
							StepDetails: actionapi.StepDetails{
								StartStep: actionapi.Step{
									Description: "",
									Title:       "",
								},
								ErrorStep: actionapi.Step{
									Description: "",
									Title:       "",
								},
								CompleteStep: actionapi.Step{
									Description: "",
									Title:       "",
								},
							},
						},
						EntityMetadata: actionapi.ActionEntityMetadata{
							Enabled: true,
							ID:      "action-123",
							Name:    "test_action",
						},
					},
				},
			},
		},
		Summary: actionapi.ActionSummary{
			Status: "success",
			Errors: []apicommon.Error{},
		},
	}

	// Set up mock expectations
	mockPreProcessor.On("PreProcessCreate", mock.AnythingOfType("*common.RequestContext"), processData).Return(testModel, nil)
	mockTranslator.On("ToAPIModel", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*translator.TranslationData"), testModel).Return(apiInput, nil)
	mockTranslator.On("ToTFModel", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*translator.TranslationData"), apiResponse).Return(testModel, nil)
	mockPostProcessor.On("PostProcessCreate", mock.AnythingOfType("*common.RequestContext"), processData, testModel).Return(nil)

	// Mock external API call
	mockExternalAPI := func(_ *common.RequestContext, _ *client.PlatformClient, apiInput *statement.StatementInputAPIModel) (*actionapi.ActionResponseAPIModel, error) {
		assert.Equal(t, "CREATE ACTION test_action", apiInput.Statement)
		return apiResponse, nil
	}

	requestContext := createTestRequestContext(version.NewBackendVersion("release-29.1.0"))
	translationData := createTestTranslationData()

	// when
	result, _, err := orchestrateWithAPIFunction(
		requestContext,
		platformClient,
		createMockResourceSchema(),
		mockPreProcessor,
		mockPostProcessor,
		processData,
		mockTranslator,
		translationData,
		mockExternalAPI,
	)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Assert on the result directly - verify orchestrator returned the expected model
	assert.Equal(t, "test_action", result.Name.ValueString())
	assert.Equal(t, "echo test", result.Command.ValueString())
	assert.Equal(t, true, result.Enabled.ValueBool())
	assert.Equal(t, int64(30), result.Timeout.ValueInt64())

	// Verify all mocks were called correctly
	mockPreProcessor.AssertExpectations(t)
	mockPostProcessor.AssertExpectations(t)
	mockTranslator.AssertExpectations(t)
}

func TestOrchestrate_PreProcessorError(t *testing.T) {
	t.Parallel()
	// given
	processData := createTestProcessData()
	platformClient := client.NewPlatformClient("http://test.com", "test-key")

	mockPreProcessor := &MockPreProcessor[*TestTFModel]{}
	mockPostProcessor := &MockPostProcessor[*TestTFModel]{}
	mockTranslator := &MockTranslator[*TestTFModel]{}

	// Mock preprocessor to return error
	mockPreProcessor.On("PreProcessCreate", mock.AnythingOfType("*common.RequestContext"), processData).Return((*TestTFModel)(nil), assert.AnError)

	// Mock external API call (won't be called due to preprocessor error)
	mockExternalAPI := func(_ *common.RequestContext, _ *client.PlatformClient, _ *statement.StatementInputAPIModel) (*actionapi.ActionResponseAPIModel, error) {
		return nil, fmt.Errorf("should not be called")
	}

	requestContext := createTestRequestContext(version.NewBackendVersion("release-29.1.0"))
	translationData := createTestTranslationData()

	// when
	result, _, err := orchestrateWithAPIFunction(
		requestContext,
		platformClient,
		createMockResourceSchema(),
		mockPreProcessor,
		mockPostProcessor,
		processData,
		mockTranslator,
		translationData,
		mockExternalAPI,
	)

	// then
	require.Error(t, err)
	require.Nil(t, result, "Result should be nil when PreProcessor fails")
	assert.Equal(t, assert.AnError, err, "Should return the exact error from PreProcessor")

	// Verify orchestrator stopped at PreProcessor and didn't call subsequent steps
	mockPreProcessor.AssertExpectations(t)
	mockPostProcessor.AssertNotCalled(t, "PostProcessCreate")
	mockTranslator.AssertNotCalled(t, "ToAPIModel")
	mockTranslator.AssertNotCalled(t, "ToTFModel")
}

func TestOrchestrate_TranslatorToAPIError(t *testing.T) {
	t.Parallel()
	// given
	processData := createTestProcessData()
	platformClient := client.NewPlatformClient("http://test.com", "test-key")

	mockPreProcessor := &MockPreProcessor[*TestTFModel]{}
	mockPostProcessor := &MockPostProcessor[*TestTFModel]{}
	mockTranslator := &MockTranslator[*TestTFModel]{}

	testModel := createTestModel()

	// Set up expectations
	mockPreProcessor.On("PreProcessCreate", mock.AnythingOfType("*common.RequestContext"), processData).Return(testModel, nil)
	mockTranslator.On("ToAPIModel", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*translator.TranslationData"), testModel).Return(nil, assert.AnError)

	// Mock external API call (won't be called due to translator error)
	mockExternalAPI := func(_ *common.RequestContext, _ *client.PlatformClient, _ *statement.StatementInputAPIModel) (*actionapi.ActionResponseAPIModel, error) {
		return nil, fmt.Errorf("should not be called")
	}

	requestContext := createTestRequestContext(version.NewBackendVersion("release-29.1.0"))
	translationData := createTestTranslationData()

	// when
	result, _, err := orchestrateWithAPIFunction(
		requestContext,
		platformClient,
		createMockResourceSchema(),
		mockPreProcessor,
		mockPostProcessor,
		processData,
		mockTranslator,
		translationData,
		mockExternalAPI,
	)

	// then
	require.Error(t, err)
	require.Nil(t, result, "Result should be nil when Translator fails")
	assert.Equal(t, assert.AnError, err, "Should return the exact error from Translator")

	// Verify orchestrator processed PreProcessor but stopped at Translator
	mockPreProcessor.AssertExpectations(t)
	mockTranslator.AssertExpectations(t)

	// PostProcessor and ToTFModel should not be called when Translator.ToAPIModel fails
	mockPostProcessor.AssertNotCalled(t, "PostProcessCreate")
	mockTranslator.AssertNotCalled(t, "ToTFModel")
}

func TestOrchestrate_ExternalAPIError(t *testing.T) {
	t.Parallel()
	// given
	processData := createTestProcessData()
	platformClient := client.NewPlatformClient("http://test.com", "test-key")

	mockPreProcessor := &MockPreProcessor[*TestTFModel]{}
	mockPostProcessor := &MockPostProcessor[*TestTFModel]{}
	mockTranslator := &MockTranslator[*TestTFModel]{}

	testModel := createTestModel()
	apiInput := &statement.StatementInputAPIModel{Statement: "CREATE ACTION test_action"}

	// Set up expectations
	mockPreProcessor.On("PreProcessCreate", mock.AnythingOfType("*common.RequestContext"), processData).Return(testModel, nil)
	mockTranslator.On("ToAPIModel", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*translator.TranslationData"), testModel).Return(apiInput, nil)

	// Mock API to return error
	mockExternalAPI := func(_ *common.RequestContext, _ *client.PlatformClient, _ *statement.StatementInputAPIModel) (*actionapi.ActionResponseAPIModel, error) {
		return nil, assert.AnError
	}

	requestContext := createTestRequestContext(version.NewBackendVersion("release-29.1.0"))
	translationData := createTestTranslationData()

	// when
	result, _, err := orchestrateWithAPIFunction(
		requestContext,
		platformClient,
		createMockResourceSchema(),
		mockPreProcessor,
		mockPostProcessor,
		processData,
		mockTranslator,
		translationData,
		mockExternalAPI,
	)

	// then
	require.Error(t, err)
	require.Nil(t, result, "Result should be nil when external API fails")
	assert.Equal(t, assert.AnError, err, "Should return the exact error from external API")

	// Verify orchestrator processed PreProcessor and Translator.ToAPIModel but stopped at external API
	mockPreProcessor.AssertExpectations(t)
	mockTranslator.AssertExpectations(t)

	// PostProcessor and ToTFModel should not be called when external API fails
	mockPostProcessor.AssertNotCalled(t, "PostProcessCreate")
	mockTranslator.AssertNotCalled(t, "ToTFModel")
}

func TestOrchestrate_APIBusinessError(t *testing.T) {
	t.Parallel()
	// given
	processData := createTestProcessData()
	platformClient := client.NewPlatformClient("http://test.com", "test-key")

	mockPreProcessor := &MockPreProcessor[*TestTFModel]{}
	mockPostProcessor := &MockPostProcessor[*TestTFModel]{}
	mockTranslator := &MockTranslator[*TestTFModel]{}

	testModel := createTestModel()
	apiInput := &statement.StatementInputAPIModel{Statement: "CREATE ACTION test_action"}

	// Set up expectations
	mockPreProcessor.On("PreProcessCreate", mock.AnythingOfType("*common.RequestContext"), processData).Return(testModel, nil)
	mockTranslator.On("ToAPIModel", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*translator.TranslationData"), testModel).Return(apiInput, nil)

	// Mock API to return business error
	mockExternalAPI := func(_ *common.RequestContext, _ *client.PlatformClient, _ *statement.StatementInputAPIModel) (*actionapi.ActionResponseAPIModel, error) {
		apiResponse := &actionapi.ActionResponseAPIModel{
			Summary: actionapi.ActionSummary{
				Status: "error",
				Errors: []apicommon.Error{
					{Type: "DUPLICATE_NAME", Message: "Action name already exists"},
				},
			},
		}
		return nil, fmt.Errorf("API business error: %s", apiResponse.GetErrors())
	}

	requestContext := createTestRequestContext(version.NewBackendVersion("release-29.1.0"))
	translationData := createTestTranslationData()

	// when
	result, _, err := orchestrateWithAPIFunction(
		requestContext,
		platformClient,
		createMockResourceSchema(),
		mockPreProcessor,
		mockPostProcessor,
		processData,
		mockTranslator,
		translationData,
		mockExternalAPI,
	)

	// then
	require.Error(t, err)
	require.Nil(t, result, "Result should be nil when external API returns business errors")

	// Verify the error contains the business error details
	assert.Contains(t, err.Error(), "API business error", "Error should indicate it's from API")
	assert.Contains(t, err.Error(), "DUPLICATE_NAME", "Error should contain specific error type")
	assert.Contains(t, err.Error(), "Action name already exists", "Error should contain error message")

	// Verify orchestrator processed PreProcessor and Translator.ToAPIModel but stopped at external API
	mockPreProcessor.AssertExpectations(t)
	mockTranslator.AssertExpectations(t)

	// PostProcessor and ToTFModel should not be called when business errors occur
	mockPostProcessor.AssertNotCalled(t, "PostProcessCreate")
	mockTranslator.AssertNotCalled(t, "ToTFModel")
}

func TestDetermineAPIVersion(t *testing.T) {
	tests := []struct {
		name           string
		backendVersion *version.BackendVersion
		expectedAPI    common.APIVersion
	}{
		{
			name:           "version below threshold (28.5.0) should use V1",
			backendVersion: version.NewBackendVersion("release-28.5.0"),
			expectedAPI:    common.V1,
		},
		{
			name:           "version below threshold (27.10.2) should use V1",
			backendVersion: version.NewBackendVersion("release-27.10.2"),
			expectedAPI:    common.V1,
		},
		{
			name:           "version below threshold (28.99.99) should use V1",
			backendVersion: version.NewBackendVersion("release-28.99.99"),
			expectedAPI:    common.V1,
		},
		{
			name:           "version below threshold (29.0.0) should use V1",
			backendVersion: version.NewBackendVersion("release-29.0.0"),
			expectedAPI:    common.V1,
		},
		{
			name:           "version equal to threshold (29.1.0) should use V2",
			backendVersion: version.NewBackendVersion("release-29.1.0"),
			expectedAPI:    common.V2,
		},
		{
			name:           "version above threshold (30.1.5) should use V2",
			backendVersion: version.NewBackendVersion("release-30.1.5"),
			expectedAPI:    common.V2,
		},
		{
			name:           "version above threshold (35.2.1) should use V2",
			backendVersion: version.NewBackendVersion("release-35.2.1"),
			expectedAPI:    common.V2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineAPIVersion(tt.backendVersion)
			if result != tt.expectedAPI {
				t.Errorf("determineAPIVersion() = %v, want %v", result, tt.expectedAPI)
			}
		})
	}
}
