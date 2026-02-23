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
	"fmt"
	"regexp"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	bottf "terraform/terraform-provider/provider/tf/resource/bot/model"
)

// BotTranslatorCommon provides common functionality for bot translators across API versions
type BotTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (t *BotTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *bottf.BotTFModel) (*statement.StatementInputAPIModel, error) {
	var stmt string

	switch requestContext.Operation {
	case common.Create:
		stmt = t.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = t.buildReadStatement(tfModel)
	case common.Update:
		stmt = t.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt = t.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}

	apiModel := &statement.StatementInputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (t *BotTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *bottf.BotTFModel) string {
	return t.buildBotStatement(requestContext, translationData, "define_bot", tfModel)
}

func (t *BotTranslatorCommon) buildReadStatement(tfModel *bottf.BotTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_bot_class(bot_name=\"%s\")", name)
}

func (t *BotTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *bottf.BotTFModel) string {
	return t.buildBotStatement(requestContext, translationData, "update_bot", tfModel)
}

func (t *BotTranslatorCommon) buildDeleteStatement(tfModel *bottf.BotTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_bot(bot_name=\"%s\")", name)
}

// parseCommandStatement parses the command using the regex pattern from the schema
// Pattern: "^\\s*if\\s*(?P<alarm_statement>.*?)\\s*then\\s*(?P<action_statement>.*?)\\s*fi\\s*$"
// Mimics the old provider behavior: if no match, returns empty strings for both components
func (t *BotTranslatorCommon) parseCommandStatement(command string) (alarmStatement string, actionStatement string) {
	// Regex pattern from the old provider schema
	pattern := `^\s*if\s*(?P<alarm_statement>.*?)\s*then\s*(?P<action_statement>.*?)\s*fi\s*$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(command)
	if len(matches) == 3 {
		alarmStatement = matches[1]  // First capture group: alarm_statement
		actionStatement = matches[2] // Second capture group: action_statement
	}
	// If no match (len(matches) < 3), both remain empty strings - same as old provider behavior

	return alarmStatement, actionStatement
}

func (t *BotTranslatorCommon) buildBotStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *bottf.BotTFModel) string {
	// Build the bot statement from the TF model using the builder pattern
	// Used for both define_bot (create) and update_bot (update) operations

	// Parse the command to extract alarm_statement and action_statement
	alarmStatement, actionStatement := t.parseCommandStatement(tfModel.Command.ValueString())

	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("bot_name", tfModel.Name.ValueString(), "name")

	// Only set alarm_statement and action_statement if they're not empty (matches old provider behavior)
	if alarmStatement != "" {
		builder = builder.SetCommandField("alarm_statement", alarmStatement, "")
	}
	if actionStatement != "" {
		builder = builder.SetCommandField("action_statement", actionStatement, "")
	}

	return builder.
		SetStringField("description", tfModel.Description.ValueString(), "description").
		SetField("enabled", tfModel.Enabled.ValueBool(), "enabled").
		SetStringField("family", tfModel.Family.ValueString(), "family").
		SetStringField("trigger_source", tfModel.TriggerSource.ValueString(), "trigger_source").
		SetStringField("external_trigger_id", tfModel.TriggerID.ValueString(), "trigger_id").
		SetStringField("alarm_resource_query", tfModel.AlarmResourceQuery.ValueString(), "alarm_resource_query").
		SetStringField("communication_workspace", tfModel.CommunicationWorkspace.ValueString(), "communication_workspace").
		SetStringField("communication_channel", tfModel.CommunicationChannel.ValueString(), "communication_channel").
		SetStringField("integration_name", tfModel.IntegrationName.ValueString(), "integration_name").
		Build()
}
