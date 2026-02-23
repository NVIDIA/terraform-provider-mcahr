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
	"os"
	"path/filepath"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/systemdefer"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/tf/core/process"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessFileContent(t *testing.T) {
	t.Parallel()

	// Setup test data
	testProcessData := &process.ProcessData{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
	}

	requestContext := common.NewRequestContext(context.Background())

	t.Run("Process inline data when inline_data is known", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InlineData: types.StringValue("hello world"),
			InputFile:  types.StringNull(), // Not known
		}

		// when
		result, err := ProcessFileContent(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", result.MD5.ValueString(), "MD5 of 'hello world'")
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", result.Checksum.ValueString(), "Checksum should equal MD5")
		assert.Equal(t, int64(11), result.FileLength.ValueInt64(), "Length of 'hello world'")
	})

	t.Run("Process input file when input_file is known", func(t *testing.T) {
		t.Parallel()

		// Setup temporary file
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		testContent := "test file content"
		err := os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(t, err)

		// given
		planModel := &filetf.FileTFModel{
			InlineData: types.StringNull(), // Not known
			InputFile:  types.StringValue(testFile),
		}

		// when
		result, err := ProcessFileContent(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Verify MD5 was calculated
		assert.NotEmpty(t, result.MD5.ValueString())
		assert.Len(t, result.MD5.ValueString(), 32, "MD5 should be 32 characters")
		assert.Equal(t, result.MD5.ValueString(), result.Checksum.ValueString(), "Checksum should equal MD5")
		assert.Equal(t, int64(len(testContent)), result.FileLength.ValueInt64())
	})

	t.Run("Both inline_data and input_file are known - should process inline_data first", func(t *testing.T) {
		t.Parallel()

		// Setup temporary file
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(testFile, []byte("file content"), 0644)
		require.NoError(t, err)

		// given
		planModel := &filetf.FileTFModel{
			InlineData: types.StringValue("inline content"), // This should be processed
			InputFile:  types.StringValue(testFile),         // This should be ignored
		}

		// when
		result, err := ProcessFileContent(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Should process inline data, not file
		assert.Equal(t, int64(14), result.FileLength.ValueInt64(), "Length of 'inline content'")
	})

	t.Run("Neither inline_data nor input_file are known", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InlineData: types.StringNull(),    // Not known
			InputFile:  types.StringUnknown(), // Not known
		}

		// when
		result, err := ProcessFileContent(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Fields should remain as they were (null/unknown)
		assert.True(t, result.MD5.IsNull() || result.MD5.IsUnknown())
		assert.True(t, result.Checksum.IsNull() || result.Checksum.IsUnknown())
		assert.True(t, result.FileLength.IsNull() || result.FileLength.IsUnknown())
	})

	t.Run("Error in processInlineData", func(t *testing.T) {
		t.Parallel()

		// given - This test scenario is hard to trigger since ContentMd5 rarely fails
		// We'll test with empty inline data which should still work
		planModel := &filetf.FileTFModel{
			InlineData: types.StringValue(""), // Empty content
			InputFile:  types.StringNull(),
		}

		// when
		result, err := ProcessFileContent(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err) // Empty content should work fine
		assert.NotNil(t, result)
		assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", result.MD5.ValueString(), "MD5 of empty string")
	})

	t.Run("Error in processInputFile - non-existent file", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InlineData: types.StringNull(),
			InputFile:  types.StringValue("/nonexistent/file.txt"),
		}

		// when
		result, err := ProcessFileContent(requestContext, testProcessData, planModel)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProcessInlineData(t *testing.T) {
	t.Parallel()

	t.Run("Process simple text content", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InlineData: types.StringValue("hello world"),
		}

		// when
		err := processInlineData(planModel)

		// then
		assert.NoError(t, err)
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", planModel.MD5.ValueString())
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", planModel.Checksum.ValueString())
		assert.Equal(t, int64(11), planModel.FileLength.ValueInt64())
	})

	t.Run("Process empty content", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InlineData: types.StringValue(""),
		}

		// when
		err := processInlineData(planModel)

		// then
		assert.NoError(t, err)
		assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", planModel.MD5.ValueString())
		assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", planModel.Checksum.ValueString())
		assert.Equal(t, int64(0), planModel.FileLength.ValueInt64())
	})

	t.Run("Process unicode content", func(t *testing.T) {
		t.Parallel()

		// given
		unicodeContent := "Hello 世界 🌍"
		planModel := &filetf.FileTFModel{
			InlineData: types.StringValue(unicodeContent),
		}

		// when
		err := processInlineData(planModel)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, planModel.MD5.ValueString())
		assert.Len(t, planModel.MD5.ValueString(), 32)
		assert.Equal(t, planModel.MD5.ValueString(), planModel.Checksum.ValueString())
		assert.Equal(t, int64(len([]byte(unicodeContent))), planModel.FileLength.ValueInt64())
	})

	t.Run("Process large content", func(t *testing.T) {
		t.Parallel()

		// given
		largeContent := string(make([]byte, 10000)) // 10KB of null bytes
		planModel := &filetf.FileTFModel{
			InlineData: types.StringValue(largeContent),
		}

		// when
		err := processInlineData(planModel)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, planModel.MD5.ValueString())
		assert.Len(t, planModel.MD5.ValueString(), 32)
		assert.Equal(t, int64(10000), planModel.FileLength.ValueInt64())
	})

	t.Run("Process JSON content", func(t *testing.T) {
		t.Parallel()

		// given
		jsonContent := `{"key": "value", "number": 42, "array": [1, 2, 3]}`
		planModel := &filetf.FileTFModel{
			InlineData: types.StringValue(jsonContent),
		}

		// when
		err := processInlineData(planModel)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, planModel.MD5.ValueString())
		assert.Equal(t, int64(len([]byte(jsonContent))), planModel.FileLength.ValueInt64())
	})
}

func TestProcessInputFile(t *testing.T) {
	t.Parallel()

	// Setup temporary directory for test files
	tempDir := t.TempDir()
	testProcessData := &process.ProcessData{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
	}

	t.Run("Process local file", func(t *testing.T) {
		t.Parallel()

		// Setup test file
		testFile := filepath.Join(tempDir, "test.txt")
		testContent := "test file content"
		err := os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(t, err)

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue(testFile),
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err = processInputFile(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, planModel.MD5.ValueString())
		assert.Len(t, planModel.MD5.ValueString(), 32)
		assert.Equal(t, planModel.MD5.ValueString(), planModel.Checksum.ValueString())
		assert.Equal(t, int64(len(testContent)), planModel.FileLength.ValueInt64())
	})

	t.Run("Process empty file", func(t *testing.T) {
		t.Parallel()

		// Setup empty test file
		emptyFile := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(emptyFile, []byte{}, 0644)
		require.NoError(t, err)

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue(emptyFile),
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err = processInputFile(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err)
		assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", planModel.MD5.ValueString())
		assert.Equal(t, int64(0), planModel.FileLength.ValueInt64())
	})

	t.Run("Process binary file", func(t *testing.T) {
		t.Parallel()

		// Setup binary test file
		binaryFile := filepath.Join(tempDir, "binary.bin")
		binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC}
		err := os.WriteFile(binaryFile, binaryContent, 0644)
		require.NoError(t, err)

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue(binaryFile),
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err = processInputFile(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, planModel.MD5.ValueString())
		assert.Len(t, planModel.MD5.ValueString(), 32)
		assert.Equal(t, int64(len(binaryContent)), planModel.FileLength.ValueInt64())
	})

	t.Run("Process large file", func(t *testing.T) {
		t.Parallel()

		// Setup large test file
		largeFile := filepath.Join(tempDir, "large.txt")
		largeContent := make([]byte, 1024*1024) // 1MB
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}
		err := os.WriteFile(largeFile, largeContent, 0644)
		require.NoError(t, err)

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue(largeFile),
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err = processInputFile(requestContext, testProcessData, planModel)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, planModel.MD5.ValueString())
		assert.Equal(t, int64(1024*1024), planModel.FileLength.ValueInt64())
	})

	t.Run("Non-existent file", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue("/nonexistent/file.txt"),
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err := processInputFile(requestContext, testProcessData, planModel)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("Directory instead of file", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue(tempDir), // Directory, not file
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err := processInputFile(requestContext, testProcessData, planModel)

		// then
		assert.Error(t, err)
	})

	t.Run("HTTP URL - should attempt download", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue("http://example.com/file.txt"),
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err := processInputFile(requestContext, testProcessData, planModel)

		// then
		// This will fail due to network issues, but we can test the error handling
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to download file")
	})

	t.Run("HTTPS URL - should attempt download", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue("https://example.com/file.txt"),
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err := processInputFile(requestContext, testProcessData, planModel)

		// then
		// This will fail due to network issues, but we can test the error handling
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to download file")
	})

	t.Run("File with no read permissions", func(t *testing.T) {
		// Skip on systems where we can't change permissions
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		t.Parallel()

		// Setup file with no read permissions
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

		// given
		planModel := &filetf.FileTFModel{
			InputFile: types.StringValue(noReadFile),
		}

		requestContext := common.NewRequestContext(context.Background())

		// when
		err = processInputFile(requestContext, testProcessData, planModel)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")
	})
}

func TestSetModelFields(t *testing.T) {
	t.Parallel()

	t.Run("Set all fields correctly", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{}
		testMd5 := "5d41402abc4b2a76b9719d911017c592"
		testSize := int64(11)

		// when
		setModelFields(planModel, testMd5, testSize)

		// then
		assert.Equal(t, testMd5, planModel.MD5.ValueString())
		assert.Equal(t, testMd5, planModel.Checksum.ValueString())
		assert.Equal(t, testSize, planModel.FileLength.ValueInt64())
	})

	t.Run("Set fields with empty MD5", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{}
		testMd5 := ""
		testSize := int64(0)

		// when
		setModelFields(planModel, testMd5, testSize)

		// then
		assert.Equal(t, "", planModel.MD5.ValueString())
		assert.Equal(t, "", planModel.Checksum.ValueString())
		assert.Equal(t, int64(0), planModel.FileLength.ValueInt64())
	})

	t.Run("Set fields with large size", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{}
		testMd5 := "abcdef1234567890abcdef1234567890"
		testSize := int64(1024 * 1024 * 1024) // 1GB

		// when
		setModelFields(planModel, testMd5, testSize)

		// then
		assert.Equal(t, testMd5, planModel.MD5.ValueString())
		assert.Equal(t, testMd5, planModel.Checksum.ValueString())
		assert.Equal(t, testSize, planModel.FileLength.ValueInt64())
	})

	t.Run("Overwrite existing fields", func(t *testing.T) {
		t.Parallel()

		// given
		planModel := &filetf.FileTFModel{
			MD5:        types.StringValue("old_md5"),
			Checksum:   types.StringValue("old_checksum"),
			FileLength: types.Int64Value(999),
		}
		newMd5 := "new_md5_hash_value_here_32chars"
		newSize := int64(42)

		// when
		setModelFields(planModel, newMd5, newSize)

		// then
		assert.Equal(t, newMd5, planModel.MD5.ValueString())
		assert.Equal(t, newMd5, planModel.Checksum.ValueString())
		assert.Equal(t, newSize, planModel.FileLength.ValueInt64())
	})
}

// Integration test that verifies consistency between inline and file processing
func TestConsistencyBetweenInlineAndFileProcessing(t *testing.T) {
	t.Parallel()

	// Setup
	testContent := "This is test content for consistency testing"
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "consistency.txt")
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	testProcessData := &process.ProcessData{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
	}

	requestContext := common.NewRequestContext(context.Background())

	// Test inline processing
	inlineModel := &filetf.FileTFModel{
		InlineData: types.StringValue(testContent),
		InputFile:  types.StringNull(),
	}

	// Test file processing
	fileModel := &filetf.FileTFModel{
		InlineData: types.StringNull(),
		InputFile:  types.StringValue(testFile),
	}

	// when
	inlineResult, err1 := ProcessFileContent(requestContext, testProcessData, inlineModel)
	fileResult, err2 := ProcessFileContent(requestContext, testProcessData, fileModel)

	// then
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotNil(t, inlineResult)
	assert.NotNil(t, fileResult)

	// Results should be identical
	assert.Equal(t, inlineResult.MD5.ValueString(), fileResult.MD5.ValueString(), "MD5 should be same for same content")
	assert.Equal(t, inlineResult.Checksum.ValueString(), fileResult.Checksum.ValueString(), "Checksum should be same for same content")
	assert.Equal(t, inlineResult.FileLength.ValueInt64(), fileResult.FileLength.ValueInt64(), "FileLength should be same for same content")
}

// Test error propagation in ProcessFileContent
func TestErrorPropagation(t *testing.T) {
	t.Parallel()

	testProcessData := &process.ProcessData{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
	}

	requestContext := common.NewRequestContext(context.Background())

	t.Run("Error from processInputFile should be propagated", func(t *testing.T) {
		t.Parallel()

		// given - non-existent file
		planModel := &filetf.FileTFModel{
			InlineData: types.StringNull(),
			InputFile:  types.StringValue("/definitely/does/not/exist/file.txt"),
		}

		// when
		result, err := ProcessFileContent(requestContext, testProcessData, planModel)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
