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

package process

import (
	"context"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	dashboardtf "terraform/terraform-provider/provider/tf/resource/dashboard/model"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestPostProcessJsonFullFields_WithValidJSON(t *testing.T) {
	t.Parallel()

	tfModel := &dashboardtf.DashboardTFModel{
		Groups: types.StringValue(`[
			{
				"name": "g1",
				"tags": ["cloud_provider", "release_tag"]
			}
		]`),
		GroupsFull: types.StringValue(`[
			{
				"name": "g1",
				"tags": ["cloud_provider", "release_tag"]
			}
		]`),
		Values: types.StringValue(`[
			{
				"color": "#78909c",
				"values": ["aws"]
			}
		]`),
		ValuesFull: types.StringValue(`[
			{
				"color": "#78909c",
				"values": ["aws"]
			}
		]`),
	}

	backendVersion := &version.BackendVersion{
		Major: 2,
		Minor: 0,
		Patch: 0,
	}

	requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)

	err := postProcessJsonFullFields(requestContext, tfModel)
	assert.NoError(t, err)

	// Verify that GroupsFull and ValuesFull are populated
	assert.False(t, tfModel.GroupsFull.IsNull())
	assert.False(t, tfModel.ValuesFull.IsNull())

	// The exact JSON content may vary due to processing, but should be valid JSON
	assert.NotEqual(t, "", tfModel.GroupsFull.ValueString())
	assert.NotEqual(t, "", tfModel.ValuesFull.ValueString())
}

func TestPostProcessJsonFullFields_WithEmptyJSON(t *testing.T) {
	t.Parallel()

	tfModel := &dashboardtf.DashboardTFModel{
		Groups:     types.StringValue("[]"),
		GroupsFull: types.StringValue("[]"),
		Values:     types.StringValue("[]"),
		ValuesFull: types.StringValue("[]"),
	}

	backendVersion := &version.BackendVersion{
		Major: 2,
		Minor: 0,
		Patch: 0,
	}

	requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)

	err := postProcessJsonFullFields(requestContext, tfModel)
	assert.NoError(t, err)

	// Should still populate _full fields even with empty arrays
	assert.Equal(t, "[]", tfModel.GroupsFull.ValueString())
	assert.Equal(t, "[]", tfModel.ValuesFull.ValueString())
}

func TestPostProcessJsonFullFields_WithNullValues(t *testing.T) {
	t.Parallel()

	tfModel := &dashboardtf.DashboardTFModel{
		Groups:     types.StringNull(),
		GroupsFull: types.StringNull(),
		Values:     types.StringNull(),
		ValuesFull: types.StringNull(),
	}

	backendVersion := &version.BackendVersion{
		Major: 2,
		Minor: 0,
		Patch: 0,
	}

	requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)

	err := postProcessJsonFullFields(requestContext, tfModel)
	assert.NoError(t, err)

	// With null input values, the _full values should remain null
	assert.True(t, tfModel.GroupsFull.IsNull())
	assert.True(t, tfModel.ValuesFull.IsNull())
}

func TestPostProcessJsonFullFields_WithInvalidJSON(t *testing.T) {
	t.Parallel()

	tfModel := &dashboardtf.DashboardTFModel{
		Groups:     types.StringValue(`[{"invalid": json}]`),
		GroupsFull: types.StringValue(`[{"invalid": json}]`),
		Values:     types.StringValue("[]"),
		ValuesFull: types.StringValue("[]"),
	}

	backendVersion := &version.BackendVersion{
		Major: 2,
		Minor: 0,
		Patch: 0,
	}

	requestContext := common.NewRequestContext(context.Background()).WithBackendVersion(backendVersion)

	err := postProcessJsonFullFields(requestContext, tfModel)
	assert.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "invalid character")
	}
}
