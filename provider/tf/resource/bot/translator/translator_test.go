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

func TestBotTranslator_ToTFModel_Success(t *testing.T) {
	// Given
	translator := &BotTranslator{}
	apiModel := createFullBotResponseV2()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "full_bot", result.Name.ValueString())
	assert.Equal(t, "if full_alarm then full_action('/tmp') fi", result.Command.ValueString())
	assert.Equal(t, "<description>", result.Description.ValueString())
	assert.True(t, result.Enabled.ValueBool())
	assert.Equal(t, "custom", result.Family.ValueString())
	assert.Equal(t, "trigger_source", result.TriggerSource.ValueString())
	assert.Equal(t, "<external_trigger_id>", result.TriggerID.ValueString())
	assert.Equal(t, "host", result.AlarmResourceQuery.ValueString())
	assert.Equal(t, "<communication_workspace>", result.CommunicationWorkspace.ValueString())
	assert.Equal(t, "<communication_channel>", result.CommunicationChannel.ValueString())
	assert.Equal(t, "<integration_name>", result.IntegrationName.ValueString())
}

func TestBotTranslator_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &BotTranslator{}
	apiModel := createMinimalBotResponseV2()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
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

	// Verify that empty optional fields are empty strings
	assert.Equal(t, "", result.Description.ValueString())
	assert.Equal(t, "", result.Family.ValueString())
	assert.Equal(t, "", result.TriggerSource.ValueString())
	assert.Equal(t, "", result.TriggerID.ValueString())
	assert.Equal(t, "", result.AlarmResourceQuery.ValueString())
	assert.Equal(t, "", result.CommunicationWorkspace.ValueString())
	assert.Equal(t, "", result.CommunicationChannel.ValueString())
	assert.Equal(t, "", result.IntegrationName.ValueString())
}

func TestBotTranslator_ToTFModel_EmptyConfigurations(t *testing.T) {
	// Given
	translator := &BotTranslator{}
	apiModel := &botapi.BotResponseAPIModel{
		Output: botapi.BotOutput{
			Configurations: botapi.BotConfigurations{
				Items: []botapi.ConfigurationItem{}, // Empty list
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "no bot configurations found")
}

func TestBotTranslator_ToTFModel_NilInput(t *testing.T) {
	// Given
	translator := &BotTranslator{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// Helper function to create a full V2 API response for testing
func createFullBotResponseV2() *botapi.BotResponseAPIModel {
	return &botapi.BotResponseAPIModel{
		Output: botapi.BotOutput{
			Configurations: botapi.BotConfigurations{
				Items: []botapi.ConfigurationItem{
					{
						Config: botapi.BotConfig{
							TriggerSource:       "trigger_source",
							AlarmResourceQuery:  "host",
							IntegrationName:     "<integration_name>",
							TriggerEntityName:   "full_alarm",
							ExecutionEntityName: "full_action('/tmp')",
							CommunicationDest: botapi.BotCommunicationDestination{
								Channel:   "<communication_channel>",
								Workspace: "<communication_workspace>",
							},
							TriggerID: "<external_trigger_id>",
						},
						EntityMetadata: botapi.BotEntityMetadata{
							Enabled:     true,
							Name:        "full_bot",
							Family:      "custom",
							Description: "<description>",
						},
					},
				},
			},
		},
		Summary: botapi.BotSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}

// Helper function to create a minimal V2 API response for testing
func createMinimalBotResponseV2() *botapi.BotResponseAPIModel {
	return &botapi.BotResponseAPIModel{
		Output: botapi.BotOutput{
			Configurations: botapi.BotConfigurations{
				Items: []botapi.ConfigurationItem{
					{
						Config: botapi.BotConfig{
							TriggerSource:       "",
							AlarmResourceQuery:  "",
							IntegrationName:     "",
							TriggerEntityName:   "simple_alarm",
							ExecutionEntityName: "simple_action",
							CommunicationDest: botapi.BotCommunicationDestination{
								Channel:   "",
								Workspace: "",
							},
							TriggerID: "",
						},
						EntityMetadata: botapi.BotEntityMetadata{
							Enabled:     false,
							Name:        "minimal_bot",
							Family:      "",
							Description: "",
						},
					},
				},
			},
		},
		Summary: botapi.BotSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}
