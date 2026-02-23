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

package alarm

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	alarmapi "terraform/terraform-provider/provider/external_api/resources/alarms"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	alarmtf "terraform/terraform-provider/provider/tf/resource/alarm/model"
	alarms "terraform/terraform-provider/provider/tf/resource/alarm/process"
	alarmschema "terraform/terraform-provider/provider/tf/resource/alarm/schema"
	"terraform/terraform-provider/provider/tf/resource/alarm/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "alarm"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AlarmResource{}
var _ resource.ResourceWithConfigure = &AlarmResource{}
var _ resource.ResourceWithImportState = &AlarmResource{}
var _ coreresource.ConfigurableResource = &AlarmResource{}

func NewAlarmResource() resource.Resource {
	return &AlarmResource{
		schema: &alarmschema.AlarmSchema{},
	}
}

type AlarmResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *alarmschema.AlarmSchema
}

// SetClient implements ConfigurableResource
func (r *AlarmResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *AlarmResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *AlarmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for alarm resource
func (r *AlarmResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*alarmtf.AlarmTFModel,
	*alarmapi.AlarmResponseAPIModelV1,
	*alarmapi.AlarmResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*alarmtf.AlarmTFModel,
		*alarmapi.AlarmResponseAPIModelV1,
		*alarmapi.AlarmResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &alarms.AlarmPreProcessor{},
		PostProcessor:        &alarms.AlarmPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.AlarmTranslatorV1{},
		TranslatorV2:         &translator.AlarmTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *AlarmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *AlarmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

func (r *AlarmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *AlarmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *AlarmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *AlarmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *AlarmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *AlarmResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &alarmtf.AlarmTFModel{})
}
