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

func TestUploadFileHttps_Success(t *testing.T) {
	t.Parallel()

	// given - create a test file
	testContent := "test file content"
	tempFile, err := os.CreateTemp("", "test-upload-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write([]byte(testContent)); err != nil {
		t.Fatalf("failed to write test content: %v", err)
	}
	tempFile.Close() // Close to allow reading

	// Mock server that accepts uploads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT method, got %s", r.Method)
		}
		if r.Header.Get("x-ms-blob-type") != "BlockBlob" {
			t.Errorf("expected x-ms-blob-type header to be BlockBlob, got %s", r.Header.Get("x-ms-blob-type"))
		}

		// Read and verify content
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		if string(body) != testContent {
			t.Errorf("expected body %s, got %s", testContent, string(body))
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := httpclient.NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())

	// when
	err = UploadFileHttps(requestContext, client, tempFile.Name(), server.URL)

	// then
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUploadFileHttps_FileNotFound(t *testing.T) {
	t.Parallel()

	// given
	client := httpclient.NewHTTPClient(30 * time.Second)
	requestContext := common.NewRequestContext(context.Background())
	nonExistentFile := "/path/that/does/not/exist.txt"

	// when
	err := UploadFileHttps(requestContext, client, nonExistentFile, "http://example.com")

	// then
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
	expectedErrorPrefix := "couldn't open local upload file"
	if !strings.Contains(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected error to contain %q, got %q", expectedErrorPrefix, err.Error())
	}
}

func TestUploadFileHttps_HTTPError(t *testing.T) {
	t.Parallel()

	// given - create a test file
	tempFile, err := os.CreateTemp("", "test-upload-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write([]byte("test")); err != nil {
		t.Fatalf("failed to write test content: %v", err)
	}
	tempFile.Close()

	// Mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
	}))
	defer server.Close()

	client := httpclient.NewHTTPClient(0)
	requestContext := common.NewRequestContext(context.Background())

	// when
	err = UploadFileHttps(requestContext, client, tempFile.Name(), server.URL)

	// then
	if err == nil {
		t.Error("expected error for HTTP 500 response, got nil")
	}
	expectedErrorPrefix := "couldn't upload file, status:"
	if !strings.Contains(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected error to contain %q, got %q", expectedErrorPrefix, err.Error())
	}
}

func TestUploadFileHttps_SuccessWithStatus200(t *testing.T) {
	t.Parallel()

	// given - create a test file
	tempFile, err := os.CreateTemp("", "test-upload-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write([]byte("test")); err != nil {
		t.Fatalf("failed to write test content: %v", err)
	}
	tempFile.Close()

	// Mock server that returns 200 OK (some APIs return 200 instead of 201)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := httpclient.NewHTTPClient(0)
	requestContext := common.NewRequestContext(context.Background())

	// when
	err = UploadFileHttps(requestContext, client, tempFile.Name(), server.URL)

	// then
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUploadFileHttpsFromString_Success(t *testing.T) {
	t.Parallel()

	// given
	testContent := "test string content"
	var receivedContent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		receivedContent = string(body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := httpclient.NewHTTPClient(0)
	requestContext := common.NewRequestContext(context.Background())

	// when
	err := UploadFileHttpsFromString(requestContext, client, testContent, server.URL)

	// then
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if receivedContent != testContent {
		t.Errorf("expected content %s, got %s", testContent, receivedContent)
	}
}

func TestUploadFileHttpsFromString_CreateTempFileError(t *testing.T) {
	t.Parallel()

	// given
	client := httpclient.NewHTTPClient(0)
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
	err := UploadFileHttpsFromString(requestContext, client, "test", "http://example.com")

	// then
	if err == nil {
		t.Error("expected error when temp file creation fails, got nil")
	}
	expectedErrorPrefix := "couldn't create local upload file"
	if !strings.Contains(err.Error(), expectedErrorPrefix) {
		t.Errorf("expected error to contain %q, got %q", expectedErrorPrefix, err.Error())
	}
}
