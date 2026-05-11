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

package integrations

import (
	"terraform/terraform-provider/provider/common"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	mapbuilder "terraform/terraform-provider/provider/tf/resource/integration/adapter/map_builder"
	modelupdater "terraform/terraform-provider/provider/tf/resource/integration/adapter/model_updater"

	"terraform/terraform-provider/provider/tf/resource/integration/adapter/utils"
	"terraform/terraform-provider/provider/tf/resource/integration/model"
)

type SmtpDataAdapter struct{}

// Ensure the adapter implements the IntegrationDataAdapter interface
var _ adapterinterface.IntegrationDataAdapter = &SmtpDataAdapter{}

func smtpTfModelFieldNames() []string {
	return []string{"smtp_host", "smtp_port", "username", "password", "sender", "max_emails_per_day", "max_emails_per_second"}
}

func (a *SmtpDataAdapter) DataFieldNames() []string {
	return smtpTfModelFieldNames()
}

func (a *SmtpDataAdapter) TFModelFieldNames() []string {
	return smtpTfModelFieldNames()
}

func (a *SmtpDataAdapter) MapToTFModel(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, integrationData map[string]interface{}, tfModel *model.IntegrationTFModel) {

	modelupdater.NewModelUpdater(options, tfModel).
		UpdateStringField("smtp_host", &tfModel.SmtpHost, utils.GetStringOrEmpty(requestContext, integrationData, "smtp_host")).
		UpdateInt64Field("smtp_port", &tfModel.SmtpPort, utils.GetInt64OrZero(requestContext, integrationData, "smtp_port")).
		UpdateStringField("username", &tfModel.Username, utils.GetStringOrEmpty(requestContext, integrationData, "username")).
		UpdateStringField("password", &tfModel.Password, utils.GetStringOrEmpty(requestContext, integrationData, "password")).
		UpdateStringField("sender", &tfModel.Sender, utils.GetStringOrEmpty(requestContext, integrationData, "sender")).
		UpdateInt64Field("max_emails_per_day", &tfModel.MaxEmailsPerDay, utils.GetInt64OrZero(requestContext, integrationData, "max_emails_per_day")).
		UpdateInt64Field("max_emails_per_second", &tfModel.MaxEmailsPerSecond, utils.GetInt64OrZero(requestContext, integrationData, "max_emails_per_second"))
}

func (a *SmtpDataAdapter) TFModelToMap(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, tfModel *model.IntegrationTFModel) map[string]interface{} {

	return mapbuilder.NewMapBuilder(options.BackendVersion, options.CompatibilityOptions).
		SetField("smtp_host", "smtp_host", tfModel.SmtpHost.ValueString()).
		SetField("smtp_port", "smtp_port", tfModel.SmtpPort.ValueInt64()).
		SetField("username", "username", tfModel.Username.ValueString()).
		SetField("password", "password", tfModel.Password.ValueString()).
		SetField("sender", "sender", tfModel.Sender.ValueString()).
		SetField("max_emails_per_day", "max_emails_per_day", tfModel.MaxEmailsPerDay.ValueInt64()).
		SetField("max_emails_per_second", "max_emails_per_second", tfModel.MaxEmailsPerSecond.ValueInt64()).
		Build()
}
