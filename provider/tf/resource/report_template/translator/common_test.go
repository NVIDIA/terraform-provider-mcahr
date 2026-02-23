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
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	reporttemplatetf "terraform/terraform-provider/provider/tf/resource/report_template/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportTemplateTranslatorCommon_ToAPIModel_Create(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorCommon{}
	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name:       types.StringValue("test_report_template"),
		Blocks:     types.StringValue(`[{"title":"Block Name","resource_query":"host"}]`),
		BlocksFull: types.StringValue(`[{"title":"Block Name","resource_query":"host"}]`),
		Links:      types.StringValue(`[{"label":"test","report_template_name":"other_template"}]`),
		LinksFull:  types.StringValue(`[{"label":"test","report_template_name":"other_template"}]`),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	expectedStatement := `define_report_template(report_template_name="test_report_template", blocks="W3sidGl0bGUiOiJCbG9jayBOYW1lIiwicmVzb3VyY2VfcXVlcnkiOiJob3N0In1d", links="W3sibGFiZWwiOiJ0ZXN0IiwicmVwb3J0X3RlbXBsYXRlX25hbWUiOiJvdGhlcl90ZW1wbGF0ZSJ9XQ==")`
	assert.Equal(t, expectedStatement, result.Statement)
}

func TestReportTemplateTranslatorCommon_ToAPIModel_Update(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorCommon{}
	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name:       types.StringValue("test_report_template"),
		Blocks:     types.StringValue(`[{"title":"Updated Block","resource_query":"host"}]`),
		BlocksFull: types.StringValue(`[{"title":"Updated Block","resource_query":"host"}]`),
		Links:      types.StringValue(`[]`),
		LinksFull:  types.StringValue(`[]`),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Update).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	// Empty array "[]" encodes to "W10="
	expectedStatement := `update_report_template(report_template_name="test_report_template", blocks="W3sidGl0bGUiOiJVcGRhdGVkIEJsb2NrIiwicmVzb3VyY2VfcXVlcnkiOiJob3N0In1d", links="W10=")`
	assert.Equal(t, expectedStatement, result.Statement)
}

func TestReportTemplateTranslatorCommon_ToAPIModel_Read(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorCommon{}
	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name: types.StringValue("test_report_template"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Read).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	expectedStatement := `get_report_template_class(report_template_name="test_report_template")`
	assert.Equal(t, expectedStatement, result.Statement)
}

func TestReportTemplateTranslatorCommon_ToAPIModel_Delete(t *testing.T) {
	t.Parallel()
	// given
	translator := &ReportTemplateTranslatorCommon{}
	tfModel := &reporttemplatetf.ReportTemplateTFModel{
		Name: types.StringValue("test_report_template"),
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Delete).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)

	expectedStatement := `delete_report_template(report_template_name="test_report_template")`
	assert.Equal(t, expectedStatement, result.Statement)
}
