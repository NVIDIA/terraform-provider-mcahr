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
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	actiontf "terraform/terraform-provider/provider/tf/resource/action/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getActionCommonParameters returns the common parameter string for define/update operations
func getActionCommonParameters() string {
	return `action_name="test-action", ` +
		`command="echo 'test'", ` +
		`enabled=true, ` +
		`timeout=5000, ` +
		`description="Test description", ` +
		`res_env_var="TEST_VAR", ` +
		`resource_query="hosts", ` +
		`shell="/bin/bash", ` +
		`allowed_resources_query="allowed_hosts", ` +
		`communication_workspace="ops-workspace", ` +
		`communication_channel="alerts-channel", ` +
		`start_title_template="started test action", ` +
		`start_short_template="started test short", ` +
		`complete_title_template="completed test action", ` +
		`complete_short_template="completed test short", ` +
		`error_title_template="failed test action", ` +
		`error_short_template="failed test short", ` +
		`params=["param1", "param2"], ` +
		`resource_tags_to_export=["tag1", "tag2"], ` +
		`file_deps=["file1.sh", "file2.py"], ` +
		`allowed_entities=["entity1", "entity2"], ` +
		`editors=["user1", "user2"]`
}

func TestActionTranslatorCommon_ToAPIModel(t *testing.T) {
	// Given
	translator := &ActionTranslatorCommon{}

	// Create test TF model
	tfModel := &actiontf.ActionTFModel{
		Name:                   types.StringValue("test-action"),
		Command:                types.StringValue("echo 'test'"),
		Enabled:                types.BoolValue(true),
		Timeout:                types.Int64Value(5000),
		Description:            types.StringValue("Test description"),
		ResEnvVar:              types.StringValue("TEST_VAR"),
		ResourceQuery:          types.StringValue("hosts"),
		Shell:                  types.StringValue("/bin/bash"),
		AllowedResourcesQuery:  types.StringValue("allowed_hosts"),
		CommunicationWorkspace: types.StringValue("ops-workspace"),
		CommunicationChannel:   types.StringValue("alerts-channel"),
		StartTitleTemplate:     types.StringValue("started test action"),
		StartShortTemplate:     types.StringValue("started test short"),
		CompleteTitleTemplate:  types.StringValue("completed test action"),
		CompleteShortTemplate:  types.StringValue("completed test short"),
		ErrorTitleTemplate:     types.StringValue("failed test action"),
		ErrorShortTemplate:     types.StringValue("failed test short"),
	}

	// Initialize sets with test values
	tfModel.Params, _ = types.ListValue(types.StringType, []attr.Value{
		types.StringValue("param1"),
		types.StringValue("param2"),
	})
	tfModel.ResourceTagsToExport, _ = types.ListValue(types.StringType, []attr.Value{
		types.StringValue("tag1"),
		types.StringValue("tag2"),
	})
	tfModel.FileDeps, _ = types.ListValue(types.StringType, []attr.Value{
		types.StringValue("file1.sh"),
		types.StringValue("file2.py"),
	})
	tfModel.AllowedEntities, _ = types.ListValue(types.StringType, []attr.Value{
		types.StringValue("entity1"),
		types.StringValue("entity2"),
	})
	tfModel.Editors, _ = types.ListValue(types.StringType, []attr.Value{
		types.StringValue("user1"),
		types.StringValue("user2"),
	})

	tests := []struct {
		name      string
		operation common.CrudOperation
		expected  string
	}{
		{
			name:      "Create operation",
			operation: common.Create,
			expected:  fmt.Sprintf("define_action(%s)", getActionCommonParameters()),
		},
		{
			name:      "Read operation",
			operation: common.Read,
			expected:  `get_action_class(action_name="test-action")`,
		},
		{
			name:      "Update operation",
			operation: common.Update,
			expected:  fmt.Sprintf("update_action(%s)", getActionCommonParameters()),
		},
		{
			name:      "Delete operation",
			operation: common.Delete,
			expected:  `delete_action(action_name="test-action")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			requestContext := common.NewRequestContext(context.Background()).WithOperation(tt.operation).WithAPIVersion(common.V2)
			translationData := &coretranslator.TranslationData{}
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

func TestActionTranslatorCommon_ToAPIModel_UnsupportedOperation(t *testing.T) {
	// Given
	translator := &ActionTranslatorCommon{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.CrudOperation(999)).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}
	tfModel := &actiontf.ActionTFModel{
		Name: types.StringValue("test-action"),
	}

	// When
	// Test with an invalid operation (cast to avoid compile error)
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}
