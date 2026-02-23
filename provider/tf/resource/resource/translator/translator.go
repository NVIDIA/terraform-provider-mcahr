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
	resourceapi "terraform/terraform-provider/provider/external_api/resources/resources"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	resourcetf "terraform/terraform-provider/provider/tf/resource/resource/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceTranslator handles translation between TF models and V2 API models
type ResourceTranslator struct {
	ResourceTranslatorCommon
}

// ToTFModel converts a V2 API model to a TF model
func (t *ResourceTranslator) ToTFModel(ctx *common.RequestContext, data *coretranslator.TranslationData, apiModel *resourceapi.ResourceResponseAPIModel) (*resourcetf.ResourceTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Symbols.Items) == 0 {
		return nil, fmt.Errorf("no resource configurations found in V2 API response")
	}

	configItem := apiModel.Output.Symbols.Items[0]

	tfModel := &resourcetf.ResourceTFModel{
		Name:        types.StringValue(configItem.Name),
		Description: types.StringValue(configItem.Description),
		Value:       types.StringValue(configItem.Formula),
	}

	params := coretranslator.ParseStringArray(configItem.Attributes.Params)
	tfModel.Params, _ = types.ListValueFrom(ctx.Context, types.StringType, params)

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V2 backend
func (t *ResourceTranslator) ToAPIModel(ctx *common.RequestContext, data *coretranslator.TranslationData, tfModel *resourcetf.ResourceTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(ctx, data, tfModel)
}
