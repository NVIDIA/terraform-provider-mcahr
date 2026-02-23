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
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"terraform/terraform-provider/provider/common"
	httpclient "terraform/terraform-provider/provider/external_api/client/http"
)

func TestDownloadFileHttpsToTemp_Success(t *testing.T) {
	t.Parallel()

	// given
	expectedContent := "test file content for download"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedContent))
	}))
	defer server.Close()

	requestContext := common.NewRequestContext(context.Background())
	client := httpclient.NewHTTPClient(30 * time.Second)
	destPattern := "test-download-*.txt"

	// when
	fileName, err := DownloadFileHttpsToTemp(requestContext, client, server.URL, destPattern)

	// then
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if fileName == "" {
		t.Error("expected non-empty file name")
	}

	// Verify file exists and has correct content
	defer os.Remove(fileName)
	content, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("failed to read downloaded file: %v", err)
	}

	if string(content) != expectedContent {
		t.Errorf("expected content %s, got %s", expectedContent, string(content))
	}

	// Verify filename matches pattern
	if !strings.Contains(fileName, "test-download-") || !strings.HasSuffix(fileName, ".txt") {
		t.Errorf("filename %s doesn't match expected pattern %s", fileName, destPattern)
	}
}

func TestDownloadFileHttpsToTemp_HTTPError(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
	}))
	defer server.Close()

	requestContext := common.NewRequestContext(context.Background())
	client := httpclient.NewHTTPClient(30 * time.Second)

	// when
	fileName, err := DownloadFileHttpsToTemp(requestContext, client, server.URL, "test-*.txt")

	// then
	if err == nil {
		t.Error("expected error for HTTP 500 response, got nil")
	}
	if fileName != "" {
		t.Errorf("expected empty filename on error, got %s", fileName)
	}

	expectedErrorPrefix := "failed to download file from url. Status:"
	if !strings.Contains(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected error to contain %q, got %q", expectedErrorPrefix, err.Error())
	}
}

func TestDownloadFileHttpsToTemp_NotFoundError(t *testing.T) {
	t.Parallel()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
	defer server.Close()

	requestContext := common.NewRequestContext(context.Background())
	client := httpclient.NewHTTPClient(30 * time.Second)

	// when
	fileName, err := DownloadFileHttpsToTemp(requestContext, client, server.URL, "test-*.txt")

	// then
	if err == nil {
		t.Error("expected error for HTTP 404 response, got nil")
	}
	if fileName != "" {
		t.Errorf("expected empty filename on error, got %s", fileName)
	}

	expectedErrorPrefix := "failed to download file from url. Status:"
	if !strings.Contains(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected error to contain %q, got %q", expectedErrorPrefix, err.Error())
	}
}

func TestDownloadFileHttpsToTemp_HTTPClientError(t *testing.T) {
	t.Parallel()

	// given
	client := httpclient.NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())
	invalidURL := "http://invalid-host-that-does-not-exist.test:9999/path"

	// when
	fileName, err := DownloadFileHttpsToTemp(requestContext, client, invalidURL, "test-*.txt")

	// then
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
	if fileName != "" {
		t.Errorf("expected empty filename on error, got %s", fileName)
	}

	expectedErrorPrefix := "couldn't open download url"
	if !strings.Contains(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected error to contain %q, got %q", expectedErrorPrefix, err.Error())
	}
}

func TestDownloadFileHttpsToTemp_CreateTempFileError(t *testing.T) {
	// NOTE: Cannot use t.Parallel() here because this test modifies global state (TMPDIR env var)
	// which would cause race conditions with other parallel tests using os.CreateTemp()

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	client := httpclient.NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())

	// Mock os.CreateTemp to fail by setting TMPDIR to non-existent directory
	originalTmpDir := os.Getenv("TMPDIR")
	os.Setenv("TMPDIR", "/non/existent/directory")
	defer func() {
		if originalTmpDir != "" {
			os.Setenv("TMPDIR", originalTmpDir)
		} else {
			os.Unsetenv("TMPDIR")
		}
	}()

	// Test if we can actually trigger the temp file error
	if _, testErr := os.CreateTemp("", "test-*"); testErr == nil {
		t.Skip("Environment allows temp file creation even with invalid TMPDIR, skipping error test")
		return
	}

	// when
	fileName, err := DownloadFileHttpsToTemp(requestContext, client, server.URL, "test-*.txt")

	// then
	if err == nil {
		t.Error("expected error when temp file creation fails, got nil")
	}
	if fileName != "" {
		t.Errorf("expected empty filename on error, got %s", fileName)
	}

	expectedErrorPrefix := "couldn't create local download file"
	if !strings.Contains(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected error to contain %q, got %q", expectedErrorPrefix, err.Error())
	}
}

func TestDownloadFileHttpsToTemp_LargeFile(t *testing.T) {
	t.Parallel()

	// Check if temp file creation works in this environment
	if _, err := os.CreateTemp("", "test-*"); err != nil {
		t.Skipf("Cannot create temp files in this environment: %v", err)
		return
	}

	// given - simulate a larger file download
	largeContent := strings.Repeat("This is a test line for large file download.\n", 1000)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Use io.Copy to simulate streaming large content
		reader := strings.NewReader(largeContent)
		io.Copy(w, reader)
	}))
	defer server.Close()

	requestContext := common.NewRequestContext(context.Background())
	client := httpclient.NewHTTPClient(30 * time.Second)

	// when
	fileName, err := DownloadFileHttpsToTemp(requestContext, client, server.URL, "large-test-*.txt")

	// then
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if fileName == "" {
		t.Error("expected non-empty file name")
	}

	// Verify file exists and has correct content
	defer os.Remove(fileName)
	content, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("failed to read downloaded file: %v", err)
	}

	if string(content) != largeContent {
		t.Errorf("downloaded content length mismatch: expected %d bytes, got %d bytes", len(largeContent), len(content))
	}
}

func TestDownloadFileHttpsToTemp_EmptyFile(t *testing.T) {
	t.Parallel()

	// Check if temp file creation works in this environment
	if _, err := os.CreateTemp("", "test-*"); err != nil {
		t.Skipf("Cannot create temp files in this environment: %v", err)
		return
	}

	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Write empty response
	}))
	defer server.Close()

	client := httpclient.NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())

	// when
	fileName, err := DownloadFileHttpsToTemp(requestContext, client, server.URL, "empty-test-*.txt")

	// then
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if fileName == "" {
		t.Error("expected non-empty file name")
	}

	// Verify file exists and is empty
	defer os.Remove(fileName)
	content, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("failed to read downloaded file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("expected empty file, got %d bytes", len(content))
	}
}
