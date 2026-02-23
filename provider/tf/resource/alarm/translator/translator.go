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
	alarmapi "terraform/terraform-provider/provider/external_api/resources/alarms"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	alarmtf "terraform/terraform-provider/provider/tf/resource/alarm/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AlarmTranslator struct {
	AlarmTranslatorCommon
}

var _ translator.Translator[*alarmtf.AlarmTFModel, *alarmapi.AlarmResponseAPIModel] = &AlarmTranslator{}

func (a *AlarmTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *alarmapi.AlarmResponseAPIModel) (*alarmtf.AlarmTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no configurations found in API response")
	}

	// Get the first configuration item, current implementation only supports one alarm to be returned by the API
	configItem := apiModel.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	tfModel := &alarmtf.AlarmTFModel{
		Name:             types.StringValue(metadata.Name),
		FireQuery:        types.StringValue(config.FireQuery),
		ClearQuery:       types.StringValue(config.ClearQuery),
		Description:      types.StringValue(metadata.Description),
		ResourceQuery:    types.StringValue(config.ResourceQuery),
		ResourceType:     types.StringValue(config.ResourceType),
		CheckIntervalSec: types.Int64Value(config.CheckIntervalSec),
		Family:           types.StringValue(metadata.Family),
		Enabled:          types.BoolValue(metadata.Enabled),

		// Map step details to template fields
		FireTitleTemplate:    types.StringValue(config.StepDetails.FireStep.Title),
		FireShortTemplate:    types.StringValue(config.StepDetails.FireStep.Description),
		ResolveTitleTemplate: types.StringValue(config.StepDetails.ClearStep.Title),
		ResolveShortTemplate: types.StringValue(config.StepDetails.ClearStep.Description),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V2 backend
func (a *AlarmTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *alarmtf.AlarmTFModel) (*statement.StatementInputAPIModel, error) {
	return a.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
