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

package upload

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"

	externalfile "terraform/terraform-provider/provider/external_api/file"
)

func HandleFileUploadProcess(requestContext *common.RequestContext, data *process.ProcessData, planModel *filetf.FileTFModel) error {

	// Get presigned URL
	presignedUrl, err := GetFilePresignedPut(requestContext, data.Client, requestContext.APIVersion, planModel.Name.ValueString())
	if err != nil {
		return fmt.Errorf("failed to get presigned URL: %s", err)
	}

	// Upload file
	err = uploadFile(requestContext, data, planModel, presignedUrl)
	if err != nil {
		return fmt.Errorf("failed to upload file: %s", err)
	}

	return nil
}

func uploadFile(requestContext *common.RequestContext, data *process.ProcessData, planModel *filetf.FileTFModel, presignedPutUrl string) error {

	if common.IsAttrKnown(planModel.InlineData) {
		// Upload the string data from inline_data
		return externalfile.UploadFileHttpsFromString(
			requestContext,
			data.Client.GetHttpClient(),
			planModel.InlineData.ValueString(),
			presignedPutUrl,
		)
	}

	inputFile := getInputFile(data, planModel)

	// Upload the file from input_file
	return externalfile.UploadFileHttps(
		requestContext,
		data.Client.GetHttpClient(),
		inputFile,
		presignedPutUrl,
	)
}

func getInputFile(data *process.ProcessData, planModel *filetf.FileTFModel) string {
	if data.StringArgs["downloaded_file_path"] != "" {
		return data.StringArgs["downloaded_file_path"]
	}
	return planModel.InputFile.ValueString()
}
