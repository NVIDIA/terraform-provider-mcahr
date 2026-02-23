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
	"terraform/terraform-provider/provider/external_api/resources/statement"
	timetriggerapi "terraform/terraform-provider/provider/external_api/resources/time_triggers"
	"terraform/terraform-provider/provider/tf/core/translator"
	timetriggertf "terraform/terraform-provider/provider/tf/resource/time_trigger/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TimeTriggerTranslator handles translation for TimeTriggerResponseAPIModel (V2)
type TimeTriggerTranslator struct {
	TimeTriggerTranslatorCommon
}

var _ translator.Translator[*timetriggertf.TimeTriggerTFModel, *timetriggerapi.TimeTriggerResponseAPIModel] = &TimeTriggerTranslator{}

func (t *TimeTriggerTranslator) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *timetriggerapi.TimeTriggerResponseAPIModel) (*timetriggertf.TimeTriggerTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	if len(apiModel.Output.Configurations.Items) == 0 {
		return nil, fmt.Errorf("no time trigger configurations found in API response")
	}

	// Get the first configuration item (current implementation only supports one time trigger to be returned by the API)
	configItem := apiModel.Output.Configurations.Items[0]
	config := configItem.Config
	metadata := configItem.EntityMetadata

	tfModel := &timetriggertf.TimeTriggerTFModel{
		Name:      types.StringValue(metadata.Name),
		FireQuery: types.StringValue(config.FireQuery),
		StartDate: types.StringValue(config.StartDate),
		EndDate:   types.StringValue(config.EndDate),
		Enabled:   types.BoolValue(metadata.Enabled),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V2 backend
func (t *TimeTriggerTranslator) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *timetriggertf.TimeTriggerTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
