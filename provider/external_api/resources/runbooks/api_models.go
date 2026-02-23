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

package runbooks

import (
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"
)

// CommunicationDestination represents the communication settings for a runbook
type CommunicationDestination struct {
	Channel   string `json:"channel"`
	Workspace string `json:"workspace"`
}

// CommunicationFilters represents the communication notification filters
type CommunicationFilters struct {
	ExecutionNotifications bool `json:"execution_notifications"`
	CudNotifications       bool `json:"cud_notifications"`
	ApprovalNotifications  bool `json:"approval_notifications"`
}

// RunbookConfig represents the configuration settings for a runbook
type RunbookConfig struct {
	Params                   []customattribute.ParamJson         `json:"params"`
	Labels                   []string                            `json:"labels"`
	Approvers                []string                            `json:"approvers"`
	AllowedEntities          []string                            `json:"allowed_entities"`
	AllowedResourcesQuery    string                              `json:"allowed_resources_query"`
	Editors                  []string                            `json:"editors"`
	SecretNames              []string                            `json:"secret_names"`
	ExternalParams           []customattribute.ExternalParamJson `json:"external_params"`
	IsRunOutputPersisted     bool                                `json:"is_run_output_persisted"`
	FilterResourceToAction   bool                                `json:"filter_resource_to_action"`
	Cells                    []customattribute.CellJsonAPI       `json:"cells"`
	TimeoutMs                int64                               `json:"timeout_ms"`
	CommunicationDestination CommunicationDestination            `json:"communication_destination"`
	CommunicationFilters     CommunicationFilters                `json:"communication_filters"`
	Category                 string                              `json:"category"`
	ParamsGroups             ParamsGroups                        `json:"params_groups"`
}

type ParamsGroups struct {
	Required []string `json:"required" tfsdk:"required"`
	Optional []string `json:"optional" tfsdk:"optional"`
	Exported []string `json:"exported" tfsdk:"exported"`
	External []string `json:"external" tfsdk:"external"`
}

// RunbookEntityMetadata represents the metadata for a runbook entity
type RunbookEntityMetadata struct {
	Enabled     bool   `json:"enabled"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     int    `json:"version"`
	Description string `json:"description"`
}

// ConfigurationItem represents a single configuration item in the API response
type ConfigurationItem struct {
	Config         RunbookConfig         `json:"config"`
	EntityMetadata RunbookEntityMetadata `json:"entity_metadata"`
}

// RunbookConfigurations represents the configurations section of the API response
type RunbookConfigurations struct {
	Items []ConfigurationItem `json:"items"`
}

// RunbookOutput represents the output section of the API response
type RunbookOutput struct {
	Configurations RunbookConfigurations `json:"configurations"`
}

// Function represents a function in the summary
type Function struct {
	StartedAtMs  int64                  `json:"started_at_ms"`
	FinishedAtMs int64                  `json:"finished_at_ms"`
	EnvVarIn     map[string]interface{} `json:"env_var_in"`
}

// RunbookSummary represents the summary section of the API response
type RunbookSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// GetStatus returns the status from the summary
func (s RunbookSummary) GetStatus() string {
	return s.Status
}

// GetErrors returns the errors from the summary
func (s RunbookSummary) GetErrors() []apicommon.Error {
	return s.Errors
}

// RunbookRequest represents the request payload for runbook operations
type RunbookRequest struct {
	Statement string `json:"statement"`
}

// RunbookResponse represents the top-level response wrapper
type RunbookResponse struct {
	Data RunbookResponseAPIModel `json:"data"`
}

// RunbookResponseAPIModel represents the complete API response for runbook operations
type RunbookResponseAPIModel struct {
	Output  RunbookOutput  `json:"output"`
	Summary RunbookSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (a RunbookResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(a.Summary.Status, a.Summary.Errors)
}
