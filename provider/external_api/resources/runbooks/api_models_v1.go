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

// RunbookResponseAPIModelV1 represents the V1 API response format for /api/v1/execute endpoint
// Unified model that handles define_notebook, update_notebook, get_notebook_class, and delete_notebook responses
type RunbookResponseAPIModelV1 struct {
	DefineNotebook   *RunbookContainerV1       `json:"define_notebook,omitempty"`
	UpdateNotebook   *RunbookContainerV1       `json:"update_notebook,omitempty"`
	GetNotebookClass *RunbookContainerV1       `json:"get_notebook_class,omitempty"`
	DeleteNotebook   *RunbookContainerV1       `json:"delete_notebook,omitempty"`
	Errors           *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// RunbookContainerV1 represents the common structure for all V1 runbook operation containers
type RunbookContainerV1 struct {
	NotebookClasses []NotebookClassV1 `json:"notebook_classes"`
	Error           apicommon.ErrorV1 `json:"error"`
	Errors          []string          `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure RunbookContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &RunbookContainerV1{}

// GetNestedError returns the nested error structure
func (c *RunbookContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *RunbookContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// NotebookClassV1 represents a notebook class (runbook) in V1 response
type NotebookClassV1 struct {
	Enabled                             bool                                `json:"enabled"`
	Name                                string                              `json:"name"`
	Description                         string                              `json:"description"`
	Params                              []customattribute.ParamJson         `json:"params"`
	Labels                              []string                            `json:"labels"`
	Cells                               []customattribute.CellJsonAPI       `json:"cells"`
	Editors                             []string                            `json:"editors"`
	AllowedEntities                     []string                            `json:"allowed_entities"`
	AllowedResourcesQuery               string                              `json:"allowed_resources_query"`
	ExternalParams                      []customattribute.ExternalParamJson `json:"external_params"`
	Approvers                           []string                            `json:"approvers"`
	TimeoutMs                           int64                               `json:"timeout_ms"`
	IsRunOutputPersisted                bool                                `json:"is_run_output_persisted"`
	FilterResourceToAction              bool                                `json:"filter_resource_to_action"`
	SecretNames                         []string                            `json:"secret_names"`
	CommunicationWorkspace              string                              `json:"communication_workspace"`
	CommunicationChannel                string                              `json:"communication_channel"`
	CommunicationCudNotifications       bool                                `json:"communication_cud_notifications"`
	CommunicationApprovalNotifications  bool                                `json:"communication_approval_notifications"`
	CommunicationExecutionNotifications bool                                `json:"communication_execution_notifications"`
	Category                            string                              `json:"category"`
}

// GetContainer returns the runbook container regardless of operation type
func (r RunbookResponseAPIModelV1) GetContainer() *RunbookContainerV1 {
	if r.DefineNotebook != nil {
		return r.DefineNotebook
	}
	if r.UpdateNotebook != nil {
		return r.UpdateNotebook
	}
	if r.GetNotebookClass != nil {
		return r.GetNotebookClass
	}
	if r.DeleteNotebook != nil {
		return r.DeleteNotebook
	}
	return nil
}

// GetErrors returns a string representation of the V1 API response errors
func (r RunbookResponseAPIModelV1) GetErrors() string {
	container := r.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(r.Errors, container)
}
