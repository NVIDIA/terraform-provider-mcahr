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

package system_settings

import (
	"context"

	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	systemsettingsapi "terraform/terraform-provider/provider/external_api/resources/system_settings"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	systemsettingstf "terraform/terraform-provider/provider/tf/resource/system_settings/model"
	systemsettingsprocess "terraform/terraform-provider/provider/tf/resource/system_settings/process"
	schema "terraform/terraform-provider/provider/tf/resource/system_settings/schema"
	"terraform/terraform-provider/provider/tf/resource/system_settings/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "system_settings"

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                 = &SystemSettingsResource{}
	_ resource.ResourceWithConfigure    = &SystemSettingsResource{}
	_ resource.ResourceWithImportState  = &SystemSettingsResource{}
	_ coreresource.ConfigurableResource = &SystemSettingsResource{}
)

type SystemSettingsResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *schema.SystemSettingsSchema
}

// NewSystemSettingsResource creates a new system_settings resource
func NewSystemSettingsResource() resource.Resource {
	return &SystemSettingsResource{
		schema: &schema.SystemSettingsSchema{},
	}
}

// Metadata returns the resource type name
func (r *SystemSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

// Schema returns the schema for the system_settings resource
func (r *SystemSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema.GetSchema()
}

// SetClient implements ConfigurableResource
func (r *SystemSettingsResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *SystemSettingsResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

// Configure prepares the resource for CRUD operations
func (r *SystemSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

// getCRUDParams returns the common CRUD operation parameters for system_settings resource
func (r *SystemSettingsResource) getCRUDParams() *coreresource.CRUDOperationParams[*systemsettingstf.SystemSettingsTFModel, *systemsettingsapi.SystemSettingsResponseAPIModelV1, *systemsettingsapi.SystemSettingsResponseAPIModel] {
	return &coreresource.CRUDOperationParams[*systemsettingstf.SystemSettingsTFModel, *systemsettingsapi.SystemSettingsResponseAPIModelV1, *systemsettingsapi.SystemSettingsResponseAPIModel]{
		ResourceType:         resourceType,
		BackendVersion:       r.backendVersion,
		Client:               r.client,
		PreProcessor:         &systemsettingsprocess.SystemSettingsPreProcessor{},
		PostProcessor:        &systemsettingsprocess.SystemSettingsPostProcessor{},
		Schema:               r.schema,
		TranslatorV1:         &translator.SystemSettingsTranslatorV1{},
		TranslatorV2:         &translator.SystemSettingsTranslator{},
		CompatibilityOptions: r.schema.GetCompatibilityOptions(),
	}
}

func (r *SystemSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	coreresource.ExecuteCreate(ctx, req, resp, r.getCRUDParams())
}

func (r *SystemSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	coreresource.ExecuteRead(ctx, req, resp, r.getCRUDParams())
}

func (r *SystemSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	coreresource.ExecuteUpdate(ctx, req, resp, r.getCRUDParams())
}

func (r *SystemSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// For system_settings, delete is a no-op. We return without error and without making any API calls.
	// Terraform will automatically remove the resource from state, but the actual system configuration
	// remains unchanged on the server (which is the desired behavior for global system settings).

}

func (r *SystemSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *SystemSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &systemsettingstf.SystemSettingsTFModel{})
}
