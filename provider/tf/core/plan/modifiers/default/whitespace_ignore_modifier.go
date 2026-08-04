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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// ignoreWhitespaceModifier is a custom plan modifier that ignores whitespace differences
type ignoreWhitespaceModifier struct{}

func (m ignoreWhitespaceModifier) Description(ctx context.Context) string {
	return "Ignores whitespace differences in strings"
}

func (m ignoreWhitespaceModifier) MarkdownDescription(ctx context.Context) string {
	return "Ignores whitespace differences in strings"
}

func (m ignoreWhitespaceModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {

	if IsPlanOrStateUnknown(req.PlanValue, req.StateValue) {
		// Cannot compare whitespace when plan or state values are unknown/null
		// State can be unknown/null: during initial resource creation
		// Plan can be unknown/null: when computed from other resources or using functions with unknown inputs
		return
	}

	stateNormalized := removeWhitespace(req.StateValue.ValueString())
	planNormalized := removeWhitespace(req.PlanValue.ValueString())

	if stateNormalized == planNormalized {
		resp.PlanValue = req.StateValue
	}
}

// IgnoreWhitespaceModifier returns a plan modifier that ignores whitespace differences
func IgnoreWhitespaceModifier() planmodifier.String {
	return ignoreWhitespaceModifier{}
}

// normalizeWhitespace collapses runs of whitespace to a single space and trims
// the ends, so values differing only in formatting compare equal.
//
// This used to delete every space instead, which is not the same thing: it made
// `echo a b` and `echo ab` indistinguishable, so an edit that changed only the
// spacing inside a command produced no plan diff and the backend kept running
// the old command. The attributes using this modifier -- action, runbook,
// time_trigger and system_settings -- are exactly the command-bearing ones,
// where silently discarding an operator's edit matters most.
func removeWhitespace(str string) string {
	return strings.Join(strings.Fields(str), " ")
}
