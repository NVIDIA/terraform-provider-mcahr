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

package principals

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// PrincipalClassV1 represents a principal class in V1 API responses
type PrincipalClassV1 struct {
	Enabled              bool   `json:"enabled"`
	Name                 string `json:"name"`
	Deleted              bool   `json:"deleted"`
	Identity             string `json:"identity"`
	ActionLimit          int    `json:"action_limit"`
	ExecuteLimit         int    `json:"execute_limit"`
	ViewLimit            int    `json:"view_limit"` // removed in release-29.0.0
	ConfigurePermission  int    `json:"configure_permission"`
	AdministerPermission int    `json:"administer_permission"`
	IDPName              string `json:"idp_name"`
}

// PrincipalContainerV1 represents the container for principal operations in V1 API
type PrincipalContainerV1 struct {
	PrincipalClasses []PrincipalClassV1 `json:"principal_classes"`
	Error            apicommon.ErrorV1  `json:"error"`
	Errors           []string           `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure PrincipalContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &PrincipalContainerV1{}

// GetNestedError returns the nested error structure
func (c *PrincipalContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *PrincipalContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// PrincipalResponseAPIModelV1 represents the complete API response for principal operations (V1)
type PrincipalResponseAPIModelV1 struct {
	GetPrincipalClass *PrincipalContainerV1     `json:"get_principal_class"`
	DefinePrincipal   *PrincipalContainerV1     `json:"define_principal"`
	UpdatePrincipal   *PrincipalContainerV1     `json:"update_principal"`
	DeletePrincipal   *PrincipalContainerV1     `json:"delete_principal"`
	Errors            *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// GetContainer returns the appropriate container based on the operation type
func (p PrincipalResponseAPIModelV1) GetContainer() *PrincipalContainerV1 {
	if p.GetPrincipalClass != nil {
		return p.GetPrincipalClass
	}
	if p.DefinePrincipal != nil {
		return p.DefinePrincipal
	}
	if p.UpdatePrincipal != nil {
		return p.UpdatePrincipal
	}
	if p.DeletePrincipal != nil {
		return p.DeletePrincipal
	}
	return nil
}

// GetErrors returns a string representation of the API response errors
func (p PrincipalResponseAPIModelV1) GetErrors() string {
	container := p.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(p.Errors, container)
}
