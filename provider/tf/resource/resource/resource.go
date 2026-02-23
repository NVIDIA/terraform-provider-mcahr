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

package resource

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	resourceapi "terraform/terraform-provider/provider/external_api/resources/resources"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"

	resourcetf "terraform/terraform-provider/provider/tf/resource/resource/model"
	resourceproc "terraform/terraform-provider/provider/tf/resource/resource/process"
	resourceschema "terraform/terraform-provider/provider/tf/resource/resource/schema"
	"terraform/terraform-provider/provider/tf/resource/resource/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "resource"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithConfigure = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}
var _ coreresource.ConfigurableResource = &Resource{}

func NewResource() resource.Resource {
	return &Resource{
		schema: &resourceschema.ResourceSchema{},
	}
}

type Resource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *resourceschema.ResourceSchema
}

// SetClient implements ConfigurableResource
func (r *Resource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *Resource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for resource
func (r *Resource) getCRUDParams() *coreresource.CRUDOperationParams[
	*resourcetf.ResourceTFModel,
	*resourceapi.ResourceResponseAPIModelV1,
	*resourceapi.ResourceResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*resourcetf.ResourceTFModel,
		*resourceapi.ResourceResponseAPIModelV1,
		*resourceapi.ResourceResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &resourceproc.ResourcePreProcessor{},
		PostProcessor:        &resourceproc.ResourcePostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.ResourceTranslatorV1{},
		TranslatorV2:         &translator.ResourceTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &resourcetf.ResourceTFModel{})
}
