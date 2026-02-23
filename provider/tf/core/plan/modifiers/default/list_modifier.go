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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func EmptyListModifier() planmodifier.List {
	return defaultListModifier{defaultValue: []interface{}{}}
}

func DefaultListModifier(defaultValue []interface{}) planmodifier.List {
	return defaultListModifier{defaultValue: defaultValue}
}

type defaultListModifier struct {
	defaultValue []interface{}
}

func (m defaultListModifier) Description(_ context.Context) string {
	return Description()
}

func (m defaultListModifier) MarkdownDescription(_ context.Context) string {
	return MarkdownDescription()
}

func (m defaultListModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {

	// If the plan/config value is known, do nothing.
	if IsPlanOrConfigKnown(req.PlanValue, req.ConfigValue) {
		return
	}

	// Extract element type from the plan value's type
	listType := req.PlanValue.Type(ctx).(basetypes.ListType)
	elementType := listType.ElemType

	listValue, diags := types.ListValueFrom(ctx, elementType, m.defaultValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = listValue
}
