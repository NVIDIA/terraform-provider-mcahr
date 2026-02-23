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

package adapterinterface

import (
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/tf/resource/integration/model"
)

type IntegrationDataAdapterOptions struct {
	BackendVersion       *version.BackendVersion
	CompatibilityOptions map[string]attribute.CompatibilityOptions
}

type IntegrationDataAdapter interface {
	// Names from the config/TF model
	TFModelFieldNames() []string
	// Names from the integration data (API request/response)
	DataFieldNames() []string
	MapToTFModel(requestContext *common.RequestContext, options *IntegrationDataAdapterOptions, integrationData map[string]interface{}, tfModel *model.IntegrationTFModel)
	TFModelToMap(requestContext *common.RequestContext, options *IntegrationDataAdapterOptions, tfModel *model.IntegrationTFModel) map[string]interface{}
}
