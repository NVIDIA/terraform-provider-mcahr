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
	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"
)

// ReportTemplateData represents the report template data in V2 API
type ReportTemplateData struct {
	Name   string                      `json:"name"`
	Blocks []customattribute.BlockJson `json:"blocks"`
	Links  []customattribute.LinkJson  `json:"links"`
}

// ReportTemplateEntityMetadata represents the metadata for report template entity in V2 API
type ReportTemplateEntityMetadata struct {
	Name string `json:"name"`
}

// ReportTemplateConfigurationItem represents a single configuration item in V2 API
type ReportTemplateConfigurationItem struct {
	Config         ReportTemplateData           `json:"config"`
	EntityMetadata ReportTemplateEntityMetadata `json:"entity_metadata"`
}

// ReportTemplateConfigurations represents the configurations structure in V2 API
type ReportTemplateConfigurations struct {
	Items []ReportTemplateConfigurationItem `json:"items"`
}

// ReportTemplateOutput represents the output structure in V2 API responses
type ReportTemplateOutput struct {
	Configurations ReportTemplateConfigurations `json:"configurations"`
}

// ReportTemplateSummary represents the summary structure in V2 API responses
type ReportTemplateSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// ReportTemplateResponseAPIModel represents the V2 API response structure
type ReportTemplateResponseAPIModel struct {
	Output  ReportTemplateOutput  `json:"output"`
	Summary ReportTemplateSummary `json:"summary"`
}

// GetErrors returns a string representation of the V2 API response errors
func (r ReportTemplateResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(r.Summary.Status, r.Summary.Errors)
}
