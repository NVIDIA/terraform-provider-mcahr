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
	"net/http"
	"strings"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/external_api/resources"
	"terraform/terraform-provider/provider/external_api/resources/statement"
)

var v1Endpoint = "/api/v1/execute"
var v2Endpoint = "/api/v1/statements/execute"

// PlatformClientInterface defines the interface for executing HTTP requests
type PlatformClientInterface interface {
	ExecuteRequest(requestContext *common.RequestContext, request *client.PlatformClientRequest) (*client.PlatformClientResponse, error)
}

func CallExternalAPI[API resources.APIModel](requestContext *common.RequestContext, client *client.PlatformClient, apiObject *statement.InputAPIModel) (API, error) {
	return CallExternalAPIWithClient[API](requestContext, client, apiObject)
}

// CallExternalAPIWithClient calls the external API with a custom client implementation (useful for testing)
func CallExternalAPIWithClient[API resources.APIModel](requestContext *common.RequestContext, clientInterface PlatformClientInterface, apiObject *statement.InputAPIModel) (API, error) {
	var nilAPI API
	var request *client.PlatformClientRequest
	var err error

	isResourceCreatedViaOpStatements := common.IsResourceCreatedViaOpStatements(requestContext.ResourceType)
	if isResourceCreatedViaOpStatements {
		request, err = createRequestForStatementExecute(apiObject)
	} else {
		request, err = createRequestForCustomAPI(requestContext, apiObject)
	}

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

	// Skip error checking only for custom REST resource DELETE operations (they may return 204 NoContent)
	skipErrorCheck := !isResourceCreatedViaOpStatements && requestContext.Operation == common.Delete
	if !common.IsNil(apiResponse) && !skipErrorCheck {
		apiBusinessErrors := apiResponse.GetErrors()
		if apiBusinessErrors != "" {
			return nilAPI, fmt.Errorf("API response errors: %s", apiBusinessErrors)
		}
	}

	return apiResponse, nil
}

func createRequestForStatementExecute(apiObject *statement.InputAPIModel) (*client.PlatformClientRequest, error) {

	body, err := json.Marshal(apiObject)
	if err != nil {
		return nil, err
	}

	// Select endpoint based on backend version
	method, endpoint, err := getMethodAndEndpoint(apiObject.APIVersion)
	if err != nil {
		return nil, err
	}

	return &client.PlatformClientRequest{
		Method:   method,
		Endpoint: endpoint,
		Body:     bytes.NewReader(body),
	}, nil
}

func createRequestForCustomAPI(requestContext *common.RequestContext, apiObject *statement.InputAPIModel) (*client.PlatformClientRequest, error) {
	apiPayload := apiObject.ApiPayload
	method, endpoint, err := getMethodAndEndpointForCustomAPI(requestContext, apiPayload)
	if err != nil {
		return nil, err
	}

	request := &client.PlatformClientRequest{
		Method:   method,
		Endpoint: endpoint,
	}

	if requestContext.Operation != common.Delete {
		request.Body = strings.NewReader(apiPayload)
	}

	return request, nil
}

// getEndpoint returns the appropriate API endpoint based on the backend version
func getMethodAndEndpoint(version common.APIVersion) (string, string, error) {
	switch version {
	case common.V1:
		return "POST", v1Endpoint, nil
	case common.V2:
		return "POST", v2Endpoint, nil
	default:
		return "", "", fmt.Errorf("unknown API version: %v", version)
	}
}

func getMethodAndEndpointForCustomAPI(requestContext *common.RequestContext, apiPayload string) (string, string, error) {
	var method, endpoint string
	var err error

	switch requestContext.ResourceType {
	case "smtp_subscription":
		endpoint, err = getSmtpSubscriptionEndpoint(requestContext, apiPayload)
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", fmt.Errorf("unknown resource type: %v", requestContext.ResourceType)
	}

	switch requestContext.Operation {
	case common.Create:
		method = "POST"
	case common.Read:
		method = "GET"
	case common.Update:
		method = "PATCH"
	case common.Delete:
		method = "DELETE"
	default:
		return "", "", fmt.Errorf("unknown operation: %v", requestContext.Operation)
	}

	return method, endpoint, nil
}

func getSmtpSubscriptionEndpoint(requestContext *common.RequestContext, apiPayload string) (string, error) {
	switch requestContext.Operation {
	case common.Create:
		return "/api/v1/integrations/smtp/subscriptions", nil
	case common.Read, common.Update, common.Delete:
		unmarshaledPayload := make(map[string]interface{})
		err := json.Unmarshal([]byte(apiPayload), &unmarshaledPayload)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal apiPayload: %w", err)
		}
		// if id is not present, return an error
		if _, ok := unmarshaledPayload["id"]; !ok {
			return "", fmt.Errorf("id is not present in the apiPayload")
		}
		subscriptionID := unmarshaledPayload["id"].(string)
		return fmt.Sprintf("/api/v1/integrations/smtp/subscriptions/%s", subscriptionID), nil
	default:
		return "", fmt.Errorf("unknown operation: %v for resource type: %v", requestContext.Operation, requestContext.ResourceType)
	}
}

func processResponse[API resources.APIModel](resp *client.PlatformClientResponse) (API, error) {
	var apiResponse API

	// If response code is NoContent (204), there is no response body to unmarshal
	// This typically happens for DELETE operations
	if resp.Response.StatusCode == http.StatusNoContent || len(resp.Body) == 0 {
		return apiResponse, nil
	}

	err := json.Unmarshal(resp.Body, &apiResponse)
	if err != nil {
		var nilAPI API // return the zero value of API (which is nil)
		return nilAPI, err
	}

	return apiResponse, nil
}
