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

type FileTranslator struct {
	FileTranslatorCommon
}

var _ translator.Translator[*filetf.FileTFModel, *fileapi.FileResponseAPIModel] = &FileTranslator{}

func (f *FileTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *fileapi.FileResponseAPIModel) (*filetf.FileTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no file configurations found in V2 API response")
	}

	configItem := apiModel.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	tfModel := &filetf.FileTFModel{
		Name:            types.StringValue(metadata.Name),
		DestinationPath: types.StringValue(config.Path),
		Description:     types.StringValue(metadata.Description),
		ResourceQuery:   types.StringValue(config.ResourceQuery),
		Enabled:         types.BoolValue(metadata.Enabled),
		FileData:        types.StringValue(config.FileData),
		FileLength:      types.Int64Value(int64(len(config.FileData))),
		Checksum:        types.StringValue(config.Checksum),
		Mode:            types.StringValue(config.Mode),
		Owner:           types.StringValue(config.Owner),
	}

	return tfModel, nil

}

// ToAPIModel converts a TF model to an API model for V2 backend
func (f *FileTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *filetf.FileTFModel) (*statement.StatementInputAPIModel, error) {
	return f.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
