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
	integrationapi "terraform/terraform-provider/provider/external_api/resources/integrations"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	adapter "terraform/terraform-provider/provider/tf/resource/integration/adapter"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IntegrationTranslatorV1 handles translation for IntegrationResponseAPIModelV1
type IntegrationTranslatorV1 struct {
	IntegrationTranslatorCommon
}

var _ translator.Translator[*integrationtf.IntegrationTFModel, *integrationapi.IntegrationResponseAPIModelV1] = &IntegrationTranslatorV1{}

func (a *IntegrationTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *integrationapi.IntegrationResponseAPIModelV1) (*integrationtf.IntegrationTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the integration container regardless of operation type (define_integration, update_integration, get_integration_class)
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no integration container found in V1 API response")
	}

	if len(container.IntegrationClasses) == 0 {
		return nil, fmt.Errorf("no integration classes found in V1 API response")
	}

	// Get the first integration class, current implementation only supports one integration to be returned by the API
	integrationClass := container.IntegrationClasses[0]

	// Handle base fields
	tfModel := &integrationtf.IntegrationTFModel{
		Name:            types.StringValue(integrationClass.Name),
		ServiceName:     types.StringValue(integrationClass.ServiceName),
		SerialNumber:    types.StringValue(integrationClass.SerialNumber),
		Enabled:         types.BoolValue(integrationClass.Enabled),
		PermissionsUser: types.StringValue(integrationClass.PermissionsUser),
	}

	// Add integration data fields (specific to the integration type) to the TF model
	integrationAdapterOptions := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       requestContext.BackendVersion,
		CompatibilityOptions: translationData.CompatibilityOptions,
	}
	adapter.JSONToTFData(requestContext, integrationAdapterOptions, integrationClass.Params, tfModel)

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (a *IntegrationTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *integrationtf.IntegrationTFModel) (*statement.StatementInputAPIModel, error) {
	return a.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
