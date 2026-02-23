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
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/tf/core/translator"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileTranslatorCommon_ToAPIModelWithVersion_V1(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslatorCommon{}
	backendVersion := &version.BackendVersion{Major: 1, Minor: 0, Patch: 0}
	tfModel := &filetf.FileTFModel{
		Name:            types.StringValue("test_file"),
		Enabled:         types.BoolValue(true),
		DestinationPath: types.StringValue("/tmp/test.txt"),
		ResourceQuery:   types.StringValue("host=test"),
		Description:     types.StringValue("Test description"),
		FileLength:      types.Int64Value(12),
		Checksum:        types.StringValue("abc123"),
		Mode:            types.StringValue("644"),
		Owner:           types.StringValue("root"),
	}

	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).
		WithAPIVersion(common.V1).
		WithBackendVersion(backendVersion)
	translationData := &translator.TranslationData{}

	// when
	result, err := fileTranslator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, common.V1, result.APIVersion)
	assert.Equal(t, "define_file(name=\"test_file\", destination_path=\"/tmp/test.txt\", resource_query=\"host=test\", file_length=12, checksum=\"abc123\", description=\"Test description\", mode=\"644\", owner=\"root\", enabled=1)", result.Statement)
}

func TestFileTranslatorCommon_ToAPIModelWithVersion_V2(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslatorCommon{}
	backendVersion := &version.BackendVersion{Major: 2, Minor: 0, Patch: 0}
	tfModel := &filetf.FileTFModel{
		Name:            types.StringValue("test_file"),
		Enabled:         types.BoolValue(true),
		DestinationPath: types.StringValue("/tmp/test.txt"),
		ResourceQuery:   types.StringValue("host=test"),
	}

	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).
		WithAPIVersion(common.V2).
		WithBackendVersion(backendVersion)
	translationData := &translator.TranslationData{}

	// when
	result, err := fileTranslator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, common.V2, result.APIVersion)
	assert.Equal(t, "define_file(name=\"test_file\", destination_path=\"/tmp/test.txt\", resource_query=\"host=test\", file_length=0, checksum=\"\", description=\"\", mode=\"\", owner=\"\", enabled=1)", result.Statement)
}

func TestFileTranslatorCommon_BuildCreateStatement(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslatorCommon{}
	backendVersion := &version.BackendVersion{Major: 2, Minor: 0, Patch: 0}
	tfModel := &filetf.FileTFModel{
		Name:            types.StringValue("test_file"),
		Enabled:         types.BoolValue(true),
		DestinationPath: types.StringValue("/tmp/test.txt"),
		ResourceQuery:   types.StringValue("host=test"),
		Description:     types.StringValue("Test description"),
		FileLength:      types.Int64Value(12),
		Checksum:        types.StringValue("abc123"),
		Mode:            types.StringValue("644"),
		Owner:           types.StringValue("root"),
	}

	requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)
	translationData := &translator.TranslationData{}

	// when
	result := fileTranslator.buildCreateStatement(requestContext, translationData, tfModel)

	// then
	assert.Equal(t, "define_file(name=\"test_file\", destination_path=\"/tmp/test.txt\", resource_query=\"host=test\", file_length=12, checksum=\"abc123\", description=\"Test description\", mode=\"644\", owner=\"root\", enabled=1)", result)
}

func TestFileTranslatorCommon_BuildReadStatement(t *testing.T) {
	t.Parallel()

	// given
	translator := &FileTranslatorCommon{}
	tfModel := &filetf.FileTFModel{
		Name: types.StringValue("test_file"),
	}

	// when
	result := translator.buildReadStatement(tfModel)

	// then
	expectedStatement := `get_file_class(file_name="test_file")`
	assert.Equal(t, expectedStatement, result)
}

func TestFileTranslatorCommon_BuildReadStatement_EmptyName(t *testing.T) {
	t.Parallel()

	// given
	translator := &FileTranslatorCommon{}
	tfModel := &filetf.FileTFModel{
		Name: types.StringValue(""),
	}

	// when
	result := translator.buildReadStatement(tfModel)

	// then
	expectedStatement := `get_file_class(file_name="")`
	assert.Equal(t, expectedStatement, result)
}

func TestFileTranslatorCommon_BuildUpdateStatement(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslatorCommon{}
	backendVersion := &version.BackendVersion{Major: 2, Minor: 0, Patch: 0}
	tfModel := &filetf.FileTFModel{
		Name:            types.StringValue("test_file"),
		Enabled:         types.BoolValue(true),
		DestinationPath: types.StringValue("/tmp/updated.txt"),
		ResourceQuery:   types.StringValue("host=updated"),
		Description:     types.StringValue("Updated description"),
		FileLength:      types.Int64Value(15),
	}

	requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)
	translationData := &translator.TranslationData{}

	// when
	result := fileTranslator.buildUpdateStatement(requestContext, translationData, tfModel)

	// then
	assert.Equal(t, "update_file(name=\"test_file\", destination_path=\"/tmp/updated.txt\", resource_query=\"host=updated\", file_length=15, checksum=\"\", description=\"Updated description\", mode=\"\", owner=\"\", enabled=true)", result)
}

func TestFileTranslatorCommon_BuildDeleteStatement(t *testing.T) {
	t.Parallel()

	// given
	translator := &FileTranslatorCommon{}
	tfModel := &filetf.FileTFModel{
		Name: types.StringValue("test_file"),
	}

	// when
	result := translator.buildDeleteStatement(tfModel)

	// then
	expectedStatement := `delete_file(name="test_file")`
	assert.Equal(t, expectedStatement, result)
}

func TestFileTranslatorCommon_BuildDeleteStatement_SpecialCharacters(t *testing.T) {
	t.Parallel()

	// given
	translator := &FileTranslatorCommon{}
	tfModel := &filetf.FileTFModel{
		Name: types.StringValue("test_file_with_\"quotes\""),
	}

	// when
	result := translator.buildDeleteStatement(tfModel)

	// then
	// The string should be properly escaped
	assert.Equal(t, "delete_file(name=\"test_file_with_\"quotes\"\")", result)
}

func TestFileTranslatorCommon_BuildFileStatement(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslatorCommon{}
	backendVersion := &version.BackendVersion{Major: 2, Minor: 0, Patch: 0}
	tfModel := &filetf.FileTFModel{
		Name:            types.StringValue("test_file"),
		Enabled:         types.BoolValue(true),
		DestinationPath: types.StringValue("/tmp/test.txt"),
		ResourceQuery:   types.StringValue("host=test"),
		Description:     types.StringValue("Test description"),
		FileLength:      types.Int64Value(12),
		Checksum:        types.StringValue("abc123"),
		Mode:            types.StringValue("644"),
		Owner:           types.StringValue("root"),
	}

	requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)
	translationData := &translator.TranslationData{}

	// when
	result := fileTranslator.buildFileStatement(requestContext, translationData, "custom_operation", tfModel)

	// then
	assert.Equal(t, "custom_operation(name=\"test_file\", destination_path=\"/tmp/test.txt\", resource_query=\"host=test\", file_length=12, checksum=\"abc123\", description=\"Test description\", mode=\"644\", owner=\"root\", enabled=true)", result)
}

func TestFileTranslatorCommon_UnsupportedOperation(t *testing.T) {
	t.Parallel()

	// given
	fileTranslator := &FileTranslatorCommon{}
	backendVersion := &version.BackendVersion{Major: 2, Minor: 0, Patch: 0}
	tfModel := &filetf.FileTFModel{
		Name: types.StringValue("test_file"),
	}

	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.CrudOperation(999)).
		WithAPIVersion(common.V1).
		WithBackendVersion(backendVersion)
	translationData := &translator.TranslationData{}

	// when
	result, err := fileTranslator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}

func TestFileTranslatorCommon_EmptyStringHandling(t *testing.T) {
	t.Parallel()

	// Test with empty strings in various fields
	// given
	fileTranslator := &FileTranslatorCommon{}
	backendVersion := &version.BackendVersion{Major: 2, Minor: 0, Patch: 0}
	tfModel := &filetf.FileTFModel{
		Name:            types.StringValue("test_file"),
		Enabled:         types.BoolValue(false),
		DestinationPath: types.StringValue(""), // Empty string
		ResourceQuery:   types.StringValue(""), // Empty string
		Description:     types.StringValue(""), // Empty string
		FileLength:      types.Int64Value(0),   // Zero value
		Checksum:        types.StringValue(""), // Empty string
		Mode:            types.StringValue(""), // Empty string
		Owner:           types.StringValue(""), // Empty string
	}

	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).
		WithAPIVersion(common.V2).
		WithBackendVersion(backendVersion)
	translationData := &translator.TranslationData{}

	// when
	result, err := fileTranslator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should contain empty string fields as ""
	assert.Equal(t, "define_file(name=\"test_file\", destination_path=\"\", resource_query=\"\", file_length=0, checksum=\"\", description=\"\", mode=\"\", owner=\"\", enabled=0)", result.Statement)
}
