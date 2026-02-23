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
	"fmt"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	timetriggertf "terraform/terraform-provider/provider/tf/resource/time_trigger/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTimeTriggerCommonParameters returns the common parameter string for define/update operations
func getTimeTriggerCommonParameters() string {
	return `time_trigger_name="test-time-trigger", ` +
		`fire_query="start time_trigger_action(action_name=\"test_action\")", ` +
		`start_date="2024-01-01T00:00:00Z", ` +
		`end_date="2024-12-31T23:59:59Z", ` +
		`enabled=true`
}

func TestTimeTriggerTranslatorCommon_ToAPIModel(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslatorCommon{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// Create test TF model
	tfModel := &timetriggertf.TimeTriggerTFModel{
		Name:      types.StringValue("test-time-trigger"),
		FireQuery: types.StringValue("start time_trigger_action(action_name=\"test_action\")"),
		StartDate: types.StringValue("2024-01-01T00:00:00Z"),
		EndDate:   types.StringValue("2024-12-31T23:59:59Z"),
		Enabled:   types.BoolValue(true),
	}

	tests := []struct {
		name      string
		operation common.CrudOperation
		expected  string
	}{
		{
			name:      "Create operation",
			operation: common.Create,
			expected:  fmt.Sprintf("define_time_trigger(%s)", getTimeTriggerCommonParameters()),
		},
		{
			name:      "Read operation",
			operation: common.Read,
			expected:  `get_time_trigger_class(time_trigger_name="test-time-trigger")`,
		},
		{
			name:      "Update operation",
			operation: common.Update,
			expected:  fmt.Sprintf("update_time_trigger(%s)", getTimeTriggerCommonParameters()),
		},
		{
			name:      "Delete operation",
			operation: common.Delete,
			expected:  `delete_time_trigger(time_trigger_name="test-time-trigger")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update requestContext with the specific operation for this test case
			testRequestContext := requestContext.WithOperation(tt.operation)

			// When
			result, err := translator.ToAPIModelWithVersion(testRequestContext, translationData, tfModel)

			// Then
			require.NoError(t, err)
			require.NotNil(t, result)

			// Verify exact statement match
			assert.Equal(t, tt.expected, result.Statement)

			// Verify BackendVersion is set correctly
			assert.Equal(t, common.V2, result.APIVersion)
		})
	}
}

func TestTimeTriggerTranslatorCommon_ToAPIModel_UnsupportedOperation(t *testing.T) {
	t.Parallel()

	// Given
	translator := &TimeTriggerTranslatorCommon{}
	tfModel := &timetriggertf.TimeTriggerTFModel{
		Name: types.StringValue("test-time-trigger"),
	}
	// Use an invalid operation (999 is not a valid CrudOperation)
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.CrudOperation(999)).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}

func TestTimeTriggerTranslatorCommon_ToAPIModel_MinimalFields(t *testing.T) {
	t.Parallel()

	// Given
	translator := &TimeTriggerTranslatorCommon{}
	tfModel := &timetriggertf.TimeTriggerTFModel{
		Name:      types.StringValue("minimal-trigger"),
		FireQuery: types.StringValue("start simple_action()"),
		StartDate: types.StringValue(""),
		EndDate:   types.StringValue(""),
		Enabled:   types.BoolValue(false),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify exact statement match
	expected := `define_time_trigger(time_trigger_name="minimal-trigger", fire_query="start simple_action()", start_date="", end_date="", enabled=false)`
	assert.Equal(t, expected, result.Statement)
}

func TestTimeTriggerTranslatorCommon_ToAPIModel_SpecialCharacters(t *testing.T) {
	t.Parallel()

	// Given
	translator := &TimeTriggerTranslatorCommon{}
	tfModel := &timetriggertf.TimeTriggerTFModel{
		Name:      types.StringValue("trigger-with-special-chars"),
		FireQuery: types.StringValue("start action_with_quotes(param=\"value with spaces\")"),
		StartDate: types.StringValue("2024-01-01T00:00:00Z"),
		EndDate:   types.StringValue("2024-12-31T23:59:59Z"),
		Enabled:   types.BoolValue(true),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify exact statement match with proper escaping
	expected := `define_time_trigger(time_trigger_name="trigger-with-special-chars", fire_query="start action_with_quotes(param=\"value with spaces\")", start_date="2024-01-01T00:00:00Z", end_date="2024-12-31T23:59:59Z", enabled=true)`
	assert.Equal(t, expected, result.Statement)
}

func TestTimeTriggerTranslatorCommon_BuildStatements(t *testing.T) {
	t.Parallel()

	// Given
	translator := &TimeTriggerTranslatorCommon{}
	tfModel := &timetriggertf.TimeTriggerTFModel{
		Name:      types.StringValue("test-statements"),
		FireQuery: types.StringValue("start test_action()"),
		StartDate: types.StringValue("2024-01-01T00:00:00Z"),
		EndDate:   types.StringValue("2024-12-31T23:59:59Z"),
		Enabled:   types.BoolValue(true),
	}
	backendVersion := &version.BackendVersion{Version: "2.0.0"}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2).WithBackendVersion(backendVersion)

	tests := []struct {
		name     string
		buildFn  func() string
		expected string
	}{
		{
			name: "buildCreateStatement",
			buildFn: func() string {
				return translator.buildCreateStatement(requestContext, &coretranslator.TranslationData{}, tfModel)
			},
			expected: `define_time_trigger(time_trigger_name="test-statements", fire_query="start test_action()", start_date="2024-01-01T00:00:00Z", end_date="2024-12-31T23:59:59Z", enabled=true)`,
		},
		{
			name: "buildUpdateStatement",
			buildFn: func() string {
				return translator.buildUpdateStatement(requestContext, &coretranslator.TranslationData{}, tfModel)
			},
			expected: `update_time_trigger(time_trigger_name="test-statements", fire_query="start test_action()", start_date="2024-01-01T00:00:00Z", end_date="2024-12-31T23:59:59Z", enabled=true)`,
		},
		{
			name: "buildReadStatement",
			buildFn: func() string {
				return translator.buildReadStatement(tfModel)
			},
			expected: `get_time_trigger_class(time_trigger_name="test-statements")`,
		},
		{
			name: "buildDeleteStatement",
			buildFn: func() string {
				return translator.buildDeleteStatement(tfModel)
			},
			expected: `delete_time_trigger(time_trigger_name="test-statements")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			result := tt.buildFn()

			// Then
			require.NotEmpty(t, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}
