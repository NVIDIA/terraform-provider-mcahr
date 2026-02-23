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
	"terraform/terraform-provider/provider/common"
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	"testing"

	actionapi "terraform/terraform-provider/provider/external_api/resources/actions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestActionResponseV1() *actionapi.ActionResponseAPIModelV1 {
	return &actionapi.ActionResponseAPIModelV1{
		DefineAction: &actionapi.ActionContainerV1{
			Error: apicommon.ErrorV1{
				Type:             "OK",
				Message:          "",
				ValidationErrors: []apicommon.ValidationError{},
			},

			ActionClasses: []actionapi.ActionClassV1{
				{
					Timeout:               30000,
					Command:               "echo 'test command'",
					Enabled:               true,
					Name:                  "test-action",
					Description:           "Test action description",
					Shell:                 "/bin/bash",
					Params:                `["param1", "param2"]`,
					ResourceQuery:         "resource.type = \"server\"",
					ResourceTagsToExport:  `["tag1", "tag2"]`,
					ResEnvVar:             "RESULT_VAR",
					FileDeps:              `["file1.sh", "file2.sh"]`,
					AllowedEntities:       []string{"entity1", "entity2"},
					AllowedResourcesQuery: "allowed.type = \"resource\"",
					Editors:               []string{"editor1", "editor2"},
					StartStepClass: actionapi.StepClassV1{
						TitleTemplate: "Test action started",
						ShortTemplate: "Action started",
					},
					ErrorStepClass: actionapi.StepClassV1{
						TitleTemplate: "Test action failed",
						ShortTemplate: "Action failed",
					},
					CompleteStepClass: actionapi.StepClassV1{
						TitleTemplate: "Test action completed",
						ShortTemplate: "Action completed",
					},
					Communication: actionapi.CommunicationV1{
						Channel:   "alerts",
						Workspace: "main",
					},
				},
			},
		},
	}
}

func createMinimalActionResponseV1() *actionapi.ActionResponseAPIModelV1 {
	return &actionapi.ActionResponseAPIModelV1{
		DefineAction: &actionapi.ActionContainerV1{
			Error: apicommon.ErrorV1{
				Type:             "OK",
				Message:          "",
				ValidationErrors: []apicommon.ValidationError{},
			},

			ActionClasses: []actionapi.ActionClassV1{
				{
					Timeout:               5000,
					Command:               "echo 'minimal'",
					Enabled:               true,
					Name:                  "minimal_action",
					Description:           "",
					Shell:                 "",
					Params:                "",
					ResourceQuery:         "",
					ResourceTagsToExport:  "",
					ResEnvVar:             "",
					FileDeps:              "",
					AllowedEntities:       []string{},
					AllowedResourcesQuery: "",
					Editors:               []string{},
					StartStepClass: actionapi.StepClassV1{
						TitleTemplate: "",
						ShortTemplate: "",
					},
					ErrorStepClass: actionapi.StepClassV1{
						TitleTemplate: "",
						ShortTemplate: "",
					},
					CompleteStepClass: actionapi.StepClassV1{
						TitleTemplate: "",
						ShortTemplate: "",
					},
					Communication: actionapi.CommunicationV1{
						Channel:   "",
						Workspace: "",
					},
				},
			},
		},
	}
}

func TestActionTranslatorV1_ToTFModel(t *testing.T) {
	// Given
	translator := &ActionTranslatorV1{}
	apiModel := createTestActionResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test-action", result.Name.ValueString())
	assert.Equal(t, "echo 'test command'", result.Command.ValueString())
	assert.True(t, result.Enabled.ValueBool())
	assert.Equal(t, int64(30000), result.Timeout.ValueInt64())
	assert.Equal(t, "Test action description", result.Description.ValueString())
	assert.Equal(t, "RESULT_VAR", result.ResEnvVar.ValueString())
	assert.Equal(t, "resource.type = \"server\"", result.ResourceQuery.ValueString())
	assert.Equal(t, "/bin/bash", result.Shell.ValueString())
	assert.Equal(t, "allowed.type = \"resource\"", result.AllowedResourcesQuery.ValueString())
	assert.Equal(t, "main", result.CommunicationWorkspace.ValueString())
	assert.Equal(t, "alerts", result.CommunicationChannel.ValueString())

	// Verify list fields by extracting values
	var params []string
	result.Params.ElementsAs(context.Background(), &params, false)
	assert.Equal(t, []string{"param1", "param2"}, params)

	var resourceTags []string
	result.ResourceTagsToExport.ElementsAs(context.Background(), &resourceTags, false)
	assert.Equal(t, []string{"tag1", "tag2"}, resourceTags)

	var fileDeps []string
	result.FileDeps.ElementsAs(context.Background(), &fileDeps, false)
	assert.Equal(t, []string{"file1.sh", "file2.sh"}, fileDeps)

	var allowedEntities []string
	result.AllowedEntities.ElementsAs(context.Background(), &allowedEntities, false)
	assert.Equal(t, []string{"entity1", "entity2"}, allowedEntities)

	var editors []string
	result.Editors.ElementsAs(context.Background(), &editors, false)
	assert.Equal(t, []string{"editor1", "editor2"}, editors)

	// Verify template fields mapped from step details
	assert.Equal(t, "Test action started", result.StartTitleTemplate.ValueString())
	assert.Equal(t, "Action started", result.StartShortTemplate.ValueString())
	assert.Equal(t, "Test action completed", result.CompleteTitleTemplate.ValueString())
	assert.Equal(t, "Action completed", result.CompleteShortTemplate.ValueString())
	assert.Equal(t, "Test action failed", result.ErrorTitleTemplate.ValueString())
	assert.Equal(t, "Action failed", result.ErrorShortTemplate.ValueString())
}

func TestActionTranslatorV1_ToTFModel_Nil(t *testing.T) {
	// Given
	translator := &ActionTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// Then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestActionTranslatorV1_ToTFModel_EmptyConfigurations(t *testing.T) {
	// Given
	translator := &ActionTranslatorV1{}
	apiModel := &actionapi.ActionResponseAPIModelV1{
		DefineAction: &actionapi.ActionContainerV1{
			ActionClasses: []actionapi.ActionClassV1{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no action classes found")
}

func TestActionTranslatorV1_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &ActionTranslatorV1{}
	apiModel := createMinimalActionResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify basic fields
	assert.Equal(t, "minimal_action", result.Name.ValueString())
	assert.Equal(t, "echo 'minimal'", result.Command.ValueString())
	assert.True(t, result.Enabled.ValueBool())
	assert.Equal(t, int64(5000), result.Timeout.ValueInt64())

	// Verify that empty optional fields are empty strings (preserving API response values)
	assert.Equal(t, "", result.Description.ValueString())
	assert.Equal(t, "", result.ResEnvVar.ValueString())
	assert.Equal(t, "", result.ResourceQuery.ValueString())
	assert.Equal(t, "", result.Shell.ValueString())
	assert.Equal(t, "", result.AllowedResourcesQuery.ValueString())
	assert.Equal(t, "", result.CommunicationWorkspace.ValueString())
	assert.Equal(t, "", result.CommunicationChannel.ValueString())

	// Verify that template fields are empty strings (from step details)
	assert.Equal(t, "", result.StartTitleTemplate.ValueString())
	assert.Equal(t, "", result.StartShortTemplate.ValueString())
	assert.Equal(t, "", result.CompleteTitleTemplate.ValueString())
	assert.Equal(t, "", result.CompleteShortTemplate.ValueString())
	assert.Equal(t, "", result.ErrorTitleTemplate.ValueString())
	assert.Equal(t, "", result.ErrorShortTemplate.ValueString())

	// Verify that empty lists are empty sets
	assert.Equal(t, 0, len(result.Params.Elements()))
	assert.Equal(t, 0, len(result.ResourceTagsToExport.Elements()))
	assert.Equal(t, 0, len(result.FileDeps.Elements()))
	assert.Equal(t, 0, len(result.AllowedEntities.Elements()))
	assert.Equal(t, 0, len(result.Editors.Elements()))
}
