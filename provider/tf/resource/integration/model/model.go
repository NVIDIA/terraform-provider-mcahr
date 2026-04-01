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

package model

import (
	core "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ core.TFModel = &IntegrationTFModel{}

type IntegrationTFModel struct {
	// Required fields
	Name         types.String `tfsdk:"name" json:"name"`
	ServiceName  types.String `tfsdk:"service_name" json:"service_name"`
	SerialNumber types.String `tfsdk:"serial_number" json:"serial_number"`

	// Optional core fields
	Enabled         types.Bool   `tfsdk:"enabled" json:"enabled,omitempty"`
	PermissionsUser types.String `tfsdk:"permissions_user" json:"permissions_user,omitempty"`

	// Shared integration fields (used by multiple integration types)
	APIUrl  types.String `tfsdk:"api_url" json:"api_url,omitempty"`
	APIKey  types.String `tfsdk:"api_key" json:"api_key,omitempty"`
	IDPName types.String `tfsdk:"idp_name" json:"idp_name,omitempty"`

	CacheTTL     types.Int64 `tfsdk:"cache_ttl" json:"cache_ttl,omitempty"`
	CacheTTLMs   types.Int64 `tfsdk:"cache_ttl_ms" json:"cache_ttl_ms,omitempty"`
	APIRateLimit types.Int64 `tfsdk:"api_rate_limit" json:"api_rate_limit,omitempty"`

	// Alertmanager-specific fields
	ExternalUrl  types.String `tfsdk:"external_url" json:"external_url,omitempty"`
	PayloadPaths types.List   `tfsdk:"payload_paths" json:"payload_paths,omitempty"`

	// Azure Active Directory-specific fields
	TenantID     types.String `tfsdk:"tenant_id" json:"tenant_id,omitempty"`
	ClientID     types.String `tfsdk:"client_id" json:"client_id,omitempty"`
	ClientSecret types.String `tfsdk:"client_secret" json:"client_secret,omitempty"`

	// Google Cloud Identity-specific fields
	Subject     types.String `tfsdk:"subject" json:"subject,omitempty"`
	Credentials types.String `tfsdk:"credentials" json:"credentials,omitempty"`

	// BCM Connectivity-specific fields
	APICertificate types.String `tfsdk:"api_certificate" json:"api_certificate,omitempty"`

	// NVault-specific fields
	Address     types.String `tfsdk:"address" json:"address,omitempty"`
	Namespace   types.String `tfsdk:"namespace" json:"namespace,omitempty"`
	RoleName    types.String `tfsdk:"role_name" json:"role_name,omitempty"`
	JWTAuthPath types.String `tfsdk:"jwt_auth_path" json:"jwt_auth_path,omitempty"`
}

// GetName returns the integration name
func (i *IntegrationTFModel) GetName() string {
	return i.Name.ValueString()
}
