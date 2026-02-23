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

func MaybeGetFromDeprecatedInt64Modifier(deprecatedFieldName string) planmodifier.Int64 {
	return maybeGetFromDeprecatedInt64Modifier{
		deprecatedFieldName: deprecatedFieldName,
	}
}

type maybeGetFromDeprecatedInt64Modifier struct {
	deprecatedFieldName string
}

func (m maybeGetFromDeprecatedInt64Modifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m maybeGetFromDeprecatedInt64Modifier) MarkdownDescription(_ context.Context) string {
	return "If this value is not provided, and the deprecated field has a non-empty value, the plan will use the deprecated field's value."
}

func (m maybeGetFromDeprecatedInt64Modifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {

	// If the config value is known, do nothing.
	if common.IsAttrKnown(req.ConfigValue) {
		return
	}

	var deprecatedValue types.Int64
	diags := req.Config.GetAttribute(ctx, path.Root(m.deprecatedFieldName), &deprecatedValue)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if !common.IsAttrKnown(deprecatedValue) {
		// Skip if deprecated value is not available
		return
	}

	resp.PlanValue = deprecatedValue
}
