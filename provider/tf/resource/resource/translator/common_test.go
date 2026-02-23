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
	"terraform/terraform-provider/provider/common"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	resourcetf "terraform/terraform-provider/provider/tf/resource/resource/model"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceTranslatorCommon_ToAPIModel_Create(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorCommon{}
	tfModel := &resourcetf.ResourceTFModel{
		Name:        types.StringValue("test_resource"),
		Description: types.StringValue("Test resource"),
		Value:       types.StringValue("host | pod | app='test'"),
		Params:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("param1"), types.StringValue("param2")}),
	}

	// when
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, common.V2, result.APIVersion)
	assert.Equal(t, `define_resource(key="test_resource", val="host | pod | app='test'", description="Test resource", params=["param1", "param2"])`, result.Statement)
}

func TestResourceTranslatorCommon_ToAPIModel_Read(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorCommon{}
	tfModel := &resourcetf.ResourceTFModel{
		Name: types.StringValue("test_resource"),
	}

	// when
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Read).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, common.V2, result.APIVersion)
	assert.Equal(t, `list resources | name = "test_resource"`, result.Statement)
}

func TestResourceTranslatorCommon_ToAPIModel_Update(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorCommon{}
	tfModel := &resourcetf.ResourceTFModel{
		Name:        types.StringValue("test_resource"),
		Description: types.StringValue("Updated resource"),
		Value:       types.StringValue("host | pod | app='updated'"),
		Params:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("param1"), types.StringValue("param2"), types.StringValue("param3")}),
	}

	// when
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Update).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, common.V2, result.APIVersion)
	assert.Equal(t, `update_resource(resource_name="test_resource", value="host | pod | app='updated'", description="Updated resource", params=["param1", "param2", "param3"])`, result.Statement)
}

func TestResourceTranslatorCommon_ToAPIModel_Delete(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorCommon{}
	tfModel := &resourcetf.ResourceTFModel{
		Name: types.StringValue("test_resource"),
	}

	// when
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Delete).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, common.V2, result.APIVersion)
	assert.Equal(t, `delete test_resource`, result.Statement)
}

func TestResourceTranslatorCommon_ToAPIModel_UnsupportedOperation(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorCommon{}
	tfModel := &resourcetf.ResourceTFModel{
		Name: types.StringValue("test_resource"),
	}

	// when
	requestContext := common.NewRequestContext(context.Background()).WithOperation(999).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}
	result, err := translator.ToAPIModelWithVersion(requestContext, translationData, tfModel)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operation")
}
