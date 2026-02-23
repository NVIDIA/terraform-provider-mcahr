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
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"
)

type FilePostProcessor struct{}

var _ process.PostProcessor[*filetf.FileTFModel] = &FilePostProcessor{}

func (p *FilePostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, result *filetf.FileTFModel) error {
	return setFieldsFromPrevious(requestContext, data, data.CreateRequest.Config, result)
}

func (p *FilePostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, result *filetf.FileTFModel) error {
	return setFieldsFromPrevious(requestContext, data, data.ReadRequest.State, result)
}

func (p *FilePostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, result *filetf.FileTFModel) error {
	return setFieldsFromPrevious(requestContext, data, data.UpdateRequest.Config, result)
}

func (p *FilePostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, result *filetf.FileTFModel) error {
	return nil
}

//
// Custom postprocess functions
//

func setFieldsFromPrevious(requestContext *common.RequestContext, data *process.ProcessData, config process.Getter, tfModel *filetf.FileTFModel) error {

	var configModel filetf.FileTFModel
	diags := config.Get(requestContext.Context, &configModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get config model: %s", diags.Errors())
	}

	// Restore values from plan/state for fields that are not returned by the backend
	tfModel.MD5 = configModel.MD5
	tfModel.InputFile = configModel.InputFile
	tfModel.InlineData = configModel.InlineData

	return nil
}
