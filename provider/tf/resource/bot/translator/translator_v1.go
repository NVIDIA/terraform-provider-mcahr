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

	"terraform/terraform-provider/provider/common"
	botapi "terraform/terraform-provider/provider/external_api/resources/bots"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	bottf "terraform/terraform-provider/provider/tf/resource/bot/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BotTranslatorV1 handles translation between TF models and V1 API models for bot resources
type BotTranslatorV1 struct {
	BotTranslatorCommon
}

var _ translator.Translator[*bottf.BotTFModel, *botapi.BotResponseAPIModelV1] = &BotTranslatorV1{}

// ToTFModel converts a V1 API model to a TF model
func (t *BotTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *botapi.BotResponseAPIModelV1) (*bottf.BotTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the bot container regardless of operation type (define_bot, update_bot, get_bot_class, delete_bot)
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no bot container found in V1 API response")
	}

	if len(container.BotClasses) == 0 {
		return nil, fmt.Errorf("no bot classes found in V1 API response")
	}

	// Get the first bot class, current implementation only supports one bot to be returned by the API
	botClass := container.BotClasses[0]

	// Build TF model from V1 bot class
	tfModel := &bottf.BotTFModel{
		Name:                   types.StringValue(botClass.Name),
		Command:                types.StringValue(t.buildCommand(botClass)),
		Description:            types.StringValue(botClass.Description),
		Enabled:                types.BoolValue(botClass.Enabled),
		Family:                 types.StringValue(botClass.ConfigData.Family),
		TriggerSource:          types.StringValue(botClass.TriggerSource),
		TriggerID:              types.StringValue(botClass.ExternalTriggerID),
		AlarmResourceQuery:     types.StringValue(botClass.AlarmResourceQuery),
		CommunicationWorkspace: types.StringValue(botClass.Communication.Workspace),
		CommunicationChannel:   types.StringValue(botClass.Communication.Channel),
		IntegrationName:        types.StringValue(botClass.IntegrationName),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (t *BotTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *bottf.BotTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}

// buildCommand reconstructs the command from alarm and action statements
func (t *BotTranslatorV1) buildCommand(botClass botapi.BotClassV1) string {
	return fmt.Sprintf("if %s then %s fi", botClass.AlarmStatement, botClass.ActionStatement)
}
