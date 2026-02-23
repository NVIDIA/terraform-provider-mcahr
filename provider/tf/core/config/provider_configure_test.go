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

package config

import (
	"testing"

	"terraform/terraform-provider/provider/external_api/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfiguredProviderData_Success(t *testing.T) {
	t.Parallel()

	// given
	expectedClient := &client.PlatformClient{}
	expectedProviderData := &FrameworkProviderData{
		Client: expectedClient,
	}

	req := resource.ConfigureRequest{
		ProviderData: expectedProviderData,
	}
	resp := &resource.ConfigureResponse{}

	// when
	result := ReadConfiguredProviderData(req, resp)

	// then
	require.NotNil(t, result, "expected non-nil result")
	assert.Same(t, expectedProviderData, result, "expected same provider data instance")
	assert.Same(t, expectedClient, result.Client, "expected same client instance")
	assert.False(t, resp.Diagnostics.HasError(), "expected no diagnostics errors")
}

func TestReadConfiguredProviderData_WrongType(t *testing.T) {
	t.Parallel()

	// given
	wrongTypeData := "not a framework provider data"
	req := resource.ConfigureRequest{
		ProviderData: wrongTypeData,
	}
	resp := &resource.ConfigureResponse{}

	// when
	result := ReadConfiguredProviderData(req, resp)

	// then
	assert.Nil(t, result, "expected nil result")
	assert.True(t, resp.Diagnostics.HasError(), "expected diagnostics error")

	errors := resp.Diagnostics.Errors()
	require.Len(t, errors, 1, "expected exactly 1 error")

	assert.Equal(t, "Unexpected Resource Configure Type", errors[0].Summary(), "error summary mismatch")
	assert.NotEmpty(t, errors[0].Detail(), "expected non-empty error detail")
	assert.Contains(t, errors[0].Detail(), "string", "expected error detail to mention actual type")
}

func TestReadConfiguredProviderData_NilProviderData(t *testing.T) {
	t.Parallel()

	// given
	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}

	// when
	result := ReadConfiguredProviderData(req, resp)

	// then
	assert.Nil(t, result, "expected nil result")
	assert.True(t, resp.Diagnostics.HasError(), "expected diagnostics error")

	errors := resp.Diagnostics.Errors()
	require.Len(t, errors, 1, "expected exactly 1 error")
	assert.Equal(t, "Unexpected Resource Configure Type", errors[0].Summary(), "error summary mismatch")
}

func TestEnsureClientConfigured_Success(t *testing.T) {
	t.Parallel()

	// given
	client := &client.PlatformClient{}
	diags := &diag.Diagnostics{}

	// when
	result := EnsureClientConfigured(client, diags)

	// then
	assert.True(t, result, "expected true result")
	assert.False(t, diags.HasError(), "expected no diagnostics errors")
}

func TestEnsureClientConfigured_NilClient(t *testing.T) {
	t.Parallel()

	// given
	var client *client.PlatformClient = nil
	diags := &diag.Diagnostics{}

	// when
	result := EnsureClientConfigured(client, diags)

	// then
	assert.False(t, result, "expected false result")
	assert.True(t, diags.HasError(), "expected diagnostics error")

	errors := diags.Errors()
	require.Len(t, errors, 1, "expected exactly 1 error")

	assert.Equal(t, "Unconfigured Client", errors[0].Summary(), "error summary mismatch")
	assert.Equal(t, "Expected configured PlatformClient. Please report this issue to the provider developers.", errors[0].Detail(), "error detail mismatch")
}
