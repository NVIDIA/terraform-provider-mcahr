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

package time_triggers

import (
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	timetrigger "terraform/terraform-provider/provider/tf/resource/time_trigger/model"
)

type TimeTriggerPreProcessor struct {
	base process.BasePreProcessor[*timetrigger.TimeTriggerTFModel]
}

var _ process.PreProcessor[*timetrigger.TimeTriggerTFModel] = &TimeTriggerPreProcessor{}

func (p *TimeTriggerPreProcessor) PreProcessCreate(requestContext *common.RequestContext, data *process.ProcessData) (*timetrigger.TimeTriggerTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.CreateRequest.Plan, &timetrigger.TimeTriggerTFModel{})
}

func (p *TimeTriggerPreProcessor) PreProcessRead(requestContext *common.RequestContext, data *process.ProcessData) (*timetrigger.TimeTriggerTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.ReadRequest.State, &timetrigger.TimeTriggerTFModel{})
}

func (p *TimeTriggerPreProcessor) PreProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData) (*timetrigger.TimeTriggerTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.UpdateRequest.Plan, &timetrigger.TimeTriggerTFModel{})
}

func (p *TimeTriggerPreProcessor) PreProcessDelete(requestContext *common.RequestContext, data *process.ProcessData) (*timetrigger.TimeTriggerTFModel, error) {
	return p.base.ExtractForDelete(requestContext, data, &timetrigger.TimeTriggerTFModel{})
}
