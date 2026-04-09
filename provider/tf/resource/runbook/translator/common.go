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
	"encoding/json"
	"fmt"

	"terraform/terraform-provider/provider/common"
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"
	converters "terraform/terraform-provider/provider/tf/resource/runbook/translator/object_converters"
)

// ActionTranslatorCommon contains shared functionality between V1 and V2 translators
type RunbookTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (r *RunbookTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *runbooktf.RunbookTFModel) (apiModel *statement.StatementInputAPIModel, err error) {
	var stmt string

	switch requestContext.Operation {
	case common.Create:
		stmt, err = r.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = r.buildReadStatement(tfModel)
	case common.Update:
		stmt, err = r.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt = r.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}

	if err != nil {
		return nil, err
	}

	apiModel = &statement.StatementInputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (r *RunbookTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *runbooktf.RunbookTFModel) (string, error) {
	return r.buildRunbookStatement(requestContext, translationData, "define_notebook", tfModel)
}

func (r *RunbookTranslatorCommon) buildReadStatement(tfModel *runbooktf.RunbookTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_notebook_class(notebook_name=\"%s\")", name)
}

func (r *RunbookTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *runbooktf.RunbookTFModel) (string, error) {
	return r.buildRunbookStatement(requestContext, translationData, "update_notebook", tfModel)
}

func (r *RunbookTranslatorCommon) buildDeleteStatement(tfModel *runbooktf.RunbookTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_notebook(notebook_name=\"%s\")", name)
}

func (r *RunbookTranslatorCommon) buildRunbookStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *runbooktf.RunbookTFModel) (string, error) {
	// Build the runbook statement from the TF model using the builder pattern
	// Used for both define_notebook (create) and update_notebook (update) operations

	ctx := requestContext.Context
	// Handle base and set fields
	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("notebook_name", tfModel.Name.ValueString(), "name").
		SetField("enabled", tfModel.Enabled.ValueBool(), "enabled").
		SetField("timeout_ms", tfModel.TimeoutMs.ValueInt64(), "timeout_ms").
		SetStringField("description", tfModel.Description.ValueString(), "description").
		SetStringField("allowed_resources_query", tfModel.AllowedResourcesQuery.ValueString(), "allowed_resources_query").
		SetStringField("communication_workspace", tfModel.CommunicationWorkspace.ValueString(), "communication_workspace").
		SetStringField("communication_channel", tfModel.CommunicationChannel.ValueString(), "communication_channel").
		SetStringField("category", tfModel.Category.ValueString(), "category").
		SetField("is_run_output_persisted", tfModel.IsRunOutputPersisted.ValueBool(), "is_run_output_persisted").
		SetField("filter_resource_to_action", tfModel.FilterResourceToAction.ValueBool(), "filter_resource_to_action").
		SetField("communication_cud_notifications", tfModel.CommunicationCudNotifications.ValueBool(), "communication_cud_notifications").
		SetField("communication_approval_notifications", tfModel.CommunicationApprovalNotifications.ValueBool(), "communication_approval_notifications").
		SetField("communication_execution_notifications", tfModel.CommunicationExecutionNotifications.ValueBool(), "communication_execution_notifications").
		SetArrayField("allowed_entities", utils.ListSliceFromTFModel(ctx, tfModel.AllowedEntities), "allowed_entities").
		SetArrayField("approvers", utils.ListSliceFromTFModel(ctx, tfModel.Approvers), "approvers").
		SetArrayField("labels", utils.ListSliceFromTFModel(ctx, tfModel.Labels), "labels").
		SetArrayField("editors", utils.ListSliceFromTFModel(ctx, tfModel.Editors), "editors").
		SetArrayField("secret_names", utils.ListSliceFromTFModel(ctx, tfModel.SecretNames), "secret_names")

	jsonParamsGroups, err := buildParamsGroupsJSON(requestContext, tfModel)
	if err != nil {
		return "", fmt.Errorf("failed to build params_groups JSON: %v", err)
	}
	builder.SetField("params_groups", jsonParamsGroups, "params_groups")

	apiCells, err := buildCellsForStatement(requestContext, tfModel)
	if err != nil {
		return "", fmt.Errorf("failed to build cells for statement: %v", err)
	}
	builder.SetField("cells", apiCells, "cells")

	apiParams, err := buildParamsForStatement(requestContext, tfModel)
	if err != nil {
		return "", fmt.Errorf("failed to build params for statement: %v", err)
	}
	builder.SetField("params", apiParams, "params")

	apiExternalParams, err := buildExternalParamsForStatement(requestContext, tfModel)
	if err != nil {
		return "", fmt.Errorf("failed to build external_params for statement: %v", err)
	}
	builder.SetField("external_params", apiExternalParams, "external_params")

	return builder.Build(), nil
}

func buildCellsForStatement(requestContext *common.RequestContext, tfModel *runbooktf.RunbookTFModel) (string, error) {
	if !tfModel.CellsList.IsNull() && !tfModel.CellsList.IsUnknown() {
		internalCells, err := converters.CellsListToInternalCells(requestContext.Context, tfModel.CellsList)
		if err != nil {
			return "", fmt.Errorf("failed to convert cells_list to internal model: %v", err)
		}
		return customattribute.InternalCellsToBase64APIModel(requestContext, internalCells)
	}

	return customattribute.MapCellsToAPIModel(requestContext, tfModel.Cells.ValueString())
}

func buildParamsForStatement(requestContext *common.RequestContext, tfModel *runbooktf.RunbookTFModel) (string, error) {
	if !tfModel.ParamsList.IsNull() && !tfModel.ParamsList.IsUnknown() {
		params, err := converters.ParamsListToInternal(requestContext.Context, tfModel.ParamsList)
		if err != nil {
			return "", fmt.Errorf("failed to convert params_list to internal model: %v", err)
		}
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return "", fmt.Errorf("failed to marshal params: %v", err)
		}
		return string(jsonBytes), nil
	}
	return tfModel.ParamsFull.ValueString(), nil
}

func buildExternalParamsForStatement(requestContext *common.RequestContext, tfModel *runbooktf.RunbookTFModel) (string, error) {
	if !tfModel.ExternalParamsList.IsNull() && !tfModel.ExternalParamsList.IsUnknown() {
		params, err := converters.ExternalParamsListToInternal(requestContext.Context, tfModel.ExternalParamsList)
		if err != nil {
			return "", fmt.Errorf("failed to convert external_params_list to internal model: %v", err)
		}
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return "", fmt.Errorf("failed to marshal external_params: %v", err)
		}
		return string(jsonBytes), nil
	}
	return tfModel.ExternalParamsFull.ValueString(), nil
}

func buildParamsGroupsJSON(requestContext *common.RequestContext, tfModel *runbooktf.RunbookTFModel) (string, error) {

	apiParamsGroups, err := converters.ParamsGroupsFromTFModel(requestContext, tfModel.ParamsGroups)
	if err != nil {
		return "", fmt.Errorf("failed to convert params_groups: %v", err)
	}

	jsonParamsGroups, err := json.Marshal(apiParamsGroups)
	if err != nil {
		return "", fmt.Errorf("failed to marshal params_groups: %v", err)
	}
	return string(jsonParamsGroups), nil
}
