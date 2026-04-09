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
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/migration"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	builder.AddAttribute("blocks", schema.StringAttribute{
		MarkdownDescription: "The JSON encoded blocks of the report template.",
		DeprecationMessage:  "Use blocks_list instead. The blocks attribute encodes blocks as a single JSON string, which causes Terraform to show full-string diffs even for small changes. blocks_list uses native Terraform list types for proper per-element diffs. This attribute will be removed in a future version.",
		Optional:            true,
		Computed:            true,
		Validators: []validator.String{
			stringvalidator.ConflictsWith(path.MatchRoot("blocks_list")),
		},
	})

	builder.AddAttribute("blocks_list", schema.ListNestedAttribute{
		MarkdownDescription: "The blocks of the report template as a native Terraform list. Provides better plan changes and drift detection than the deprecated `blocks` JSON string. Cannot be used together with `blocks`.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			migration.DefaultListWithDeprecatedConflict(path.MatchRoot("blocks")),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"title": schema.StringAttribute{
					MarkdownDescription: "The title of the block.",
					Required:            true,
				},
				"resource_query": schema.StringAttribute{
					MarkdownDescription: "The resource query for the block.",
					Required:            true,
				},
				"group_by_tag": schema.StringAttribute{
					MarkdownDescription: "The tag to group resources by.",
					Required:            true,
				},
				"breakdown_by_tag": schema.StringAttribute{
					MarkdownDescription: "The tag to break down resources by.",
					Required:            true,
				},
				"view_mode": schema.StringAttribute{
					MarkdownDescription: "The view mode (COUNT or PERCENTAGE).",
					Optional:            true,
					Computed:            true,
					PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				},
				"include_resources_without_group_tag": schema.BoolAttribute{
					MarkdownDescription: "Whether to include resources without the group tag.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
				"include_other_breakdown_tag_values": schema.BoolAttribute{
					MarkdownDescription: "Whether to include other breakdown tag values.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
				"other_tags_to_export": schema.ListAttribute{
					MarkdownDescription: "Additional tags to export.",
					Optional:            true,
					Computed:            true,
					ElementType:         types.StringType,
					Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				},
				"group_by_tag_order": schema.SingleNestedAttribute{
					MarkdownDescription: "The ordering configuration for group-by tags.",
					Optional:            true,
					Computed:            true,
					PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "The ordering type (DEFAULT, BY_TOTAL_ASC, BY_TOTAL_DESC, CUSTOM).",
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"values": schema.ListAttribute{
							MarkdownDescription: "Custom ordering values.",
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
							Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
						},
					},
				},
				"breakdown_tags_values": schema.ListNestedAttribute{
					MarkdownDescription: "Breakdown tag value configurations.",
					Required:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"color": schema.StringAttribute{
								MarkdownDescription: "The color of the breakdown value.",
								Required:            true,
							},
							"label": schema.StringAttribute{
								MarkdownDescription: "The label for the breakdown value.",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString(""),
							},
							"values": schema.ListAttribute{
								MarkdownDescription: "The values in the breakdown group.",
								Required:            true,
								ElementType:         types.StringType,
							},
						},
					},
				},
				"resources_breakdown": schema.ListNestedAttribute{
					MarkdownDescription: "Resources breakdown configurations.",
					Optional:            true,
					Computed:            true,
					Default:             listdefault.StaticValue(types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{"group_by_value": types.StringType, "breakdown_values": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"value": types.StringType, "count": types.Int64Type}}}}}, []attr.Value{})),
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"group_by_value": schema.StringAttribute{
								MarkdownDescription: "The group-by value.",
								Required:            true,
							},
							"breakdown_values": schema.ListNestedAttribute{
								MarkdownDescription: "The breakdown values.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"value": schema.StringAttribute{
											MarkdownDescription: "The breakdown value.",
											Required:            true,
										},
										"count": schema.Int64Attribute{
											MarkdownDescription: "The count for the breakdown value.",
											Required:            true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	builder.AddAttribute("blocks_full", schema.StringAttribute{
		MarkdownDescription: "Complete blocks configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	builder.AddAttribute("links", schema.StringAttribute{
		MarkdownDescription: "The JSON encoded links of a report template with other report templates.",
		DeprecationMessage:  "Use links_list instead. The links attribute encodes links as a single JSON string, which causes Terraform to show full-string diffs even for small changes. links_list uses native Terraform list types for proper per-element diffs. This attribute will be removed in a future version.",
		Optional:            true,
		Computed:            true,
		Validators: []validator.String{
			stringvalidator.ConflictsWith(path.MatchRoot("links_list")),
		},
	})

	builder.AddAttribute("links_list", schema.ListNestedAttribute{
		MarkdownDescription: "The links of a report template with other report templates as a native Terraform list. Provides better plan changes and drift detection than the deprecated `links` JSON string. Cannot be used together with `links`.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			migration.DefaultListWithDeprecatedConflict(path.MatchRoot("links")),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"label": schema.StringAttribute{
					MarkdownDescription: "The label for the link.",
					Required:            true,
				},
				"report_template_name": schema.StringAttribute{
					MarkdownDescription: "The name of the linked report template.",
					Required:            true,
				},
			},
		},
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
