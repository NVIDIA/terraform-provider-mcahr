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

package secrets

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// NVaultSecretResponseAPIModelV1 represents the response structure for V1 nvault secret API calls
type NVaultSecretResponseAPIModelV1 struct {
	DefineSecret *NVaultSecretContainerV1  `json:"define_secret,omitempty"`
	GetSecret    *NVaultSecretContainerV1  `json:"get_secret,omitempty"`
	UpdateSecret *NVaultSecretContainerV1  `json:"update_secret,omitempty"`
	DeleteSecret *NVaultSecretContainerV1  `json:"delete_secret,omitempty"`
	Errors       *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// NVaultSecretContainerV1 represents the nvault secret container in V1 API responses
type NVaultSecretContainerV1 struct {
	Secrets []NVaultSecretV1  `json:"secrets,omitempty"`
	Error   apicommon.ErrorV1 `json:"error,omitempty"`
	Errors  []string          `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure NVaultSecretContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &NVaultSecretContainerV1{}

// GetNestedError returns the nested error structure
func (c *NVaultSecretContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *NVaultSecretContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// NVaultSecretV1 represents a nvault secret in V1 API responses
type NVaultSecretV1 struct {
	Name       string             `json:"name,omitempty"`
	SecretInfo NVaultSecretInfoV1 `json:"secret_info,omitempty"`
}

// NVaultSecretInfoV1 represents nvault secret information in V1 API responses
type NVaultSecretInfoV1 struct {
	IntegrationName string `json:"integration_name,omitempty"`
	VaultSecretPath string `json:"vault_secret_path,omitempty"`
	VaultSecretKey  string `json:"vault_secret_key,omitempty"`
}

// GetContainer returns the appropriate nvault secret container from the V1 API response
func (r NVaultSecretResponseAPIModelV1) GetContainer() *NVaultSecretContainerV1 {
	if r.DefineSecret != nil {
		return r.DefineSecret
	}
	if r.GetSecret != nil {
		return r.GetSecret
	}
	if r.UpdateSecret != nil {
		return r.UpdateSecret
	}
	if r.DeleteSecret != nil {
		return r.DeleteSecret
	}
	return nil
}

// GetErrors returns a formatted error string from the V1 API response
func (r NVaultSecretResponseAPIModelV1) GetErrors() string {
	container := r.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(r.Errors, container)
}
