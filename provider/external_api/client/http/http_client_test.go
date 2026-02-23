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
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"terraform/terraform-provider/provider/common"
	"testing"
	"time"
)

func TestHTTPClientExecuteSuccess(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization header, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("User-Agent") != "terraform" {
			t.Errorf("expected terraform User-Agent, got %s", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Accept") != "*/*" {
			t.Errorf("expected Accept header, got %s", r.Header.Get("Accept"))
		}
		if r.Header.Get("Custom-Header") != "custom-value" {
			t.Errorf("expected Custom-Header, got %s", r.Header.Get("Custom-Header"))
		}

		// Check request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		expectedBody := `{"test": "data"}`
		if string(body) != expectedBody {
			t.Errorf("expected body %s, got %s", expectedBody, string(body))
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())
	requestBody := strings.NewReader(`{"test": "data"}`)
	httpReq := &HTTPRequest{
		Method: "POST",
		URL:    server.URL + "/test",
		Body:   requestBody,
		Headers: map[string]string{
			"Authorization": "Bearer test-token",
			"Content-Type":  "application/json",
			"Custom-Header": "custom-value",
		},
	}

	// when
	resp, err := client.Execute(requestContext, httpReq)

	// then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Response == nil {
		t.Fatalf("expected non-nil Response")
	}
	if resp.Response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.Response.StatusCode)
	}
	if resp.Response.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type header, got %s", resp.Response.Header.Get("Content-Type"))
	}
	expectedBody := `{"success": true}`
	if string(resp.Body) != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, string(resp.Body))
	}
}

func TestHTTPClientExecuteError(t *testing.T) {
	t.Parallel()

	// given
	client := NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())

	// Use invalid URL to trigger error
	httpReq := &HTTPRequest{
		Method:  "GET",
		URL:     "http://invalid-url-that-does-not-exist.local",
		Body:    nil,
		Headers: map[string]string{},
	}

	// when
	_, err := client.Execute(requestContext, httpReq)

	// then
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
	if !strings.Contains(err.Error(), "HTTP request failed") {
		t.Errorf("expected error message to contain context, got: %v", err)
	}
}

func TestHTTPClientExecuteTimeout(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Sleep longer than client timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	requestContext := common.NewRequestContext(context.Background())
	client := NewHTTPClient(10 * time.Millisecond) // Very short timeout
	httpReq := &HTTPRequest{
		Method:  "GET",
		URL:     server.URL,
		Body:    nil,
		Headers: map[string]string{},
	}

	// when
	_, err := client.Execute(requestContext, httpReq)

	// then
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "HTTP request failed") {
		t.Errorf("expected error message to contain context, got: %v", err)
	}
}

func TestHTTPClientExecuteInvalidRequestCreation(t *testing.T) {
	t.Parallel()

	// given
	client := NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())

	// Use invalid method to trigger request creation error
	httpReq := &HTTPRequest{
		Method:  "INVALID METHOD WITH SPACES",
		URL:     "http://example.com",
		Body:    nil,
		Headers: map[string]string{},
	}

	// when
	_, err := client.Execute(requestContext, httpReq)

	// then
	if err == nil {
		t.Fatal("expected error for invalid method")
	}
	if !strings.Contains(err.Error(), "failed to create HTTP request") {
		t.Errorf("expected error message about request creation, got: %v", err)
	}
}

func TestHTTPClientExecuteBodyReadError(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	}))
	defer server.Close()

	client := NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())

	// Use a reader that will cause an error when trying to read the request body
	errorReader := &errorReader{err: io.ErrUnexpectedEOF}

	httpReq := &HTTPRequest{
		Method:  "POST",
		URL:     server.URL,
		Body:    errorReader,
		Headers: map[string]string{},
	}

	// when
	_, err := client.Execute(requestContext, httpReq)

	// then
	if err == nil {
		t.Fatal("expected error from error reader")
	}
}

// errorReader is a helper for testing error conditions
type errorReader struct {
	err error
}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, er.err
}

func TestHTTPClientResponseBodyReadError(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send headers first
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("partial"))

		// Hijack the connection and close it abruptly
		if hijacker, ok := w.(http.Hijacker); ok {
			conn, _, err := hijacker.Hijack()
			if err == nil {
				conn.Close() // Close connection during body transmission
			}
		}
	}))
	defer server.Close()

	requestContext := common.NewRequestContext(context.Background())
	client := NewHTTPClient(30 * time.Second)
	httpReq := &HTTPRequest{
		Method:  "GET",
		URL:     server.URL,
		Body:    nil,
		Headers: map[string]string{},
	}

	// when
	_, err := client.Execute(requestContext, httpReq)

	// then
	if err == nil {
		t.Fatal("expected error from connection close during body read")
	}

	// Should be an error about reading response body
	if !strings.Contains(err.Error(), "failed to read response body") &&
		!strings.Contains(err.Error(), "HTTP request failed") {
		t.Errorf("expected response body read error, got: %v", err)
	}
}
