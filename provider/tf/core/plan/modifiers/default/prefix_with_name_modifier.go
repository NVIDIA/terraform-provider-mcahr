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

package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type prefixWithNameModifier struct {
	prefixValue string
}

// PrefixWithNameModifier returns a plan modifier that sets a default value by
// concatenating a prefix with the resource's "name" attribute when the field
// is not explicitly configured.
//
// The resulting value will be: "{prefix} {resource_name}"
//
// Example:
//   - PrefixWithNameModifier("fired") with name="cpu_alarm" -> "fired cpu_alarm"
//   - PrefixWithNameModifier("cleared") with name="cpu_alarm" -> "cleared cpu_alarm"
func PrefixWithNameModifier(prefixValue string) planmodifier.String {
	return prefixWithNameModifier{prefixValue: prefixValue}
}

func (m prefixWithNameModifier) Description(_ context.Context) string {
	return Description()
}

func (m prefixWithNameModifier) MarkdownDescription(_ context.Context) string {
	return MarkdownDescription()
}

func (m prefixWithNameModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the plan/config value is known, do nothing.
	if IsPlanOrConfigKnown(req.PlanValue, req.ConfigValue) {
		return
	}

	// Get the resource name from the "name" attribute
	var resourceName types.String
	diags := req.Plan.GetAttribute(ctx, path.Root("name"), &resourceName)
	if diags.HasError() {
		return
	}

	// Skip if name is not available
	if resourceName.IsNull() || resourceName.IsUnknown() {
		return
	}

	// Concatenate prefix with resource name
	finalValue := fmt.Sprintf("%s %s", m.prefixValue, resourceName.ValueString())
	resp.PlanValue = types.StringValue(finalValue)
}
