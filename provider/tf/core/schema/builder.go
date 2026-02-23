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

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"terraform/terraform-provider/provider/common/version"
)

type SchemaBuilder struct {
	// The object that is being built
	Schema schema.Schema

	// Building arguments
	PlatformVersion *version.BackendVersion
}

func NewSchemaBuilder() *SchemaBuilder {

	return &SchemaBuilder{
		Schema: schema.Schema{
			Attributes: make(map[string]schema.Attribute),
		},
	}
}

func (b *SchemaBuilder) Build() schema.Schema {
	return b.Schema
}

func (b *SchemaBuilder) AddMarkdownDescription(description string) {
	b.Schema.MarkdownDescription = description
}

func (b *SchemaBuilder) AddAttribute(name string, attribute schema.Attribute) {
	b.Schema.Attributes[name] = attribute
}
