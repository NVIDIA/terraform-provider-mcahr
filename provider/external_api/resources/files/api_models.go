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

package files

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// FileConfigV2 represents the file configuration in V2 API responses
type FileConfigV2 struct {
	Owner         string `json:"owner"`
	Mode          string `json:"mode"`
	Path          string `json:"path"`
	Checksum      string `json:"checksum"`
	ResourceQuery string `json:"resource_query"`
	FileData      string `json:"file_data"`
	URI           string `json:"uri"`
	PresignedUri  string `json:"presigned_uri"`
}

// FileMetadataV2 represents the file metadata in V2 API responses
type FileMetadataV2 struct {
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ConfigurationItem represents a single configuration item in the V2 API response
type ConfigurationItem struct {
	Config         FileConfigV2   `json:"config"`
	EntityType     string         `json:"entity_type"`
	EntityMetadata FileMetadataV2 `json:"entity_metadata"`
}

// Configurations represents the configurations section of the V2 API response
type Configurations struct {
	Count int                 `json:"count"`
	Items []ConfigurationItem `json:"items"`
}

// FileOutput represents the output section of the V2 API response
type FileOutput struct {
	Configurations Configurations `json:"configurations"`
}

// FileSummary represents the summary section of the V2 API response
type FileSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// FileResponseAPIModel represents the complete API response for file operations (V2)
type FileResponseAPIModel struct {
	Output  FileOutput  `json:"output"`
	Summary FileSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (f FileResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(f.Summary.Status, f.Summary.Errors)
}
