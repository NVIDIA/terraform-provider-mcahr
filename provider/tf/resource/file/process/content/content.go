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

package content

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ProcessFileContent(requestContext *common.RequestContext, data *process.ProcessData, planModel *filetf.FileTFModel) (*filetf.FileTFModel, error) {

	if common.IsAttrKnown(planModel.InlineData) {
		if err := processInlineData(planModel); err != nil {
			return nil, err
		}
	} else if common.IsAttrKnown(planModel.InputFile) {
		if err := processInputFile(requestContext, data, planModel); err != nil {
			return nil, err
		}
	}

	return planModel, nil
}

func processInlineData(planModel *filetf.FileTFModel) error {

	contentBytes := []byte(planModel.InlineData.ValueString())

	md5, err := ContentMd5(contentBytes)
	if err != nil {
		return err
	}

	contentSize := ContentSize(contentBytes)

	setModelFields(planModel, md5, contentSize)

	return nil
}

func processInputFile(requestContext *common.RequestContext, data *process.ProcessData, planModel *filetf.FileTFModel) error {

	filePath, err := maybeDownloadFile(requestContext, data, planModel.InputFile.ValueString())
	if err != nil {
		return fmt.Errorf("failed to download file: %s", err)
	}

	md5, err := FileMd5(filePath)
	if err != nil {
		return err
	}

	contentSize, err := FileSize(filePath)
	if err != nil {
		return err
	}

	setModelFields(planModel, md5, contentSize)

	return nil
}

func setModelFields(planModel *filetf.FileTFModel, md5 string, contentSize int64) {
	planModel.MD5 = types.StringValue(md5)
	planModel.Checksum = types.StringValue(md5)
	planModel.FileLength = types.Int64Value(contentSize)
}
