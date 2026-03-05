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

package system_settings

import (
	"terraform/terraform-provider/provider/common/attribute"
	defaultmodifiers "terraform/terraform-provider/provider/tf/core/plan/modifiers/default"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SystemSettingsSchema struct{}

var _ coreschema.ResourceSchema = &SystemSettingsSchema{}

// GetSchema returns the schema definition for the system_settings resource
func (s *SystemSettingsSchema) GetSchema() schema.Schema {

	builder := coreschema.NewSchemaBuilder()

	builder.AddMarkdownDescription("System-level settings. Note: there must only be one instance of this terraform resource named 'system_settings'.")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The name/symbol for the object within backend and the op language (must be unique, only alphanumeric/underscore). For system_settings, this must be 'system_settings'.",
		Validators: []validator.String{
			validators.NameValidator(),
			validators.ExactValueValidator("system_settings"),
		},
	})

	// Access Control attributes
	builder.AddAttribute("administrator_grants_create_user", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "If administrator is allowed to grant create user permissions.",
	})

	builder.AddAttribute("administrator_grants_create_user_token", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "If administrator is allowed to grant create user token permissions.",
	})

	builder.AddAttribute("administrator_grants_regenerate_user_token", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "If administrator is allowed to grant regenerate user token permissions.",
	})

	builder.AddAttribute("administrator_grants_read_user_token", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "If administrator is allowed to grant read user token permissions.",
	})

	// Runbook/Approval attributes
	builder.AddAttribute("approval_feature_enabled", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "If the approval feature is enabled.",
	})

	builder.AddAttribute("runbook_ad_hoc_approval_request_enabled", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If runbook ad-hoc approval requests are enabled.",
		PlanModifiers: []planmodifier.Bool{
			defaultmodifiers.DefaultBoolModifier(true),
		},
	})

	builder.AddAttribute("runbook_approval_request_expiry_time", schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Runbook approval request expiry time in minutes.",
		PlanModifiers: []planmodifier.Int64{
			defaultmodifiers.DefaultInt64Modifier(60),
		},
	})

	builder.AddAttribute("run_approval_expiry_time", schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Run approval expiry time in minutes.",
		PlanModifiers: []planmodifier.Int64{
			defaultmodifiers.DefaultInt64Modifier(60),
		},
	})

	builder.AddAttribute("approval_editable_allowed_resource_query_enabled", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "If approval editable allowed resource query is enabled.",
	})

	builder.AddAttribute("approval_allow_individual_notification", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If approval individual notifications are allowed.",
		PlanModifiers: []planmodifier.Bool{
			defaultmodifiers.DefaultBoolModifier(true),
		},
	})

	builder.AddAttribute("approval_optional_request_ticket_url", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If approval request ticket URL is optional.",
		PlanModifiers: []planmodifier.Bool{
			defaultmodifiers.DefaultBoolModifier(false),
		},
	})

	builder.AddAttribute("time_trigger_permissions_user", schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The user for time trigger permissions.",
		PlanModifiers: []planmodifier.String{
			defaultmodifiers.IgnoreWhitespaceModifier(),
		},
	})

	builder.AddAttribute("parallel_runs_fired_by_time_triggers", schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Number of parallel runs fired by time triggers.",
		PlanModifiers: []planmodifier.Int64{
			defaultmodifiers.DefaultInt64Modifier(10),
		},
	})

	// Audit attributes
	builder.AddAttribute("external_audit_storage_enabled", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "If external audit storage is enabled.",
	})

	builder.AddAttribute("external_audit_storage_type", schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("ELASTIC"),
		MarkdownDescription: "The type of external audit storage.",
		PlanModifiers:       []planmodifier.String{defaultmodifiers.IgnoreWhitespaceModifier()},
	})

	builder.AddAttribute("external_audit_storage_batch_period_sec", schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(5),
		MarkdownDescription: "External audit storage batch period in seconds.",
	})

	// General attributes
	builder.AddAttribute("environment_name", schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "The environment name.",
	})

	builder.AddAttribute("environment_name_background", schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The background color for the environment name.",
		PlanModifiers: []planmodifier.String{
			defaultmodifiers.DefaultStringModifier("#EF5350"),
			defaultmodifiers.IgnoreWhitespaceModifier(),
		},
	})

	builder.AddAttribute("param_value_max_length", schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Maximum length for parameter values.",
		PlanModifiers: []planmodifier.Int64{
			defaultmodifiers.DefaultInt64Modifier(10000),
		},
	})

	builder.AddAttribute("maintenance_mode_enabled", schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If maintenance mode is enabled.",
		PlanModifiers: []planmodifier.Bool{
			defaultmodifiers.DefaultBoolModifier(false),
		},
	})

	builder.AddAttribute("allowed_tags", schema.ListAttribute{
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
		MarkdownDescription: "Set of allowed tags.",
		PlanModifiers:       []planmodifier.List{},
	})

	builder.AddAttribute("skipped_tags", schema.ListAttribute{
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
		MarkdownDescription: "Set of skipped tags.",
		PlanModifiers:       []planmodifier.List{},
	})

	builder.AddAttribute("managed_secrets", schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The type of managed secrets.",
		PlanModifiers: []planmodifier.String{
			defaultmodifiers.DefaultStringModifier("LOCAL"),
			defaultmodifiers.IgnoreWhitespaceModifier(),
		},
	})

	return builder.Build()
}

func (s *SystemSettingsSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{
		"runbook_ad_hoc_approval_request_enabled": {MinVersion: "release-25.0.0"},
		"runbook_approval_request_expiry_time":    {MinVersion: "release-25.0.0"},
		"run_approval_expiry_time":                {MinVersion: "release-25.0.0"},
		"approval_allow_individual_notification":  {MinVersion: "release-17.0.0"},
		"approval_optional_request_ticket_url":    {MinVersion: "release-17.0.0"},
		"time_trigger_permissions_user":           {MinVersion: "release-19.1.0"},
		"parallel_runs_fired_by_time_triggers":    {MinVersion: "release-25.0.0"},
		"environment_name_background":             {MinVersion: "release-18.0.0"},
		"param_value_max_length":                  {MinVersion: "release-19.0.0"},
		"maintenance_mode_enabled":                {MinVersion: "release-25.1.0"},
		"allowed_tags":                            {MinVersion: "release-27.2.0"},
		"skipped_tags":                            {MinVersion: "release-27.2.0"},
		"managed_secrets":                         {MinVersion: "release-28.1.0"},
	}
}

func (s *SystemSettingsSchema) GetFieldComparisonRules() map[string]coreschema.FieldComparisonRule {
	return coreschema.DefaultFieldComparisonRules()
}
