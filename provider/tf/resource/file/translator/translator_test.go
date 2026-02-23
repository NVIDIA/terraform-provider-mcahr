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
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	fileapi "terraform/terraform-provider/provider/external_api/resources/files"
	"terraform/terraform-provider/provider/tf/core/translator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileTranslator_ToTFModel_Success(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslator{}
	apiModel := &fileapi.FileResponseAPIModel{
		Output: fileapi.FileOutput{
			Configurations: fileapi.Configurations{
				Count: 1,
				Items: []fileapi.ConfigurationItem{
					{
						Config: fileapi.FileConfigV2{
							Path:          "/tmp/test.txt",
							ResourceQuery: "host=test",
							FileData:      "test content",
							Checksum:      "abc123",
							Mode:          "644",
							Owner:         "root",
						},
						EntityMetadata: fileapi.FileMetadataV2{
							Name:        "test_file",
							Description: "Test file description",
							Enabled:     true,
						},
					},
				},
			},
		},
		Summary: fileapi.FileSummary{
			Status: "success",
			Errors: []apicommon.Error{},
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := fileTranslator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_file", result.Name.ValueString())
	assert.Equal(t, "/tmp/test.txt", result.DestinationPath.ValueString())
	assert.Equal(t, "Test file description", result.Description.ValueString())
	assert.Equal(t, "host=test", result.ResourceQuery.ValueString())
	assert.True(t, result.Enabled.ValueBool())
	assert.Equal(t, "test content", result.FileData.ValueString())
	assert.Equal(t, int64(len("test content")), result.FileLength.ValueInt64())
	assert.Equal(t, "abc123", result.Checksum.ValueString())
	assert.Equal(t, "644", result.Mode.ValueString())
	assert.Equal(t, "root", result.Owner.ValueString())
}

func TestFileTranslator_ToTFModel_NilModel(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslator{}
	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := fileTranslator.ToTFModel(requestContext, translationData, nil)

	// then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFileTranslator_ToTFModel_NoConfigurations(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslator{}
	apiModel := &fileapi.FileResponseAPIModel{
		Output: fileapi.FileOutput{
			Configurations: fileapi.Configurations{
				Count: 0,
				Items: []fileapi.ConfigurationItem{},
			},
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := fileTranslator.ToTFModel(requestContext, translationData, apiModel)

	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no file configurations found in V2 API response")
}

// Test cases with empty/null values
func TestFileTranslator_ToTFModel_EmptyValues(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslator{}
	apiModel := &fileapi.FileResponseAPIModel{
		Output: fileapi.FileOutput{
			Configurations: fileapi.Configurations{
				Count: 1,
				Items: []fileapi.ConfigurationItem{
					{
						Config: fileapi.FileConfigV2{
							Path:          "",
							ResourceQuery: "",
							FileData:      "",
							Checksum:      "",
							Mode:          "",
							Owner:         "",
						},
						EntityMetadata: fileapi.FileMetadataV2{
							Name:        "test_file",
							Description: "",
							Enabled:     false,
						},
					},
				},
			},
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &translator.TranslationData{}

	// when
	result, err := fileTranslator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_file", result.Name.ValueString())
	assert.Equal(t, "", result.DestinationPath.ValueString())
	assert.Equal(t, "", result.Description.ValueString())
	assert.Equal(t, "", result.ResourceQuery.ValueString())
	assert.False(t, result.Enabled.ValueBool())
	assert.Equal(t, "", result.FileData.ValueString())
	assert.Equal(t, int64(0), result.FileLength.ValueInt64())
	assert.Equal(t, "", result.Checksum.ValueString())
	assert.Equal(t, "", result.Mode.ValueString())
	assert.Equal(t, "", result.Owner.ValueString())
}
