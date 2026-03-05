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

package orchestrator

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	externalapi "terraform/terraform-provider/provider/external_api"
	"terraform/terraform-provider/provider/external_api/client"
	api "terraform/terraform-provider/provider/external_api/resources"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	model "terraform/terraform-provider/provider/tf/core/model"
	"terraform/terraform-provider/provider/tf/core/process"
	"terraform/terraform-provider/provider/tf/core/process/compatibility"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/translator"
	apiresponsediff "terraform/terraform-provider/provider/tf/core/warnings/api_response_diff"
)

// API_V2_THRESHOLD_VERSION defines the minimum version that requires API V2
// Versions greater than or equal to this threshold use V2, others use V1
var apiV2ThresholdVersion = version.NewBackendVersion("release-29.1.0")

// ExternalAPIFunc represents the external API call function signature
type ExternalAPIFunc[API api.APIModel] func(*common.RequestContext, *client.PlatformClient, *statement.StatementInputAPIModel) (API, error)

// OrchestrateByAPIVersion provides version-aware orchestration for Terraform resources
// without maintaining any state. It automatically selects the appropriate translator
// (V1 or V2) based on the API version and delegates to the core Orchestrate function.
//
// Parameters:
// • operation: The CRUD operation to perform (Create, Read, Update, Delete)
// • client: HTTP client configured for the external platform API
// • processData: Container with Terraform request/response context
// • preprocessor: Component to extract and prepare data from Terraform context
// • postprocessor: Component to handle operation-specific post-processing
// • schema: Schema for the resource
// • translatorV1: Translator implementation for V1 API
// • translatorV2: Translator implementation for V2 API
//
// Returns:
// • TF: The final Terraform model representing the resource state after the operation
// • error: Any error that occurred during the orchestration process
//
// # Usage Example:
//
//	return core.OrchestrateByAPIVersion(
//		operation,
//		r.client,
//		processData,
//		&actions.ActionPreProcessor{},
//		&actions.ActionPostProcessor{},
//		&translator.ActionTranslatorV1{},
//		&translator.ActionTranslator{},
//	)
func OrchestrateByAPIVersion[TF model.TFModel, API_V1 api.APIModel, API_V2 api.APIModel, PRE process.PreProcessor[TF], POST process.PostProcessor[TF]](
	requestContext *common.RequestContext,
	client *client.PlatformClient,
	schema coreschema.ResourceSchema,
	preprocessor PRE,
	postprocessor POST,
	processData *process.ProcessData,
	translatorV1 translator.Translator[TF, API_V1],
	translatorV2 translator.Translator[TF, API_V2],
	translationData *translator.TranslationData,
) (TF, error) {

	apiVersion := determineAPIVersion(requestContext.BackendVersion)

	if apiVersion == common.V1 {
		return Orchestrate(
			requestContext.WithAPIVersion(common.V1),
			client,
			schema,
			preprocessor,
			postprocessor,
			processData,
			translatorV1,
			translationData,
		)
	} else {
		return Orchestrate(
			requestContext.WithAPIVersion(common.V2),
			client,
			schema,
			preprocessor,
			postprocessor,
			processData,
			translatorV2,
			translationData,
		)
	}
}

// Orchestrate is a convenience function that uses the default external API implementation
func Orchestrate[TF model.TFModel, API api.APIModel, PRE process.PreProcessor[TF], POST process.PostProcessor[TF], TRANS translator.Translator[TF, API]](
	requestContext *common.RequestContext,
	client *client.PlatformClient,
	schema coreschema.ResourceSchema,
	preProcessor PRE,
	postProcessor POST,
	processData *process.ProcessData,
	trans TRANS,
	translationData *translator.TranslationData,
) (TF, error) {

	res, operationAPIWasCalled, err := orchestrateWithAPIFunction(requestContext, client, schema, preProcessor, postProcessor, processData, trans, translationData, externalapi.CallExternalAPI[API])

	if operationAPIWasCalled && err != nil {
		// If the API was called, we need to cleanup objects from the remote platform
		cleanup(requestContext, client, schema, preProcessor, postProcessor, processData, trans, translationData, externalapi.CallExternalAPI[API])
	}

	return res, err
}

// orchestrateWithAPIFunction coordinates the complete lifecycle of a CRUD operation for Terraform resources.
//
// This function provides a unified workflow for Create, Read, Update, and Delete operations
// by orchestrating the interaction between preprocessing, API communication, translation,
// and post-processing components.
//
// # Workflow:
//
// 1. **Pre-processing**: Extracts and prepares data from Terraform requests
// 2. **Translation**: Converts TF models to API models for external communication
// 3. **API Communication**: Executes HTTP requests against the external platform
// 4. **Response Translation**: Converts API responses back to TF models (except deletes)
// 5. **Post-processing**: Handles operation-specific cleanup and finalization
//
// # Generic Type Parameters:
//
// • TF: Terraform model type (e.g., ActionTFModel) representing the resource model
// • API: API model type (e.g., ActionAPIModel) representing external API model
// • PRE: PreProcessor implementation for extracting data from Terraform requests
// • POST: PostProcessor implementation for handling operation-specific logic
// • TRANS: Translator implementation for converting between TF and API models

// # Parameters:
//
// • crudOperation: The type of operation (Create, Read, Update, Delete)
// • client: HTTP client configured for the external platform API
// • processData: Container with Terraform request/response context
// • schema: Schema for the resource
// • preProcessor: Component to extract and prepare data from Terraform context
// • postProcessor: Component to handle operation-specific post-processing
// • translator: Component to convert between TF and API model formats
// • externalAPIFunc: Function to call the external API (allows dependency injection for testing)

// # Returns:
//
// • TF: The final Terraform model representing the resource state after the operation
// • error: Any error that occurred during the orchestration process
//
// # Design Benefits:
//
// • **Consistency**: All CRUD operations follow the same orchestrated workflow
// • **Modularity**: Each component has a single responsibility and can be tested independently
// • **Type Safety**: Generic constraints ensure compile-time type checking
// • **Extensibility**: New resource types can reuse the orchestration by providing implementations
// • **Error Handling**: Centralized error propagation with clear failure points
// • **Testability**: External API can be mocked via dependency injection
func orchestrateWithAPIFunction[TF model.TFModel, API api.APIModel, PRE process.PreProcessor[TF], POST process.PostProcessor[TF], TRANS translator.Translator[TF, API]](
	requestContext *common.RequestContext,
	client *client.PlatformClient,
	schema coreschema.ResourceSchema,
	preProcessor PRE,
	postProcessor POST,
	processData *process.ProcessData,
	trans TRANS,
	translationData *translator.TranslationData,
	externalAPIFunc ExternalAPIFunc[API],
) (TF, bool, error) {
	var nilTF TF
	operationAPICallWasSuccesfull := false

	// Pre-process the request
	tfObject, err := callPreProcessor(requestContext, preProcessor, processData)
	if common.HasErrorOrNil(err, tfObject) {
		return nilTF, operationAPICallWasSuccesfull, err
	}

	// Convert the Terraform model to an API model
	apiObject, err := trans.ToAPIModel(requestContext, translationData, tfObject)
	if common.HasErrorOrNil(err, apiObject) {
		return nilTF, operationAPICallWasSuccesfull, err
	}

	// Call the external API
	apiResponse, err := externalAPIFunc(requestContext, client, apiObject)
	if common.HasErrorOrNil(err, apiResponse) {
		return nilTF, operationAPICallWasSuccesfull, err
	}

	operationAPICallWasSuccesfull = true

	// Convert the API model back to a Terraform model
	tfObject, err = trans.ToTFModel(requestContext, translationData, apiResponse)
	if common.HasErrorOrNil(err, tfObject) {
		return nilTF, operationAPICallWasSuccesfull, err
	}

	// For Create/Update operations: ensure API response matches plan expectations
	err = ensurePlanConsistency(requestContext, processData, schema, tfObject)
	if err != nil {
		return nilTF, operationAPICallWasSuccesfull, err
	}

	// Post-process the response for resource-specific exceptions
	// (e.g., computed fields, deprecated field mappings, _full fields for drift detection)
	err = callPostProcessor(requestContext, postProcessor, processData, tfObject)
	if common.HasErrorOrNil(err, tfObject) {
		return nilTF, operationAPICallWasSuccesfull, err
	}

	// Post-process the response to nullify incompatible fields
	err = compatibility.PostProcess(requestContext, tfObject, schema.GetCompatibilityOptions())
	if err != nil {
		return nilTF, operationAPICallWasSuccesfull, err
	}

	return tfObject, operationAPICallWasSuccesfull, nil
}

// ensurePlanConsistency validates API response against plan and restores plan values.
// This prevents "inconsistent result after apply" errors by ensuring the tfObject matches plan expectations.
// Only runs for Create/Update operations. Read/Delete operations preserve API response for drift detection.
func ensurePlanConsistency[TF model.TFModel](requestContext *common.RequestContext, processData *process.ProcessData, schema coreschema.ResourceSchema, tfObject TF) error {
	// Only apply for Create/Update operations
	// Read/Delete operations preserve API response for drift detection
	if requestContext.Operation != common.Create && requestContext.Operation != common.Update {
		return nil
	}

	// Check for API response differences and warn user
	// This detects when the API normalized or modified user input
	err := apiresponsediff.CheckPlanVsApiResponseDelta(requestContext, processData, schema, tfObject)
	if err != nil {
		return err
	}

	// Restore all fields from plan to prevent "inconsistent result after apply" errors
	// This is the default behavior - individual PostProcessors can override specific fields if needed
	return process.RestoreAllFieldsFromPlan(requestContext, processData, tfObject)
}

func callPreProcessor[TF model.TFModel](requestContext *common.RequestContext, preProcessor process.PreProcessor[TF], processData *process.ProcessData) (TF, error) {
	switch requestContext.Operation {
	case common.Create:
		return preProcessor.PreProcessCreate(requestContext, processData)
	case common.Read:
		return preProcessor.PreProcessRead(requestContext, processData)
	case common.Update:
		return preProcessor.PreProcessUpdate(requestContext, processData)
	case common.Delete:
		return preProcessor.PreProcessDelete(requestContext, processData)
	default:
		// Return the zero value of TF pointer type (which is nil)
		var tfObject TF
		return tfObject, fmt.Errorf("unsupported CRUD operation: %v", requestContext.Operation)
	}
}

func callPostProcessor[TF model.TFModel](requestContext *common.RequestContext, postProcessor process.PostProcessor[TF], processData *process.ProcessData, tfObject TF) error {
	switch requestContext.Operation {
	case common.Create:
		return postProcessor.PostProcessCreate(requestContext, processData, tfObject)
	case common.Read:
		return postProcessor.PostProcessRead(requestContext, processData, tfObject)
	case common.Update:
		return postProcessor.PostProcessUpdate(requestContext, processData, tfObject)
	case common.Delete:
		return postProcessor.PostProcessDelete(requestContext, processData, tfObject)
	default:
		return fmt.Errorf("unsupported CRUD operation: %v", requestContext.Operation)
	}
}

// determineAPIVersion determines which API version to use based on the backend version
// If the version is greater than or equal to API_V2_THRESHOLD_VERSION, use V2, otherwise use V1
func determineAPIVersion(backendVersion *version.BackendVersion) common.APIVersion {

	// Use V2 if backend version >= threshold version, otherwise V1
	if version.CompareVersions(backendVersion, apiV2ThresholdVersion) >= 0 {
		return common.V2
	}

	return common.V1
}

func cleanup[TF model.TFModel, API api.APIModel](
	requestContext *common.RequestContext,
	client *client.PlatformClient,
	schema coreschema.ResourceSchema,
	preProcessor process.PreProcessor[TF],
	postProcessor process.PostProcessor[TF],
	processData *process.ProcessData,
	trans translator.Translator[TF, API],
	translationData *translator.TranslationData,
	externalAPIFunc ExternalAPIFunc[API],
) {

	maybeDeleteRemoteResource(
		requestContext,
		client,
		schema,
		preProcessor,
		postProcessor,
		processData,
		trans,
		translationData,
		externalAPIFunc,
	)
}

func maybeDeleteRemoteResource[TF model.TFModel, API api.APIModel](
	requestContext *common.RequestContext,
	client *client.PlatformClient,
	schema coreschema.ResourceSchema,
	preProcessor process.PreProcessor[TF],
	postProcessor process.PostProcessor[TF],
	processData *process.ProcessData,
	trans translator.Translator[TF, API],
	translationData *translator.TranslationData,
	externalAPIFunc ExternalAPIFunc[API],
) {

	if requestContext.Operation == common.Create {

		// The preprocessor will extract the resource to delete from the create request
		_, _, err := orchestrateWithAPIFunction(
			requestContext.WithOperation(common.Delete),
			client,
			schema,
			preProcessor,
			postProcessor,
			processData,
			trans,
			translationData,
			externalAPIFunc,
		)
		if err != nil {
			processData.CreateResponse.Diagnostics.AddError("Failed to cleanup remote resource after encountering an error on create. The resource was not deleted because: ", err.Error())
		}
	}
}
