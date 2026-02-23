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

type TimeTriggerPostProcessor struct{}

var _ process.PostProcessor[*timetrigger.TimeTriggerTFModel] = &TimeTriggerPostProcessor{}

func (p *TimeTriggerPostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, model *timetrigger.TimeTriggerTFModel) error {
	// No additional processing needed for time trigger create
	return nil
}

func (p *TimeTriggerPostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, model *timetrigger.TimeTriggerTFModel) error {
	// No additional processing needed for time trigger read
	return nil
}

func (p *TimeTriggerPostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, model *timetrigger.TimeTriggerTFModel) error {
	// No additional processing needed for time trigger update
	return nil
}

func (p *TimeTriggerPostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, model *timetrigger.TimeTriggerTFModel) error {
	// No additional processing needed for time trigger delete
	return nil
}
