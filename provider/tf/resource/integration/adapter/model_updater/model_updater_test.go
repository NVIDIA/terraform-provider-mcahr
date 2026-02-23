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

package modelupdater

import (
	"testing"

	"terraform/terraform-provider/provider/common/attribute"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	"terraform/terraform-provider/provider/tf/resource/integration/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"terraform/terraform-provider/provider/common/version"
)

func TestNewModelUpdater(t *testing.T) {
	tests := []struct {
		name    string
		options *adapterinterface.IntegrationDataAdapterOptions
		tfModel *model.IntegrationTFModel
	}{
		{
			name: "with valid options and model",
			options: &adapterinterface.IntegrationDataAdapterOptions{
				BackendVersion: version.NewBackendVersion("release-1.2.3"),
				CompatibilityOptions: map[string]attribute.CompatibilityOptions{
					"field1": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
				},
			},
			tfModel: &model.IntegrationTFModel{},
		},
		{
			name: "with nil backend version",
			options: &adapterinterface.IntegrationDataAdapterOptions{
				BackendVersion:       nil,
				CompatibilityOptions: map[string]attribute.CompatibilityOptions{},
			},
			tfModel: &model.IntegrationTFModel{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updater := NewModelUpdater(tt.options, tt.tfModel)

			assert.NotNil(t, updater)
			assert.NotNil(t, updater.attrCompatibilityChecker)
			assert.NotNil(t, updater.tfModel)
			assert.NotNil(t, updater.fields)
			assert.Equal(t, 0, len(updater.fields))
		})
	}
}

func TestModelUpdater_UpdateStringField_Compatible(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"api_url": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		APIUrl: types.StringValue("old_value"),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.StringValue("new_value")

	result := updater.UpdateStringField("api_url", &tfModel.APIUrl, newValue)

	// Should return updater for chaining
	assert.Equal(t, updater, result)

	// Field should be updated
	assert.Equal(t, "new_value", tfModel.APIUrl.ValueString())
}

func TestModelUpdater_UpdateStringField_Incompatible_BelowMinVersion(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.0.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"api_url": {MinVersion: "release-1.5.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		APIUrl: types.StringValue("old_value"),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.StringValue("new_value")

	result := updater.UpdateStringField("api_url", &tfModel.APIUrl, newValue)

	// Should return updater for chaining
	assert.Equal(t, updater, result)

	// Field should NOT be updated
	assert.Equal(t, "old_value", tfModel.APIUrl.ValueString())
}

func TestModelUpdater_UpdateStringField_Incompatible_AboveMaxVersion(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-2.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"api_url": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		APIUrl: types.StringValue("old_value"),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.StringValue("new_value")

	result := updater.UpdateStringField("api_url", &tfModel.APIUrl, newValue)

	// Should return updater for chaining
	assert.Equal(t, updater, result)

	// Field should NOT be updated
	assert.Equal(t, "old_value", tfModel.APIUrl.ValueString())
}

func TestModelUpdater_UpdateStringField_NoCompatibilityOptions(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{},
	}

	tfModel := &model.IntegrationTFModel{
		APIUrl: types.StringValue("old_value"),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.StringValue("new_value")

	updater.UpdateStringField("api_url", &tfModel.APIUrl, newValue)

	// Field should be updated (no restrictions)
	assert.Equal(t, "new_value", tfModel.APIUrl.ValueString())
}

func TestModelUpdater_UpdateInt64Field_Compatible(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"cache_ttl": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		CacheTTL: types.Int64Value(100),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.Int64Value(200)

	result := updater.UpdateInt64Field("cache_ttl", &tfModel.CacheTTL, newValue)

	// Should return updater for chaining
	assert.Equal(t, updater, result)

	// Field should be updated
	assert.Equal(t, int64(200), tfModel.CacheTTL.ValueInt64())
}

func TestModelUpdater_UpdateInt64Field_Incompatible(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.0.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"cache_ttl": {MinVersion: "release-1.5.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		CacheTTL: types.Int64Value(100),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.Int64Value(200)

	updater.UpdateInt64Field("cache_ttl", &tfModel.CacheTTL, newValue)

	// Field should NOT be updated
	assert.Equal(t, int64(100), tfModel.CacheTTL.ValueInt64())
}

func TestModelUpdater_UpdateBoolField_Compatible(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"enabled": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		Enabled: types.BoolValue(false),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.BoolValue(true)

	result := updater.UpdateBoolField("enabled", &tfModel.Enabled, newValue)

	// Should return updater for chaining
	assert.Equal(t, updater, result)

	// Field should be updated
	assert.Equal(t, true, tfModel.Enabled.ValueBool())
}

func TestModelUpdater_UpdateBoolField_Incompatible(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-3.0.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"enabled": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		Enabled: types.BoolValue(false),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.BoolValue(true)

	updater.UpdateBoolField("enabled", &tfModel.Enabled, newValue)

	// Field should NOT be updated
	assert.Equal(t, false, tfModel.Enabled.ValueBool())
}

func TestModelUpdater_UpdateSetField_Compatible(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"payload_paths": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		},
	}

	oldSet, _ := types.ListValueFrom(nil, types.StringType, []string{"old1", "old2"})
	tfModel := &model.IntegrationTFModel{
		PayloadPaths: oldSet,
	}

	newSet, _ := types.ListValueFrom(nil, types.StringType, []string{"new1", "new2", "new3"})

	updater := NewModelUpdater(options, tfModel)
	result := updater.UpdateSetField("payload_paths", &tfModel.PayloadPaths, newSet)

	// Should return updater for chaining
	assert.Equal(t, updater, result)

	// Field should be updated
	var resultValues []string
	tfModel.PayloadPaths.ElementsAs(nil, &resultValues, false)
	assert.ElementsMatch(t, []string{"new1", "new2", "new3"}, resultValues)
}

func TestModelUpdater_UpdateSetField_Incompatible(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.0.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"payload_paths": {MinVersion: "release-1.5.0", MaxVersion: "release-2.0.0"},
		},
	}

	oldSet, _ := types.ListValueFrom(nil, types.StringType, []string{"old1", "old2"})
	tfModel := &model.IntegrationTFModel{
		PayloadPaths: oldSet,
	}

	newSet, _ := types.ListValueFrom(nil, types.StringType, []string{"new1", "new2", "new3"})

	updater := NewModelUpdater(options, tfModel)
	updater.UpdateSetField("payload_paths", &tfModel.PayloadPaths, newSet)

	// Field should NOT be updated
	var resultValues []string
	tfModel.PayloadPaths.ElementsAs(nil, &resultValues, false)
	assert.ElementsMatch(t, []string{"old1", "old2"}, resultValues)
}

func TestModelUpdater_MethodChaining(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"api_url":   {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
			"cache_ttl": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
			"enabled":   {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		APIUrl:   types.StringValue("old_url"),
		CacheTTL: types.Int64Value(100),
		Enabled:  types.BoolValue(false),
	}

	updater := NewModelUpdater(options, tfModel)

	// Test method chaining
	updater.
		UpdateStringField("api_url", &tfModel.APIUrl, types.StringValue("new_url")).
		UpdateInt64Field("cache_ttl", &tfModel.CacheTTL, types.Int64Value(200)).
		UpdateBoolField("enabled", &tfModel.Enabled, types.BoolValue(true))

	assert.Equal(t, "new_url", tfModel.APIUrl.ValueString())
	assert.Equal(t, int64(200), tfModel.CacheTTL.ValueInt64())
	assert.Equal(t, true, tfModel.Enabled.ValueBool())
}

func TestModelUpdater_MultipleFields_MixedCompatibility(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"api_url":   {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"}, // Compatible
			"cache_ttl": {MinVersion: "release-2.5.0", MaxVersion: "release-3.0.0"}, // Incompatible (too old)
			"enabled":   {MinVersion: "release-1.0.0", MaxVersion: "release-1.2.0"}, // Incompatible (too new)
		},
	}

	tfModel := &model.IntegrationTFModel{
		APIUrl:   types.StringValue("old_url"),
		CacheTTL: types.Int64Value(100),
		Enabled:  types.BoolValue(false),
	}

	updater := NewModelUpdater(options, tfModel)
	updater.
		UpdateStringField("api_url", &tfModel.APIUrl, types.StringValue("new_url")).
		UpdateInt64Field("cache_ttl", &tfModel.CacheTTL, types.Int64Value(200)).
		UpdateBoolField("enabled", &tfModel.Enabled, types.BoolValue(true))

	// Only compatible field should be updated
	assert.Equal(t, "new_url", tfModel.APIUrl.ValueString())
	assert.Equal(t, int64(100), tfModel.CacheTTL.ValueInt64()) // Not updated
	assert.Equal(t, false, tfModel.Enabled.ValueBool())        // Not updated
}

func TestModelUpdater_UpdateStringField_NilBackendVersion(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion: nil,
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{
			"api_url": {MinVersion: "release-1.0.0", MaxVersion: "release-2.0.0"},
		},
	}

	tfModel := &model.IntegrationTFModel{
		APIUrl: types.StringValue("old_value"),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.StringValue("new_value")

	updater.UpdateStringField("api_url", &tfModel.APIUrl, newValue)

	// Field should be updated (nil version assumes latest)
	assert.Equal(t, "new_value", tfModel.APIUrl.ValueString())
}

func TestModelUpdater_UpdateStringField_NullValue(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{},
	}

	tfModel := &model.IntegrationTFModel{
		APIUrl: types.StringValue("old_value"),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.StringNull()

	updater.UpdateStringField("api_url", &tfModel.APIUrl, newValue)

	// Field should be updated to null
	assert.True(t, tfModel.APIUrl.IsNull())
}

func TestModelUpdater_UpdateInt64Field_ZeroValue(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{},
	}

	tfModel := &model.IntegrationTFModel{
		CacheTTL: types.Int64Value(100),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.Int64Value(0)

	updater.UpdateInt64Field("cache_ttl", &tfModel.CacheTTL, newValue)

	// Field should be updated to zero
	assert.Equal(t, int64(0), tfModel.CacheTTL.ValueInt64())
}

func TestModelUpdater_UpdateBoolField_SameValue(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{},
	}

	tfModel := &model.IntegrationTFModel{
		Enabled: types.BoolValue(true),
	}

	updater := NewModelUpdater(options, tfModel)
	newValue := types.BoolValue(true)

	updater.UpdateBoolField("enabled", &tfModel.Enabled, newValue)

	// Field should still be true
	assert.Equal(t, true, tfModel.Enabled.ValueBool())
}

func TestModelUpdater_UpdateSetField_EmptySet(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{},
	}

	oldSet, _ := types.ListValueFrom(nil, types.StringType, []string{"old1", "old2"})
	tfModel := &model.IntegrationTFModel{
		PayloadPaths: oldSet,
	}

	emptySet, _ := types.ListValueFrom(nil, types.StringType, []string{})

	updater := NewModelUpdater(options, tfModel)
	updater.UpdateSetField("payload_paths", &tfModel.PayloadPaths, emptySet)

	// Field should be updated to empty set
	var resultValues []string
	tfModel.PayloadPaths.ElementsAs(nil, &resultValues, false)
	assert.Empty(t, resultValues)
}

func TestModelUpdater_UpdateMultipleFieldTypes(t *testing.T) {
	options := &adapterinterface.IntegrationDataAdapterOptions{
		BackendVersion:       version.NewBackendVersion("release-1.5.0"),
		CompatibilityOptions: map[string]attribute.CompatibilityOptions{},
	}

	tfModel := &model.IntegrationTFModel{
		Name:         types.StringValue("old_name"),
		APIUrl:       types.StringValue("old_url"),
		CacheTTL:     types.Int64Value(100),
		CacheTTLMs:   types.Int64Value(100000),
		Enabled:      types.BoolValue(false),
		PayloadPaths: types.ListNull(types.StringType),
	}

	newSet, _ := types.ListValueFrom(nil, types.StringType, []string{"path1", "path2"})

	updater := NewModelUpdater(options, tfModel)
	updater.
		UpdateStringField("name", &tfModel.Name, types.StringValue("new_name")).
		UpdateStringField("api_url", &tfModel.APIUrl, types.StringValue("new_url")).
		UpdateInt64Field("cache_ttl", &tfModel.CacheTTL, types.Int64Value(200)).
		UpdateInt64Field("cache_ttl_ms", &tfModel.CacheTTLMs, types.Int64Value(200000)).
		UpdateBoolField("enabled", &tfModel.Enabled, types.BoolValue(true)).
		UpdateSetField("payload_paths", &tfModel.PayloadPaths, newSet)

	assert.Equal(t, "new_name", tfModel.Name.ValueString())
	assert.Equal(t, "new_url", tfModel.APIUrl.ValueString())
	assert.Equal(t, int64(200), tfModel.CacheTTL.ValueInt64())
	assert.Equal(t, int64(200000), tfModel.CacheTTLMs.ValueInt64())
	assert.Equal(t, true, tfModel.Enabled.ValueBool())

	var resultValues []string
	tfModel.PayloadPaths.ElementsAs(nil, &resultValues, false)
	assert.ElementsMatch(t, []string{"path1", "path2"}, resultValues)
}
