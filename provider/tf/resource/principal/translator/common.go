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
	principaltf "terraform/terraform-provider/provider/tf/resource/principal/model"
)

// PrincipalTranslatorCommon provides common functionality for principal translators across API versions
type PrincipalTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (t *PrincipalTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *principaltf.PrincipalTFModel) (*statement.StatementInputAPIModel, error) {
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

func (t *PrincipalTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *principaltf.PrincipalTFModel) string {
	return t.buildPrincipalStatement(requestContext, translationData, "define_principal", tfModel)
}

func (t *PrincipalTranslatorCommon) buildReadStatement(tfModel *principaltf.PrincipalTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_principal_class(name=\"%s\")", name)
}

func (t *PrincipalTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *principaltf.PrincipalTFModel) string {
	return t.buildPrincipalStatement(requestContext, translationData, "update_principal", tfModel)
}

func (t *PrincipalTranslatorCommon) buildDeleteStatement(tfModel *principaltf.PrincipalTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_principal(principal_name=\"%s\")", name)
}

func (t *PrincipalTranslatorCommon) buildPrincipalStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *principaltf.PrincipalTFModel) string {
	// Build the principal statement from the TF model using the builder pattern
	// Used for both define_principal (create) and update_principal (update) operations

	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("principal_name", tfModel.Name.ValueString(), "name").
		SetStringField("identity", tfModel.Identity.ValueString(), "identity").
		SetField("action_limit", tfModel.ActionLimit.ValueInt64(), "action_limit").
		SetField("execute_limit", tfModel.ExecuteLimit.ValueInt64(), "execute_limit").
		SetField("view_limit", tfModel.ViewLimit.ValueInt64(), "view_limit"). // removed in release-29.0.0
		SetField("configure_permission", utils.BoolToInt(tfModel.ConfigurePermission.ValueBool()), "configure_permission").
		SetField("administer_permission", utils.BoolToInt(tfModel.AdministerPermission.ValueBool()), "administer_permission").
		SetStringField("idp_name", tfModel.IDPName.ValueString(), "idp_name")

	return builder.Build()
}
