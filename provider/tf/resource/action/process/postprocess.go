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

package actions

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	action "terraform/terraform-provider/provider/tf/resource/action/model"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type ActionPostProcessor struct{}

var _ process.PostProcessor[*action.ActionTFModel] = &ActionPostProcessor{}

func (p *ActionPostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, tf *action.ActionTFModel) error {
	// Preserve original long template values from config for deprecated fields
	return p.preserveOriginalLongTemplatesFromConfig(requestContext, data.CreateRequest.Config, tf)
}

func (p *ActionPostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, tf *action.ActionTFModel) error {
	// For read operations, preserve state values for deprecated long templates
	return p.preserveOriginalLongTemplatesFromState(requestContext, data.ReadRequest.State, tf)
}

func (p *ActionPostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, tf *action.ActionTFModel) error {
	// Preserve original long template values from config for deprecated fields
	return p.preserveOriginalLongTemplatesFromConfig(requestContext, data.UpdateRequest.Config, tf)
}

func (p *ActionPostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, tf *action.ActionTFModel) error {
	// No need to preserve values for delete operations
	return nil
}

// preserveOriginalLongTemplatesFromState extracts the original long template values from state
// and sets them back into the TF model to preserve user-configured values for deprecated fields
func (p *ActionPostProcessor) preserveOriginalLongTemplatesFromState(requestContext *common.RequestContext, source tfsdk.State, tf *action.ActionTFModel) error {
	var originalModel action.ActionTFModel

	// Get the original values from state
	diags := source.Get(requestContext.Context, &originalModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original values from state: %s", diags.Errors())
	}

	return p.applyLongTemplateValues(&originalModel, tf)
}

// preserveOriginalLongTemplatesFromConfig extracts the original long template values from config
// and sets them back into the TF model to preserve user-configured values for deprecated fields
func (p *ActionPostProcessor) preserveOriginalLongTemplatesFromConfig(requestContext *common.RequestContext, source tfsdk.Config, tf *action.ActionTFModel) error {
	var originalModel action.ActionTFModel

	// Get the original values from config
	diags := source.Get(requestContext.Context, &originalModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original values from config: %s", diags.Errors())
	}

	return p.applyLongTemplateValues(&originalModel, tf)
}

// applyLongTemplateValues copies the deprecated long template values from source to target model
// if they are not null in the source
func (p *ActionPostProcessor) applyLongTemplateValues(sourceModel, targetModel *action.ActionTFModel) error {
	// Preserve the original deprecated long template values if they were set by user
	if !sourceModel.StartLongTemplate.IsNull() {
		targetModel.StartLongTemplate = sourceModel.StartLongTemplate
	}

	if !sourceModel.ErrorLongTemplate.IsNull() {
		targetModel.ErrorLongTemplate = sourceModel.ErrorLongTemplate
	}

	if !sourceModel.CompleteLongTemplate.IsNull() {
		targetModel.CompleteLongTemplate = sourceModel.CompleteLongTemplate
	}

	return nil
}
