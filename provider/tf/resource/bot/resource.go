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

package bot

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	botapi "terraform/terraform-provider/provider/external_api/resources/bots"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	bottf "terraform/terraform-provider/provider/tf/resource/bot/model"
	botprocess "terraform/terraform-provider/provider/tf/resource/bot/process"
	botschema "terraform/terraform-provider/provider/tf/resource/bot/schema"
	"terraform/terraform-provider/provider/tf/resource/bot/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "bot"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &BotResource{}
	_ resource.ResourceWithConfigure    = &BotResource{}
	_ resource.ResourceWithImportState  = &BotResource{}
	_ coreresource.ConfigurableResource = &BotResource{}
)

// NewBotResource creates a new bot resource.
func NewBotResource() resource.Resource {
	return &BotResource{
		schema: &botschema.BotSchema{},
	}
}

// BotResource defines the resource implementation.
type BotResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *botschema.BotSchema
}

func (r *BotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *BotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *BotResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *BotResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *BotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for bot resource
func (r *BotResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*bottf.BotTFModel,
	*botapi.BotResponseAPIModelV1,
	*botapi.BotResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*bottf.BotTFModel,
		*botapi.BotResponseAPIModelV1,
		*botapi.BotResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &botprocess.BotPreProcessor{},
		PostProcessor:        &botprocess.BotPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.BotTranslatorV1{},
		TranslatorV2:         &translator.BotTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *BotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *BotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *BotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *BotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *BotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *BotResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &bottf.BotTFModel{})
}
