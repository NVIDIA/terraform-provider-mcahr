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
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type PrincipalSchema struct{}

var _ coreschema.ResourceSchema = &PrincipalSchema{}

func (s *PrincipalSchema) GetSchema() schema.Schema {

	builder := coreschema.NewSchemaBuilder()
	builder.AddMarkdownDescription("Principal. An authorization group (e.g. Okta groups). Note: Admin privilege (in the platform) is required to create principal objects.")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name/symbol for the object within the platform and the op language (must be unique, only alphanumeric/underscore).",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	builder.AddAttribute("identity", schema.StringAttribute{
		MarkdownDescription: "The email address or provider's (e.g. Okta) group-name for a permissions group.",
		Required:            true,
	})

	// Optional attributes
	builder.AddAttribute("action_limit", schema.Int64Attribute{
		MarkdownDescription: "The number of simultaneous actions allowed for a permissions group.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	})

	builder.AddAttribute("execute_limit", schema.Int64Attribute{
		MarkdownDescription: "The number of simultaneous linux (shell) commands allowed for a permissions group.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	})

	builder.AddAttribute("configure_permission", schema.BoolAttribute{
		MarkdownDescription: "If a permissions group is allowed to perform \"configure\" actions.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	})

	builder.AddAttribute("administer_permission", schema.BoolAttribute{
		MarkdownDescription: "If a permissions group is allowed to perform \"administer\" actions.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	})

	builder.AddAttribute("idp_name", schema.StringAttribute{
		MarkdownDescription: "The Identity Provider's name.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("view_limit", schema.Int64Attribute{
		MarkdownDescription: "**Deprecated** The number of simultaneous metrics allowed for a permissions group.",
		Optional:            true,
		DeprecationMessage:  core.DeprecatedFieldMessage,
	})

	return builder.Build()
}

func (s *PrincipalSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{
		"idp_name": {
			MinVersion: "release-22.0.0",
		},
		"view_limit": {
			MaxVersion: "release-28.99.999",
		},
	}
}

func (s *PrincipalSchema) GetFieldComparisonRules() map[string]coreschema.FieldComparisonRule {
	return coreschema.DefaultFieldComparisonRules()
}
