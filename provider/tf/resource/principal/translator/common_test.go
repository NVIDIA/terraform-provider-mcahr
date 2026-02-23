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
	"testing"

	"terraform/terraform-provider/provider/common"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	principaltf "terraform/terraform-provider/provider/tf/resource/principal/model"
	"terraform/terraform-provider/provider/tf/resource/principal/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrincipalTranslatorCommon_ToAPIModel(t *testing.T) {
	tests := []struct {
		name      string
		operation common.CrudOperation
	}{
		{"Create operation", common.Create},
		{"Read operation", common.Read},
		{"Update operation", common.Update},
		{"Delete operation", common.Delete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			translator := &PrincipalTranslatorCommon{}
			tfModel := &principaltf.PrincipalTFModel{
				Name:                 types.StringValue("test_principal"),
				Identity:             types.StringValue("test@example.com"),
				ActionLimit:          types.Int64Value(50),
				ExecuteLimit:         types.Int64Value(25),
				ViewLimit:            types.Int64Value(50),
				ConfigurePermission:  types.BoolValue(true),
				AdministerPermission: types.BoolValue(false),
				IDPName:              types.StringValue("okta"),
			}
			requestContext := common.NewRequestContext(context.Background()).WithOperation(tt.operation).WithAPIVersion(common.V1)
			translationData := &coretranslator.TranslationData{}

			// When
			result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

			// Then
			assert.NoError(t, err)
			require.NotNil(t, result)

			switch tt.operation {
			case common.Create:
				expected := "define_principal(" +
					"principal_name=\"test_principal\", " +
					"identity=\"test@example.com\", " +
					"action_limit=50, " +
					"execute_limit=25, " +
					"view_limit=50, " +
					"configure_permission=1, " +
					"administer_permission=0, " +
					"idp_name=\"okta\")"
				assert.Equal(t, expected, result.Statement)
			case common.Update:
				expected := "update_principal(" +
					"principal_name=\"test_principal\", " +
					"identity=\"test@example.com\", " +
					"action_limit=50, " +
					"execute_limit=25, " +
					"view_limit=50, " +
					"configure_permission=1, " +
					"administer_permission=0, " +
					"idp_name=\"okta\")"
				assert.Equal(t, expected, result.Statement)
			case common.Read:
				assert.Equal(t, "get_principal_class(name=\"test_principal\")", result.Statement)
			case common.Delete:
				assert.Equal(t, "delete_principal(principal_name=\"test_principal\")", result.Statement)
			}
		})
	}
}

func TestPrincipalTranslatorCommon_ToAPIModel_UnsupportedOperation(t *testing.T) {
	// Given
	translator := &PrincipalTranslatorCommon{}
	tfModel := &principaltf.PrincipalTFModel{
		Name: types.StringValue("test_principal"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.CrudOperation(999)).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}

func TestPrincipalTranslatorCommon_buildCreateStatement_MinimalFields(t *testing.T) {
	// Given
	translator := &PrincipalTranslatorCommon{}
	tfModel := &principaltf.PrincipalTFModel{
		Name:                 types.StringValue("minimal_principal"),
		Identity:             types.StringValue("minimal@example.com"),
		ActionLimit:          types.Int64Value(0),
		ExecuteLimit:         types.Int64Value(0),
		ViewLimit:            types.Int64Value(0),
		ConfigurePermission:  types.BoolValue(false),
		AdministerPermission: types.BoolValue(false),
		IDPName:              types.StringValue(""),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	principalSchema := schema.PrincipalSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: principalSchema.GetCompatibilityOptions()}

	// When
	statement := translator.buildCreateStatement(requestContext, translationData, tfModel)

	// Then
	expected := "define_principal(" +
		"principal_name=\"minimal_principal\", " +
		"identity=\"minimal@example.com\", " +
		"action_limit=0, " +
		"execute_limit=0, " +
		"view_limit=0, " +
		"configure_permission=0, " +
		"administer_permission=0, " +
		"idp_name=\"\")"
	assert.Equal(t, expected, statement)
}

func TestPrincipalTranslatorCommon_buildCreateStatement_AllFields(t *testing.T) {
	// Given
	translator := &PrincipalTranslatorCommon{}
	tfModel := &principaltf.PrincipalTFModel{
		Name:                 types.StringValue("full_principal"),
		Identity:             types.StringValue("full@example.com"),
		ActionLimit:          types.Int64Value(100),
		ExecuteLimit:         types.Int64Value(50),
		ViewLimit:            types.Int64Value(50),
		ConfigurePermission:  types.BoolValue(true),
		AdministerPermission: types.BoolValue(true),
		IDPName:              types.StringValue("azure"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	principalSchema := schema.PrincipalSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: principalSchema.GetCompatibilityOptions()}

	// When
	statement := translator.buildCreateStatement(requestContext, translationData, tfModel)

	// Then
	expected := "define_principal(" +
		"principal_name=\"full_principal\", " +
		"identity=\"full@example.com\", " +
		"action_limit=100, " +
		"execute_limit=50, " +
		"view_limit=50, " +
		"configure_permission=1, " +
		"administer_permission=1, " +
		"idp_name=\"azure\")"
	assert.Equal(t, expected, statement)
}
