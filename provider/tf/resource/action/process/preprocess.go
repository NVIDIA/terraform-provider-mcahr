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
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	action "terraform/terraform-provider/provider/tf/resource/action/model"
)

type ActionPreProcessor struct {
	base process.BasePreProcessor[*action.ActionTFModel]
}

var _ process.PreProcessor[*action.ActionTFModel] = &ActionPreProcessor{} // check that ActionPreProcessor implements PreProcessor

func (p *ActionPreProcessor) PreProcessCreate(requestContext *common.RequestContext, data *process.ProcessData) (*action.ActionTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.CreateRequest.Plan, &action.ActionTFModel{})
}

func (p *ActionPreProcessor) PreProcessRead(requestContext *common.RequestContext, data *process.ProcessData) (*action.ActionTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.ReadRequest.State, &action.ActionTFModel{})
}

func (p *ActionPreProcessor) PreProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData) (*action.ActionTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.UpdateRequest.Plan, &action.ActionTFModel{})
}

func (p *ActionPreProcessor) PreProcessDelete(requestContext *common.RequestContext, data *process.ProcessData) (*action.ActionTFModel, error) {
	return p.base.ExtractForDelete(requestContext, data, &action.ActionTFModel{})
}
