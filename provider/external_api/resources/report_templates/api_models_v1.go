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

package report_templates

import (
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
)

// ReportTemplateClassV1 represents a single report template class in V1 API
type ReportTemplateClassV1 struct {
	Name   string `json:"name"`
	Blocks string `json:"blocks"`
	Links  string `json:"links"`
}

// ReportTemplateContainerV1 represents the container for report template operations in V1 API
type ReportTemplateContainerV1 struct {
	ReportTemplateClasses []ReportTemplateClassV1 `json:"report_template_classes"`
	Error                 apicommon.ErrorV1       `json:"error"`
	Errors                []string                `json:"errors,omitempty"` // Direct errors array for some error cases
}

// Ensure ReportTemplateContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &ReportTemplateContainerV1{}

// GetNestedError returns the nested error structure
func (c *ReportTemplateContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array
func (c *ReportTemplateContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// ReportTemplateResponseAPIModelV1 represents the V1 API response structure
type ReportTemplateResponseAPIModelV1 struct {
	DefineReportTemplate   *ReportTemplateContainerV1 `json:"define_report_template,omitempty"`
	UpdateReportTemplate   *ReportTemplateContainerV1 `json:"update_report_template,omitempty"`
	GetReportTemplateClass *ReportTemplateContainerV1 `json:"get_report_template_class,omitempty"`
	DeleteReportTemplate   *ReportTemplateContainerV1 `json:"delete_report_template,omitempty"`
	Errors                 *apicommon.SyntaxErrorsV1  `json:"errors,omitempty"` // Top-level syntax errors
}

// GetContainer returns the appropriate container based on the operation type
func (r ReportTemplateResponseAPIModelV1) GetContainer() *ReportTemplateContainerV1 {
	if r.DefineReportTemplate != nil {
		return r.DefineReportTemplate
	}
	if r.UpdateReportTemplate != nil {
		return r.UpdateReportTemplate
	}
	if r.GetReportTemplateClass != nil {
		return r.GetReportTemplateClass
	}
	if r.DeleteReportTemplate != nil {
		return r.DeleteReportTemplate
	}
	return nil
}

// GetErrors returns a string representation of the V1 API response errors
func (r ReportTemplateResponseAPIModelV1) GetErrors() string {
	container := r.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(r.Errors, container)
}
