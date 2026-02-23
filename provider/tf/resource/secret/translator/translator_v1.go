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
	secretapi "terraform/terraform-provider/provider/external_api/resources/secrets"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	secrettf "terraform/terraform-provider/provider/tf/resource/secret/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NVaultSecretTranslatorV1 handles translation between TF models and V1 API models for nvault secret resources
type NVaultSecretTranslatorV1 struct {
	NVaultSecretTranslatorCommon
}

var _ translator.Translator[*secrettf.NVaultSecretTFModel, *secretapi.NVaultSecretResponseAPIModelV1] = &NVaultSecretTranslatorV1{}

// ToTFModel converts a V1 API model to a TF model
func (t *NVaultSecretTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *secretapi.NVaultSecretResponseAPIModelV1) (*secrettf.NVaultSecretTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no nvault secret container found in V1 API response")
	}

	if len(container.Secrets) == 0 {
		return nil, fmt.Errorf("no nvault secrets found in V1 API response")
	}

	// Get the first secret, current implementation only supports one nvault secret to be returned by the API
	secret := container.Secrets[0]

	tfModel := &secrettf.NVaultSecretTFModel{
		Name:            types.StringValue(secret.Name),
		VaultSecretPath: types.StringValue(secret.SecretInfo.VaultSecretPath),
		VaultSecretKey:  types.StringValue(secret.SecretInfo.VaultSecretKey),
		IntegrationName: types.StringValue(secret.SecretInfo.IntegrationName),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to a V1 API model
func (t *NVaultSecretTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *secrettf.NVaultSecretTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
