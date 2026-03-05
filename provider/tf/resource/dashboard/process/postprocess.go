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

package process

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	customattribute "terraform/terraform-provider/provider/external_api/resources/dashboards/custom_attribute"
	"terraform/terraform-provider/provider/tf/core/process"
	dashboardtf "terraform/terraform-provider/provider/tf/resource/dashboard/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DashboardPostProcessor struct{}

var _ process.PostProcessor[*dashboardtf.DashboardTFModel] = &DashboardPostProcessor{}

func (p *DashboardPostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *dashboardtf.DashboardTFModel) error {
	// No post-processing needed for create - orchestrator's RestoreAllFieldsFromPlan handles restoration
	return nil
}

func (p *DashboardPostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, tfModel *dashboardtf.DashboardTFModel) error {
	// Process JSON fields first to populate _full attributes from API response for drift detection
	if err := postProcessJsonFullFields(requestContext, tfModel); err != nil {
		return err
	}

	// For READ, restore base fields from state (skip _full variants for drift detection)
	// Keep _full fields from API response to enable drift detection
	return restoreBaseFieldsFromState(requestContext, data.ReadRequest.State, tfModel)
}

func (p *DashboardPostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *dashboardtf.DashboardTFModel) error {
	// No post-processing needed for update - orchestrator's RestoreAllFieldsFromPlan handles restoration
	return nil
}

func (p *DashboardPostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, tfModel *dashboardtf.DashboardTFModel) error {
	// No post-processing needed for delete operation
	return nil
}

// postProcessJsonFullFields processes all _full JSON fields with version-aware logic
func postProcessJsonFullFields(requestContext *common.RequestContext, tfModel *dashboardtf.DashboardTFModel) error {

	var err error

	tfModel.GroupsFull, err = postProcessJsonFullField[*customattribute.GroupJson](&tfModel.GroupsFull, requestContext.BackendVersion)
	if err != nil {
		return err
	}

	tfModel.ValuesFull, err = postProcessJsonFullField[*customattribute.ValueJson](&tfModel.ValuesFull, requestContext.BackendVersion)
	if err != nil {
		return err
	}

	return nil
}

func postProcessJsonFullField[T common.JsonConfigurable](fullField *types.String, backendVersion *version.BackendVersion) (types.String, error) {

	if fullField.IsNull() || fullField.IsUnknown() {
		return *fullField, nil
	}

	fullFieldString := fullField.ValueString()

	// Remarshal the full field to apply the custom struct tags (like min_version, max_version, etc.)
	// and set the default values for the fields that are not present in the JSON
	// See customattribute structs for more details
	// This is necessary to avoid any TF errors in case the backend returns values that are not supported by it (not common, but possible in case of backend bugs)
	overriddenValues, err := common.RemarshalListWithConfig[T](fullFieldString, common.JsonConfig{BackendVersion: backendVersion})
	if err != nil {
		return types.StringNull(), err
	}

	return types.StringValue(overriddenValues), nil
}

// restoreBaseFieldsFromState restores base fields from state (excludes _full computed variants)
// Used during READ to preserve user input while allowing _full fields to reflect API state for drift detection
func restoreBaseFieldsFromState(requestContext *common.RequestContext, source process.Getter, tfModel *dashboardtf.DashboardTFModel) error {

	var originalModel dashboardtf.DashboardTFModel
	diags := source.Get(requestContext.Context, &originalModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original model: %s", diags.Errors())
	}

	// Restore base fields from state (skip _full variants to enable drift detection)
	tfModel.Groups = originalModel.Groups
	tfModel.Values = originalModel.Values

	return nil
}
