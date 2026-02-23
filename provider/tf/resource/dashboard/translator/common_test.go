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

package translator

import (
	"context"
	"fmt"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	dashboardtf "terraform/terraform-provider/provider/tf/resource/dashboard/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getMinimalDashboardParameters returns parameters with only required fields
func getMinimalDashboardParameters() string {
	config := `{"groups":[],"identifiers":[],"other_tags":[],"resource_query":"host","values":[]}`
	return fmt.Sprintf(`dashboard_name="test-dashboard", dashboard_type="test-type", dashboard_configuration=%s`, utils.EncodeBase64(config))
}

// getFullDashboardParameters returns parameters with all fields populated
func getFullDashboardParameters() string {
	config := `{"groups":[{"name":"group1","tags":["tag1","tag2"]}],"identifiers":["id1","id2"],"other_tags":["tag1","tag2"],"resource_query":"host","values":[{"color":"red","values":["value1","value2"]}]}`
	return fmt.Sprintf(`dashboard_name="test-dashboard", dashboard_type="test-type", dashboard_configuration=%s`, utils.EncodeBase64(config))
}

func TestDashboardTranslatorCommon_ToAPIModel_Minimal(t *testing.T) {
	// Given
	translator := &DashboardTranslatorCommon{}

	// Create minimal test TF model
	tfModel := &dashboardtf.DashboardTFModel{
		Name:          types.StringValue("test-dashboard"),
		DashboardType: types.StringValue("test-type"),
		ResourceQuery: types.StringValue("host"),
		Groups:        types.StringValue("[]"),
		GroupsFull:    types.StringValue("[]"),
		Values:        types.StringValue("[]"),
		ValuesFull:    types.StringValue("[]"),
		OtherTags:     types.ListValueMust(types.StringType, []attr.Value{}),
		Identifiers:   types.ListValueMust(types.StringType, []attr.Value{}),
	}

	tests := []struct {
		name      string
		operation common.CrudOperation
		expected  string
	}{
		{
			name:      "Create operation",
			operation: common.Create,
			expected:  fmt.Sprintf("define_dashboard(%s)", getMinimalDashboardParameters()),
		},
		{
			name:      "Read operation",
			operation: common.Read,
			expected:  `get_dashboard_class(dashboard_name="test-dashboard")`,
		},
		{
			name:      "Update operation",
			operation: common.Update,
			expected:  fmt.Sprintf("update_dashboard(%s)", getMinimalDashboardParameters()),
		},
		{
			name:      "Delete operation",
			operation: common.Delete,
			expected:  `delete_dashboard(dashboard_name="test-dashboard")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			requestContext := common.NewRequestContext(context.Background()).
				WithOperation(tt.operation).
				WithAPIVersion(common.V2).
				WithBackendVersion(&version.BackendVersion{Major: 2, Minor: 0, Patch: 0})
			translationData := &utils.TranslationData{}

			result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

			// Then
			require.NoError(t, err)
			require.NotNil(t, result)

			// Verify exact statement match
			assert.Equal(t, tt.expected, result.Statement)

			// Verify BackendVersion is set correctly
			assert.Equal(t, common.V2, result.APIVersion)
		})
	}
}

func TestDashboardTranslatorCommon_ToAPIModel_Full(t *testing.T) {
	// Given
	translator := &DashboardTranslatorCommon{}

	// Create full test TF model
	tfModel := &dashboardtf.DashboardTFModel{
		Name:          types.StringValue("test-dashboard"),
		DashboardType: types.StringValue("test-type"),
		ResourceQuery: types.StringValue("host"),
		Groups:        types.StringValue(`[{"name":"group1","tags":["tag1","tag2"]}]`),
		GroupsFull:    types.StringValue(`[{"name":"group1","tags":["tag1","tag2"]}]`),
		Values:        types.StringValue(`[{"color":"red","values":["value1","value2"]}]`),
		ValuesFull:    types.StringValue(`[{"color":"red","values":["value1","value2"]}]`),
		OtherTags:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("tag1"), types.StringValue("tag2")}),
		Identifiers:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("id1"), types.StringValue("id2")}),
	}

	tests := []struct {
		name      string
		operation common.CrudOperation
		expected  string
	}{
		{
			name:      "Create operation",
			operation: common.Create,
			expected:  fmt.Sprintf("define_dashboard(%s)", getFullDashboardParameters()),
		},
		{
			name:      "Read operation",
			operation: common.Read,
			expected:  `get_dashboard_class(dashboard_name="test-dashboard")`,
		},
		{
			name:      "Update operation",
			operation: common.Update,
			expected:  fmt.Sprintf("update_dashboard(%s)", getFullDashboardParameters()),
		},
		{
			name:      "Delete operation",
			operation: common.Delete,
			expected:  `delete_dashboard(dashboard_name="test-dashboard")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			requestContext := common.NewRequestContext(context.Background()).
				WithOperation(tt.operation).
				WithAPIVersion(common.V2).
				WithBackendVersion(&version.BackendVersion{Major: 2, Minor: 0, Patch: 0})
			translationData := &utils.TranslationData{}

			result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

			// Then
			require.NoError(t, err)
			require.NotNil(t, result)

			// Verify exact statement match
			assert.Equal(t, tt.expected, result.Statement)

			// Verify BackendVersion is set correctly
			assert.Equal(t, common.V2, result.APIVersion)
		})
	}
}

func TestDashboardTranslatorCommon_ToAPIModel_UnsupportedOperation(t *testing.T) {
	// Given
	translator := &DashboardTranslatorCommon{}
	tfModel := &dashboardtf.DashboardTFModel{
		Name: types.StringValue("test-dashboard"),
	}

	requestContext := common.NewRequestContext(context.Background()).
		WithOperation(common.CrudOperation(999)).
		WithAPIVersion(common.V2).
		WithBackendVersion(&version.BackendVersion{Major: 2, Minor: 0, Patch: 0})
	translationData := &utils.TranslationData{}

	// When
	// Test with an invalid operation (cast to avoid compile error)
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}
