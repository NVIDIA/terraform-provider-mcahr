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

package principals

import apicommon "terraform/terraform-provider/provider/external_api/resources/common"

// PrincipalData represents the principal data in V2 API responses
type PrincipalData struct {
	Name                 string `json:"name"`
	Identity             string `json:"identity"`
	IDPName              string `json:"idp_name"`
	ActionLimit          int    `json:"action_limit"`
	ExecuteLimit         int    `json:"execute_limit"`
	ViewLimit            int    `json:"view_limit"` // removed in release-29.0.0
	ConfigurePermission  int    `json:"configure_permission"`
	AdministerPermission int    `json:"administer_permission"`
}

// AccessControlItem represents a single access control item in the V2 API response
type AccessControlItem struct {
	Data PrincipalData `json:"data"`
}

// AccessControl represents the access_control section of the V2 API response
type AccessControl struct {
	Items []AccessControlItem `json:"items"`
}

// PrincipalOutput represents the output section of the V2 API response
type PrincipalOutput struct {
	AccessControl AccessControl `json:"access_control"`
}

// PrincipalSummary represents the summary section of the V2 API response
type PrincipalSummary struct {
	Status string            `json:"status"`
	Errors []apicommon.Error `json:"errors"`
}

// PrincipalResponseAPIModel represents the complete API response for principal operations (V2)
type PrincipalResponseAPIModel struct {
	Output  PrincipalOutput  `json:"output"`
	Summary PrincipalSummary `json:"summary"`
}

// GetErrors returns a string representation of the API response status and any errors
func (p PrincipalResponseAPIModel) GetErrors() string {
	return apicommon.FormatErrors(p.Summary.Status, p.Summary.Errors)
}
