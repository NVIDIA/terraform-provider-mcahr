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

package helper

import (
	"terraform/terraform-provider/provider/common"
	externalapi "terraform/terraform-provider/provider/external_api"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/external_api/resources"
	"terraform/terraform-provider/provider/external_api/resources/statement"
)

func RunOpCommand[API resources.APIModel](requestContext *common.RequestContext, client *client.PlatformClient, apiVersion common.APIVersion, command string) (API, error) {

	apiRequest := statement.StatementInputAPIModel{
		Statement:  command,
		APIVersion: apiVersion,
	}

	apiResponse, err := externalapi.CallExternalAPI[API](requestContext, client, &apiRequest)
	if err != nil {
		var nilAPI API // return the zero value of API (which is nil)
		return nilAPI, err
	}

	return apiResponse, nil
}
