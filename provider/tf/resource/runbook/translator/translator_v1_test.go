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
	runbookapi "terraform/terraform-provider/provider/external_api/resources/runbooks"
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	runbooktf "terraform/terraform-provider/provider/tf/resource/runbook/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunbookTranslatorV1_ToTFModel(t *testing.T) {
	tests := []struct {
		name        string
		apiModel    *runbookapi.RunbookResponseAPIModelV1
		expectError bool
		expectNil   bool
		validate    func(t *testing.T, tfModel *runbooktf.RunbookTFModel)
	}{
		{
			name:        "Nil API model",
			apiModel:    nil,
			expectError: false,
			expectNil:   true,
		},
		{
			name: "Valid define_notebook response",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				DefineNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{
						{
							Name:                                "test_runbook",
							Enabled:                             true,
							Description:                         "Test description",
							TimeoutMs:                           5000,
							AllowedResourcesQuery:               "resource_query",
							CommunicationWorkspace:              "workspace",
							CommunicationChannel:                "channel",
							Category:                            "test_category",
							IsRunOutputPersisted:                true,
							FilterResourceToAction:              false,
							CommunicationCudNotifications:       true,
							CommunicationApprovalNotifications:  false,
							CommunicationExecutionNotifications: true,
							AllowedEntities:                     []string{"entity1", "entity2"},
							Approvers:                           []string{"user1", "user2"},
							Labels:                              []string{"label1", "label2"},
							Editors:                             []string{"editor1"},
							SecretNames:                         []string{"secret1"},
							Cells: []customattribute.CellJsonAPI{
								{
									Type:    "OP_LANG",
									Content: "print('hello')",
									Name:    "cell1",
									Enabled: true,
								},
							},
							Params: []customattribute.ParamJson{
								{
									Name:  "param1",
									Value: "value1",
								},
							},
							ExternalParams: []customattribute.ExternalParamJson{
								{
									Name:   "ext1",
									Source: "api",
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "test_runbook", tfModel.Name.ValueString())
				assert.Equal(t, true, tfModel.Enabled.ValueBool())
				assert.Equal(t, "Test description", tfModel.Description.ValueString())
				assert.Equal(t, int64(5000), tfModel.TimeoutMs.ValueInt64())
				assert.Equal(t, "resource_query", tfModel.AllowedResourcesQuery.ValueString())
				assert.Equal(t, "workspace", tfModel.CommunicationWorkspace.ValueString())
				assert.Equal(t, "channel", tfModel.CommunicationChannel.ValueString())
				assert.Equal(t, "test_category", tfModel.Category.ValueString())
				assert.Equal(t, true, tfModel.IsRunOutputPersisted.ValueBool())
				assert.Equal(t, false, tfModel.FilterResourceToAction.ValueBool())
				assert.Equal(t, true, tfModel.CommunicationCudNotifications.ValueBool())
				assert.Equal(t, false, tfModel.CommunicationApprovalNotifications.ValueBool())
				assert.Equal(t, true, tfModel.CommunicationExecutionNotifications.ValueBool())

				// Check sets
				assert.Equal(t, 2, len(tfModel.AllowedEntities.Elements()))
				allowedEntitiesElements := tfModel.AllowedEntities.Elements()
				assert.Contains(t, allowedEntitiesElements, types.StringValue("entity1"))
				assert.Contains(t, allowedEntitiesElements, types.StringValue("entity2"))

				assert.Equal(t, 2, len(tfModel.Approvers.Elements()))
				approversElements := tfModel.Approvers.Elements()
				assert.Contains(t, approversElements, types.StringValue("user1"))
				assert.Contains(t, approversElements, types.StringValue("user2"))

				assert.Equal(t, 2, len(tfModel.Labels.Elements()))
				labelsElements := tfModel.Labels.Elements()
				assert.Contains(t, labelsElements, types.StringValue("label1"))
				assert.Contains(t, labelsElements, types.StringValue("label2"))

				assert.Equal(t, 1, len(tfModel.Editors.Elements()))
				editorsElements := tfModel.Editors.Elements()
				assert.Contains(t, editorsElements, types.StringValue("editor1"))

				assert.Equal(t, 1, len(tfModel.SecretNames.Elements()))
				secretNamesElements := tfModel.SecretNames.Elements()
				assert.Contains(t, secretNamesElements, types.StringValue("secret1"))

				// Check JSON fields
				assert.Equal(t, tfModel.Cells.ValueString(), "[{\"description\":\"\",\"enabled\":true,\"name\":\"cell1\",\"op\":\"print('hello')\",\"secret_aware\":false}]")
				assert.Equal(t, tfModel.Params.ValueString(), "[{\"description\":\"\",\"export\":false,\"name\":\"param1\",\"param_type\":\"\",\"required\":false,\"value\":\"value1\"}]")
				assert.Equal(t, tfModel.ExternalParams.ValueString(), "[{\"description\":\"\",\"export\":false,\"json_path\":\"\",\"name\":\"ext1\",\"param_type\":\"\",\"source\":\"api\",\"value\":\"\"}]")
			},
		},
		{
			name: "Valid update_notebook response",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				UpdateNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{
						{
							Name:    "updated_runbook",
							Enabled: false,
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "updated_runbook", tfModel.Name.ValueString())
				assert.Equal(t, false, tfModel.Enabled.ValueBool())
			},
		},
		{
			name: "Valid get_notebook_class response",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				GetNotebookClass: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{
						{
							Name:        "get_runbook",
							Enabled:     true,
							Description: "Retrieved runbook",
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "get_runbook", tfModel.Name.ValueString())
				assert.Equal(t, true, tfModel.Enabled.ValueBool())
				assert.Equal(t, "Retrieved runbook", tfModel.Description.ValueString())
			},
		},
		{
			name:     "No container in response",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				// All fields nil
			},
			expectError: true,
		},
		{
			name: "Empty notebook classes",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				DefineNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{},
				},
			},
			expectError: true,
		},
		{
			name: "Multiple notebook classes (uses first)",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				DefineNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{
						{Name: "first", Enabled: true},
						{Name: "second", Enabled: false},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "first", tfModel.Name.ValueString())
				assert.Equal(t, true, tfModel.Enabled.ValueBool())
			},
		},
		{
			name: "Empty sets and null values",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				DefineNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{
						{
							Name:            "minimal",
							Enabled:         false,
							AllowedEntities: []string{},
							Approvers:       nil,
							Cells:           []customattribute.CellJsonAPI{},
							Params:          nil,
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "minimal", tfModel.Name.ValueString())
				assert.Equal(t, false, tfModel.Enabled.ValueBool())
				assert.Equal(t, 0, len(tfModel.AllowedEntities.Elements()))
				assert.Equal(t, 0, len(tfModel.Approvers.Elements()))
				assert.Equal(t, "[]", tfModel.Cells.ValueString())
				assert.Equal(t, "null", tfModel.Params.ValueString())
			},
		},
		{
			name: "Category field populated",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				DefineNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{
						{
							Name:     "category_test",
							Category: "infrastructure",
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "category_test", tfModel.Name.ValueString())
				assert.Equal(t, "infrastructure", tfModel.Category.ValueString())
			},
		},
		{
			name: "Category field empty",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				DefineNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{
						{
							Name:     "no_category",
							Category: "",
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "no_category", tfModel.Name.ValueString())
				assert.Equal(t, "", tfModel.Category.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			translator := &RunbookTranslatorV1{}
			requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
			translationData := &coretranslator.TranslationData{}

			// when
			tfModel, err := translator.ToTFModel(requestContext, translationData, tt.apiModel)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.expectNil {
				assert.Nil(t, tfModel)
			} else if !tt.expectError {
				require.NotNil(t, tfModel)
				if tt.validate != nil {
					tt.validate(t, tfModel)
				}
			}
		})
	}
}

func TestRunbookTranslatorV1_JSONFieldHandling(t *testing.T) {
	// given
	translator := &RunbookTranslatorV1{}
	apiModel := &runbookapi.RunbookResponseAPIModelV1{
		DefineNotebook: &runbookapi.RunbookContainerV1{
			NotebookClasses: []runbookapi.NotebookClassV1{
				{
					Name: "json_test",
					Cells: []customattribute.CellJsonAPI{
						{
							Type:    "OP_LANG",
							Content: "print('op cell')",
							Name:    "op_cell",
						},
						{
							Type:    "MARKDOWN",
							Content: "# Markdown",
							Name:    "md_cell",
						},
					},
					Params: []customattribute.ParamJson{
						{
							Name:      "param1",
							Value:     "value1",
							Required:  true,
							ParamType: "PARAM",
						},
					},
					ExternalParams: []customattribute.ExternalParamJson{
						{
							Name:      "ext1",
							Source:    "api_endpoint",
							JsonPath:  "$.data",
							ParamType: "EXTERNAL",
						},
					},
				},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	tfModel, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, tfModel)

	// Check cells were converted to internal model
	assert.Contains(t, tfModel.Cells.ValueString(), `"op":"print('op cell')"`)
	assert.Contains(t, tfModel.Cells.ValueString(), `"md":"# Markdown"`)
	assert.Equal(t, tfModel.Cells.ValueString(), tfModel.CellsFull.ValueString())

	// Check params
	assert.Contains(t, tfModel.Params.ValueString(), "param1")
	assert.Equal(t, tfModel.Params.ValueString(), tfModel.ParamsFull.ValueString())

	// Check external params
	assert.Contains(t, tfModel.ExternalParams.ValueString(), "ext1")
	assert.Equal(t, tfModel.ExternalParams.ValueString(), tfModel.ExternalParamsFull.ValueString())
}

func TestRunbookTranslatorV1_GetContainer(t *testing.T) {
	tests := []struct {
		name     string
		apiModel *runbookapi.RunbookResponseAPIModelV1
		expected *runbookapi.RunbookContainerV1
	}{
		{
			name: "DefineNotebook container",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				DefineNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{{Name: "define"}},
				},
			},
			expected: &runbookapi.RunbookContainerV1{
				NotebookClasses: []runbookapi.NotebookClassV1{{Name: "define"}},
			},
		},
		{
			name: "UpdateNotebook container",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				UpdateNotebook: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{{Name: "update"}},
				},
			},
			expected: &runbookapi.RunbookContainerV1{
				NotebookClasses: []runbookapi.NotebookClassV1{{Name: "update"}},
			},
		},
		{
			name: "GetNotebookClass container",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{
				GetNotebookClass: &runbookapi.RunbookContainerV1{
					NotebookClasses: []runbookapi.NotebookClassV1{{Name: "get"}},
				},
			},
			expected: &runbookapi.RunbookContainerV1{
				NotebookClasses: []runbookapi.NotebookClassV1{{Name: "get"}},
			},
		},
		{
			name:     "No container",
			apiModel: &runbookapi.RunbookResponseAPIModelV1{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			container := tt.apiModel.GetContainer()

			// then
			if tt.expected == nil {
				assert.Nil(t, container)
			} else {
				require.NotNil(t, container)
				assert.Equal(t, tt.expected.NotebookClasses[0].Name, container.NotebookClasses[0].Name)
			}
		})
	}
}

func TestRunbookTranslatorV1_ComplexCellsJSON(t *testing.T) {
	// given
	translator := &RunbookTranslatorV1{}
	apiModel := &runbookapi.RunbookResponseAPIModelV1{
		DefineNotebook: &runbookapi.RunbookContainerV1{
			NotebookClasses: []runbookapi.NotebookClassV1{
				{
					Name: "complex_cells",
					Cells: []customattribute.CellJsonAPI{
						{
							Type:        "OP_LANG",
							Content:     "import json\ndata = {'key': 'value'}\nprint(json.dumps(data))",
							Name:        "json_cell",
							Enabled:     true,
							SecretAware: true,
							Description: "JSON processing cell",
						},
						{
							CellType:    "MARKDOWN", // Using CellType instead of Type
							Content:     "## Complex Markdown\n\n- Item 1\n- Item 2\n\n```python\ncode_block()\n```",
							Name:        "complex_md",
							Enabled:     false,
							SecretAware: false,
						},
					},
				},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	tfModel, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, tfModel)

	// Verify complex content is preserved
	cellsStr := tfModel.Cells.ValueString()
	assert.Contains(t, cellsStr, "import json")
	assert.Contains(t, cellsStr, "json.dumps(data)")
	assert.Contains(t, cellsStr, "## Complex Markdown")
	assert.Contains(t, cellsStr, "code_block()")

	// Verify cell metadata
	assert.Contains(t, cellsStr, "json_cell")
	assert.Contains(t, cellsStr, "complex_md")
}
