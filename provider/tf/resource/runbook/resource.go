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

package runbook

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	runbookapi "terraform/terraform-provider/provider/external_api/resources/runbooks"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"
	plan "terraform/terraform-provider/provider/tf/resource/runbook/plan/modifiers"
	datavalidator "terraform/terraform-provider/provider/tf/resource/runbook/plan/validators/data"
	runbookprocess "terraform/terraform-provider/provider/tf/resource/runbook/process"
	"terraform/terraform-provider/provider/tf/resource/runbook/schema"
	translator "terraform/terraform-provider/provider/tf/resource/runbook/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RunbookResource{}
var _ resource.ResourceWithConfigure = &RunbookResource{}
var _ resource.ResourceWithImportState = &RunbookResource{}
var _ resource.ResourceWithValidateConfig = &RunbookResource{}
var _ resource.ResourceWithModifyPlan = &RunbookResource{}
var _ coreresource.ConfigurableResource = &RunbookResource{}

func NewRunbookResource(typeName string) resource.Resource {
	return &RunbookResource{
		typeName: typeName,
		schema:   &schema.RunbookSchema{},
	}
}

type RunbookResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	typeName       string
	schema         *schema.RunbookSchema
}

func (r *RunbookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, r.typeName)
}

func (r *RunbookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *RunbookResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *RunbookResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *RunbookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for principal resource
func (r *RunbookResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*runbooktf.RunbookTFModel,
	*runbookapi.RunbookResponseAPIModelV1,
	*runbookapi.RunbookResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*runbooktf.RunbookTFModel,
		*runbookapi.RunbookResponseAPIModelV1,
		*runbookapi.RunbookResponseAPIModel,
	]{
		ResourceType:         r.typeName,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &runbookprocess.RunbookPreProcessor{},
		PostProcessor:        &runbookprocess.RunbookPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.RunbookTranslatorV1{},
		TranslatorV2:         &translator.RunbookTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *RunbookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *RunbookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *RunbookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *RunbookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *RunbookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, r.typeName)
}

func (r *RunbookResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	plan.ModifyPlan(ctx, req, resp, r.schema, r.backendVersion)
}

func (r *RunbookResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	datavalidator.ApplyDataValidators(ctx, req, resp)
}
