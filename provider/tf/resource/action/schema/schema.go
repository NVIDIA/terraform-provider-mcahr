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
	"terraform/terraform-provider/provider/tf/core"
	defaultmodifiers "terraform/terraform-provider/provider/tf/core/plan/modifiers/default"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	schemabuilder "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ActionSchema struct{}

var _ coreschema.ResourceSchema = &ActionSchema{}

// GetActionSchema returns the schema definition for the Action resource.
//
// This schema defines all the attributes available for configuring an action,
// including required fields (name, command), optional configuration fields,
// template fields for notifications, and computed fields (resource_id).
func (s *ActionSchema) GetSchema() schema.Schema {

	builder := schemabuilder.NewSchemaBuilder()

	builder.AddMarkdownDescription("Action resource for executing commands")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name of the action",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	builder.AddAttribute("command", schema.StringAttribute{
		MarkdownDescription: "The command to execute for this action",
		Required:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	// Optional string attributes with default empty string
	builder.AddAttribute("description", schema.StringAttribute{
		MarkdownDescription: "Description of the action",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.EmptyStringModifier()},
	})

	builder.AddAttribute("res_env_var", schema.StringAttribute{
		MarkdownDescription: "Resource environment variable",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.EmptyStringModifier()},
	})

	builder.AddAttribute("resource_query", schema.StringAttribute{
		MarkdownDescription: "Query to identify resources",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.EmptyStringModifier(), defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	builder.AddAttribute("shell", schema.StringAttribute{
		MarkdownDescription: "Shell to use for command execution",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.EmptyStringModifier()},
	})

	// Template attributes for notifications
	builder.AddAttribute("start_title_template", schema.StringAttribute{
		MarkdownDescription: "Title template for start notifications",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.PrefixWithNameModifier("started")},
	})

	builder.AddAttribute("start_short_template", schema.StringAttribute{
		MarkdownDescription: "Short template for start notifications",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.EmptyStringModifier()},
	})

	builder.AddAttribute("start_long_template", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** Long template for start notifications (deprecated - server controlled)",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	builder.AddAttribute("error_title_template", schema.StringAttribute{
		MarkdownDescription: "Title template for error notifications",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.PrefixWithNameModifier("failed")},
	})

	builder.AddAttribute("error_short_template", schema.StringAttribute{
		MarkdownDescription: "Short template for error notifications",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.EmptyStringModifier()},
	})

	builder.AddAttribute("error_long_template", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** Long template for error notifications (deprecated - server controlled)",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	builder.AddAttribute("complete_title_template", schema.StringAttribute{
		MarkdownDescription: "Title template for completion notifications",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.PrefixWithNameModifier("completed")},
	})

	builder.AddAttribute("complete_short_template", schema.StringAttribute{
		MarkdownDescription: "Short template for completion notifications",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.EmptyStringModifier()},
	})

	builder.AddAttribute("complete_long_template", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** Long template for completion notifications (deprecated - server controlled)",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	// Communication attributes
	builder.AddAttribute("communication_workspace", schema.StringAttribute{
		MarkdownDescription: "Communication workspace",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			defaultmodifiers.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("communication_channel", schema.StringAttribute{
		MarkdownDescription: "Communication channel",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			defaultmodifiers.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("allowed_resources_query", schema.StringAttribute{
		MarkdownDescription: "Query for allowed resources",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.EmptyStringModifier(), defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	// Boolean and numeric attributes
	builder.AddAttribute("enabled", schema.BoolAttribute{
		MarkdownDescription: "Whether the action is enabled",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.Bool{defaultmodifiers.DefaultBoolModifier(false)},
	})

	builder.AddAttribute("timeout", schema.Int64Attribute{
		MarkdownDescription: "Timeout for action execution",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.Int64{defaultmodifiers.DefaultInt64Modifier(60000)},
	})

	// List attributes
	builder.AddAttribute("params", schema.ListAttribute{
		MarkdownDescription: "Action parameters",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers:       []planmodifier.List{defaultmodifiers.EmptyListModifier()},
	})

	builder.AddAttribute("resource_tags_to_export", schema.ListAttribute{
		MarkdownDescription: "Resource tags to export",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers:       []planmodifier.List{defaultmodifiers.EmptyListModifier()},
	})

	builder.AddAttribute("file_deps", schema.ListAttribute{
		MarkdownDescription: "File dependencies",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers:       []planmodifier.List{defaultmodifiers.EmptyListModifier()},
	})

	builder.AddAttribute("allowed_entities", schema.ListAttribute{
		MarkdownDescription: "Allowed entities",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers:       []planmodifier.List{defaultmodifiers.EmptyListModifier()},
	})

	builder.AddAttribute("editors", schema.ListAttribute{
		MarkdownDescription: "Editors of the action",
		Optional:            true,
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			defaultmodifiers.EmptyListModifier(),
		},
	})

	return builder.Build()
}

func (s *ActionSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{
		"communication_workspace": {
			MinVersion: "release-14.1.0",
		},
		"communication_channel": {
			MinVersion: "release-14.1.0",
		},
		"editors": {
			MinVersion: "release-18.0.0",
		},
	}
}
