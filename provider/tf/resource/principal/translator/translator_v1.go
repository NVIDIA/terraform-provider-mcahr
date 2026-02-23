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
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	principaltf "terraform/terraform-provider/provider/tf/resource/principal/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PrincipalTranslatorV1 handles translation between TF models and V1 API models for principal resources
type PrincipalTranslatorV1 struct {
	PrincipalTranslatorCommon
}

var _ translator.Translator[*principaltf.PrincipalTFModel, *principalapi.PrincipalResponseAPIModelV1] = &PrincipalTranslatorV1{}

// ToTFModel converts a V1 API model to a TF model
func (t *PrincipalTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *principalapi.PrincipalResponseAPIModelV1) (*principaltf.PrincipalTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the principal container regardless of operation type (define_principal, update_principal, get_principal_class, delete_principal)
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no principal container found in V1 API response")
	}

	if len(container.PrincipalClasses) == 0 {
		return nil, fmt.Errorf("no principal classes found in V1 API response")
	}

	// Get the first principal class, current implementation only supports one principal to be returned by the API
	principalClass := container.PrincipalClasses[0]

	// Build TF model from V1 principal class
	tfModel := &principaltf.PrincipalTFModel{
		Name:                 types.StringValue(principalClass.Name),
		Identity:             types.StringValue(principalClass.Identity),
		ActionLimit:          types.Int64Value(int64(principalClass.ActionLimit)),
		ExecuteLimit:         types.Int64Value(int64(principalClass.ExecuteLimit)),
		ConfigurePermission:  types.BoolValue(utils.IntToBool(principalClass.ConfigurePermission)),
		AdministerPermission: types.BoolValue(utils.IntToBool(principalClass.AdministerPermission)),
		IDPName:              types.StringValue(principalClass.IDPName),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (t *PrincipalTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *principaltf.PrincipalTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
