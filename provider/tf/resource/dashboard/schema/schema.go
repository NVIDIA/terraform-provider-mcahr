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

package schema

import (
	"terraform/terraform-provider/provider/common/attribute"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DashboardSchema struct{}

var _ coreschema.ResourceSchema = &DashboardSchema{}

func (s *DashboardSchema) GetSchema() schema.Schema {

	builder := coreschema.NewSchemaBuilder()

	builder.AddMarkdownDescription("Dashboard resource for creating dashboards with groups and values")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name of the dashboard",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	builder.AddAttribute("dashboard_type", schema.StringAttribute{
		MarkdownDescription: "The type of the dashboard",
		Required:            true,
	})

	// Optional attributes
	builder.AddAttribute("resource_query", schema.StringAttribute{
		MarkdownDescription: "The resource query for the dashboard",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	// JSON attributes
	builder.AddAttribute("groups", schema.StringAttribute{
		MarkdownDescription: "The JSON encoded groups of the dashboard",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("[]"),
	})

	builder.AddAttribute("groups_full", schema.StringAttribute{
		MarkdownDescription: "Complete groups configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	builder.AddAttribute("values", schema.StringAttribute{
		MarkdownDescription: "The JSON encoded values of the dashboard",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("[]"),
	})

	builder.AddAttribute("values_full", schema.StringAttribute{
		MarkdownDescription: "Complete values configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	// Set attributes
	builder.AddAttribute("other_tags", schema.ListAttribute{
		MarkdownDescription: "Additional tags for the dashboard",
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
		Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
	})

	builder.AddAttribute("identifiers", schema.ListAttribute{
		MarkdownDescription: "Identifiers for the dashboard",
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
		Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
	})

	return builder.Build()
}

func (s *DashboardSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{}
}

func (s *DashboardSchema) GetFieldComparisonRules() map[string]coreschema.FieldComparisonRule {
	return map[string]coreschema.FieldComparisonRule{
		// Skip minimal fields - they have _full variants for comparison
		// The minimal fields will always differ (API adds defaults) which is expected
		"groups": {
			Behavior: coreschema.SkipComparison,
			Reason:   "Has groups_full variant. Use groups_full for detecting API modifications.",
		},
		"values": {
			Behavior: coreschema.SkipComparison,
			Reason:   "Has values_full variant. Use values_full for detecting API modifications.",
		},

		// NOTE: groups_full and values_full are NOT skipped
		// They will be compared to detect real API modifications
	}
}
