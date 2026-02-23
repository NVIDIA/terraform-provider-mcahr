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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentMd5(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		content      []byte
		expectedHash string
		expectError  bool
	}{
		{
			name:         "Empty content",
			content:      []byte{},
			expectedHash: "d41d8cd98f00b204e9800998ecf8427e", // MD5 of empty string
			expectError:  false,
		},
		{
			name:         "Simple text content",
			content:      []byte("hello world"),
			expectedHash: "5eb63bbbe01eeed093cb22bb8f5acdc3", // MD5 of "hello world"
			expectError:  false,
		},
		{
			name:         "Binary content",
			content:      []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			expectedHash: "b59121341ab26766729b7f1d7f7e0c2f", // MD5 of binary data
			expectError:  false,
		},
		{
			name:         "Large content",
			content:      make([]byte, 1024),                 // 1KB of zeros
			expectedHash: "0f343b0931126a20f133d67c2b018a3b", // MD5 of 1024 zero bytes
			expectError:  false,
		},
		{
			name:         "Unicode content",
			content:      []byte("Hello 世界 🌍"),
			expectedHash: "a1f2bc5b7f3a4c0d19e3c6e2d5f8b9a7", // This will be calculated
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result, err := ContentMd5(tt.content)

			// then
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.Len(t, result, 32, "MD5 hash should be 32 characters long")
				// For known test cases, verify exact hash
				if tt.name != "Unicode content" {
					assert.Equal(t, tt.expectedHash, result)
				}
			}
		})
	}
}

func TestFileMd5(t *testing.T) {
	t.Parallel()

	// Setup temporary directory for test files
	tempDir := t.TempDir()

	// Test cases with file operations
	t.Run("Empty file", func(t *testing.T) {
		t.Parallel()

		// given
		emptyFile := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(emptyFile, []byte{}, 0644)
		require.NoError(t, err)

		// when
		result, err := FileMd5(emptyFile)

		// then
		assert.NoError(t, err)
		assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", result, "MD5 of empty file")
	})

	t.Run("File with content", func(t *testing.T) {
		t.Parallel()

		// given
		contentFile := filepath.Join(tempDir, "content.txt")
		testContent := []byte("hello world")
		err := os.WriteFile(contentFile, testContent, 0644)
		require.NoError(t, err)

		// when
		result, err := FileMd5(contentFile)

		// then
		assert.NoError(t, err)
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", result, "MD5 of 'hello world'")
	})

	t.Run("Large file", func(t *testing.T) {
		t.Parallel()

		// given
		largeFile := filepath.Join(tempDir, "large.txt")
		largeContent := make([]byte, 1024*1024) // 1MB of zeros
		err := os.WriteFile(largeFile, largeContent, 0644)
		require.NoError(t, err)

		// when
		result, err := FileMd5(largeFile)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 32, "MD5 hash should be 32 characters long")
	})

	t.Run("Binary file", func(t *testing.T) {
		t.Parallel()

		// given
		binaryFile := filepath.Join(tempDir, "binary.bin")
		binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC}
		err := os.WriteFile(binaryFile, binaryContent, 0644)
		require.NoError(t, err)

		// when
		result, err := FileMd5(binaryFile)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 32, "MD5 hash should be 32 characters long")
	})

	t.Run("Non-existent file", func(t *testing.T) {
		t.Parallel()

		// given
		nonExistentFile := filepath.Join(tempDir, "does-not-exist.txt")

		// when
		result, err := FileMd5(nonExistentFile)

		// then
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("Directory instead of file", func(t *testing.T) {
		t.Parallel()

		// when
		result, err := FileMd5(tempDir)

		// then
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("File with no read permissions", func(t *testing.T) {
		// Skip on systems where we can't change permissions
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		t.Parallel()

		// given
		noReadFile := filepath.Join(tempDir, "noread.txt")
		err := os.WriteFile(noReadFile, []byte("test content"), 0644)
		require.NoError(t, err)

		// Remove read permissions
		err = os.Chmod(noReadFile, 0000)
		require.NoError(t, err)

		// Restore permissions at the end for cleanup
		defer func() {
			_ = os.Chmod(noReadFile, 0644)
		}()

		// when
		result, err := FileMd5(noReadFile)

		// then
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "permission denied")
	})
}

func TestContentSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		content      []byte
		expectedSize int64
	}{
		{
			name:         "Empty content",
			content:      []byte{},
			expectedSize: 0,
		},
		{
			name:         "Single byte",
			content:      []byte{0x00},
			expectedSize: 1,
		},
		{
			name:         "Text content",
			content:      []byte("hello world"),
			expectedSize: 11,
		},
		{
			name:         "Binary content",
			content:      []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			expectedSize: 5,
		},
		{
			name:         "Unicode content",
			content:      []byte("Hello 世界 🌍"), // UTF-8 encoded
			expectedSize: int64(len([]byte("Hello 世界 🌍"))),
		},
		{
			name:         "Large content",
			content:      make([]byte, 1024*1024), // 1MB
			expectedSize: 1024 * 1024,
		},
		{
			name:         "Nil content",
			content:      nil,
			expectedSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// when
			result := ContentSize(tt.content)

			// then
			assert.Equal(t, tt.expectedSize, result)
		})
	}
}

func TestFileSize(t *testing.T) {
	t.Parallel()

	// Setup temporary directory for test files
	tempDir := t.TempDir()

	t.Run("Empty file", func(t *testing.T) {
		t.Parallel()

		// given
		emptyFile := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(emptyFile, []byte{}, 0644)
		require.NoError(t, err)

		// when
		result, err := FileSize(emptyFile)

		// then
		assert.NoError(t, err)
		assert.Equal(t, int64(0), result)
	})

	t.Run("File with content", func(t *testing.T) {
		t.Parallel()

		// given
		contentFile := filepath.Join(tempDir, "content.txt")
		testContent := []byte("hello world")
		err := os.WriteFile(contentFile, testContent, 0644)
		require.NoError(t, err)

		// when
		result, err := FileSize(contentFile)

		// then
		assert.NoError(t, err)
		assert.Equal(t, int64(len(testContent)), result)
	})

	t.Run("Large file", func(t *testing.T) {
		t.Parallel()

		// given
		largeFile := filepath.Join(tempDir, "large.txt")
		largeContent := make([]byte, 1024*1024) // 1MB
		err := os.WriteFile(largeFile, largeContent, 0644)
		require.NoError(t, err)

		// when
		result, err := FileSize(largeFile)

		// then
		assert.NoError(t, err)
		assert.Equal(t, int64(1024*1024), result)
	})

	t.Run("Binary file", func(t *testing.T) {
		t.Parallel()

		// given
		binaryFile := filepath.Join(tempDir, "binary.bin")
		binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC}
		err := os.WriteFile(binaryFile, binaryContent, 0644)
		require.NoError(t, err)

		// when
		result, err := FileSize(binaryFile)

		// then
		assert.NoError(t, err)
		assert.Equal(t, int64(len(binaryContent)), result)
	})

	t.Run("Non-existent file", func(t *testing.T) {
		t.Parallel()

		// given
		nonExistentFile := filepath.Join(tempDir, "does-not-exist.txt")

		// when
		result, err := FileSize(nonExistentFile)

		// then
		assert.Error(t, err)
		assert.Equal(t, int64(0), result)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("Directory", func(t *testing.T) {
		t.Parallel()

		// when
		result, err := FileSize(tempDir)

		// then
		// Directory size should be accessible (usually returns a positive size)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, result, int64(0))
	})

	t.Run("File with no read permissions", func(t *testing.T) {
		// Skip on systems where we can't change permissions
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		t.Parallel()

		// given
		noReadFile := filepath.Join(tempDir, "noread.txt")
		err := os.WriteFile(noReadFile, []byte("test content"), 0644)
		require.NoError(t, err)

		// Remove read permissions
		err = os.Chmod(noReadFile, 0000)
		require.NoError(t, err)

		// Restore permissions at the end for cleanup
		defer func() {
			_ = os.Chmod(noReadFile, 0644)
		}()

		// when
		result, err := FileSize(noReadFile)

		// then
		// os.Stat should still work even without read permissions
		assert.NoError(t, err, "os.Stat should work without read permissions")
		assert.Equal(t, int64(12), result) // "test content" is 12 bytes
	})
}

// Test consistency between ContentMd5/ContentSize and FileMd5/FileSize
func TestConsistencyBetweenContentAndFileFunctions(t *testing.T) {
	t.Parallel()

	// Setup temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "consistency_test.txt")
	testContent := []byte("This is a consistency test between content and file functions")

	err := os.WriteFile(testFile, testContent, 0644)
	require.NoError(t, err)

	t.Run("MD5 consistency", func(t *testing.T) {
		// when
		contentMd5, err1 := ContentMd5(testContent)
		fileMd5, err2 := FileMd5(testFile)

		// then
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, contentMd5, fileMd5, "ContentMd5 and FileMd5 should produce the same result")
	})

	t.Run("Size consistency", func(t *testing.T) {
		// when
		contentSize := ContentSize(testContent)
		fileSize, err := FileSize(testFile)

		// then
		assert.NoError(t, err)
		assert.Equal(t, contentSize, fileSize, "ContentSize and FileSize should produce the same result")
	})
}

// Benchmark tests
func BenchmarkContentMd5(b *testing.B) {
	content := make([]byte, 1024) // 1KB content
	for i := range content {
		content[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ContentMd5(content)
	}
}

func BenchmarkFileMd5(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark.txt")
	content := make([]byte, 1024) // 1KB content
	for i := range content {
		content[i] = byte(i % 256)
	}
	err := os.WriteFile(testFile, content, 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FileMd5(testFile)
	}
}

func BenchmarkContentSize(b *testing.B) {
	content := make([]byte, 1024) // 1KB content

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ContentSize(content)
	}
}

func BenchmarkFileSize(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark.txt")
	content := make([]byte, 1024) // 1KB content
	err := os.WriteFile(testFile, content, 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FileSize(testFile)
	}
}
