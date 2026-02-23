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

// AlarmTranslatorV1 handles translation between TF models and V1 API models for alarm resources
type AlarmTranslatorV1 struct {
	AlarmTranslatorCommon
}

var _ translator.Translator[*alarmtf.AlarmTFModel, *alarmapi.AlarmResponseAPIModelV1] = &AlarmTranslatorV1{}

// ToTFModel converts a V1 API model to a TF model
func (t *AlarmTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *alarmapi.AlarmResponseAPIModelV1) (*alarmtf.AlarmTFModel, error) {
	if apiModel == nil {
		return nil, nil
	}

	// Get the alarm container regardless of operation type (define_alarm, update_alarm, get_alarm_class, delete_alarm)
	container := apiModel.GetContainer()
	if container == nil {
		return nil, fmt.Errorf("no alarm container found in V1 API response")
	}

	if len(container.AlarmClasses) == 0 {
		return nil, fmt.Errorf("no alarm classes found in V1 API response")
	}

	// Get the first alarm class, current implementation only supports one alarm to be returned by the API
	alarmClass := container.AlarmClasses[0]

	// Build TF model from V1 alarm class and response-level fields
	tfModel := &alarmtf.AlarmTFModel{
		Name:             types.StringValue(alarmClass.Name),
		FireQuery:        types.StringValue(alarmClass.FireQuery),
		ClearQuery:       types.StringValue(alarmClass.ClearQuery),
		Description:      types.StringValue(alarmClass.Description),
		ResourceQuery:    types.StringValue(alarmClass.ResourceQuery),
		ResourceType:     types.StringValue(alarmClass.ResourceType),
		CheckIntervalSec: types.Int64Value(alarmClass.CheckIntervalSec),
		Enabled:          types.BoolValue(alarmClass.Enabled),
	}

	// Map family from config data
	tfModel.Family = types.StringValue(alarmClass.ConfigData.Family)

	// Map step classes to template fields
	tfModel.FireTitleTemplate = types.StringValue(alarmClass.FireStepClass.TitleTemplate)
	tfModel.FireShortTemplate = types.StringValue(alarmClass.FireStepClass.ShortTemplate)
	tfModel.ResolveTitleTemplate = types.StringValue(alarmClass.ClearStepClass.TitleTemplate)
	tfModel.ResolveShortTemplate = types.StringValue(alarmClass.ClearStepClass.ShortTemplate)

	return tfModel, nil
}

// ToAPIModel converts a TF model to a V1 API statement model
func (t *AlarmTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *alarmtf.AlarmTFModel) (*statement.StatementInputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
