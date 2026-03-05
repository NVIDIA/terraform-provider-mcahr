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
	"fmt"
	"strings"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/process"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// addWarningsToDiagnostics adds warnings to the Terraform diagnostics when plan
// values differ from API response values.
func addWarningsToDiagnostics(requestContext *common.RequestContext, processData *process.ProcessData, differences []fieldDifference) {
	diagnostics := getDiagnostics(requestContext, processData)

	if diagnostics == nil {
		return
	}

	summary := fmt.Sprintf("API modified %d field(s)", len(differences))

	var details strings.Builder
	details.WriteString("Plan vs API response:\n\n")

	for _, diff := range differences {
		details.WriteString(fmt.Sprintf("%s:\n", diff.FieldName))
		constructDifferenceMessage(&details, diff.PlanValue, diff.ResponseValue)
		details.WriteString("\n")
	}

	details.WriteString("State will use Plan values.")

	diagnostics.AddWarning(summary, details.String())
}

// getDiagnostics returns the appropriate diagnostics collection for the current operation.
func getDiagnostics(requestContext *common.RequestContext, processData *process.ProcessData) *diag.Diagnostics {
	switch requestContext.Operation {
	case common.Create:
		return &processData.CreateResponse.Diagnostics
	case common.Update:
		return &processData.UpdateResponse.Diagnostics
	default:
		return nil
	}
}

// constructDifferenceMessage formats the difference between plan and response values.
// Uses different display strategies based on value length:
//   - Short values (<=80 chars): inline display
//   - Medium values (<=300 chars): full multi-line display
//   - Long values (>300 chars): smart diff showing only differences with context
func constructDifferenceMessage(builder *strings.Builder, planValue, responseValue string) {
	const maxInlineLength = 80
	const maxFullDisplayLength = 300

	planLen := len(planValue)
	responseLen := len(responseValue)

	planDisplay := planValue
	responseDisplay := responseValue
	if planValue == "" {
		planDisplay = "(empty)"
	}
	if responseValue == "" {
		responseDisplay = "(empty)"
	}

	bothShort := planLen <= maxInlineLength && responseLen <= maxInlineLength

	if bothShort {
		builder.WriteString(fmt.Sprintf("  Plan: %s\n", planDisplay))
		builder.WriteString(fmt.Sprintf("  API:  %s\n", responseDisplay))
		return
	}

	maxLen := max(planLen, responseLen)

	if maxLen <= maxFullDisplayLength {
		builder.WriteString(fmt.Sprintf("  Plan:\n    %s\n", planDisplay))
		builder.WriteString(fmt.Sprintf("  API:\n    %s\n", responseDisplay))
		return
	}

	showDiffWithContext(builder, planValue, responseValue)
}

// showDiffWithContext displays a smart diff for long strings.
// Shows both complete values with context, using clear separators.
func showDiffWithContext(builder *strings.Builder, planValue, responseValue string) {
	commonPrefixLen := commonPrefixLength(planValue, responseValue)
	commonSuffixLen := commonSuffixLength(planValue, responseValue)

	planLen := len(planValue)
	responseLen := len(responseValue)

	// Get context windows around the difference
	const maxContextChars = 50
	contextStart := max(0, commonPrefixLen-maxContextChars)
	contextEnd := min(max(planLen, responseLen)-commonSuffixLen+maxContextChars, max(planLen, responseLen))

	// Extract the relevant portions with context
	planDisplay := extractWithEllipsis(planValue, contextStart, min(contextEnd, planLen))
	responseDisplay := extractWithEllipsis(responseValue, contextStart, min(contextEnd, responseLen))

	// Show both values clearly separated
	builder.WriteString("  Plan:\n")
	builder.WriteString("    ")
	builder.WriteString(planDisplay)
	builder.WriteString("\n\n")

	builder.WriteString("  API:\n")
	builder.WriteString("    ")
	builder.WriteString(responseDisplay)
	builder.WriteString("\n")
}

// extractWithEllipsis extracts a substring and adds ellipsis if truncated.
func extractWithEllipsis(value string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(value) {
		end = len(value)
	}

	result := value[start:end]
	truncated := truncateDiff(result)

	var output string
	if start > 0 {
		output = "..." + truncated
	} else {
		output = truncated
	}

	if end < len(value) {
		output = output + "..."
	}

	return output
}

// commonPrefixLength returns the length of the common prefix between two strings.
// Uses chunked comparison for performance: compares chunkSize bytes at a time,
// then falls back to byte-by-byte comparison for the remainder.
//
// Performance: O(n) where n is the prefix length, but faster
// than naive byte-by-byte comparison due to optimized substring comparisons.
//
// Parameters:
//   - a, b: strings to compare
//   - chunkSize: number of bytes to compare at once (recommended: 8-32)
func commonPrefixLengthWithChunk(a, b string, chunkSize int) int {
	if chunkSize < 1 {
		chunkSize = 1
	}

	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	i := 0

	for i+chunkSize <= minLen {
		if a[i:i+chunkSize] != b[i:i+chunkSize] {
			break
		}
		i += chunkSize
	}

	for i < minLen && a[i] == b[i] {
		i++
	}

	return i
}

// commonSuffixLength returns the length of the common suffix between two strings.
// Uses chunked comparison for performance: compares chunkSize bytes at a time,
// then falls back to byte-by-byte comparison for the remainder.
//
// Parameters:
//   - a, b: strings to compare
//   - chunkSize: number of bytes to compare at once (see commonPrefixLengthWithChunk docs)
func commonSuffixLengthWithChunk(a, b string, chunkSize int) int {
	if chunkSize < 1 {
		chunkSize = 1
	}

	aLen := len(a)
	bLen := len(b)
	minLen := aLen
	if bLen < minLen {
		minLen = bLen
	}

	i := 0

	for i+chunkSize <= minLen {
		aStart := aLen - i - chunkSize
		bStart := bLen - i - chunkSize
		if a[aStart:aStart+chunkSize] != b[bStart:bStart+chunkSize] {
			break
		}
		i += chunkSize
	}

	for i < minLen {
		if a[aLen-1-i] != b[bLen-1-i] {
			break
		}
		i++
	}

	return i
}

// commonPrefixLength returns the length of the common prefix between two strings.
// Uses a default chunk size of 16 bytes for optimal performance on most inputs.
func commonPrefixLength(a, b string) int {
	return commonPrefixLengthWithChunk(a, b, 16)
}

// commonSuffixLength returns the length of the common suffix between two strings.
// Uses a default chunk size of 16 bytes for optimal performance on most inputs.
func commonSuffixLength(a, b string) int {
	return commonSuffixLengthWithChunk(a, b, 16)
}

// truncateDiff returns a string, potentially truncated if longer than 200 chars.
// Shows first 100 and last 100 characters with total length in between.
func truncateDiff(s string) string {
	if len(s) <= 200 {
		return s
	}
	return fmt.Sprintf("%s ... (%d chars) ... %s", s[:100], len(s), s[len(s)-100:])
}

// extractSubstring returns the substring between start and end positions.
// Uses half-open interval [start, end): includes start, excludes end.
// Returns empty string if start >= end.
func extractSubstring(value string, start, end int) string {
	if start >= end {
		return ""
	}
	return value[start:end]
}
