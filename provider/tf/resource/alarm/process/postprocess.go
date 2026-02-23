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

package alarms

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	alarm "terraform/terraform-provider/provider/tf/resource/alarm/model"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type AlarmPostProcessor struct{}

var _ process.PostProcessor[*alarm.AlarmTFModel] = &AlarmPostProcessor{}

func (p *AlarmPostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, tf *alarm.AlarmTFModel) error {
	// Preserve original deprecated field values from config
	return p.preserveOriginalDeprecatedFieldsFromConfig(requestContext, data.CreateRequest.Config, tf)
}

func (p *AlarmPostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, tf *alarm.AlarmTFModel) error {
	// For read operations, preserve state values for deprecated fields
	return p.preserveOriginalDeprecatedFieldsFromState(requestContext, data.ReadRequest.State, tf)
}

func (p *AlarmPostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, tf *alarm.AlarmTFModel) error {
	// Preserve original deprecated field values from config
	return p.preserveOriginalDeprecatedFieldsFromConfig(requestContext, data.UpdateRequest.Config, tf)
}

func (p *AlarmPostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, tf *alarm.AlarmTFModel) error {
	// No need to preserve values for delete operations
	return nil
}

// preserveOriginalDeprecatedFieldsFromState extracts the original deprecated field values from state
// and sets them back into the TF model to preserve user-configured values for deprecated fields
func (p *AlarmPostProcessor) preserveOriginalDeprecatedFieldsFromState(requestContext *common.RequestContext, state tfsdk.State, tf *alarm.AlarmTFModel) error {
	var originalModel alarm.AlarmTFModel

	// Get the original values from state
	diags := state.Get(requestContext.Context, &originalModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original values from state: %s", diags.Errors())
	}

	return p.applyDeprecatedFieldValues(&originalModel, tf)
}

// preserveOriginalDeprecatedFieldsFromConfig extracts the original deprecated field values from config
// and sets them back into the TF model to preserve user-configured values for deprecated fields
func (p *AlarmPostProcessor) preserveOriginalDeprecatedFieldsFromConfig(requestContext *common.RequestContext, config tfsdk.Config, tf *alarm.AlarmTFModel) error {
	var originalModel alarm.AlarmTFModel

	// Get the original values from config
	diags := config.Get(requestContext.Context, &originalModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original values from config: %s", diags.Errors())
	}

	return p.applyDeprecatedFieldValues(&originalModel, tf)
}

// applyDeprecatedFieldValues copies all deprecated field values from source to target model
// if they are not null in the source
func (p *AlarmPostProcessor) applyDeprecatedFieldValues(sourceModel, targetModel *alarm.AlarmTFModel) error {
	// Preserve the original deprecated field values if they were set by user

	// Deprecated query fields
	if !sourceModel.MuteQuery.IsNull() {
		targetModel.MuteQuery = sourceModel.MuteQuery
	}

	if !sourceModel.RaiseFor.IsNull() {
		targetModel.RaiseFor = sourceModel.RaiseFor
	}

	// Deprecated condition fields
	if !sourceModel.MetricName.IsNull() {
		targetModel.MetricName = sourceModel.MetricName
	}

	if !sourceModel.ConditionType.IsNull() {
		targetModel.ConditionType = sourceModel.ConditionType
	}

	if !sourceModel.ConditionValue.IsNull() {
		targetModel.ConditionValue = sourceModel.ConditionValue
	}

	// Deprecated template fields
	if !sourceModel.FireLongTemplate.IsNull() {
		targetModel.FireLongTemplate = sourceModel.FireLongTemplate
	}

	if !sourceModel.ResolveLongTemplate.IsNull() {
		targetModel.ResolveLongTemplate = sourceModel.ResolveLongTemplate
	}

	return nil
}
