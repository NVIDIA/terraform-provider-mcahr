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
	defaultmodifiers "terraform/terraform-provider/provider/tf/core/plan/modifiers/default"
	deprecation "terraform/terraform-provider/provider/tf/core/plan/modifiers/deprecation"
	nulls "terraform/terraform-provider/provider/tf/core/plan/modifiers/null"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	schemabuilder "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type BotSchema struct{}

var _ coreschema.ResourceSchema = &BotSchema{}

func (s *BotSchema) GetSchema() schema.Schema {

	builder := schemabuilder.NewSchemaBuilder()

	builder.AddMarkdownDescription("An automation that ties an Action to an Alert.")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name/symbol for the object within backend and the op language (must be unique, only alphanumeric/underscore).",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	builder.AddAttribute("command", schema.StringAttribute{
		MarkdownDescription: "The bot command in the format 'if <alarm> then <action> fi'.",
		Required:            true,
	})

	// Optional attributes
	builder.AddAttribute("description", schema.StringAttribute{
		MarkdownDescription: "Description of the bot.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("event_type", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** Use `trigger_source` instead.",
		Optional:            true,
		Computed:            true,
		DeprecationMessage:  "use `trigger_source` instead.",
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
		},
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.ConflictsWith(path.MatchRoot("trigger_source")),
		},
	})

	builder.AddAttribute("trigger_source", schema.StringAttribute{
		MarkdownDescription: "The source of the trigger. It's required when linking an **external** trigger to an execution entity (for example: `ALERTMANAGER`). If it's not set, then the source of the internal trigger entity will be saved in the state file.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			deprecation.MaybeGetFromDeprecatedStringModifier("event_type"),
		},
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})

	builder.AddAttribute("monitor_id", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** Use `trigger_id` instead.",
		Optional:            true,
		Computed:            true,
		DeprecationMessage:  "use `trigger_id` instead.",
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
		},
		Validators: []validator.String{
			stringvalidator.ConflictsWith(path.MatchRoot("trigger_id")),
		},
	})

	builder.AddAttribute("trigger_id", schema.StringAttribute{
		MarkdownDescription: "The ID of the trigger. It's required when linking an **external** trigger to an execution entity. If it's not set, then the id of the internal trigger entity will be saved in the state file.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			deprecation.MaybeGetFromDeprecatedStringModifier("monitor_id"),
		},
	})

	builder.AddAttribute("alarm_resource_query", schema.StringAttribute{
		MarkdownDescription: "Resource query for alarm context.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("communication_workspace", schema.StringAttribute{
		MarkdownDescription: "Communication workspace for notifications.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("communication_channel", schema.StringAttribute{
		MarkdownDescription: "Communication channel for notifications.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("integration_name", schema.StringAttribute{
		MarkdownDescription: "Integration name for external systems.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("enabled", schema.BoolAttribute{
		MarkdownDescription: "Whether the bot is enabled.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	})

	builder.AddAttribute("family", schema.StringAttribute{
		MarkdownDescription: "The family/category for the bot.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("custom"),
		PlanModifiers:       []planmodifier.String{defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	return builder.Build()
}

func (s *BotSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{
		"communication_workspace": {
			MinVersion: "release-14.1.0",
		},
		"communication_channel": {
			MinVersion: "release-14.1.0",
		},
		"integration_name": {
			MinVersion: "release-15.0.0",
		},
	}
}
