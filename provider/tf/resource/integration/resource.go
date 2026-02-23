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

package integration

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	integrationapi "terraform/terraform-provider/provider/external_api/resources/integrations"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"
	integrationprocess "terraform/terraform-provider/provider/tf/resource/integration/process"
	integrationschema "terraform/terraform-provider/provider/tf/resource/integration/schema"
	"terraform/terraform-provider/provider/tf/resource/integration/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "integration"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IntegrationResource{}
var _ resource.ResourceWithConfigure = &IntegrationResource{}
var _ resource.ResourceWithImportState = &IntegrationResource{}
var _ coreresource.ConfigurableResource = &IntegrationResource{}

func NewIntegrationResource() resource.Resource {
	return &IntegrationResource{
		schema: &integrationschema.IntegrationSchema{},
	}
}

type IntegrationResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *integrationschema.IntegrationSchema
}

func (r *IntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *IntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *IntegrationResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *IntegrationResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *IntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for principal resource
func (r *IntegrationResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*integrationtf.IntegrationTFModel,
	*integrationapi.IntegrationResponseAPIModelV1,
	*integrationapi.IntegrationResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*integrationtf.IntegrationTFModel,
		*integrationapi.IntegrationResponseAPIModelV1,
		*integrationapi.IntegrationResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &integrationprocess.IntegrationPreProcessor{},
		PostProcessor:        &integrationprocess.IntegrationPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.IntegrationTranslatorV1{},
		TranslatorV2:         &translator.IntegrationTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *IntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *IntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *IntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *IntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *IntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *IntegrationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &integrationtf.IntegrationTFModel{})
}
