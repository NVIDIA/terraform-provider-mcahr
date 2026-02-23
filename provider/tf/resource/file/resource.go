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

package file

import (
	"context"

	"terraform/terraform-provider/provider/common/systemdefer"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	fileapi "terraform/terraform-provider/provider/external_api/resources/files"
	"terraform/terraform-provider/provider/tf/core/plan/modifiers/compatibility"
	coreresource "terraform/terraform-provider/provider/tf/core/resource"
	filetf "terraform/terraform-provider/provider/tf/resource/file/model"
	fileprocess "terraform/terraform-provider/provider/tf/resource/file/process"
	schema "terraform/terraform-provider/provider/tf/resource/file/schema"
	"terraform/terraform-provider/provider/tf/resource/file/translator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const resourceType = "file"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &FileResource{}
	_ resource.ResourceWithConfigure    = &FileResource{}
	_ resource.ResourceWithImportState  = &FileResource{}
	_ coreresource.ConfigurableResource = &FileResource{}
)

func NewFileResource() resource.Resource {
	return &FileResource{
		schema: &schema.FileSchema{},
	}
}

// FileResource defines the resource implementation.
type FileResource struct {
	client         *client.PlatformClient
	backendVersion *version.BackendVersion
	schema         *schema.FileSchema
}

// SetClient implements ConfigurableResource
func (r *FileResource) SetClient(client *client.PlatformClient) {
	r.client = client
}

// SetBackendVersion implements ConfigurableResource
func (r *FileResource) SetBackendVersion(version *version.BackendVersion) {
	r.backendVersion = version
}

// getCRUDParams returns the common CRUD operation parameters for file resource
func (r *FileResource) getCRUDParams() *coreresource.CRUDOperationParams[*filetf.FileTFModel, *fileapi.FileResponseAPIModelV1, *fileapi.FileResponseAPIModel] {
	return &coreresource.CRUDOperationParams[*filetf.FileTFModel, *fileapi.FileResponseAPIModelV1, *fileapi.FileResponseAPIModel]{
		DeferFunctionList: systemdefer.NewDeferFunctionList(),
		StringArgs:        map[string]string{},
		ResourceType:      resourceType,
		BackendVersion:    r.backendVersion,
		Client:            r.client,
		PreProcessor:      &fileprocess.FilePreProcessor{},
		PostProcessor:     &fileprocess.FilePostProcessor{},
		Schema:            r.schema,
		TranslatorV1:      &translator.FileTranslatorV1{},
		TranslatorV2:      &translator.FileTranslator{},
	}
}

func (r *FileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	coreresource.SetMetadata(req, resp, resourceType)
}

func (r *FileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = (&schema.FileSchema{}).GetSchema()
}

func (r *FileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	coreresource.Configure(ctx, req, resp, r)
}

func (r *FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	params := r.getCRUDParams()
	defer params.DeferFunctionList.ExecuteAll()
	coreresource.ExecuteCreate(ctx, req, resp, params)
}

func (r *FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	params := r.getCRUDParams()
	defer params.DeferFunctionList.ExecuteAll()
	coreresource.ExecuteRead(ctx, req, resp, params)
}

func (r *FileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	params := r.getCRUDParams()
	defer params.DeferFunctionList.ExecuteAll()
	coreresource.ExecuteUpdate(ctx, req, resp, params)
}

func (r *FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	params := r.getCRUDParams()
	defer params.DeferFunctionList.ExecuteAll()
	coreresource.ExecuteDelete(ctx, req, resp, params)
}

func (r *FileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	coreresource.ExecuteImportState(ctx, req, resp, resourceType)
}

func (r *FileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compatibility.ApplyResourceCompatibilityModifiers(ctx, &req, resp, r.schema, r.backendVersion, &filetf.FileTFModel{})
}
