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
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ExtractVersionData(version string) (major int64, minor int64, patch int64, err error) {
	verRe := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

	if verRe.MatchString(version) {
		match := verRe.FindStringSubmatch(version)

		major, errMajor := strconv.ParseInt(match[1], Base, BitSize)
		minor, errMinor := strconv.ParseInt(match[2], Base, BitSize)
		patch, errPatch := strconv.ParseInt(match[3], Base, BitSize)

		if errMajor != nil || errMinor != nil || errPatch != nil {
			return 0, 0, 0, fmt.Errorf("couldn't parse backend version number in string '%s'", version)
		}

		return major, minor, patch, nil
	}

	return 0, 0, 0, fmt.Errorf("couldn't find backend version number in string '%s'", version)
}

// CompareVersions compares two backend versions
// Returns: -1 if first < second, 0 if equal, 1 if first > second
func CompareVersions(first *BackendVersion, second *BackendVersion) int {

	if first.Major != second.Major {
		if first.Major < second.Major {
			return -1
		}
		return 1
	}

	if first.Minor != second.Minor {
		if first.Minor < second.Minor {
			return -1
		}
		return 1
	}

	if first.Patch < second.Patch {
		return -1
	} else if first.Patch > second.Patch {
		return 1
	}

	return 0
}

func IsFieldSupported(backendVersion *BackendVersion, minVersion, maxVersion string, skipVersionPrefixes []string) bool {

	if backendVersion == nil {
		// If no version info, assume latest version and include all fields
		return true
	}

	// Check if the backend version is in the skip version prefixes
	for _, skipVersionPrefix := range skipVersionPrefixes {
		if strings.HasPrefix(backendVersion.Version, skipVersionPrefix) {
			return false
		}
	}

	// Check minimum version requirement
	if minVersion != "" {
		minVersionRecord := NewBackendVersion(minVersion)
		if minVersionRecord != nil && CompareVersions(backendVersion, minVersionRecord) < 0 {
			return false
		}
	}

	// Check maximum version requirement (field deprecated/removed)
	if maxVersion != "" {
		maxVersionRecord := NewBackendVersion(maxVersion)
		if maxVersionRecord != nil && CompareVersions(backendVersion, maxVersionRecord) > 0 {
			return false
		}
	}

	return true
}
