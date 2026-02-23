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

package secret

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	secretapi "terraform/terraform-provider/provider/external_api/resources/secrets"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	secrettf "terraform/terraform-provider/provider/tf/resource/secret/model"
	"terraform/terraform-provider/provider/tf/resource/secret/schema"

	"terraform/terraform-provider/provider/tf/resource/secret/process"
	"terraform/terraform-provider/provider/tf/resource/secret/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "nvault_secret"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &NVaultSecretResource{}
	_ resource.ResourceWithConfigure    = &NVaultSecretResource{}
	_ resource.ResourceWithImportState  = &NVaultSecretResource{}
	_ coreresource.ConfigurableResource = &NVaultSecretResource{}
)

// NewNVaultSecretResource creates a new nvault secret resource.
func NewNVaultSecretResource() resource.Resource {
	return &NVaultSecretResource{
		schema: &schema.NVaultSecretSchema{},
	}
}

// NVaultSecretResource defines the resource implementation.
type NVaultSecretResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *schema.NVaultSecretSchema
}

func (r *NVaultSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *NVaultSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *NVaultSecretResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *NVaultSecretResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *NVaultSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for nvault secret resource
func (r *NVaultSecretResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*secrettf.NVaultSecretTFModel,
	*secretapi.NVaultSecretResponseAPIModelV1,
	*secretapi.NVaultSecretResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*secrettf.NVaultSecretTFModel,
		*secretapi.NVaultSecretResponseAPIModelV1,
		*secretapi.NVaultSecretResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		Schema:               r.schema,
		PreProcessor:         &process.NVaultSecretPreProcessor{},
		PostProcessor:        &process.NVaultSecretPostProcessor{},
		TranslatorV1:         &translator.NVaultSecretTranslatorV1{},
		TranslatorV2:         &translator.NVaultSecretTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *NVaultSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *NVaultSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *NVaultSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *NVaultSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *NVaultSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *NVaultSecretResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &secrettf.NVaultSecretTFModel{})
}
