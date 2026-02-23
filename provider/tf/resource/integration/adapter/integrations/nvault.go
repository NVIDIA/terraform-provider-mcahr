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

type NvaultDataAdapter struct{}

// Ensure the adapter implements the IntegrationDataAdapter interface
var _ adapterinterface.IntegrationDataAdapter = &NvaultDataAdapter{}

func nvaultTfModelFieldNames() []string {
	return []string{"address", "namespace", "role_name", "jwt_auth_path"}
}
func (a *NvaultDataAdapter) DataFieldNames() []string {
	return nvaultTfModelFieldNames()
}

func (a *NvaultDataAdapter) TFModelFieldNames() []string {
	return nvaultTfModelFieldNames()
}

func (a *NvaultDataAdapter) MapToTFModel(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, integrationData map[string]interface{}, tfModel *model.IntegrationTFModel) {

	modelupdater.NewModelUpdater(options, tfModel).
		UpdateStringField("address", &tfModel.Address, utils.GetStringOrEmpty(requestContext, integrationData, "address")).
		UpdateStringField("namespace", &tfModel.Namespace, utils.GetStringOrEmpty(requestContext, integrationData, "namespace")).
		UpdateStringField("role_name", &tfModel.RoleName, utils.GetStringOrEmpty(requestContext, integrationData, "role_name")).
		UpdateStringField("jwt_auth_path", &tfModel.JWTAuthPath, utils.GetStringOrEmpty(requestContext, integrationData, "jwt_auth_path"))
}

func (a *NvaultDataAdapter) TFModelToMap(requestContext *common.RequestContext, options *adapterinterface.IntegrationDataAdapterOptions, tfModel *model.IntegrationTFModel) map[string]interface{} {

	return mapbuilder.NewMapBuilder(options.BackendVersion, options.CompatibilityOptions).
		SetField("address", "address", tfModel.Address.ValueString()).
		SetField("namespace", "namespace", tfModel.Namespace.ValueString()).
		SetField("role_name", "role_name", tfModel.RoleName.ValueString()).
		SetField("jwt_auth_path", "jwt_auth_path", tfModel.JWTAuthPath.ValueString()).
		Build()
}
