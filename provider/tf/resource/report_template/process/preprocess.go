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
	reporttemplatetf "terraform/terraform-provider/provider/tf/resource/report_template/model"
)

type ReportTemplatePreProcessor struct {
	base process.BasePreProcessor[*reporttemplatetf.ReportTemplateTFModel]
}

var _ process.PreProcessor[*reporttemplatetf.ReportTemplateTFModel] = &ReportTemplatePreProcessor{}

func (p *ReportTemplatePreProcessor) PreProcessCreate(requestContext *common.RequestContext, data *process.ProcessData) (*reporttemplatetf.ReportTemplateTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.CreateRequest.Plan, &reporttemplatetf.ReportTemplateTFModel{})
}

func (p *ReportTemplatePreProcessor) PreProcessRead(requestContext *common.RequestContext, data *process.ProcessData) (*reporttemplatetf.ReportTemplateTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.ReadRequest.State, &reporttemplatetf.ReportTemplateTFModel{})
}

func (p *ReportTemplatePreProcessor) PreProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData) (*reporttemplatetf.ReportTemplateTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.UpdateRequest.Plan, &reporttemplatetf.ReportTemplateTFModel{})
}

func (p *ReportTemplatePreProcessor) PreProcessDelete(requestContext *common.RequestContext, data *process.ProcessData) (*reporttemplatetf.ReportTemplateTFModel, error) {
	return p.base.ExtractForDelete(requestContext, data, &reporttemplatetf.ReportTemplateTFModel{})
}
