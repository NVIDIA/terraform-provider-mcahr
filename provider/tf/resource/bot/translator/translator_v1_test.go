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

	botapi "terraform/terraform-provider/provider/external_api/resources/bots"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFullBotResponseV1() *botapi.BotResponseAPIModelV1 {
	return &botapi.BotResponseAPIModelV1{
		GetBotClass: &botapi.BotContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			BotClasses: []botapi.BotClassV1{
				{
					Name:        "test_bot",
					Description: "Test bot description",
					Enabled:     true,
					ConfigData: botapi.ConfigDataV1{
						Family: "monitoring",
					},
					Communication: botapi.CommunicationV1{
						Workspace: "ops-workspace",
						Channel:   "alerts-channel",
					},
					AlarmStatement:     "cpu_alarm",
					ActionStatement:    "restart_service",
					EventType:          "trigger_source",
					TriggerSource:      "trigger_source",
					MonitorID:          "alert_group_456",
					ExternalTriggerID:  "alert_group_456",
					AlarmResourceQuery: "host",
					IntegrationName:    "alertmanager",
				},
			},
		},
	}
}

func createMinimalBotResponseV1() *botapi.BotResponseAPIModelV1 {
	return &botapi.BotResponseAPIModelV1{
		GetBotClass: &botapi.BotContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			BotClasses: []botapi.BotClassV1{
				{
					Name:        "minimal_bot",
					Description: "",
					Enabled:     false,
					ConfigData: botapi.ConfigDataV1{
						Family: "",
					},
					Communication: botapi.CommunicationV1{
						Workspace: "",
						Channel:   "",
					},
					AlarmStatement:  "simple_alarm",
					ActionStatement: "simple_action",
				},
			},
		},
	}
}

func TestBotTranslatorV1_ToTFModel_Success(t *testing.T) {
	t.Parallel()

	// Given
	translator := &BotTranslatorV1{}
	apiModel := createFullBotResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_bot", result.Name.ValueString())
	assert.Equal(t, "if cpu_alarm then restart_service fi", result.Command.ValueString())
	assert.Equal(t, "Test bot description", result.Description.ValueString())
	assert.True(t, result.Enabled.ValueBool())
	assert.Equal(t, "monitoring", result.Family.ValueString())
	assert.Equal(t, "ops-workspace", result.CommunicationWorkspace.ValueString())
	assert.Equal(t, "alerts-channel", result.CommunicationChannel.ValueString())
	assert.Equal(t, "trigger_source", result.TriggerSource.ValueString())
	assert.Equal(t, "alert_group_456", result.TriggerID.ValueString())
	assert.Equal(t, "host", result.AlarmResourceQuery.ValueString())
	assert.Equal(t, "alertmanager", result.IntegrationName.ValueString())
}

func TestBotTranslatorV1_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &BotTranslatorV1{}
	apiModel := createMinimalBotResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify basic fields
	assert.Equal(t, "minimal_bot", result.Name.ValueString())
	assert.Equal(t, "if simple_alarm then simple_action fi", result.Command.ValueString())
	assert.False(t, result.Enabled.ValueBool())

	// Verify that empty optional fields are empty strings (preserving API response values)
	assert.Equal(t, "", result.Description.ValueString())
	assert.Equal(t, "", result.Family.ValueString())
	assert.Equal(t, "", result.CommunicationWorkspace.ValueString())
	assert.Equal(t, "", result.CommunicationChannel.ValueString())
	assert.Equal(t, "", result.TriggerSource.ValueString())
	assert.Equal(t, "", result.TriggerID.ValueString())
	assert.Equal(t, "", result.AlarmResourceQuery.ValueString())
	assert.Equal(t, "", result.IntegrationName.ValueString())
}

func TestBotTranslatorV1_ToTFModel_EmptyBotClasses(t *testing.T) {
	t.Parallel()

	// Given - V1 API response with empty bot classes (translator-level validation)
	apiModel := &botapi.BotResponseAPIModelV1{
		GetBotClass: &botapi.BotContainerV1{
			Error: apicommon.ErrorV1{
				Type: "OK",
			},
			BotClasses: []botapi.BotClassV1{}, // Empty list
		},
	}

	translator := &BotTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	tfModel, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	require.Nil(t, tfModel)
	assert.Contains(t, err.Error(), "no bot classes found")
}

func TestBotTranslatorV1_ToTFModel_NilInput(t *testing.T) {
	t.Parallel()

	// Given
	translator := &BotTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}
