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

package report_template

import (
	"terraform/terraform-provider/provider/common/attribute"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type ReportTemplateSchema struct{}

var _ coreschema.ResourceSchema = &ReportTemplateSchema{}

// GetSchema returns the schema for the report template resource
func (s *ReportTemplateSchema) GetSchema() schema.Schema {
	builder := coreschema.NewSchemaBuilder()

	builder.AddMarkdownDescription("Report template resource for creating report templates with blocks and links")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name of the report template",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	// JSON attributes
	builder.AddAttribute("blocks", schema.StringAttribute{
		MarkdownDescription: "The JSON encoded blocks of the report template",
		Required:            true,
	})

	builder.AddAttribute("blocks_full", schema.StringAttribute{
		MarkdownDescription: "Complete blocks configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	builder.AddAttribute("links", schema.StringAttribute{
		MarkdownDescription: "The JSON encoded links of a report template with other report templates",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("[]"),
	})

	builder.AddAttribute("links_full", schema.StringAttribute{
		MarkdownDescription: "Complete links configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	return builder.Build()
}

func (s *ReportTemplateSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{}
}

func (s *ReportTemplateSchema) GetFieldComparisonRules() map[string]coreschema.FieldComparisonRule {
	return map[string]coreschema.FieldComparisonRule{
		// Skip minimal fields - they have _full variants for comparison
		// The minimal fields will always differ (API adds defaults) which is expected
		"blocks": {
			Behavior: coreschema.SkipComparison,
			Reason:   "Has blocks_full variant. Use blocks_full for detecting API modifications.",
		},
		"links": {
			Behavior: coreschema.SkipComparison,
			Reason:   "Has links_full variant. Use links_full for detecting API modifications.",
		},

		// NOTE: blocks_full and links_full are NOT skipped
		// They will be compared to detect real API modifications
	}
}
