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

// ActionTranslatorV1 handles translation for ActionResponseAPIModelV1
type ActionTranslatorV1 struct {
	ActionTranslatorCommon
}

var _ translator.Translator[*actiontf.ActionTFModel, *actionapi.ActionResponseAPIModelV1] = &ActionTranslatorV1{}

func (a *ActionTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *actionapi.ActionResponseAPIModelV1) (*actiontf.ActionTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the action container regardless of operation type (define_action, update_action, get_action_class)
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no action container found in V1 API response")
	}

	if len(container.ActionClasses) == 0 {
		return nil, fmt.Errorf("no action classes found in V1 API response")
	}

	// Get the first action class, current implementation only supports one action to be returned by the API
	actionClass := container.ActionClasses[0]

	tfModel := &actiontf.ActionTFModel{
		Name:                   types.StringValue(actionClass.Name),
		Command:                types.StringValue(actionClass.Command),
		Enabled:                types.BoolValue(actionClass.Enabled),
		Timeout:                types.Int64Value(actionClass.Timeout),
		Description:            types.StringValue(actionClass.Description),
		ResEnvVar:              types.StringValue(actionClass.ResEnvVar),
		ResourceQuery:          types.StringValue(actionClass.ResourceQuery),
		Shell:                  types.StringValue(actionClass.Shell),
		AllowedResourcesQuery:  types.StringValue(actionClass.AllowedResourcesQuery),
		CommunicationWorkspace: types.StringValue(actionClass.Communication.Workspace),
		CommunicationChannel:   types.StringValue(actionClass.Communication.Channel),
		// Map step classes to template fields
		StartTitleTemplate:    types.StringValue(actionClass.StartStepClass.TitleTemplate),
		StartShortTemplate:    types.StringValue(actionClass.StartStepClass.ShortTemplate),
		ErrorTitleTemplate:    types.StringValue(actionClass.ErrorStepClass.TitleTemplate),
		ErrorShortTemplate:    types.StringValue(actionClass.ErrorStepClass.ShortTemplate),
		CompleteTitleTemplate: types.StringValue(actionClass.CompleteStepClass.TitleTemplate),
		CompleteShortTemplate: types.StringValue(actionClass.CompleteStepClass.ShortTemplate),
	}

	// Handle set fields
	ctx := requestContext.Context
	tfModel.Params, _ = types.ListValueFrom(ctx, types.StringType, utils.ParseStringArray(actionClass.Params))
	tfModel.ResourceTagsToExport, _ = types.ListValueFrom(ctx, types.StringType, utils.ParseStringArray(actionClass.ResourceTagsToExport))
	tfModel.FileDeps, _ = types.ListValueFrom(ctx, types.StringType, utils.ParseStringArray(actionClass.FileDeps))
	tfModel.AllowedEntities, _ = types.ListValueFrom(ctx, types.StringType, actionClass.AllowedEntities)
	tfModel.Editors, _ = types.ListValueFrom(ctx, types.StringType, actionClass.Editors)

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (a *ActionTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *actiontf.ActionTFModel) (*statement.StatementInputAPIModel, error) {
	return a.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
