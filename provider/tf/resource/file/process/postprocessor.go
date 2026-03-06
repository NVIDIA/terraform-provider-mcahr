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
	"terraform/terraform-provider/provider/common/log"
	filesapi "terraform/terraform-provider/provider/external_api/resources/files"
	corehelper "terraform/terraform-provider/provider/tf/core/helper"
	"terraform/terraform-provider/provider/tf/core/process"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"
	"terraform/terraform-provider/provider/tf/resource/file/process/upload"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FilePostProcessor struct{}

var _ process.PostProcessor[*filetf.FileTFModel] = &FilePostProcessor{}

func (p *FilePostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, result *filetf.FileTFModel) error {
	if isV1DeferredUpload(data) {
		if err := handleV1DeferredUpload(requestContext, data, result); err != nil {
			return fmt.Errorf("failed to upload file: %s", err)
		}
	}
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

// handleV1DeferredUpload performs the file upload that was deferred during preprocessing
// for V1 backends. At this point, define_file has been called and the file object exists,
// so we can obtain the presigned URL, upload the content, then update file_data and enabled.
func handleV1DeferredUpload(requestContext *common.RequestContext, data *process.ProcessData, result *filetf.FileTFModel) error {
	name := result.Name.ValueString()

	log.LogInfo(requestContext, "Performing V1 deferred file upload", map[string]any{
		"file_name": name,
	})

	// Re-extract plan to get inline_data/input_file for the upload
	var planModel filetf.FileTFModel
	diags := data.CreateRequest.Plan.Get(requestContext.Context, &planModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get plan model for V1 upload: %s", diags.Errors())
	}

	// File now exists on backend; get presigned URL and upload
	err := upload.HandleFileUploadProcess(requestContext, data, &planModel)
	if err != nil {
		return fmt.Errorf("failed to upload file content: %s", err)
	}

	// Get the storage URI so we can set file_data
	uri, err := getFileAttribute(requestContext, data, name, "uri")
	if err != nil {
		return fmt.Errorf("failed to get file URI: %s", err)
	}

	originalEnabled, err := getV1OriginalEnabled(data)
	if err != nil {
		return err
	}

	// Update file_data and restore enabled to original value
	fileDataValue := fmt.Sprintf(":%s", uri)
	updateStmt := fmt.Sprintf(
		"update_file(name=\"%s\", file_data=\"%s\", enabled=%v)",
		name, fileDataValue, originalEnabled,
	)
	_, err = corehelper.RunOpCommand[*filesapi.FileResponseAPIModelV1](
		requestContext, data.Client, common.V1, updateStmt,
	)
	if err != nil {
		return fmt.Errorf("failed to update file_data after upload: %s", err)
	}

	// Update the result model to reflect the final state
	result.FileData = types.StringValue(fileDataValue)
	result.Enabled = types.BoolValue(originalEnabled)

	log.LogInfo(requestContext, "V1 deferred file upload completed", map[string]any{
		"file_name": name,
	})

	return nil
}

// getFileAttribute retrieves a file attribute from the V1 backend using get_file_attribute.
func getFileAttribute(requestContext *common.RequestContext, data *process.ProcessData, fileName, fieldName string) (string, error) {
	statement := fmt.Sprintf("get_file_attribute(name=\"%s\", field_name=\"%s\")", fileName, fieldName)
	response, err := corehelper.RunOpCommand[*upload.GetFilePresignedPutAPIModelV1](
		requestContext, data.Client, common.V1, statement,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get file attribute '%s': %s", fieldName, err)
	}
	value := response.PresingedPut
	if value == "" {
		return "", fmt.Errorf("file attribute '%s' is empty", fieldName)
	}
	return value, nil
}

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
