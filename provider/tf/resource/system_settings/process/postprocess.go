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
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	systemsettingstf "terraform/terraform-provider/provider/tf/resource/system_settings/model"
)

// SystemSettingsPostProcessor handles postprocessing for system_settings resources
type SystemSettingsPostProcessor struct{}

var _ process.PostProcessor[*systemsettingstf.SystemSettingsTFModel] = &SystemSettingsPostProcessor{}

// PostProcessCreate performs postprocessing for create operations
func (p *SystemSettingsPostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *systemsettingstf.SystemSettingsTFModel) error {
	// For system_settings, we don't modify the model
	// State updates should happen in the resource, not in post processors as per requirements
	return nil
}

// PostProcessRead performs postprocessing for read operations
func (p *SystemSettingsPostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, tfModel *systemsettingstf.SystemSettingsTFModel) error {
	return nil
}

// PostProcessUpdate performs postprocessing for update operations
func (p *SystemSettingsPostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *systemsettingstf.SystemSettingsTFModel) error {
	return nil
}

// PostProcessDelete performs postprocessing for delete operations
func (p *SystemSettingsPostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, tfModel *systemsettingstf.SystemSettingsTFModel) error {
	return nil
}
