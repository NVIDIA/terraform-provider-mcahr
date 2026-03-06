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
	"strconv"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"
	"terraform/terraform-provider/provider/tf/resource/file/process/content"
	"terraform/terraform-provider/provider/tf/resource/file/process/upload"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FilePreProcessor struct {
	base process.BasePreProcessor[*filetf.FileTFModel]
}

var _ process.PreProcessor[*filetf.FileTFModel] = &FilePreProcessor{}

func (p *FilePreProcessor) PreProcessCreate(requestContext *common.RequestContext, data *process.ProcessData) (*filetf.FileTFModel, error) {
	if requestContext.APIVersion == common.V1 {
		return p.handleV1Create(requestContext, data)
	}
	return p.handleCreateOrUpdate(requestContext, data, data.CreateRequest.Plan)
}

// handleV1Create handles file creation for V1 backends (< release-29.1.0).
// On V1 backends, the file object must exist before we can obtain a presigned URL,
// so we defer the upload to post-processing (after define_file creates the object).
// We also force enabled=false because V1 backends reject define_file without file_data
// when enabled=true.
func (p *FilePreProcessor) handleV1Create(requestContext *common.RequestContext, data *process.ProcessData) (*filetf.FileTFModel, error) {
	model, err := p.base.ExtractFrom(requestContext, data.CreateRequest.Plan, &filetf.FileTFModel{})
	if err != nil {
		return nil, err
	}

	result, err := content.ProcessFileContent(requestContext, data, model)
	if err != nil {
		return nil, err
	}

	data.StringArgs[V1DeferredUploadKey] = "true"
	data.StringArgs[V1OriginalEnabledKey] = strconv.FormatBool(result.Enabled.ValueBool())

	if result.Enabled.ValueBool() {
		result.Enabled = types.BoolValue(false)
	}

	return result, nil
}

const (
	V1DeferredUploadKey  = "v1_deferred_upload"
	V1OriginalEnabledKey = "v1_original_enabled"
)

func isV1DeferredUpload(data *process.ProcessData) bool {
	return data.StringArgs[V1DeferredUploadKey] == "true"
}

func getV1OriginalEnabled(data *process.ProcessData) (bool, error) {
	val, ok := data.StringArgs[V1OriginalEnabledKey]
	if !ok {
		return false, fmt.Errorf("missing original enabled value for V1 deferred upload")
	}
	return strconv.ParseBool(val)
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

func (p *FilePreProcessor) handleCreateOrUpdate(requestContext *common.RequestContext, data *process.ProcessData, planGetter process.Getter) (*filetf.FileTFModel, error) {

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
