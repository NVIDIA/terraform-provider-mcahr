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

func TestResourceTranslator_ToTFModel_Success(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslator{}
	apiModel := &resourceapi.ResourceResponseAPIModel{
		Output: resourceapi.ResourceOutput{
			Symbols: resourceapi.ResourceConfigurations{
				Items: []resourceapi.ResourceConfigurationItem{
					{
						Name:        "test_resource",
						Description: "Test resource",
						Formula:     "host | pod | app='test'",
						Attributes: resourceapi.ResourceAttributes{
							Params: `["param1", "param2"]`,
						},
					},
				},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
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

func TestResourceTranslator_ToTFModel_NilModel(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslator{}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, nil)

	// then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestResourceTranslator_ToTFModel_NoConfigurations(t *testing.T) {
	t.Parallel()
	// given
	translator := &ResourceTranslator{}
	apiModel := &resourceapi.ResourceResponseAPIModel{
		Output: resourceapi.ResourceOutput{
			Symbols: resourceapi.ResourceConfigurations{
				Items: []resourceapi.ResourceConfigurationItem{},
			},
		},
	}
	requestContext := common.NewRequestContext(context.Background()).WithOperation(common.Create).WithAPIVersion(common.V2)
	translationData := &coretranslator.TranslationData{}

	// when
	result, err := translator.ToTFModel(requestContext, translationData, apiModel)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no resource configurations found")
}
