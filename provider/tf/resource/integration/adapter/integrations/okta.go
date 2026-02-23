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

type OktaDataAdapter struct{}

// Ensure the adapter implements the IntegrationDataAdapter interface
var _ adapterinterface.IntegrationDataAdapter = &OktaDataAdapter{}

func oktaTfModelFieldNames() []string {
	return []string{"idp_name", "cache_ttl_ms", "api_rate_limit"}
}

func (a *OktaDataAdapter) DataFieldNames() []string {
	return append(oktaTfModelFieldNames(), "api_token", "url")
}

func (a *OktaDataAdapter) TFModelFieldNames() []string {
	return append(oktaTfModelFieldNames(), "api_key", "api_url")
}

func (a *OktaDataAdapter) MapToTFModel(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, integrationData map[string]interface{}, tfModel *model.IntegrationTFModel) {

	modelupdater.NewModelUpdater(options, tfModel).
		UpdateStringField("idp_name", &tfModel.IDPName, utils.GetStringOrEmpty(requestContext, integrationData, "idp_name")).
		UpdateInt64Field("cache_ttl_ms", &tfModel.CacheTTLMs, utils.GetInt64OrZero(requestContext, integrationData, "cache_ttl_ms")).
		UpdateInt64Field("api_rate_limit", &tfModel.APIRateLimit, utils.GetInt64OrZero(requestContext, integrationData, "api_rate_limit")).
		UpdateStringField("api_token", &tfModel.APIKey, utils.GetStringOrEmpty(requestContext, integrationData, "api_token")).
		UpdateStringField("url", &tfModel.APIUrl, utils.GetStringOrEmpty(requestContext, integrationData, "url"))
}

func (a *OktaDataAdapter) TFModelToMap(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, tfModel *model.IntegrationTFModel) map[string]interface{} {

	return mapbuilder.NewMapBuilder(options.BackendVersion, options.CompatibilityOptions).
		SetField("idp_name", "idp_name", tfModel.IDPName.ValueString()).
		SetField("cache_ttl_ms", "cache_ttl_ms", tfModel.CacheTTLMs.ValueInt64()).
		SetField("api_rate_limit", "api_rate_limit", tfModel.APIRateLimit.ValueInt64()).
		SetField("api_token", "api_key", tfModel.APIKey.ValueString()).
		SetField("url", "api_url", tfModel.APIUrl.ValueString()).
		Build()
}
