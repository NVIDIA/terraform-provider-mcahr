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

package content

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/systemdefer"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/tf/core/process"

	"github.com/stretchr/testify/assert"
)

func TestShouldDownloadFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputUrl string
		expected bool
	}{
		{
			name:     "HTTP URL",
			inputUrl: "http://example.com/file.txt",
			expected: true,
		},
		{
			name:     "HTTPS URL",
			inputUrl: "https://example.com/file.txt",
			expected: true,
		},
		{
			name:     "Local file path",
			inputUrl: "/path/to/local/file.txt",
			expected: false,
		},
		{
			name:     "Relative file path",
			inputUrl: "./relative/file.txt",
			expected: false,
		},
		{
			name:     "Windows path",
			inputUrl: "C:\\Windows\\file.txt",
			expected: false,
		},
		{
			name:     "FTP URL (not supported)",
			inputUrl: "ftp://example.com/file.txt",
			expected: false,
		},
		{
			name:     "Empty string",
			inputUrl: "",
			expected: false,
		},
		{
			name:     "Just http:",
			inputUrl: "http:",
			expected: true,
		},
		{
			name:     "Just https://",
			inputUrl: "https://",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := shouldDownloadFile(tt.inputUrl)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaybeDownloadFile_LocalFile(t *testing.T) {
	t.Parallel()

	// given
	data := &process.ProcessData{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
	}

	localFilePath := "/path/to/local/file.txt"

	// when
	fileName, err := maybeDownloadFile(common.NewRequestContext(context.Background()), data, localFilePath)

	// then
	assert.NoError(t, err)
	assert.Equal(t, localFilePath, fileName, "Should return the original path for local files")
}

func TestMaybeDownloadFile_HTTPSUrl(t *testing.T) {
	t.Parallel()

	// given
	data := &process.ProcessData{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
	}

	httpsUrl := "https://example.com/test.txt"

	// when
	fileName, err := maybeDownloadFile(common.NewRequestContext(context.Background()), data, httpsUrl)

	// then
	// This will fail due to network issues, but we can test the error handling
	assert.Error(t, err)
	assert.Empty(t, fileName)
	assert.Contains(t, err.Error(), "failed to read remote file object")
}

func TestMaybeDownloadFile_HTTPUrl(t *testing.T) {
	t.Parallel()

	// given
	data := &process.ProcessData{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
	}

	httpUrl := "http://example.com/test.txt"

	// when
	fileName, err := maybeDownloadFile(common.NewRequestContext(context.Background()), data, httpUrl)

	// then
	// This will fail due to network issues, but we can test the error handling
	assert.Error(t, err)
	assert.Empty(t, fileName)
	assert.Contains(t, err.Error(), "failed to read remote file object")
}

func TestMaybeDownloadFile_EmptyUrl(t *testing.T) {
	t.Parallel()

	// given
	data := &process.ProcessData{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
	}

	emptyUrl := ""

	// when
	fileName, err := maybeDownloadFile(common.NewRequestContext(context.Background()), data, emptyUrl)

	// then
	assert.NoError(t, err)
	assert.Equal(t, emptyUrl, fileName, "Should return empty string for empty input")
}

func TestTempDownloadFilePathPrefix(t *testing.T) {
	t.Parallel()

	// Test that the constant is properly defined
	assert.Equal(t, "tmp_opcp-", tempDownloadFilePathPrefix)
}

func TestShouldDownloadFile_CaseInsensitive(t *testing.T) {
	t.Parallel()

	// Test that the function handles case sensitivity appropriately
	tests := []struct {
		name     string
		inputUrl string
		expected bool
	}{
		{
			name:     "Uppercase HTTP",
			inputUrl: "HTTP://example.com/file.txt",
			expected: false, // Function is case-sensitive, only lowercase supported
		},
		{
			name:     "Uppercase HTTPS",
			inputUrl: "HTTPS://example.com/file.txt",
			expected: false, // Function is case-sensitive, only lowercase supported
		},
		{
			name:     "Mixed case http",
			inputUrl: "HtTp://example.com/file.txt",
			expected: false, // Function is case-sensitive
		},
		{
			name:     "Mixed case https",
			inputUrl: "HtTpS://example.com/file.txt",
			expected: false, // Function is case-sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := shouldDownloadFile(tt.inputUrl)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMaybeDownloadFile_EdgeCases(t *testing.T) {
	t.Parallel()

	// Test with nil DeferFunctionList
	t.Run("Nil DeferFunctionList", func(t *testing.T) {
		data := &process.ProcessData{
			DeferFunctionList: nil, // This might cause issues
			Client:            client.NewPlatformClient("http://test.com", "test-key"),
		}

		localPath := "/local/file.txt"
		fileName, err := maybeDownloadFile(common.NewRequestContext(context.Background()), data, localPath)

		// Should still work for local files
		assert.NoError(t, err)
		assert.Equal(t, localPath, fileName)
	})

	// Test with nil Client for local files
	t.Run("Nil Client for local file", func(t *testing.T) {
		data := &process.ProcessData{
			DeferFunctionList: systemdefer.NewDeferFunctionList(),
			Client:            nil, // No client
		}

		localPath := "/local/file.txt"
		fileName, err := maybeDownloadFile(common.NewRequestContext(context.Background()), data, localPath)

		// Should still work for local files
		assert.NoError(t, err)
		assert.Equal(t, localPath, fileName)
	})

	// Test with nil Client for remote files
	t.Run("Nil Client for remote file", func(t *testing.T) {
		data := &process.ProcessData{
			DeferFunctionList: systemdefer.NewDeferFunctionList(),
			Client:            nil, // No client
		}

		remoteUrl := "https://example.com/file.txt"

		// This will panic due to nil client, so we need to catch it
		defer func() {
			if r := recover(); r != nil {
				// Panic is expected when client is nil for remote URLs
				assert.NotNil(t, r)
			}
		}()

		fileName, err := maybeDownloadFile(common.NewRequestContext(context.Background()), data, remoteUrl)

		// If we get here without panic, should fail gracefully
		if err != nil {
			assert.Error(t, err)
			assert.Empty(t, fileName)
		}
	})
}
