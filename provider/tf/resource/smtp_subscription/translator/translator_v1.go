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
	"terraform/terraform-provider/provider/common"
	subscriptionsapi "terraform/terraform-provider/provider/external_api/resources/smtp_subscription"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	smtp_subscriptiontf "terraform/terraform-provider/provider/tf/resource/smtp_subscription/model"
)

// SmtpSubscriptionTranslatorV1 handles translation for SmtpSubscriptionResponseAPIModelV1
type SmtpSubscriptionTranslatorV1 struct {
	SmtpSubscriptionTranslatorCommon
}

var _ translator.Translator[*smtp_subscriptiontf.SmtpSubscriptionTFModel, *subscriptionsapi.SmtpSubscriptionResponseAPIModelV1] = &SmtpSubscriptionTranslatorV1{}

func (t *SmtpSubscriptionTranslatorV1) ToTFModel(requestContext *common.RequestContext, translationData *translator.TranslationData, apiModel *subscriptionsapi.SmtpSubscriptionResponseAPIModelV1) (*smtp_subscriptiontf.SmtpSubscriptionTFModel, error) {
	// SmtpSubscriptionResponseAPIModelV1 is an alias for SmtpSubscriptionResponseAPIModel
	return t.ToTFModelFromResponse(requestContext, apiModel)
}

// ToAPIModel converts a TF model to an API model for V1 backend
func (t *SmtpSubscriptionTranslatorV1) ToAPIModel(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *smtp_subscriptiontf.SmtpSubscriptionTFModel) (*statement.InputAPIModel, error) {
	return t.ToAPIModelWithVersion(requestContext, translationData, tfModel)
}
