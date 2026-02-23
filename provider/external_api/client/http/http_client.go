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

package http_client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/log"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

type HTTPRequest struct {
	Method        string
	URL           string
	Body          io.Reader
	Headers       map[string]string
	ContentLength int64
}

type HTTPResponse struct {
	Response *http.Response
	Body     []byte
}

var defaultHeaders = map[string]string{
	"User-Agent": "terraform",
	"Accept":     "*/*",
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{Timeout: timeout},
	}
}

func (h *HTTPClient) createRequest(requestContext *common.RequestContext, httpReq *HTTPRequest) (*http.Request, error) {
	req, err := http.NewRequestWithContext(requestContext.Context, httpReq.Method, httpReq.URL, httpReq.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	return req, nil
}

func (h *HTTPClient) readRequestBody(requestContext *common.RequestContext, httpReq *HTTPRequest) ([]byte, error) {

	if httpReq.Body != nil {
		body, err := io.ReadAll(httpReq.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		// Create new reader with the body content for the actual request
		httpReq.Body = bytes.NewReader(body)
	}
	return nil, nil
}

func (h *HTTPClient) Execute(requestContext *common.RequestContext, httpReq *HTTPRequest) (*HTTPResponse, error) {
	// It's recommended to use this method instead of ExecuteRaw
	// if you can afford to read the entire body in memory (both request and response)
	// It also closes the response body automatically

	debugEnabled := log.IsDebugEnabled(requestContext)
	var loggingRequestBody []byte
	if debugEnabled {
		// Only read and preserve request body if debug logging is enabled
		var err error
		loggingRequestBody, err = h.readRequestBody(requestContext, httpReq)
		if err != nil {
			return nil, err
		}
	}

	req, err := h.createRequest(requestContext, httpReq)
	if err != nil {
		return nil, err
	}

	if debugEnabled {
		h.logRequestDetails(requestContext, req, loggingRequestBody)
	}

	resp, err := h.doExecuteRequest(requestContext, req, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if debugEnabled {
		h.logResponseDetails(requestContext, resp, responseBody)
	}

	return &HTTPResponse{
		Response: resp,
		Body:     responseBody,
	}, nil
}

func (h *HTTPClient) ExecuteRaw(requestContext *common.RequestContext, httpReq *HTTPRequest) (*http.Response, error) {
	// WARNING: Use this method with caution, it doesn't close the response body

	req, err := h.createRequest(requestContext, httpReq)
	if err != nil {
		return nil, err
	}

	return h.doExecuteRequest(requestContext, req, httpReq)
}

func (h *HTTPClient) doExecuteRequest(requestContext *common.RequestContext, req *http.Request, httpReq *HTTPRequest) (*http.Response, error) {

	h.setHeaders(req, httpReq.Headers)
	h.setContentLength(req, httpReq.ContentLength)

	h.logRequest(requestContext, req)

	start := time.Now()
	resp, err := h.client.Do(req)
	duration := time.Since(start)

	h.logResponse(requestContext, req, resp, err, duration)

	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// WARNING: resp.Body is not closed here, so the caller is responsible for closing it
	return resp, nil

}

func (h *HTTPClient) setHeaders(req *http.Request, headers map[string]string) {
	// Add extra headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set default headers if not already set
	for key, value := range defaultHeaders {
		if req.Header.Get(key) == "" {
			req.Header.Set(key, value)
		}
	}
}

func (h *HTTPClient) setContentLength(req *http.Request, contentLength int64) {
	if contentLength > 0 {
		req.ContentLength = contentLength
	}
}
