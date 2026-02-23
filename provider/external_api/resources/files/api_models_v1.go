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

// FileResponseAPIModelV1 represents the complete API response for file operations (V1)
type FileResponseAPIModelV1 struct {
	// Operation-specific containers only
	GetFileClass *FileContainerV1          `json:"get_file_class,omitempty"`
	DefineFile   *FileContainerV1          `json:"define_file,omitempty"`
	UpdateFile   *FileContainerV1          `json:"update_file,omitempty"`
	DeleteFile   *FileContainerV1          `json:"delete_file,omitempty"`
	Errors       *apicommon.SyntaxErrorsV1 `json:"errors,omitempty"` // Top-level syntax errors
}

// FileContainerV1 represents the container for file operations in V1 API
type FileContainerV1 struct {
	Error       apicommon.ErrorV1 `json:"error"`
	Errors      []string          `json:"errors,omitempty"` // Direct errors array for some error cases
	FileClasses []FileClassV1     `json:"file_classes"`
}

// Ensure ActionContainerV1 implements V1ErrorContainer interface
var _ apicommon.V1ErrorContainer = &FileContainerV1{}

// GetNestedError returns the nested error structure
func (c *FileContainerV1) GetNestedError() apicommon.ErrorV1 {
	return c.Error
}

// GetDirectErrors returns the direct validation errors array (none for actions)
func (c *FileContainerV1) GetDirectErrors() []string {
	return c.Errors
}

// FileClassV1 represents a file class in V1 API responses
type FileClassV1 struct {
	Enabled         bool   `json:"enabled"`
	Name            string `json:"name"`
	Owner           string `json:"owner"`
	Mode            string `json:"mode"`
	Checksum        string `json:"checksum"`
	Description     string `json:"description"`
	ResourceQuery   string `json:"resource_query"`
	DestinationPath string `json:"destination_path"`
	FileLength      *int   `json:"file_length"`
	FileData        string `json:"file_data"`
}

// GetContainer returns the appropriate container based on the operation type
func (f FileResponseAPIModelV1) GetContainer() *FileContainerV1 {
	if f.GetFileClass != nil {
		return f.GetFileClass
	}
	if f.DefineFile != nil {
		return f.DefineFile
	}
	if f.UpdateFile != nil {
		return f.UpdateFile
	}
	if f.DeleteFile != nil {
		return f.DeleteFile
	}
	return nil
}

// GetErrors returns a string representation of the API response errors
func (f FileResponseAPIModelV1) GetErrors() string {
	container := f.GetContainer()
	if container == nil {
		return "No container found in API response"
	}

	return apicommon.FormatV1ErrorsFromContainer(f.Errors, container)
}
