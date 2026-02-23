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
	runbookapi "terraform/terraform-provider/provider/external_api/resources/runbooks"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RunbookTranslatorV1 handles translation for RunbookResponseAPIModelV1
type RunbookTranslatorV1 struct {
	RunbookTranslatorCommon
}

var _ translator.Translator[*runbooktf.RunbookTFModel, *runbookapi.RunbookResponseAPIModelV1] = &RunbookTranslatorV1{}

func (r *RunbookTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *runbookapi.RunbookResponseAPIModelV1) (*runbooktf.RunbookTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	// Get the runbook container regardless of operation type (define_notebook, update_notebook, get_notebook_class)
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no runbook container found in V1 API response")
	}

	if len(container.NotebookClasses) == 0 {
		return nil, fmt.Errorf("no runbook classes found in V1 API response")
	}

	// Get the first notebook class, current implementation only supports one runbook to be returned by the API
	notebookClass := container.NotebookClasses[0]

	// Handle base fields
	tfModel := &runbooktf.RunbookTFModel{
		Name:                                types.StringValue(notebookClass.Name),
		Enabled:                             types.BoolValue(notebookClass.Enabled),
		Description:                         types.StringValue(notebookClass.Description),
		TimeoutMs:                           types.Int64Value(notebookClass.TimeoutMs),
		AllowedResourcesQuery:               types.StringValue(notebookClass.AllowedResourcesQuery),
		CommunicationWorkspace:              types.StringValue(notebookClass.CommunicationWorkspace),
		CommunicationChannel:                types.StringValue(notebookClass.CommunicationChannel),
		Category:                            types.StringValue(notebookClass.Category),
		IsRunOutputPersisted:                types.BoolValue(notebookClass.IsRunOutputPersisted),
		FilterResourceToAction:              types.BoolValue(notebookClass.FilterResourceToAction),
		CommunicationCudNotifications:       types.BoolValue(notebookClass.CommunicationCudNotifications),
		CommunicationApprovalNotifications:  types.BoolValue(notebookClass.CommunicationApprovalNotifications),
		CommunicationExecutionNotifications: types.BoolValue(notebookClass.CommunicationExecutionNotifications),
	}

	// Handle JSON fields
	err := toTFModelJsonFields(tfModel, notebookClass.Cells, notebookClass.Params, notebookClass.ExternalParams)
	if err != nil {
		return nil, err
	}

	// Handle set fields
	tfModel.AllowedEntities, _ = types.ListValueFrom(requestContext.Context, types.StringType, notebookClass.AllowedEntities)
	tfModel.Approvers, _ = types.ListValueFrom(requestContext.Context, types.StringType, notebookClass.Approvers)
	tfModel.Labels, _ = types.ListValueFrom(requestContext.Context, types.StringType, notebookClass.Labels)
	tfModel.Editors, _ = types.ListValueFrom(requestContext.Context, types.StringType, notebookClass.Editors)
	tfModel.SecretNames, _ = types.ListValueFrom(requestContext.Context, types.StringType, notebookClass.SecretNames)

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (r *RunbookTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *runbooktf.RunbookTFModel) (*statement.StatementInputAPIModel, error) {
	return r.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
