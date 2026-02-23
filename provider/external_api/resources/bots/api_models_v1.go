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

// BotResponseAPIModelV1 represents the response structure for V1 bot API calls
type BotResponseAPIModelV1 struct {
	GetBotClass *BotContainerV1           `json:"get_bot_class,omitempty"`
	DefineBot   *BotContainerV1           `json:"define_bot,omitempty"`
	UpdateBot   *BotContainerV1           `json:"update_bot,omitempty"`
	DeleteBot   *BotContainerV1           `json:"delete_bot,omitempty"`
	Errors      *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// BotContainerV1 represents the bot container in V1 API responses
type BotContainerV1 struct {
	BotClasses []BotClassV1      `json:"bot_classes,omitempty"`
	Error      apicommon.ErrorV1 `json:"error,omitempty"`
	Errors     []string          `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure BotContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &BotContainerV1{}

// GetNestedError returns the nested error structure
func (c *BotContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *BotContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// BotClassV1 represents a bot class in V1 API responses
type BotClassV1 struct {
	Name               string          `json:"name,omitempty"`
	Description        string          `json:"description,omitempty"`
	Enabled            bool            `json:"enabled,omitempty"`
	ConfigData         ConfigDataV1    `json:"config_data"`
	Communication      CommunicationV1 `json:"communication"`
	AlarmStatement     string          `json:"alarm_statement,omitempty"`
	ActionStatement    string          `json:"action_statement,omitempty"`
	TriggerSource      string          `json:"trigger_source,omitempty"`
	ExternalTriggerID  string          `json:"external_trigger_id,omitempty"`
	AlarmResourceQuery string          `json:"alarm_resource_query,omitempty"`
	EventType          string          `json:"event_type,omitempty"`
	MonitorID          string          `json:"monitor_id,omitempty"`
	IntegrationName    string          `json:"integration_name,omitempty"`
}

// ConfigDataV1 represents config data structure in V1 API responses
type ConfigDataV1 struct {
	Family string `json:"family,omitempty"`
}

// CommunicationV1 represents communication settings in V1 response
type CommunicationV1 struct {
	Channel   string `json:"channel"`
	Workspace string `json:"workspace"`
}

// GetContainer returns the appropriate bot container from the V1 API response
func (r BotResponseAPIModelV1) GetContainer() *BotContainerV1 {
	if r.GetBotClass != nil {
		return r.GetBotClass
	}
	if r.DefineBot != nil {
		return r.DefineBot
	}
	if r.UpdateBot != nil {
		return r.UpdateBot
	}
	if r.DeleteBot != nil {
		return r.DeleteBot
	}
	return nil
}

// GetErrors returns a formatted error string from the V1 API response
func (r BotResponseAPIModelV1) GetErrors() string {
	container := r.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(r.Errors, container)
}
