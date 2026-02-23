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

var _ core.TFModel = &PrincipalTFModel{}

// PrincipalTFModel represents the Terraform model for principal resources
type PrincipalTFModel struct {
	Name                 types.String `tfsdk:"name"`
	Identity             types.String `tfsdk:"identity"`
	IDPName              types.String `tfsdk:"idp_name"`
	ActionLimit          types.Int64  `tfsdk:"action_limit"`
	ExecuteLimit         types.Int64  `tfsdk:"execute_limit"`
	AdministerPermission types.Bool   `tfsdk:"administer_permission"`
	ConfigurePermission  types.Bool   `tfsdk:"configure_permission"`
	ViewLimit            types.Int64  `tfsdk:"view_limit"`
}

// GetName returns the name of the principal resource
func (p *PrincipalTFModel) GetName() string {
	return p.Name.ValueString()
}
