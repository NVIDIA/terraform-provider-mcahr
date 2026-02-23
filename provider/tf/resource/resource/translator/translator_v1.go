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
	utils "terraform/terraform-provider/provider/tf/core/translator"
	resourcetf "terraform/terraform-provider/provider/tf/resource/resource/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceTranslatorV1 handles translation between TF models and V1 API models
type ResourceTranslatorV1 struct {
	ResourceTranslatorCommon
}

// ToTFModel converts a V1 API model to a TF model
func (t *ResourceTranslatorV1) ToTFModel(ctx *common.RequestContext, data *coretranslator.TranslationData, apiModel *resourceapi.ResourceResponseAPIModelV1) (*resourcetf.ResourceTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no resource container found in V1 API response")
	}

	if len(container.Symbol) == 0 {
		return nil, fmt.Errorf("no resource symbols found in V1 API response")
	}

	symbol := container.Symbol[0]

	tfModel := &resourcetf.ResourceTFModel{
		Name:        types.StringValue(symbol.Name),
		Description: types.StringValue(symbol.Attributes.Description),
		Value:       types.StringValue(symbol.Formula),
	}

	// Handle params from attributes (it's a JSON string)
	params := utils.ParseStringArray(symbol.Attributes.Params)
	tfModel.Params, _ = types.ListValueFrom(ctx.Context, types.StringType, params)

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (t *ResourceTranslatorV1) ToAPIModel(ctx *common.RequestContext, data *coretranslator.TranslationData, tfModel *resourcetf.ResourceTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(ctx, data, tfModel)
}
