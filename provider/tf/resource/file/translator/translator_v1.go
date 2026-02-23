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

package translator

import (
	"fmt"

	"terraform/terraform-provider/provider/common"
	fileapi "terraform/terraform-provider/provider/external_api/resources/files"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FileTranslatorV1 handles translation between TF models and V1 API models for file resources
type FileTranslatorV1 struct {
	FileTranslatorCommon
}

var _ translator.Translator[*filetf.FileTFModel, *fileapi.FileResponseAPIModelV1] = &FileTranslatorV1{}

// ToTFModel converts a V1 API model to a TF model
func (t *FileTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *fileapi.FileResponseAPIModelV1) (*filetf.FileTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the file container regardless of operation type (define_file, update_file, get_file_class, delete_file)
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no file container found in V1 API response")
	}

	if len(container.FileClasses) == 0 {
		return nil, fmt.Errorf("no file classes found in V1 API response")
	}

	// Get the first file class, current implementation only supports one file to be returned by the API
	fileClass := container.FileClasses[0]

	// Build TF model from V1 file class
	tfModel := &filetf.FileTFModel{
		Name:            types.StringValue(fileClass.Name),
		Enabled:         types.BoolValue(fileClass.Enabled),
		DestinationPath: types.StringValue(fileClass.DestinationPath),
		ResourceQuery:   types.StringValue(fileClass.ResourceQuery),
		Checksum:        types.StringValue(fileClass.Checksum),
		FileData:        types.StringValue(fileClass.FileData),
		Description:     types.StringValue(fileClass.Description),
		Mode:            types.StringValue(fileClass.Mode),
		Owner:           types.StringValue(fileClass.Owner),
	}

	// Set file length if available
	if fileClass.FileLength != nil {
		tfModel.FileLength = types.Int64Value(int64(*fileClass.FileLength))
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (t *FileTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *filetf.FileTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
