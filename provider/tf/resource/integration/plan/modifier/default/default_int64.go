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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EmptyInt64Modifier() planmodifier.Int64 {
	return DefaultInt64Modifier(0)
}

// DefaultInt64Modifier creates a service-aware modifier that sets a specific default int64 value.
// The default will only be applied to fields that are compatible with the current service_name.
//
// Usage:
//
//	"cache_ttl_ms": schema.Int64Attribute{
//	    Optional: true,
//	    PlanModifiers: []planmodifier.Int64{
//	        defaults.DefaultInt64Modifier(300000), // 5 minutes - only for caching services
//	    },
//	}
func DefaultInt64Modifier(defaultValue int64) planmodifier.Int64 {
	return defaultInt64Modifier{
		defaultValue: defaultValue,
	}
}

type defaultInt64Modifier struct {
	defaultValue int64
}

func (m defaultInt64Modifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m defaultInt64Modifier) MarkdownDescription(_ context.Context) string {
	return "If the value is not set, the default value will be enforced."
}

func (m defaultInt64Modifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {

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

	resp.PlanValue = types.Int64Value(m.defaultValue)
}
