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
	"terraform/terraform-provider/provider/tf/core/validators"
	"terraform/terraform-provider/provider/tf/resource/resource/model"
)

// ResourceTranslatorCommon contains shared logic for resource translators
type ResourceTranslatorCommon struct{}

// ToAPIModelWithVersion creates a statement API model for the given TF model and API version
func (t *ResourceTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, tfModel *model.ResourceTFModel) (*statement.InputAPIModel, error) {
	var stmt string
	var err error

	switch requestContext.Operation {
	case common.Create:
		stmt, err = t.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = t.buildReadStatement(tfModel)
	case common.Update:
		stmt, err = t.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt, err = t.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}
	if err != nil {
		return nil, err
	}

	apiModel := &statement.InputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (t *ResourceTranslatorCommon) buildResourceStatement(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, prefix string, nameField string, valueField string, tfModel *model.ResourceTFModel) (string, error) {
	params, err := utils.ListSliceFromTFModel(requestContext.Context, tfModel.Params)
	if err != nil {
		return "", fmt.Errorf("attribute \"params\": %w", err)
	}

	builder := utils.NewStatementBuilder(prefix, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField(nameField, tfModel.Name.ValueString(), "name").
		SetStringField(valueField, tfModel.Value.ValueString(), "value").
		SetStringField("description", tfModel.Description.ValueString(), "description").
		SetArrayField("params", params, "params")

	return builder.Build(), nil
}

func (t *ResourceTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, tfModel *model.ResourceTFModel) (string, error) {
	return t.buildResourceStatement(requestContext, translationData, "define_resource", "key", "val", tfModel)
}

func (t *ResourceTranslatorCommon) buildReadStatement(tfModel *model.ResourceTFModel) string {
	return fmt.Sprintf("list resources | name = %s", utils.EscapeString(tfModel.Name.ValueString()))
}

func (t *ResourceTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *coretranslator.TranslationData, tfModel *model.ResourceTFModel) (string, error) {
	return t.buildResourceStatement(requestContext, translationData, "update_resource", "resource_name", "value", tfModel)
}

// buildDeleteStatement builds `delete <name>`.
//
// The op-lang delete verb takes a bare identifier, not a quoted string, so this
// is the one statement site that cannot use EscapeString -- quoting would change
// what the backend receives. The name is validated instead: anything outside
// NameValidator's alphabet cannot carry statement syntax, and the schema already
// applies that same validator to `name`, so no configuration that applies today
// starts failing.
func (t *ResourceTranslatorCommon) buildDeleteStatement(tfModel *model.ResourceTFModel) (string, error) {
	name := tfModel.Name.ValueString()

	if !validators.IsValidName(name) {
		return "", fmt.Errorf(
			"resource name %q cannot be used in a delete statement: it must contain only "+
				"alphanumeric characters and underscores and must not start with a digit",
			name)
	}

	return fmt.Sprintf("delete %s", name), nil
}
