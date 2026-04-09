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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunbookTranslator_ToTFModel(t *testing.T) {
	tests := []struct {
		name        string
		apiModel    *runbookapi.RunbookResponseAPIModel
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
			name: "No configurations",
			apiModel: &runbookapi.RunbookResponseAPIModel{
				Output: runbookapi.RunbookOutput{
					Configurations: runbookapi.RunbookConfigurations{
						Items: []runbookapi.ConfigurationItem{},
					},
				},
			},
			expectError: true,
		},
		{
			name: "Valid configuration",
			apiModel: &runbookapi.RunbookResponseAPIModel{
				Output: runbookapi.RunbookOutput{
					Configurations: runbookapi.RunbookConfigurations{
						Items: []runbookapi.ConfigurationItem{
							{
								EntityMetadata: runbookapi.RunbookEntityMetadata{
									Name:        "test_runbook",
									Enabled:     true,
									Description: "Test description",
								},
								Config: runbookapi.RunbookConfig{
									TimeoutMs:              5000,
									AllowedResourcesQuery:  "resource_query",
									Category:               "test_category",
									IsRunOutputPersisted:   true,
									FilterResourceToAction: false,
									CommunicationDestination: runbookapi.CommunicationDestination{
										Workspace: "workspace",
										Channel:   "channel",
									},
									CommunicationFilters: runbookapi.CommunicationFilters{
										CudNotifications:       true,
										ApprovalNotifications:  false,
										ExecutionNotifications: true,
									},
									AllowedEntities: []string{"entity1", "entity2"},
									Approvers:       []string{"user1", "user2"},
									Labels:          []string{"label1", "label2"},
									Editors:         []string{"editor1"},
									SecretNames:     []string{"secret1"},
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
				assert.Equal(t, []attr.Value{types.StringValue("entity1"), types.StringValue("entity2")}, tfModel.AllowedEntities.Elements())
				assert.Equal(t, []attr.Value{types.StringValue("user1"), types.StringValue("user2")}, tfModel.Approvers.Elements())
				assert.Equal(t, []attr.Value{types.StringValue("label1"), types.StringValue("label2")}, tfModel.Labels.Elements())
				assert.Equal(t, []attr.Value{types.StringValue("editor1")}, tfModel.Editors.Elements())
				assert.Equal(t, []attr.Value{types.StringValue("secret1")}, tfModel.SecretNames.Elements())

				// Check JSON fields
				assert.Equal(t, tfModel.Cells.ValueString(), "[{\"description\":\"\",\"enabled\":true,\"name\":\"cell1\",\"op\":\"print('hello')\",\"secret_aware\":false}]")
				// Check cells_list
				assert.False(t, tfModel.CellsList.IsNull())
				require.Equal(t, 1, len(tfModel.CellsList.Elements()))
				cellObj := tfModel.CellsList.Elements()[0].(types.Object)
				assert.Equal(t, "print('hello')", cellObj.Attributes()["op"].(types.String).ValueString())
				assert.Equal(t, "cell1", cellObj.Attributes()["name"].(types.String).ValueString())
				assert.Equal(t, true, cellObj.Attributes()["enabled"].(types.Bool).ValueBool())
				assert.Equal(t, false, cellObj.Attributes()["secret_aware"].(types.Bool).ValueBool())
				assert.Equal(t, "", cellObj.Attributes()["description"].(types.String).ValueString())
				assert.True(t, cellObj.Attributes()["md"].(types.String).IsNull())
				assert.Equal(t, tfModel.Params.ValueString(), "[{\"description\":\"\",\"export\":false,\"name\":\"param1\",\"param_type\":\"\",\"required\":false,\"value\":\"value1\"}]")
				assert.Equal(t, tfModel.ExternalParams.ValueString(), "[{\"description\":\"\",\"export\":false,\"json_path\":\"\",\"name\":\"ext1\",\"param_type\":\"\",\"source\":\"api\",\"value\":\"\"}]")
			},
		},
		{
			name: "Multiple configurations (uses first)",
			apiModel: &runbookapi.RunbookResponseAPIModel{
				Output: runbookapi.RunbookOutput{
					Configurations: runbookapi.RunbookConfigurations{
						Items: []runbookapi.ConfigurationItem{
							{
								EntityMetadata: runbookapi.RunbookEntityMetadata{
									Name:    "first",
									Enabled: true,
								},
								Config: runbookapi.RunbookConfig{
									TimeoutMs: 1000,
								},
							},
							{
								EntityMetadata: runbookapi.RunbookEntityMetadata{
									Name:    "second",
									Enabled: false,
								},
								Config: runbookapi.RunbookConfig{
									TimeoutMs: 2000,
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "first", tfModel.Name.ValueString())
				assert.Equal(t, true, tfModel.Enabled.ValueBool())
				assert.Equal(t, int64(1000), tfModel.TimeoutMs.ValueInt64())
			},
		},
		{
			name: "Minimal configuration",
			apiModel: &runbookapi.RunbookResponseAPIModel{
				Output: runbookapi.RunbookOutput{
					Configurations: runbookapi.RunbookConfigurations{
						Items: []runbookapi.ConfigurationItem{
							{
								EntityMetadata: runbookapi.RunbookEntityMetadata{
									Name:    "minimal",
									Enabled: false,
								},
								Config: runbookapi.RunbookConfig{
									TimeoutMs:                0,
									CommunicationDestination: runbookapi.CommunicationDestination{},
									CommunicationFilters:     runbookapi.CommunicationFilters{},
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "minimal", tfModel.Name.ValueString())
				assert.Equal(t, false, tfModel.Enabled.ValueBool())
				assert.Equal(t, int64(0), tfModel.TimeoutMs.ValueInt64())
				assert.Equal(t, "", tfModel.CommunicationWorkspace.ValueString())
				assert.Equal(t, "", tfModel.CommunicationChannel.ValueString())
			},
		},
		{
			name: "Empty arrays and null values",
			apiModel: &runbookapi.RunbookResponseAPIModel{
				Output: runbookapi.RunbookOutput{
					Configurations: runbookapi.RunbookConfigurations{
						Items: []runbookapi.ConfigurationItem{
							{
								EntityMetadata: runbookapi.RunbookEntityMetadata{
									Name: "empty_arrays",
								},
								Config: runbookapi.RunbookConfig{
									AllowedEntities: []string{},
									Approvers:       nil,
									Cells:           []customattribute.CellJsonAPI{},
									Params:          nil,
									ExternalParams:  []customattribute.ExternalParamJson{},
								},
							},
						},
					},
				},
			},
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "empty_arrays", tfModel.Name.ValueString())
				assert.Equal(t, 0, len(tfModel.AllowedEntities.Elements()))
				assert.Equal(t, 0, len(tfModel.Approvers.Elements()))
				assert.Equal(t, "[]", tfModel.Cells.ValueString())
				assert.Equal(t, 0, len(tfModel.CellsList.Elements()))
				assert.Equal(t, "null", tfModel.Params.ValueString())
				assert.Equal(t, "[]", tfModel.ExternalParams.ValueString())
			},
		},
		{
			name: "Category field populated",
			apiModel: &runbookapi.RunbookResponseAPIModel{
				Output: runbookapi.RunbookOutput{
					Configurations: runbookapi.RunbookConfigurations{
						Items: []runbookapi.ConfigurationItem{
							{
								EntityMetadata: runbookapi.RunbookEntityMetadata{
									Name: "category_test",
								},
								Config: runbookapi.RunbookConfig{
									Category: "infrastructure",
								},
							},
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
			apiModel: &runbookapi.RunbookResponseAPIModel{
				Output: runbookapi.RunbookOutput{
					Configurations: runbookapi.RunbookConfigurations{
						Items: []runbookapi.ConfigurationItem{
							{
								EntityMetadata: runbookapi.RunbookEntityMetadata{
									Name: "no_category",
								},
								Config: runbookapi.RunbookConfig{
									Category: "",
								},
							},
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
			translator := &RunbookTranslator{}
			requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
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

func TestToTFModelJsonFields(t *testing.T) {
	tests := []struct {
		name           string
		cells          []customattribute.CellJsonAPI
		params         []customattribute.ParamJson
		externalParams []customattribute.ExternalParamJson
		expectError    bool
		validate       func(t *testing.T, tfModel *runbooktf.RunbookTFModel)
	}{
		{
			name: "Valid JSON fields",
			cells: []customattribute.CellJsonAPI{
				{
					Type:    "OP_LANG",
					Content: "print('test')",
					Name:    "cell1",
				},
			},
			params: []customattribute.ParamJson{
				{
					Name:  "param1",
					Value: "value1",
				},
			},
			externalParams: []customattribute.ExternalParamJson{
				{
					Name:   "ext1",
					Source: "api",
				},
			},
			expectError: false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				// Check cells
				assert.Equal(t, "[{\"description\":\"\",\"enabled\":false,\"name\":\"cell1\",\"op\":\"print('test')\",\"secret_aware\":false}]", tfModel.Cells.ValueString())
				assert.Equal(t, tfModel.Cells.ValueString(), tfModel.CellsFull.ValueString())
				// Check cells_list
				assert.False(t, tfModel.CellsList.IsNull())
				require.Equal(t, 1, len(tfModel.CellsList.Elements()))
				cellObj := tfModel.CellsList.Elements()[0].(types.Object)
				assert.Equal(t, "print('test')", cellObj.Attributes()["op"].(types.String).ValueString())
				assert.Equal(t, "cell1", cellObj.Attributes()["name"].(types.String).ValueString())
				assert.Equal(t, false, cellObj.Attributes()["enabled"].(types.Bool).ValueBool())
				assert.Equal(t, false, cellObj.Attributes()["secret_aware"].(types.Bool).ValueBool())
				assert.True(t, cellObj.Attributes()["md"].(types.String).IsNull())

				// Check params
				assert.Equal(t, "[{\"description\":\"\",\"export\":false,\"name\":\"param1\",\"param_type\":\"\",\"required\":false,\"value\":\"value1\"}]", tfModel.Params.ValueString())
				assert.Equal(t, tfModel.Params.ValueString(), tfModel.ParamsFull.ValueString())

				// Check external params
				assert.Equal(t, "[{\"description\":\"\",\"export\":false,\"json_path\":\"\",\"name\":\"ext1\",\"param_type\":\"\",\"source\":\"api\",\"value\":\"\"}]", tfModel.ExternalParams.ValueString())
				assert.Equal(t, tfModel.ExternalParams.ValueString(), tfModel.ExternalParamsFull.ValueString())
			},
		},
		{
			name:           "Empty arrays",
			cells:          []customattribute.CellJsonAPI{},
			params:         []customattribute.ParamJson{},
			externalParams: []customattribute.ExternalParamJson{},
			expectError:    false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "[]", tfModel.Cells.ValueString())
				assert.Equal(t, 0, len(tfModel.CellsList.Elements()))
				assert.Equal(t, "[]", tfModel.Params.ValueString())
				assert.Equal(t, "[]", tfModel.ExternalParams.ValueString())
			},
		},
		{
			name:           "Nil arrays",
			cells:          nil,
			params:         nil,
			externalParams: nil,
			expectError:    false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.Equal(t, "[]", tfModel.Cells.ValueString())
				assert.True(t, tfModel.CellsList.IsNull())
				assert.Equal(t, "null", tfModel.Params.ValueString())
				assert.Equal(t, "null", tfModel.ExternalParams.ValueString())
			},
		},
		{
			name: "Complex cells with both types",
			cells: []customattribute.CellJsonAPI{
				{
					Type:        "OP_LANG",
					Content:     "import os\nprint(os.environ)",
					Name:        "env_cell",
					Enabled:     true,
					SecretAware: true,
				},
				{
					Type:    "MARKDOWN",
					Content: "# Documentation\n\nThis is a **markdown** cell with `code`",
					Name:    "doc_cell",
					Enabled: false,
				},
			},
			params:         []customattribute.ParamJson{},
			externalParams: []customattribute.ExternalParamJson{},
			expectError:    false,
			validate: func(t *testing.T, tfModel *runbooktf.RunbookTFModel) {
				assert.False(t, tfModel.CellsList.IsNull())
				require.Equal(t, 2, len(tfModel.CellsList.Elements()))
				cell0 := tfModel.CellsList.Elements()[0].(types.Object)
				assert.Equal(t, "env_cell", cell0.Attributes()["name"].(types.String).ValueString())
				assert.Equal(t, "import os\nprint(os.environ)", cell0.Attributes()["op"].(types.String).ValueString())
				assert.Equal(t, true, cell0.Attributes()["enabled"].(types.Bool).ValueBool())
				assert.Equal(t, true, cell0.Attributes()["secret_aware"].(types.Bool).ValueBool())
				assert.True(t, cell0.Attributes()["md"].(types.String).IsNull())
				cell1 := tfModel.CellsList.Elements()[1].(types.Object)
				assert.Equal(t, "doc_cell", cell1.Attributes()["name"].(types.String).ValueString())
				assert.Equal(t, "# Documentation\n\nThis is a **markdown** cell with `code`", cell1.Attributes()["md"].(types.String).ValueString())
				assert.Equal(t, false, cell1.Attributes()["enabled"].(types.Bool).ValueBool())
				assert.Equal(t, false, cell1.Attributes()["secret_aware"].(types.Bool).ValueBool())
				assert.True(t, cell1.Attributes()["op"].(types.String).IsNull())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tfModel := &runbooktf.RunbookTFModel{}

			// when
			err := toTFModelJsonFields(tfModel, tt.cells, tt.params, tt.externalParams)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, tfModel)
				}
			}
		})
	}
}

func TestRunbookTranslator_CompleteV2Response(t *testing.T) {
	// given
	translator := &RunbookTranslator{}
	apiModel := &runbookapi.RunbookResponseAPIModel{
		Output: runbookapi.RunbookOutput{
			Configurations: runbookapi.RunbookConfigurations{
				Items: []runbookapi.ConfigurationItem{
					{
						EntityMetadata: runbookapi.RunbookEntityMetadata{
							Name:        "complete_v2",
							Enabled:     true,
							Description: "Complete V2 response test",
						},
						Config: runbookapi.RunbookConfig{
							TimeoutMs:              10000,
							AllowedResourcesQuery:  "complex_query",
							Category:               "production",
							IsRunOutputPersisted:   true,
							FilterResourceToAction: true,
							CommunicationDestination: runbookapi.CommunicationDestination{
								Workspace: "slack_workspace",
								Channel:   "#alerts",
							},
							CommunicationFilters: runbookapi.CommunicationFilters{
								CudNotifications:       true,
								ApprovalNotifications:  true,
								ExecutionNotifications: false,
							},
							AllowedEntities: []string{"entity1", "entity2", "entity3"},
							Approvers:       []string{"approver1", "approver2"},
							Labels:          []string{"production", "critical", "automated"},
							Editors:         []string{"editor1", "editor2", "editor3"},
							SecretNames:     []string{"api_key", "db_password"},
							Cells: []customattribute.CellJsonAPI{
								{
									Type:        "OP_LANG",
									Content:     "# Complex script\nimport requests\nresponse = requests.get('https://api.example.com')\nprint(response.json())",
									Name:        "api_call",
									Enabled:     true,
									SecretAware: true,
									Description: "API call cell",
								},
								{
									Type:        "MARKDOWN",
									Content:     "## Results\n\nThe API call above retrieves data from our endpoint.",
									Name:        "results_doc",
									Enabled:     true,
									SecretAware: false,
								},
							},
							Params: []customattribute.ParamJson{
								{
									Name:        "api_url",
									Value:       "https://api.example.com",
									Required:    true,
									Export:      true,
									ParamType:   "PARAM",
									Description: "API endpoint URL",
								},
								{
									Name:        "timeout",
									Value:       "30",
									Required:    false,
									Export:      false,
									ParamType:   "PARAM",
									Description: "Request timeout",
								},
							},
							ExternalParams: []customattribute.ExternalParamJson{
								{
									Name:        "auth_token",
									Value:       "",
									Source:      "secrets_manager",
									JsonPath:    "$.auth.token",
									Export:      false,
									ParamType:   "EXTERNAL",
									Description: "Authentication token",
								},
							},
						},
					},
				},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	tfModel, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, tfModel)

	// Verify all fields
	assert.Equal(t, "complete_v2", tfModel.Name.ValueString())
	assert.Equal(t, true, tfModel.Enabled.ValueBool())
	assert.Equal(t, "Complete V2 response test", tfModel.Description.ValueString())
	assert.Equal(t, int64(10000), tfModel.TimeoutMs.ValueInt64())
	assert.Equal(t, "complex_query", tfModel.AllowedResourcesQuery.ValueString())
	assert.Equal(t, "slack_workspace", tfModel.CommunicationWorkspace.ValueString())
	assert.Equal(t, "#alerts", tfModel.CommunicationChannel.ValueString())
	assert.Equal(t, "production", tfModel.Category.ValueString())
	assert.Equal(t, true, tfModel.IsRunOutputPersisted.ValueBool())
	assert.Equal(t, true, tfModel.FilterResourceToAction.ValueBool())
	assert.Equal(t, true, tfModel.CommunicationCudNotifications.ValueBool())
	assert.Equal(t, true, tfModel.CommunicationApprovalNotifications.ValueBool())
	assert.Equal(t, false, tfModel.CommunicationExecutionNotifications.ValueBool())

	// Verify sets
	assert.Equal(t, 3, len(tfModel.AllowedEntities.Elements()))
	assert.Equal(t, 2, len(tfModel.Approvers.Elements()))
	assert.Equal(t, 3, len(tfModel.Labels.Elements()))
	assert.Equal(t, 3, len(tfModel.Editors.Elements()))
	assert.Equal(t, 2, len(tfModel.SecretNames.Elements()))

	// Verify JSON fields
	assert.Equal(t, `[{"description":"API call cell","enabled":true,"name":"api_call","op":"# Complex script\nimport requests\nresponse = requests.get('https://api.example.com')\nprint(response.json())","secret_aware":true},{"description":"","enabled":true,"md":"## Results\n\nThe API call above retrieves data from our endpoint.","name":"results_doc","secret_aware":false}]`, tfModel.Cells.ValueString())

	// Verify cells_list
	assert.False(t, tfModel.CellsList.IsNull())
	require.Equal(t, 2, len(tfModel.CellsList.Elements()))
	cell0 := tfModel.CellsList.Elements()[0].(types.Object)
	assert.Equal(t, "api_call", cell0.Attributes()["name"].(types.String).ValueString())
	assert.Equal(t, "# Complex script\nimport requests\nresponse = requests.get('https://api.example.com')\nprint(response.json())", cell0.Attributes()["op"].(types.String).ValueString())
	assert.True(t, cell0.Attributes()["md"].(types.String).IsNull())
	assert.Equal(t, true, cell0.Attributes()["enabled"].(types.Bool).ValueBool())
	assert.Equal(t, true, cell0.Attributes()["secret_aware"].(types.Bool).ValueBool())
	assert.Equal(t, "API call cell", cell0.Attributes()["description"].(types.String).ValueString())
	cell1 := tfModel.CellsList.Elements()[1].(types.Object)
	assert.Equal(t, "results_doc", cell1.Attributes()["name"].(types.String).ValueString())
	assert.Equal(t, "## Results\n\nThe API call above retrieves data from our endpoint.", cell1.Attributes()["md"].(types.String).ValueString())
	assert.True(t, cell1.Attributes()["op"].(types.String).IsNull())
	assert.Equal(t, true, cell1.Attributes()["enabled"].(types.Bool).ValueBool())
	assert.Equal(t, false, cell1.Attributes()["secret_aware"].(types.Bool).ValueBool())
	assert.Equal(t, "", cell1.Attributes()["description"].(types.String).ValueString())

	assert.Equal(t, `[{"description":"API endpoint URL","export":true,"name":"api_url","param_type":"PARAM","required":true,"value":"https://api.example.com"},{"description":"Request timeout","export":false,"name":"timeout","param_type":"PARAM","required":false,"value":"30"}]`, tfModel.Params.ValueString())
	assert.Equal(t, `[{"description":"Authentication token","export":false,"json_path":"$.auth.token","name":"auth_token","param_type":"EXTERNAL","source":"secrets_manager","value":""}]`, tfModel.ExternalParams.ValueString())
}
