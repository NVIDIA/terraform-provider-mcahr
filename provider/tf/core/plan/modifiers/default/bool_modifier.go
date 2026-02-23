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

func EmptyBoolModifier() planmodifier.Bool {
	return defaultBoolModifier{defaultValue: false}
}

func DefaultBoolModifier(defaultValue bool) planmodifier.Bool {
	return defaultBoolModifier{defaultValue: defaultValue}
}

type defaultBoolModifier struct {
	defaultValue bool
}

func (m defaultBoolModifier) Description(_ context.Context) string {
	return Description()
}

func (m defaultBoolModifier) MarkdownDescription(_ context.Context) string {
	return MarkdownDescription()
}

func (m defaultBoolModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {

	// If the plan/config value is known, do nothing.
	if IsPlanOrConfigKnown(req.PlanValue, req.ConfigValue) {
		return
	}

	resp.PlanValue = types.BoolValue(m.defaultValue)
}
