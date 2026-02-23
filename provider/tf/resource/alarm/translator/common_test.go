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
	alarmtf "terraform/terraform-provider/provider/tf/resource/alarm/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getAlarmCommonParameters returns the common parameter string for define/update operations
func getAlarmCommonParameters() string {
	return `alarm_name="test-alarm", ` +
		`fire_query="cpu_usage > 80", ` +
		`clear_query="cpu_usage < 70", ` +
		`description="Test alarm description", ` +
		`resource_query="host", ` +
		`resource_type="HOST", ` +
		`check_interval_sec=60, ` +
		`family="custom", ` +
		`enabled=true, ` +
		`fire_title_template="fired test alarm", ` +
		`fire_short_template="fired test short", ` +
		`resolve_title_template="cleared test alarm", ` +
		`resolve_short_template="cleared test short"`
}

func TestAlarmTranslatorCommon_ToAPIModel(t *testing.T) {
	// Given
	translator := &AlarmTranslatorCommon{}

	// Create test TF model
	tfModel := &alarmtf.AlarmTFModel{
		Name:                 types.StringValue("test-alarm"),
		FireQuery:            types.StringValue("cpu_usage > 80"),
		ClearQuery:           types.StringValue("cpu_usage < 70"),
		Description:          types.StringValue("Test alarm description"),
		ResourceQuery:        types.StringValue("host"),
		ResourceType:         types.StringValue("HOST"),
		CheckIntervalSec:     types.Int64Value(60),
		Family:               types.StringValue("custom"),
		Enabled:              types.BoolValue(true),
		FireTitleTemplate:    types.StringValue("fired test alarm"),
		FireShortTemplate:    types.StringValue("fired test short"),
		ResolveTitleTemplate: types.StringValue("cleared test alarm"),
		ResolveShortTemplate: types.StringValue("cleared test short"),
	}

	tests := []struct {
		name      string
		operation common.CrudOperation
		expected  string
	}{
		{
			name:      "Create operation",
			operation: common.Create,
			expected:  fmt.Sprintf("define_alarm(%s)", getAlarmCommonParameters()),
		},
		{
			name:      "Read operation",
			operation: common.Read,
			expected:  `get_alarm_class(alarm_name="test-alarm")`,
		},
		{
			name:      "Update operation",
			operation: common.Update,
			expected:  fmt.Sprintf("update_alarm(%s)", getAlarmCommonParameters()),
		},
		{
			name:      "Delete operation",
			operation: common.Delete,
			expected:  `delete_alarm(alarm_name="test-alarm")`,
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

func TestAlarmTranslatorCommon_ToAPIModel_UnsupportedOperation(t *testing.T) {
	// Given
	translator := &AlarmTranslatorCommon{}
	tfModel := &alarmtf.AlarmTFModel{
		Name: types.StringValue("test-alarm"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.CrudOperation(999)).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// When
	// Test with an invalid operation (cast to avoid compile error)
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// Then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}
