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

type NVaultSecretTranslator struct {
	NVaultSecretTranslatorCommon
}

var _ translator.Translator[*secrettf.NVaultSecretTFModel, *secretapi.NVaultSecretResponseAPIModel] = &NVaultSecretTranslator{}

func (s *NVaultSecretTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *secretapi.NVaultSecretResponseAPIModel) (*secrettf.NVaultSecretTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no nvault secret configurations found in API response")
	}

	// Get the first configuration item, current implementation only supports one nvault secret to be returned by the API
	configItem := apiModel.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	tfModel := &secrettf.NVaultSecretTFModel{
		Name:            types.StringValue(metadata.Name),
		VaultSecretPath: types.StringValue(config.ExternalValue.VaultSecretPath),
		VaultSecretKey:  types.StringValue(config.ExternalValue.VaultSecretKey),
		IntegrationName: types.StringValue(config.ExternalValue.IntegrationName),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to a V2 API model
func (s *NVaultSecretTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *secrettf.NVaultSecretTFModel) (*statement.StatementInputAPIModel, error) {
	return s.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
