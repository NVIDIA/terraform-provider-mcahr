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

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareVersions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		first    *BackendVersion
		second   *BackendVersion
		expected int
	}{
		{
			name:     "equal versions",
			first:    &BackendVersion{Major: 1, Minor: 2, Patch: 3},
			second:   &BackendVersion{Major: 1, Minor: 2, Patch: 3},
			expected: 0,
		},
		{
			name:     "first major version higher",
			first:    &BackendVersion{Major: 2, Minor: 0, Patch: 0},
			second:   &BackendVersion{Major: 1, Minor: 9, Patch: 9},
			expected: 1,
		},
		{
			name:     "first major version lower",
			first:    &BackendVersion{Major: 1, Minor: 9, Patch: 9},
			second:   &BackendVersion{Major: 2, Minor: 0, Patch: 0},
			expected: -1,
		},
		{
			name:     "first minor version higher",
			first:    &BackendVersion{Major: 1, Minor: 5, Patch: 0},
			second:   &BackendVersion{Major: 1, Minor: 2, Patch: 9},
			expected: 1,
		},
		{
			name:     "first minor version lower",
			first:    &BackendVersion{Major: 1, Minor: 2, Patch: 9},
			second:   &BackendVersion{Major: 1, Minor: 5, Patch: 0},
			expected: -1,
		},
		{
			name:     "first patch version higher",
			first:    &BackendVersion{Major: 1, Minor: 2, Patch: 5},
			second:   &BackendVersion{Major: 1, Minor: 2, Patch: 3},
			expected: 1,
		},
		{
			name:     "first patch version lower",
			first:    &BackendVersion{Major: 1, Minor: 2, Patch: 3},
			second:   &BackendVersion{Major: 1, Minor: 2, Patch: 5},
			expected: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// when
			result := CompareVersions(tc.first, tc.second)

			// then
			assert.Equal(t, tc.expected, result, "version comparison result mismatch for %v vs %v", tc.first, tc.second)
		})
	}
}

func TestExtractBackendVersion_ReleasePrefix(t *testing.T) {
	t.Parallel()

	// given
	version := "release-2.1.0"

	// when
	major, minor, patch, err := ExtractBackendVersion(version)

	// then
	require.NoError(t, err, "expected no error")
	assert.Equal(t, int64(2), major, "major version mismatch")
	assert.Equal(t, int64(1), minor, "minor version mismatch")
	assert.Equal(t, int64(0), patch, "patch version mismatch")
}

func TestExtractBackendVersion_ArmReleasePrefix(t *testing.T) {
	t.Parallel()

	// given
	version := "arm-release-3.2.1"

	// when
	major, minor, patch, err := ExtractBackendVersion(version)

	// then
	require.NoError(t, err, "expected no error")
	assert.Equal(t, int64(3), major, "major version mismatch")
	assert.Equal(t, int64(2), minor, "minor version mismatch")
	assert.Equal(t, int64(1), patch, "patch version mismatch")
}

func TestExtractBackendVersion_ReleaseInvalidFormat(t *testing.T) {
	t.Parallel()

	// given
	version := "release-invalid"

	// when
	major, minor, patch, err := ExtractBackendVersion(version)

	// then
	require.NoError(t, err, "expected no error even if format is invalid (returns 0.0.0)")
	assert.Equal(t, int64(0), major, "expected zero major")
	assert.Equal(t, int64(0), minor, "expected zero minor")
	assert.Equal(t, int64(0), patch, "expected zero patch")
}

func TestExtractBackendVersion_ArmReleaseInvalidFormat(t *testing.T) {
	t.Parallel()

	// given
	version := "arm-release-invalid"

	// when
	major, minor, patch, err := ExtractBackendVersion(version)

	// then
	require.NoError(t, err, "expected no error even if format is invalid (returns 0.0.0)")
	assert.Equal(t, int64(0), major, "expected zero major")
	assert.Equal(t, int64(0), minor, "expected zero minor")
	assert.Equal(t, int64(0), patch, "expected zero patch")
}

func TestExtractBackendVersion_OtherPrefix(t *testing.T) {
	t.Parallel()

	// given
	testCases := []string{
		"unknown-1.2.3",
		"master-1.2.3",
		"arm-master-1.2.3",
		"random-prefix-1.2.3",
		"no-prefix-1.2.3",
		"prefix-release-not-at-start", // e.g. "my-release-1.2.3"
	}

	for _, version := range testCases {
		t.Run(version, func(t *testing.T) {
			// when
			major, minor, patch, err := ExtractBackendVersion(version)

			// then
			require.NoError(t, err, "expected no error for non-release prefixes (dev build)")
			assert.Equal(t, int64(9999), major, "major version should be 9999")
			assert.Equal(t, int64(9999), minor, "minor version should be 9999")
			assert.Equal(t, int64(9999), patch, "patch version should be 9999")
		})
	}
}

func TestExtractVersionData_ValidVersion(t *testing.T) {
	t.Parallel()

	// given
	version := "release-1.2.3-extra-info"

	// when
	major, minor, patch, err := ExtractVersionData(version)

	// then
	require.NoError(t, err, "expected no error")
	assert.Equal(t, int64(1), major, "major version mismatch")
	assert.Equal(t, int64(2), minor, "minor version mismatch")
	assert.Equal(t, int64(3), patch, "patch version mismatch")
}

func TestExtractVersionData_SimpleVersion(t *testing.T) {
	t.Parallel()

	// given
	version := "1.2.3"

	// when
	major, minor, patch, err := ExtractVersionData(version)

	// then
	require.NoError(t, err, "expected no error")
	assert.Equal(t, int64(1), major, "major version mismatch")
	assert.Equal(t, int64(2), minor, "minor version mismatch")
	assert.Equal(t, int64(3), patch, "patch version mismatch")
}

func TestExtractVersionData_NoVersionPattern(t *testing.T) {
	t.Parallel()

	// given
	version := "no-version-here"

	// when
	major, minor, patch, err := ExtractVersionData(version)

	// then
	require.Error(t, err, "expected error for invalid version pattern")
	assert.Equal(t, int64(0), major, "expected zero major")
	assert.Equal(t, int64(0), minor, "expected zero minor")
	assert.Equal(t, int64(0), patch, "expected zero patch")
	assert.Equal(t, "couldn't find backend version number in string 'no-version-here'", err.Error(), "error message mismatch")
}

func TestExtractVersionData_PartialVersion(t *testing.T) {
	t.Parallel()

	// given
	version := "release-1.2"

	// when
	major, minor, patch, err := ExtractVersionData(version)

	// then
	require.Error(t, err, "expected error for partial version pattern")
	assert.Equal(t, int64(0), major, "expected zero major")
	assert.Equal(t, int64(0), minor, "expected zero minor")
	assert.Equal(t, int64(0), patch, "expected zero patch")
}

func TestExtractVersionData_MultipleVersions(t *testing.T) {
	t.Parallel()

	// given
	version := "release-1.2.3-and-4.5.6"

	// when
	major, minor, patch, err := ExtractVersionData(version)

	// then
	require.NoError(t, err, "expected no error")
	// Should extract the first version pattern
	assert.Equal(t, int64(1), major, "should extract first major version")
	assert.Equal(t, int64(2), minor, "should extract first minor version")
	assert.Equal(t, int64(3), patch, "should extract first patch version")
}

func TestExtractVersionData_LargeVersionNumbers(t *testing.T) {
	t.Parallel()

	// given
	version := "release-999.888.777"

	// when
	major, minor, patch, err := ExtractVersionData(version)

	// then
	require.NoError(t, err, "expected no error")
	assert.Equal(t, int64(999), major, "major version mismatch")
	assert.Equal(t, int64(888), minor, "minor version mismatch")
	assert.Equal(t, int64(777), patch, "patch version mismatch")
}

func TestExtractVersionData_ZeroVersions(t *testing.T) {
	t.Parallel()

	// given
	version := "release-0.0.0"

	// when
	major, minor, patch, err := ExtractVersionData(version)

	// then
	require.NoError(t, err, "expected no error")
	assert.Equal(t, int64(0), major, "major version mismatch")
	assert.Equal(t, int64(0), minor, "minor version mismatch")
	assert.Equal(t, int64(0), patch, "patch version mismatch")
}

func TestExtractVersionData_ParseIntError(t *testing.T) {
	t.Parallel()

	// given
	// Use a number larger than int64 max value to cause ParseInt to fail
	version := "release-99999999999999999999.1.0"

	// when
	major, minor, patch, err := ExtractVersionData(version)

	// then
	require.Error(t, err, "expected error for ParseInt failure")
	assert.Equal(t, int64(0), major, "expected zero major")
	assert.Equal(t, int64(0), minor, "expected zero minor")
	assert.Equal(t, int64(0), patch, "expected zero patch")
	assert.Equal(t, "couldn't parse backend version number in string 'release-99999999999999999999.1.0'", err.Error(), "error message mismatch")
}
