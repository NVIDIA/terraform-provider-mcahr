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
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	"terraform/terraform-provider/provider/tf/resource/resource/model"
)

// ResourceTranslatorCommon contains shared logic for resource translators
type ResourceTranslatorCommon struct{}

// ToAPIModelWithVersion creates a statement API model for the given TF model and API version
func (t *ResourceTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, tfModel *model.ResourceTFModel) (*statement.StatementInputAPIModel, error) {
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

func (t *ResourceTranslatorCommon) buildResourceStatement(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, prefix string, nameField string, valueField string, tfModel *model.ResourceTFModel) string {
	builder := utils.NewStatementBuilder(prefix, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField(nameField, tfModel.Name.ValueString(), "name").
		SetStringField(valueField, tfModel.Value.ValueString(), "value").
		SetStringField("description", tfModel.Description.ValueString(), "description").
		SetArrayField("params", utils.ListSliceFromTFModel(requestContext.Context, tfModel.Params), "params")

	return builder.Build()
}

func (t *ResourceTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, tfModel *model.ResourceTFModel) string {
	return t.buildResourceStatement(requestContext, translationData, "define_resource", "key", "val", tfModel)
}

func (t *ResourceTranslatorCommon) buildReadStatement(tfModel *model.ResourceTFModel) string {
	return fmt.Sprintf("list resources | name = \"%s\"", tfModel.Name.ValueString())
}

func (t *ResourceTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, tfModel *model.ResourceTFModel) string {
	return t.buildResourceStatement(requestContext, translationData, "update_resource", "resource_name", "value", tfModel)
}

func (t *ResourceTranslatorCommon) buildDeleteStatement(tfModel *model.ResourceTFModel) string {
	return fmt.Sprintf("delete %s", tfModel.Name.ValueString())
}
