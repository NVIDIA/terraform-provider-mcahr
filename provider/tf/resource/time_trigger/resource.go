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

package time_trigger

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	timetriggerapi "terraform/terraform-provider/provider/external_api/resources/time_triggers"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	timetriggertf "terraform/terraform-provider/provider/tf/resource/time_trigger/model"
	timetriggers "terraform/terraform-provider/provider/tf/resource/time_trigger/process"
	"terraform/terraform-provider/provider/tf/resource/time_trigger/schema"
	"terraform/terraform-provider/provider/tf/resource/time_trigger/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "time_trigger"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TimeTriggerResource{}
var _ resource.ResourceWithConfigure = &TimeTriggerResource{}
var _ resource.ResourceWithImportState = &TimeTriggerResource{}
var _ coreresource.ConfigurableResource = &TimeTriggerResource{}

func NewTimeTriggerResource() resource.Resource {
	return &TimeTriggerResource{
		schema: &schema.TimeTriggerSchema{},
	}
}

// TimeTriggerResource defines the resource implementation.
type TimeTriggerResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *schema.TimeTriggerSchema
}

// SetClient implements ConfigurableResource
func (r *TimeTriggerResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *TimeTriggerResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

func (r *TimeTriggerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for time_trigger resource
func (r *TimeTriggerResource) getCRUDParams() *coreresource.CRUDOperationParams[
	*timetriggertf.TimeTriggerTFModel,
	*timetriggerapi.TimeTriggerResponseAPIModelV1,
	*timetriggerapi.TimeTriggerResponseAPIModel,
] {
	return &coreresource.CRUDOperationParams[
		*timetriggertf.TimeTriggerTFModel,
		*timetriggerapi.TimeTriggerResponseAPIModelV1,
		*timetriggerapi.TimeTriggerResponseAPIModel,
	]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		Schema:               r.schema,
		PreProcessor:         &timetriggers.TimeTriggerPreProcessor{},
		PostProcessor:        &timetriggers.TimeTriggerPostProcessor{},
		TranslatorV1:         &translator.TimeTriggerTranslatorV1{},
		TranslatorV2:         &translator.TimeTriggerTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *TimeTriggerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *TimeTriggerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

func (r *TimeTriggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *TimeTriggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *TimeTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *TimeTriggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	coreresource.ExecuteDelete(ctx, req, resp, r.getCRUDParams())
}

func (r *TimeTriggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *TimeTriggerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &timetriggertf.TimeTriggerTFModel{})
}
