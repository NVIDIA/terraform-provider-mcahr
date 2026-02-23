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

type DatadogDataAdapter struct{}

// Ensure the adapter implements the IntegrationDataAdapter interface
var _ adapterinterface.IntegrationDataAdapter = &DatadogDataAdapter{}

func datadogTfModelFieldNames() []string {
	return []string{"api_key", "api_url", "site_url", "app_key", "webhook_name"}
}

func (a *DatadogDataAdapter) DataFieldNames() []string {
	return datadogTfModelFieldNames()
}

func (a *DatadogDataAdapter) TFModelFieldNames() []string {
	return datadogTfModelFieldNames()
}

func (a *DatadogDataAdapter) MapToTFModel(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, integrationData map[string]interface{}, tfModel *model.IntegrationTFModel) {

	modelupdater.NewModelUpdater(options, tfModel).
		UpdateStringField("api_key", &tfModel.APIKey, utils.GetStringOrEmpty(requestContext, integrationData, "api_key")).
		UpdateStringField("api_url", &tfModel.APIUrl, utils.GetStringOrEmpty(requestContext, integrationData, "api_url")).
		UpdateStringField("site_url", &tfModel.SiteUrl, utils.GetStringOrEmpty(requestContext, integrationData, "site_url")).
		UpdateStringField("app_key", &tfModel.AppKey, utils.GetStringOrEmpty(requestContext, integrationData, "app_key")).
		UpdateStringField("webhook_name", &tfModel.WebhookName, utils.GetStringOrEmpty(requestContext, integrationData, "webhook_name"))
}

func (a *DatadogDataAdapter) TFModelToMap(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, tfModel *model.IntegrationTFModel) map[string]interface{} {

	return mapbuilder.NewMapBuilder(options.BackendVersion, options.CompatibilityOptions).
		SetField("api_key", "api_key", tfModel.APIKey.ValueString()).
		SetField("api_url", "api_url", tfModel.APIUrl.ValueString()).
		SetField("site_url", "site_url", tfModel.SiteUrl.ValueString()).
		SetField("app_key", "app_key", tfModel.AppKey.ValueString()).
		SetField("webhook_name", "webhook_name", tfModel.WebhookName.ValueString()).
		Build()
}
