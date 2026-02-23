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
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/external_api/resources"
	filesapi "terraform/terraform-provider/provider/external_api/resources/files"
	corehelper "terraform/terraform-provider/provider/tf/core/helper"
)

type GetFilePresignedPutAPIModelV1 struct {
	PresingedPut string `json:"get_file_attribute"`
}

var _ resources.APIModel = &GetFilePresignedPutAPIModelV1{}

func (t GetFilePresignedPutAPIModelV1) GetErrors() string {
	return ""
}

func GetFilePresignedPut(requestContext *common.RequestContext, client *client.PlatformClient, apiVersion common.APIVersion, fileObjectName string) (string, error) {

	statement := fmt.Sprintf("get_file_attribute(name=\"%s\", field_name=\"presigned_put\")", fileObjectName)

	presignedUrl, err := getFilePresignedPutByApiVersion(requestContext, client, apiVersion, statement)
	if err != nil {
		return "", err
	}

	if presignedUrl == "" {
		return "", fmt.Errorf("presigned URL is empty")
	}

	return presignedUrl, nil
}

func getFilePresignedPutByApiVersion(requestContext *common.RequestContext, client *client.PlatformClient, apiVersion common.APIVersion, statement string) (string, error) {
	switch apiVersion {
	case common.V1:
		return getFilePresignedPutV1(requestContext, client, statement)
	case common.V2:
		return getFilePresignedPutV2(requestContext, client, statement)
	default:
		return "", fmt.Errorf("unknown API version: %v", apiVersion)
	}
}

func getFilePresignedPutV1(requestContext *common.RequestContext, client *client.PlatformClient, statement string) (string, error) {
	presignedUrlResponse, err := corehelper.RunOpCommand[*GetFilePresignedPutAPIModelV1](requestContext, client, common.V1, statement)
	if err != nil {
		return "", presignedPutErrorResponse(err)
	}

	return presignedUrlResponse.PresingedPut, nil
}

func getFilePresignedPutV2(requestContext *common.RequestContext, client *client.PlatformClient, statement string) (string, error) {
	presignedUrlResponse, err := corehelper.RunOpCommand[*filesapi.FileResponseAPIModel](requestContext, client, common.V2, statement)
	if err != nil {
		return "", presignedPutErrorResponse(err)
	}

	// Validate response structure step by step
	if presignedUrlResponse == nil {
		return "", fmt.Errorf("presigned URL response is nil")
	}

	if presignedUrlResponse.Output.Configurations.Items == nil {
		return "", fmt.Errorf("file configurations items is nil")
	}

	if len(presignedUrlResponse.Output.Configurations.Items) == 0 {
		return "", fmt.Errorf("no file configurations returned from API")
	}

	config := presignedUrlResponse.Output.Configurations.Items[0].Config
	if config.PresignedUri == "" {
		return "", fmt.Errorf("presigned URI is empty in API response")
	}

	return config.PresignedUri, nil
}

func presignedPutErrorResponse(err error) error {
	return fmt.Errorf("failed to get presigned URL for file upload: %s", err)
}
