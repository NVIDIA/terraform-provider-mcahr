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

// TimeTriggerTranslatorV1 handles translation for TimeTriggerResponseAPIModelV1
type TimeTriggerTranslatorV1 struct {
	TimeTriggerTranslatorCommon
}

var _ translator.Translator[*timetriggertf.TimeTriggerTFModel, *timetriggerapi.TimeTriggerResponseAPIModelV1] = &TimeTriggerTranslatorV1{}

func (t *TimeTriggerTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *timetriggerapi.TimeTriggerResponseAPIModelV1) (*timetriggertf.TimeTriggerTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the time trigger container regardless of operation type (define_time_trigger, update_time_trigger, get_time_trigger_class)
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no time trigger container found in V1 API response")
	}

	if len(container.TimeTriggerClasses) == 0 {
		return nil, fmt.Errorf("no time trigger classes found in V1 API response")
	}

	// Get the first time trigger class, current implementation only supports one time trigger to be returned by the API
	timeTriggerClass := container.TimeTriggerClasses[0]

	tfModel := &timetriggertf.TimeTriggerTFModel{
		Name:      types.StringValue(timeTriggerClass.Name),
		FireQuery: types.StringValue(timeTriggerClass.FireQuery),
		StartDate: types.StringValue(timeTriggerClass.StartDate),
		EndDate:   types.StringValue(timeTriggerClass.EndDate),
		Enabled:   types.BoolValue(timeTriggerClass.Enabled),
	}

	return tfModel, nil
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (t *TimeTriggerTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *timetriggertf.TimeTriggerTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
