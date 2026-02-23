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
	"terraform/terraform-provider/provider/common"
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	"testing"

	alarmapi "terraform/terraform-provider/provider/external_api/resources/alarms"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFullAlarmResponseFullV1() *alarmapi.AlarmResponseAPIModelV1 {
	return &alarmapi.AlarmResponseAPIModelV1{
		GetAlarmClass: &alarmapi.AlarmContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			AlarmClasses: []alarmapi.AlarmClassV1{
				{
					Name:             "test_alarm",
					Description:      "Test alarm description",
					ResourceType:     "HOST",
					Enabled:          true,
					ResourceQuery:    "host",
					FireQuery:        "cpu_usage > 80",
					ClearQuery:       "cpu_usage < 70",
					CheckIntervalSec: 60,
					ConfigData: alarmapi.ConfigDataV1{
						Family: "custom",
					},
					FireStepClass: alarmapi.StepClassV1{
						TitleTemplate: "fired test alarm",
						ShortTemplate: "fired test short",
					},
					ClearStepClass: alarmapi.StepClassV1{
						TitleTemplate: "cleared test alarm",
						ShortTemplate: "cleared test short",
					},
				},
			},
		},
	}
}

func createMinimalAlarmResponseV1() *alarmapi.AlarmResponseAPIModelV1 {
	return &alarmapi.AlarmResponseAPIModelV1{
		GetAlarmClass: &alarmapi.AlarmContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			AlarmClasses: []alarmapi.AlarmClassV1{
				{
					Name:             "minimal_alarm",
					Description:      "",
					ResourceType:     "",
					Enabled:          false,
					ResourceQuery:    "",
					FireQuery:        "cpu > 90",
					ClearQuery:       "cpu < 85",
					CheckIntervalSec: 1,
					ConfigData: alarmapi.ConfigDataV1{
						Family: "",
					},
					FireStepClass: alarmapi.StepClassV1{
						TitleTemplate: "",
						ShortTemplate: "",
					},
					ClearStepClass: alarmapi.StepClassV1{
						TitleTemplate: "",
						ShortTemplate: "",
					},
				},
			},
		},
	}
}

func TestAlarmTranslatorV1_ToTFModel_Success(t *testing.T) {
	t.Parallel()

	// Given
	translator := &AlarmTranslatorV1{}
	apiModel := createFullAlarmResponseFullV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_alarm", result.Name.ValueString())
	assert.Equal(t, "Test alarm description", result.Description.ValueString())
	assert.Equal(t, "HOST", result.ResourceType.ValueString())
	assert.True(t, result.Enabled.ValueBool())
	assert.Equal(t, "custom", result.Family.ValueString())
	assert.Equal(t, "host", result.ResourceQuery.ValueString())
	assert.Equal(t, "cpu_usage > 80", result.FireQuery.ValueString())
	assert.Equal(t, "cpu_usage < 70", result.ClearQuery.ValueString())
	assert.Equal(t, int64(60), result.CheckIntervalSec.ValueInt64())

	// Verify template fields mapped from step classes
	assert.Equal(t, "fired test alarm", result.FireTitleTemplate.ValueString())
	assert.Equal(t, "fired test short", result.FireShortTemplate.ValueString())
	assert.Equal(t, "cleared test alarm", result.ResolveTitleTemplate.ValueString())
	assert.Equal(t, "cleared test short", result.ResolveShortTemplate.ValueString())

}

func TestAlarmTranslatorV1_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &AlarmTranslatorV1{}
	apiModel := createMinimalAlarmResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify basic fields
	assert.Equal(t, "minimal_alarm", result.Name.ValueString())
	assert.Equal(t, "cpu > 90", result.FireQuery.ValueString())
	assert.Equal(t, "cpu < 85", result.ClearQuery.ValueString())
	assert.False(t, result.Enabled.ValueBool())
	assert.Equal(t, int64(1), result.CheckIntervalSec.ValueInt64())

	// Verify that empty optional fields are empty strings (preserving API response values)
	assert.Equal(t, "", result.Description.ValueString())
	assert.Equal(t, "", result.ResourceQuery.ValueString())
	assert.Equal(t, "", result.ResourceType.ValueString())
	assert.Equal(t, "", result.Family.ValueString())

	// Verify that template fields are empty strings (from step classes)
	assert.Equal(t, "", result.FireTitleTemplate.ValueString())
	assert.Equal(t, "", result.FireShortTemplate.ValueString())
	assert.Equal(t, "", result.ResolveTitleTemplate.ValueString())
	assert.Equal(t, "", result.ResolveShortTemplate.ValueString())
}

func TestAlarmTranslatorV1_ToTFModel_EmptyAlarmClasses(t *testing.T) {
	t.Parallel()

	// Given - V1 API response with empty alarm classes (translator-level validation)
	apiModel := &alarmapi.AlarmResponseAPIModelV1{
		GetAlarmClass: &alarmapi.AlarmContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			AlarmClasses: []alarmapi.AlarmClassV1{}, // Empty list
		},
	}

	translator := &AlarmTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	tfModel, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	require.Nil(t, tfModel)
	assert.Contains(t, err.Error(), "no alarm classes found")
}

func TestAlarmTranslatorV1_ToTFModel_NilInput(t *testing.T) {
	t.Parallel()

	// Given
	translator := &AlarmTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	tfModel, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	require.NoError(t, err)
	require.Nil(t, tfModel)
}
