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
	planmodifiers "terraform/terraform-provider/provider/tf/core/plan/modifiers/default"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type TimeTriggerSchema struct{}

var _ coreschema.ResourceSchema = &TimeTriggerSchema{}

// GetTimeTriggerSchema returns the schema definition for the TimeTrigger resource.
//
// This schema defines all the attributes available for configuring a time trigger,
// including required fields (name, fire_query) and optional configuration fields.
func (s *TimeTriggerSchema) GetSchema() schema.Schema {

	builder := coreschema.NewSchemaBuilder()

	builder.AddMarkdownDescription("A condition that triggers Notebooks.")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name/symbol for the trigger within backend and the op language (must be unique, only alphanumeric/underscore).",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	builder.AddAttribute("fire_query", schema.StringAttribute{
		MarkdownDescription: "The trigger condition for the TimeTrigger (e.g. 'every 5m').",
		Required:            true,
		PlanModifiers:       []planmodifier.String{planmodifiers.IgnoreWhitespaceModifier()},
	})

	// Optional attributes
	builder.AddAttribute("start_date", schema.StringAttribute{
		MarkdownDescription: "When the trigger condition starts firing (defaults to current time when provider starts). The accepted format is ISO8601, e.g. '2024-02-17T08:08:01'.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})

	builder.AddAttribute("end_date", schema.StringAttribute{
		MarkdownDescription: "When the trigger condition stops firing. (defaults to unset, e.g. no stop date). The accepted format is ISO8601, e.g. '2029-02-17T08:08:01'.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("enabled", schema.BoolAttribute{
		MarkdownDescription: "If the trigger is currently enabled or disabled.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	})

	return builder.Build()
}

func (s *TimeTriggerSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{}
}

func (s *TimeTriggerSchema) GetFieldComparisonRules() map[string]coreschema.FieldComparisonRule {
	return coreschema.DefaultFieldComparisonRules()
}
