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

package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"terraform/terraform-provider/provider/common"
	commonlog "terraform/terraform-provider/provider/common/log"
	http_client "terraform/terraform-provider/provider/external_api/client/http"
)

type PlatformClient struct {
	baseURL    string
	apiToken   string
	httpClient *http_client.HTTPClient
}

type PlatformClientRequest struct {
	Method   string
	Endpoint string
	Body     io.Reader
}

type PlatformClientResponse struct {
	Response *http.Response
	Body     []byte
}

var (
	// TODO: Maybe make these configurable
	clientTimeoutSeconds  = 90
	maxRetries            = 5
	rateLimitDelaySeconds = 15
)

func (c *PlatformClient) GetHttpClient() *http_client.HTTPClient {
	return c.httpClient
}
func NewPlatformClient(baseURL, apiToken string) *PlatformClient {
	httpClient := http_client.NewHTTPClient(time.Second * time.Duration(clientTimeoutSeconds))

	return &PlatformClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		apiToken:   apiToken,
		httpClient: httpClient,
	}
}

// ExecuteRequest executes a request with the provided RequestContext wrapper for enhanced logging and metadata
func (c *PlatformClient) ExecuteRequest(requestCtx *common.RequestContext, request *PlatformClientRequest) (resp *PlatformClientResponse, err error) {
	// Use the subsystem logging with metadata from RequestContext
	logFields := map[string]interface{}{
		"method":   request.Method,
		"endpoint": request.Endpoint,
	}

	// Use the logging manager for consistent logging
	commonlog.LogDebug(requestCtx, "Executing platform client request", logFields)

	resp, err = c.executeWithRetries(requestCtx, request)
	if err != nil {
		logFields["error"] = err.Error()
		commonlog.LogError(requestCtx, "Platform client request failed", logFields)
	}

	return resp, err
}

func (c *PlatformClient) executeWithRetries(requestCtx *common.RequestContext, request *PlatformClientRequest) (*PlatformClientResponse, error) {

	// Buffer the body up front. request.Body is a single-use io.Reader that
	// the first attempt drains, so every retry after it used to send an empty
	// body -- a write silently degrading into a no-op rather than failing.
	var body []byte
	if request.Body != nil {
		var err error
		body, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
	}

	for attempt := 0; ; attempt++ {
		resp, err := c.executeRetriableRequest(requestCtx, request, body)
		if err == nil {
			return resp, nil
		}

		retry, err := handleRetries(attempt, err)
		if !retry || err != nil {
			return resp, err
		}
	}
}

func (c *PlatformClient) executeRetriableRequest(requestCtx *common.RequestContext, request *PlatformClientRequest, body []byte) (platformResp *PlatformClientResponse, err error) {

	// Every attempt gets its own reader over the buffered bytes. Requests that
	// started with no body keep having none, so Content-Length is unaffected.
	var bodyReader io.Reader
	if request.Body != nil {
		bodyReader = bytes.NewReader(body)
	}

	resp, err := c.httpClient.Execute(requestCtx, &http_client.HTTPRequest{
		Method: request.Method,
		URL:    c.baseURL + request.Endpoint,
		Body:   bodyReader,
		Headers: map[string]string{
			"Authorization": "Bearer " + c.apiToken,
			"Content-Type":  "application/json; charset=utf-8",
		},
	})
	if err != nil {
		return nil, err
	}

	if err := maybeMapStatusCodeError(resp.Response); err != nil {
		return nil, err
	}

	return &PlatformClientResponse{
		Response: resp.Response,
		Body:     resp.Body,
	}, nil
}
