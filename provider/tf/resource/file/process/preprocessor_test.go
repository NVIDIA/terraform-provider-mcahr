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
	"math/big"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/systemdefer"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/tf/core/process"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── test helpers ────────────────────────────────────────────────────────────

var fileAttrTypes = map[string]tftypes.Type{
	"name":             tftypes.String,
	"destination_path": tftypes.String,
	"resource_query":   tftypes.String,
	"description":      tftypes.String,
	"enabled":          tftypes.Bool,
	"input_file":       tftypes.String,
	"inline_data":      tftypes.String,
	"md5":              tftypes.String,
	"file_data":        tftypes.String,
	"file_length":      tftypes.Number,
	"checksum":         tftypes.String,
	"mode":             tftypes.String,
	"owner":            tftypes.String,
}

var fileFrameworkSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"name":             schema.StringAttribute{Required: true},
		"destination_path": schema.StringAttribute{Required: true},
		"resource_query":   schema.StringAttribute{Required: true},
		"description":      schema.StringAttribute{Optional: true, Computed: true},
		"enabled":          schema.BoolAttribute{Optional: true, Computed: true},
		"input_file":       schema.StringAttribute{Optional: true, Computed: true},
		"inline_data":      schema.StringAttribute{Optional: true, Computed: true},
		"md5":              schema.StringAttribute{Optional: true, Computed: true},
		"file_data":        schema.StringAttribute{Computed: true},
		"file_length":      schema.Int64Attribute{Computed: true},
		"checksum":         schema.StringAttribute{Computed: true},
		"mode":             schema.StringAttribute{Optional: true, Computed: true},
		"owner":            schema.StringAttribute{Optional: true, Computed: true},
	},
}

type planValues struct {
	name            string
	destinationPath string
	resourceQuery   string
	description     *string
	enabled         bool
	inputFile       *string
	inlineData      *string
	md5             *string
	fileData        *string
	fileLength      *int64
	checksum        *string
	mode            *string
	owner           *string
}

func maybeString(s *string) tftypes.Value {
	if s == nil {
		return tftypes.NewValue(tftypes.String, nil)
	}
	return tftypes.NewValue(tftypes.String, *s)
}

func maybeNumber(n *int64) tftypes.Value {
	if n == nil {
		return tftypes.NewValue(tftypes.Number, nil)
	}
	return tftypes.NewValue(tftypes.Number, new(big.Float).SetInt64(*n))
}

func str(s string) *string { return &s }
func i64(n int64) *int64   { return &n }

func buildTestPlan(t *testing.T, v planValues) tfsdk.Plan {
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
	return tfsdk.Plan{Raw: raw, Schema: fileFrameworkSchema}
}

func newTestProcessData(t *testing.T, plan tfsdk.Plan) *process.ProcessData {
	t.Helper()
	return &process.ProcessData{
		Client:            client.NewPlatformClient("http://test.com", "test-key"),
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs:        map[string]string{},
		CreateRequest:     &resource.CreateRequest{Plan: plan},
		CreateResponse:    &resource.CreateResponse{},
	}
}

func newTestProcessDataDeadEndClient(t *testing.T, plan tfsdk.Plan) *process.ProcessData {
	t.Helper()
	// Points to a port that refuses connections — any real HTTP call will fail fast
	return &process.ProcessData{
		Client:            client.NewPlatformClient("http://127.0.0.1:1", "test-key"),
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs:        map[string]string{},
		CreateRequest:     &resource.CreateRequest{Plan: plan},
		CreateResponse:    &resource.CreateResponse{},
	}
}

func v1Ctx() *common.RequestContext {
	return common.NewRequestContext(context.Background()).
		WithOperation(common.Create).
		WithAPIVersion(common.V1).
		WithBackendVersion(version.NewBackendVersion("release-28.4.0"))
}

func v2Ctx() *common.RequestContext {
	return common.NewRequestContext(context.Background()).
		WithOperation(common.Create).
		WithAPIVersion(common.V2).
		WithBackendVersion(version.NewBackendVersion("release-29.1.0"))
}

// ─── isV1DeferredUpload ───────────────────────────────────────────────────────

func TestIsV1DeferredUpload(t *testing.T) {
	t.Parallel()

	t.Run("returns true when flag is set to 'true'", func(t *testing.T) {
		t.Parallel()
		data := &process.ProcessData{StringArgs: map[string]string{V1DeferredUploadKey: "true"}}
		assert.True(t, isV1DeferredUpload(data))
	})

	t.Run("returns false when flag is absent", func(t *testing.T) {
		t.Parallel()
		data := &process.ProcessData{StringArgs: map[string]string{}}
		assert.False(t, isV1DeferredUpload(data))
	})

	t.Run("returns false when flag has any value other than 'true'", func(t *testing.T) {
		t.Parallel()
		for _, v := range []string{"yes", "1", "True", "TRUE", ""} {
			data := &process.ProcessData{StringArgs: map[string]string{V1DeferredUploadKey: v}}
			assert.False(t, isV1DeferredUpload(data), "expected false for value %q", v)
		}
	})
}

// ─── getV1OriginalEnabled ─────────────────────────────────────────────────────

func TestGetV1OriginalEnabled(t *testing.T) {
	t.Parallel()

	t.Run("returns true when stored as 'true'", func(t *testing.T) {
		t.Parallel()
		data := &process.ProcessData{StringArgs: map[string]string{V1OriginalEnabledKey: "true"}}
		got, err := getV1OriginalEnabled(data)
		require.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("returns false when stored as 'false'", func(t *testing.T) {
		t.Parallel()
		data := &process.ProcessData{StringArgs: map[string]string{V1OriginalEnabledKey: "false"}}
		got, err := getV1OriginalEnabled(data)
		require.NoError(t, err)
		assert.False(t, got)
	})

	t.Run("returns error when key is missing", func(t *testing.T) {
		t.Parallel()
		data := &process.ProcessData{StringArgs: map[string]string{}}
		_, err := getV1OriginalEnabled(data)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing original enabled value")
	})
}

// ─── handleV1Create ───────────────────────────────────────────────────────────

func TestHandleV1Create_SetsDeferredUploadFlag(t *testing.T) {
	t.Parallel()

	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello world"),
	})
	data := newTestProcessData(t, plan)

	p := &FilePreProcessor{}
	result, err := p.handleV1Create(v1Ctx(), data)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "true", data.StringArgs[V1DeferredUploadKey])
}

func TestHandleV1Create_ForcesEnabledFalseWhenOriginallyTrue(t *testing.T) {
	t.Parallel()

	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello world"),
	})
	data := newTestProcessData(t, plan)

	p := &FilePreProcessor{}
	result, err := p.handleV1Create(v1Ctx(), data)

	require.NoError(t, err)
	assert.False(t, result.Enabled.ValueBool(),
		"enabled must be forced to false so define_file does not require file_data")
	assert.Equal(t, "true", data.StringArgs[V1OriginalEnabledKey],
		"original enabled=true must be saved for restore after upload")
}

func TestHandleV1Create_KeepsEnabledFalseWhenOriginallyFalse(t *testing.T) {
	t.Parallel()

	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: false, inlineData: str("hello world"),
	})
	data := newTestProcessData(t, plan)

	p := &FilePreProcessor{}
	result, err := p.handleV1Create(v1Ctx(), data)

	require.NoError(t, err)
	assert.False(t, result.Enabled.ValueBool())
	assert.Equal(t, "false", data.StringArgs[V1OriginalEnabledKey])
}

func TestHandleV1Create_ComputesContentMetadata(t *testing.T) {
	t.Parallel()

	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello world"),
	})
	data := newTestProcessData(t, plan)

	p := &FilePreProcessor{}
	result, err := p.handleV1Create(v1Ctx(), data)

	require.NoError(t, err)
	// MD5, checksum and file_length must be computed even though upload is deferred
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", result.MD5.ValueString(), "MD5 of 'hello world'")
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", result.Checksum.ValueString())
	assert.Equal(t, int64(11), result.FileLength.ValueInt64())
}

func TestHandleV1Create_DoesNotAttemptUpload(t *testing.T) {
	t.Parallel()

	// Use a dead-end client: any actual HTTP call (e.g. to get the presigned URL)
	// would return a connection-refused error, surfacing as a test failure.
	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello world"),
	})
	data := newTestProcessDataDeadEndClient(t, plan)

	p := &FilePreProcessor{}
	_, err := p.handleV1Create(v1Ctx(), data)

	// Must succeed — no HTTP calls should be made during V1 preprocessing
	assert.NoError(t, err)
}

// ─── PreProcessCreate routing ─────────────────────────────────────────────────

func TestPreProcessCreate_RoutesToV1PathForV1APIVersion(t *testing.T) {
	t.Parallel()

	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello world"),
	})
	// Dead-end client would fail if V2 path (which calls upload) were chosen
	data := newTestProcessDataDeadEndClient(t, plan)

	p := &FilePreProcessor{}
	_, err := p.PreProcessCreate(v1Ctx(), data)

	assert.NoError(t, err)
	assert.Equal(t, "true", data.StringArgs[V1DeferredUploadKey],
		"V1 path must set the deferred-upload flag")
}

func TestPreProcessCreate_RoutesToV2PathForV2APIVersion(t *testing.T) {
	t.Parallel()

	plan := buildTestPlan(t, planValues{
		name: "my_file", destinationPath: "/tmp/f.txt", resourceQuery: "host",
		enabled: true, inlineData: str("hello world"),
	})
	// Dead-end client: V2 path will try to upload (content processing + upload)
	// and will fail — that's expected here; we just verify routing by the absence
	// of the V1 deferred-upload flag and the presence of an upload error.
	data := newTestProcessDataDeadEndClient(t, plan)

	p := &FilePreProcessor{}
	_, err := p.PreProcessCreate(v2Ctx(), data)

	// V2 path attempts upload, which fails with dead-end client
	assert.Error(t, err, "V2 path should attempt upload and fail with dead-end client")
	assert.Empty(t, data.StringArgs[V1DeferredUploadKey],
		"V2 path must NOT set the deferred-upload flag")
}
