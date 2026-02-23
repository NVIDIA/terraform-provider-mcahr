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
	dashboardtf "terraform/terraform-provider/provider/tf/resource/dashboard/model"
)

type DashboardPreProcessor struct {
	base process.BasePreProcessor[*dashboardtf.DashboardTFModel]
}

var _ process.PreProcessor[*dashboardtf.DashboardTFModel] = &DashboardPreProcessor{}

func (p *DashboardPreProcessor) PreProcessCreate(requestContext *common.RequestContext, data *process.ProcessData) (*dashboardtf.DashboardTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.CreateRequest.Plan, &dashboardtf.DashboardTFModel{})
}

func (p *DashboardPreProcessor) PreProcessRead(requestContext *common.RequestContext, data *process.ProcessData) (*dashboardtf.DashboardTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.ReadRequest.State, &dashboardtf.DashboardTFModel{})
}

func (p *DashboardPreProcessor) PreProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData) (*dashboardtf.DashboardTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.UpdateRequest.Plan, &dashboardtf.DashboardTFModel{})
}

func (p *DashboardPreProcessor) PreProcessDelete(requestContext *common.RequestContext, data *process.ProcessData) (*dashboardtf.DashboardTFModel, error) {
	return p.base.ExtractForDelete(requestContext, data, &dashboardtf.DashboardTFModel{})
}
