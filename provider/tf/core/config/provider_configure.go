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

package config

import (
	"fmt"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type FrameworkProviderData struct {
	Client         *client.PlatformClient
	BackendVersion *version.BackendVersion
}

func ReadConfiguredProviderData(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *FrameworkProviderData {

	providerData, ok := req.ProviderData.(*FrameworkProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *provider.FrameworkProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return providerData
}

func EnsureClientConfigured(client *client.PlatformClient, diags *diag.Diagnostics) bool {
	if client == nil {
		diags.AddError(
			"Unconfigured Client",
			"Expected configured PlatformClient. Please report this issue to the provider developers.",
		)
		return false
	}
	return true
}
