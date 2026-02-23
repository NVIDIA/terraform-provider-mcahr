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

package dashboards

import (
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
)

// DashboardClassV1 represents a dashboard class in V1 API
type DashboardClassV1 struct {
	Name          string `json:"name"`
	Configuration string `json:"configuration"`
	DashboardType string `json:"dashboard_type"`
}

// DashboardContainerV1 represents the container for dashboard operations in V1 API
type DashboardContainerV1 struct {
	DashboardClasses []DashboardClassV1 `json:"dashboard_classes"`
	Error            *apicommon.ErrorV1 `json:"error"`
	Errors           []string           `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure DashboardContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &DashboardContainerV1{}

// GetNestedError returns the nested error structure
func (c *DashboardContainerV1) GetNestedError() apicommon.ErrorV1 {
	return *c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *DashboardContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// DashboardResponseAPIModelV1 represents the V1 API response structure
type DashboardResponseAPIModelV1 struct {
	DefineDashboard   *DashboardContainerV1     `json:"define_dashboard,omitempty"`
	UpdateDashboard   *DashboardContainerV1     `json:"update_dashboard,omitempty"`
	GetDashboardClass *DashboardContainerV1     `json:"get_dashboard_class,omitempty"`
	DeleteDashboard   *DashboardContainerV1     `json:"delete_dashboard,omitempty"`
	Errors            *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// GetContainer returns the appropriate container based on the operation type
func (r *DashboardResponseAPIModelV1) GetContainer() *DashboardContainerV1 {
	if r.DefineDashboard != nil {
		return r.DefineDashboard
	}
	if r.UpdateDashboard != nil {
		return r.UpdateDashboard
	}
	if r.GetDashboardClass != nil {
		return r.GetDashboardClass
	}
	if r.DeleteDashboard != nil {
		return r.DeleteDashboard
	}
	return nil
}

// GetErrors returns a string representation of the V1 API response errors
func (r DashboardResponseAPIModelV1) GetErrors() string {
	container := r.GetContainer()
	if container == nil {
		return "No container found in API response"
	}
	return apicommon.FormatV1ErrorsFromContainer(r.Errors, container)
}
