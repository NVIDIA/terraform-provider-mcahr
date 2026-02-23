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
	customattribute "terraform/terraform-provider/provider/external_api/resources/dashboards/custom_attribute"
)

// DashboardConfig represents the configuration settings for a dashboard
type DashboardConfig struct {
	DashboardType string                            `json:"dashboard_type"`
	Configuration customattribute.ConfigurationJson `json:"configuration"`
}

// DashboardEntityMetadata represents the metadata for dashboard entity in V2 API
type DashboardEntityMetadata struct {
	Name string `json:"name"`
}

// ConfigurationItem represents a single configuration item in V2 API
type ConfigurationItem struct {
	Config         DashboardConfig         `json:"config"`
	EntityMetadata DashboardEntityMetadata `json:"entity_metadata"`
}

// DashboardConfigurations represents the configurations structure in V2 API
type DashboardConfigurations struct {
	Items []ConfigurationItem `json:"items"`
}

// DashboardOutput represents the output structure in V2 API responses
type DashboardOutput struct {
	Configurations DashboardConfigurations `json:"configurations"`
}

// DashboardSummary represents the summary structure in V2 API responses
type DashboardSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// DashboardResponseAPIModel represents the V2 API response structure
type DashboardResponseAPIModel struct {
	Output  DashboardOutput  `json:"output"`
	Summary DashboardSummary `json:"summary"`
}

// GetErrors returns a string representation of the V2 API response errors
func (r DashboardResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(r.Summary.Status, r.Summary.Errors)
}
