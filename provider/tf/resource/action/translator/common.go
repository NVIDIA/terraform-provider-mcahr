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
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	actiontf "terraform/terraform-provider/provider/tf/resource/action/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ActionTranslatorCommon contains shared functionality between V1 and V2 translators
type ActionTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (a *ActionTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *actiontf.ActionTFModel) (*statement.InputAPIModel, error) {
	var stmt string

	var err error

	switch requestContext.Operation {
	case common.Create:
		stmt, err = a.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = a.buildReadStatement(tfModel)
	case common.Update:
		stmt, err = a.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt = a.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}
	if err != nil {
		return nil, err
	}

	apiModel := &statement.InputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (a *ActionTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *actiontf.ActionTFModel) (string, error) {
	return a.buildActionStatement(requestContext, translationData, "define_action", tfModel)
}

func (a *ActionTranslatorCommon) buildReadStatement(tfModel *actiontf.ActionTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_action_class(action_name=%s)", utils.EscapeString(name))
}

func (a *ActionTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *actiontf.ActionTFModel) (string, error) {
	return a.buildActionStatement(requestContext, translationData, "update_action", tfModel)
}

func (a *ActionTranslatorCommon) buildDeleteStatement(tfModel *actiontf.ActionTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_action(action_name=%s)", utils.EscapeString(name))
}

func (a *ActionTranslatorCommon) buildActionStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *actiontf.ActionTFModel) (string, error) {
	// Build the action statement from the TF model using the builder pattern
	// Used for both define_action (create) and update_action (update) operations

	ctx := requestContext.Context

	// Extract the list attributes up front. A conversion failure must abort
	// the apply: silently sending an empty allowed_entities or editors list
	// would clear the action's authorization instead.
	lists, err := utils.ListSlicesFromTFModel(ctx, map[string]types.List{
		"params":                  tfModel.Params,
		"resource_tags_to_export": tfModel.ResourceTagsToExport,
		"file_deps":               tfModel.FileDeps,
		"allowed_entities":        tfModel.AllowedEntities,
		"editors":                 tfModel.Editors,
	})
	if err != nil {
		return "", err
	}

	return utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("action_name", tfModel.Name.ValueString(), "name").
		SetStringField("command", tfModel.Command.ValueString(), "command").
		SetField("enabled", tfModel.Enabled.ValueBool(), "enabled").
		SetField("timeout", tfModel.Timeout.ValueInt64(), "timeout").
		SetStringField("description", tfModel.Description.ValueString(), "description").
		SetStringField("res_env_var", tfModel.ResEnvVar.ValueString(), "res_env_var").
		SetStringField("resource_query", tfModel.ResourceQuery.ValueString(), "resource_query").
		SetStringField("shell", tfModel.Shell.ValueString(), "shell").
		SetStringField("allowed_resources_query", tfModel.AllowedResourcesQuery.ValueString(), "allowed_resources_query").
		SetStringField("communication_workspace", tfModel.CommunicationWorkspace.ValueString(), "communication_workspace").
		SetStringField("communication_channel", tfModel.CommunicationChannel.ValueString(), "communication_channel").
		SetStringField("start_title_template", tfModel.StartTitleTemplate.ValueString(), "start_title_template").
		SetStringField("start_short_template", tfModel.StartShortTemplate.ValueString(), "start_short_template").
		SetStringField("complete_title_template", tfModel.CompleteTitleTemplate.ValueString(), "complete_title_template").
		SetStringField("complete_short_template", tfModel.CompleteShortTemplate.ValueString(), "complete_short_template").
		SetStringField("error_title_template", tfModel.ErrorTitleTemplate.ValueString(), "error_title_template").
		SetStringField("error_short_template", tfModel.ErrorShortTemplate.ValueString(), "error_short_template").
		SetArrayField("params", lists["params"], "params").
		SetArrayField("resource_tags_to_export", lists["resource_tags_to_export"], "resource_tags_to_export").
		SetArrayField("file_deps", lists["file_deps"], "file_deps").
		SetArrayField("allowed_entities", lists["allowed_entities"], "allowed_entities").
		SetArrayField("editors", lists["editors"], "editors").
		Build(), nil
}
