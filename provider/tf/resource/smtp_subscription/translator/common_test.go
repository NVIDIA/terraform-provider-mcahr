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

package translator

import (
	"context"
	"encoding/json"
	"testing"

	"terraform/terraform-provider/provider/common"
	subscriptionsapi "terraform/terraform-provider/provider/external_api/resources/smtp_subscription"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	smtp_subscriptiontf "terraform/terraform-provider/provider/tf/resource/smtp_subscription/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test TF model matching the terraform example
func createTestSmtpSubscriptionTFModel(t *testing.T) *smtp_subscriptiontf.SmtpSubscriptionTFModel {
	ctx := context.Background()

	// Recipients
	recipients := []string{"foo@bar.com", "bar@baz.com"}
	recipientsList, diags := types.ListValueFrom(ctx, types.StringType, recipients)
	require.False(t, diags.HasError(), "Failed to create recipients list")

	// Filters - must create as object list
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
	require.False(t, diags.HasError(), "Failed to create filters list")

	return &smtp_subscriptiontf.SmtpSubscriptionTFModel{
		ID:              types.StringValue("subscription-123"),
		Name:            types.StringValue("smtp_subscription"),
		IntegrationName: types.StringValue("smtp_integration"),
		Recipients:      recipientsList,
		Filters:         filtersList,
		Enabled:         types.BoolValue(true),
	}
}

func TestSmtpSubscriptionTranslatorCommon_ToAPIModel_Create(t *testing.T) {
	// Given
	translator := &SmtpSubscriptionTranslatorCommon{}
	tfModel := createTestSmtpSubscriptionTFModel(t)
	// remove the ID since CREATE operation should not include it
	tfModel.ID = types.StringNull()

	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).
		WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.ApiPayload)

	// Parse the payload to verify structure
	var payload map[string]interface{}
	err = json.Unmarshal([]byte(result.ApiPayload), &payload)
	require.NoError(t, err)

	// Verify basic fields
	assert.Equal(t, "smtp_subscription", payload["name"])
	assert.Equal(t, "smtp_integration", payload["integration_name"])
	assert.Equal(t, true, payload["enabled"])

	// Verify recipients is an array of strings (not attr.Value objects)
	recipients, ok := payload["recipients"].([]interface{})
	require.True(t, ok, "recipients should be an array")
	require.Len(t, recipients, 2)
	assert.Equal(t, "foo@bar.com", recipients[0])
	assert.Equal(t, "bar@baz.com", recipients[1])

	// Verify filters is an array of objects
	filters, ok := payload["filters"].([]interface{})
	require.True(t, ok, "filters should be an array")
	require.Len(t, filters, 4)

	// Verify first filter structure
	filter0, ok := filters[0].(map[string]interface{})
	require.True(t, ok, "filter should be an object")
	assert.Equal(t, "ACTION", filter0["category"])
	assert.Equal(t, "TRIGGER", filter0["type"])
	assert.Equal(t, "EXECUTING", filter0["status"])

	// Create should NOT include ID (empty string from struct)
	id, hasID := payload["id"]
	if hasID {
		assert.Empty(t, id, "Create operation should have empty ID")
	}
}

func TestSmtpSubscriptionTranslatorCommon_ToAPIModel_Update(t *testing.T) {
	// Given
	translator := &SmtpSubscriptionTranslatorCommon{}
	tfModel := createTestSmtpSubscriptionTFModel(t)

	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Update).
		WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(result.ApiPayload), &payload)
	require.NoError(t, err)

	// Update SHOULD include ID
	assert.Equal(t, "subscription-123", payload["id"])
}

func TestSmtpSubscriptionTranslatorCommon_ToAPIModel_Delete(t *testing.T) {
	// Given
	translator := &SmtpSubscriptionTranslatorCommon{}
	tfModel := createTestSmtpSubscriptionTFModel(t)

	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Delete).
		WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(result.ApiPayload), &payload)
	require.NoError(t, err)

	// Delete SHOULD include ID
	assert.Equal(t, "subscription-123", payload["id"])
}

func TestSmtpSubscriptionTranslatorCommon_RecipientsFormat(t *testing.T) {
	// Given - test that recipients serialize correctly as string array
	ctx := context.Background()
	translator := &SmtpSubscriptionTranslatorCommon{}

	// Create recipients list
	recipientsList, _ := types.ListValueFrom(ctx, types.StringType, []string{"a@test.com", "b@test.com"})

	// Create empty filters list
	filterAttrTypes := map[string]attr.Type{
		"category": types.StringType,
		"type":     types.StringType,
		"status":   types.StringType,
	}
	filtersList, _ := types.ListValue(types.ObjectType{AttrTypes: filterAttrTypes}, []attr.Value{})

	tfModel := &smtp_subscriptiontf.SmtpSubscriptionTFModel{
		Name:            types.StringValue("test"),
		IntegrationName: types.StringValue("integration"),
		Recipients:      recipientsList,
		Filters:         filtersList,
		Enabled:         types.BoolValue(true),
	}

	requestContext := common.NewRequestContext(ctx).
		WithOperation(common.Create).
		WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)

	// The JSON should contain recipients as ["a@test.com", "b@test.com"]
	// NOT as [{"value":"a@test.com"}, {"value":"b@test.com"}]
	assert.Contains(t, result.ApiPayload, `"recipients":["a@test.com","b@test.com"]`)
}

func TestSmtpSubscriptionTranslatorCommon_FiltersFormat(t *testing.T) {
	// Given - test that filters serialize correctly as object array
	ctx := context.Background()
	translator := &SmtpSubscriptionTranslatorCommon{}

	// Create empty recipients list
	recipientsList, _ := types.ListValueFrom(ctx, types.StringType, []string{})

	// Create filters list with one filter
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

	tfModel := &smtp_subscriptiontf.SmtpSubscriptionTFModel{
		Name:            types.StringValue("test"),
		IntegrationName: types.StringValue("integration"),
		Recipients:      recipientsList,
		Filters:         filtersList,
		Enabled:         types.BoolValue(true),
	}

	requestContext := common.NewRequestContext(ctx).
		WithOperation(common.Create).
		WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)

	// Parse to verify structure
	var payload map[string]interface{}
	json.Unmarshal([]byte(result.ApiPayload), &payload)

	filters := payload["filters"].([]interface{})
	require.Len(t, filters, 1)

	filter := filters[0].(map[string]interface{})
	assert.Equal(t, "ACTION", filter["category"])
	assert.Equal(t, "TRIGGER", filter["type"])
	assert.Equal(t, "EXECUTING", filter["status"])
}

// ============================================================================
// ToTFModel Tests (API Response -> TF Model)
// ============================================================================

// Helper function to create a test API response matching real API format
func createTestSmtpSubscriptionAPIResponse() *subscriptionsapi.SmtpSubscriptionResponseAPIModel {
	return &subscriptionsapi.SmtpSubscriptionResponseAPIModel{
		ID:              "f901c70d-5a30-4b81-97f1-68b631746a69",
		Name:            "smtp_subscription",
		IntegrationName: "smtp_integration",
		Recipients:      []string{"foo@bar.com", "bar@baz.com"},
		Filters: []subscriptionsapi.SmtpSubscriptionFilter{
			{Category: "ACTION", Type: "TRIGGER", Status: "EXECUTING"},
			{Category: "ALARM", Type: "TRIGGER", Status: "EXECUTING"},
			{Category: "BOT", Type: "TRIGGER", Status: "EXECUTING"},
			{Category: "TIME_TRIGGER", Type: "TRIGGER", Status: "EXECUTING"},
		},
		Enabled:       true,
		CreatedBy:     "user@example.com",
		UpdatedBy:     "user@example.com",
		CreatedTimeMs: 1767876664408,
		UpdatedTimeMs: 1767876664408,
	}
}

func createMinimalSmtpSubscriptionAPIResponse() *subscriptionsapi.SmtpSubscriptionResponseAPIModel {
	return &subscriptionsapi.SmtpSubscriptionResponseAPIModel{
		ID:              "minimal-id",
		Name:            "minimal_subscription",
		IntegrationName: "minimal_integration",
		Recipients:      []string{"test@test.com"},
		Filters:         []subscriptionsapi.SmtpSubscriptionFilter{},
		Enabled:         false,
	}
}

func TestSmtpSubscriptionTranslatorCommon_ToTFModel_Success(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SmtpSubscriptionTranslatorCommon{}
	apiModel := createTestSmtpSubscriptionAPIResponse()
	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Read).
		WithAPIVersion(common.V2)

	// When
	result, err := translator.ToTFModelFromResponse(requestContext, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify basic fields
	assert.Equal(t, "f901c70d-5a30-4b81-97f1-68b631746a69", result.ID.ValueString())
	assert.Equal(t, "smtp_subscription", result.Name.ValueString())
	assert.Equal(t, "smtp_integration", result.IntegrationName.ValueString())
	assert.True(t, result.Enabled.ValueBool())

	// Verify recipients
	var recipients []string
	result.Recipients.ElementsAs(context.Background(), &recipients, false)
	assert.Equal(t, []string{"foo@bar.com", "bar@baz.com"}, recipients)

	// Verify filters
	var filters []smtp_subscriptiontf.FilterTFModel
	result.Filters.ElementsAs(context.Background(), &filters, false)
	require.Len(t, filters, 4)
	assert.Equal(t, "ACTION", filters[0].Category.ValueString())
	assert.Equal(t, "TRIGGER", filters[0].Type.ValueString())
	assert.Equal(t, "EXECUTING", filters[0].Status.ValueString())
}

func TestSmtpSubscriptionTranslatorCommon_ToTFModel_NilInput(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SmtpSubscriptionTranslatorCommon{}
	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Read).
		WithAPIVersion(common.V2)

	// When
	result, err := translator.ToTFModelFromResponse(requestContext, nil)

	// Then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestSmtpSubscriptionTranslatorCommon_ToTFModel_EmptyFilters(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SmtpSubscriptionTranslatorCommon{}
	apiModel := createMinimalSmtpSubscriptionAPIResponse()
	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Read).
		WithAPIVersion(common.V2)

	// When
	result, err := translator.ToTFModelFromResponse(requestContext, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal-id", result.ID.ValueString())
	assert.Equal(t, "minimal_subscription", result.Name.ValueString())
	assert.False(t, result.Enabled.ValueBool())

	// Verify empty filters list
	var filters []smtp_subscriptiontf.FilterTFModel
	result.Filters.ElementsAs(context.Background(), &filters, false)
	assert.Empty(t, filters)
}

func TestSmtpSubscriptionTranslatorCommon_ToTFModel_SingleRecipient(t *testing.T) {
	t.Parallel()
	// Given
	translator := &SmtpSubscriptionTranslatorCommon{}
	apiModel := &subscriptionsapi.SmtpSubscriptionResponseAPIModel{
		ID:              "single-recipient-id",
		Name:            "single_recipient_sub",
		IntegrationName: "integration",
		Recipients:      []string{"only@one.com"},
		Filters: []subscriptionsapi.SmtpSubscriptionFilter{
			{Category: "ACTION", Type: "CREATE", Status: "SUCCEEDED"},
		},
		Enabled: true,
	}
	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Read).
		WithAPIVersion(common.V2)

	// When
	result, err := translator.ToTFModelFromResponse(requestContext, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	var recipients []string
	result.Recipients.ElementsAs(context.Background(), &recipients, false)
	assert.Equal(t, []string{"only@one.com"}, recipients)

	var filters []smtp_subscriptiontf.FilterTFModel
	result.Filters.ElementsAs(context.Background(), &filters, false)
	require.Len(t, filters, 1)
	assert.Equal(t, "ACTION", filters[0].Category.ValueString())
	assert.Equal(t, "CREATE", filters[0].Type.ValueString())
	assert.Equal(t, "SUCCEEDED", filters[0].Status.ValueString())
}

// ============================================================================
// Translator Interface Tests (V1 and V2 wrappers)
// These verify the translator structs correctly delegate to common logic
// ============================================================================

func TestSmtpSubscriptionTranslator_ToTFModel(t *testing.T) {
	t.Parallel()
	// Given - V2 translator
	translator := &SmtpSubscriptionTranslator{}
	apiModel := createTestSmtpSubscriptionAPIResponse()
	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Read).
		WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "smtp_subscription", result.Name.ValueString())
}

func TestSmtpSubscriptionTranslatorV1_ToTFModel(t *testing.T) {
	t.Parallel()
	// Given - V1 translator (uses same response type as V2 for REST API)
	translator := &SmtpSubscriptionTranslatorV1{}
	apiModel := createTestSmtpSubscriptionAPIResponse()
	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Read).
		WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "smtp_subscription", result.Name.ValueString())
}

func TestSmtpSubscriptionTranslator_ToAPIModel(t *testing.T) {
	t.Parallel()
	// Given - V2 translator
	translator := &SmtpSubscriptionTranslator{}
	tfModel := createTestSmtpSubscriptionTFModel(t)
	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).
		WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModel(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ApiPayload)
}

func TestSmtpSubscriptionTranslatorV1_ToAPIModel(t *testing.T) {
	t.Parallel()
	// Given - V1 translator
	translator := &SmtpSubscriptionTranslatorV1{}
	tfModel := createTestSmtpSubscriptionTFModel(t)
	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).
		WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModel(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ApiPayload)
}
