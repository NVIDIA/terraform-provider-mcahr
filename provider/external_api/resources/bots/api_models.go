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

package bots

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// BotConfig represents the configuration settings for a bot
type BotConfig struct {
	TriggerEntityName   string                      `json:"trigger_entity_name"`
	ExecutionEntityName string                      `json:"execution_entity_name"`
	CommunicationDest   BotCommunicationDestination `json:"communication_destination"`
	TriggerID           string                      `json:"trigger_id"`
	TriggerSource       string                      `json:"trigger_source"`
	AlarmResourceQuery  string                      `json:"alarm_resource_query"`
	IntegrationName     string                      `json:"integration_name"`
}

// BotCommunicationDestination represents the communication settings for a bot
type BotCommunicationDestination struct {
	Channel   string `json:"channel"`
	Workspace string `json:"workspace"`
}

// BotEntityMetadata represents the metadata for a bot entity
type BotEntityMetadata struct {
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	Family      string `json:"family"`
	Description string `json:"description"`
}

// ConfigurationItem represents a single configuration item in the API response
type ConfigurationItem struct {
	Config         BotConfig         `json:"config"`
	EntityMetadata BotEntityMetadata `json:"entity_metadata"`
}

// BotConfigurations represents the configurations section of the API response
type BotConfigurations struct {
	Items []ConfigurationItem `json:"items"`
}

// BotOutput represents the output section of the API response
type BotOutput struct {
	Configurations BotConfigurations `json:"configurations"`
}

// BotSummary represents the summary section of the API response
type BotSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// BotResponseAPIModel represents the complete API response for bot operations (V2)
type BotResponseAPIModel struct {
	Output  BotOutput  `json:"output"`
	Summary BotSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (b BotResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(b.Summary.Status, b.Summary.Errors)
}
