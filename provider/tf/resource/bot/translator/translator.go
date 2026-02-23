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

type BotTranslator struct {
	BotTranslatorCommon
}

var _ translator.Translator[*bottf.BotTFModel, *botapi.BotResponseAPIModel] = &BotTranslator{}

// buildCommand reconstructs the command from alarm and action statements
func (b *BotTranslator) buildCommand(alarmStatement, actionStatement string) string {
	return fmt.Sprintf("if %s then %s fi", alarmStatement, actionStatement)
}

func (b *BotTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *botapi.BotResponseAPIModel) (*bottf.BotTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no bot configurations found in API response")
	}

	// Get the first configuration item, current implementation only supports one bot to be returned by the API
	configItem := apiModel.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	// Extract alarm and action statements from the compound command
	alarmStatement := config.TriggerEntityName
	actionStatement := config.ExecutionEntityName

	// Reconstruct the command from statements
	command := b.buildCommand(alarmStatement, actionStatement)

	tfModel := &bottf.BotTFModel{
		Name:                   types.StringValue(metadata.Name),
		Command:                types.StringValue(command),
		Description:            types.StringValue(metadata.Description),
		Enabled:                types.BoolValue(metadata.Enabled),
		Family:                 types.StringValue(metadata.Family),
		TriggerSource:          types.StringValue(config.TriggerSource),
		TriggerID:              types.StringValue(config.TriggerID),
		AlarmResourceQuery:     types.StringValue(config.AlarmResourceQuery),
		CommunicationWorkspace: types.StringValue(config.CommunicationDest.Workspace),
		CommunicationChannel:   types.StringValue(config.CommunicationDest.Channel),
		IntegrationName:        types.StringValue(config.IntegrationName),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V2 backend
func (b *BotTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *bottf.BotTFModel) (*statement.StatementInputAPIModel, error) {
	return b.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
