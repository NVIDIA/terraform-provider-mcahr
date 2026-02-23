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
	actionapi "terraform/terraform-provider/provider/external_api/resources/actions"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	actiontf "terraform/terraform-provider/provider/tf/resource/action/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ActionTranslator struct {
	ActionTranslatorCommon
}

var _ translator.Translator[*actiontf.ActionTFModel, *actionapi.ActionResponseAPIModel] = &ActionTranslator{}

func (a *ActionTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *actionapi.ActionResponseAPIModel) (*actiontf.ActionTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	if apiModel.Output.Configurations.Count == 0 || len(apiModel.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no configurations found in API response")
	}

	// Get the first configuration item, current implementation only supports one action to be returned by the API
	configItem := apiModel.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	tfModel := &actiontf.ActionTFModel{
		Name:                   types.StringValue(metadata.Name),
		Command:                types.StringValue(config.CommandText),
		Enabled:                types.BoolValue(metadata.Enabled),
		Timeout:                types.Int64Value(config.Timeout),
		Description:            types.StringValue(metadata.Description),
		ResEnvVar:              types.StringValue(config.ResEnvVar),
		ResourceQuery:          types.StringValue(config.ResourceQuery),
		Shell:                  types.StringValue(config.Shell),
		AllowedResourcesQuery:  types.StringValue(config.AllowedResourcesQuery),
		CommunicationWorkspace: types.StringValue(config.CommunicationDest.Workspace),
		CommunicationChannel:   types.StringValue(config.CommunicationDest.Channel),
		// Map step details to template fields
		StartTitleTemplate:    types.StringValue(config.StepDetails.StartStep.Title),
		StartShortTemplate:    types.StringValue(config.StepDetails.StartStep.Description),
		ErrorTitleTemplate:    types.StringValue(config.StepDetails.ErrorStep.Title),
		ErrorShortTemplate:    types.StringValue(config.StepDetails.ErrorStep.Description),
		CompleteTitleTemplate: types.StringValue(config.StepDetails.CompleteStep.Title),
		CompleteShortTemplate: types.StringValue(config.StepDetails.CompleteStep.Description),
	}

	// Handle set fields
	ctx := requestContext.Context
	tfModel.Params, _ = types.ListValueFrom(ctx, types.StringType, utils.ParseStringArray(config.Params))
	tfModel.ResourceTagsToExport, _ = types.ListValueFrom(ctx, types.StringType, utils.ParseStringArray(config.ResourceTagsToExport))
	tfModel.FileDeps, _ = types.ListValueFrom(ctx, types.StringType, utils.ParseStringArray(config.FileDeps))
	tfModel.AllowedEntities, _ = types.ListValueFrom(ctx, types.StringType, config.AllowedEntities)
	tfModel.Editors, _ = types.ListValueFrom(ctx, types.StringType, config.Editors)

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V2 backend
func (a *ActionTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *actiontf.ActionTFModel) (*statement.StatementInputAPIModel, error) {
	return a.ToAPIModelWithVersion(requestContext, translationData, tfModel)

}
