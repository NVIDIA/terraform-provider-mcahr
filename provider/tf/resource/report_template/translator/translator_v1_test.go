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
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportTemplateTranslatorV1_ToTFModel_Success(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorV1{}
	apiModel := createFullReportTemplateResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test_report_template", result.Name.ValueString())
	expectedBlocksJSON := `[{` +
		`"breakdown_by_tag":"tag_4",` +
		`"breakdown_tags_values":[{"color":"#AAAAAA","label":"label_0","values":["passed","skipped","failed"]}],` +
		`"group_by_tag":"tag_4",` +
		`"group_by_tag_order":{"type":"DEFAULT","values":[]},` +
		`"include_other_breakdown_tag_values":true,` +
		`"include_resources_without_group_tag":false,` +
		`"other_tags_to_export":["other_tag_1","other_tag_2"],` +
		`"resource_query":"host",` +
		`"resources_breakdown":[{"group_by_value":"tag_0","breakdown_values":[{"value":"value","count":1}]}],` +
		`"title":"Block Name",` +
		`"view_mode":"PERCENTAGE"` +
		`}]`
	assert.Equal(t, expectedBlocksJSON, result.Blocks.ValueString())

	expectedLinksJSON := `[{` +
		`"label":"minimal-report-full",` +
		`"report_template_name":"minimal_report_template"` +
		`}]`
	assert.Equal(t, expectedLinksJSON, result.Links.ValueString())
}

func TestReportTemplateTranslatorV1_ToTFModel_MinimalData(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorV1{}
	apiModel := createMinimalReportTemplateResponseV1()
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal_report_template", result.Name.ValueString())
	expectedBlocksJSON := `[{` +
		`"breakdown_by_tag":"",` +
		`"breakdown_tags_values":[],` +
		`"group_by_tag":"",` +
		`"group_by_tag_order":{"type":"DEFAULT","values":[]},` +
		`"include_other_breakdown_tag_values":false,` +
		`"include_resources_without_group_tag":false,` +
		`"other_tags_to_export":[],` +
		`"resource_query":"host",` +
		`"resources_breakdown":[],` +
		`"title":"Minimal Block",` +
		`"view_mode":"COUNT"` +
		`}]`
	assert.Equal(t, expectedBlocksJSON, result.Blocks.ValueString())
	assert.Equal(t, "[]", result.Links.ValueString())
}

func TestReportTemplateTranslatorV1_ToTFModel_NilModel(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestReportTemplateTranslatorV1_ToTFModel_NoContainer(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorV1{}
	apiModel := &reporttemplateapi.ReportTemplateResponseAPIModelV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no report template container found")
}

func TestReportTemplateTranslatorV1_ToTFModel_NoClasses(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorV1{}
	apiModel := &reporttemplateapi.ReportTemplateResponseAPIModelV1{
		GetReportTemplateClass: &reporttemplateapi.ReportTemplateContainerV1{
			ReportTemplateClasses: []reporttemplateapi.ReportTemplateClassV1{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no report template classes found")
}

// Helper functions for creating test data

func createFullReportTemplateResponseV1() *reporttemplateapi.ReportTemplateResponseAPIModelV1 {
	return &reporttemplateapi.ReportTemplateResponseAPIModelV1{
		GetReportTemplateClass: &reporttemplateapi.ReportTemplateContainerV1{
			Error: apicommon.ErrorV1{
				Type:    "OK",
				Message: "",
			},
			ReportTemplateClasses: []reporttemplateapi.ReportTemplateClassV1{
				{
					Name: "test_report_template",
					Blocks: `[{` +
						`"breakdown_by_tag":"tag_4",` +
						`"breakdown_tags_values":[{"color":"#AAAAAA","label":"label_0","values":["passed","skipped","failed"]}],` +
						`"group_by_tag":"tag_4",` +
						`"group_by_tag_order":{"type":"DEFAULT","values":[]},` +
						`"include_other_breakdown_tag_values":true,` +
						`"include_resources_without_group_tag":false,` +
						`"other_tags_to_export":["other_tag_1","other_tag_2"],` +
						`"resource_query":"host",` +
						`"resources_breakdown":[{"group_by_value":"tag_0","breakdown_values":[{"value":"value","count":1}]}],` +
						`"title":"Block Name",` +
						`"view_mode":"PERCENTAGE"` +
						`}]`,
					Links: `[{"label":"minimal-report-full","report_template_name":"minimal_report_template"}]`,
				},
			},
		},
	}
}

func createMinimalReportTemplateResponseV1() *reporttemplateapi.ReportTemplateResponseAPIModelV1 {
	return &reporttemplateapi.ReportTemplateResponseAPIModelV1{
		DefineReportTemplate: &reporttemplateapi.ReportTemplateContainerV1{
			Error: apicommon.ErrorV1{
				Type:    "OK",
				Message: "",
			},
			ReportTemplateClasses: []reporttemplateapi.ReportTemplateClassV1{
				{
					Name:   "minimal_report_template",
					Blocks: `[{"title":"Minimal Block","resource_query":"host"}]`,
					Links:  `[]`,
				},
			},
		},
	}
}
