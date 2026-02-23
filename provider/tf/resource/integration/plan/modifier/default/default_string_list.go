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

package defaults

import (
	"context"
	"terraform/terraform-provider/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EmptyStringListModifier() planmodifier.List {
	return DefaultStringListModifier([]string{})
}

// DefaultStringListModifier creates a service-aware modifier that sets specific default list values.
// The defaults will only be applied to fields that are compatible with the current service_name.
//
// Usage:
//
//	"payload_paths": schema.ListAttribute{
//	    Optional: true,
//	    ElementType: types.StringType,
//	    PlanModifiers: []planmodifier.List{
//	        defaults.DefaultStringListModifier([]string{"alerts.receiver", "alerts.status"}), // Only for Alertmanager
//	    },
//	}
func DefaultStringListModifier(defaultValue []string) planmodifier.List {
	return defaultStringListModifier{
		defaultValue: defaultValue,
	}
}

type defaultStringListModifier struct {
	defaultValue []string
}

func (m defaultStringListModifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m defaultStringListModifier) MarkdownDescription(_ context.Context) string {
	return "If the value is not set, the default value will be enforced."
}

func (m defaultStringListModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {

	// If the config value is known, do nothing.
	if common.IsAttrKnown(req.ConfigValue) {
		return
	}

	var serviceName types.String
	diags := req.Plan.GetAttribute(ctx, path.Root("service_name"), &serviceName)
	if diags.HasError() {
		return
	}
	if !common.IsAttrKnown(serviceName) {
		// Skip if service name is not available
		return
	}

	if !IsServiceNameCompatible(serviceName.ValueString(), GetAttributeName(req.Path)) {
		// Skip if service name is not compatible
		return
	}

	defaultValues := make([]attr.Value, len(m.defaultValue))
	for i, value := range m.defaultValue {
		defaultValues[i] = types.StringValue(value)
	}

	resp.PlanValue = types.ListValueMust(types.StringType, defaultValues)
}
