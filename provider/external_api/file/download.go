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

package file

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"

	"terraform/terraform-provider/provider/common"
	httpclient "terraform/terraform-provider/provider/external_api/client/http"
)

func DownloadFileHttpsToTemp(requestContext *common.RequestContext, client *httpclient.HTTPClient, sourceUrl string, destPattern string) (fileName string, err error) {

	resp, err := client.ExecuteRaw(requestContext, &httpclient.HTTPRequest{
		Method: "GET",
		URL:    sourceUrl,
		Body:   nil,
	})
	if err != nil {
		return "", fmt.Errorf("couldn't open download url '%s'\n", sourceUrl)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file from url. Status: %s", resp.Status)
	}

	tempFile, err := os.CreateTemp("", destPattern)
	if err != nil {
		return "", fmt.Errorf("couldn't create local download file '%s'\n", destPattern)
	}
	defer tempFile.Close()

	// Use buffered writer for better performance
	bufferedWriter := bufio.NewWriter(tempFile)
	defer bufferedWriter.Flush()
	// Use buffered reader for the response
	bufferedReader := bufio.NewReader(resp.Body)

	// NOTE: This processes a block at a time, which is important for large files and mem usage
	_, err = io.Copy(bufferedWriter, bufferedReader)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("couldn't process download data from url '%s'\n", sourceUrl)
	}

	return tempFile.Name(), nil
}
