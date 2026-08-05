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

package smtp_subscription

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	subscriptionsapi "terraform/terraform-provider/provider/external_api/resources/smtp_subscription"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	smtp_subscriptiontf "terraform/terraform-provider/provider/tf/resource/smtp_subscription/model"
	smtp_subscriptionprocess "terraform/terraform-provider/provider/tf/resource/smtp_subscription/process"
	"terraform/terraform-provider/provider/tf/resource/smtp_subscription/schema"
	"terraform/terraform-provider/provider/tf/resource/smtp_subscription/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "smtp_subscription"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &SmtpSubscriptionResource{}
	_ resource.ResourceWithConfigure    = &SmtpSubscriptionResource{}
	_ resource.ResourceWithImportState  = &SmtpSubscriptionResource{}
	_ coreresource.ConfigurableResource = &SmtpSubscriptionResource{}
)

// NewSmtpSubscriptionResource creates a new smtp subscription resource.
func NewSmtpSubscriptionResource() resource.Resource {
	return &SmtpSubscriptionResource{
		schema: &schema.SmtpSubscriptionSchema{},
	}
}

// SmtpSubscriptionResource defines the resource implementation.
type SmtpSubscriptionResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *schema.SmtpSubscriptionSchema
}

func (r *SmtpSubscriptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *SmtpSubscriptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *SmtpSubscriptionResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *SmtpSubscriptionResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *SmtpSubscriptionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for smtp_subscription resource
func (r *SmtpSubscriptionResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*smtp_subscriptiontf.SmtpSubscriptionTFModel,
	*subscriptionsapi.SmtpSubscriptionResponseAPIModelV1,
	*subscriptionsapi.SmtpSubscriptionResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*smtp_subscriptiontf.SmtpSubscriptionTFModel,
		*subscriptionsapi.SmtpSubscriptionResponseAPIModelV1,
		*subscriptionsapi.SmtpSubscriptionResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &smtp_subscriptionprocess.SmtpSubscriptionPreProcessor{},
		PostProcessor:        &smtp_subscriptionprocess.SmtpSubscriptionPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.SmtpSubscriptionTranslatorV1{},
		TranslatorV2:         &translator.SmtpSubscriptionTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *SmtpSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *SmtpSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *SmtpSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *SmtpSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *SmtpSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType, r.schema)
}

func (r *SmtpSubscriptionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &smtp_subscriptiontf.SmtpSubscriptionTFModel{})
}
