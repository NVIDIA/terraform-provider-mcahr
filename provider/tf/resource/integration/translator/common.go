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
	adapter "terraform/terraform-provider/provider/tf/resource/integration/adapter"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"
)

// IntegrationTranslatorCommon contains shared functionality between V1 and V2 translators
type IntegrationTranslatorCommon struct{}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (a *IntegrationTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *integrationtf.IntegrationTFModel) (apiModel *statement.StatementInputAPIModel, err error) {
	var stmt string

	switch requestContext.Operation {
	case common.Create:
		stmt, err = a.buildCreateStatement(requestContext, translationData, tfModel)
		if err != nil {
			return nil, err
		}
	case common.Read:
		stmt = a.buildReadStatement(tfModel)
	case common.Update:
		stmt, err = a.buildUpdateStatement(requestContext, translationData, tfModel)
		if err != nil {
			return nil, err
		}
	case common.Delete:
		stmt = a.buildDeleteStatement(tfModel)
	default:
		err = fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}
	if err != nil {
		return nil, err
	}

	apiModel = &statement.StatementInputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (a *IntegrationTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *integrationtf.IntegrationTFModel) (string, error) {
	return a.buildIntegrationStatement(requestContext, translationData, "define_integration", tfModel)
}

func (a *IntegrationTranslatorCommon) buildReadStatement(tfModel *integrationtf.IntegrationTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_integration_class(integration_name=\"%s\")", name)
}

func (a *IntegrationTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *integrationtf.IntegrationTFModel) (string, error) {
	return a.buildIntegrationStatement(requestContext, translationData, "update_integration", tfModel)
}

func (a *IntegrationTranslatorCommon) buildDeleteStatement(tfModel *integrationtf.IntegrationTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_integration(integration_name=\"%s\")", name)
}

func (a *IntegrationTranslatorCommon) buildIntegrationStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *integrationtf.IntegrationTFModel) (string, error) {
	// Build the integration statement from the TF model using the builder pattern
	// Used for both define_integration (create) and update_integration (update) operations

	adapterOptions := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       requestContext.BackendVersion,
		CompatibilityOptions: translationData.CompatibilityOptions,
	}
	jsonParams, err := adapter.TFDataToJSON(requestContext, adapterOptions, tfModel)
	if err != nil {
		return "", err
	}

	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions).
		SetStringField("integration_name", tfModel.Name.ValueString(), "name").
		SetStringField("serial_number", tfModel.SerialNumber.ValueString(), "serial_number").
		SetField("enabled", tfModel.Enabled.ValueBool(), "enabled").
		SetStringField("permissions_user", tfModel.PermissionsUser.ValueString(), "permissions_user").
		SetField("params", jsonParams, "")

	if statementName == "define_integration" {
		builder.SetStringField("service_name", tfModel.ServiceName.ValueString(), "service_name")
	}

	return builder.Build(), nil
}
