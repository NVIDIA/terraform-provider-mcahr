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
	reporttemplateapi "terraform/terraform-provider/provider/external_api/resources/report_templates"
	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportTemplateTranslator_ToTFModel_Success(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslator{}
	apiModel := createFullReportTemplateResponseV2()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_report_template", result.Name.ValueString())
	// Note: Fields are alphabetically ordered due to map iteration in ApplyCustomStructTags
	expectedBlocksJSON := `[{` +
		`"breakdown_by_tag":"service",` +
		`"breakdown_tags_values":[{"color":"#FF0000","label":"Production","values":["prod","production"]}],` +
		`"group_by_tag":"environment",` +
		`"group_by_tag_order":{"type":"CUSTOM","values":["prod","staging"]},` +
		`"include_other_breakdown_tag_values":false,` +
		`"include_resources_without_group_tag":true,` +
		`"other_tags_to_export":["tag1","tag2"],` +
		`"resource_query":"host",` +
		`"resources_breakdown":[{"group_by_value":"environment","breakdown_values":[{"value":"prod","count":1}]}],` +
		`"title":"Block Name",` +
		`"view_mode":"COUNT"` +
		`}]`
	assert.Equal(t, expectedBlocksJSON, result.Blocks.ValueString())

	expectedLinksJSON := `[{` +
		`"label":"test",` +
		`"report_template_name":"other_template"` +
		`}]`
	assert.Equal(t, expectedLinksJSON, result.Links.ValueString())
}

func TestReportTemplateTranslator_ToTFModel_MinimalData(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslator{}
	apiModel := createMinimalReportTemplateResponseV2()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal_report_template", result.Name.ValueString())
	// Note: All fields included with defaults due to custom marshaling
	expectedMinimalBlocksJSON := `[{` +
		`"breakdown_by_tag":"",` +
		`"breakdown_tags_values":null,` +
		`"group_by_tag":"",` +
		`"group_by_tag_order":{"type":"","values":null},` +
		`"include_other_breakdown_tag_values":false,` +
		`"include_resources_without_group_tag":false,` +
		`"other_tags_to_export":null,` +
		`"resource_query":"host",` +
		`"resources_breakdown":null,` +
		`"title":"Minimal Block",` +
		`"view_mode":""` +
		`}]`
	assert.Equal(t, expectedMinimalBlocksJSON, result.Blocks.ValueString())
	assert.Equal(t, "[]", result.Links.ValueString())
}

func TestReportTemplateTranslator_ToTFModel_NilModel(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslator{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestReportTemplateTranslator_ToTFModel_NoConfigurations(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslator{}
	apiModel := &reporttemplateapi.ReportTemplateResponseAPIModel{
		Output: reporttemplateapi.ReportTemplateOutput{
			Configurations: reporttemplateapi.ReportTemplateConfigurations{
				Items: []reporttemplateapi.ReportTemplateConfigurationItem{},
			},
		},
		Summary: reporttemplateapi.ReportTemplateSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no report template configurations found")
}

func TestReportTemplateTranslator_GetErrors_WithErrors(t *testing.T) {
	t.Parallel()
	// given
	apiModel := &reporttemplateapi.ReportTemplateResponseAPIModel{
		Output: reporttemplateapi.ReportTemplateOutput{
			Configurations: reporttemplateapi.ReportTemplateConfigurations{
				Items: []reporttemplateapi.ReportTemplateConfigurationItem{},
			},
		},
		Summary: reporttemplateapi.ReportTemplateSummary{
			Status: "OP_FAILED",
			Errors: []apicommon.Error{
				{
					Message: "Class not found for name \"full_report_template\"",
					Type:    "EXECUTION_ERROR",
				},
			},
		},
	}

	// when
	result := apiModel.GetErrors()

	// then
	expectedError := "Status: OP_FAILED; Errors: EXECUTION_ERROR: Class not found for name \"full_report_template\""
	assert.Equal(t, expectedError, result)
}

func TestReportTemplateTranslator_GetErrors_NoErrors(t *testing.T) {
	t.Parallel()
	// given
	apiModel := &reporttemplateapi.ReportTemplateResponseAPIModel{
		Summary: reporttemplateapi.ReportTemplateSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}

	// when
	result := apiModel.GetErrors()

	// then
	assert.Equal(t, "", result)
}

// Helper functions for creating test data

func createFullReportTemplateResponseV2() *reporttemplateapi.ReportTemplateResponseAPIModel {
	return &reporttemplateapi.ReportTemplateResponseAPIModel{
		Output: reporttemplateapi.ReportTemplateOutput{
			Configurations: reporttemplateapi.ReportTemplateConfigurations{
				Items: []reporttemplateapi.ReportTemplateConfigurationItem{
					{
						Config: reporttemplateapi.ReportTemplateData{
							Name: "test_report_template",
							Blocks: []customattribute.BlockJson{
								{
									Title:                           "Block Name",
									ResourceQuery:                   "host",
									GroupByTag:                      "environment",
									BreakdownByTag:                  "service",
									ViewMode:                        "COUNT",
									IncludeResourcesWithoutGroupTag: true,
									IncludeOtherBreakdownTagValues:  false,
									OtherTagsToExport:               []string{"tag1", "tag2"},
									BreakdownTagsValues: []customattribute.BreakdownTagValue{
										{
											Label:  "Production",
											Values: []string{"prod", "production"},
											Color:  "#FF0000",
										},
									},
									ResourcesBreakdown: []customattribute.ResourcesBreakdown{
										{
											BreakdownValues: []customattribute.BreakdownValue{
												{
													Count: 1,
													Value: "prod",
												},
											},
											GroupByValue: "environment",
										},
									},
									GroupByTagOrder: customattribute.GroupByTagOrder{
										Type:   "CUSTOM",
										Values: []string{"prod", "staging"},
									},
								},
							},
							Links: []customattribute.LinkJson{
								{
									Label:              "test",
									ReportTemplateName: "other_template",
								},
							},
						},
						EntityMetadata: reporttemplateapi.ReportTemplateEntityMetadata{
							Name: "test_report_template",
						},
					},
				},
			},
		},
		Summary: reporttemplateapi.ReportTemplateSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}

func createMinimalReportTemplateResponseV2() *reporttemplateapi.ReportTemplateResponseAPIModel {
	return &reporttemplateapi.ReportTemplateResponseAPIModel{
		Output: reporttemplateapi.ReportTemplateOutput{
			Configurations: reporttemplateapi.ReportTemplateConfigurations{
				Items: []reporttemplateapi.ReportTemplateConfigurationItem{
					{
						Config: reporttemplateapi.ReportTemplateData{
							Name: "minimal_report_template",
							Blocks: []customattribute.BlockJson{
								{
									Title:         "Minimal Block",
									ResourceQuery: "host",
								},
							},
							Links: []customattribute.LinkJson{},
						},
						EntityMetadata: reporttemplateapi.ReportTemplateEntityMetadata{
							Name: "minimal_report_template",
						},
					},
				},
			},
		},
		Summary: reporttemplateapi.ReportTemplateSummary{
			Status: "OP_COMPLETED",
			Errors: []apicommon.Error{},
		},
	}
}
