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

package modelupdater

import (
	"terraform/terraform-provider/provider/common/attribute"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"
	"terraform/terraform-provider/provider/tf/resource/integration/model"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ModelUpdater struct {
	attrCompatibilityChecker *attribute.CompatibilityChecker
	tfModel                  *model.IntegrationTFModel

	fields map[string]any
}

func NewModelUpdater(options *adapterinterface.IntegrationDataAdapterOptions, tfModel *model.IntegrationTFModel) *ModelUpdater {

	attrCompatibilityChecker := attribute.NewCompatibilityChecker(options.BackendVersion, options.CompatibilityOptions)

	return &ModelUpdater{
		attrCompatibilityChecker: attrCompatibilityChecker,
		tfModel:                  tfModel,
		fields:                   make(map[string]any),
	}
}

func (b *ModelUpdater) UpdateStringField(tfFieldName string, modelValue *basetypes.StringValue, mapValue basetypes.StringValue) *ModelUpdater {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		*modelValue = mapValue
	}
	return b
}

func (b *ModelUpdater) UpdateInt64Field(tfFieldName string, modelValue *basetypes.Int64Value, mapValue basetypes.Int64Value) *ModelUpdater {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		*modelValue = mapValue
	}
	return b
}

func (b *ModelUpdater) UpdateBoolField(tfFieldName string, modelValue *basetypes.BoolValue, mapValue basetypes.BoolValue) *ModelUpdater {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		*modelValue = mapValue
	}
	return b
}

func (b *ModelUpdater) UpdateSetField(tfFieldName string, modelValue *basetypes.ListValue, mapValue basetypes.ListValue) *ModelUpdater {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		*modelValue = mapValue
	}
	return b
}
