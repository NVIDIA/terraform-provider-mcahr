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
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/migration"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
		MarkdownDescription: "The JSON encoded groups of the dashboard.",
		DeprecationMessage:  "Use groups_list instead. The groups attribute encodes groups as a single JSON string, which causes Terraform to show full-string diffs even for small changes. groups_list uses native Terraform list types for proper per-element diffs. This attribute will be removed in a future version.",
		Optional:            true,
		Computed:            true,
		Validators: []validator.String{
			stringvalidator.ConflictsWith(path.MatchRoot("groups_list")),
		},
	})

	builder.AddAttribute("groups_list", schema.ListNestedAttribute{
		MarkdownDescription: "The groups of the dashboard as a native Terraform list. Provides better plan changes and drift detection than the deprecated `groups` JSON string. Cannot be used together with `groups`.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			migration.DefaultListWithDeprecatedConflict(path.MatchRoot("groups")),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: "The name of the group.",
					Required:            true,
				},
				"tags": schema.ListAttribute{
					MarkdownDescription: "The tags for the group.",
					Required:            true,
					ElementType:         types.StringType,
				},
			},
		},
	})

	builder.AddAttribute("groups_full", schema.StringAttribute{
		MarkdownDescription: "Complete groups configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	builder.AddAttribute("values", schema.StringAttribute{
		MarkdownDescription: "The JSON encoded values of the dashboard.",
		DeprecationMessage:  "Use values_list instead. The values attribute encodes values as a single JSON string, which causes Terraform to show full-string diffs even for small changes. values_list uses native Terraform list types for proper per-element diffs. This attribute will be removed in a future version.",
		Optional:            true,
		Computed:            true,
		Validators: []validator.String{
			stringvalidator.ConflictsWith(path.MatchRoot("values_list")),
		},
	})

	builder.AddAttribute("values_list", schema.ListNestedAttribute{
		MarkdownDescription: "The values of the dashboard as a native Terraform list. Provides better plan changes and drift detection than the deprecated `values` JSON string. Cannot be used together with `values`.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			migration.DefaultListWithDeprecatedConflict(path.MatchRoot("values")),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"color": schema.StringAttribute{
					MarkdownDescription: "The color of the value group.",
					Required:            true,
				},
				"values": schema.ListAttribute{
					MarkdownDescription: "The values in the group.",
					Required:            true,
					ElementType:         types.StringType,
				},
			},
		},
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
