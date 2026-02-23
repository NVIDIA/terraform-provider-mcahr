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

package alarms

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// FireStep represents fire step configuration in API response
type FireStep struct {
	Description string `json:"description"`
	Title       string `json:"title"`
}

// ClearStep represents clear step configuration in API response
type ClearStep struct {
	Description string `json:"description"`
	Title       string `json:"title"`
}

// StepDetails represents the step details configuration
type StepDetails struct {
	FireStep  FireStep  `json:"fire_step"`
	ClearStep ClearStep `json:"clear_step"`
}

// AlarmConfig represents the configuration settings for an alarm
type AlarmConfig struct {
	FireQuery        string      `json:"fire_query"`
	ClearQuery       string      `json:"clear_query"`
	ResourceQuery    string      `json:"resource_query"`
	ResourceType     string      `json:"resource_type"`
	CheckIntervalSec int64       `json:"check_interval_sec"`
	StepDetails      StepDetails `json:"step_details"`
}

// AlarmEntityMetadata represents the metadata for an alarm entity
type AlarmEntityMetadata struct {
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Family      string `json:"family"`
}

// ConfigurationItem represents a single configuration item in the API response
type ConfigurationItem struct {
	Config         AlarmConfig         `json:"config"`
	EntityMetadata AlarmEntityMetadata `json:"entity_metadata"`
}

// AlarmConfigurations represents the configurations section of the API response
type AlarmConfigurations struct {
	Count int                 `json:"count"`
	Items []ConfigurationItem `json:"items"`
}

// AlarmOutput represents the output section of the API response
type AlarmOutput struct {
	Configurations AlarmConfigurations `json:"configurations"`
}

// AlarmSummary represents the summary section of the API response
type AlarmSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// AlarmResponseAPIModel represents the complete API response for alarm operations (V2)
type AlarmResponseAPIModel struct {
	Output  AlarmOutput  `json:"output"`
	Summary AlarmSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (a AlarmResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(a.Summary.Status, a.Summary.Errors)
}
