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

// ActionResponseAPIModelV1 represents the V1 API response format for /api/v1/execute endpoint
// Unified model that handles define_action, update_action, get_action_class, and delete_action responses
type ActionResponseAPIModelV1 struct {
	DefineAction   *ActionContainerV1        `json:"define_action,omitempty"`
	UpdateAction   *ActionContainerV1        `json:"update_action,omitempty"`
	GetActionClass *ActionContainerV1        `json:"get_action_class,omitempty"`
	DeleteAction   *ActionContainerV1        `json:"delete_action,omitempty"`
	Errors         *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// ActionContainerV1 represents the common structure for all V1 action operation containers
type ActionContainerV1 struct {
	ActionClasses []ActionClassV1   `json:"action_classes"`
	Error         apicommon.ErrorV1 `json:"error"`
	Errors        []string          `json:"errors,omitempty"` // Direct errors array for some error cases

}

// Ensure ActionContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &ActionContainerV1{}

// GetNestedError returns the nested error structure
func (c *ActionContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array (none for actions)
func (c *ActionContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// ActionClassV1 represents an action class in V1 response
// Simplified to only include fields needed by the translator
type ActionClassV1 struct {
	Timeout               int64           `json:"timeout"`
	Command               string          `json:"command"`
	Enabled               bool            `json:"enabled"`
	Name                  string          `json:"name"`
	Description           string          `json:"description"`
	Shell                 string          `json:"shell"`
	Params                string          `json:"params"`
	ResourceQuery         string          `json:"resource_query"`
	ResourceTagsToExport  string          `json:"resource_tags_to_export"`
	StartStepClass        StepClassV1     `json:"start_step_class"`
	CompleteStepClass     StepClassV1     `json:"complete_step_class"`
	ErrorStepClass        StepClassV1     `json:"error_step_class"`
	ResEnvVar             string          `json:"res_env_var"`
	FileDeps              string          `json:"file_deps"`
	AllowedEntities       []string        `json:"allowed_entities"`
	AllowedResourcesQuery string          `json:"allowed_resources_query"`
	Communication         CommunicationV1 `json:"communication"`
	Editors               []string        `json:"editors"`
}

// StepClassV1 represents a step class in V1 response
// Simplified to only include fields needed by the translator
type StepClassV1 struct {
	TitleTemplate string `json:"title_template"`
	ShortTemplate string `json:"short_template"`
}

// CommunicationV1 represents communication settings in V1 response
type CommunicationV1 struct {
	Channel   string `json:"channel"`
	Workspace string `json:"workspace"`
}

// GetContainer returns the action container regardless of operation type
func (a ActionResponseAPIModelV1) GetContainer() *ActionContainerV1 {
	if a.DefineAction != nil {
		return a.DefineAction
	}
	if a.UpdateAction != nil {
		return a.UpdateAction
	}
	if a.GetActionClass != nil {
		return a.GetActionClass
	}
	if a.DeleteAction != nil {
		return a.DeleteAction
	}
	return nil
}

// GetErrors returns a string representation of the V1 API response errors
func (a ActionResponseAPIModelV1) GetErrors() string {
	container := a.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(a.Errors, container)
}
