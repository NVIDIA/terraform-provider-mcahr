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

package backend_version

import (
	"fmt"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/log"
	"terraform/terraform-provider/provider/common/version"
	externalapi "terraform/terraform-provider/provider/external_api"
	"terraform/terraform-provider/provider/external_api/resources/statement"
)

// BackendVersionProvider provides functionality to fetch backend version from the API
type BackendVersionProvider struct {
	client externalapi.PlatformClientInterface
}

// NewBackendVersionProvider creates a new backend version provider
func NewBackendVersionProvider(client externalapi.PlatformClientInterface) *BackendVersionProvider {
	return &BackendVersionProvider{
		client: client,
	}
}

type FetchResult struct {
	Version      *version.BackendVersion
	Error        error
	IsParseError bool
}

// FetchBackendVersion fetches the backend version from the API
func (p *BackendVersionProvider) FetchBackendVersion(requestContext *common.RequestContext) FetchResult {
	log.LogInfo(requestContext, "Fetching backend version from API", nil)

	// Create the statement to get backend version
	apiModel := &statement.StatementInputAPIModel{
		Statement:  "backend_version",
		APIVersion: common.V1, // Use V1 endpoint for backend_version
	}

	// Call the API
	response, err := externalapi.CallExternalAPIWithClient[*BackendVersionResponseAPIModelV1](requestContext, p.client, apiModel)
	if err != nil {
		log.LogError(requestContext, "Failed to fetch backend version", map[string]any{"error": err.Error()})
		return FetchResult{Error: fmt.Errorf("failed to fetch backend version: %w", err)}
	}

	backendImageTag := response.GetBackendVersion()
	if backendImageTag == "" {
		log.LogWarn(requestContext, "Backend image tag is empty in API response", nil)
		return FetchResult{Error: fmt.Errorf("backend image tag is empty in API response")}
	}

	log.LogInfo(requestContext, "Successfully fetched backend version", map[string]any{"backend_image_tag": backendImageTag})

	// Create and return the BackendVersion
	backendVersion := version.NewBackendVersion(backendImageTag)
	if backendVersion == nil {
		return FetchResult{
			Error:        fmt.Errorf("failed to parse backend version: %s", backendImageTag),
			IsParseError: true,
		}
	}

	return FetchResult{Version: backendVersion}
}

// FetchBackendVersionWithFallback fetches the backend version from API with a fallback value
func (p *BackendVersionProvider) FetchBackendVersionWithFallback(requestContext *common.RequestContext, fallbackVersion string) (*version.BackendVersion, error) {
	result := p.FetchBackendVersion(requestContext)
	if result.Error != nil {
		if result.IsParseError {
			log.LogWarn(requestContext, "Failed to parse backend version from API response, using fallback", map[string]any{
				"error":            result.Error.Error(),
				"fallback_version": fallbackVersion,
			})

			fallbackBackendVersion := version.NewBackendVersion(fallbackVersion)
			if fallbackBackendVersion == nil {
				log.LogError(requestContext, "Fallback version is also invalid", map[string]any{"fallback_version": fallbackVersion})
				return nil, fmt.Errorf("both API response and fallback version are invalid: %v", result.Error)
			}

			return fallbackBackendVersion, nil
		}

		// For other API errors, return the error
		return nil, result.Error
	}

	return result.Version, nil
}
