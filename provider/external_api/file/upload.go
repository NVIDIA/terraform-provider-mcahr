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
	"fmt"
	"net/http"
	"os"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/log"
	httpclient "terraform/terraform-provider/provider/external_api/client/http"
)

func UploadFileHttps(requestContext *common.RequestContext, client *httpclient.HTTPClient, sourceFile string, destUrl string) error {
	// Open file
	file, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("couldn't open local upload file '%s'", sourceFile)
	}
	defer file.Close()

	// Get file stats (size)
	stat, err := os.Stat(sourceFile)
	if err != nil {
		return fmt.Errorf("couldn't stat file to upload: %s", err.Error())
	}
	fileSize := stat.Size()

	// Upload file
	response, err := client.ExecuteRaw(requestContext, &httpclient.HTTPRequest{
		Method: http.MethodPut,
		URL:    destUrl,
		Body:   file,
		Headers: map[string]string{
			"x-ms-blob-type": "BlockBlob", // only used by Azure, ignored by S3
		},
		ContentLength: fileSize,
	})
	if err != nil {
		return fmt.Errorf("couldn't upload file: %s", err.Error())
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != 201 && response.StatusCode != 200 {
		var body []byte
		response.Body.Read(body)
		return fmt.Errorf("couldn't upload file, status: %s, message: %v", response.Status, string(body))
	}

	// Log success
	log.LogInfo(requestContext, fmt.Sprintf("Uploaded file '%s' (%d bytes) status: %v - %v\n", sourceFile, fileSize, response.StatusCode, http.StatusText(response.StatusCode)), nil)

	return nil
}

func UploadFileHttpsFromString(requestContext *common.RequestContext, client *httpclient.HTTPClient, fileData string, destUrl string) error {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		return fmt.Errorf("couldn't create local upload file\n")
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Write data to the temporary file
	if _, err := tempFile.Write([]byte(fileData)); err != nil {
		return fmt.Errorf("couldn't write data to local upload file\n")
	}

	// Upload the temporary file
	return UploadFileHttps(requestContext, client, tempFile.Name(), destUrl)
}
