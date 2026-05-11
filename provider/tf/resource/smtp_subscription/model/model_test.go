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

package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmtpSubscriptionTFModel_GetName(t *testing.T) {
	// Given
	model := &SmtpSubscriptionTFModel{
		Name: types.StringValue("test_subscription"),
	}

	// When
	name := model.GetName()

	// Then
	assert.Equal(t, "test_subscription", name)
}

func TestSmtpSubscriptionTFModel_WithRecipients(t *testing.T) {
	// Given - create a model with recipients list
	ctx := context.Background()
	recipients := []string{"foo@bar.com", "bar@baz.com"}

	// Create the recipients list using types.ListValueFrom
	recipientsList, diags := types.ListValueFrom(ctx, types.StringType, recipients)
	require.False(t, diags.HasError(), "Failed to create recipients list")

	model := &SmtpSubscriptionTFModel{
		Name:            types.StringValue("test_subscription"),
		IntegrationName: types.StringValue("smtp_integration"),
		Recipients:      recipientsList,
		Enabled:         types.BoolValue(true),
	}

	// When - extract recipients from the model
	var extractedRecipients []string
	diags = model.Recipients.ElementsAs(ctx, &extractedRecipients, false)

	// Then
	require.False(t, diags.HasError(), "Failed to extract recipients")
	assert.Equal(t, []string{"foo@bar.com", "bar@baz.com"}, extractedRecipients)
}

func TestSmtpSubscriptionTFModel_WithFilters(t *testing.T) {
	// Given - create a model with filters list
	ctx := context.Background()

	// Define the filter object type (must match schema)
	filterAttrTypes := map[string]attr.Type{
		"category": types.StringType,
		"type":     types.StringType,
		"status":   types.StringType,
	}
	filterObjectType := types.ObjectType{AttrTypes: filterAttrTypes}

	// Create filter objects
	filter1, diags := types.ObjectValue(filterAttrTypes, map[string]attr.Value{
		"category": types.StringValue("ACTION"),
		"type":     types.StringValue("TRIGGER"),
		"status":   types.StringValue("EXECUTING"),
	})
	require.False(t, diags.HasError(), "Failed to create filter1")

	filter2, diags := types.ObjectValue(filterAttrTypes, map[string]attr.Value{
		"category": types.StringValue("ALARM"),
		"type":     types.StringValue("TRIGGER"),
		"status":   types.StringValue("EXECUTING"),
	})
	require.False(t, diags.HasError(), "Failed to create filter2")

	// Create the filters list
	filtersList, diags := types.ListValue(filterObjectType, []attr.Value{filter1, filter2})
	require.False(t, diags.HasError(), "Failed to create filters list")

	model := &SmtpSubscriptionTFModel{
		Name:            types.StringValue("test_subscription"),
		IntegrationName: types.StringValue("smtp_integration"),
		Filters:         filtersList,
		Enabled:         types.BoolValue(true),
	}

	// When - extract filters from the model using FilterTFModel struct
	var extractedFilters []FilterTFModel
	diags = model.Filters.ElementsAs(ctx, &extractedFilters, false)

	// Then
	require.False(t, diags.HasError(), "Failed to extract filters")
	require.Len(t, extractedFilters, 2)

	assert.Equal(t, "ACTION", extractedFilters[0].Category.ValueString())
	assert.Equal(t, "TRIGGER", extractedFilters[0].Type.ValueString())
	assert.Equal(t, "EXECUTING", extractedFilters[0].Status.ValueString())

	assert.Equal(t, "ALARM", extractedFilters[1].Category.ValueString())
	assert.Equal(t, "TRIGGER", extractedFilters[1].Type.ValueString())
	assert.Equal(t, "EXECUTING", extractedFilters[1].Status.ValueString())
}

func TestSmtpSubscriptionTFModel_FullModel(t *testing.T) {
	// Given - create a complete model matching terraform resource example
	ctx := context.Background()

	// Recipients
	recipients := []string{"foo@bar.com", "bar@baz.com"}
	recipientsList, diags := types.ListValueFrom(ctx, types.StringType, recipients)
	require.False(t, diags.HasError())

	// Filters
	filterAttrTypes := map[string]attr.Type{
		"category": types.StringType,
		"type":     types.StringType,
		"status":   types.StringType,
	}
	filterObjectType := types.ObjectType{AttrTypes: filterAttrTypes}

	filterValues := []struct {
		Category string
		Type     string
		Status   string
	}{
		{"ACTION", "TRIGGER", "EXECUTING"},
		{"ALARM", "TRIGGER", "EXECUTING"},
		{"BOT", "TRIGGER", "EXECUTING"},
		{"TIME_TRIGGER", "TRIGGER", "EXECUTING"},
	}

	filterObjects := make([]attr.Value, len(filterValues))
	for i, f := range filterValues {
		filterObj, diags := types.ObjectValue(filterAttrTypes, map[string]attr.Value{
			"category": types.StringValue(f.Category),
			"type":     types.StringValue(f.Type),
			"status":   types.StringValue(f.Status),
		})
		require.False(t, diags.HasError())
		filterObjects[i] = filterObj
	}

	filtersList, diags := types.ListValue(filterObjectType, filterObjects)
	require.False(t, diags.HasError())

	// Create full model
	model := &SmtpSubscriptionTFModel{
		ID:              types.StringValue("subscription-123"),
		Name:            types.StringValue("smtp_subscription"),
		IntegrationName: types.StringValue("smtp_integration"),
		Recipients:      recipientsList,
		Filters:         filtersList,
		Enabled:         types.BoolValue(true),
	}

	// Then - verify all fields
	assert.Equal(t, "subscription-123", model.ID.ValueString())
	assert.Equal(t, "smtp_subscription", model.Name.ValueString())
	assert.Equal(t, "smtp_integration", model.IntegrationName.ValueString())
	assert.True(t, model.Enabled.ValueBool())

	// Extract and verify recipients
	var extractedRecipients []string
	model.Recipients.ElementsAs(ctx, &extractedRecipients, false)
	assert.Equal(t, []string{"foo@bar.com", "bar@baz.com"}, extractedRecipients)

	// Extract and verify filters
	var extractedFilters []FilterTFModel
	model.Filters.ElementsAs(ctx, &extractedFilters, false)
	require.Len(t, extractedFilters, 4)
	assert.Equal(t, "ACTION", extractedFilters[0].Category.ValueString())
	assert.Equal(t, "ALARM", extractedFilters[1].Category.ValueString())
	assert.Equal(t, "BOT", extractedFilters[2].Category.ValueString())
	assert.Equal(t, "TIME_TRIGGER", extractedFilters[3].Category.ValueString())
}

// TestExtractRecipientsAsStrings demonstrates the correct way to extract recipients as []string
// This is what the translator should do
func TestExtractRecipientsAsStrings(t *testing.T) {
	ctx := context.Background()

	// Create a recipients list like terraform would
	recipientsList, _ := types.ListValueFrom(ctx, types.StringType, []string{"a@b.com", "c@d.com"})

	// WRONG way - Elements() returns []attr.Value, not []string
	elements := recipientsList.Elements()
	assert.IsType(t, []attr.Value{}, elements)

	// CORRECT way - use ElementsAs to extract as []string
	var recipientsStrings []string
	diags := recipientsList.ElementsAs(ctx, &recipientsStrings, false)
	require.False(t, diags.HasError())
	assert.Equal(t, []string{"a@b.com", "c@d.com"}, recipientsStrings)
}

// TestExtractFiltersAsStructs demonstrates the correct way to extract filters
// This is what the translator should do
func TestExtractFiltersAsStructs(t *testing.T) {
	ctx := context.Background()

	// Create filters like terraform would
	filterAttrTypes := map[string]attr.Type{
		"category": types.StringType,
		"type":     types.StringType,
		"status":   types.StringType,
	}
	filterObjectType := types.ObjectType{AttrTypes: filterAttrTypes}

	filter1, _ := types.ObjectValue(filterAttrTypes, map[string]attr.Value{
		"category": types.StringValue("ACTION"),
		"type":     types.StringValue("TRIGGER"),
		"status":   types.StringValue("EXECUTING"),
	})

	filtersList, _ := types.ListValue(filterObjectType, []attr.Value{filter1})

	// WRONG way - Elements() returns []attr.Value
	elements := filtersList.Elements()
	assert.IsType(t, []attr.Value{}, elements)

	// CORRECT way - use ElementsAs with FilterTFModel struct
	var filtersTF []FilterTFModel
	diags := filtersList.ElementsAs(ctx, &filtersTF, false)
	require.False(t, diags.HasError())
	require.Len(t, filtersTF, 1)
	assert.Equal(t, "ACTION", filtersTF[0].Category.ValueString())
	assert.Equal(t, "TRIGGER", filtersTF[0].Type.ValueString())
	assert.Equal(t, "EXECUTING", filtersTF[0].Status.ValueString())
}
