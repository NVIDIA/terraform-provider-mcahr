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

// ResourceAttributes represents the attributes nested in the configuration item
type ResourceAttributes struct {
	Params string `json:"params,omitempty"` // JSON string like "[\"param1\",\"param2\"]"
}

// ResourceConfigurationItem represents a single configuration item in V2 API
type ResourceConfigurationItem struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Formula     string             `json:"formula"`
	Attributes  ResourceAttributes `json:"attributes"`
}

// ResourceConfigurations represents the configurations structure in V2 API
type ResourceConfigurations struct {
	Items []ResourceConfigurationItem `json:"items"`
}

// ResourceOutput represents the output structure in V2 API responses
type ResourceOutput struct {
	Symbols ResourceConfigurations `json:"symbols"`
}

// ResourceSummary represents the summary structure in V2 API responses
type ResourceSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// ResourceResponseAPIModel represents the V2 API response structure
type ResourceResponseAPIModel struct {
	Output  ResourceOutput  `json:"output"`
	Summary ResourceSummary `json:"summary"`
}

// GetErrors returns a string representation of the V2 API response errors
func (r ResourceResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(r.Summary.Status, r.Summary.Errors)
}
