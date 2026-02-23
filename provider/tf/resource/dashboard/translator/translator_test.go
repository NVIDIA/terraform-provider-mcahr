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

func TestDashboardTranslator_ToTFModel(t *testing.T) {
	// Given
	translator := &DashboardTranslator{}
	apiModel := &dashboardapi.DashboardResponseAPIModel{
		Output: dashboardapi.DashboardOutput{
			Configurations: dashboardapi.DashboardConfigurations{
				Items: []dashboardapi.ConfigurationItem{
					{
						Config: dashboardapi.DashboardConfig{
							DashboardType: "test-type",
							Configuration: customattribute.ConfigurationJson{
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
							},
						},
						EntityMetadata: dashboardapi.DashboardEntityMetadata{
							Name: "test-dashboard",
						},
					},
				},
			},
		},
		Summary: dashboardapi.DashboardSummary{
			Status: "success",
			Errors: []apicommon.Error{},
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.NoError(t, err)
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

func TestDashboardTranslator_ToTFModel_EmptyOptionalFields(t *testing.T) {
	// Given
	translator := &DashboardTranslator{}
	apiModel := &dashboardapi.DashboardResponseAPIModel{
		Output: dashboardapi.DashboardOutput{
			Configurations: dashboardapi.DashboardConfigurations{
				Items: []dashboardapi.ConfigurationItem{
					{
						Config: dashboardapi.DashboardConfig{
							DashboardType: "minimal-type",
							Configuration: customattribute.ConfigurationJson{
								ResourceQuery: "",
								Groups:        []customattribute.GroupJson{},
								Values:        []customattribute.ValueJson{},
								OtherTags:     []string{},
								Identifiers:   []string{},
							},
						},
						EntityMetadata: dashboardapi.DashboardEntityMetadata{
							Name: "minimal-dashboard",
						},
					},
				},
			},
		},
		Summary: dashboardapi.DashboardSummary{
			Status: "success",
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.NoError(t, err)
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

func TestDashboardTranslator_ToTFModel_Nil(t *testing.T) {
	// Given
	dashboardTranslator := &DashboardTranslator{}
	requestContext := common.NewRequestContext(context.Background())
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := dashboardTranslator.ToTFModel(requestContext, translationData, nil)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestDashboardTranslator_ToTFModel_NoConfigurations(t *testing.T) {
	// Given
	dashboardTranslator := &DashboardTranslator{}
	apiModel := &dashboardapi.DashboardResponseAPIModel{
		Output: dashboardapi.DashboardOutput{
			Configurations: dashboardapi.DashboardConfigurations{
				Items: []dashboardapi.ConfigurationItem{},
			},
		},
	}

	requestContext := common.NewRequestContext(context.Background())
	translationData := &coretranslator.TranslationData{}

	// When
	result, err := dashboardTranslator.ToTFModel(requestContext, translationData, apiModel)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no configurations found")
}
