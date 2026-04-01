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
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	bottf "terraform/terraform-provider/provider/tf/resource/bot/model"
	"terraform/terraform-provider/provider/tf/resource/bot/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBotTranslatorCommon_ToAPIModel(t *testing.T) {
	tests := []struct {
		name      string
		operation common.CrudOperation
	}{
		{"Create operation", common.Create},
		{"Read operation", common.Read},
		{"Update operation", common.Update},
		{"Delete operation", common.Delete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			translator := &BotTranslatorCommon{}
			tfModel := &bottf.BotTFModel{
				Name:                   types.StringValue("test_bot"),
				Command:                types.StringValue("if cpu_alarm then restart_action fi"),
				Description:            types.StringValue("Test bot description"),
				Enabled:                types.BoolValue(true),
				Family:                 types.StringValue("custom"),
				TriggerSource:          types.StringValue("trigger_source"),
				TriggerID:              types.StringValue("trigger_123"),
				AlarmResourceQuery:     types.StringValue("host"),
				CommunicationWorkspace: types.StringValue("ops-workspace"),
				CommunicationChannel:   types.StringValue("alerts-channel"),
				IntegrationName:        types.StringValue("alertmanager"),
			}
			requestContext := common.NewRequestContext(context.Background()).WithOperation(tt.operation).WithAPIVersion(common.V1)
			translationData := &coretranslator.TranslationData{}

			// When
			result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

			// Then
			assert.NoError(t, err)
			require.NotNil(t, result)

			switch tt.operation {
			case common.Create:
				expected := "define_bot(" +
					"bot_name=\"test_bot\", " +
					"alarm_statement=\"cpu_alarm\", " +
					"action_statement=\"restart_action\", " +
					"description=\"Test bot description\", " +
					"enabled=true, " +
					"family=\"custom\", " +
					"trigger_source=\"trigger_source\", " +
					"external_trigger_id=\"trigger_123\", " +
					"alarm_resource_query=\"host\", " +
					"communication_workspace=\"ops-workspace\", " +
					"communication_channel=\"alerts-channel\", " +
					"integration_name=\"alertmanager\")"
				assert.Equal(t, expected, result.Statement)
			case common.Update:
				expected := "update_bot(" +
					"bot_name=\"test_bot\", " +
					"alarm_statement=\"cpu_alarm\", " +
					"action_statement=\"restart_action\", " +
					"description=\"Test bot description\", " +
					"enabled=true, " +
					"family=\"custom\", " +
					"trigger_source=\"trigger_source\", " +
					"external_trigger_id=\"trigger_123\", " +
					"alarm_resource_query=\"host\", " +
					"communication_workspace=\"ops-workspace\", " +
					"communication_channel=\"alerts-channel\", " +
					"integration_name=\"alertmanager\")"
				assert.Equal(t, expected, result.Statement)
			case common.Read:
				assert.Equal(t, "get_bot_class(bot_name=\"test_bot\")", result.Statement)
			case common.Delete:
				assert.Equal(t, "delete_bot(bot_name=\"test_bot\")", result.Statement)
			}
		})
	}
}

func TestBotTranslatorCommon_ToAPIModel_UnsupportedOperation(t *testing.T) {
	// Given
	translator := &BotTranslatorCommon{}
	tfModel := &bottf.BotTFModel{
		Name: types.StringValue("test_bot"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.CrudOperation(999)).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}

func TestBotTranslatorCommon_buildCreateStatement_MinimalFields(t *testing.T) {
	// Given
	translator := &BotTranslatorCommon{}
	tfModel := &bottf.BotTFModel{
		Name:    types.StringValue("minimal_bot"),
		Command: types.StringValue("if simple_alarm then simple_action fi"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	botSchema := schema.BotSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: botSchema.GetCompatibilityOptions()}

	// When
	statement := translator.buildCreateStatement(requestContext, translationData, tfModel)

	// Then
	expected := "define_bot(" +
		"bot_name=\"minimal_bot\", " +
		"alarm_statement=\"simple_alarm\", " +
		"action_statement=\"simple_action\", " +
		"description=\"\", " +
		"enabled=false, " +
		"family=\"\", " +
		"trigger_source=\"\", " +
		"external_trigger_id=\"\", " +
		"alarm_resource_query=\"\", " +
		"communication_workspace=\"\", " +
		"communication_channel=\"\", " +
		"integration_name=\"\")"
	assert.Equal(t, expected, statement)
}

func TestBotTranslatorCommon_buildCreateStatement_AllFields(t *testing.T) {
	// Given
	translator := &BotTranslatorCommon{}
	tfModel := &bottf.BotTFModel{
		Name:                   types.StringValue("full_bot"),
		Command:                types.StringValue("if cpu_alarm then restart_action fi"),
		Description:            types.StringValue("Full featured bot"),
		TriggerSource:          types.StringValue("trigger_source"),
		TriggerID:              types.StringValue("monitor_123"),
		AlarmResourceQuery:     types.StringValue("host"),
		CommunicationWorkspace: types.StringValue("workspace"),
		CommunicationChannel:   types.StringValue("channel"),
		IntegrationName:        types.StringValue("integration"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	botSchema := schema.BotSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: botSchema.GetCompatibilityOptions()}

	// When
	statement := translator.buildCreateStatement(requestContext, translationData, tfModel)

	// Then
	expected := "define_bot(" +
		"bot_name=\"full_bot\", " +
		"alarm_statement=\"cpu_alarm\", " +
		"action_statement=\"restart_action\", " +
		"description=\"Full featured bot\", " +
		"enabled=false, " +
		"family=\"\", " +
		"trigger_source=\"trigger_source\", " +
		"external_trigger_id=\"monitor_123\", " +
		"alarm_resource_query=\"host\", " +
		"communication_workspace=\"workspace\", " +
		"communication_channel=\"channel\", " +
		"integration_name=\"integration\")"
	assert.Equal(t, expected, statement)
}
