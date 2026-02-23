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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPrefixWithNameModifier_ConfigValueSetShouldNotModify(t *testing.T) {
	t.Parallel()

	// Given - When config value is set, modifier should not modify it
	modifier := PrefixWithNameModifier("fired")
	ctx := context.Background()

	req := planmodifier.StringRequest{
		PlanValue:   types.StringValue("custom title"),
		ConfigValue: types.StringValue("custom title"), // Config is set - should return early
		StateValue:  types.StringNull(),
		// Plan is nil, but that's ok because IsPlanOrConfigKnown returns true
	}

	resp := &planmodifier.StringResponse{
		PlanValue: req.PlanValue,
	}

	// When
	modifier.PlanModifyString(ctx, req, resp)

	// Then - Should not modify the value (returns early due to known config)
	expected := types.StringValue("custom title")
	if !resp.PlanValue.Equal(expected) {
		t.Errorf("Expected plan value %v, got %v", expected, resp.PlanValue)
	}
}

func TestPrefixWithNameModifier_PlanValueKnownShouldNotModify(t *testing.T) {
	t.Parallel()

	// Given - When plan value is known, modifier should not modify it
	modifier := PrefixWithNameModifier("fired")
	ctx := context.Background()

	req := planmodifier.StringRequest{
		PlanValue:   types.StringValue("existing value"),
		ConfigValue: types.StringNull(),
		StateValue:  types.StringNull(),
		// Plan is nil, but that's ok because IsPlanOrConfigKnown returns true
	}

	resp := &planmodifier.StringResponse{
		PlanValue: req.PlanValue,
	}

	// When
	modifier.PlanModifyString(ctx, req, resp)

	// Then - Should not modify the value (returns early due to known plan)
	expected := types.StringValue("existing value")
	if !resp.PlanValue.Equal(expected) {
		t.Errorf("Expected plan value %v, got %v", expected, resp.PlanValue)
	}
}

func TestPrefixWithNameModifier_GeneratesPrefixedValue(t *testing.T) {
	t.Parallel()

	// Given - A resource plan with a name attribute
	modifier := PrefixWithNameModifier("fired")
	ctx := context.Background()

	// Create a plan with a name attribute for testing
	// Create a simple schema with a name attribute
	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"title": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	// Create the plan data with the name value
	planData := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":  tftypes.String,
			"title": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"name":  tftypes.NewValue(tftypes.String, "cpu_alarm"),
		"title": tftypes.NewValue(tftypes.String, nil), // null value for the field we're testing
	})

	planWithName := tfsdk.Plan{
		Raw:    planData,
		Schema: testSchema,
	}

	req := planmodifier.StringRequest{
		PlanValue:   types.StringNull(), // Unknown value - should trigger modification
		ConfigValue: types.StringNull(), // Not configured - should trigger modification
		StateValue:  types.StringNull(),
		Plan:        planWithName,
	}

	resp := &planmodifier.StringResponse{
		PlanValue: req.PlanValue,
	}

	// When
	modifier.PlanModifyString(ctx, req, resp)

	// Then - Should generate the prefixed value
	expected := types.StringValue("fired cpu_alarm")
	if !resp.PlanValue.Equal(expected) {
		t.Errorf("Expected plan value %v, got %v", expected, resp.PlanValue)
	}
	if resp.Diagnostics.HasError() {
		t.Fatalf("Unexpected error: %v", resp.Diagnostics)
	}
}
