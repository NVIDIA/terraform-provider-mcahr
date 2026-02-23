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

package resources

import (
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
)

// ResourceAttributesV1 represents the attributes nested in the symbol
type ResourceAttributesV1 struct {
	Description string `json:"description"`
	Params      string `json:"params"` // JSON string like "[\"param1\",\"param2\"]"
	Name        string `json:"name,omitempty"`
	Value       string `json:"value,omitempty"`
}

// ResourceSymbolV1 represents a single resource symbol in V1 API
type ResourceSymbolV1 struct {
	Name       string               `json:"name"`
	Formula    string               `json:"formula"`
	Attributes ResourceAttributesV1 `json:"attributes"`
}

// ResourceContainerV1Single represents the container for create/update/delete operations
type ResourceContainerV1Single struct {
	Symbol ResourceSymbolV1  `json:"symbol"`
	Error  apicommon.ErrorV1 `json:"error"`
	Errors []string          `json:"errors,omitempty"`
}

// ResourceContainerV1List represents the container for list operations
type ResourceContainerV1List struct {
	Symbol []ResourceSymbolV1 `json:"symbol"`
	Error  apicommon.ErrorV1  `json:"error"`
	Errors []string           `json:"errors,omitempty"`
}

// ResourceResponseAPIModelV1 represents the V1 API response structure
type ResourceResponseAPIModelV1 struct {
	DefineResource *ResourceContainerV1Single `json:"define_resource,omitempty"`
	UpdateResource *ResourceContainerV1Single `json:"update_resource,omitempty"`
	ListType       *ResourceContainerV1List   `json:"list_type,omitempty"`
	DeleteResource *ResourceContainerV1Single `json:"delete_resource,omitempty"`
	Errors         *apicommon.SyntaxErrorsV1  `json:"errors,omitempty"`
}

// GetContainer returns the appropriate container as ResourceContainerV1List
// For single-object responses, it wraps the symbol in an array
func (r ResourceResponseAPIModelV1) GetContainer() *ResourceContainerV1List {
	if r.DefineResource != nil {
		return &ResourceContainerV1List{
			Symbol: []ResourceSymbolV1{r.DefineResource.Symbol},
			Error:  r.DefineResource.Error,
			Errors: r.DefineResource.Errors,
		}
	}
	if r.UpdateResource != nil {
		return &ResourceContainerV1List{
			Symbol: []ResourceSymbolV1{r.UpdateResource.Symbol},
			Error:  r.UpdateResource.Error,
			Errors: r.UpdateResource.Errors,
		}
	}
	if r.ListType != nil {
		return r.ListType
	}
	if r.DeleteResource != nil {
		return &ResourceContainerV1List{
			Symbol: []ResourceSymbolV1{r.DeleteResource.Symbol},
			Error:  r.DeleteResource.Error,
			Errors: r.DeleteResource.Errors,
		}
	}
	return nil
}

// GetErrors returns a string representation of the V1 API response errors
func (r ResourceResponseAPIModelV1) GetErrors() string {
	container := r.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(r.Errors, container)
}

// Ensure containers implement V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &ResourceContainerV1Single{}
var _ apicommon.V1ErrorContainer = &ResourceContainerV1List{}

// GetNestedError returns the nested error structure for single container
func (c *ResourceContainerV1Single) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array for single container
func (c *ResourceContainerV1Single) GetDirectErrors() []string {
	return c.Errors
}

// GetNestedError returns the nested error structure for list container
func (c *ResourceContainerV1List) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array for list container
func (c *ResourceContainerV1List) GetDirectErrors() []string {
	return c.Errors
}
