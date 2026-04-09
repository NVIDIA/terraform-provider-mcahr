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
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	"terraform/terraform-provider/provider/tf/resource/runbook/model"
	"terraform/terraform-provider/provider/tf/resource/runbook/schema"
	converters "terraform/terraform-provider/provider/tf/resource/runbook/translator/object_converters"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunbookTranslatorCommon_ToAPIModelWithVersion(t *testing.T) {

	runbookSchema := schema.RunbookSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: runbookSchema.GetCompatibilityOptions()}

	tests := []struct {
		name           string
		tfModel        *model.RunbookTFModel
		operation      common.CrudOperation
		backendVersion *version.BackendVersion
		apiVersion     common.APIVersion
		expectError    bool
		validate       func(t *testing.T, apiModel *statement.StatementInputAPIModel)
	}{
		{
			name: "Create operation",
			tfModel: &model.RunbookTFModel{
				Name:            types.StringValue("test_runbook"),
				AllowedEntities: types.ListNull(types.StringType),
				Approvers:       types.ListNull(types.StringType),
				Labels:          types.ListNull(types.StringType),
				Editors:         types.ListNull(types.StringType),
				SecretNames:     types.ListNull(types.StringType),
				Cells:           types.StringValue(`[{"op": "print('test')", "name": "cell1"}]`),
			},
			operation:      common.Create,
			backendVersion: version.NewBackendVersion("release-29.0.0"),
			apiVersion:     common.V1,
			expectError:    false,
			validate: func(t *testing.T, apiModel *statement.StatementInputAPIModel) {
				assert.Equal(t, "define_notebook(notebook_name=\"test_runbook\", enabled=false, timeout_ms=0, description=\"\", allowed_resources_query=\"\", communication_workspace=\"\", communication_channel=\"\", is_run_output_persisted=false, filter_resource_to_action=false, communication_cud_notifications=false, communication_approval_notifications=false, communication_execution_notifications=false, allowed_entities=[], approvers=[], labels=[], editors=[], secret_names=[], cells=\"W3siY29udGVudCI6InByaW50KCd0ZXN0JykiLCJlbmFibGVkIjp0cnVlLCJuYW1lIjoiY2VsbDEiLCJzZWNyZXRfYXdhcmUiOmZhbHNlLCJ0eXBlIjoiT1BfTEFORyJ9XQ==\", params=, external_params=)", apiModel.Statement)
				assert.Equal(t, common.V1, apiModel.APIVersion)
			},
		},
		{
			name: "Read operation",
			tfModel: &model.RunbookTFModel{
				Name:            types.StringValue("test_runbook"),
				AllowedEntities: types.ListNull(types.StringType),
				Approvers:       types.ListNull(types.StringType),
				Labels:          types.ListNull(types.StringType),
				Editors:         types.ListNull(types.StringType),
				SecretNames:     types.ListNull(types.StringType),
			},
			operation:      common.Read,
			backendVersion: version.NewBackendVersion("release-29.1.0"),
			apiVersion:     common.V2,
			expectError:    false,
			validate: func(t *testing.T, apiModel *statement.StatementInputAPIModel) {
				assert.Equal(t, `get_notebook_class(notebook_name="test_runbook")`, apiModel.Statement)
				assert.Equal(t, common.V2, apiModel.APIVersion)
			},
		},
		{
			name: "Update operation",
			tfModel: &model.RunbookTFModel{
				Name:            types.StringValue("test_runbook"),
				AllowedEntities: types.ListNull(types.StringType),
				Approvers:       types.ListNull(types.StringType),
				Labels:          types.ListNull(types.StringType),
				Editors:         types.ListNull(types.StringType),
				SecretNames:     types.ListNull(types.StringType),
				Cells:           types.StringValue(`[{"op": "print('test')", "name": "cell1"}]`),
			},
			operation:      common.Update,
			backendVersion: version.NewBackendVersion("release-29.0.0"),
			apiVersion:     common.V1,
			expectError:    false,
			validate: func(t *testing.T, apiModel *statement.StatementInputAPIModel) {
				assert.Equal(t, "update_notebook(notebook_name=\"test_runbook\", enabled=false, timeout_ms=0, description=\"\", allowed_resources_query=\"\", communication_workspace=\"\", communication_channel=\"\", is_run_output_persisted=false, filter_resource_to_action=false, communication_cud_notifications=false, communication_approval_notifications=false, communication_execution_notifications=false, allowed_entities=[], approvers=[], labels=[], editors=[], secret_names=[], cells=\"W3siY29udGVudCI6InByaW50KCd0ZXN0JykiLCJlbmFibGVkIjp0cnVlLCJuYW1lIjoiY2VsbDEiLCJzZWNyZXRfYXdhcmUiOmZhbHNlLCJ0eXBlIjoiT1BfTEFORyJ9XQ==\", params=, external_params=)", apiModel.Statement)
				assert.Equal(t, common.V1, apiModel.APIVersion)
			},
		},
		{
			name: "Delete operation",
			tfModel: &model.RunbookTFModel{
				Name:            types.StringValue("test_runbook"),
				AllowedEntities: types.ListNull(types.StringType),
				Approvers:       types.ListNull(types.StringType),
				Labels:          types.ListNull(types.StringType),
				Editors:         types.ListNull(types.StringType),
				SecretNames:     types.ListNull(types.StringType),
			},
			operation:      common.Delete,
			backendVersion: version.NewBackendVersion("release-29.1.0"),
			apiVersion:     common.V2,
			expectError:    false,
			validate: func(t *testing.T, apiModel *statement.StatementInputAPIModel) {
				assert.Equal(t, `delete_notebook(notebook_name="test_runbook")`, apiModel.Statement)
				assert.Equal(t, common.V2, apiModel.APIVersion)
			},
		},
		{
			name: "Unsupported operation",
			tfModel: &model.RunbookTFModel{
				Name:            types.StringValue("test_runbook"),
				AllowedEntities: types.ListNull(types.StringType),
				Approvers:       types.ListNull(types.StringType),
				Labels:          types.ListNull(types.StringType),
				Editors:         types.ListNull(types.StringType),
				SecretNames:     types.ListNull(types.StringType),
			},
			operation:      common.CrudOperation(99), // Invalid operation
			backendVersion: version.NewBackendVersion("release-29.0.0"),
			apiVersion:     common.V1,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			translator := &RunbookTranslatorCommon{}
			requestContext := common.NewRequestContext(context.Background()).WithOperation(tt.operation).WithBackendVersion(tt.backendVersion).WithAPIVersion(tt.apiVersion)

			// when
			apiModel, err := translator.ToAPIModelWithVersion(requestContext, translationData, tt.tfModel)

			// then
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, apiModel)
			} else {
				require.NoError(t, err)
				require.NotNil(t, apiModel)
				if tt.validate != nil {
					tt.validate(t, apiModel)
				}
			}
		})
	}
}

func TestRunbookTranslatorCommon_BuildCreateStatement(t *testing.T) {
	// given
	translator := &RunbookTranslatorCommon{}
	tfModel := createTestTFModel()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	runbookSchema := schema.RunbookSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: runbookSchema.GetCompatibilityOptions()}

	// when
	stmt, err := translator.buildCreateStatement(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	assert.Contains(t, stmt, "define_notebook")
	assert.Contains(t, stmt, "notebook_name=\"test_runbook\"")
	assert.Contains(t, stmt, "enabled=true")
	assert.Contains(t, stmt, "timeout_ms=5000")
	assert.Contains(t, stmt, "description=\"Test description\"")
	assert.Contains(t, stmt, "cells=")
}

func TestRunbookTranslatorCommon_BuildReadStatement(t *testing.T) {
	// given
	translator := &RunbookTranslatorCommon{}
	tfModel := &model.RunbookTFModel{
		Name: types.StringValue("read_test"),
	}

	// when
	stmt := translator.buildReadStatement(tfModel)

	// then
	assert.Equal(t, `get_notebook_class(notebook_name="read_test")`, stmt)
}

func TestRunbookTranslatorCommon_BuildUpdateStatement(t *testing.T) {
	// given
	translator := &RunbookTranslatorCommon{}
	tfModel := createTestTFModel()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Update).WithAPIVersion(common.V1)
	runbookSchema := schema.RunbookSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: runbookSchema.GetCompatibilityOptions()}

	// when
	stmt, err := translator.buildUpdateStatement(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	assert.Contains(t, stmt, "update_notebook")
	assert.Contains(t, stmt, "notebook_name=\"test_runbook\"")
	assert.Contains(t, stmt, "enabled=true")
	assert.Contains(t, stmt, "timeout_ms=5000")
}

func TestRunbookTranslatorCommon_BuildDeleteStatement(t *testing.T) {
	// given
	translator := &RunbookTranslatorCommon{}
	tfModel := &model.RunbookTFModel{
		Name: types.StringValue("delete_test"),
	}

	// when
	stmt := translator.buildDeleteStatement(tfModel)

	// then
	assert.Equal(t, `delete_notebook(notebook_name="delete_test")`, stmt)
}

func TestRunbookTranslatorCommon_BuildRunbookStatement(t *testing.T) {

	runbookSchema := schema.RunbookSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: runbookSchema.GetCompatibilityOptions()}

	tests := []struct {
		name          string
		statementName string
		tfModel       *model.RunbookTFModel
		expectError   bool
		validate      func(t *testing.T, stmt string)
	}{
		{
			name:          "Complete model for define_notebook",
			statementName: "define_notebook",
			tfModel:       createTestTFModel(),
			expectError:   false,
			validate: func(t *testing.T, stmt string) {
				// Check all fields are present
				assert.Contains(t, stmt, "define_notebook(")
				assert.Contains(t, stmt, "notebook_name=\"test_runbook\"")
				assert.Contains(t, stmt, "enabled=true")
				assert.Contains(t, stmt, "timeout_ms=5000")
				assert.Contains(t, stmt, "description=\"Test description\"")
				assert.Contains(t, stmt, "allowed_resources_query=\"resource_query\"")
				assert.Contains(t, stmt, "communication_workspace=\"workspace\"")
				assert.Contains(t, stmt, "category=\"\"")
				assert.Contains(t, stmt, "communication_channel=\"channel\"")
				assert.Contains(t, stmt, "is_run_output_persisted=true")
				assert.Contains(t, stmt, "filter_resource_to_action=false")
				assert.Contains(t, stmt, "communication_cud_notifications=true")
				assert.Contains(t, stmt, "communication_approval_notifications=false")
				assert.Contains(t, stmt, "communication_execution_notifications=true")
				assert.Contains(t, stmt, "allowed_entities=[\"entity1\", \"entity2\"]")
				assert.Contains(t, stmt, "approvers=[\"user1\", \"user2\"]")
				assert.Contains(t, stmt, "labels=[\"label1\", \"label2\"]")
				assert.Contains(t, stmt, "editors=[\"editor1\"]")
				assert.Contains(t, stmt, "secret_names=[\"secret1\"]")
			},
		},
		{
			name:          "Minimal model",
			statementName: "update_notebook",
			tfModel: &model.RunbookTFModel{
				Name:               types.StringValue("minimal"),
				Enabled:            types.BoolValue(false),
				TimeoutMs:          types.Int64Value(1000),
				Cells:              types.StringValue("[]"),
				CellsFull:          types.StringValue("[]"),
				ParamsFull:         types.StringValue("[]"),
				ExternalParamsFull: types.StringValue("[]"),
				AllowedEntities:    types.ListNull(types.StringType),
				Approvers:          types.ListNull(types.StringType),
				Labels:             types.ListNull(types.StringType),
				Editors:            types.ListNull(types.StringType),
				SecretNames:        types.ListNull(types.StringType),
			},
			expectError: false,
			validate: func(t *testing.T, stmt string) {
				assert.Contains(t, stmt, "update_notebook(")
				assert.Contains(t, stmt, "notebook_name=\"minimal\"")
				assert.Contains(t, stmt, "enabled=false")
				assert.Contains(t, stmt, "timeout_ms=1000")
			},
		},
		{
			name:          "Invalid cells JSON",
			statementName: "define_notebook",
			tfModel: &model.RunbookTFModel{
				Name:            types.StringValue("invalid"),
				Cells:           types.StringValue("{invalid json}"),
				AllowedEntities: types.ListNull(types.StringType),
				Approvers:       types.ListNull(types.StringType),
				Labels:          types.ListNull(types.StringType),
				Editors:         types.ListNull(types.StringType),
				SecretNames:     types.ListNull(types.StringType),
			},
			expectError: true,
		},
		{
			name:          "Empty sets",
			statementName: "define_notebook",
			tfModel: &model.RunbookTFModel{
				Name:            types.StringValue("empty_sets"),
				Enabled:         types.BoolValue(true),
				TimeoutMs:       types.Int64Value(2000),
				Cells:           types.StringValue("[]"),
				CellsFull:       types.StringValue("[]"),
				AllowedEntities: types.ListNull(types.StringType),
				Approvers:       types.ListNull(types.StringType),
				Labels:          types.ListNull(types.StringType),
				Editors:         types.ListNull(types.StringType),
				SecretNames:     types.ListNull(types.StringType),
			},
			expectError: false,
			validate: func(t *testing.T, stmt string) {
				assert.Contains(t, stmt, "allowed_entities=[]")
				assert.Contains(t, stmt, "approvers=[]")
			},
		},
		{
			name:          "cells_list path used when set",
			statementName: "define_notebook",
			tfModel: func() *model.RunbookTFModel {
				ctx := context.Background()
				cellObj, _ := types.ObjectValue(
					converters.CellsListAttrTypes,
					map[string]attr.Value{
						"op": types.StringValue("host | limit 1"), "md": types.StringNull(),
						"name": types.StringValue("unnamed"), "enabled": types.BoolValue(true),
						"secret_aware": types.BoolValue(false), "description": types.StringValue(""),
					},
				)
				cellsList, _ := types.ListValue(converters.CellsListObjectType, []attr.Value{cellObj})
				_ = ctx
				return &model.RunbookTFModel{
					Name:               types.StringValue("cells_list_test"),
					Enabled:            types.BoolValue(true),
					TimeoutMs:          types.Int64Value(1000),
					Cells:              types.StringValue("[]"),
					CellsFull:          types.StringValue("[]"),
					CellsList:          cellsList,
					ParamsFull:         types.StringValue("[]"),
					ExternalParamsFull: types.StringValue("[]"),
					AllowedEntities:    types.ListNull(types.StringType),
					Approvers:          types.ListNull(types.StringType),
					Labels:             types.ListNull(types.StringType),
					Editors:            types.ListNull(types.StringType),
					SecretNames:        types.ListNull(types.StringType),
				}
			}(),
			expectError: false,
			validate: func(t *testing.T, stmt string) {
				assert.Contains(t, stmt, "cells=")
				assert.NotContains(t, stmt, `cells="W10="`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			translator := &RunbookTranslatorCommon{}
			requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)

			// when
			stmt, err := translator.buildRunbookStatement(requestContext, translationData, tt.statementName, tt.tfModel)

			// then
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, stmt)
				if tt.validate != nil {
					tt.validate(t, stmt)
				}
			}
		})
	}
}

func TestRunbookTranslatorCommon_CellsHandling(t *testing.T) {
	// given
	translator := &RunbookTranslatorCommon{}
	tfModel := &model.RunbookTFModel{
		Name:      types.StringValue("cells_test"),
		Enabled:   types.BoolValue(true),
		TimeoutMs: types.Int64Value(1000),
		Cells: types.StringValue(`[
			{"op": "print('hello')", "name": "cell1", "enabled": true},
			{"md": "# Header", "name": "cell2", "enabled": false}
		]`),
		CellsFull:          types.StringValue("[]"),
		ParamsFull:         types.StringValue("[]"),
		ExternalParamsFull: types.StringValue("[]"),
		AllowedEntities:    types.ListNull(types.StringType),
		Approvers:          types.ListNull(types.StringType),
		Labels:             types.ListNull(types.StringType),
		Editors:            types.ListNull(types.StringType),
		SecretNames:        types.ListNull(types.StringType),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	runbookSchema := schema.RunbookSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: runbookSchema.GetCompatibilityOptions()}

	// when
	stmt, err := translator.buildCreateStatement(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	assert.Contains(t, stmt, "cells=")
	// The cells should be base64 encoded in the API call
}

// Helper function to create a complete test TF model
func createTestTFModel() *model.RunbookTFModel {
	ctx := context.Background()

	allowedEntities, _ := types.ListValueFrom(ctx, types.StringType, []string{"entity1", "entity2"})
	approvers, _ := types.ListValueFrom(ctx, types.StringType, []string{"user1", "user2"})
	labels, _ := types.ListValueFrom(ctx, types.StringType, []string{"label1", "label2"})
	editors, _ := types.ListValueFrom(ctx, types.StringType, []string{"editor1"})
	secretNames, _ := types.ListValueFrom(ctx, types.StringType, []string{"secret1"})

	return &model.RunbookTFModel{
		Name:                                types.StringValue("test_runbook"),
		Enabled:                             types.BoolValue(true),
		Description:                         types.StringValue("Test description"),
		TimeoutMs:                           types.Int64Value(5000),
		AllowedResourcesQuery:               types.StringValue("resource_query"),
		CommunicationWorkspace:              types.StringValue("workspace"),
		CommunicationChannel:                types.StringValue("channel"),
		IsRunOutputPersisted:                types.BoolValue(true),
		FilterResourceToAction:              types.BoolValue(false),
		CommunicationCudNotifications:       types.BoolValue(true),
		CommunicationApprovalNotifications:  types.BoolValue(false),
		CommunicationExecutionNotifications: types.BoolValue(true),
		AllowedEntities:                     allowedEntities,
		Approvers:                           approvers,
		Labels:                              labels,
		Editors:                             editors,
		SecretNames:                         secretNames,
		Cells:                               types.StringValue(`[{"op": "print('test')", "name": "cell1"}]`),
		CellsFull:                           types.StringValue(`[{"op": "print('test')", "name": "cell1"}]`),
		Params:                              types.StringValue(`[{"name": "param1", "value": "value1"}]`),
		ParamsFull:                          types.StringValue(`[{"name": "param1", "value": "value1"}]`),
		ExternalParams:                      types.StringValue(`[{"name": "ext1", "source": "api"}]`),
		ExternalParamsFull:                  types.StringValue(`[{"name": "ext1", "source": "api"}]`),
	}
}

func TestRunbookTranslatorCommon_NullAndUnknownValues(t *testing.T) {
	// given
	translator := &RunbookTranslatorCommon{}
	tfModel := &model.RunbookTFModel{
		Name:                   types.StringValue("null_test"),
		Enabled:                types.BoolValue(true),
		TimeoutMs:              types.Int64Value(1000),
		Description:            types.StringNull(),
		AllowedResourcesQuery:  types.StringUnknown(),
		CommunicationWorkspace: types.StringNull(),
		CommunicationChannel:   types.StringUnknown(),
		Cells:                  types.StringValue("[]"),
		CellsFull:              types.StringValue("[]"),
		ParamsFull:             types.StringNull(),
		ExternalParamsFull:     types.StringUnknown(),
		AllowedEntities:        types.ListNull(types.StringType),
		Approvers:              types.ListNull(types.StringType),
		Labels:                 types.ListNull(types.StringType),
		Editors:                types.ListNull(types.StringType),
		SecretNames:            types.ListNull(types.StringType),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	runbookSchema := schema.RunbookSchema{}
	translationData := &coretranslator.TranslationData{CompatibilityOptions: runbookSchema.GetCompatibilityOptions()}

	// when
	stmt, err := translator.buildCreateStatement(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	assert.Contains(t, stmt, "notebook_name=\"null_test\"")
	assert.Contains(t, stmt, "enabled=true")
	assert.Contains(t, stmt, "timeout_ms=1000")
	assert.Contains(t, stmt, "description=\"\"") // Null becomes empty string
}
