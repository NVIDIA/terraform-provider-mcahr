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

package mapbuilder

import (
	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/common/version"
)

type MapBuilder struct {
	attrCompatibilityChecker *attribute.CompatibilityChecker

	fields map[string]any
}

func NewMapBuilder(backendVersion *version.BackendVersion, compatibilityOptions map[string]attribute.CompatibilityOptions) *MapBuilder {

	attrCompatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, compatibilityOptions)

	return &MapBuilder{
		attrCompatibilityChecker: attrCompatibilityChecker,
		fields:                   make(map[string]any),
	}
}

func (b *MapBuilder) SetField(mapFieldName string, tfFieldName string, value any) *MapBuilder {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		b.fields[mapFieldName] = value
	}
	return b
}

func (b *MapBuilder) Build() map[string]any {
	return b.fields
}
