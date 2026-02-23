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

package backend_version

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPlatformClient implements the external_api.PlatformClientInterface for testing
type MockPlatformClient struct {
	mock.Mock
}

func (m *MockPlatformClient) ExecuteRequest(requestContext *common.RequestContext, request *client.PlatformClientRequest) (*client.PlatformClientResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*client.PlatformClientResponse), args.Error(1)
}

func TestFetchBackendVersion_SuccessWithMockedAPI(t *testing.T) {
	// Given
	mockClient := &MockPlatformClient{}
	provider := NewBackendVersionProvider(mockClient)

	// Mock successful API response (V1 format)
	expectedResponse := `{
		"backend_version": "release-30.2.1"
	}`
	mockResponse := &client.PlatformClientResponse{
		Body: []byte(expectedResponse),
	}
	mockClient.On("ExecuteRequest", mock.Anything).Return(mockResponse, nil)

	// When
	result := provider.FetchBackendVersion(common.NewRequestContext(context.Background()))

	// Then
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Version)
	assert.Equal(t, "release-30.2.1", result.Version.Version)
	assert.Equal(t, int64(30), result.Version.Major)
	assert.Equal(t, int64(2), result.Version.Minor)
	assert.Equal(t, int64(1), result.Version.Patch)
	assert.False(t, result.IsParseError)
	mockClient.AssertExpectations(t)
}

func TestFetchBackendVersion_ParseError(t *testing.T) {
	// Given
	mockClient := &MockPlatformClient{}
	provider := NewBackendVersionProvider(mockClient)

	// Mock successful API response but with invalid version (V1 format)
	// Use "release-" prefix but no valid numbers to trigger parse error
	expectedResponse := `{
		"backend_version": "release-no-numbers"
	}`
	mockResponse := &client.PlatformClientResponse{
		Body: []byte(expectedResponse),
	}
	mockClient.On("ExecuteRequest", mock.Anything).Return(mockResponse, nil)

	// When
	result := provider.FetchBackendVersion(common.NewRequestContext(context.Background()))
	// Then
	assert.NotNil(t, result.Error)
	assert.True(t, result.IsParseError)
	assert.Nil(t, result.Version)
	assert.Equal(t, "failed to parse backend version: release-no-numbers", result.Error.Error())
	mockClient.AssertExpectations(t)
}

func TestFetchBackendVersionWithFallback_ErrorWhenFallbackInvalid(t *testing.T) {
	// Given
	mockClient := &MockPlatformClient{}
	provider := NewBackendVersionProvider(mockClient)

	// Mock API response with invalid version format
	expectedResponse := `{
		"backend_version": "release-no-numbers"
	}`
	mockResponse := &client.PlatformClientResponse{
		Body: []byte(expectedResponse),
	}
	mockClient.On("ExecuteRequest", mock.Anything).Return(mockResponse, nil)

	// When
	result, err := provider.FetchBackendVersionWithFallback(common.NewRequestContext(context.Background()), "release-also-invalid")

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "both API response and fallback version are invalid: failed to parse backend version: release-no-numbers", err.Error())
	mockClient.AssertExpectations(t)
}

func TestFetchBackendVersion_SyntaxError(t *testing.T) {
	// Given
	mockClient := &MockPlatformClient{}
	provider := NewBackendVersionProvider(mockClient)

	// Mock API response with syntax error (V1 format)
	expectedResponse := `{
		"errors": {
			"$": [
				"Undefined symbol backend_versiond"
			]
		}
	}`
	mockResponse := &client.PlatformClientResponse{
		Body: []byte(expectedResponse),
	}
	mockClient.On("ExecuteRequest", mock.Anything).Return(mockResponse, nil)

	// When
	result := provider.FetchBackendVersion(common.NewRequestContext(context.Background()))

	// Then
	assert.NotNil(t, result.Error)
	assert.Nil(t, result.Version)
	assert.False(t, result.IsParseError)
	assert.Equal(t, "failed to fetch backend version: API response errors: Syntax Errors: Undefined symbol backend_versiond", result.Error.Error())
	mockClient.AssertExpectations(t)
}

func TestFetchBackendVersionWithFallback_DoesNotUseFallbackOnSyntaxError(t *testing.T) {
	// Given
	mockClient := &MockPlatformClient{}
	provider := NewBackendVersionProvider(mockClient)

	// Mock API response with syntax error (NOT a parse error)
	expectedResponse := `{
		"errors": {
			"$": [
				"Undefined symbol backend_versiond"
			]
		}
	}`
	mockResponse := &client.PlatformClientResponse{
		Body: []byte(expectedResponse),
	}
	mockClient.On("ExecuteRequest", mock.Anything).Return(mockResponse, nil)

	// When
	result, err := provider.FetchBackendVersionWithFallback(common.NewRequestContext(context.Background()), "release-29.0.0")

	// Then
	// Should return the API error and NOT use the fallback
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to fetch backend version: API response errors: Syntax Errors: Undefined symbol backend_versiond", err.Error())
	mockClient.AssertExpectations(t)
}

func TestFetchBackendVersionWithFallback_DoesNotUseFallbackOnSuccess(t *testing.T) {
	// Given
	mockClient := &MockPlatformClient{}
	provider := NewBackendVersionProvider(mockClient)

	// Mock successful API response
	expectedResponse := `{
		"backend_version": "release-31.0.5"
	}`
	mockResponse := &client.PlatformClientResponse{
		Body: []byte(expectedResponse),
	}
	mockClient.On("ExecuteRequest", mock.Anything).Return(mockResponse, nil)

	// When
	result, err := provider.FetchBackendVersionWithFallback(common.NewRequestContext(context.Background()), "release-29.0.0")

	// Then
	// Should return the API response and NOT use the fallback
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Verify it used the API response, not the fallback
	assert.Equal(t, "release-31.0.5", result.Version)
	assert.Equal(t, int64(31), result.Major)
	assert.Equal(t, int64(0), result.Minor)
	assert.Equal(t, int64(5), result.Patch)
	mockClient.AssertExpectations(t)
}

func TestFetchBackendVersionWithFallback_UsesFallbackOnlyForParseError(t *testing.T) {
	// Given
	mockClient := &MockPlatformClient{}
	provider := NewBackendVersionProvider(mockClient)

	// Mock API response with a version string that fails to parse (IsParseError = true)
	// Use "release-" prefix but no valid numbers to trigger parse error
	expectedResponse := `{
		"backend_version": "release-invalid-numbers"
	}`
	mockResponse := &client.PlatformClientResponse{
		Body: []byte(expectedResponse),
	}
	mockClient.On("ExecuteRequest", mock.Anything).Return(mockResponse, nil)

	// When
	result, err := provider.FetchBackendVersionWithFallback(common.NewRequestContext(context.Background()), "release-30.1.2")

	// Then
	// Should use the fallback because it's a parse error
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "release-30.1.2", result.Version)
	assert.Equal(t, int64(30), result.Major)
	assert.Equal(t, int64(1), result.Minor)
	assert.Equal(t, int64(2), result.Patch)
	mockClient.AssertExpectations(t)
}
