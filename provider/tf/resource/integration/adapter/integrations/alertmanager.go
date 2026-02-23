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
	mapbuilder "terraform/terraform-provider/provider/tf/resource/integration/adapter/map_builder"
	modelupdater "terraform/terraform-provider/provider/tf/resource/integration/adapter/model_updater"

	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	"terraform/terraform-provider/provider/tf/resource/integration/adapter/utils"
	"terraform/terraform-provider/provider/tf/resource/integration/model"
)

type AlertmanagerDataAdapter struct{}

// Ensure the adapter implements the IntegrationDataAdapter interface
var _ adapterinterface.IntegrationDataAdapter = &AlertmanagerDataAdapter{}

func alertManagerTfModelFieldNames() []string {
	return []string{"external_url", "payload_paths"}
}

func (a *AlertmanagerDataAdapter) DataFieldNames() []string {
	return alertManagerTfModelFieldNames()
}

func (a *AlertmanagerDataAdapter) TFModelFieldNames() []string {
	return alertManagerTfModelFieldNames()
}

func (a *AlertmanagerDataAdapter) MapToTFModel(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, integrationData map[string]interface{}, tfModel *model.IntegrationTFModel) {

	modelupdater.NewModelUpdater(options, tfModel).
		UpdateStringField("external_url", &tfModel.ExternalUrl, utils.GetStringOrEmpty(requestContext, integrationData, "external_url")).
		UpdateSetField("payload_paths", &tfModel.PayloadPaths, utils.GetStringListOrEmpty(requestContext, integrationData, "payload_paths"))
}

func (a *AlertmanagerDataAdapter) TFModelToMap(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, tfModel *model.IntegrationTFModel) map[string]interface{} {

	return mapbuilder.NewMapBuilder(options.BackendVersion, options.CompatibilityOptions).
		SetField("external_url", "external_url", tfModel.ExternalUrl.ValueString()).
		SetField("payload_paths", "payload_paths", utils.StringListTFModel(requestContext, tfModel.PayloadPaths)).
		Build()
}
