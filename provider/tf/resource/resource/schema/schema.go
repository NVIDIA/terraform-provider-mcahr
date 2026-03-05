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

// ResourceSchema defines the schema for the resource type
type ResourceSchema struct{}

var _ coreschema.ResourceSchema = &ResourceSchema{}

// GetCompatibilityOptions returns compatibility options for the resource
func (s *ResourceSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{}
}

// GetSchema returns the schema structure
func (s *ResourceSchema) GetSchema() schema.Schema {
	builder := coreschema.NewSchemaBuilder()

	builder.AddMarkdownDescription("A server or compute resource in the system (e.g. host, pod, container).")

	// Required attributes
	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name/symbol for the resource within the system (must be unique, only alphanumeric/underscore).",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Validators: []validator.String{
			validators.NameValidator(),
		},
	})

	builder.AddAttribute("value", schema.StringAttribute{
		MarkdownDescription: "The Op statement that defines a Resource.",
		Required:            true,
		PlanModifiers: []planmodifier.String{
			defaultmodifiers.IgnoreWhitespaceModifier(),
			stringplanmodifier.UseStateForUnknown(),
		},
	})

	// Optional attributes
	builder.AddAttribute("description", schema.StringAttribute{
		MarkdownDescription: "A user-friendly explanation of the resource.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
	})

	builder.AddAttribute("params", schema.ListAttribute{
		MarkdownDescription: "Named variables to pass to the resource.",
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
		Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
	})

	return builder.Build()
}

func (s *ResourceSchema) GetFieldComparisonRules() map[string]coreschema.FieldComparisonRule {
	return coreschema.DefaultFieldComparisonRules()
}
