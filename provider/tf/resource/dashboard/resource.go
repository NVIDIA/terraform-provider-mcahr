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

package dashboard

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	dashboardapi "terraform/terraform-provider/provider/external_api/resources/dashboards"
	compatibility "terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	dashboardtf "terraform/terraform-provider/provider/tf/resource/dashboard/model"
	dashboardplan "terraform/terraform-provider/provider/tf/resource/dashboard/plan/modifiers"
	dashboardprocess "terraform/terraform-provider/provider/tf/resource/dashboard/process"
	dashboardschema "terraform/terraform-provider/provider/tf/resource/dashboard/schema"
	"terraform/terraform-provider/provider/tf/resource/dashboard/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "dashboard"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &DashboardResource{}
	_ resource.ResourceWithConfigure    = &DashboardResource{}
	_ resource.ResourceWithImportState  = &DashboardResource{}
	_ resource.ResourceWithModifyPlan   = &DashboardResource{}
	_ coreresource.ConfigurableResource = &DashboardResource{}
)

// NewDashboardResource creates a new dashboard resource.
func NewDashboardResource() resource.Resource {
	return &DashboardResource{
		schema: &dashboardschema.DashboardSchema{},
	}
}

// DashboardResource defines the resource implementation.
type DashboardResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *dashboardschema.DashboardSchema
}

func (r *DashboardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *DashboardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *DashboardResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *DashboardResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *DashboardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for dashboard resource
func (r *DashboardResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*dashboardtf.DashboardTFModel,
	*dashboardapi.DashboardResponseAPIModelV1,
	*dashboardapi.DashboardResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*dashboardtf.DashboardTFModel,
		*dashboardapi.DashboardResponseAPIModelV1,
		*dashboardapi.DashboardResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &dashboardprocess.DashboardPreProcessor{},
		PostProcessor:        &dashboardprocess.DashboardPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.DashboardTranslatorV1{},
		TranslatorV2:         &translator.DashboardTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *DashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *DashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *DashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *DashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *DashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *DashboardResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	dashboardplan.ModifyPlan(ctx, req, resp, r.backendVersion)

	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &dashboardtf.DashboardTFModel{})
}
