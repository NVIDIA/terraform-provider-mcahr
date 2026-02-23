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

package process

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	bottf "terraform/terraform-provider/provider/tf/resource/bot/model"
)

type BotPostProcessor struct{}

var _ process.PostProcessor[*bottf.BotTFModel] = &BotPostProcessor{}

func (p *BotPostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, result *bottf.BotTFModel) error {
	return setDeprecatedFields(requestContext, data.CreateRequest.Config, result)
}

func (p *BotPostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, result *bottf.BotTFModel) error {
	return setDeprecatedFields(requestContext, data.ReadRequest.State, result)
}

func (p *BotPostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, result *bottf.BotTFModel) error {
	return setDeprecatedFields(requestContext, data.UpdateRequest.Config, result)
}

func (p *BotPostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, result *bottf.BotTFModel) error {
	return nil
}

// If the user has deprecated fields set, then set them to the value returned by the API
// Otherwise, leave them as is (Null)
func setDeprecatedFields(requestContext *common.RequestContext, getter process.Getter, tfModel *bottf.BotTFModel) error {

	// Get the original values from config/state (before calling the API)
	var configModel bottf.BotTFModel
	diags := getter.Get(requestContext.Context, &configModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original model: %s", diags.Errors())
	}

	// Set deprecated fields to the value returned by the API

	if common.IsAttrKnown(configModel.EventType) {
		tfModel.EventType = tfModel.TriggerSource
	}

	if common.IsAttrKnown(configModel.MonitorID) {
		tfModel.MonitorID = tfModel.TriggerID
	}

	return nil
}
