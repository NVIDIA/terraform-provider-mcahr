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

// IntegrationItem represents a single integration item in the V2 API response
type IntegrationItem struct {
	Name            string                 `json:"name"`
	Enabled         bool                   `json:"enabled"`
	SerialNumber    string                 `json:"serial_number"`
	PermissionsUser string                 `json:"permissions_user"`
	IntegrationData map[string]interface{} `json:"integration_data"`
	IntegrationType string                 `json:"integration_type"`
}

// IntegrationConfigurations represents the integrations section of the API response
type IntegrationConfigurations struct {
	Items []IntegrationItem `json:"items"`
}

// IntegrationOutput represents the output section of the API response
type IntegrationOutput struct {
	Integrations IntegrationConfigurations `json:"integrations"`
}

// IntegrationSummary represents the summary section of the API response
type IntegrationSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// IntegrationResponseAPIModel represents the complete API response for integration operations (V2)
type IntegrationResponseAPIModel struct {
	Output  IntegrationOutput  `json:"output"`
	Summary IntegrationSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (i IntegrationResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(i.Summary.Status, i.Summary.Errors)
}
