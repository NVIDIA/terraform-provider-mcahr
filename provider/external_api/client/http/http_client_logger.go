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
	"net/http"
	"strings"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/log"
	"time"
)

// securityHeaders are headers that should be excluded from debug logs for security reasons
var securityHeaders = []string{
	"authorization",
	"x-api-key",
	"api-key",
	"apikey",
	"x-auth-token",
	"auth-token",
	"x-access-token",
	"access-token",
	"bearer",
	"cookie",
	"set-cookie",
	"x-csrf-token",
	"csrf-token",
	"x-session-token",
	"session-token",
}

// filterSecurityHeaders returns a filtered map of headers excluding security-sensitive headers
func filterSecurityHeaders(headers http.Header) map[string]string {
	filtered := make(map[string]string)
	for key, values := range headers {
		keyLower := strings.ToLower(key)

		// Check if this header is in the security headers list
		if isSecurityHeader(keyLower) {
			// Mask security headers
			filtered[key] = "***MASKED***"
		} else {
			// Join multiple values with comma if present
			filtered[key] = strings.Join(values, ", ")
		}
	}
	return filtered
}

// isSecurityHeader checks if a header name matches any security header pattern
func isSecurityHeader(keyLower string) bool {
	for _, secHeader := range securityHeaders {
		if keyLower == secHeader || strings.Contains(keyLower, secHeader) {
			return true
		}
	}
	return false
}

// logRequest logs basic HTTP request information at info level
func (h *HTTPClient) logRequest(requestCtx *common.RequestContext, req *http.Request) {
	logFields := map[string]any{
		"method":   req.Method,
		"endpoint": req.URL.Path,
	}

	log.LogInfo(requestCtx, "HTTP Request", logFields)
}

// logRequestDetails logs detailed HTTP request information at debug level
func (h *HTTPClient) logRequestDetails(requestCtx *common.RequestContext, req *http.Request, body []byte) {
	logFields := map[string]any{
		"method":   req.Method,
		"url":      req.URL.String(),
		"endpoint": req.URL.Path,
		"headers":  filterSecurityHeaders(req.Header),
	}

	// Add query parameters if present
	if len(req.URL.Query()) > 0 {
		logFields["query_params"] = req.URL.Query()
	}

	// Add request body if present
	if len(body) > 0 {
		logFields["request_body"] = string(body)
		logFields["request_body_size"] = len(body)
	}

	log.LogDebug(requestCtx, "HTTP Request Details", logFields)
}

// logResponse logs basic HTTP response information at info/error level
func (h *HTTPClient) logResponse(requestCtx *common.RequestContext, req *http.Request, resp *http.Response, err error, duration time.Duration) {
	logFields := map[string]any{
		"method":      req.Method,
		"endpoint":    req.URL.Path,
		"duration_ms": duration.Milliseconds(),
	}

	if err != nil {
		logFields["error"] = err.Error()
		log.LogError(requestCtx, "HTTP Request Failed", logFields)
		return
	}

	logFields["status_code"] = resp.StatusCode
	logFields["status"] = resp.Status

	log.LogInfo(requestCtx, "HTTP Response", logFields)
}

// logResponseDetails logs detailed HTTP response information at debug level
func (h *HTTPClient) logResponseDetails(requestCtx *common.RequestContext, resp *http.Response, body []byte) {
	logFields := map[string]any{
		"status_code": resp.StatusCode,
		"status":      resp.Status,
		"headers":     filterSecurityHeaders(resp.Header),
	}

	// Add response body
	if len(body) > 0 {
		logFields["response_body"] = string(body)
		logFields["response_body_size"] = len(body)
	}

	// Add content type if present
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		logFields["content_type"] = contentType
	}

	log.LogDebug(requestCtx, "HTTP Response Details", logFields)
}
