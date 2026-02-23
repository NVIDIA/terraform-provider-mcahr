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

var _ core.TFModel = &NVaultSecretTFModel{}

// NVaultSecretTFModel represents the Terraform model for nvault secret resources
type NVaultSecretTFModel struct {
	Name            types.String `tfsdk:"name"`
	VaultSecretPath types.String `tfsdk:"vault_secret_path"`
	VaultSecretKey  types.String `tfsdk:"vault_secret_key"`
	IntegrationName types.String `tfsdk:"integration_name"`
}

// GetName returns the name of the nvault secret resource
func (s *NVaultSecretTFModel) GetName() string {
	return s.Name.ValueString()
}
