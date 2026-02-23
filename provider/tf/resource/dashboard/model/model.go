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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DashboardTFModel represents the terraform configuration for a dashboard
type DashboardTFModel struct {

	// Required attributes
	Name          types.String `tfsdk:"name"`
	DashboardType types.String `tfsdk:"dashboard_type"`

	// Optional attributes
	ResourceQuery types.String `tfsdk:"resource_query"`

	// JSON attributes
	Groups     types.String `tfsdk:"groups"`
	GroupsFull types.String `tfsdk:"groups_full"`
	Values     types.String `tfsdk:"values"`
	ValuesFull types.String `tfsdk:"values_full"`

	// List attributes
	OtherTags   types.List `tfsdk:"other_tags"`
	Identifiers types.List `tfsdk:"identifiers"`
}

// GetName returns the name of the dashboard resource
func (d DashboardTFModel) GetName() string {
	return d.Name.ValueString()
}
