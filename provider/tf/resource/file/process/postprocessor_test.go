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

package process

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/systemdefer"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/tf/core/process"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildTestConfig builds a tfsdk.Config matching the file schema with the given plan values.
func buildTestConfig(t *testing.T, v planValues) tfsdk.Config {
	t.Helper()
	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: fileAttrTypes}, map[string]tftypes.Value{
		"name":             tftypes.NewValue(tftypes.String, v.name),
		"destination_path": tftypes.NewValue(tftypes.String, v.destinationPath),
		"resource_query":   tftypes.NewValue(tftypes.String, v.resourceQuery),
		"description":      maybeString(v.description),
		"enabled":          tftypes.NewValue(tftypes.Bool, v.enabled),
		"input_file":       maybeString(v.inputFile),
		"inline_data":      maybeString(v.inlineData),
		"md5":              maybeString(v.md5),
		"file_data":        maybeString(v.fileData),
		"file_length":      maybeNumber(v.fileLength),
		"checksum":         maybeString(v.checksum),
		"mode":             maybeString(v.mode),
		"owner":            maybeString(v.owner),
	})
	return tfsdk.Config{Raw: raw, Schema: fileFrameworkSchema}
}

// ─── PostProcessCreate: deferred-upload flag routing ─────────────────────────

func TestPostProcessCreate_SkipsDeferredUploadWhenFlagNotSet(t *testing.T) {
	t.Parallel()

	// No V1 deferred-upload flag → handleV1DeferredUpload must NOT be called.
	// setFieldsFromPrevious will be called; it reads from the config.
	vals := planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello"), md5: str("abc123"),
	}
	plan := buildTestPlan(t, vals)
	cfg := buildTestConfig(t, vals)

	data := &process.ProcessData{
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs:        map[string]string{}, // flag NOT set
		CreateRequest:     &resource.CreateRequest{Plan: plan, Config: cfg},
		CreateResponse:    &resource.CreateResponse{},
	}

	resultModel := &filetf.FileTFModel{
		Name:       types.StringValue("my_file"),
		InlineData: types.StringNull(),
		MD5:        types.StringNull(),
		InputFile:  types.StringNull(),
	}

	requestCtx := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).WithAPIVersion(common.V1)

	p := &FilePostProcessor{}
	// Must not error and must not attempt any HTTP call
	err := p.PostProcessCreate(requestCtx, data, resultModel)
	assert.NoError(t, err)
}

func TestPostProcessCreate_FailsWhenDeferredUploadFlagSetButClientUnreachable(t *testing.T) {
	t.Parallel()

	// Flag is set but the backend is unreachable → handleV1DeferredUpload is
	// called and should surface a network error (presigned URL fetch fails).
	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello"),
	})
	data := &process.ProcessData{
		Client:            client.NewPlatformClient("http://127.0.0.1:1", "test-key"),
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs: map[string]string{
			V1DeferredUploadKey:  "true",
			V1OriginalEnabledKey: "true",
		},
		CreateRequest:  &resource.CreateRequest{Plan: plan},
		CreateResponse: &resource.CreateResponse{},
	}

	resultModel := &filetf.FileTFModel{
		Name: types.StringValue("my_file"),
	}

	requestCtx := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).WithAPIVersion(common.V1)

	p := &FilePostProcessor{}
	err := p.PostProcessCreate(requestCtx, data, resultModel)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to upload file")
}

// ─── PostProcessCreate: setFieldsFromPrevious restoration ────────────────────

func TestPostProcessCreate_RestoresTFOnlyFieldsFromPlan(t *testing.T) {
	t.Parallel()

	// After a successful create, inline_data / input_file / md5 must be
	// restored from the plan because the API never returns them.
	inlineData := "hello world"
	md5 := "5eb63bbbe01eeed093cb22bb8f5acdc3"

	vals := planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: false, inlineData: &inlineData, md5: &md5,
	}
	plan := buildTestPlan(t, vals)
	cfg := buildTestConfig(t, vals)

	data := &process.ProcessData{
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs:        map[string]string{}, // no deferred upload
		CreateRequest:     &resource.CreateRequest{Plan: plan, Config: cfg},
		CreateResponse:    &resource.CreateResponse{},
	}

	// API response model — inline_data / md5 are empty (never returned by API)
	resultModel := &filetf.FileTFModel{
		Name:       types.StringValue("my_file"),
		InlineData: types.StringNull(),
		InputFile:  types.StringNull(),
		MD5:        types.StringNull(),
	}

	requestCtx := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).WithAPIVersion(common.V2)

	p := &FilePostProcessor{}
	err := p.PostProcessCreate(requestCtx, data, resultModel)

	require.NoError(t, err)
	assert.Equal(t, inlineData, resultModel.InlineData.ValueString(),
		"inline_data must be restored from plan")
	assert.Equal(t, md5, resultModel.MD5.ValueString(),
		"md5 must be restored from plan")
}

// ─── handleV1DeferredUpload: result model fields updated on success ───────────
// These tests cannot exercise the happy path without a live backend, so we
// verify the failure modes are clearly reported.

func TestHandleV1DeferredUpload_ReturnsErrorWhenPresignedURLFetchFails(t *testing.T) {
	t.Parallel()

	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello world"),
	})
	data := &process.ProcessData{
		Client:            client.NewPlatformClient("http://127.0.0.1:1", "test-key"),
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs: map[string]string{
			V1DeferredUploadKey:  "true",
			V1OriginalEnabledKey: "true",
		},
		CreateRequest:  &resource.CreateRequest{Plan: plan},
		CreateResponse: &resource.CreateResponse{},
	}

	result := &filetf.FileTFModel{Name: types.StringValue("my_file")}

	requestCtx := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).WithAPIVersion(common.V1)

	err := handleV1DeferredUpload(requestCtx, data, result)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to upload file content",
		"error must mention the upload step that failed")
}

func TestHandleV1DeferredUpload_DoesNotMutateResultOnFailure(t *testing.T) {
	t.Parallel()

	// When the upload fails the result model must not be partially mutated.
	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello world"),
	})
	data := &process.ProcessData{
		Client:            client.NewPlatformClient("http://127.0.0.1:1", "test-key"),
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs: map[string]string{
			V1DeferredUploadKey:  "true",
			V1OriginalEnabledKey: "true",
		},
		CreateRequest:  &resource.CreateRequest{Plan: plan},
		CreateResponse: &resource.CreateResponse{},
	}

	originalFileData := types.StringValue(":s3://bucket/original")
	originalEnabled := types.BoolValue(false)

	result := &filetf.FileTFModel{
		Name:     types.StringValue("my_file"),
		FileData: originalFileData,
		Enabled:  originalEnabled,
	}

	requestCtx := common.NewRequestContext(context.Background()).
		WithOperation(common.Create).WithAPIVersion(common.V1)

	_ = handleV1DeferredUpload(requestCtx, data, result)

	// Result should not have been mutated before the error
	assert.Equal(t, originalFileData, result.FileData)
	assert.Equal(t, originalEnabled, result.Enabled)
}

// ─── getFileAttribute ─────────────────────────────────────────────────────────

func TestGetFileAttribute_ReturnsErrorWhenClientUnreachable(t *testing.T) {
	t.Parallel()

	data := &process.ProcessData{
		Client:            client.NewPlatformClient("http://127.0.0.1:1", "test-key"),
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs:        map[string]string{},
	}

	requestCtx := common.NewRequestContext(context.Background()).WithAPIVersion(common.V1)

	_, err := getFileAttribute(requestCtx, data, "my_file", "uri")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get file attribute 'uri'")
}
