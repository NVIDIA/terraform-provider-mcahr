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

package externalapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/external_api/resources/statement"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Generic test models that implement the APIModel interface
type GenericSuccessResponse struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Status  string            `json:"status"`
	Data    map[string]string `json:"data"`
	Count   int               `json:"count"`
	Enabled bool              `json:"enabled"`
}

func (g GenericSuccessResponse) GetErrors() string {
	if g.Status == "error" {
		return "Status: " + g.Status
	}
	return ""
}

type GenericErrorResponse struct {
	Status string         `json:"status"`
	Errors []GenericError `json:"errors"`
}

type GenericError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (g GenericErrorResponse) GetErrors() string {
	if len(g.Errors) == 0 {
		return ""
	}

	result := "Status: " + g.Status + "; Errors: "
	for i, err := range g.Errors {
		if i > 0 {
			result += ", "
		}
		result += err.Code + ": " + err.Message
	}
	return result
}

// V1-style generic response (unified container pattern)
type GenericResponseV1 struct {
	Operation *GenericContainerV1 `json:"operation,omitempty"`
}

type GenericContainerV1 struct {
	Error GenericErrorV1  `json:"error"`
	Count *int            `json:"count,omitempty"`
	Items []GenericItemV1 `json:"items"`
}

type GenericErrorV1 struct {
	Type    string  `json:"type"`
	Message *string `json:"message"`
}

type GenericItemV1 struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Enabled  bool              `json:"enabled"`
	Settings map[string]string `json:"settings"`
}

func (g GenericResponseV1) GetErrors() string {
	if g.Operation == nil || g.Operation.Error.Type == "OK" || g.Operation.Error.Type == "" {
		return ""
	}

	result := "Error Type: " + g.Operation.Error.Type
	if g.Operation.Error.Message != nil && *g.Operation.Error.Message != "" {
		result += "; Message: " + *g.Operation.Error.Message
	}
	return result
}

// MockPlatformClient implements a mock for testing
type MockPlatformClient struct {
	mock.Mock
}

// Verify that PlatformClient implements our interface (compile-time check)
var _ PlatformClientInterface = (*client.PlatformClient)(nil)

func (m *MockPlatformClient) ExecuteRequest(requestContext *common.RequestContext, request *client.PlatformClientRequest) (*client.PlatformClientResponse, error) {
	args := m.Called(requestContext, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.PlatformClientResponse), args.Error(1)
}

// Verify that our mock implements the interface (compile-time check)
var _ PlatformClientInterface = (*MockPlatformClient)(nil)

func TestCallExternalAPI_SuccessfulResponseV2(t *testing.T) {
	// Given
	mockClient := new(MockPlatformClient)

	requestContext := common.NewRequestContext(context.Background())

	expectedResponse := GenericSuccessResponse{
		ID:     "test-resource-123",
		Name:   "test-resource",
		Status: "success",
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Count:   42,
		Enabled: true,
	}

	responseBody, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	platformResp := &client.PlatformClientResponse{
		Response: &http.Response{StatusCode: 200},
		Body:     responseBody,
	}
	mockClient.On("ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*client.PlatformClientRequest")).Return(platformResp, nil)

	apiObject := &statement.StatementInputAPIModel{
		Statement:  "get_resource(name=\"test-resource\")",
		APIVersion: common.V2,
	}

	// When
	result, err := CallExternalAPIWithClient[*GenericSuccessResponse](requestContext, mockClient, apiObject)

	// Then
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test-resource-123", result.ID)
	assert.Equal(t, "test-resource", result.Name)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, 42, result.Count)
	assert.True(t, result.Enabled)
	assert.Equal(t, "value1", result.Data["key1"])

	// Verify correct V2 endpoint was called
	mockClient.AssertCalled(t, "ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.MatchedBy(func(req *client.PlatformClientRequest) bool {
		return req.Method == "POST" && req.Endpoint == "/api/v1/statements/execute"
	}))
	mockClient.AssertExpectations(t)
}

func TestCallExternalAPI_SuccessfulResponseV1(t *testing.T) {
	// Given
	mockClient := new(MockPlatformClient)

	requestContext := common.NewRequestContext(context.Background())

	expectedResponse := GenericResponseV1{
		Operation: &GenericContainerV1{
			Error: GenericErrorV1{
				Type:    "OK",
				Message: nil,
			},
			Count: intPtr(2),
			Items: []GenericItemV1{
				{
					ID:      "item-456",
					Name:    "test-item",
					Enabled: true,
					Settings: map[string]string{
						"setting1": "value1",
						"setting2": "value2",
					},
				},
				{
					ID:      "item-789",
					Name:    "another-item",
					Enabled: false,
					Settings: map[string]string{
						"setting1": "different_value",
					},
				},
			},
		},
	}

	responseBody, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	platformResp := &client.PlatformClientResponse{
		Response: &http.Response{StatusCode: 200},
		Body:     responseBody,
	}
	mockClient.On("ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*client.PlatformClientRequest")).Return(platformResp, nil)

	apiObject := &statement.StatementInputAPIModel{
		Statement:  "operation(param=\"value\")",
		APIVersion: common.V1,
	}

	// When
	result, err := CallExternalAPIWithClient[*GenericResponseV1](requestContext, mockClient, apiObject)

	// Then
	assert.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Operation)
	assert.Equal(t, "OK", result.Operation.Error.Type)
	assert.Equal(t, 2, *result.Operation.Count)
	assert.Len(t, result.Operation.Items, 2)
	assert.Equal(t, "item-456", result.Operation.Items[0].ID)
	assert.Equal(t, "test-item", result.Operation.Items[0].Name)
	assert.True(t, result.Operation.Items[0].Enabled)
	assert.Equal(t, "value1", result.Operation.Items[0].Settings["setting1"])

	// Verify correct V1 endpoint was called
	mockClient.AssertCalled(t, "ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.MatchedBy(func(req *client.PlatformClientRequest) bool {
		return req.Method == "POST" && req.Endpoint == "/api/v1/execute"
	}))
	mockClient.AssertExpectations(t)
}

func TestCallExternalAPI_ComplexDataResponseV2(t *testing.T) {
	// Given
	mockClient := new(MockPlatformClient)

	requestContext := common.NewRequestContext(context.Background())

	expectedResponse := GenericSuccessResponse{
		ID:     "operation-303",
		Name:   "complex-response",
		Status: "success",
		Data: map[string]string{
			"metadata1": "value1",
			"metadata2": "value2",
			"info":      "additional_data",
		},
		Count:   42,
		Enabled: true,
	}

	responseBody, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	platformResp := &client.PlatformClientResponse{
		Response: &http.Response{StatusCode: 200},
		Body:     responseBody,
	}
	mockClient.On("ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*client.PlatformClientRequest")).Return(platformResp, nil)

	apiObject := &statement.StatementInputAPIModel{
		Statement:  "complex_operation(param=value)",
		APIVersion: common.V2,
	}

	// When
	result, err := CallExternalAPIWithClient[*GenericSuccessResponse](requestContext, mockClient, apiObject)

	// Then
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "operation-303", result.ID)
	assert.Equal(t, "complex-response", result.Name)
	assert.Equal(t, 42, result.Count)
	assert.Equal(t, "value1", result.Data["metadata1"])
	assert.Equal(t, "additional_data", result.Data["info"])

	mockClient.AssertCalled(t, "ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.MatchedBy(func(req *client.PlatformClientRequest) bool {
		return req.Method == "POST" && req.Endpoint == "/api/v1/statements/execute"
	}))
	mockClient.AssertExpectations(t)
}

func TestCallExternalAPI_BusinessErrorResponseV2(t *testing.T) {
	// Given
	mockClient := new(MockPlatformClient)

	requestContext := common.NewRequestContext(context.Background())

	expectedResponse := GenericErrorResponse{
		Status: "error",
		Errors: []GenericError{
			{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid parameter value",
			},
		},
	}

	responseBody, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	platformResp := &client.PlatformClientResponse{
		Response: &http.Response{StatusCode: 400},
		Body:     responseBody,
	}
	mockClient.On("ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*client.PlatformClientRequest")).Return(platformResp, nil)

	apiObject := &statement.StatementInputAPIModel{
		Statement:  "create_invalid_item(name=\"\")",
		APIVersion: common.V2,
	}

	// When
	result, err := CallExternalAPIWithClient[*GenericErrorResponse](requestContext, mockClient, apiObject)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API response errors")
	assert.Contains(t, err.Error(), "VALIDATION_ERROR")
	assert.Contains(t, err.Error(), "Invalid parameter value")

	mockClient.AssertExpectations(t)
}

func TestCallExternalAPI_BusinessErrorResponseV1(t *testing.T) {
	// Given
	mockClient := new(MockPlatformClient)

	requestContext := common.NewRequestContext(context.Background())

	expectedResponse := GenericResponseV1{
		Operation: &GenericContainerV1{
			Error: GenericErrorV1{
				Type:    "DUPLICATE_ERROR",
				Message: stringPtr("Item with this name already exists"),
			},
			Count: intPtr(0),
			Items: []GenericItemV1{},
		},
	}

	responseBody, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	platformResp := &client.PlatformClientResponse{
		Response: &http.Response{StatusCode: 400},
		Body:     responseBody,
	}
	mockClient.On("ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*client.PlatformClientRequest")).Return(platformResp, nil)

	apiObject := &statement.StatementInputAPIModel{
		Statement:  "create_item(name=\"duplicate-name\")",
		APIVersion: common.V1,
	}

	// When
	result, err := CallExternalAPIWithClient[*GenericResponseV1](requestContext, mockClient, apiObject)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API response errors")
	assert.Contains(t, err.Error(), "DUPLICATE_ERROR")
	assert.Contains(t, err.Error(), "Item with this name already exists")

	mockClient.AssertExpectations(t)
}

func TestCallExternalAPI_MultipleErrorsResponseV2(t *testing.T) {
	// Given
	mockClient := new(MockPlatformClient)

	requestContext := common.NewRequestContext(context.Background())

	expectedResponse := GenericErrorResponse{
		Status: "error",
		Errors: []GenericError{
			{
				Code:    "PERMISSION_ERROR",
				Message: "Access denied",
			},
			{
				Code:    "RESOURCE_ERROR",
				Message: "Resource not found",
			},
		},
	}

	responseBody, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	platformResp := &client.PlatformClientResponse{
		Response: &http.Response{StatusCode: 403},
		Body:     responseBody,
	}
	mockClient.On("ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*client.PlatformClientRequest")).Return(platformResp, nil)

	apiObject := &statement.StatementInputAPIModel{
		Statement:  "access_restricted_resource()",
		APIVersion: common.V2,
	}

	// When
	result, err := CallExternalAPIWithClient[*GenericErrorResponse](requestContext, mockClient, apiObject)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	errorMessage := err.Error()
	assert.Contains(t, errorMessage, "API response errors")
	assert.Contains(t, errorMessage, "PERMISSION_ERROR")
	assert.Contains(t, errorMessage, "Access denied")
	assert.Contains(t, errorMessage, "RESOURCE_ERROR")
	assert.Contains(t, errorMessage, "Resource not found")

	mockClient.AssertExpectations(t)
}

func TestCallExternalAPI_InvalidJSONResponse(t *testing.T) {
	// Given
	mockClient := new(MockPlatformClient)

	requestContext := common.NewRequestContext(context.Background())

	platformResp := &client.PlatformClientResponse{
		Response: &http.Response{StatusCode: 200},
		Body:     []byte("invalid json"),
	}
	mockClient.On("ExecuteRequest", mock.AnythingOfType("*common.RequestContext"), mock.AnythingOfType("*client.PlatformClientRequest")).Return(platformResp, nil)

	apiObject := &statement.StatementInputAPIModel{
		Statement:  "test_operation()",
		APIVersion: common.V2,
	}

	// When
	result, err := CallExternalAPIWithClient[*GenericSuccessResponse](requestContext, mockClient, apiObject)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid character")

	mockClient.AssertExpectations(t)
}

// Helper functions to create pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
