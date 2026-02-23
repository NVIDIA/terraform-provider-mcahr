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

package compatibility

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/common/version"
	core "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test TFModel implementation for testing
type TestTFModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Timeout     types.Int64  `tfsdk:"timeout"`
	Labels      types.List   `tfsdk:"labels"`
}

var _ core.TFModel = &TestTFModel{}

func (t *TestTFModel) GetName() string {
	return t.Name.ValueString()
}

func TestPostProcess_NullifyIncompatibleAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tfModel      *TestTFModel
		options      map[string]attribute.CompatibilityOptions
		expectNulled map[string]bool // field name -> should be nulled
	}{
		{
			name: "No options - all compatible",
			tfModel: &TestTFModel{
				Name:        types.StringValue("test"),
				Description: types.StringValue("test description"),
				Enabled:     types.BoolValue(true),
				Timeout:     types.Int64Value(30),
			},
			options:      map[string]attribute.CompatibilityOptions{},
			expectNulled: map[string]bool{},
		},
		{
			name: "Some attributes incompatible",
			tfModel: &TestTFModel{
				Name:        types.StringValue("test"),
				Description: types.StringValue("test description"),
				Enabled:     types.BoolValue(true),
				Timeout:     types.Int64Value(30),
			},
			options: map[string]attribute.CompatibilityOptions{
				"description": {MinVersion: "release-2.0.0"}, // Will be incompatible
				"timeout":     {MinVersion: "release-1.0.0"}, // Will be compatible
			},
			expectNulled: map[string]bool{
				"description": true,
				"timeout":     false,
			},
		},
		{
			name: "All configured attributes incompatible",
			tfModel: &TestTFModel{
				Name:        types.StringValue("test"),
				Description: types.StringValue("test description"),
				Enabled:     types.BoolValue(true),
			},
			options: map[string]attribute.CompatibilityOptions{
				"description": {MinVersion: "release-2.0.0"},
				"enabled":     {MinVersion: "release-3.0.0"},
			},
			expectNulled: map[string]bool{
				"description": true,
				"enabled":     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			backendVersion := version.NewBackendVersion("release-1.5.0")
			require.NotNil(t, backendVersion, "Failed to create backend version")

			requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)

			// when
			err := PostProcess(requestContext, tt.tfModel, tt.options)

			// then
			require.NoError(t, err)

			// Check that expected fields are nulled
			for fieldName, shouldBeNull := range tt.expectNulled {
				switch fieldName {
				case "description":
					if shouldBeNull {
						assert.True(t, tt.tfModel.Description.IsNull(), "Description should be null")
					} else {
						assert.False(t, tt.tfModel.Description.IsNull(), "Description should not be null")
					}
				case "enabled":
					if shouldBeNull {
						assert.True(t, tt.tfModel.Enabled.IsNull(), "Enabled should be null")
					} else {
						assert.False(t, tt.tfModel.Enabled.IsNull(), "Enabled should not be null")
					}
				case "timeout":
					if shouldBeNull {
						assert.True(t, tt.tfModel.Timeout.IsNull(), "Timeout should be null")
					} else {
						assert.False(t, tt.tfModel.Timeout.IsNull(), "Timeout should not be null")
					}
				}
			}

			// Name should never be nulled (not in options)
			assert.False(t, tt.tfModel.Name.IsNull(), "Name should never be null")
		})
	}
}

func TestSetFieldToNil_AllTypes(t *testing.T) {
	t.Parallel()

	// Test model with different field types
	testModel := &TestTFModel{
		Name:        types.StringValue("test"),
		Description: types.StringValue("test description"),
		Enabled:     types.BoolValue(true),
		Timeout:     types.Int64Value(30),
		Labels:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("label1")}),
	}

	// Create compatibility options that will make all fields incompatible
	options := map[string]attribute.CompatibilityOptions{
		"name":        {MinVersion: "release-2.0.0"},
		"description": {MinVersion: "release-2.0.0"},
		"enabled":     {MinVersion: "release-2.0.0"},
		"timeout":     {MinVersion: "release-2.0.0"},
		"labels":      {MinVersion: "release-2.0.0"},
	}

	backendVersion := version.NewBackendVersion("release-1.0.0") // Lower version
	require.NotNil(t, backendVersion, "Failed to create backend version")

	requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)

	// when - generic function will infer the type from the concrete model
	err := PostProcess(requestContext, testModel, options)

	// then
	require.NoError(t, err)

	// All fields should be null
	assert.True(t, testModel.Name.IsNull(), "Name should be null")
	assert.True(t, testModel.Description.IsNull(), "Description should be null")
	assert.True(t, testModel.Enabled.IsNull(), "Enabled should be null")
	assert.True(t, testModel.Timeout.IsNull(), "Timeout should be null")
	assert.True(t, testModel.Labels.IsNull(), "Labels should be null")
}
