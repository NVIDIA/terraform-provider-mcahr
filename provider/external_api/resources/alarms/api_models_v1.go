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

// AlarmResponseAPIModelV1 represents the response structure for V1 alarm API calls
type AlarmResponseAPIModelV1 struct {
	GetAlarmClass *AlarmContainerV1         `json:"get_alarm_class,omitempty"`
	DefineAlarm   *AlarmContainerV1         `json:"define_alarm,omitempty"`
	UpdateAlarm   *AlarmContainerV1         `json:"update_alarm,omitempty"`
	DeleteAlarm   *AlarmContainerV1         `json:"delete_alarm,omitempty"`
	Errors        *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// AlarmContainerV1 represents the alarm container in V1 API responses
type AlarmContainerV1 struct {
	AlarmClasses []AlarmClassV1    `json:"alarm_classes,omitempty"`
	Error        apicommon.ErrorV1 `json:"error,omitempty"`
	Errors       []string          `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure AlarmContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &AlarmContainerV1{}

// GetNestedError returns the nested error structure
func (c *AlarmContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *AlarmContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// AlarmClassV1 represents an alarm class in V1 API responses
type AlarmClassV1 struct {
	Name             string       `json:"name,omitempty"`
	Description      string       `json:"description,omitempty"`
	ResourceType     string       `json:"resource_type,omitempty"`
	Enabled          bool         `json:"enabled,omitempty"`
	ConfigData       ConfigDataV1 `json:"config_data"`
	FireStepClass    StepClassV1  `json:"fire_step_class"`
	ClearStepClass   StepClassV1  `json:"clear_step_class"`
	ResourceQuery    string       `json:"resource_query,omitempty"`
	FireQuery        string       `json:"fire_query,omitempty"`
	ClearQuery       string       `json:"clear_query,omitempty"`
	CheckIntervalSec int64        `json:"check_interval_sec,omitempty"`
}

// ConfigDataV1 represents config data structure in V1 API responses
type ConfigDataV1 struct {
	Family string `json:"family,omitempty"`
}

// StepClassV1 represents step class structure in V1 API responses
type StepClassV1 struct {
	TitleTemplate string `json:"title_template,omitempty"`
	ShortTemplate string `json:"short_template,omitempty"`
}

// GetContainer returns the appropriate alarm container from the V1 API response
func (r AlarmResponseAPIModelV1) GetContainer() *AlarmContainerV1 {
	if r.GetAlarmClass != nil {
		return r.GetAlarmClass
	}
	if r.DefineAlarm != nil {
		return r.DefineAlarm
	}
	if r.UpdateAlarm != nil {
		return r.UpdateAlarm
	}
	if r.DeleteAlarm != nil {
		return r.DeleteAlarm
	}
	return nil
}

// GetErrors returns a formatted error string from the V1 API response
func (r AlarmResponseAPIModelV1) GetErrors() string {
	container := r.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(r.Errors, container)
}
