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
	"strings"

	"terraform/terraform-provider/provider/common"
	integrationapi "terraform/terraform-provider/provider/external_api/resources/integrations"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	adapter "terraform/terraform-provider/provider/tf/resource/integration/adapter"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IntegrationTranslator struct {
	IntegrationTranslatorCommon
}

var _ translator.Translator[*integrationtf.IntegrationTFModel, *integrationapi.IntegrationResponseAPIModel] = &IntegrationTranslator{}

func (a *IntegrationTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *integrationapi.IntegrationResponseAPIModel) (*integrationtf.IntegrationTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Integrations.Items) == 0 {
		return nil, fmt.Errorf("no integrations found in API response")
	}

	tfModel := &integrationtf.IntegrationTFModel{
		Name:            types.StringValue(apiModel.Output.Integrations.Items[0].Name),
		ServiceName:     types.StringValue(strings.ToLower(apiModel.Output.Integrations.Items[0].IntegrationType)),
		SerialNumber:    types.StringValue(apiModel.Output.Integrations.Items[0].SerialNumber),
		Enabled:         types.BoolValue(apiModel.Output.Integrations.Items[0].Enabled),
		PermissionsUser: types.StringValue(apiModel.Output.Integrations.Items[0].PermissionsUser),

		// Set complex types to null to avoid incorrect initialization of the complex attributes
		// These will be overriden by the adapter if needed
		// "set" types need to be initialized with the element type to be valid
		PayloadPaths: types.ListNull(types.StringType),
	}

	// Add integration data fields (specific to the integration type) to the TF model
	integrationAdapterOptions := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       requestContext.BackendVersion,
		CompatibilityOptions: translationData.CompatibilityOptions,
	}
	adapter.MapToTFData(requestContext, integrationAdapterOptions, apiModel.Output.Integrations.Items[0].IntegrationData, tfModel)

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V2 backend
func (a *IntegrationTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *integrationtf.IntegrationTFModel) (*statement.StatementInputAPIModel, error) {
	return a.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
