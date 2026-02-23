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

package common

import (
	"context"
	"terraform/terraform-provider/provider/common/version"
)

// RequestContext wraps a context.Context with additional generic data
// that can be passed through the entire request pipeline
type RequestContext struct {
	// Context is the underlying Go context for cancellation, deadlines, and values
	Context context.Context

	// ResourceType indicates which resource type initiated this request
	ResourceType string

	// Operation indicates what CRUD operation is being performed
	Operation CrudOperation

	// BackendVersion indicates the version of the platform backend
	BackendVersion *version.BackendVersion

	// APIVersion indicates the version of the Execute API
	APIVersion APIVersion
}

// NewRequestContext creates a new RequestContext with the given context
func NewRequestContext(ctx context.Context) *RequestContext {
	return &RequestContext{
		Context: ctx,
	}
}

// WithResourceType sets the resource type and returns the RequestContext for chaining
func (rc *RequestContext) WithResourceType(resourceType string) *RequestContext {
	rc.ResourceType = resourceType
	return rc
}

// WithOperation sets the operation and returns the RequestContext for chaining
func (rc *RequestContext) WithOperation(operation CrudOperation) *RequestContext {
	rc.Operation = operation
	return rc
}

// WithBackendVersion sets the backend version and returns the RequestContext for chaining
func (rc *RequestContext) WithBackendVersion(backendVersion *version.BackendVersion) *RequestContext {
	rc.BackendVersion = backendVersion
	return rc
}

// WithAPIVersion sets the API version and returns the RequestContext for chaining
func (rc *RequestContext) WithAPIVersion(apiVersion APIVersion) *RequestContext {
	rc.APIVersion = apiVersion
	return rc
}
