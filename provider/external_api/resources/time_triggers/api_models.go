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

package time_triggers

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// TimeTriggerConfig represents the configuration settings for a time trigger
type TimeTriggerConfig struct {
	FireQuery string `json:"fire_query"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// TimeTriggerEntityMetadata represents the metadata for a time trigger entity
type TimeTriggerEntityMetadata struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

// ConfigurationItem represents a single configuration item in the API response
type ConfigurationItem struct {
	Config         TimeTriggerConfig         `json:"config"`
	EntityMetadata TimeTriggerEntityMetadata `json:"entity_metadata"`
}

// TimeTriggerConfigurations represents the configurations section of the API response
type TimeTriggerConfigurations struct {
	Items []ConfigurationItem `json:"items"`
}

// TimeTriggerOutput represents the output section of the API response
type TimeTriggerOutput struct {
	Configurations TimeTriggerConfigurations `json:"configurations"`
}

// TimeTriggerSummary represents the summary section of the API response
type TimeTriggerSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// TimeTriggerResponseAPIModel represents the complete API response for time trigger operations (V2)
type TimeTriggerResponseAPIModel struct {
	Output  TimeTriggerOutput  `json:"output"`
	Summary TimeTriggerSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (t TimeTriggerResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(t.Summary.Status, t.Summary.Errors)
}
