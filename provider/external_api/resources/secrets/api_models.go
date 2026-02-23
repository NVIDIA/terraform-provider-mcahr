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

package secrets

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// NVaultSecretConfig represents the configuration settings for a nvault secret
type NVaultSecretConfig struct {
	ExternalValue NVaultSecretExternalValue `json:"external_value"`
}

// NVaultSecretExternalValue represents the external value configuration for a nvault secret
type NVaultSecretExternalValue struct {
	IntegrationName string `json:"integration_name"`
	VaultSecretPath string `json:"vault_secret_path"`
	VaultSecretKey  string `json:"vault_secret_key"`
}

// NVaultSecretEntityMetadata represents the metadata for a nvault secret entity
type NVaultSecretEntityMetadata struct {
	Name string `json:"name"`
}

// NVaultConfigurationItem represents a single configuration item in the API response
type NVaultConfigurationItem struct {
	Config         NVaultSecretConfig         `json:"config"`
	EntityMetadata NVaultSecretEntityMetadata `json:"entity_metadata"`
}

// NVaultSecretConfigurations represents the configurations section of the API response
type NVaultSecretConfigurations struct {
	Items []NVaultConfigurationItem `json:"items"`
}

// NVaultSecretOutput represents the output section of the API response
type NVaultSecretOutput struct {
	Configurations NVaultSecretConfigurations `json:"configurations"`
}

// NVaultSecretSummary represents the summary section of the API response
type NVaultSecretSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// NVaultSecretResponseAPIModel represents the complete API response for nvault secret operations (V2)
type NVaultSecretResponseAPIModel struct {
	Output  NVaultSecretOutput  `json:"output"`
	Summary NVaultSecretSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (s NVaultSecretResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(s.Summary.Status, s.Summary.Errors)
}
