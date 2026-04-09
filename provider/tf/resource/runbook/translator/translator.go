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
	runbookapi "terraform/terraform-provider/provider/external_api/resources/runbooks"
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"
	converters "terraform/terraform-provider/provider/tf/resource/runbook/translator/object_converters"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RunbookTranslator struct {
	RunbookTranslatorCommon
}

var _ translator.Translator[*runbooktf.RunbookTFModel, *runbookapi.RunbookResponseAPIModel] = &RunbookTranslator{}

func (r *RunbookTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *runbookapi.RunbookResponseAPIModel) (*runbooktf.RunbookTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no configurations found in API response")
	}

	// Get the first configuration item, current implementation only supports one runbook to be returned by the API
	configItem := apiModel.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	// Handle base fields
	tfModel := &runbooktf.RunbookTFModel{
		Name:                                types.StringValue(metadata.Name),
		Enabled:                             types.BoolValue(metadata.Enabled),
		Description:                         types.StringValue(metadata.Description),
		TimeoutMs:                           types.Int64Value(config.TimeoutMs),
		AllowedResourcesQuery:               types.StringValue(config.AllowedResourcesQuery),
		CommunicationWorkspace:              types.StringValue(config.CommunicationDestination.Workspace),
		CommunicationChannel:                types.StringValue(config.CommunicationDestination.Channel),
		Category:                            types.StringValue(config.Category),
		IsRunOutputPersisted:                types.BoolValue(config.IsRunOutputPersisted),
		FilterResourceToAction:              types.BoolValue(config.FilterResourceToAction),
		CommunicationCudNotifications:       types.BoolValue(config.CommunicationFilters.CudNotifications),
		CommunicationApprovalNotifications:  types.BoolValue(config.CommunicationFilters.ApprovalNotifications),
		CommunicationExecutionNotifications: types.BoolValue(config.CommunicationFilters.ExecutionNotifications),
	}

	// Handle JSON fields
	err := toTFModelJsonFields(tfModel, config.Cells, config.Params, config.ExternalParams)
	if err != nil {
		return nil, err
	}

	// Handle set fields
	tfModel.AllowedEntities, _ = types.ListValueFrom(requestContext.Context, types.StringType, config.AllowedEntities)
	tfModel.Approvers, _ = types.ListValueFrom(requestContext.Context, types.StringType, config.Approvers)
	tfModel.Labels, _ = types.ListValueFrom(requestContext.Context, types.StringType, config.Labels)
	tfModel.Editors, _ = types.ListValueFrom(requestContext.Context, types.StringType, config.Editors)
	tfModel.SecretNames, _ = types.ListValueFrom(requestContext.Context, types.StringType, config.SecretNames)

	// Handle object fields
	tfModel.ParamsGroups, _ = converters.ParamsGroupsToTFModel(requestContext, config.ParamsGroups)

	return tfModel, nil
}

func toTFModelJsonFields(tfModel *runbooktf.RunbookTFModel, cells []customattribute.CellJsonAPI, params []customattribute.ParamJson, externalParams []customattribute.ExternalParamJson) error {
	// The original attributes (i.e. Cells, Params, ExternalParams) can be overridden by the postprocessor
	// but if config is empty then this API response value is used

	// Cells
	cellsJson, err := customattribute.MapCellsToInternalModel(cells)
	if err != nil {
		return fmt.Errorf("failed to map cells to internal model: %v", err)
	}
	tfModel.Cells = types.StringValue(cellsJson)
	tfModel.CellsFull = types.StringValue(cellsJson)

	// CellsList
	cellsList, diags := converters.CellsListFromAPICells(cells)
	if diags.HasError() {
		return fmt.Errorf("failed to convert cells to cells_list: %s", diags.Errors())
	}
	tfModel.CellsList = cellsList

	// Params
	paramsJson, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %v", err)
	}
	tfModel.Params = types.StringValue(string(paramsJson))
	tfModel.ParamsFull = types.StringValue(string(paramsJson))

	// ParamsList
	paramsList, paramsDiags := converters.ParamsListFromAPI(params)
	if paramsDiags.HasError() {
		return fmt.Errorf("failed to convert params to params_list: %s", paramsDiags.Errors())
	}
	tfModel.ParamsList = paramsList

	// External Params
	externalParamsJson, err := json.Marshal(externalParams)
	if err != nil {
		return fmt.Errorf("failed to marshal external_params: %v", err)
	}
	tfModel.ExternalParams = types.StringValue(string(externalParamsJson))
	tfModel.ExternalParamsFull = types.StringValue(string(externalParamsJson))

	// ExternalParamsList
	externalParamsList, epDiags := converters.ExternalParamsListFromAPI(externalParams)
	if epDiags.HasError() {
		return fmt.Errorf("failed to convert external_params to external_params_list: %s", epDiags.Errors())
	}
	tfModel.ExternalParamsList = externalParamsList

	return nil
}

func (r *RunbookTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *runbooktf.RunbookTFModel) (apiModel *statement.StatementInputAPIModel, err error) {
	return r.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
