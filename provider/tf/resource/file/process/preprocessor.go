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
	corecommon "terraform/terraform-provider/provider/tf/core/common"
	"terraform/terraform-provider/provider/tf/core/process"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"
	"terraform/terraform-provider/provider/tf/resource/file/process/content"
	"terraform/terraform-provider/provider/tf/resource/file/process/upload"
)

type FilePreProcessor struct {
	base process.BasePreProcessor[*filetf.FileTFModel]
}

var _ process.PreProcessor[*filetf.FileTFModel] = &FilePreProcessor{}

func (p *FilePreProcessor) PreProcessCreate(requestContext *common.RequestContext, data *process.ProcessData) (*filetf.FileTFModel, error) {
	return p.handleCreateOrUpdate(requestContext, data, data.CreateRequest.Plan)
}

func (p *FilePreProcessor) PreProcessRead(requestContext *common.RequestContext, data *process.ProcessData) (*filetf.FileTFModel, error) {
	return p.base.ExtractFrom(requestContext, data.ReadRequest.State, &filetf.FileTFModel{})
}

func (p *FilePreProcessor) PreProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData) (*filetf.FileTFModel, error) {
	return p.handleCreateOrUpdate(requestContext, data, data.UpdateRequest.Plan)
}

func (p *FilePreProcessor) PreProcessDelete(requestContext *common.RequestContext, data *process.ProcessData) (*filetf.FileTFModel, error) {
	return p.base.ExtractForDelete(requestContext, data, &filetf.FileTFModel{})
}

//
// Custom preprocess functions
//

func (p *FilePreProcessor) handleCreateOrUpdate(requestContext *common.RequestContext, data *process.ProcessData, planGetter corecommon.Getter) (*filetf.FileTFModel, error) {

	model, err := p.base.ExtractFrom(requestContext, planGetter, &filetf.FileTFModel{})
	if err != nil {
		return nil, err
	}

	result, err := content.ProcessFileContent(requestContext, data, model)
	if err != nil {
		return nil, err
	}

	err = upload.HandleFileUploadProcess(requestContext, data, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
