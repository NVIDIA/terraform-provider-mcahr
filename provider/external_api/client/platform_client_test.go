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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"terraform/terraform-provider/provider/common"
	"testing"
	"time"
)

func TestPlatformClientExecuteRequestSuccess(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			t.Errorf("expected JSON content type, got %s", r.Header.Get("Content-Type"))
		}
		expectedAuth := "Bearer test-api-token"
		if r.Header.Get("Authorization") != expectedAuth {
			t.Errorf("expected Authorization %s, got %s", expectedAuth, r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	apiToken := "test-api-token"
	platformClient := NewPlatformClient(server.URL, apiToken)
	request := &PlatformClientRequest{
		Method:   "POST",
		Endpoint: "/test",
		Body:     strings.NewReader(`{"test": "data"}`),
	}

	// when
	response, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)

	// then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.Response.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", response.Response.StatusCode)
	}
	expectedBody := `{"result": "success"}`
	if string(response.Body) != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, string(response.Body))
	}
}

func TestPlatformClientExecuteRequestAuthErrorNoRetry(t *testing.T) {
	t.Parallel()

	// given
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}))
	defer server.Close()

	apiToken := "test-api-token"
	platformClient := NewPlatformClient(server.URL, apiToken)
	request := &PlatformClientRequest{
		Method:   "GET",
		Endpoint: "/test",
		Body:     nil,
	}

	// when
	_, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)

	// then
	if err == nil {
		t.Fatal("expected error for 401 status")
	}
	// Should only make one request (no retries for auth errors)
	if requestCount != 1 {
		t.Errorf("expected 1 request (no retries), got %d", requestCount)
	}
	// Error should be an AuthenticationError
	if _, ok := err.(*AuthenticationError); !ok {
		t.Errorf("expected AuthenticationError, got %T", err)
	}
}

func TestPlatformClientExecuteRequestServerErrorRetry(t *testing.T) {
	t.Parallel()

	// given
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount <= 2 {
			// First two requests - return 500
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		} else {
			// Third request - return success
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"result": "success after retry"}`))
		}
	}))
	defer server.Close()

	apiToken := "test-api-token"
	platformClient := NewPlatformClient(server.URL, apiToken)
	request := &PlatformClientRequest{
		Method:   "GET",
		Endpoint: "/test",
		Body:     nil,
	}

	// when
	response, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)

	// then
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}
	if requestCount != 3 {
		t.Errorf("expected 3 requests (2 failures + 1 success), got %d", requestCount)
	}
	if response.Response.StatusCode != 200 {
		t.Errorf("expected final status 200, got %d", response.Response.StatusCode)
	}
	expectedBody := `{"result": "success after retry"}`
	if string(response.Body) != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, string(response.Body))
	}
}

func TestPlatformClientExecuteRequestClientErrorNoRetry(t *testing.T) {
	t.Parallel()

	// given
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		// Always return 400 (client error - should not retry)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	apiToken := "test-api-token"
	platformClient := NewPlatformClient(server.URL, apiToken)
	request := &PlatformClientRequest{
		Method:   "GET",
		Endpoint: "/test",
		Body:     nil,
	}

	// when
	_, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)

	// then
	if err == nil {
		t.Fatal("expected error for 400 status")
	}

	// Should only make one request (no retries for client errors)
	if requestCount != 1 {
		t.Errorf("expected 1 request (no retries), got %d", requestCount)
	}

	// Error should be a ClientError
	if _, ok := err.(*ClientError); !ok {
		t.Errorf("expected ClientError, got %T", err)
	}
}

func TestPlatformClientMaybeMapStatusCode(t *testing.T) {
	t.Parallel()

	// given
	testCases := []struct {
		statusCode    int
		expectedError string
	}{
		{200, ""},
		{201, ""},
		{299, ""},
		{301, "*client.HTTPError"}, // Redirect - default case
		{302, "*client.HTTPError"}, // Redirect - default case
		{304, "*client.HTTPError"}, // Not Modified - default case
		{400, "*client.ClientError"},
		{401, "*client.AuthenticationError"},
		{403, "*client.AuthorizationError"},
		{404, "*client.ClientError"},
		{429, "*client.RateLimitError"},
		{500, "*client.ServerError"},
		{502, "*client.ServerError"},
		{503, "*client.ServerError"},
		{600, "*client.HTTPError"}, // Custom status - default case
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("status_%d", tc.statusCode), func(t *testing.T) {
			t.Parallel()

			// given
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(fmt.Sprintf("Status %d", tc.statusCode)))
			}))
			defer server.Close()

			apiToken := "test-api-token"
			platformClient := NewPlatformClient(server.URL, apiToken)
			request := &PlatformClientRequest{
				Method:   "GET",
				Endpoint: "/test",
				Body:     nil,
			}

			// when
			_, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)

			// then
			if tc.expectedError == "" {
				if err != nil {
					t.Errorf("expected no error for status %d, got %v", tc.statusCode, err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error for status %d", tc.statusCode)
				} else {
					errorType := fmt.Sprintf("%T", err)
					if errorType != tc.expectedError {
						t.Errorf("expected error type %s for status %d, got %s", tc.expectedError, tc.statusCode, errorType)
					}
				}
			}
		})
	}
}

func TestPlatformClientExecuteRequestRateLimitRetry(t *testing.T) {
	t.Parallel()

	// given
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount == 1 {
			// First request - return 429 (rate limit)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limited"))
		} else {
			// Second request - return success
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"result": "success after rate limit"}`))
		}
	}))
	defer server.Close()

	apiToken := "test-api-token"
	platformClient := NewPlatformClient(server.URL, apiToken)
	request := &PlatformClientRequest{
		Method:   "GET",
		Endpoint: "/test",
		Body:     nil,
	}

	// when
	start := time.Now()
	response, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)
	elapsed := time.Since(start)

	// then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if requestCount != 2 {
		t.Errorf("expected 2 requests (1 rate limited + 1 success), got %d", requestCount)
	}

	if response.Response.StatusCode != 200 {
		t.Errorf("expected final status 200, got %d", response.Response.StatusCode)
	}

	// Should have waited for rate limit delay (15 seconds)
	expectedMinDuration := 14 * time.Second // Allow some tolerance
	if elapsed < expectedMinDuration {
		t.Errorf("expected to wait at least %v for rate limit, but only waited %v", expectedMinDuration, elapsed)
	}
}

func TestPlatformClientExecuteRequestMaxRetries(t *testing.T) {
	t.Parallel()

	// given
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		// Always return 500 to trigger retries
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	apiToken := "test-api-token"
	platformClient := NewPlatformClient(server.URL, apiToken)
	request := &PlatformClientRequest{
		Method:   "GET",
		Endpoint: "/test",
		Body:     nil,
	}

	// when
	_, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)

	// then
	if err == nil {
		t.Fatal("expected error after max retries")
	}

	// Should make maxRetries + 1 requests (initial + retries)
	// maxRetries is 5, so should make 6 total requests
	expectedRequests := 6
	if requestCount != expectedRequests {
		t.Errorf("expected %d requests (initial + max retries), got %d", expectedRequests, requestCount)
	}

	// Final error should be ServerError
	if _, ok := err.(*ServerError); !ok {
		t.Errorf("expected ServerError after max retries, got %T", err)
	}
}

func TestPlatformClientBaseURLTrimming(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the URL path is correct (no double slashes)
		expectedPath := "/test"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	apiToken := "test-api-token"
	// Test with trailing slash in baseURL
	platformClient := NewPlatformClient(server.URL+"/", apiToken)
	request := &PlatformClientRequest{
		Method:   "GET",
		Endpoint: "/test",
		Body:     nil,
	}

	// when
	_, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)

	// then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPlatformClientExecuteRequestHTTPClientError(t *testing.T) {
	t.Parallel()

	// given
	apiToken := "test-api-token"
	// Use invalid URL to trigger HTTP client error
	platformClient := NewPlatformClient("http://invalid-host-that-does-not-exist.local", apiToken)
	request := &PlatformClientRequest{
		Method:   "GET",
		Endpoint: "/test",
		Body:     nil,
	}

	// when
	_, err := platformClient.ExecuteRequest(common.NewRequestContext(context.Background()), request)

	// then
	if err == nil {
		t.Fatal("expected error for invalid host")
	}

	if !strings.Contains(err.Error(), "HTTP request failed") {
		t.Errorf("expected HTTP request failed error, got: %v", err)
	}
}
