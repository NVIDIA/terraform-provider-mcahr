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
)

func EmptyInt64Modifier() planmodifier.Int64 {
	return defaultInt64Modifier{defaultValue: 0}
}

func DefaultInt64Modifier(defaultValue int64) planmodifier.Int64 {
	return defaultInt64Modifier{defaultValue: defaultValue}
}

type defaultInt64Modifier struct {
	defaultValue int64
}

func (m defaultInt64Modifier) Description(_ context.Context) string {
	return Description()
}

func (m defaultInt64Modifier) MarkdownDescription(_ context.Context) string {
	return MarkdownDescription()
}

func (m defaultInt64Modifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {

	// If the plan/config value is known, do nothing.
	if IsPlanOrConfigKnown(req.PlanValue, req.ConfigValue) {
		return
	}

	resp.PlanValue = types.Int64Value(m.defaultValue)
}
