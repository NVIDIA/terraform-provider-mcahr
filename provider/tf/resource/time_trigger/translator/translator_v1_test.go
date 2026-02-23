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

func createTestTimeTriggerResponseV1() *timetriggerapi.TimeTriggerResponseAPIModelV1 {
	return &timetriggerapi.TimeTriggerResponseAPIModelV1{
		GetTimeTriggerClass: &timetriggerapi.TimeTriggerContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			TimeTriggerClasses: []timetriggerapi.TimeTriggerClassV1{
				{
					Name:      "test-trigger",
					Enabled:   true,
					FireQuery: "every 5m",
					StartDate: "2024-02-29T08:00:00",
					EndDate:   "2100-02-28T08:00:00",
				},
			},
		},
	}
}

func createMinimalTimeTriggerResponseV1() *timetriggerapi.TimeTriggerResponseAPIModelV1 {
	return &timetriggerapi.TimeTriggerResponseAPIModelV1{
		DefineTimeTrigger: &timetriggerapi.TimeTriggerContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			TimeTriggerClasses: []timetriggerapi.TimeTriggerClassV1{
				{
					Name:      "minimal-trigger",
					Enabled:   false,
					FireQuery: "every 1h",
					StartDate: "2024-02-29T08:00:00",
					EndDate:   "",
				},
			},
		},
	}
}

func TestTimeTriggerTranslatorV1_ToTFModel_Success(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslatorV1{}
	apiModel := createTestTimeTriggerResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
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

func TestTimeTriggerTranslatorV1_ToTFModel_Minimal(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslatorV1{}
	apiModel := createMinimalTimeTriggerResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal-trigger", result.Name.ValueString())
	assert.Equal(t, "every 1h", result.FireQuery.ValueString())
	assert.Equal(t, "2024-02-29T08:00:00", result.StartDate.ValueString())
	assert.Equal(t, "", result.EndDate.ValueString())
	assert.False(t, result.Enabled.ValueBool())
}

func TestTimeTriggerTranslatorV1_ToTFModel_NilInput(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestTimeTriggerTranslatorV1_ToTFModel_NoContainer(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslatorV1{}
	apiModel := &timetriggerapi.TimeTriggerResponseAPIModelV1{
		// No containers set
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no time trigger container found in V1 API response")
}

func TestTimeTriggerTranslatorV1_ToTFModel_EmptyTimeTriggerClasses(t *testing.T) {
	// Given
	translator := &TimeTriggerTranslatorV1{}
	apiModel := &timetriggerapi.TimeTriggerResponseAPIModelV1{
		GetTimeTriggerClass: &timetriggerapi.TimeTriggerContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			TimeTriggerClasses: []timetriggerapi.TimeTriggerClassV1{}, // Empty
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no time trigger classes found in V1 API response")
}
