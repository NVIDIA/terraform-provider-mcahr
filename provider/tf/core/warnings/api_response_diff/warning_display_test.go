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

package apiresponsediff

import (
	"strings"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddWarningsToDiagnostics_CreateOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Create,
	}

	processData := &process.ProcessData{
		CreateResponse: &resource.CreateResponse{},
	}

	differences := []fieldDifference{
		{
			FieldName:     "name",
			PlanValue:     "original",
			ResponseValue: "modified",
		},
	}

	addWarningsToDiagnostics(requestContext, processData, differences)

	assert.True(t, processData.CreateResponse.Diagnostics.HasError() || len(processData.CreateResponse.Diagnostics.Warnings()) > 0)
	warnings := processData.CreateResponse.Diagnostics.Warnings()
	require.Len(t, warnings, 1)
	assert.Contains(t, warnings[0].Summary(), "API modified 1 field(s)")
	assert.Contains(t, warnings[0].Detail(), "name")
}

func TestAddWarningsToDiagnostics_UpdateOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Update,
	}

	processData := &process.ProcessData{
		UpdateResponse: &resource.UpdateResponse{},
	}

	differences := []fieldDifference{
		{
			FieldName:     "description",
			PlanValue:     "config-desc",
			ResponseValue: "api-desc",
		},
	}

	addWarningsToDiagnostics(requestContext, processData, differences)

	warnings := processData.UpdateResponse.Diagnostics.Warnings()
	require.Len(t, warnings, 1)
	assert.Contains(t, warnings[0].Detail(), "description")
}

func TestAddWarningsToDiagnostics_NilDiagnostics(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Read,
	}

	processData := &process.ProcessData{}

	differences := []fieldDifference{
		{
			FieldName:     "name",
			PlanValue:     "original",
			ResponseValue: "modified",
		},
	}

	addWarningsToDiagnostics(requestContext, processData, differences)
}

func TestAddWarningsToDiagnostics_MultipleDifferences(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Create,
	}

	processData := &process.ProcessData{
		CreateResponse: &resource.CreateResponse{},
	}

	differences := []fieldDifference{
		{
			FieldName:     "field1",
			PlanValue:     "value1",
			ResponseValue: "response1",
		},
		{
			FieldName:     "field2",
			PlanValue:     "value2",
			ResponseValue: "response2",
		},
		{
			FieldName:     "field3",
			PlanValue:     "value3",
			ResponseValue: "response3",
		},
	}

	addWarningsToDiagnostics(requestContext, processData, differences)

	warnings := processData.CreateResponse.Diagnostics.Warnings()
	require.Len(t, warnings, 1)
	assert.Contains(t, warnings[0].Summary(), "3 field(s)")
	assert.Contains(t, warnings[0].Detail(), "field1")
	assert.Contains(t, warnings[0].Detail(), "field2")
	assert.Contains(t, warnings[0].Detail(), "field3")
}

func TestGetDiagnostics_CreateOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Create,
	}

	processData := &process.ProcessData{
		CreateResponse: &resource.CreateResponse{},
	}

	diag := getDiagnostics(requestContext, processData)

	assert.NotNil(t, diag)
}

func TestGetDiagnostics_UpdateOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Update,
	}

	processData := &process.ProcessData{
		UpdateResponse: &resource.UpdateResponse{},
	}

	diag := getDiagnostics(requestContext, processData)

	assert.NotNil(t, diag)
}

func TestGetDiagnostics_ReadOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Read,
	}

	processData := &process.ProcessData{}

	diag := getDiagnostics(requestContext, processData)

	assert.Nil(t, diag)
}

func TestGetDiagnostics_DeleteOperation(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Delete,
	}

	processData := &process.ProcessData{}

	diag := getDiagnostics(requestContext, processData)

	assert.Nil(t, diag)
}

func TestConstructDifferenceMessage_ShortValues(t *testing.T) {
	var builder strings.Builder
	planValue := "short"
	responseValue := "value"

	constructDifferenceMessage(&builder, planValue, responseValue)

	result := builder.String()
	assert.Contains(t, result, "Plan:")
	assert.Contains(t, result, "API:")
	assert.Contains(t, result, planValue)
	assert.Contains(t, result, responseValue)
}

func TestConstructDifferenceMessage_EmptyPlanValue(t *testing.T) {
	var builder strings.Builder
	configValue := ""
	responseValue := "some-value"

	constructDifferenceMessage(&builder, configValue, responseValue)

	result := builder.String()
	assert.Contains(t, result, "(empty)")
	assert.Contains(t, result, "some-value")
}

func TestConstructDifferenceMessage_EmptyResponseValue(t *testing.T) {
	var builder strings.Builder
	configValue := "some-value"
	responseValue := ""

	constructDifferenceMessage(&builder, configValue, responseValue)

	result := builder.String()
	assert.Contains(t, result, "some-value")
	assert.Contains(t, result, "(empty)")
}

func TestConstructDifferenceMessage_MediumValues(t *testing.T) {
	var builder strings.Builder
	planValue := strings.Repeat("a", 150)
	responseValue := strings.Repeat("b", 150)

	constructDifferenceMessage(&builder, planValue, responseValue)

	result := builder.String()
	assert.Contains(t, result, "Plan:")
	assert.Contains(t, result, "API:")
	assert.Contains(t, result, planValue)
	assert.Contains(t, result, responseValue)
}

func TestConstructDifferenceMessage_LongValues(t *testing.T) {
	var builder strings.Builder
	planValue := strings.Repeat("x", 400)
	responseValue := strings.Repeat("y", 400)

	constructDifferenceMessage(&builder, planValue, responseValue)

	result := builder.String()
	// For long values, we show with context
	assert.Contains(t, result, "Plan:")
	assert.Contains(t, result, "API:")
}

func TestCommonPrefixLength_IdenticalStrings(t *testing.T) {
	result := commonPrefixLength("hello", "hello")
	assert.Equal(t, 5, result)
}

func TestCommonPrefixLength_NoCommonPrefix(t *testing.T) {
	result := commonPrefixLength("abc", "xyz")
	assert.Equal(t, 0, result)
}

func TestCommonPrefixLength_PartialMatch(t *testing.T) {
	result := commonPrefixLength("hello world", "hello universe")
	assert.Equal(t, 6, result)
}

func TestCommonPrefixLength_EmptyStrings(t *testing.T) {
	result := commonPrefixLength("", "")
	assert.Equal(t, 0, result)
}

func TestCommonPrefixLength_OneEmpty(t *testing.T) {
	result := commonPrefixLength("hello", "")
	assert.Equal(t, 0, result)
}

func TestCommonPrefixLength_LongStrings(t *testing.T) {
	prefix := strings.Repeat("a", 100)
	result := commonPrefixLength(prefix+"xxx", prefix+"yyy")
	assert.Equal(t, 100, result)
}

func TestCommonSuffixLength_IdenticalStrings(t *testing.T) {
	result := commonSuffixLength("hello", "hello")
	assert.Equal(t, 5, result)
}

func TestCommonSuffixLength_NoCommonSuffix(t *testing.T) {
	result := commonSuffixLength("abc", "xyz")
	assert.Equal(t, 0, result)
}

func TestCommonSuffixLength_PartialMatch(t *testing.T) {
	result := commonSuffixLength("world hello", "universe hello")
	assert.Equal(t, 6, result)
}

func TestCommonSuffixLength_EmptyStrings(t *testing.T) {
	result := commonSuffixLength("", "")
	assert.Equal(t, 0, result)
}

func TestCommonSuffixLength_OneEmpty(t *testing.T) {
	result := commonSuffixLength("hello", "")
	assert.Equal(t, 0, result)
}

func TestCommonSuffixLength_LongStrings(t *testing.T) {
	suffix := strings.Repeat("z", 100)
	result := commonSuffixLength("xxx"+suffix, "yyy"+suffix)
	assert.Equal(t, 100, result)
}

func TestCommonPrefixLengthWithChunk_DefaultChunkSize(t *testing.T) {
	result := commonPrefixLengthWithChunk("hello world", "hello universe", 16)
	assert.Equal(t, 6, result)
}

func TestCommonPrefixLengthWithChunk_SmallChunkSize(t *testing.T) {
	result := commonPrefixLengthWithChunk("hello world", "hello universe", 2)
	assert.Equal(t, 6, result)
}

func TestCommonPrefixLengthWithChunk_InvalidChunkSize(t *testing.T) {
	result := commonPrefixLengthWithChunk("hello", "hello", 0)
	assert.Equal(t, 5, result)
}

func TestCommonPrefixLengthWithChunk_NegativeChunkSize(t *testing.T) {
	result := commonPrefixLengthWithChunk("hello", "hello", -5)
	assert.Equal(t, 5, result)
}

func TestCommonSuffixLengthWithChunk_DefaultChunkSize(t *testing.T) {
	result := commonSuffixLengthWithChunk("world hello", "universe hello", 16)
	assert.Equal(t, 6, result)
}

func TestCommonSuffixLengthWithChunk_SmallChunkSize(t *testing.T) {
	result := commonSuffixLengthWithChunk("world hello", "universe hello", 2)
	assert.Equal(t, 6, result)
}

func TestCommonSuffixLengthWithChunk_InvalidChunkSize(t *testing.T) {
	result := commonSuffixLengthWithChunk("hello", "hello", 0)
	assert.Equal(t, 5, result)
}

func TestCommonSuffixLengthWithChunk_NegativeChunkSize(t *testing.T) {
	result := commonSuffixLengthWithChunk("hello", "hello", -5)
	assert.Equal(t, 5, result)
}

func TestTruncateDiff_ShortString(t *testing.T) {
	input := "short string"
	result := truncateDiff(input)
	assert.Equal(t, input, result)
}

func TestTruncateDiff_ExactlyAtLimit(t *testing.T) {
	input := strings.Repeat("a", 200)
	result := truncateDiff(input)
	assert.Equal(t, input, result)
}

func TestTruncateDiff_LongString(t *testing.T) {
	input := strings.Repeat("x", 500)
	result := truncateDiff(input)

	assert.Contains(t, result, "...")
	assert.Contains(t, result, "(500 chars)")
	assert.Less(t, len(result), len(input))
}

func TestTruncateDiff_VeryLongString(t *testing.T) {
	input := strings.Repeat("y", 10000)
	result := truncateDiff(input)

	assert.Contains(t, result, "...")
	assert.Contains(t, result, "(10000 chars)")
	assert.True(t, strings.HasPrefix(result, input[:100]))
	assert.True(t, strings.HasSuffix(result, input[len(input)-100:]))
}

func TestExtractSubstring_ValidRange(t *testing.T) {
	input := "hello world"
	result := extractSubstring(input, 0, 5)
	assert.Equal(t, "hello", result)
}

func TestExtractSubstring_MiddleRange(t *testing.T) {
	input := "hello world"
	result := extractSubstring(input, 6, 11)
	assert.Equal(t, "world", result)
}

func TestExtractSubstring_StartEqualsEnd(t *testing.T) {
	input := "hello world"
	result := extractSubstring(input, 5, 5)
	assert.Equal(t, "", result)
}

func TestExtractSubstring_StartGreaterThanEnd(t *testing.T) {
	input := "hello world"
	result := extractSubstring(input, 10, 5)
	assert.Equal(t, "", result)
}

func TestExtractSubstring_FullString(t *testing.T) {
	input := "hello world"
	result := extractSubstring(input, 0, len(input))
	assert.Equal(t, input, result)
}

func TestShowDiffWithContext_CompletelyDifferent(t *testing.T) {
	var builder strings.Builder
	planValue := "aaaaaaaaaa"
	responseValue := "bbbbbbbbbb"

	showDiffWithContext(&builder, planValue, responseValue)

	result := builder.String()
	assert.Contains(t, result, "Plan:")
	assert.Contains(t, result, "API:")
	assert.Contains(t, result, planValue)
	assert.Contains(t, result, responseValue)
}

func TestShowDiffWithContext_PrefixMatch(t *testing.T) {
	var builder strings.Builder
	planValue := "prefix_aaa"
	responseValue := "prefix_bbb"

	showDiffWithContext(&builder, planValue, responseValue)

	result := builder.String()
	assert.Contains(t, result, "prefix")
	assert.Contains(t, result, "Plan:")
	assert.Contains(t, result, "API:")
	assert.Contains(t, result, "aaa")
	assert.Contains(t, result, "bbb")
}

func TestShowDiffWithContext_SuffixMatch(t *testing.T) {
	var builder strings.Builder
	planValue := "aaa_suffix"
	responseValue := "bbb_suffix"

	showDiffWithContext(&builder, planValue, responseValue)

	result := builder.String()
	assert.Contains(t, result, "suffix")
	assert.Contains(t, result, "Plan:")
	assert.Contains(t, result, "API:")
	assert.Contains(t, result, "aaa")
	assert.Contains(t, result, "bbb")
}

func TestShowDiffWithContext_LongPrefix(t *testing.T) {
	var builder strings.Builder
	prefix := strings.Repeat("x", 100)
	configValue := prefix + "aaa"
	responseValue := prefix + "bbb"

	showDiffWithContext(&builder, configValue, responseValue)

	result := builder.String()
	assert.Contains(t, result, "...")
}

func TestAddWarningsToDiagnostics_EmptyValues(t *testing.T) {
	requestContext := &common.RequestContext{
		Operation: common.Create,
	}

	processData := &process.ProcessData{
		CreateResponse: &resource.CreateResponse{},
	}

	differences := []fieldDifference{
		{
			FieldName:     "field",
			PlanValue:     "",
			ResponseValue: "value",
		},
	}

	addWarningsToDiagnostics(requestContext, processData, differences)

	warnings := processData.CreateResponse.Diagnostics.Warnings()
	require.Len(t, warnings, 1)
	assert.Contains(t, warnings[0].Detail(), "(empty)")
}
