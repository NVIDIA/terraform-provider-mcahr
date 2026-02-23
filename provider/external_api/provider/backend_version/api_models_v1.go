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

package backend_version

import (
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
)

// BackendVersionResponseAPIModelV1 represents the API response for backend_version statement (V1)
type BackendVersionResponseAPIModelV1 struct {
	BackendVersion string                    `json:"backend_version"`
	Errors         *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"`
}

// GetErrors returns a string representation of any errors
func (b BackendVersionResponseAPIModelV1) GetErrors() string {
	if b.Errors != nil {
		return apicommon.FormatSyntaxErrors(b.Errors.Root)
	}
	return ""
}

// GetBackendVersion extracts the backend image tag from the response
func (b BackendVersionResponseAPIModelV1) GetBackendVersion() string {
	return b.BackendVersion
}
