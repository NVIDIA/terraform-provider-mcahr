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
	"context"
	"encoding/json"
	"fmt"

	"terraform/terraform-provider/provider/common"
	subscriptionsapi "terraform/terraform-provider/provider/external_api/resources/smtp_subscription"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	smtp_subscriptiontf "terraform/terraform-provider/provider/tf/resource/smtp_subscription/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SmtpSubscriptionTranslatorCommon contains shared functionality between V1 and V2 translators
type SmtpSubscriptionTranslatorCommon struct{}

// ToTFModelFromResponse converts an API response to a TF model
// Shared by both V1 and V2 translators since SMTP subscriptions use the same REST response format
func (t *SmtpSubscriptionTranslatorCommon) ToTFModelFromResponse(requestContext *common.RequestContext, apiModel *subscriptionsapi.SmtpSubscriptionResponseAPIModel) (*smtp_subscriptiontf.SmtpSubscriptionTFModel, error) {

	if apiModel == nil {
		return nil, nil
	}

	ctx := requestContext.Context

	// Build recipients list from API response
	recipients, _ := types.ListValueFrom(ctx, types.StringType, apiModel.Recipients)

	// Build filters list from API response
	filters := toTFFilters(apiModel.Filters)

	tfModel := &smtp_subscriptiontf.SmtpSubscriptionTFModel{
		ID:              types.StringValue(apiModel.ID),
		Name:            types.StringValue(apiModel.Name),
		IntegrationName: types.StringValue(apiModel.IntegrationName),
		Recipients:      recipients,
		Filters:         filters,
		Enabled:         types.BoolValue(apiModel.Enabled),
	}

	return tfModel, nil
}

// toTFFilters converts API filters to TF list
func toTFFilters(apiFilters []subscriptionsapi.SmtpSubscriptionFilter) types.List {
	filterAttrTypes := map[string]attr.Type{
		"category": types.StringType,
		"type":     types.StringType,
		"status":   types.StringType,
	}
	filterObjectType := types.ObjectType{AttrTypes: filterAttrTypes}

	filterElements := make([]attr.Value, len(apiFilters))
	for i, f := range apiFilters {
		filterObj, _ := types.ObjectValue(filterAttrTypes, map[string]attr.Value{
			"category": types.StringValue(f.Category),
			"type":     types.StringValue(f.Type),
			"status":   types.StringValue(f.Status),
		})
		filterElements[i] = filterObj
	}
	filters, _ := types.ListValue(filterObjectType, filterElements)
	return filters
}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (t *SmtpSubscriptionTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *smtp_subscriptiontf.SmtpSubscriptionTFModel) (*statement.InputAPIModel, error) {
	// Build API request using the existing API model types
	apiRequest, err := toAPIRequest(requestContext.Context, tfModel)
	if err != nil {
		return nil, err
	}

	// Add ID for Update/Delete operations
	switch requestContext.Operation {
	case common.Update, common.Delete:
		apiRequest.ID = tfModel.ID.ValueString()
	}

	marshalledPayload, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, err
	}

	return &statement.InputAPIModel{
		ApiPayload: string(marshalledPayload),
		APIVersion: requestContext.APIVersion,
	}, nil
}

// toAPIRequest converts TF model to API request struct
// Using the existing API model types for automatic JSON serialization
func toAPIRequest(ctx context.Context, tfModel *smtp_subscriptiontf.SmtpSubscriptionTFModel) (*subscriptionsapi.SmtpSubscriptionUpdateRequest, error) {
	// Extract recipients as []string. A conversion failure must abort rather
	// than send an empty recipients list, which would silently stop delivery.
	recipients, err := utils.ListSliceFromTFModel(ctx, tfModel.Recipients)
	if err != nil {
		return nil, fmt.Errorf("attribute \"recipients\": %w", err)
	}

	// Extract filters and convert to API type
	filters, err := toAPIFilters(ctx, tfModel.Filters)
	if err != nil {
		return nil, err
	}

	// Use pointer for optional bool
	enabled := tfModel.Enabled.ValueBool()

	return &subscriptionsapi.SmtpSubscriptionUpdateRequest{
		ID:              tfModel.ID.ValueString(),
		Name:            tfModel.Name.ValueString(),
		IntegrationName: tfModel.IntegrationName.ValueString(),
		Recipients:      recipients,
		Filters:         filters,
		Enabled:         &enabled,
	}, nil
}

// toAPIFilters converts TF filter models to API filter structs
func toAPIFilters(ctx context.Context, filtersList types.List) ([]subscriptionsapi.SmtpSubscriptionFilter, error) {
	var filtersTF []smtp_subscriptiontf.FilterTFModel
	if !filtersList.IsNull() && !filtersList.IsUnknown() {
		if diags := filtersList.ElementsAs(ctx, &filtersTF, false); diags.HasError() {
			return nil, fmt.Errorf("attribute \"filters\": %s", diags.Errors())
		}
	}

	filters := make([]subscriptionsapi.SmtpSubscriptionFilter, len(filtersTF))
	for i, f := range filtersTF {
		filters[i] = subscriptionsapi.SmtpSubscriptionFilter{
			Category: f.Category.ValueString(),
			Type:     f.Type.ValueString(),
			Status:   f.Status.ValueString(),
		}
	}
	return filters, nil
}
