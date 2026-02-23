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

package report_template

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	reporttemplateapi "terraform/terraform-provider/provider/external_api/resources/report_templates"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	reporttemplatetf "terraform/terraform-provider/provider/tf/resource/report_template/model"
	reporttemplateplan "terraform/terraform-provider/provider/tf/resource/report_template/plan/modifiers"
	reporttemplateprocess "terraform/terraform-provider/provider/tf/resource/report_template/process"
	reporttemplateschema "terraform/terraform-provider/provider/tf/resource/report_template/schema"
	"terraform/terraform-provider/provider/tf/resource/report_template/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "report_template"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &ReportTemplateResource{}
	_ resource.ResourceWithConfigure    = &ReportTemplateResource{}
	_ resource.ResourceWithImportState  = &ReportTemplateResource{}
	_ resource.ResourceWithModifyPlan   = &ReportTemplateResource{}
	_ coreresource.ConfigurableResource = &ReportTemplateResource{}
)

// NewReportTemplateResource creates a new report template resource.
func NewReportTemplateResource() resource.Resource {
	return &ReportTemplateResource{
		schema: &reporttemplateschema.ReportTemplateSchema{},
	}
}

// ReportTemplateResource defines the resource implementation.
type ReportTemplateResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *reporttemplateschema.ReportTemplateSchema
}

func (r *ReportTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *ReportTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *ReportTemplateResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *ReportTemplateResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *ReportTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for report template resource
func (r *ReportTemplateResource) getCRUDParams() *coreresource.CRUDOperationParams[*reporttemplatetf.ReportTemplateTFModel, *reporttemplateapi.ReportTemplateResponseAPIModelV1, *reporttemplateapi.ReportTemplateResponseAPIModel] {
	return &coreresource.CRUDOperationParams[*reporttemplatetf.ReportTemplateTFModel, *reporttemplateapi.ReportTemplateResponseAPIModelV1, *reporttemplateapi.ReportTemplateResponseAPIModel]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &reporttemplateprocess.ReportTemplatePreProcessor{},
		PostProcessor:        &reporttemplateprocess.ReportTemplatePostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.ReportTemplateTranslatorV1{},
		TranslatorV2:         &translator.ReportTemplateTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *ReportTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *ReportTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *ReportTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *ReportTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *ReportTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *ReportTemplateResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	reporttemplateplan.ModifyPlan(ctx, req, resp, r.backendVersion)

	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &reporttemplatetf.ReportTemplateTFModel{})
}
