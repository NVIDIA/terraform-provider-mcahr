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
	principalapi "terraform/terraform-provider/provider/external_api/resources/principals"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	translatorutils "terraform/terraform-provider/provider/tf/core/translator"
	principaltf "terraform/terraform-provider/provider/tf/resource/principal/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrincipalTranslator struct {
	PrincipalTranslatorCommon
}

var _ translatorutils.Translator[*principaltf.PrincipalTFModel, *principalapi.PrincipalResponseAPIModel] = &PrincipalTranslator{}

func (p *PrincipalTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translatorutils.TranslationData, apiModel *principalapi.PrincipalResponseAPIModel) (*principaltf.PrincipalTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	// Validate using items length as per requirement
	if len(apiModel.Output.AccessControl.Items) == 0 {
		return nil, fmt.Errorf("No principal access control items found in API response")
	}

	// Get the first access control item, current implementation only supports one principal to be returned by the API
	accessControlItem := apiModel.Output.AccessControl.Items[0]
	data := accessControlItem.Data

	tfModel := &principaltf.PrincipalTFModel{
		Name:                 types.StringValue(data.Name),
		Identity:             types.StringValue(data.Identity),
		ActionLimit:          types.Int64Value(int64(data.ActionLimit)),
		ExecuteLimit:         types.Int64Value(int64(data.ExecuteLimit)),
		ConfigurePermission:  types.BoolValue(translatorutils.IntToBool(data.ConfigurePermission)),
		AdministerPermission: types.BoolValue(translatorutils.IntToBool(data.AdministerPermission)),
		IDPName:              types.StringValue(data.IDPName),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V2 backend
func (p *PrincipalTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translatorutils.TranslationData, tfModel *principaltf.PrincipalTFModel) (*statement.StatementInputAPIModel, error) {
	return p.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
