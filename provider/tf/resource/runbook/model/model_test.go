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
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunbookTFModel_DeepCopy_CreatesNewPointer(t *testing.T) {
	// given
	original := &RunbookTFModel{
		Name:        types.StringValue("test-runbook"),
		Description: types.StringValue("Test description"),
		Enabled:     types.BoolValue(true),
		TimeoutMs:   types.Int64Value(5000),
	}

	// when
	copy := original.Copy()

	// then
	require.NotNil(t, copy)
	assert.NotSame(t, original, copy, "Copy should be a different pointer")
}

func TestRunbookTFModel_DeepCopy_WithNullFields(t *testing.T) {
	// given
	original := &RunbookTFModel{
		Name:                   types.StringValue("test-runbook"),
		Description:            types.StringNull(),
		Enabled:                types.BoolNull(),
		TimeoutMs:              types.Int64Null(),
		Cells:                  types.StringNull(),
		Params:                 types.StringNull(),
		AllowedEntities:        types.ListNull(types.StringType),
		Labels:                 types.ListNull(types.StringType),
		CommunicationChannel:   types.StringNull(),
		FilterResourceToAction: types.BoolNull(),
	}

	// when
	copy := original.Copy()

	// then
	require.NotNil(t, copy)

	// Null values are preserved
	assert.True(t, copy.Description.IsNull())
	assert.True(t, copy.Enabled.IsNull())
	assert.True(t, copy.TimeoutMs.IsNull())
	assert.True(t, copy.Cells.IsNull())
	assert.True(t, copy.Params.IsNull())
	assert.True(t, copy.AllowedEntities.IsNull())
	assert.True(t, copy.Labels.IsNull())
	assert.True(t, copy.CommunicationChannel.IsNull())
	assert.True(t, copy.FilterResourceToAction.IsNull())
}

func TestRunbookTFModel_DeepCopy_WithUnknownFields(t *testing.T) {
	// given
	original := &RunbookTFModel{
		Name:                   types.StringValue("test-runbook"),
		Description:            types.StringUnknown(),
		Enabled:                types.BoolUnknown(),
		TimeoutMs:              types.Int64Unknown(),
		Cells:                  types.StringUnknown(),
		AllowedResourcesQuery:  types.StringUnknown(),
		AllowedEntities:        types.ListUnknown(types.StringType),
		FilterResourceToAction: types.BoolUnknown(),
	}

	// when
	copy := original.Copy()

	// then
	require.NotNil(t, copy)

	// Unknown values are preserved
	assert.True(t, copy.Description.IsUnknown())
	assert.True(t, copy.Enabled.IsUnknown())
	assert.True(t, copy.TimeoutMs.IsUnknown())
	assert.True(t, copy.Cells.IsUnknown())
	assert.True(t, copy.AllowedResourcesQuery.IsUnknown())
	assert.True(t, copy.AllowedEntities.IsUnknown())
	assert.True(t, copy.FilterResourceToAction.IsUnknown())
}

func TestRunbookTFModel_DeepCopy_EmptyModel(t *testing.T) {
	// given
	original := &RunbookTFModel{}

	// when
	copy := original.Copy()

	// then
	require.NotNil(t, copy)
	assert.NotSame(t, original, copy)

	// Verify all fields are null (zero values for Terraform types)
	// String fields
	assert.True(t, copy.Name.IsNull())
	assert.True(t, copy.Description.IsNull())
	assert.True(t, copy.Cells.IsNull())
	assert.True(t, copy.CellsFull.IsNull())
	assert.True(t, copy.Params.IsNull())
	assert.True(t, copy.ParamsFull.IsNull())
	assert.True(t, copy.ExternalParams.IsNull())
	assert.True(t, copy.ExternalParamsFull.IsNull())
	assert.True(t, copy.AllowedResourcesQuery.IsNull())
	assert.True(t, copy.CommunicationWorkspace.IsNull())
	assert.True(t, copy.CommunicationChannel.IsNull())
	assert.True(t, copy.Data.IsNull())

	// Boolean fields
	assert.True(t, copy.Enabled.IsNull())
	assert.True(t, copy.CommunicationCudNotifications.IsNull())
	assert.True(t, copy.CommunicationApprovalNotifications.IsNull())
	assert.True(t, copy.CommunicationExecutionNotifications.IsNull())
	assert.True(t, copy.IsRunOutputPersisted.IsNull())
	assert.True(t, copy.FilterResourceToAction.IsNull())

	// Numeric field
	assert.True(t, copy.TimeoutMs.IsNull())

	// Set fields
	assert.True(t, copy.AllowedEntities.IsNull())
	assert.True(t, copy.Approvers.IsNull())
	assert.True(t, copy.Labels.IsNull())
	assert.True(t, copy.Editors.IsNull())
	assert.True(t, copy.SecretNames.IsNull())
}

func TestRunbookTFModel_DeepCopy_NilModel(t *testing.T) {
	// given
	var original *RunbookTFModel = nil

	// when
	copy := original.Copy()

	// then
	assert.Nil(t, copy)
}

func TestRunbookTFModel_DeepCopy_ComplexStringFields(t *testing.T) {
	// given - test with complex JSON strings
	original := &RunbookTFModel{
		Name:           types.StringValue("complex-runbook"),
		Cells:          types.StringValue(`[{"type":"op","params":{"key":"value"}},{"type":"metric"}]`),
		CellsFull:      types.StringValue(`[{"type":"op","params":{"key":"value"},"full":true}]`),
		Params:         types.StringValue(`{"timeout":5000,"retries":3,"config":{"nested":"value"}}`),
		ParamsFull:     types.StringValue(`{"timeout":5000,"retries":3,"full":true}`),
		ExternalParams: types.StringValue(`{"api_key":"secret","endpoint":"https://api.example.com"}`),
		Data:           types.StringValue(`{"metadata":{"created":"2024-01-01","author":"test"}}`),
	}

	// when
	copy := original.Copy()

	// then
	// Complex strings are correctly copied
	assert.Equal(t, original.Cells.ValueString(), copy.Cells.ValueString())
	assert.Equal(t, original.CellsFull.ValueString(), copy.CellsFull.ValueString())
	assert.Equal(t, original.Params.ValueString(), copy.Params.ValueString())
	assert.Equal(t, original.ParamsFull.ValueString(), copy.ParamsFull.ValueString())
	assert.Equal(t, original.ExternalParams.ValueString(), copy.ExternalParams.ValueString())
	assert.Equal(t, original.Data.ValueString(), copy.Data.ValueString())
}

func TestRunbookTFModel_DeepCopy_EmptySetFields(t *testing.T) {
	// given - test with empty sets
	ctx := context.Background()
	emptyLabels, _ := types.ListValue(types.StringType, []attr.Value{})

	original := &RunbookTFModel{
		Name:   types.StringValue("empty-sets-runbook"),
		Labels: emptyLabels,
	}

	// when
	copy := original.Copy()

	// then
	require.NotNil(t, copy)

	// Verify empty set is preserved
	assert.True(t, original.Labels.Equal(copy.Labels))
	var copyLabelsElements []types.String
	copy.Labels.ElementsAs(ctx, &copyLabelsElements, false)
	assert.Equal(t, 0, len(copyLabelsElements))
}

func TestRunbookTFModel_DeepCopy_MixedFieldStates(t *testing.T) {
	// given - model with mix of null, unknown, and known values
	labels, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("label1"),
	})

	original := &RunbookTFModel{
		Name:                   types.StringValue("mixed-runbook"),
		Description:            types.StringNull(),
		Enabled:                types.BoolValue(true),
		TimeoutMs:              types.Int64Unknown(),
		Cells:                  types.StringValue(`[{"type":"op"}]`),
		Params:                 types.StringNull(),
		AllowedResourcesQuery:  types.StringValue("host"),
		CommunicationWorkspace: types.StringUnknown(),
		Labels:                 labels,
		Editors:                types.ListNull(types.StringType),
	}

	// when
	copy := original.Copy()

	// then
	// All states are preserved
	assert.Equal(t, "mixed-runbook", copy.Name.ValueString())
	assert.True(t, copy.Description.IsNull())
	assert.True(t, copy.Enabled.ValueBool())
	assert.True(t, copy.TimeoutMs.IsUnknown())
	assert.Equal(t, `[{"type":"op"}]`, copy.Cells.ValueString())
	assert.True(t, copy.Params.IsNull())
	assert.Equal(t, "host", copy.AllowedResourcesQuery.ValueString())
	assert.True(t, copy.CommunicationWorkspace.IsUnknown())
	assert.True(t, original.Labels.Equal(copy.Labels))
	assert.True(t, copy.Editors.IsNull())
}

func TestRunbookTFModel_DeepCopy_ZeroValues(t *testing.T) {
	// given - test with explicit zero values (different from null)
	original := &RunbookTFModel{
		Name:      types.StringValue("zero-values"),
		Enabled:   types.BoolValue(false), // explicit false
		TimeoutMs: types.Int64Value(0),    // explicit 0
	}

	// when
	copy := original.Copy()

	// then
	// Zero values are distinct from null
	assert.False(t, copy.Enabled.IsNull())
	assert.False(t, copy.Enabled.ValueBool())
	assert.False(t, copy.TimeoutMs.IsNull())
	assert.Equal(t, int64(0), copy.TimeoutMs.ValueInt64())
}

func TestRunbookTFModel_DeepCopy_OnlyRequiredFields(t *testing.T) {
	// given
	original := &RunbookTFModel{
		Name: types.StringValue("minimal-runbook"),
	}

	// when
	copy := original.Copy()

	// then
	assert.Equal(t, "minimal-runbook", copy.Name.ValueString())

	// Optional fields should remain as zero values
	assert.True(t, copy.Description.IsNull())
	assert.True(t, copy.Enabled.IsNull())
	assert.True(t, copy.Cells.IsNull())
}

func TestRunbookTFModel_DeepCopy_SpecialCharactersInStrings(t *testing.T) {
	// given - test with special characters
	original := &RunbookTFModel{
		Name:        types.StringValue("test-runbook-with-special-chars-123_456"),
		Description: types.StringValue("Description with\nnewlines\tand\ttabs"),
		Cells:       types.StringValue(`[{"command":"echo \"hello world\""}]`),
		Data:        types.StringValue(`{"unicode":"こんにちは","emoji":"🚀"}`),
	}

	// when
	copy := original.Copy()

	// then
	require.NotNil(t, copy)

	// Special characters are preserved
	assert.Equal(t, "Description with\nnewlines\tand\ttabs", copy.Description.ValueString())
	assert.Equal(t, `[{"command":"echo \"hello world\""}]`, copy.Cells.ValueString())
	assert.Equal(t, `{"unicode":"こんにちは","emoji":"🚀"}`, copy.Data.ValueString())
}

func TestRunbookTFModel_Copy_AllFieldsPopulated(t *testing.T) {

	// given
	allowedEntities, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("entity1"),
		types.StringValue("entity2"),
		types.StringValue("entity3"),
	})
	approvers, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("approver1"),
		types.StringValue("approver2"),
	})
	labels, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("label1"),
		types.StringValue("label2"),
		types.StringValue("label3"),
		types.StringValue("label4"),
	})
	editors, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("editor1"),
		types.StringValue("editor2"),
		types.StringValue("editor3"),
	})
	secretNames, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("secret1"),
		types.StringValue("secret2"),
	})

	original := &RunbookTFModel{
		// Required field
		Name: types.StringValue("all-fields-runbook"),

		// String fields
		Description:            types.StringValue("Complete description"),
		Cells:                  types.StringValue(`[{"type":"op","id":"1"}]`),
		CellsFull:              types.StringValue(`[{"type":"op","id":"1","full":true}]`),
		Params:                 types.StringValue(`{"key":"value"}`),
		ParamsFull:             types.StringValue(`{"key":"value","full":true}`),
		ExternalParams:         types.StringValue(`{"api":"key"}`),
		ExternalParamsFull:     types.StringValue(`{"api":"key","full":true}`),
		AllowedResourcesQuery:  types.StringValue("host"),
		CommunicationWorkspace: types.StringValue("workspace-123"),
		CommunicationChannel:   types.StringValue("#channel"),
		Data:                   types.StringValue(`{"metadata":"test"}`),

		// Boolean fields
		Enabled:                             types.BoolValue(true),
		CommunicationCudNotifications:       types.BoolValue(true),
		CommunicationApprovalNotifications:  types.BoolValue(false),
		CommunicationExecutionNotifications: types.BoolValue(true),
		IsRunOutputPersisted:                types.BoolValue(true),
		FilterResourceToAction:              types.BoolValue(false),

		// Numeric field
		TimeoutMs: types.Int64Value(30000),

		// Set fields (these contain internal slices - the key test!)
		AllowedEntities: allowedEntities,
		Approvers:       approvers,
		Labels:          labels,
		Editors:         editors,
		SecretNames:     secretNames,
	}

	// when - create a copy
	copy := original.Copy()

	// then - verify copy is not nil and is a different pointer
	require.NotNil(t, copy)
	assert.NotSame(t, original, copy, "Copy should be a different pointer")

	// Verify all fields are equal initially
	assert.Equal(t, original.Name.ValueString(), copy.Name.ValueString())
	assert.Equal(t, original.Description.ValueString(), copy.Description.ValueString())
	assert.Equal(t, original.Cells.ValueString(), copy.Cells.ValueString())
	assert.Equal(t, original.CellsFull.ValueString(), copy.CellsFull.ValueString())
	assert.Equal(t, original.Params.ValueString(), copy.Params.ValueString())
	assert.Equal(t, original.ParamsFull.ValueString(), copy.ParamsFull.ValueString())
	assert.Equal(t, original.ExternalParams.ValueString(), copy.ExternalParams.ValueString())
	assert.Equal(t, original.ExternalParamsFull.ValueString(), copy.ExternalParamsFull.ValueString())
	assert.Equal(t, original.AllowedResourcesQuery.ValueString(), copy.AllowedResourcesQuery.ValueString())
	assert.Equal(t, original.CommunicationWorkspace.ValueString(), copy.CommunicationWorkspace.ValueString())
	assert.Equal(t, original.CommunicationChannel.ValueString(), copy.CommunicationChannel.ValueString())
	assert.Equal(t, original.Data.ValueString(), copy.Data.ValueString())
	assert.Equal(t, original.Enabled.ValueBool(), copy.Enabled.ValueBool())
	assert.Equal(t, original.CommunicationCudNotifications.ValueBool(), copy.CommunicationCudNotifications.ValueBool())
	assert.Equal(t, original.CommunicationApprovalNotifications.ValueBool(), copy.CommunicationApprovalNotifications.ValueBool())
	assert.Equal(t, original.CommunicationExecutionNotifications.ValueBool(), copy.CommunicationExecutionNotifications.ValueBool())
	assert.Equal(t, original.IsRunOutputPersisted.ValueBool(), copy.IsRunOutputPersisted.ValueBool())
	assert.Equal(t, original.FilterResourceToAction.ValueBool(), copy.FilterResourceToAction.ValueBool())
	assert.Equal(t, original.TimeoutMs.ValueInt64(), copy.TimeoutMs.ValueInt64())

	// Verify Set fields are equal
	assert.True(t, original.AllowedEntities.Equal(copy.AllowedEntities))
	assert.True(t, original.Approvers.Equal(copy.Approvers))
	assert.True(t, original.Labels.Equal(copy.Labels))
	assert.True(t, original.Editors.Equal(copy.Editors))
	assert.True(t, original.SecretNames.Equal(copy.SecretNames))
}
