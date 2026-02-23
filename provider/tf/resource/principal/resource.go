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

package principal

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	principalapi "terraform/terraform-provider/provider/external_api/resources/principals"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	principaltf "terraform/terraform-provider/provider/tf/resource/principal/model"
	principalprocess "terraform/terraform-provider/provider/tf/resource/principal/process"
	"terraform/terraform-provider/provider/tf/resource/principal/schema"
	"terraform/terraform-provider/provider/tf/resource/principal/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "principal"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &PrincipalResource{}
	_ resource.ResourceWithConfigure    = &PrincipalResource{}
	_ resource.ResourceWithImportState  = &PrincipalResource{}
	_ coreresource.ConfigurableResource = &PrincipalResource{}
)

func NewPrincipalResource() resource.Resource {
	return &PrincipalResource{
		schema: &schema.PrincipalSchema{},
	}
}

// PrincipalResource defines the resource implementation.
type PrincipalResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *schema.PrincipalSchema
}

func (r *PrincipalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *PrincipalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *PrincipalResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *PrincipalResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *PrincipalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for principal resource
func (r *PrincipalResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*principaltf.PrincipalTFModel,
	*principalapi.PrincipalResponseAPIModelV1,
	*principalapi.PrincipalResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*principaltf.PrincipalTFModel,
		*principalapi.PrincipalResponseAPIModelV1,
		*principalapi.PrincipalResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &principalprocess.PrincipalPreProcessor{},
		PostProcessor:        &principalprocess.PrincipalPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.PrincipalTranslatorV1{},
		TranslatorV2:         &translator.PrincipalTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *PrincipalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *PrincipalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *PrincipalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *PrincipalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *PrincipalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *PrincipalResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &principaltf.PrincipalTFModel{})
}
