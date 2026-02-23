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
	nulls "terraform/terraform-provider/provider/tf/core/plan/modifiers/null"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	schemabuilder "terraform/terraform-provider/provider/tf/core/schema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RunbookSchema struct{}

var _ coreschema.ResourceSchema = &RunbookSchema{}

// GetRunbookSchema returns the schema definition for the Runbook resource.
//
// This schema defines all the attributes available for configuring a runbook,
// including required fields (name), optional configuration fields,
// and computed fields.
func (s *RunbookSchema) GetSchema() schema.Schema {

	builder := schemabuilder.NewSchemaBuilder()

	builder.AddMarkdownDescription("Runbook resource for interactive notebooks of Op commands and user documentation")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name of the runbook",
		// This is optional here to avoid raising error when name is provided in the "data" attribute
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	})

	// Optional complex data attributes (b64json)
	builder.AddAttribute("cells", schema.StringAttribute{
		MarkdownDescription: "The data cells inside a runbook. Defined as a list of JSON objects encoded in base64. These may be either Markdown or Op commands. Shows diffs only when configuration changes.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("[]"),
	})

	builder.AddAttribute("cells_full", schema.StringAttribute{
		MarkdownDescription: "Complete cells configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	builder.AddAttribute("params", schema.StringAttribute{
		MarkdownDescription: "Named variables to pass to a runbook, encoded as JSON. Shows diffs only when configuration changes.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("[]"),
	})

	builder.AddAttribute("params_full", schema.StringAttribute{
		MarkdownDescription: "Complete parameter configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	builder.AddAttribute("external_params", schema.StringAttribute{
		MarkdownDescription: "Runbook parameters defined via JSON path used to extract the parameter's value from an external payload, encoded as base64 JSON",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("[]"),
	})

	builder.AddAttribute("external_params_full", schema.StringAttribute{
		MarkdownDescription: "Complete external parameter configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.",
		Computed:            true,
	})

	// Optional boolean attributes
	builder.AddAttribute("enabled", schema.BoolAttribute{
		MarkdownDescription: "Whether the runbook is enabled",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	})

	// Optional string attributes
	builder.AddAttribute("description", schema.StringAttribute{
		MarkdownDescription: "Description of the runbook",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("allowed_resources_query", schema.StringAttribute{
		MarkdownDescription: "The list of resources on which a runbook can run. No restriction, if left empty.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("category", schema.StringAttribute{
		MarkdownDescription: "Specifies the category for this runbook. To use categories, make sure your platform administrator has enabled the `ENABLE_RUNBOOK_CATEGORIES` setting. Once enabled, you can organize your runbooks by assigning them a category.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})

	// Communication attributes
	builder.AddAttribute("communication_workspace", schema.StringAttribute{
		MarkdownDescription: "A string value denoting the slack workspace where notifications related to the runbook should be sent to",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("communication_channel", schema.StringAttribute{
		MarkdownDescription: "A string value denoting the slack channel where notifications related to the runbook should be sent to",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	// Numeric attributes
	builder.AddAttribute("timeout_ms", schema.Int64Attribute{
		MarkdownDescription: "Maximum time to wait for runbook execution, in milliseconds",
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(60000),
	})

	// Boolean communication flags
	builder.AddAttribute("communication_cud_notifications", schema.BoolAttribute{
		MarkdownDescription: "Enables slack notifications for create/update/delete operations",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	})

	builder.AddAttribute("communication_approval_notifications", schema.BoolAttribute{
		MarkdownDescription: "Enables slack notifications for approval operations",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	})

	builder.AddAttribute("communication_execution_notifications", schema.BoolAttribute{
		MarkdownDescription: "Enables slack notifications for runbook executions",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	})

	builder.AddAttribute("is_run_output_persisted", schema.BoolAttribute{
		MarkdownDescription: "A boolean value denoting whether or not cell outputs should be persisted when running a runbook",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	})

	builder.AddAttribute("filter_resource_to_action", schema.BoolAttribute{
		MarkdownDescription: "Determines whether parameters containing resources are exported to actions",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	})

	// Set attributes
	builder.AddAttribute("allowed_entities", schema.ListAttribute{
		MarkdownDescription: "The list of users who can run a runbook. Any user can run if left empty.",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
	})

	builder.AddAttribute("approvers", schema.ListAttribute{
		MarkdownDescription: "List of users who can approve runbook execution",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
	})

	builder.AddAttribute("labels", schema.ListAttribute{
		MarkdownDescription: "A list of strings by which runbooks can be grouped",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
	})

	builder.AddAttribute("editors", schema.ListAttribute{
		MarkdownDescription: "List of users who can edit the runbook (with configure permission). Empty maps to all users.",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
	})

	builder.AddAttribute("secret_names", schema.ListAttribute{
		MarkdownDescription: "A list of strings that contains the name of the secrets that are used in the runbook.",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
	})

	builder.AddAttribute("params_groups", schema.SingleNestedAttribute{
		MarkdownDescription: "Categorized parameter lists. Defaults to null if not specified.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			nulls.NullObjectIfUnknownModifier(),
		},
		Attributes: map[string]schema.Attribute{
			"required": schema.ListAttribute{
				MarkdownDescription: "List of required parameter names.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"optional": schema.ListAttribute{
				MarkdownDescription: "List of optional parameter names.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"exported": schema.ListAttribute{
				MarkdownDescription: "List of exported parameter names.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"external": schema.ListAttribute{
				MarkdownDescription: "List of external parameter names.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	})

	builder.AddAttribute("data", schema.StringAttribute{
		MarkdownDescription: "JSON-encoded string containing the runbook data. This can be loaded from a file using the `file()` function, e.g., `data = file(\"${path.module}/runbook.json\")`. " +
			"Unlike other JSON fields (params, cells, external_params), this field only stores what the user sets and does not have a corresponding _full attribute.",
		Optional: true,
		Computed: true,
	})

	return builder.Build()
}

func (s *RunbookSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{
		"allowed_entities": {
			MinVersion: "release-12.3.0",
		},
		"approvers": {
			MinVersion: "release-12.3.0",
		},
		"allowed_resources_query": {
			MinVersion: "release-12.3.0",
		},
		"is_run_output_persisted": {
			MinVersion: "release-12.3.0",
		},
		"communication_workspace": {
			MinVersion: "release-12.5.0",
		},
		"communication_channel": {
			MinVersion: "release-12.5.0",
		},
		"editors": {
			MinVersion: "release-15.1.0",
		},
		"labels": {
			MinVersion: "release-16.0.0",
		},
		"communication_cud_notifications": {
			MinVersion: "release-17.0.0",
		},
		"communication_approval_notifications": {
			MinVersion: "release-17.0.0",
		},
		"communication_execution_notifications": {
			MinVersion: "release-17.0.0",
		},
		"filter_resource_to_action": {
			MinVersion: "release-28.0.0",
		},
		"secret_names": {
			MinVersion: "release-28.1.0",
		},
		"category": {
			MinVersion: "release-29.1.0",
		},
		"params_groups": {
			MinVersion: "release-29.1.0",
		},
	}
}
