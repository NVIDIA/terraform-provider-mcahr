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
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	"testing"

	"terraform/terraform-provider/provider/common"
	alarmapi "terraform/terraform-provider/provider/external_api/resources/alarms"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	alarmtf "terraform/terraform-provider/provider/tf/resource/alarm/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlarmTranslator_ToTFModel(t *testing.T) {
	// Given
	translator := &AlarmTranslator{}
	apiModel := &alarmapi.AlarmResponseAPIModel{
		Output: alarmapi.AlarmOutput{
			Configurations: alarmapi.AlarmConfigurations{
				Count: 1,
				Items: []alarmapi.ConfigurationItem{
					{
						Config: alarmapi.AlarmConfig{
							FireQuery:        "cpu_usage > 80",
							ClearQuery:       "cpu_usage < 70",
							ResourceQuery:    "hosts",
							ResourceType:     "HOST",
							CheckIntervalSec: 60,
							StepDetails: alarmapi.StepDetails{
								FireStep: alarmapi.FireStep{
									Title:       "High CPU Alert",
									Description: "High CPU",
								},
								ClearStep: alarmapi.ClearStep{
									Title:       "CPU Alert Cleared",
									Description: "CPU Normal",
								},
							},
						},
						EntityMetadata: alarmapi.AlarmEntityMetadata{
							Enabled:     true,
							Name:        "cpu_alarm",
							Description: "CPU monitoring alarm",
							Family:      "custom",
						},
					},
				},
			},
		},
		Summary: alarmapi.AlarmSummary{
			Status: "success",
			Errors: []apicommon.Error{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "cpu_alarm", result.Name.ValueString())
	assert.Equal(t, "cpu_usage > 80", result.FireQuery.ValueString())
	assert.Equal(t, "cpu_usage < 70", result.ClearQuery.ValueString())
	assert.Equal(t, "", result.MuteQuery.ValueString()) // Not returned in V2 API response
	assert.Equal(t, "CPU monitoring alarm", result.Description.ValueString())
	assert.Equal(t, "hosts", result.ResourceQuery.ValueString())
	assert.Equal(t, "HOST", result.ResourceType.ValueString())
	assert.Equal(t, int64(60), result.CheckIntervalSec.ValueInt64())
	assert.Equal(t, "", result.RaiseFor.ValueString()) // Not returned in V2 API response
	assert.Equal(t, "custom", result.Family.ValueString())
	assert.True(t, result.Enabled.ValueBool())

	// Test template fields
	assert.Equal(t, "High CPU Alert", result.FireTitleTemplate.ValueString())
	assert.Equal(t, "", result.FireLongTemplate.ValueString()) // Not provided in API response
	assert.Equal(t, "High CPU", result.FireShortTemplate.ValueString())
	assert.Equal(t, "CPU Alert Cleared", result.ResolveTitleTemplate.ValueString())
	assert.Equal(t, "", result.ResolveLongTemplate.ValueString()) // Not provided in API response
	assert.Equal(t, "CPU Normal", result.ResolveShortTemplate.ValueString())

	// Condition details not returned in V2 API response
	assert.Equal(t, "", result.ConditionType.ValueString())
	assert.Equal(t, "", result.ConditionValue.ValueString())
	assert.Equal(t, "", result.MetricName.ValueString())
}

func TestAlarmTranslator_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &AlarmTranslator{}
	apiModel := &alarmapi.AlarmResponseAPIModel{
		Output: alarmapi.AlarmOutput{
			Configurations: alarmapi.AlarmConfigurations{
				Count: 1,
				Items: []alarmapi.ConfigurationItem{
					{
						Config: alarmapi.AlarmConfig{
							FireQuery: "cpu_usage > 80",
							StepDetails: alarmapi.StepDetails{
								FireStep: alarmapi.FireStep{
									Title:       "",
									Description: "",
								},
								ClearStep: alarmapi.ClearStep{
									Title:       "",
									Description: "",
								},
							},
						},
						EntityMetadata: alarmapi.AlarmEntityMetadata{
							Name:        "minimal_alarm",
							Description: "",
							Family:      "",
						},
					},
				},
			},
		},
		Summary: alarmapi.AlarmSummary{
			Status: "success",
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal_alarm", result.Name.ValueString())
	assert.Equal(t, "cpu_usage > 80", result.FireQuery.ValueString())
	assert.Equal(t, "", result.ClearQuery.ValueString())
	assert.Equal(t, "", result.Description.ValueString())
}

func TestAlarmTranslator_ToTFModel_Nil(t *testing.T) {
	// Given
	translator := &AlarmTranslator{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestAlarmTranslator_ToTFModel_NoConfigurations(t *testing.T) {
	// Given
	translator := &AlarmTranslator{}
	apiModel := &alarmapi.AlarmResponseAPIModel{
		Output: alarmapi.AlarmOutput{
			Configurations: alarmapi.AlarmConfigurations{
				Count: 0,
				Items: []alarmapi.ConfigurationItem{},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no configurations found")
}

func TestAlarmTranslator_ToAPIModel(t *testing.T) {
	// Given
	translator := &AlarmTranslator{}
	tfModel := &alarmtf.AlarmTFModel{
		Name:      types.StringValue("test_alarm"),
		FireQuery: types.StringValue("cpu_usage > 80"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModel(requestContext, translationData, tfModel)

	// Then
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Contains(t, result.Statement, "define_alarm")
	assert.Contains(t, result.Statement, "test_alarm")
	assert.Contains(t, result.Statement, "cpu_usage > 80")
	assert.Equal(t, common.V2, result.APIVersion)
}
