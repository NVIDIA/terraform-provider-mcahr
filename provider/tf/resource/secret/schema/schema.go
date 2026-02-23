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
	"terraform/terraform-provider/provider/common/attribute"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	schemabuilder "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type NVaultSecretSchema struct{}

var _ coreschema.ResourceSchema = &NVaultSecretSchema{}

// GetNVaultSecretResourceSchema returns the schema for the nvault secret resource
func (s *NVaultSecretSchema) GetSchema() schema.Schema {
	builder := schemabuilder.NewSchemaBuilder()

	builder.AddMarkdownDescription("Shoreline nvault_secret. A secret managed by NVault. Creating it requires an active NVault integration to be configured and enabled.")

	builder.AddAttribute("name", schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The name/symbol for the object within Shoreline and the op language (must be unique, only alphanumeric/underscore).",
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		Validators:          []validator.String{validators.NameValidator()},
	})

	builder.AddAttribute("vault_secret_path", schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The full path in Vault where the secret is stored. It includes the mount path and the subpath. It tells Vault where to look for the secret.",
	})

	builder.AddAttribute("vault_secret_key", schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The specific key within the secret data to retrieve.",
	})

	builder.AddAttribute("integration_name", schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The name/symbol of a Shoreline integration.",
	})

	return builder.Build()
}

func (s *NVaultSecretSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{}
}
