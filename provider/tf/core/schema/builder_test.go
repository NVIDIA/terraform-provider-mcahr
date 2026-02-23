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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSchemaBuilder(t *testing.T) {
	t.Parallel()

	// when
	builder := NewSchemaBuilder()

	// then
	require.NotNil(t, builder, "expected non-nil builder")
	assert.NotNil(t, builder.Schema.Attributes, "expected attributes map to be initialized")
	assert.Empty(t, builder.Schema.Attributes, "expected empty attributes map")
}

func TestSchemaBuilder_Build(t *testing.T) {
	t.Parallel()

	// given
	builder := NewSchemaBuilder()
	builder.AddMarkdownDescription("test description")

	// when
	result := builder.Build()

	// then
	assert.Equal(t, "test description", result.MarkdownDescription, "description mismatch")
	assert.NotNil(t, result.Attributes, "expected attributes map in built schema")
}

func TestSchemaBuilder_AddMarkdownDescription(t *testing.T) {
	t.Parallel()

	// given
	builder := NewSchemaBuilder()
	description := "This is a test description"

	// when
	builder.AddMarkdownDescription(description)

	// then
	assert.Equal(t, description, builder.Schema.MarkdownDescription, "description mismatch")
}

func TestSchemaBuilder_AddAttribute(t *testing.T) {
	t.Parallel()

	// given
	builder := NewSchemaBuilder()
	attr := schema.StringAttribute{
		Required: true,
	}

	// when
	builder.AddAttribute("test_attr", attr)

	// then
	assert.Len(t, builder.Schema.Attributes, 1, "expected 1 attribute")
	assert.Contains(t, builder.Schema.Attributes, "test_attr", "expected 'test_attr' to be added to schema")
}
