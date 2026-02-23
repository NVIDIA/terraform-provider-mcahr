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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type AlarmSchema struct{}

var _ coreschema.ResourceSchema = &AlarmSchema{}

// GetAlarmSchema returns the schema definition for the Alarm resource.
//
// This schema defines all the attributes available for configuring an alarm,
// including required fields (name, fire_query), optional configuration fields,
// template fields for notifications, and computed fields.
func (s *AlarmSchema) GetSchema() schema.Schema {

	builder := schemabuilder.NewSchemaBuilder()

	builder.AddMarkdownDescription("A condition that triggers Alerts or Actions.")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name/symbol for the object within Shoreline and the op language (must be unique, only alphanumeric/underscore).",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	builder.AddAttribute("fire_query", schema.StringAttribute{
		MarkdownDescription: "The trigger condition for an Alarm (general expression) or the TimeTrigger (e.g. 'every 5m').",
		Required:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	// Optional string attributes with default empty string
	builder.AddAttribute("clear_query", schema.StringAttribute{
		MarkdownDescription: "The Alarm's resolution condition.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
		PlanModifiers:       []planmodifier.String{defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	builder.AddAttribute("mute_query", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** The Alarm's mute condition.",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	builder.AddAttribute("description", schema.StringAttribute{
		MarkdownDescription: "A user-friendly explanation of an object.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("resource_query", schema.StringAttribute{
		MarkdownDescription: "A set of Resources (e.g. host, pod, container), optionally filtered on tags or dynamic conditions.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
		PlanModifiers:       []planmodifier.String{defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	builder.AddAttribute("resource_type", schema.StringAttribute{
		MarkdownDescription: "The type of object (i.e., Alarm, Action, Bot, Resource, or File).",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("condition_type", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** Kind of check in an Alarm (e.g. above or below) vs a threshold for a Metric.",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	builder.AddAttribute("condition_value", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** Switching value (threshold) for a Metric in an Alarm.",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	builder.AddAttribute("metric_name", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** The Alarm's triggering Metric.",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	builder.AddAttribute("raise_for", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** Where an Alarm is raised (e.g., local to a resource, or global to the system).",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	builder.AddAttribute("family", schema.StringAttribute{
		MarkdownDescription: "General class for an Action or Bot (e.g., custom, standard, metric, or system check).",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("custom"),
		PlanModifiers:       []planmodifier.String{defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	// Template attributes for fire notifications
	builder.AddAttribute("fire_title_template", schema.StringAttribute{
		MarkdownDescription: "UI title of the Alarm's triggering condition.",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.PrefixWithNameModifier("fired")},
	})

	builder.AddAttribute("fire_short_template", schema.StringAttribute{
		MarkdownDescription: "The short description of the Alarm's triggering condition.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("fire_long_template", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** The long description of the Alarm's triggering condition.",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	// Template attributes for resolve notifications
	builder.AddAttribute("resolve_title_template", schema.StringAttribute{
		MarkdownDescription: "UI title of the Alarm's' resolution.",
		Optional:            true,
		Computed:            true,
		PlanModifiers:       []planmodifier.String{defaultmodifiers.PrefixWithNameModifier("cleared")},
	})

	builder.AddAttribute("resolve_short_template", schema.StringAttribute{
		MarkdownDescription: "The short description of the Alarm's resolution.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("resolve_long_template", schema.StringAttribute{
		MarkdownDescription: "**Deprecated** The long description of the Alarm's resolution.",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	// Boolean and numeric attributes
	builder.AddAttribute("enabled", schema.BoolAttribute{
		MarkdownDescription: "If the object is currently enabled or disabled.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	})

	builder.AddAttribute("check_interval_sec", schema.Int64Attribute{
		MarkdownDescription: "Interval (in seconds) between Alarm evaluations.",
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(1),
	})

	return builder.Build()
}

func (s *AlarmSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{}
}
