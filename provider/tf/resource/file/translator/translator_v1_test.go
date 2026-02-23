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

package translator

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	fileapi "terraform/terraform-provider/provider/external_api/resources/files"
	"terraform/terraform-provider/provider/tf/core/translator"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileTranslatorV1_ToTFModel_Success(t *testing.T) {
	t.Parallel()

	// given
	translatorV1 := &FileTranslatorV1{}
	fileLength := 12
	apiModel := &fileapi.FileResponseAPIModelV1{
		GetFileClass: &fileapi.FileContainerV1{
			FileClasses: []fileapi.FileClassV1{
				{
					Name:            "test_file",
					Enabled:         true,
					Owner:           "root",
					Mode:            "644",
					Checksum:        "abc123",
					Description:     "Test file description",
					ResourceQuery:   "host=test",
					DestinationPath: "/tmp/test.txt",
					FileLength:      &fileLength,
					FileData:        "test content",
				},
			},
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_file", result.Name.ValueString())
	assert.True(t, result.Enabled.ValueBool())
	assert.Equal(t, "/tmp/test.txt", result.DestinationPath.ValueString())
	assert.Equal(t, "host=test", result.ResourceQuery.ValueString())
	assert.Equal(t, "abc123", result.Checksum.ValueString())
	assert.Equal(t, "test content", result.FileData.ValueString())
	assert.Equal(t, "Test file description", result.Description.ValueString())
	assert.Equal(t, "644", result.Mode.ValueString())
	assert.Equal(t, "root", result.Owner.ValueString())
	assert.Equal(t, int64(12), result.FileLength.ValueInt64())
}

func TestFileTranslatorV1_ToTFModel_DefineFileContainer(t *testing.T) {
	t.Parallel()

	// given - test with DefineFile container instead of GetFileClass
	translatorV1 := &FileTranslatorV1{}
	apiModel := &fileapi.FileResponseAPIModelV1{
		DefineFile: &fileapi.FileContainerV1{
			FileClasses: []fileapi.FileClassV1{
				{
					Name:            "new_file",
					Enabled:         false,
					DestinationPath: "/tmp/new.txt",
					ResourceQuery:   "pod=test",
				},
			},
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "new_file", result.Name.ValueString())
	assert.False(t, result.Enabled.ValueBool())
	assert.Equal(t, "/tmp/new.txt", result.DestinationPath.ValueString())
	assert.Equal(t, "pod=test", result.ResourceQuery.ValueString())
}

func TestFileTranslatorV1_ToTFModel_NilModel(t *testing.T) {
	t.Parallel()

	// given
	translatorV1 := &FileTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := translatorV1.ToTFModel(requestContext, translationData, nil)

	// then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFileTranslatorV1_ToTFModel_NoContainer(t *testing.T) {
	t.Parallel()

	// given
	translatorV1 := &FileTranslatorV1{}
	apiModel := &fileapi.FileResponseAPIModelV1{
		// No containers set
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)

	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no file container found in V1 API response")
}

func TestFileTranslatorV1_ToTFModel_NoFileClasses(t *testing.T) {
	t.Parallel()

	// given
	translatorV1 := &FileTranslatorV1{}
	apiModel := &fileapi.FileResponseAPIModelV1{
		GetFileClass: &fileapi.FileContainerV1{
			FileClasses: []fileapi.FileClassV1{}, // Empty file classes
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)

	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no file classes found in V1 API response")
}

func TestFileTranslatorV1_ToTFModel_FileLength_Nil(t *testing.T) {
	t.Parallel()

	// given
	translatorV1 := &FileTranslatorV1{}
	apiModel := &fileapi.FileResponseAPIModelV1{
		UpdateFile: &fileapi.FileContainerV1{
			FileClasses: []fileapi.FileClassV1{
				{
					Name:       "test_file",
					FileLength: nil, // Nil file length
				},
			},
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	// FileLength should be unknown/null when nil
	assert.True(t, result.FileLength.IsNull() || result.FileLength.IsUnknown() || result.FileLength.ValueInt64() == 0)
}

func TestFileTranslatorV1_Interface_Compliance(t *testing.T) {
	// Verify that FileTranslatorV1 implements the required interface
	var _ translator.Translator[*filetf.FileTFModel, *fileapi.FileResponseAPIModelV1] = &FileTranslatorV1{}
}

func TestFileTranslatorV1_GetContainer_Priority(t *testing.T) {
	// Test the priority order of container selection in GetContainer method

	// Test GetFileClass has priority
	t.Run("GetFileClass priority", func(t *testing.T) {
		t.Parallel()

		translatorV1 := &FileTranslatorV1{}
		apiModel := &fileapi.FileResponseAPIModelV1{
			GetFileClass: &fileapi.FileContainerV1{
				FileClasses: []fileapi.FileClassV1{
					{Name: "get_file_class_test"},
				},
			},
			DefineFile: &fileapi.FileContainerV1{
				FileClasses: []fileapi.FileClassV1{
					{Name: "define_file_test"},
				},
			},
		}

		requestContext := common.NewRequestContext(context.Background())
		translationData := &translator.TranslationData{}

		result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)
		require.NoError(t, err)
		assert.Equal(t, "get_file_class_test", result.Name.ValueString())
	})

	// Test DefineFile when GetFileClass is nil
	t.Run("DefineFile fallback", func(t *testing.T) {
		t.Parallel()

		translatorV1 := &FileTranslatorV1{}
		apiModel := &fileapi.FileResponseAPIModelV1{
			DefineFile: &fileapi.FileContainerV1{
				FileClasses: []fileapi.FileClassV1{
					{Name: "define_file_test"},
				},
			},
		}

		requestContext := common.NewRequestContext(context.Background())
		translationData := &translator.TranslationData{}

		result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)
		require.NoError(t, err)
		assert.Equal(t, "define_file_test", result.Name.ValueString())
	})

	// Test UpdateFile fallback
	t.Run("UpdateFile fallback", func(t *testing.T) {
		t.Parallel()

		translatorV1 := &FileTranslatorV1{}
		apiModel := &fileapi.FileResponseAPIModelV1{
			UpdateFile: &fileapi.FileContainerV1{
				FileClasses: []fileapi.FileClassV1{
					{Name: "update_file_test"},
				},
			},
		}

		requestContext := common.NewRequestContext(context.Background())
		translationData := &translator.TranslationData{}

		result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)
		require.NoError(t, err)
		assert.Equal(t, "update_file_test", result.Name.ValueString())
	})

	// Test DeleteFile fallback
	t.Run("DeleteFile fallback", func(t *testing.T) {
		t.Parallel()

		translatorV1 := &FileTranslatorV1{}
		apiModel := &fileapi.FileResponseAPIModelV1{
			DeleteFile: &fileapi.FileContainerV1{
				FileClasses: []fileapi.FileClassV1{
					{Name: "delete_file_test"},
				},
			},
		}

		requestContext := common.NewRequestContext(context.Background())
		translationData := &translator.TranslationData{}

		result, err := translatorV1.ToTFModel(requestContext, translationData, apiModel)
		require.NoError(t, err)
		assert.Equal(t, "delete_file_test", result.Name.ValueString())
	})
}
