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

package action

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	actionapi "terraform/terraform-provider/provider/external_api/resources/actions"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	actiontf "terraform/terraform-provider/provider/tf/resource/action/model"
	actions "terraform/terraform-provider/provider/tf/resource/action/process"
	actionschema "terraform/terraform-provider/provider/tf/resource/action/schema"
	"terraform/terraform-provider/provider/tf/resource/action/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "action"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ActionResource{}
var _ resource.ResourceWithConfigure = &ActionResource{}
var _ resource.ResourceWithImportState = &ActionResource{}
var _ coreresource.ConfigurableResource = &ActionResource{}

func NewActionResource() resource.Resource {
	return &ActionResource{
		schema: &actionschema.ActionSchema{},
	}
}

type ActionResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *actionschema.ActionSchema
}

// SetClient implements ConfigurableResource
func (r *ActionResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *ActionResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *ActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for action resource
func (r *ActionResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*actiontf.ActionTFModel,
	*actionapi.ActionResponseAPIModelV1,
	*actionapi.ActionResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*actiontf.ActionTFModel,
		*actionapi.ActionResponseAPIModelV1,
		*actionapi.ActionResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &actions.ActionPreProcessor{},
		PostProcessor:        &actions.ActionPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.ActionTranslatorV1{},
		TranslatorV2:         &translator.ActionTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *ActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *ActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

func (r *ActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *ActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *ActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *ActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *ActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *ActionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &actiontf.ActionTFModel{})
}
