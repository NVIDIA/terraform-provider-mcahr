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
	"os"
	"strings"
	"terraform/terraform-provider/provider/common"
	externalfile "terraform/terraform-provider/provider/external_api/file"
	"terraform/terraform-provider/provider/tf/core/process"
)

const (
	tempDownloadFilePathPrefix = "tmp_opcp-"
)

var (
	TempDownloadFilePath = ""
)

func shouldDownloadFile(inputFile string) bool {
	return strings.HasPrefix(inputFile, "http:") || strings.HasPrefix(inputFile, "https://")
}

func maybeDownloadFile(requestContext *common.RequestContext, data *process.ProcessData, fileUrl string) (fileName string, err error) {

	if shouldDownloadFile(fileUrl) {
		tmpFilePath, err := externalfile.DownloadFileHttpsToTemp(requestContext, data.Client.GetHttpClient(), fileUrl, tempDownloadFilePathPrefix)
		if err != nil {
			return "", fmt.Errorf("failed to read remote file object %s: %s", fileUrl, err)
		}

		// defer the removal of the temporary file at the end of the resource operation flow
		data.DeferFunctionList.AddDefer(func() {
			os.Remove(tmpFilePath)
		})

		// needed for the upload operation
		data.StringArgs["downloaded_file_path"] = tmpFilePath

		return tmpFilePath, nil
	}

	// if not remote, return the original local file path
	return fileUrl, nil
}
