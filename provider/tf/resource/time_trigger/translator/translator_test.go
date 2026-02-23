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
	"testing"

	"terraform/terraform-provider/provider/common"
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	timetriggerapi "terraform/terraform-provider/provider/external_api/resources/time_triggers"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestTimeTriggerResponse() *timetriggerapi.TimeTriggerResponseAPIModel {
	return &timetriggerapi.TimeTriggerResponseAPIModel{
		Output: timetriggerapi.TimeTriggerOutput{
			Configurations: timetriggerapi.TimeTriggerConfigurations{
				Items: []timetriggerapi.ConfigurationItem{
					{
						Config: timetriggerapi.TimeTriggerConfig{
							FireQuery: "every 5m",
							StartDate: "2024-02-29T08:00:00",
							EndDate:   "2100-02-28T08:00:00",
						},
						EntityMetadata: timetriggerapi.TimeTriggerEntityMetadata{
							Enabled: true,
							Name:    "test-trigger",
						},
					},
				},
			},
		},
		Summary: timetriggerapi.TimeTriggerSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}

func createMinimalTimeTriggerResponse() *timetriggerapi.TimeTriggerResponseAPIModel {
	return &timetriggerapi.TimeTriggerResponseAPIModel{
		Output: timetriggerapi.TimeTriggerOutput{
			Configurations: timetriggerapi.TimeTriggerConfigurations{
				Items: []timetriggerapi.ConfigurationItem{
					{
						Config: timetriggerapi.TimeTriggerConfig{
							FireQuery: "every 10m",
							StartDate: "2024-02-29T08:00:00",
							EndDate:   "",
						},
						EntityMetadata: timetriggerapi.TimeTriggerEntityMetadata{
							Enabled: false,
							Name:    "minimal-trigger",
						},
					},
				},
			},
		},
		Summary: timetriggerapi.TimeTriggerSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}

func TestTimeTriggerTranslator_ToTFModel_Success(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslator{}
	apiModel := createTestTimeTriggerResponse()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test-trigger", result.Name.ValueString())
	assert.Equal(t, "every 5m", result.FireQuery.ValueString())
	assert.Equal(t, "2024-02-29T08:00:00", result.StartDate.ValueString())
	assert.Equal(t, "2100-02-28T08:00:00", result.EndDate.ValueString())
	assert.True(t, result.Enabled.ValueBool())
}

func TestTimeTriggerTranslator_ToTFModel_MinimalResponse(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslator{}
	apiModel := createMinimalTimeTriggerResponse()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal-trigger", result.Name.ValueString())
	assert.Equal(t, "every 10m", result.FireQuery.ValueString())
	assert.Equal(t, "2024-02-29T08:00:00", result.StartDate.ValueString())
	assert.Equal(t, "", result.EndDate.ValueString())
	assert.False(t, result.Enabled.ValueBool())
}

func TestTimeTriggerTranslator_ToTFModel_NilInput(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslator{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestTimeTriggerTranslator_ToTFModel_EmptyConfigurations(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslator{}
	apiModel := &timetriggerapi.TimeTriggerResponseAPIModel{
		Output: timetriggerapi.TimeTriggerOutput{
			Configurations: timetriggerapi.TimeTriggerConfigurations{
				Items: []timetriggerapi.ConfigurationItem{}, // Empty when there's an error
			},
		},
		Summary: timetriggerapi.TimeTriggerSummary{
			Status: "OP_COMPLETED", // Status is still OP_COMPLETED even with errors
			Errors: []apicommon.Error{
				{Type: "DUPLICATE_NAME", Message: "An object named full_time_trigger already exists. Please choose a globally-unique Time Trigger name to continue."},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no time trigger configurations found in API response")
}

func TestTimeTriggerTranslator_ToTFModel_EmptyItems(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslator{}
	apiModel := &timetriggerapi.TimeTriggerResponseAPIModel{
		Output: timetriggerapi.TimeTriggerOutput{
			Configurations: timetriggerapi.TimeTriggerConfigurations{
				Items: []timetriggerapi.ConfigurationItem{}, // Empty items
			},
		},
		Summary: timetriggerapi.TimeTriggerSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no time trigger configurations found in API response")
}

func TestTimeTriggerTranslator_ToTFModel_NilItems(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslator{}
	apiModel := &timetriggerapi.TimeTriggerResponseAPIModel{
		Output: timetriggerapi.TimeTriggerOutput{
			Configurations: timetriggerapi.TimeTriggerConfigurations{
				Items: nil, // Nil items
			},
		},
		Summary: timetriggerapi.TimeTriggerSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no time trigger configurations found in API response")
}
