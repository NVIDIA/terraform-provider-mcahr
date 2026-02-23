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

package actions

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// Step represents a single step configuration
type Step struct {
	Description string `json:"description"`
	Title       string `json:"title"`
}

// StepDetails represents the step details configuration for an action
type StepDetails struct {
	StartStep    Step `json:"start_step"`
	ErrorStep    Step `json:"error_step"`
	CompleteStep Step `json:"complete_step"`
}

// CommunicationDestination represents the communication settings for an action
type CommunicationDestination struct {
	Channel   string `json:"channel"`
	Workspace string `json:"workspace"`
}

// ActionConfig represents the configuration settings for an action
type ActionConfig struct {
	Timeout               int64                    `json:"timeout"`
	Shell                 string                   `json:"shell"`
	Params                string                   `json:"params"`
	ResourceQuery         string                   `json:"resource_query"`
	ResourceTagsToExport  string                   `json:"resource_tags_to_export"`
	ResEnvVar             string                   `json:"res_env_var"`
	FileDeps              string                   `json:"file_deps"`
	AllowedEntities       []string                 `json:"allowed_entities"`
	AllowedResourcesQuery string                   `json:"allowed_resources_query"`
	Editors               []string                 `json:"editors"`
	CommandText           string                   `json:"command_text"`
	CommunicationDest     CommunicationDestination `json:"communication_destination"`
	StepDetails           StepDetails              `json:"step_details"`
}

// ActionEntityMetadata represents the metadata for an action entity
type ActionEntityMetadata struct {
	Enabled     bool   `json:"enabled"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     int    `json:"version"`
	Description string `json:"description"`
}

// ConfigurationItem represents a single configuration item in the API response
type ConfigurationItem struct {
	Config         ActionConfig         `json:"config"`
	EntityMetadata ActionEntityMetadata `json:"entity_metadata"`
}

// ActionConfigurations represents the configurations section of the API response
type ActionConfigurations struct {
	Count int                 `json:"count"`
	Items []ConfigurationItem `json:"items"`
}

// ActionOutput represents the output section of the API response
type ActionOutput struct {
	Configurations ActionConfigurations `json:"configurations"`
}

// ActionSummary represents the summary section of the API response
type ActionSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// ActionRequest represents the request payload for action operations
type ActionRequest struct {
	Statement string `json:"statement"`
}

// ActionResponse represents the top-level response wrapper
type ActionResponse struct {
	Data ActionResponseAPIModel `json:"data"`
}

// ActionResponseAPIModel represents the complete API response for action operations (V2)
type ActionResponseAPIModel struct {
	Output  ActionOutput  `json:"output"`
	Summary ActionSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (a ActionResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(a.Summary.Status, a.Summary.Errors)
}
