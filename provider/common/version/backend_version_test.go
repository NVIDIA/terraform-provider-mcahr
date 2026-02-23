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

func TestBackendVersion_IsValid(t *testing.T) {
	t.Parallel()

	// given
	testCases := []struct {
		name     string
		version  *BackendVersion
		expected bool
	}{
		{
			name: "valid version with positive major",
			version: &BackendVersion{
				Version: "1.2.3",
				Major:   1,
				Minor:   2,
				Patch:   3,
			},
			expected: true,
		},
		{
			name: "valid version with zero minor and patch",
			version: &BackendVersion{
				Version: "1.0.0",
				Major:   1,
				Minor:   0,
				Patch:   0,
			},
			expected: true,
		},
		{
			name: "invalid version with zero major",
			version: &BackendVersion{
				Version: "0.1.0",
				Major:   0,
				Minor:   1,
				Patch:   0,
			},
			expected: false,
		},
		{
			name: "invalid version with negative major",
			version: &BackendVersion{
				Version: "-1.0.0",
				Major:   -1,
				Minor:   0,
				Patch:   0,
			},
			expected: false,
		},
		{
			name: "invalid version with negative minor",
			version: &BackendVersion{
				Version: "1.-1.0",
				Major:   1,
				Minor:   -1,
				Patch:   0,
			},
			expected: false,
		},
		{
			name: "invalid version with negative patch",
			version: &BackendVersion{
				Version: "1.0.-1",
				Major:   1,
				Minor:   0,
				Patch:   -1,
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// when
			result := tc.version.IsValid()

			// then
			assert.Equal(t, tc.expected, result, "IsValid() result mismatch for version %v", tc.version)
		})
	}
}

func TestNewBackendVersion_ReleasePrefix(t *testing.T) {
	t.Parallel()

	// given
	version := "release-2.1.0"

	// when
	result := NewBackendVersion(version)

	// then
	require.NotNil(t, result, "expected non-nil result")
	assert.Equal(t, version, result.Version, "version string mismatch")
	assert.Equal(t, int64(2), result.Major, "major version mismatch")
	assert.Equal(t, int64(1), result.Minor, "minor version mismatch")
	assert.Equal(t, int64(0), result.Patch, "patch version mismatch")
}

func TestNewBackendVersion_ReleaseInvalidFormat(t *testing.T) {
	t.Parallel()

	// given
	version := "release-invalid-format"

	// when
	result := NewBackendVersion(version)

	// then
	assert.Nil(t, result, "expected nil result for invalid version format with release prefix")
}

func TestNewBackendVersion_ArmReleasePrefix(t *testing.T) {
	t.Parallel()

	// given
	version := "arm-release-3.2.1"

	// when
	result := NewBackendVersion(version)

	// then
	require.NotNil(t, result, "expected non-nil result")
	assert.Equal(t, version, result.Version, "version string mismatch")
	assert.Equal(t, int64(3), result.Major, "major version mismatch")
	assert.Equal(t, int64(2), result.Minor, "minor version mismatch")
	assert.Equal(t, int64(1), result.Patch, "patch version mismatch")
}

func TestNewBackendVersion_ArmReleaseInvalidFormat(t *testing.T) {
	t.Parallel()

	// given
	version := "arm-release-invalid-format"

	// when
	result := NewBackendVersion(version)

	// then
	assert.Nil(t, result, "expected nil result for invalid version format with arm-release prefix")
}

func TestNewBackendVersion_InvalidPrefix(t *testing.T) {
	t.Parallel()

	// given
	version := "invalid-prefix-1.2.3"

	// when
	result := NewBackendVersion(version)

	// then
	// "invalid-prefix" doesn't match stable/release/arm, so it's treated as dev build (9999.9999.9999)
	require.NotNil(t, result, "expected non-nil result for dev/master-like builds")
	assert.Equal(t, int64(9999), result.Major, "major version should be 9999 for unrecognized prefix")
	assert.Equal(t, int64(9999), result.Minor, "minor version should be 9999 for unrecognized prefix")
	assert.Equal(t, int64(9999), result.Patch, "patch version should be 9999 for unrecognized prefix")
}

func TestNewBackendVersion_InvalidPrefixAndFormat(t *testing.T) {
	t.Parallel()

	// given
	version := "invalid-prefix-and-format"

	// when
	result := NewBackendVersion(version)

	// then
	// "invalid-prefix" doesn't match stable/release/arm, so it's treated as dev build (9999.9999.9999)
	require.NotNil(t, result, "expected non-nil result for dev/master-like builds")
	assert.Equal(t, int64(9999), result.Major, "major version should be 9999 for unrecognized prefix")
	assert.Equal(t, int64(9999), result.Minor, "minor version should be 9999 for unrecognized prefix")
	assert.Equal(t, int64(9999), result.Patch, "patch version should be 9999 for unrecognized prefix")
}
