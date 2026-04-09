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

package migration

import (
	"context"
	"fmt"
	"terraform/terraform-provider/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefaultListWithDeprecatedConflict returns a plan modifier for a ListNestedAttribute
// that manages the mutual exclusion with its deprecated JSON string counterpart.
//
// Behavior:
//   - If this list is set and the deprecated field is not: keep this list, null the deprecated field
//   - If this list is not set and the deprecated field is set: null this list, keep the deprecated field
//   - If neither is set: default this list to an empty list, null the deprecated field
//   - If both are set: do nothing (the schema conflict validator handles the error)
//
// This replaces the need for a schema default on the deprecated field and the manual
// mode enforcement logic in ModifyPlan. Reusable for cells_list, params_list, etc.
func DefaultListWithDeprecatedConflict(deprecatedFieldPath path.Expression) planmodifier.List {
	return defaultListWithDeprecatedConflictModifier{
		deprecatedFieldPath: deprecatedFieldPath,
	}
}

type defaultListWithDeprecatedConflictModifier struct {
	deprecatedFieldPath path.Expression
}

func (m defaultListWithDeprecatedConflictModifier) Description(_ context.Context) string {
	return fmt.Sprintf("Manages mutual exclusion with deprecated field %s: keeps this list and nulls the deprecated field when this list is set; nulls this list when the deprecated field is active; defaults to empty list when neither is set.", m.deprecatedFieldPath)
}

func (m defaultListWithDeprecatedConflictModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m defaultListWithDeprecatedConflictModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	deprecatedIsSet := isDeprecatedFieldSet(ctx, req, m.deprecatedFieldPath)
	listIsSet := common.IsAttrKnown(req.ConfigValue)

	switch {
	case listIsSet && !deprecatedIsSet:
		// List mode: keep list as-is (plan already has it)

	case !listIsSet && deprecatedIsSet:
		// Deprecated mode: null this list
		resp.PlanValue = types.ListNull(req.PlanValue.ElementType(ctx))

	case !listIsSet && !deprecatedIsSet:
		// Neither set: default to empty list (new list mode)
		resp.PlanValue = types.ListValueMust(req.PlanValue.ElementType(ctx), []attr.Value{})

		// Both set: schema ConflictsWith validator handles the error, do nothing here
	}
}

func isDeprecatedFieldSet(ctx context.Context, req planmodifier.ListRequest, deprecatedPath path.Expression) bool {
	matchedPaths, diags := req.Config.PathMatches(ctx, deprecatedPath)
	if diags.HasError() || len(matchedPaths) == 0 {
		return false
	}

	var deprecatedValue types.String
	diags = req.Config.GetAttribute(ctx, matchedPaths[0], &deprecatedValue)
	if diags.HasError() {
		return false
	}

	return common.IsAttrKnown(deprecatedValue)
}
