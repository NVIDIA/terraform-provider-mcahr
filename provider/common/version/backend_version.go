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
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	Base              int = 10
	BitSize           int = 64
	VersionLowerBound     = "0.0.0"
	VersionUpperBound     = "9999.9999.9999"
)

type BackendVersion struct {
	Version string
	Major   int64
	Minor   int64
	Patch   int64
}

func (v *BackendVersion) IsValid() bool {
	return v.Major > 0 && v.Minor >= 0 && v.Patch >= 0
}

func NewBackendVersion(version string) *BackendVersion {

	major, minor, patch, err := ExtractBackendVersion(version)
	if err != nil {
		tflog.Error(context.Background(), "Invalid backend version", map[string]any{"version": version})
		return nil
	}

	backendVersion := &BackendVersion{
		Version: version,
		Major:   major,
		Minor:   minor,
		Patch:   patch,
	}

	if !backendVersion.IsValid() {
		tflog.Error(context.Background(), "Invalid backend version", map[string]any{"version": version})
		return nil
	}

	return backendVersion
}

func ExtractBackendVersion(version string) (major int64, minor int64, patch int64, err error) {

	if strings.HasPrefix(version, "release") || strings.HasPrefix(version, "arm-release") {
		// parse out '\d+\.\d+.\d+' suffix
		major, minor, patch, err = ExtractVersionData(version)
	} else {
		// dev build, special case
		major, minor, patch = 9999, 9999, 9999
	}

	return major, minor, patch, nil
}
