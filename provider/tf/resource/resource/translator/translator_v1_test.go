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
	resourceapi "terraform/terraform-provider/provider/external_api/resources/resources"
	coretranslator "terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceTranslatorV1_ToTFModel_Success(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorV1{}
	apiModel := &resourceapi.ResourceResponseAPIModelV1{
		DefineResource: &resourceapi.ResourceContainerV1Single{
			Symbol: resourceapi.ResourceSymbolV1{
				Name:    "test_resource",
				Formula: "host | pod | app='test'",
				Attributes: resourceapi.ResourceAttributesV1{
					Description: "Test resource",
					Params:      `["param1", "param2"]`,
				},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test_resource", result.Name.ValueString())
	assert.Equal(t, "Test resource", result.Description.ValueString())
	assert.Equal(t, "host | pod | app='test'", result.Value.ValueString())
	require.False(t, result.Params.IsNull())
	params := utils.ListSliceFromTFModel(requestContext.Context, result.Params)
	assert.Equal(t, []string{"param1", "param2"}, params)
}

func TestResourceTranslatorV1_ToTFModel_NilModel(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestResourceTranslatorV1_ToTFModel_NoContainer(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorV1{}
	apiModel := &resourceapi.ResourceResponseAPIModelV1{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no resource container found")
}

func TestResourceTranslatorV1_ToTFModel_NoSymbols(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslatorV1{}
	apiModel := &resourceapi.ResourceResponseAPIModelV1{
		ListType: &resourceapi.ResourceContainerV1List{
			Symbol: []resourceapi.ResourceSymbolV1{},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Read).WithAPIVersion(common.V1)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no resource symbols found")
}
