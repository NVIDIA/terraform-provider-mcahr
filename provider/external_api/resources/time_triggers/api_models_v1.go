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

// TimeTriggerResponseAPIModelV1 represents the V1 API response format for /api/v1/execute endpoint
// Unified model that handles define_time_trigger, update_time_trigger, get_time_trigger_class, and delete_time_trigger responses
type TimeTriggerResponseAPIModelV1 struct {
	DefineTimeTrigger   *TimeTriggerContainerV1   `json:"define_time_trigger,omitempty"`
	UpdateTimeTrigger   *TimeTriggerContainerV1   `json:"update_time_trigger,omitempty"`
	GetTimeTriggerClass *TimeTriggerContainerV1   `json:"get_time_trigger_class,omitempty"`
	DeleteTimeTrigger   *TimeTriggerContainerV1   `json:"delete_time_trigger,omitempty"`
	Errors              *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// TimeTriggerContainerV1 represents the common structure for all V1 time trigger operation containers
type TimeTriggerContainerV1 struct {
	TimeTriggerClasses []TimeTriggerClassV1 `json:"time_trigger_classes"`
	Error              apicommon.ErrorV1    `json:"error"`
	Errors             []string             `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure TimeTriggerContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &TimeTriggerContainerV1{}

// GetNestedError returns the nested error structure
func (c *TimeTriggerContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *TimeTriggerContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// TimeTriggerClassV1 represents a time trigger class in V1 response
// Simplified to only include fields needed by the translator (matching TF model fields)
type TimeTriggerClassV1 struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	FireQuery string `json:"fire_query"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// GetContainer returns the time trigger container regardless of operation type
func (t TimeTriggerResponseAPIModelV1) GetContainer() *TimeTriggerContainerV1 {
	if t.DefineTimeTrigger != nil {
		return t.DefineTimeTrigger
	}
	if t.UpdateTimeTrigger != nil {
		return t.UpdateTimeTrigger
	}
	if t.GetTimeTriggerClass != nil {
		return t.GetTimeTriggerClass
	}
	if t.DeleteTimeTrigger != nil {
		return t.DeleteTimeTrigger
	}
	return nil
}

// GetErrors returns a string representation of the V1 API response errors
func (t TimeTriggerResponseAPIModelV1) GetErrors() string {
	container := t.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(t.Errors, container)
}
