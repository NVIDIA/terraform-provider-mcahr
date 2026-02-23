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
	apicommon "terraform/terraform-provider/provider/external_api/resources/common"
	dashboardapi "terraform/terraform-provider/provider/external_api/resources/dashboards"
	customattribute "terraform/terraform-provider/provider/external_api/resources/dashboards/custom_attribute"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFullDashboardResponseV1() *dashboardapi.DashboardResponseAPIModelV1 {
	config := customattribute.ConfigurationJson{
		ResourceQuery: "host",
		Groups: []customattribute.GroupJson{
			{
				Name: "group1",
				Tags: []string{"tag1", "tag2"},
			},
		},
		Values: []customattribute.ValueJson{
			{
				Color:  "red",
				Values: []string{"value1", "value2"},
			},
		},
		OtherTags:   []string{"other1", "other2"},
		Identifiers: []string{"id1", "id2"},
	}

	configJSON, _ := config.MarshalJSON()

	return &dashboardapi.DashboardResponseAPIModelV1{
		GetDashboardClass: &dashboardapi.DashboardContainerV1{
			Error: &apicommon.ErrorV1{
				Type: "OK",
			},
			DashboardClasses: []dashboardapi.DashboardClassV1{
				{
					Name:          "test-dashboard",
					DashboardType: "test-type",
					Configuration: string(configJSON),
				},
			},
		},
	}
}

func createMinimalDashboardResponseV1() *dashboardapi.DashboardResponseAPIModelV1 {
	config := customattribute.ConfigurationJson{
		ResourceQuery: "",
		Groups:        []customattribute.GroupJson{},
		Values:        []customattribute.ValueJson{},
		OtherTags:     []string{},
		Identifiers:   []string{},
	}

	configJSON, _ := config.MarshalJSON()

	return &dashboardapi.DashboardResponseAPIModelV1{
		GetDashboardClass: &dashboardapi.DashboardContainerV1{
			Error: &apicommon.ErrorV1{
				Type: "OK",
			},
			DashboardClasses: []dashboardapi.DashboardClassV1{
				{
					Name:          "minimal-dashboard",
					DashboardType: "minimal-type",
					Configuration: string(configJSON),
				},
			},
		},
	}
}

func TestDashboardTranslatorV1_ToTFModel_Success(t *testing.T) {
	t.Parallel()

	// Given
	translator := &DashboardTranslatorV1{}
	apiModel := createFullDashboardResponseV1()

	requestContext := common.NewRequestContext(context.Background())
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test-dashboard", result.Name.ValueString())
	assert.Equal(t, "test-type", result.DashboardType.ValueString())
	assert.Equal(t, "host", result.ResourceQuery.ValueString())

	// Verify Groups JSON
	expectedGroups := `[{"name":"group1","tags":["tag1","tag2"]}]`
	assert.JSONEq(t, expectedGroups, result.Groups.ValueString())
	assert.JSONEq(t, expectedGroups, result.GroupsFull.ValueString())

	// Verify Values JSON
	expectedValues := `[{"color":"red","values":["value1","value2"]}]`
	assert.JSONEq(t, expectedValues, result.Values.ValueString())
	assert.JSONEq(t, expectedValues, result.ValuesFull.ValueString())

	// Verify sets
	otherTags := []string{}
	result.OtherTags.ElementsAs(t.Context(), &otherTags, false)
	assert.ElementsMatch(t, []string{"other1", "other2"}, otherTags)

	identifiers := []string{}
	result.Identifiers.ElementsAs(t.Context(), &identifiers, false)
	assert.ElementsMatch(t, []string{"id1", "id2"}, identifiers)
}

func TestDashboardTranslatorV1_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &DashboardTranslatorV1{}
	apiModel := createMinimalDashboardResponseV1()

	requestContext := common.NewRequestContext(context.Background())
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal-dashboard", result.Name.ValueString())
	assert.Equal(t, "minimal-type", result.DashboardType.ValueString())
	assert.Equal(t, "", result.ResourceQuery.ValueString())
	assert.JSONEq(t, "[]", result.Groups.ValueString())
	assert.JSONEq(t, "[]", result.Values.ValueString())

	otherTags := []string{}
	result.OtherTags.ElementsAs(t.Context(), &otherTags, false)
	assert.Empty(t, otherTags)

	identifiers := []string{}
	result.Identifiers.ElementsAs(t.Context(), &identifiers, false)
	assert.Empty(t, identifiers)
}

func TestDashboardTranslatorV1_ToTFModel_EmptyDashboardClasses(t *testing.T) {
	t.Parallel()

	// Given - V1 API response with empty dashboard classes
	apiModel := &dashboardapi.DashboardResponseAPIModelV1{
		GetDashboardClass: &dashboardapi.DashboardContainerV1{
			Error: &apicommon.ErrorV1{
				Type: "OK",
			},
			DashboardClasses: []dashboardapi.DashboardClassV1{}, // Empty list
		},
	}

	dashboardTranslator := &DashboardTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background())
	translationData := &coretranslator.TranslationData{}

	// When
	tfModel, err := dashboardTranslator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	require.Error(t, err)
	require.Nil(t, tfModel)
	assert.Contains(t, err.Error(), "no dashboard classes found")
}

func TestDashboardTranslatorV1_ToTFModel_NilInput(t *testing.T) {
	t.Parallel()

	// Given
	dashboardTranslator := &DashboardTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background())
	translationData := &coretranslator.TranslationData{}

	// When
	tfModel, err := dashboardTranslator.ToTFModel(requestContext, translationData, nil)

	// Then
	require.NoError(t, err)
	require.Nil(t, tfModel)
}
