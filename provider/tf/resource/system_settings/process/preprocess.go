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

// SystemSettingsPreProcessor handles preprocessing for system_settings resources
type SystemSettingsPreProcessor struct {
	process.BasePreProcessor[*systemsettingstf.SystemSettingsTFModel]
}

// PreProcessCreate performs preprocessing for create operations
func (p *SystemSettingsPreProcessor) PreProcessCreate(requestContext *common.RequestContext, data *process.ProcessData) (*systemsettingstf.SystemSettingsTFModel, error) {
	return p.ExtractFrom(requestContext, data.CreateRequest.Plan, &systemsettingstf.SystemSettingsTFModel{})
}

// PreProcessRead performs preprocessing for read operations
func (p *SystemSettingsPreProcessor) PreProcessRead(requestContext *common.RequestContext, data *process.ProcessData) (*systemsettingstf.SystemSettingsTFModel, error) {
	return p.ExtractFrom(requestContext, data.ReadRequest.State, &systemsettingstf.SystemSettingsTFModel{})
}

// PreProcessUpdate performs preprocessing for update operations
func (p *SystemSettingsPreProcessor) PreProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData) (*systemsettingstf.SystemSettingsTFModel, error) {
	return p.ExtractFrom(requestContext, data.UpdateRequest.Plan, &systemsettingstf.SystemSettingsTFModel{})
}

// PreProcessDelete performs preprocessing for delete operations
func (p *SystemSettingsPreProcessor) PreProcessDelete(requestContext *common.RequestContext, data *process.ProcessData) (*systemsettingstf.SystemSettingsTFModel, error) {
	return p.ExtractForDelete(requestContext, data, &systemsettingstf.SystemSettingsTFModel{})
}
