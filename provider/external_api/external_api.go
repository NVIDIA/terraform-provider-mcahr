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

package externalapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/external_api/resources"
	"terraform/terraform-provider/provider/external_api/resources/statement"
)

var v1Endpoint = "/api/v1/execute"
var v2Endpoint = "/api/v1/statements/execute"
var METHOD = "POST"

// PlatformClientInterface defines the interface for executing HTTP requests
type PlatformClientInterface interface {
	ExecuteRequest(requestContext *common.RequestContext, request *client.PlatformClientRequest) (*client.PlatformClientResponse, error)
}

func CallExternalAPI[API resources.APIModel](requestContext *common.RequestContext, client *client.PlatformClient, apiObject *statement.StatementInputAPIModel) (API, error) {
	return CallExternalAPIWithClient[API](requestContext, client, apiObject)
}

// CallExternalAPIWithClient calls the external API with a custom client implementation (useful for testing)
func CallExternalAPIWithClient[API resources.APIModel](requestContext *common.RequestContext, clientInterface PlatformClientInterface, apiObject *statement.StatementInputAPIModel) (API, error) {

	var nilAPI API

	request, err := createRequest(apiObject)
	if err != nil {
		return nilAPI, err
	}

	resp, err := clientInterface.ExecuteRequest(requestContext, request)
	if err != nil {
		return nilAPI, err
	}

	apiResponse, err := processResponse[API](resp)
	if err != nil {
		return nilAPI, err
	}

	if !common.IsNil(apiResponse) {
		apiBusinessErrors := apiResponse.GetErrors()
		if apiBusinessErrors != "" {
			return nilAPI, fmt.Errorf("API response errors: %s", apiBusinessErrors)
		}
	}

	return apiResponse, nil
}

func createRequest(apiObject *statement.StatementInputAPIModel) (*client.PlatformClientRequest, error) {

	body, err := json.Marshal(apiObject)
	if err != nil {
		return nil, err
	}

	// Select endpoint based on backend version
	endpoint, err := getEndpoint(apiObject.APIVersion)
	if err != nil {
		return nil, err
	}

	return &client.PlatformClientRequest{
		Method:   METHOD,
		Endpoint: endpoint,
		Body:     bytes.NewReader(body),
	}, nil
}

// getEndpoint returns the appropriate API endpoint based on the backend version
func getEndpoint(version common.APIVersion) (string, error) {
	switch version {
	case common.V1:
		return v1Endpoint, nil
	case common.V2:
		return v2Endpoint, nil
	default:
		return "", fmt.Errorf("unknown API version: %v", version)
	}
}

func processResponse[API resources.APIModel](resp *client.PlatformClientResponse) (API, error) {
	var apiResponse API
	err := json.Unmarshal(resp.Body, &apiResponse)
	if err != nil {
		var nilAPI API // return the zero value of API (which is nil)
		return nilAPI, err
	}

	return apiResponse, nil
}
