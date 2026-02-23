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

package actions

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	action "terraform/terraform-provider/provider/tf/resource/action/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions to create test data

func createTestActionTFModel() *action.ActionTFModel {
	return &action.ActionTFModel{
		Name:                  types.StringValue("test_action"),
		Command:               types.StringValue("echo test"),
		Enabled:               types.BoolValue(true),
		Timeout:               types.Int64Value(30),
		Description:           types.StringValue("Test action"),
		StartTitleTemplate:    types.StringValue("Starting test"),
		StartShortTemplate:    types.StringValue("Start"),
		CompleteTitleTemplate: types.StringValue("Completed test"),
		CompleteShortTemplate: types.StringValue("Complete"),
		ErrorTitleTemplate:    types.StringValue("Error in test"),
		ErrorShortTemplate:    types.StringValue("Error"),
		// Long templates initially null/empty
		StartLongTemplate:    types.StringNull(),
		CompleteLongTemplate: types.StringNull(),
		ErrorLongTemplate:    types.StringNull(),
	}
}

func createTestActionModelWithLongTemplates() *action.ActionTFModel {
	model := createTestActionTFModel()
	// Set long template values as if they were configured by user
	model.StartLongTemplate = types.StringValue("This is a detailed start long message")
	model.CompleteLongTemplate = types.StringValue("This is a detailed completion long message")
	model.ErrorLongTemplate = types.StringValue("This is a detailed error long message")
	return model
}

// Test PostProcessDelete (simplest case - no-op)
func TestActionPostProcessor_PostProcessDelete_Success(t *testing.T) {
	t.Parallel()
	// given
	processor := &ActionPostProcessor{}
	tfModel := createTestActionTFModel()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Delete).WithAPIVersion(common.V2)

	// when
	err := processor.PostProcessDelete(requestContext, nil, tfModel)

	// then
	assert.NoError(t, err)
	// No state changes expected for delete operation
}

// Test applyLongTemplateValues directly - core business logic
func TestActionPostProcessor_applyLongTemplateValues_WithNullValues(t *testing.T) {
	t.Parallel()
	// given
	processor := &ActionPostProcessor{}
	sourceModel := createTestActionTFModel() // Has null long templates
	targetModel := createTestActionTFModel()

	// when
	err := processor.applyLongTemplateValues(sourceModel, targetModel)

	// then
	require.NoError(t, err)

	// Values should remain null since source has null values
	assert.True(t, targetModel.StartLongTemplate.IsNull())
	assert.True(t, targetModel.CompleteLongTemplate.IsNull())
	assert.True(t, targetModel.ErrorLongTemplate.IsNull())
}

func TestActionPostProcessor_applyLongTemplateValues_WithSetValues(t *testing.T) {
	t.Parallel()
	// given
	processor := &ActionPostProcessor{}
	sourceModel := createTestActionModelWithLongTemplates() // Has set long templates
	targetModel := createTestActionTFModel()                // Has null long templates

	// when
	err := processor.applyLongTemplateValues(sourceModel, targetModel)

	// then
	require.NoError(t, err)

	// Values should be copied from source
	assert.Equal(t, "This is a detailed start long message", targetModel.StartLongTemplate.ValueString())
	assert.Equal(t, "This is a detailed completion long message", targetModel.CompleteLongTemplate.ValueString())
	assert.Equal(t, "This is a detailed error long message", targetModel.ErrorLongTemplate.ValueString())
}

func TestActionPostProcessor_applyLongTemplateValues_PartiallySet(t *testing.T) {
	t.Parallel()
	// given
	processor := &ActionPostProcessor{}
	sourceModel := createTestActionTFModel()
	// Only set some long templates
	sourceModel.StartLongTemplate = types.StringValue("Only start long message is set")
	sourceModel.CompleteLongTemplate = types.StringNull() // Keep as null
	sourceModel.ErrorLongTemplate = types.StringValue("Only error long message is set")

	targetModel := createTestActionTFModel()

	// when
	err := processor.applyLongTemplateValues(sourceModel, targetModel)

	// then
	require.NoError(t, err)

	// Only the set values should be copied
	assert.Equal(t, "Only start long message is set", targetModel.StartLongTemplate.ValueString())
	assert.True(t, targetModel.CompleteLongTemplate.IsNull()) // Should remain null
	assert.Equal(t, "Only error long message is set", targetModel.ErrorLongTemplate.ValueString())
}

func TestActionPostProcessor_applyLongTemplateValues_OverwritesExistingValues(t *testing.T) {
	t.Parallel()
	// given
	processor := &ActionPostProcessor{}
	sourceModel := createTestActionModelWithLongTemplates() // Has set long templates
	targetModel := createTestActionTFModel()

	// Set target to have some initial values
	targetModel.StartLongTemplate = types.StringValue("Old start long message")
	targetModel.CompleteLongTemplate = types.StringValue("Old complete long message")
	targetModel.ErrorLongTemplate = types.StringValue("Old error long message")

	// when
	err := processor.applyLongTemplateValues(sourceModel, targetModel)

	// then
	require.NoError(t, err)

	// Values should be overwritten with source values
	assert.Equal(t, "This is a detailed start long message", targetModel.StartLongTemplate.ValueString())
	assert.Equal(t, "This is a detailed completion long message", targetModel.CompleteLongTemplate.ValueString())
	assert.Equal(t, "This is a detailed error long message", targetModel.ErrorLongTemplate.ValueString())
}

func TestActionPostProcessor_applyLongTemplateValues_PreservesNullFromSource(t *testing.T) {
	t.Parallel()
	// given
	processor := &ActionPostProcessor{}
	sourceModel := createTestActionTFModel() // Has null long templates
	targetModel := createTestActionTFModel()

	// Set target to have some initial values
	targetModel.StartLongTemplate = types.StringValue("Old start long message")
	targetModel.CompleteLongTemplate = types.StringValue("Old complete long message")
	targetModel.ErrorLongTemplate = types.StringValue("Old error long message")

	// when
	err := processor.applyLongTemplateValues(sourceModel, targetModel)

	// then
	require.NoError(t, err)

	// Values should remain unchanged since source has null values
	assert.Equal(t, "Old start long message", targetModel.StartLongTemplate.ValueString())
	assert.Equal(t, "Old complete long message", targetModel.CompleteLongTemplate.ValueString())
	assert.Equal(t, "Old error long message", targetModel.ErrorLongTemplate.ValueString())
}

func TestActionPostProcessor_applyLongTemplateValues_EmptyStringVsNull(t *testing.T) {
	t.Parallel()
	// given
	processor := &ActionPostProcessor{}
	sourceModel := createTestActionTFModel()
	// Set source to have empty string values (not null)
	sourceModel.StartLongTemplate = types.StringValue("")
	sourceModel.CompleteLongTemplate = types.StringValue("")
	sourceModel.ErrorLongTemplate = types.StringValue("")

	targetModel := createTestActionTFModel()
	targetModel.StartLongTemplate = types.StringValue("Original start long message")
	targetModel.CompleteLongTemplate = types.StringValue("Original complete long message")
	targetModel.ErrorLongTemplate = types.StringValue("Original error long message")

	// when
	err := processor.applyLongTemplateValues(sourceModel, targetModel)

	// then
	require.NoError(t, err)

	// Empty string values should overwrite existing values
	assert.Equal(t, "", targetModel.StartLongTemplate.ValueString())
	assert.Equal(t, "", targetModel.CompleteLongTemplate.ValueString())
	assert.Equal(t, "", targetModel.ErrorLongTemplate.ValueString())
}
