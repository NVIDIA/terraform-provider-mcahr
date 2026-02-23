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
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/common/log"
	"terraform/terraform-provider/provider/common/systemdefer"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	api "terraform/terraform-provider/provider/external_api/resources"
	"terraform/terraform-provider/provider/tf/core/config"
	model "terraform/terraform-provider/provider/tf/core/model"
	coreorchestrator "terraform/terraform-provider/provider/tf/core/orchestrator"
	"terraform/terraform-provider/provider/tf/core/process"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/translator"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const defaultLogLevel = hclog.Info

// SetMetadata is a common utility method for setting resource metadata
func SetMetadata(req resource.MetadataRequest, resp *resource.MetadataResponse, resourceType string) {
	resp.TypeName = req.ProviderTypeName + "_" + resourceType
}

// Configure is a common utility method for configuring framework resources with version support
func Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, r ConfigurableResource) {
	// Prevent panic if the provider data is nil
	// This can happen if the provider is not configured yet
	// and the resource Configure is called multiple times
	if req.ProviderData == nil {
		return
	}

	providerData := config.ReadConfiguredProviderData(req, resp)
	if providerData != nil {
		r.SetClient(providerData.Client)
		r.SetBackendVersion(providerData.BackendVersion)
	}
}

// CRUDOperationParams contains parameters for CRUD operations
type CRUDOperationParams[TF model.TFModel, API_V1 api.APIModel, API_V2 api.APIModel] struct {
	ResourceType         string
	BackendVersion       *version.BackendVersion
	Client               *client.PlatformClient
	PreProcessor         process.PreProcessor[TF]
	PostProcessor        process.PostProcessor[TF]
	Schema               coreschema.ResourceSchema
	TranslatorV1         translator.Translator[TF, API_V1]
	TranslatorV2         translator.Translator[TF, API_V2]
	CompatibilityOptions map[string]attribute.CompatibilityOptions
	DeferFunctionList    *systemdefer.DeferFunctionList
	StringArgs           map[string]string
}

// ExecuteCreate performs a common create operation for resources
func ExecuteCreate[TF model.TFModel, API_V1 api.APIModel, API_V2 api.APIModel](ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, params *CRUDOperationParams[TF, API_V1, API_V2]) {

	// Extract resource name from the config
	resourceName := extractResourceName[TF](ctx, req.Config)

	// Create subsystem context on-demand for this resource type with standard persistent fields
	subsystemCtx := createSubsystemContext(ctx, params, common.Create, resourceName)

	// Create request context with operation and resource metadata using the subsystem context
	requestCtx := common.NewRequestContext(subsystemCtx).
		WithOperation(common.Create).
		WithResourceType(params.ResourceType).
		WithBackendVersion(params.BackendVersion)

	if !config.EnsureClientConfigured(params.Client, &resp.Diagnostics) {
		log.LogDebug(requestCtx, "Client not configured, skipping create operation", nil)
		return
	}

	// Use consistent logging utility
	log.LogInfo(requestCtx, fmt.Sprintf("Starting %s create operation", params.ResourceType), nil)

	// Use version-aware orchestration to support both V1 and V2 APIs
	resultData, err := coreorchestrator.OrchestrateByAPIVersion(
		requestCtx,
		params.Client,
		params.Schema,
		params.PreProcessor,
		params.PostProcessor,
		&process.ProcessData{
			CreateRequest:     &req,
			CreateResponse:    resp,
			Client:            params.Client,
			DeferFunctionList: params.DeferFunctionList,
			StringArgs:        params.StringArgs,
		},
		params.TranslatorV1, // V1 translator
		params.TranslatorV2, // V2 translator
		&translator.TranslationData{
			CompatibilityOptions: params.CompatibilityOptions,
		},
	)

	if common.HasErrorOrNil(err, resultData) {
		errorMessage := getErrorMessage(err, "create")

		log.LogError(requestCtx, fmt.Sprintf("%s create operation failed", params.ResourceType), map[string]any{
			"error": errorMessage,
		})
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Failed to create %s: %s", params.ResourceType, errorMessage))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resultData)...)

	log.LogInfo(requestCtx, fmt.Sprintf("%s create operation completed successfully", params.ResourceType), map[string]any{
		fmt.Sprintf("%s_name", params.ResourceType): resultData.GetName(),
	})
}

// ExecuteRead performs a common read operation for resources
func ExecuteRead[TF model.TFModel, API_V1 api.APIModel, API_V2 api.APIModel](ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, params *CRUDOperationParams[TF, API_V1, API_V2]) {

	// Extract resource name from the state
	resourceName := extractResourceName[TF](ctx, req.State)

	// Create subsystem context on-demand for this resource type with standard persistent fields
	subsystemCtx := createSubsystemContext(ctx, params, common.Read, resourceName)

	// Create request context with operation and resource metadata using the subsystem context
	requestCtx := common.NewRequestContext(subsystemCtx).
		WithOperation(common.Read).
		WithResourceType(params.ResourceType).
		WithBackendVersion(params.BackendVersion)

	if !config.EnsureClientConfigured(params.Client, &resp.Diagnostics) {
		log.LogDebug(requestCtx, "Client not configured, skipping read operation", nil)
		return
	}

	log.LogInfo(requestCtx, fmt.Sprintf("Starting %s read operation", params.ResourceType), nil)

	// Use version-aware orchestration to support both V1 and V2 APIs
	resultData, err := coreorchestrator.OrchestrateByAPIVersion(
		requestCtx,
		params.Client,
		params.Schema,
		params.PreProcessor,
		params.PostProcessor,
		&process.ProcessData{
			ReadRequest:       &req,
			ReadResponse:      resp,
			Client:            params.Client,
			DeferFunctionList: params.DeferFunctionList,
			StringArgs:        params.StringArgs,
		},
		params.TranslatorV1, // V1 translator
		params.TranslatorV2, // V2 translator
		&translator.TranslationData{
			CompatibilityOptions: params.CompatibilityOptions,
		},
	)

	if common.HasErrorOrNil(err, resultData) {
		errorMessage := getErrorMessage(err, "read")

		log.LogError(requestCtx, fmt.Sprintf("%s read operation failed", params.ResourceType), map[string]any{
			"error": errorMessage,
		})
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Failed to read %s: %s", params.ResourceType, errorMessage))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resultData)...)

	log.LogInfo(requestCtx, fmt.Sprintf("%s read operation completed successfully", params.ResourceType), map[string]any{
		fmt.Sprintf("%s_name", params.ResourceType): resultData.GetName(),
	})
}

// ExecuteUpdate performs a common update operation for resources
func ExecuteUpdate[TF model.TFModel, API_V1 api.APIModel, API_V2 api.APIModel](ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, params *CRUDOperationParams[TF, API_V1, API_V2]) {

	// Extract resource name from the plan
	resourceName := extractResourceName[TF](ctx, req.Plan)

	// Create subsystem context on-demand for this resource type with standard persistent fields
	subsystemCtx := createSubsystemContext(ctx, params, common.Update, resourceName)

	// Create request context with operation and resource metadata using the subsystem context
	requestCtx := common.NewRequestContext(subsystemCtx).
		WithOperation(common.Update).
		WithResourceType(params.ResourceType).
		WithBackendVersion(params.BackendVersion)

	if !config.EnsureClientConfigured(params.Client, &resp.Diagnostics) {
		log.LogDebug(requestCtx, "Client not configured, skipping update operation", nil)
		return
	}

	log.LogInfo(requestCtx, fmt.Sprintf("Starting %s update operation", params.ResourceType), nil)

	// Use version-aware orchestration to support both V1 and V2 APIs
	resultData, err := coreorchestrator.OrchestrateByAPIVersion(
		requestCtx,
		params.Client,
		params.Schema,
		params.PreProcessor,
		params.PostProcessor,
		&process.ProcessData{
			UpdateRequest:     &req,
			UpdateResponse:    resp,
			Client:            params.Client,
			DeferFunctionList: params.DeferFunctionList,
			StringArgs:        params.StringArgs,
		},
		params.TranslatorV1, // V1 translator
		params.TranslatorV2, // V2 translator
		&translator.TranslationData{
			CompatibilityOptions: params.CompatibilityOptions,
		},
	)

	if common.HasErrorOrNil(err, resultData) {
		errorMessage := getErrorMessage(err, "update")

		log.LogError(requestCtx, fmt.Sprintf("%s update operation failed", params.ResourceType), map[string]any{
			"error": errorMessage,
		})
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Failed to update %s: %s", params.ResourceType, errorMessage))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resultData)...)

	log.LogInfo(requestCtx, fmt.Sprintf("%s update operation completed successfully", params.ResourceType), map[string]any{
		fmt.Sprintf("%s_name", params.ResourceType): resultData.GetName(),
	})
}

// ExecuteDelete performs a common delete operation for resources
func ExecuteDelete[TF model.TFModel, API_V1 api.APIModel, API_V2 api.APIModel](ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, params *CRUDOperationParams[TF, API_V1, API_V2]) {

	// Extract resource name from the state
	resourceName := extractResourceName[TF](ctx, req.State)

	// Create subsystem context on-demand for this resource type with standard persistent fields
	subsystemCtx := createSubsystemContext(ctx, params, common.Delete, resourceName)

	// Create request context with operation and resource metadata using the subsystem context
	requestCtx := common.NewRequestContext(subsystemCtx).
		WithOperation(common.Delete).
		WithResourceType(params.ResourceType).
		WithBackendVersion(params.BackendVersion)

	if !config.EnsureClientConfigured(params.Client, &resp.Diagnostics) {
		log.LogDebug(requestCtx, "Client not configured, skipping delete operation", nil)
		return
	}

	log.LogInfo(requestCtx, fmt.Sprintf("Starting %s delete operation", params.ResourceType), nil)

	// Use version-aware orchestration to support both V1 and V2 APIs
	resultData, err := coreorchestrator.OrchestrateByAPIVersion(
		requestCtx,
		params.Client,
		params.Schema,
		params.PreProcessor,
		params.PostProcessor,
		&process.ProcessData{
			DeleteRequest:     &req,
			DeleteResponse:    resp,
			Client:            params.Client,
			DeferFunctionList: params.DeferFunctionList,
			StringArgs:        params.StringArgs,
		},
		params.TranslatorV1, // V1 translator
		params.TranslatorV2, // V2 translator
		&translator.TranslationData{
			CompatibilityOptions: params.CompatibilityOptions,
		},
	)

	if common.HasErrorOrNil(err, resultData) {
		errorMessage := getErrorMessage(err, "delete")

		log.LogError(requestCtx, fmt.Sprintf("%s delete operation failed", params.ResourceType), map[string]any{
			"error": errorMessage,
		})
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Failed to delete %s: %s", params.ResourceType, errorMessage))
		return
	}

	log.LogInfo(requestCtx, fmt.Sprintf("%s delete operation completed successfully", params.ResourceType), map[string]any{
		fmt.Sprintf("%s_name", params.ResourceType): resultData.GetName(),
	})
}

// ExecuteImportState performs a common import state operation for resources
func ExecuteImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, resourceType string) {

	// Create subsystem context on-demand for this resource type with standard persistent fields
	// For import, we use the import ID as the resource name since that's what we're importing
	subsystemCtx := createSimpleSubsystemContext(ctx, resourceType, "import", req.ID)

	// Create request context for logging
	requestCtx := common.NewRequestContext(subsystemCtx).
		WithOperation(common.Import).
		WithResourceType(resourceType)

	log.LogInfo(requestCtx, fmt.Sprintf("Starting %s import operation", resourceType), map[string]interface{}{
		"import_id": req.ID,
	})

	// Use the import ID as the resource name for reading the resource
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)

	if resp.Diagnostics.HasError() {
		log.LogError(requestCtx, fmt.Sprintf("%s import operation failed", resourceType), map[string]any{
			"import_id": req.ID,
			"error":     "Import state passthrough failed",
		})
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Failed to import %s: Import state passthrough failed", resourceType))
		return
	}

	log.LogInfo(requestCtx, fmt.Sprintf("%s import operation completed successfully", resourceType), map[string]any{
		"import_id":                          req.ID,
		fmt.Sprintf("%s_name", resourceType): req.ID,
	})
}

func getErrorMessage(err error, operation string) string {
	if err != nil {
		return err.Error()
	} else {
		return fmt.Sprintf("Result data is nil for %s operation", operation)
	}
}
