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
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"
)

// FileTranslatorCommon provides common functionality for file translators across API versions
type FileTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (t *FileTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *filetf.FileTFModel) (*statement.StatementInputAPIModel, error) {
	var stmt string

	switch requestContext.Operation {
	case common.Create:
		stmt = t.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = t.buildReadStatement(tfModel)
	case common.Update:
		stmt = t.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt = t.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}

	apiModel := &statement.StatementInputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (t *FileTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *filetf.FileTFModel) string {
	return t.buildFileStatement(requestContext, translationData, "define_file", tfModel)
}

func (t *FileTranslatorCommon) buildReadStatement(tfModel *filetf.FileTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_file_class(file_name=\"%s\")", name)
}

func (t *FileTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *filetf.FileTFModel) string {
	return t.buildFileStatement(requestContext, translationData, "update_file", tfModel)
}

func (t *FileTranslatorCommon) buildDeleteStatement(tfModel *filetf.FileTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_file(name=\"%s\")", name)
}

func (t *FileTranslatorCommon) buildFileStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *filetf.FileTFModel) string {
	// Build the file statement from the TF model using the builder pattern
	// Used for both define_file (create) and update_file (update) operations

	// File data
	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("name", tfModel.Name.ValueString(), "name").
		SetStringField("destination_path", tfModel.DestinationPath.ValueString(), "destination_path").
		SetStringField("resource_query", tfModel.ResourceQuery.ValueString(), "resource_query").
		SetField("file_length", tfModel.FileLength.ValueInt64(), "file_length").
		SetStringField("checksum", tfModel.Checksum.ValueString(), "checksum").
		SetStringField("description", tfModel.Description.ValueString(), "description").
		SetStringField("mode", tfModel.Mode.ValueString(), "mode").
		SetStringField("owner", tfModel.Owner.ValueString(), "owner")

	// There is a bug in the backend where define accepts only int and update accepts only bool for enabled field
	if statementName == "define_file" {
		builder = builder.SetField("enabled", utils.BoolToInt(tfModel.Enabled.ValueBool()), "enabled")
	} else {
		builder = builder.SetField("enabled", tfModel.Enabled.ValueBool(), "enabled")
	}

	return builder.Build()
}
