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

package integrations

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// IntegrationResponseAPIModelV1 represents the V1 API response format for /api/v1/execute endpoint
type IntegrationResponseAPIModelV1 struct {
	GetIntegrationClass *IntegrationContainerV1   `json:"get_integration_class,omitempty"`
	DefineIntegration   *IntegrationContainerV1   `json:"define_integration,omitempty"`
	UpdateIntegration   *IntegrationContainerV1   `json:"update_integration,omitempty"`
	DeleteIntegration   *IntegrationContainerV1   `json:"delete_integration,omitempty"`
	Errors              *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// IntegrationContainerV1 represents the integration container in V1 API responses
type IntegrationContainerV1 struct {
	IntegrationClasses []IntegrationClassV1 `json:"integration_classes,omitempty"`
	Error              apicommon.ErrorV1    `json:"error,omitempty"`
	Errors             []string             `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure IntegrationContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &IntegrationContainerV1{}

// GetNestedError returns the nested error structure
func (c *IntegrationContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *IntegrationContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// IntegrationClassV1 represents an integration class in V1 API responses
type IntegrationClassV1 struct {
	Enabled         bool   `json:"enabled,omitempty"`
	Name            string `json:"name,omitempty"`
	Params          string `json:"params,omitempty"` // JSON string containing integration-specific config
	ServiceName     string `json:"service_name,omitempty"`
	SerialNumber    string `json:"serial_number,omitempty"`
	PermissionsUser string `json:"permissions_user,omitempty"`
}

// GetContainer returns the appropriate integration container from the V1 API response
func (r IntegrationResponseAPIModelV1) GetContainer() *IntegrationContainerV1 {
	if r.GetIntegrationClass != nil {
		return r.GetIntegrationClass
	}
	if r.DefineIntegration != nil {
		return r.DefineIntegration
	}
	if r.UpdateIntegration != nil {
		return r.UpdateIntegration
	}
	if r.DeleteIntegration != nil {
		return r.DeleteIntegration
	}
	return nil
}

// GetErrors returns a formatted error string from the V1 API response
func (r IntegrationResponseAPIModelV1) GetErrors() string {
	container := r.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(r.Errors, container)
}
